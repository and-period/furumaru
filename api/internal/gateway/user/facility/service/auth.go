package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/auth"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/response"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

type Auth struct {
	response.Auth
}

func NewAuth(user *entity.User, auth *auth.Auth) *Auth {
	return &Auth{
		Auth: response.Auth{
			UserID:       user.ID,
			AccessToken:  auth.AccessToken,
			RefreshToken: auth.RefreshToken,
			ExpiresIn:    auth.ExpiresIn,
			TokenType:    util.AuthTokenType,
		},
	}
}

func (a *Auth) Response() *response.Auth {
	return &a.Auth
}
