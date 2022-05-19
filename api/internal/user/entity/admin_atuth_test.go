package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/stretchr/testify/assert"
)

func TestAdminAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		adminID string
		role    AdminRole
		result  *cognito.AuthResult
		expect  *AdminAuth
	}{
		{
			name:    "success",
			adminID: "admin-id",
			role:    AdminRoleAdministrator,
			result: &cognito.AuthResult{
				IDToken:      "id-token",
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
			},
			expect: &AdminAuth{
				AdminID:      "admin-id",
				Role:         AdminRoleAdministrator,
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAdminAuth(tt.adminID, tt.role, tt.result)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
