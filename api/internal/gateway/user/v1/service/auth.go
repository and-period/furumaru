package service

import (
	"github.com/and-period/marche/api/internal/gateway/entity"
	"github.com/and-period/marche/api/internal/gateway/user/v1/response"
	"github.com/and-period/marche/api/internal/gateway/util"
)

type Auth struct {
	*response.Auth
}

func NewAuth(auth *entity.UserAuth) *Auth {
	return &Auth{
		Auth: &response.Auth{
			UserID:       auth.UserId,
			AccessToken:  auth.AccessToken,
			RefreshToken: auth.RefreshToken,
			ExpiresIn:    auth.ExpiresIn,
			TokenType:    util.AuthTokenType,
		},
	}
}

func (a *Auth) Response() *response.Auth {
	return a.Auth
}
