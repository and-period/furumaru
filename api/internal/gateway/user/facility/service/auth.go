package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/auth"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
)

type Auth struct {
	types.Auth
}

func NewAuth(userID string, auth *auth.Auth) *Auth {
	return &Auth{
		Auth: types.Auth{
			UserID:       userID,
			AccessToken:  auth.AccessToken,
			RefreshToken: auth.RefreshToken,
			ExpiresIn:    auth.ExpiresIn,
			TokenType:    util.AuthTokenType,
		},
	}
}

func (a *Auth) Response() *types.Auth {
	return &a.Auth
}
