package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			name:   "developer",
			role:   AdminRoleDeveloper,
			expect: nil,
		},
		{
			name:   "operator",
			role:   AdminRoleOperator,
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

func TestAdmin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		adminID   string
		cognitoID string
		email     string
		role      AdminRole
		expect    *Admin
	}{
		{
			name:      "success",
			adminID:   "admin-id",
			cognitoID: "cognito-id",
			email:     "test-admin@and-period.jp",
			role:      AdminRoleAdministrator,
			expect: &Admin{
				ID:        "admin-id",
				CognitoID: "cognito-id",
				Email:     "test-admin@and-period.jp",
				Role:      AdminRoleAdministrator,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAdmin(tt.adminID, tt.cognitoID, tt.email, tt.role)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
