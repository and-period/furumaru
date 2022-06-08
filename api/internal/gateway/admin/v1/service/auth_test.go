package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestAdminRole(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		role           entity.AdminRole
		expect         AdminRole
		expectString   string
		expectResponse int32
	}{
		{
			name:           "administrator",
			role:           entity.AdminRoleAdministrator,
			expect:         AdminRoleAdministrator,
			expectString:   "administrator",
			expectResponse: 1,
		},
		{
			name:           "coordinator",
			role:           entity.AdminRoleCoordinator,
			expect:         AdminRoleCoordinator,
			expectString:   "coordinator",
			expectResponse: 2,
		},
		{
			name:           "producer",
			role:           entity.AdminRoleProducer,
			expect:         AdminRoleProducer,
			expectString:   "producer",
			expectResponse: 3,
		},
		{
			name:           "unknown",
			role:           entity.AdminRoleUnknown,
			expect:         AdminRoleUnknown,
			expectString:   "unknown",
			expectResponse: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAdminRole(tt.role)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.expectString, actual.String())
			assert.Equal(t, tt.expectResponse, actual.Response())
		})
	}
}

func TestAuth(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		auth   *entity.AdminAuth
		expect *Auth
	}{
		{
			name: "success",
			auth: &entity.AdminAuth{
				AdminID:      "admin-id",
				Role:         entity.AdminRoleAdministrator,
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
			},
			expect: &Auth{
				Auth: &response.Auth{
					AdminID:      "admin-id",
					Role:         int32(AdminRoleAdministrator),
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresIn:    3600,
					TokenType:    "Bearer",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewAuth(tt.auth))
		})
	}
}

func TestAuth_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		auth   *Auth
		expect *response.Auth
	}{
		{
			name: "success",
			auth: &Auth{
				Auth: &response.Auth{
					AdminID:      "admin-id",
					Role:         int32(AdminRoleAdministrator),
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresIn:    3600,
					TokenType:    "Bearer",
				},
			},
			expect: &response.Auth{
				AdminID:      "admin-id",
				Role:         int32(AdminRoleAdministrator),
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
				TokenType:    "Bearer",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.auth.Response())
		})
	}
}
