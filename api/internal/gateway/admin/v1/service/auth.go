package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

type Auth struct {
	response.Auth
}

type AuthUser struct {
	response.AuthUser
}

func NewAuth(auth *entity.AdminAuth) *Auth {
	return &Auth{
		Auth: response.Auth{
			AdminID:      auth.AdminID,
			Type:         NewAdminType(auth.Type).Response(),
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

func (a *AuthUser) Response() *response.AuthUser {
	return &a.AuthUser
}
