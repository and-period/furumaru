package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestAuthProviderType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		providerType entity.AdminAuthProviderType
		expect       types.AuthProviderType
	}{
		{
			name:         "google",
			providerType: entity.AdminAuthProviderType(types.AuthProviderTypeGoogle),
			expect:       types.AuthProviderType(types.AuthProviderTypeGoogle),
		},
		{
			name:         "line",
			providerType: entity.AdminAuthProviderType(types.AuthProviderTypeLINE),
			expect:       types.AuthProviderType(types.AuthProviderTypeLINE),
		},
		{
			name:         "unknown",
			providerType: entity.AdminAuthProviderType(types.AuthProviderTypeUnknown),
			expect:       types.AuthProviderType(types.AuthProviderTypeUnknown),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAuthProviderType(tt.providerType)
			assert.Equal(t, tt.expect, actual.Response())
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
				Type:         entity.AdminType(types.AdminTypeAdministrator),
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
			},
			expect: &Auth{
				Auth: types.Auth{
					AdminID:      "admin-id",
					Type:         types.AdminTypeAdministrator,
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresIn:    3600,
					TokenType:    "Bearer",
				},
			},
		},
	}
	for _, tt := range tests {
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
		expect *types.Auth
	}{
		{
			name: "success",
			auth: &Auth{
				Auth: types.Auth{
					AdminID:      "admin-id",
					Type:         types.AdminTypeAdministrator,
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresIn:    3600,
					TokenType:    "Bearer",
				},
			},
			expect: &types.Auth{
				AdminID:      "admin-id",
				Type:         types.AdminTypeAdministrator,
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
				TokenType:    "Bearer",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.auth.Response())
		})
	}
}

func TestAuthProviders(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name      string
		providers entity.AdminAuthProviders
		expect    AuthProviders
		response  []*types.AuthProvider
	}{
		{
			name: "success",
			providers: entity.AdminAuthProviders{
				{
					AdminID:      "admin-id",
					ProviderType: entity.AdminAuthProviderType(types.AuthProviderTypeGoogle),
					AccountID:    "account-id",
					Email:        "test@example.com",
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
			expect: AuthProviders{
				{
					AuthProvider: types.AuthProvider{
						Type:        1,
						ConnectedAt: now.Unix(),
					},
				},
			},
			response: []*types.AuthProvider{
				{
					Type:        1,
					ConnectedAt: now.Unix(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAuthProviders(tt.providers)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.response, actual.Response())
		})
	}
}
