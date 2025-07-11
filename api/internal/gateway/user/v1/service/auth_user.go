package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

type AuthUser struct {
	response.AuthUser
}

func NewAuthUser(user *entity.User, notification *entity.UserNotification) *AuthUser {
	res := &AuthUser{
		AuthUser: response.AuthUser{
			ID:            user.ID,
			Username:      user.Member.Username,
			AccountID:     user.AccountID,
			ThumbnailURL:  user.ThumbnailURL,
			Lastname:      user.Member.Lastname,
			Firstname:     user.Member.Firstname,
			LastnameKana:  user.Member.LastnameKana,
			FirstnameKana: user.Member.FirstnameKana,
			Email:         user.Email(),
		},
	}
	if notification != nil {
		res.NotificationEnabled = !notification.Disabled
	}
	return res
}

func (u *AuthUser) Response() *response.AuthUser {
	return &u.AuthUser
}
