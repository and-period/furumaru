package entity

import (
	"github.com/and-period/marche/api/internal/user/entity"
)

// Admin - 管理者
type Admin struct {
	*entity.Admin
}

type Admins []*Admin

func NewAdmin(admin *entity.Admin) *Admin {
	return &Admin{
		Admin: admin,
	}
}

func NewAdmins(admins entity.Admins) Admins {
	res := make(Admins, len(admins))
	for i := range admins {
		res[i] = NewAdmin(admins[i])
	}
	return res
}
