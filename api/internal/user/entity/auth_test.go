package entity

import (
	"testing"

	"github.com/and-period/marche/api/pkg/cognito"
	"github.com/and-period/marche/api/proto/user"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		userID string
		result *cognito.AuthResult
		expect *Auth
	}{
		{
			name:   "success",
			userID: "user-id",
			result: &cognito.AuthResult{
				IDToken:      "id-token",
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
			},
			expect: &Auth{
				UserID:       "user-id",
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
			actual := NewAuth(tt.userID, tt.result)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAuth_Proto(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		auth   *Auth
		expect *user.Auth
	}{
		{
			name: "success",
			auth: &Auth{
				UserID:       "user-id",
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
			},
			expect: &user.Auth{
				UserId:       "user-id",
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
			assert.Equal(t, tt.expect, tt.auth.Proto())
		})
	}
}
