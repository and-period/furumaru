package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/stretchr/testify/assert"
)

func TestUserAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		userID string
		result *cognito.AuthResult
		expect *UserAuth
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
			expect: &UserAuth{
				UserID:       "user-id",
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewUserAuth(tt.userID, tt.result)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
