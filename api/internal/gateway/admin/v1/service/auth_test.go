package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestAuthProviderType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		providerType entity.AdminAuthProviderType
		expect       AuthProviderType
		response     int32
	}{
		{
			name:         "google",
			providerType: entity.AdminAuthProviderTypeGoogle,
			expect:       AuthProviderTypeGoogle,
			response:     1,
		},
		{
			name:         "line",
			providerType: entity.AdminAuthProviderTypeLINE,
			expect:       AuthProviderTypeLINE,
			response:     2,
		},
		{
			name:         "unknown",
			providerType: entity.AdminAuthProviderTypeUnknown,
			expect:       AuthProviderTypeUnknown,
			response:     0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAuthProviderType(tt.providerType)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.response, actual.Response())
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
				Type:         entity.AdminTypeAdministrator,
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
			},
			expect: &Auth{
				Auth: response.Auth{
					AdminID:      "admin-id",
					Type:         int32(AdminTypeAdministrator),
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
		expect *response.Auth
	}{
		{
			name: "success",
			auth: &Auth{
				Auth: response.Auth{
					AdminID:      "admin-id",
					Type:         int32(AdminTypeAdministrator),
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresIn:    3600,
					TokenType:    "Bearer",
				},
			},
			expect: &response.Auth{
				AdminID:      "admin-id",
				Type:         int32(AdminTypeAdministrator),
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
		response  []*response.AuthProvider
	}{
		{
			name: "success",
			providers: entity.AdminAuthProviders{
				{
					AdminID:      "admin-id",
					ProviderType: entity.AdminAuthProviderTypeGoogle,
					AccountID:    "account-id",
					Email:        "test@example.com",
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
			expect: AuthProviders{
				{
					AuthProvider: response.AuthProvider{
						Type:        1,
						ConnectedAt: now.Unix(),
					},
				},
			},
			response: []*response.AuthProvider{
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
