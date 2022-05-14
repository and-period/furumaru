package service

import (
	"github.com/and-period/marche/api/internal/gateway/admin/v1/response"
	"github.com/and-period/marche/api/internal/user/entity"
)

type AdminRole int32

const (
	AdminRoleUnknown       AdminRole = 0
	AdminRoleAdministrator AdminRole = 1 // 管理者
	AdminRoleProducer      AdminRole = 2 // 生産者
)

type Admin struct {
	*response.Admin
}

type Admins []*Admin

func NewAdminRole(role entity.AdminRole) AdminRole {
	switch role {
	case entity.AdminRoleAdministrator:
		return AdminRoleAdministrator
	case entity.AdminRoleProducer:
		return AdminRoleProducer
	default:
		return AdminRoleUnknown
	}
}

func (r AdminRole) String() string {
	switch r {
	case AdminRoleAdministrator:
		return "admin"
	case AdminRoleProducer:
		return "producer"
	default:
		return "unknown"
	}
}

func (r AdminRole) Response() int32 {
	return int32(r)
}

func NewAdmin(admin *entity.Admin) *Admin {
	return &Admin{
		Admin: &response.Admin{
			ID:            admin.ID,
			Lastname:      admin.Lastname,
			Firstname:     admin.Firstname,
			LastnameKana:  admin.LastnameKana,
			FirstnameKana: admin.FirstnameKana,
			Email:         admin.Email,
			Role:          NewAdminRole(admin.Role).Response(),
			ThumbnailURL:  admin.ThumbnailURL,
			CreatedAt:     admin.CreatedAt.Unix(),
			UpdatedAt:     admin.CreatedAt.Unix(),
		},
	}
}

func (a *Admin) Response() *response.Admin {
	return a.Admin
}

func NewAdmins(admins entity.Admins) Admins {
	res := make(Admins, len(admins))
	for i := range admins {
		res[i] = NewAdmin(admins[i])
	}
	return res
}

func (as Admins) Response() []*response.Admin {
	res := make([]*response.Admin, len(as))
	for i := range as {
		res[i] = as[i].Response()
	}
	return res
}
