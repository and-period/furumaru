package entity

import (
	"github.com/and-period/marche/api/pkg/cognito"
	"github.com/and-period/marche/api/proto/user"
)

type UserAuth struct {
	UserID       string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int32
}

func NewUserAuth(userID string, rs *cognito.AuthResult) *UserAuth {
	return &UserAuth{
		UserID:       userID,
		AccessToken:  rs.AccessToken,
		RefreshToken: rs.RefreshToken,
		ExpiresIn:    rs.ExpiresIn,
	}
}

func (a *UserAuth) Proto() *user.UserAuth {
	return &user.UserAuth{
		UserId:       a.UserID,
		AccessToken:  a.AccessToken,
		RefreshToken: a.RefreshToken,
		ExpiresIn:    a.ExpiresIn,
	}
}
