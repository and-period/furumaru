package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

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
