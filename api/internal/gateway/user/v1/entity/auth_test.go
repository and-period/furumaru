package entity

import (
	"testing"

	"github.com/and-period/marche/api/internal/gateway/entity"
	"github.com/and-period/marche/api/proto/user"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		auth   *entity.Auth
		expect *Auth
	}{
		{
			name: "success",
			auth: &entity.Auth{
				Auth: &user.Auth{
					UserId:       "user-id",
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresIn:    3600,
				},
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
			assert.Equal(t, tt.expect, NewAuth(tt.auth))
		})
	}
}
