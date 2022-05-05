package service

import uentity "github.com/and-period/marche/api/internal/user/entity"

type AdminRole int32

const (
	AdminRoleUnknown       AdminRole = 0
	AdminRoleAdministrator AdminRole = 1 // 管理者
	AdminRoleDeveloper     AdminRole = 2 // 開発者
	AdminRoleOperator      AdminRole = 3 // 運用者
)

type Auth struct{}

func NewAdminRole(role uentity.AdminRole) AdminRole {
	switch role {
	case uentity.AdminRoleAdministrator:
		return AdminRoleAdministrator
	case uentity.AdminRoleDeveloper:
		return AdminRoleDeveloper
	case uentity.AdminRoleOperator:
		return AdminRoleOperator
	default:
		return AdminRoleUnknown
	}
}

func (r *AdminRole) String() string {
	switch *r {
	case AdminRoleAdministrator:
		return "admin"
	case AdminRoleDeveloper:
		return "developer"
	case AdminRoleOperator:
		return "operator"
	default:
		return "unknown"
	}
}
