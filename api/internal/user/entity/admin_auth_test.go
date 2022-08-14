package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestAdminRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		role      int32
		expect    AdminRole
		expectErr error
	}{
		{
			name:      "success",
			role:      1,
			expect:    AdminRoleAdministrator,
			expectErr: nil,
		},
		{
			name:      "invalid role",
			role:      0,
			expect:    AdminRoleUnknown,
			expectErr: errInvalidAdminRole,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewAdminRole(tt.role)
			assert.ErrorIs(t, tt.expectErr, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAdminRole_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		role   AdminRole
		expect error
	}{
		{
			name:   "administrator",
			role:   AdminRoleAdministrator,
			expect: nil,
		},
		{
			name:   "coordinator",
			role:   AdminRoleCoordinator,
			expect: nil,
		},
		{
			name:   "producer",
			role:   AdminRoleProducer,
			expect: nil,
		},
		{
			name:   "unknown",
			role:   AdminRoleUnknown,
			expect: errInvalidAdminRole,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.role.Validate()
			assert.ErrorIs(t, err, tt.expect)
		})
	}
}

func TestAdminAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		adminID   string
		cognitoID string
		role      AdminRole
		expect    *AdminAuth
	}{
		{
			name:      "success",
			adminID:   "admin-id",
			cognitoID: "cognito-id",
			role:      AdminRoleAdministrator,
			expect: &AdminAuth{
				AdminID:   "admin-id",
				CognitoID: "cognito-id",
				Role:      AdminRoleAdministrator,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAdminAuth(tt.adminID, tt.cognitoID, tt.role)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAdminAuth_Fill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		auth   *AdminAuth
		result *cognito.AuthResult
		expect *AdminAuth
	}{
		{
			name: "success",
			auth: &AdminAuth{
				AdminID:   "admin-id",
				CognitoID: "cognito-id",
				Role:      AdminRoleAdministrator,
				CreatedAt: jst.Date(2022, 1, 1, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 18, 30, 0, 0),
			},
			result: &cognito.AuthResult{
				IDToken:      "id-token",
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
			},
			expect: &AdminAuth{
				AdminID:      "admin-id",
				CognitoID:    "cognito-id",
				Role:         AdminRoleAdministrator,
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
				CreatedAt:    jst.Date(2022, 1, 1, 18, 30, 0, 0),
				UpdatedAt:    jst.Date(2022, 1, 1, 18, 30, 0, 0),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.auth.Fill(tt.result)
			assert.Equal(t, tt.expect, tt.auth)
		})
	}
}

func TestAdminAuths_GroupByRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		auths  AdminAuths
		expect map[AdminRole]AdminAuths
	}{
		{
			name: "success",
			auths: AdminAuths{
				{
					AdminID:   "admin-id",
					CognitoID: "cognito-id",
					Role:      AdminRoleAdministrator,
				},
			},
			expect: map[AdminRole]AdminAuths{
				AdminRoleAdministrator: {
					{
						AdminID:   "admin-id",
						CognitoID: "cognito-id",
						Role:      AdminRoleAdministrator,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.auths.GroupByRole())
		})
	}
}

func TestAdminAuths_AdminIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		auths  AdminAuths
		expect []string
	}{
		{
			name: "success",
			auths: AdminAuths{
				{
					AdminID:   "admin-id",
					CognitoID: "cognito-id",
					Role:      AdminRoleAdministrator,
				},
			},
			expect: []string{"admin-id"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.auths.AdminIDs())
		})
	}
}

func TestAdminAuths_Devices(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		auths  AdminAuths
		expect []string
	}{
		{
			name: "success",
			auths: AdminAuths{
				{
					AdminID:   "admin-id",
					CognitoID: "cognito-id",
					Device:    "instance-id",
					Role:      AdminRoleAdministrator,
				},
			},
			expect: []string{"instance-id"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.auths.Devices())
		})
	}
}
