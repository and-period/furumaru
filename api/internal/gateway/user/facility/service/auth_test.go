package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/auth"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/response"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		user   *entity.User
		auth   *auth.Auth
		expect *Auth
	}{
		{
			name: "success",
			user: &entity.User{
				ID: "user-id",
			},
			auth: &auth.Auth{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
			},
			expect: &Auth{
				Auth: response.Auth{
					UserID:       "user-id",
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresIn:    3600,
					TokenType:    util.AuthTokenType,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAuth(tt.user, tt.auth)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAuth_Response(t *testing.T) {
	t.Parallel()

	now := time.Now()
	auth := &Auth{
		Auth: response.Auth{
			UserID:       "user-id",
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			ExpiresIn:    3600,
			TokenType:    "Bearer",
		},
	}

	expected := &response.Auth{
		UserID:       "user-id",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
		TokenType:    "Bearer",
	}

	assert.Equal(t, expected, auth.Response())
	assert.WithinDuration(t, now, now, time.Second)
}