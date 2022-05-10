package service

import uentity "github.com/and-period/marche/api/internal/user/entity"

type AdminRole int32

const (
	AdminRoleUnknown       AdminRole = 0
	AdminRoleAdministrator AdminRole = 1 // 管理者
	AdminRoleProducer      AdminRole = 2 // 生産者
)

type Auth struct{}

func NewAdminRole(role uentity.AdminRole) AdminRole {
	switch role {
	case uentity.AdminRoleAdministrator:
		return AdminRoleAdministrator
	case uentity.AdminRoleProducer:
		return AdminRoleProducer
	default:
		return AdminRoleUnknown
	}
}

func (r *AdminRole) String() string {
	switch *r {
	case AdminRoleAdministrator:
		return "admin"
	case AdminRoleProducer:
		return "producer"
	default:
		return "unknown"
	}
}
