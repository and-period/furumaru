package entity

import "github.com/and-period/marche/api/pkg/cognito"

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
