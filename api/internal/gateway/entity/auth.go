package entity

import "github.com/and-period/marche/api/proto/user"

type UserAuth struct {
	*user.UserAuth
}

func NewUserAuth(auth *user.UserAuth) *UserAuth {
	return &UserAuth{
		UserAuth: auth,
	}
}
