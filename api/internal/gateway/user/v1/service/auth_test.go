package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		auth   *entity.UserAuth
		expect *Auth
	}{
		{
			name: "success",
			auth: &entity.UserAuth{
				UserID:       "user-id",
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
			},
			expect: &Auth{
				Auth: &response.Auth{
					UserID:       "user-id",
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
					UserID:       "user-id",
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresIn:    3600,
					TokenType:    "Bearer",
				},
			},
			expect: &response.Auth{
				UserID:       "user-id",
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
		user   *entity.User
		expect *AuthUser
	}{
		{
			name: "success",
			user: &entity.User{
				ID:           "user-id",
				ProviderType: entity.ProviderTypeEmail,
				Email:        "test@and-period.jp",
				PhoneNumber:  "+819012345678",
				Username:     "username",
				ThumbnailURL: "https://and-period.jp/thumbnail.png",
				CreatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
				VerifiedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &AuthUser{
				AuthUser: &response.AuthUser{
					ID:           "user-id",
					Username:     "username",
					ThumbnailURL: "https://and-period.jp/thumbnail.png",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewAuthUser(tt.user))
		})
	}
}

func TestAuthUser_Response(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		auth   *AuthUser
		expect *response.AuthUser
	}{
		{
			name: "success",
			auth: &AuthUser{
				AuthUser: &response.AuthUser{
					ID:           "user-id",
					Username:     "username",
					ThumbnailURL: "https://and-period.jp/thumbnail.png",
				},
			},
			expect: &response.AuthUser{
				ID:           "user-id",
				Username:     "username",
				ThumbnailURL: "https://and-period.jp/thumbnail.png",
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
