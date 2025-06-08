//nolint:paralleltest,tparallel
package rbac

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnforcer(t *testing.T) {
	modelPath, err := generateTempFile("model.conf", modelText)
	require.NoError(t, err)
	defer os.Remove(modelPath)
	policyPath, err := generateTempFile("policy.csv", policyText)
	require.NoError(t, err)
	defer os.Remove(policyPath)

	tests := []struct {
		name       string
		modelPath  string
		policyPath string
		expect     bool
		expectErr  bool
	}{
		{
			name:       "success",
			modelPath:  modelPath,
			policyPath: policyPath,
			expect:     true,
			expectErr:  false,
		},
		{
			name:       "non initialize",
			modelPath:  "",
			policyPath: "",
			expect:     true,
			expectErr:  false,
		},
		{
			name:       "occurred error",
			modelPath:  "dummy.conf",
			policyPath: "dummy.csv",
			expect:     false,
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enforcer, err := NewEnforcer(tt.modelPath, tt.policyPath)
			require.Equal(t, tt.expectErr, err != nil, err)
			assert.Equal(t, tt.expect, enforcer != nil)
		})
	}
}

func TestEnforcer_Enforce(t *testing.T) {
	modelPath, err := generateTempFile("model.conf", modelText)
	require.NoError(t, err)
	defer os.Remove(modelPath)
	policyPath, err := generateTempFile("policy.csv", policyText)
	require.NoError(t, err)
	defer os.Remove(policyPath)
	enforcer, err := NewEnforcer(modelPath, policyPath)
	require.NoError(t, err)
	assert.NotNil(t, enforcer)

	tests := []struct {
		name      string
		group     string
		path      string
		method    string
		expect    bool
		expectErr bool
	}{
		{
			name:      "allow admin",
			group:     "admin",
			path:      "/v1/admins",
			method:    http.MethodPost,
			expect:    true,
			expectErr: false,
		},
		{
			name:      "deny admin",
			group:     "admin",
			path:      "/v1/admins/hoge/deny",
			method:    http.MethodGet,
			expect:    false,
			expectErr: false,
		},
		{
			name:      "allow developer",
			group:     "developer",
			path:      "/v1/admins",
			method:    http.MethodGet,
			expect:    true,
			expectErr: false,
		},
		{
			name:      "deny developer",
			group:     "developer",
			path:      "/v1/admins",
			method:    http.MethodPost,
			expect:    false,
			expectErr: false,
		},
		{
			name:      "allow operator",
			group:     "operator",
			path:      "/v1/users",
			method:    http.MethodGet,
			expect:    true,
			expectErr: false,
		},
		{
			name:      "deny operator",
			group:     "operator",
			path:      "/v1/admins",
			method:    http.MethodGet,
			expect:    false,
			expectErr: false,
		},
		{
			name:      "deny default",
			group:     "default",
			path:      "/v1/admins",
			method:    http.MethodGet,
			expect:    false,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			enforced, err := enforcer.Enforce(tt.group, tt.path, tt.method)
			assert.Equal(t, tt.expectErr, err != nil, err)
			assert.Equal(t, tt.expect, enforced)
		})
	}
}

func TestEnforcer_GetRolesForuser(t *testing.T) {
	modelPath, err := generateTempFile("model.conf", modelText)
	require.NoError(t, err)
	defer os.Remove(modelPath)
	policyPath, err := generateTempFile("policy.csv", policyText)
	require.NoError(t, err)
	defer os.Remove(policyPath)
	enforcer, err := NewEnforcer(modelPath, policyPath)
	require.NoError(t, err)
	assert.NotNil(t, enforcer)

	tests := []struct {
		name      string
		group     string
		domain    []string
		expect    []string
		expectErr bool
	}{
		{
			name:      "admin",
			group:     "admin",
			domain:    []string{},
			expect:    []string{"admin_write", "admin_read", "user_write", "user_read"},
			expectErr: false,
		},
		{
			name:      "developer",
			group:     "developer",
			domain:    []string{},
			expect:    []string{"admin_read", "user_read"},
			expectErr: false,
		},
		{
			name:      "operator",
			group:     "operator",
			domain:    []string{},
			expect:    []string{"user_read"},
			expectErr: false,
		},
		{
			name:      "default",
			group:     "default",
			domain:    []string{},
			expect:    []string{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rules, err := enforcer.GetRolesForUser(tt.group, tt.domain...)
			assert.Equal(t, tt.expectErr, err != nil, err)
			assert.ElementsMatch(t, tt.expect, rules)
		})
	}
}

func generateTempFile(name, content string) (string, error) {
	f, err := os.CreateTemp("", name)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b := []byte(content)
	if _, err := f.Write(b); err != nil {
		return "", err
	}
	return f.Name(), nil
}

const modelText = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
`

const policyText = `
p, admin_write, /v1/admins, POST, allow
p, admin_write, /v1/admins/*, PATCH, allow
p, admin_read, /v1/admins, GET, allow
p, admin_read, /v1/admins/*, GET, allow
p, admin_read, /v1/admins/*/deny, GET, deny
p, user_write, /v1/users, POST, allow
p, user_write, /v1/users/*, (PATCH)|(DELETE), allow
p, user_write, /v1/users/*, PUT, deny
p, user_read, /v1/users, GET, allow

g, admin, admin_write
g, admin, admin_read
g, admin, user_write
g, admin, user_read
g, developer, admin_read
g, developer, user_read
g, operator, user_read
`
