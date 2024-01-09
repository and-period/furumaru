package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

type Auth struct {
	response.Auth
}

func NewAuth(auth *entity.UserAuth) *Auth {
	return &Auth{
		Auth: response.Auth{
			UserID:       auth.UserID,
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

type AuthUser struct {
	response.AuthUser
}

func NewAuthUser(user *entity.User) *AuthUser {
	return &AuthUser{
		AuthUser: response.AuthUser{
			ID:            user.ID,
			Username:      user.Member.Username,
			AccountID:     user.Member.AccountID,
			ThumbnailURL:  user.Member.ThumbnailURL,
			Thumbnails:    NewImages(user.Member.Thumbnails).Response(),
			Lastname:      user.Member.Lastname,
			Firstname:     user.Member.Firstname,
			LastnameKana:  user.Member.LastnameKana,
			FirstnameKana: user.Member.FirstnameKana,
			Email:         user.Email(),
		},
	}
}

func (u *AuthUser) Response() *response.AuthUser {
	return &u.AuthUser
}
