package entity

import (
	"github.com/and-period/marche/api/pkg/cognito"
	"github.com/and-period/marche/api/proto/user"
)

type Auth struct {
	UserID       string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int32
}

func NewAuth(userID string, rs *cognito.AuthResult) *Auth {
	return &Auth{
		UserID:       userID,
		AccessToken:  rs.AccessToken,
		RefreshToken: rs.RefreshToken,
		ExpiresIn:    rs.ExpiresIn,
	}
}

func (a *Auth) Proto() *user.Auth {
	return &user.Auth{
		UserId:       a.UserID,
		AccessToken:  a.AccessToken,
		RefreshToken: a.RefreshToken,
		ExpiresIn:    a.ExpiresIn,
	}
}
