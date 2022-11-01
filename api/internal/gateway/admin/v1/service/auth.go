package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

type AdminRole int32

const (
	AdminRoleUnknown       AdminRole = 0
	AdminRoleAdministrator AdminRole = 1 // 管理者
	AdminRoleCoordinator   AdminRole = 2 // コーディネータ
	AdminRoleProducer      AdminRole = 3 // 生産者
)

type Auth struct {
	response.Auth
}

func NewAdminRole(role entity.AdminRole) AdminRole {
	switch role {
	case entity.AdminRoleAdministrator:
		return AdminRoleAdministrator
	case entity.AdminRoleCoordinator:
		return AdminRoleCoordinator
	case entity.AdminRoleProducer:
		return AdminRoleProducer
	default:
		return AdminRoleUnknown
	}
}

func (r AdminRole) String() string {
	switch r {
	case AdminRoleAdministrator:
		return "administrator"
	case AdminRoleCoordinator:
		return "coordinator"
	case AdminRoleProducer:
		return "producer"
	default:
		return "unknown"
	}
}

func (r AdminRole) IsCoordinator() bool {
	return r == AdminRoleCoordinator
}

func (r AdminRole) Response() int32 {
	return int32(r)
}

func NewAuth(auth *entity.AdminAuth) *Auth {
	return &Auth{
		Auth: response.Auth{
			AdminID:      auth.AdminID,
			Role:         NewAdminRole(auth.Role).Response(),
			AccessToken:  auth.AccessToken,
			RefreshToken: auth.RefreshToken,
			ExpiresIn:    auth.ExpiresIn,
			TokenType:    util.AuthTokenType,
		},
	}
}

func (a *Auth) Response() *response.Auth {
	return &a.Auth
}
