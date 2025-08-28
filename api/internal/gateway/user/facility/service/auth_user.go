package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

type AuthUser struct {
	response.AuthUser
}

func NewAuthUser(user *entity.User) *AuthUser {
	return &AuthUser{
		AuthUser: response.AuthUser{
			ID:            user.ID,
			Firstname:     user.FacilityUser.Firstname,
			Lastname:      user.FacilityUser.Lastname,
			LastnameKana:  user.FacilityUser.LastnameKana,
			FirstnameKana: user.FacilityUser.FirstnameKana,
			Email:         user.FacilityUser.Email,
			PhoneNumber:   user.FacilityUser.PhoneNumber,
			LastCheckInAt: jst.Unix(user.FacilityUser.LastCheckInAt),
		},
	}
}

func (u *AuthUser) Response() *response.AuthUser {
	return &u.AuthUser
}
