package entity

import "github.com/and-period/marche/api/proto/user"

type Auth struct {
	*user.Auth
}

func NewAuth(auth *user.Auth) *Auth {
	return &Auth{
		Auth: auth,
	}
}
