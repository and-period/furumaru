package entity

import (
	"testing"

	"github.com/and-period/marche/api/proto/user"
	"github.com/stretchr/testify/assert"
)

func TestUserAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		auth   *user.UserAuth
		expect *UserAuth
	}{
		{
			name: "success",
			auth: &user.UserAuth{
				UserId:       "user-id",
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
			},
			expect: &UserAuth{
				UserAuth: &user.UserAuth{
					UserId:       "user-id",
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresIn:    3600,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewUserAuth(tt.auth))
		})
	}
}
