package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
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

func TestAuthUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admin  *entity.Admin
		expect *AuthUser
	}{
		{
			name: "success",
			admin: &entity.Admin{
				ID:            "admin-id",
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				StoreName:     "&.農園",
				Email:         "test-admin01@and-period.jp",
				PhoneNumber:   "+819012345678",
				Role:          entity.AdminRoleAdministrator,
				ThumbnailURL:  "https://and-period.jp",
				CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &AuthUser{
				AuthUser: &response.AuthUser{
					ID:            "admin-id",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					StoreName:     "&.農園",
					ThumbnailURL:  "https://and-period.jp",
					Role:          int32(AdminRoleAdministrator),
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewAuthUser(tt.admin))
		})
	}
}

func TestAuthUser_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		user   *AuthUser
		expect *response.AuthUser
	}{
		{
			name: "success",
			user: &AuthUser{
				AuthUser: &response.AuthUser{
					ID:            "admin-id",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					StoreName:     "&.農園",
					ThumbnailURL:  "https://and-period.jp",
					Role:          int32(AdminRoleAdministrator),
				},
			},
			expect: &response.AuthUser{
				ID:            "admin-id",
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://and-period.jp",
				Role:          int32(AdminRoleAdministrator),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.user.Response())
		})
	}
}
