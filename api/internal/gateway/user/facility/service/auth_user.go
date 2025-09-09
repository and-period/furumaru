package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

type AuthUser struct {
	types.AuthUser
}

func NewAuthUser(user *entity.User) *AuthUser {
	return &AuthUser{
		AuthUser: types.AuthUser{
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

func (u *AuthUser) Response() *types.AuthUser {
	return &u.AuthUser
}
