package service

import (
	"strings"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

// AdminRole - 管理者ロール
type AdminRole int32

const (
	AdminRoleUnknown       AdminRole = 0
	AdminRoleAdministrator AdminRole = 1 // 管理者
	AdminRoleCoordinator   AdminRole = 2 // コーディネータ
	AdminRoleProducer      AdminRole = 3 // 生産者
)

// AdminStatus - 管理者ステータス
type AdminStatus int32

const (
	AdminStatusUnknown     AdminStatus = 0
	AdminStatusInvited     AdminStatus = 1 // 招待中
	AdminStatusActivated   AdminStatus = 2 // 有効
	AdminStatusDeactivated AdminStatus = 3 // 無効
)

type Admin struct {
	response.Admin
}

type Admins []*Admin

//nolint:staticcheck
func NewAdminRole(role entity.LegacyAdminRole) AdminRole {
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

func NewAdminStatus(status entity.AdminStatus) AdminStatus {
	switch status {
	case entity.AdminStatusInvited:
		return AdminStatusInvited
	case entity.AdminStatusActivated:
		return AdminStatusActivated
	case entity.AdminStatusDeactivated:
		return AdminStatusDeactivated
	default:
		return AdminStatusUnknown
	}
}

func (s AdminStatus) Response() int32 {
	return int32(s)
}

func NewAdmin(admin *entity.Admin) *Admin {
	return &Admin{
		Admin: response.Admin{
			ID:            admin.ID,
			Role:          admin.Role,
			Lastname:      admin.Lastname,
			Firstname:     admin.Firstname,
			LastnameKana:  admin.LastnameKana,
			FirstnameKana: admin.FirstnameKana,
			Email:         admin.Email,
			CreatedAt:     admin.CreatedAt.Unix(),
			UpdatedAt:     admin.UpdatedAt.Unix(),
		},
	}
}

func (a *Admin) Name() string {
	return strings.TrimSpace(strings.Join([]string{a.Lastname, a.Firstname}, " "))
}

func (a *Admin) Response() *response.Admin {
	return &a.Admin
}

func NewAdmins(admins entity.Admins) Admins {
	res := make(Admins, len(admins))
	for i := range admins {
		res[i] = NewAdmin(admins[i])
	}
	return res
}

func (as Admins) Map() map[string]*Admin {
	res := make(map[string]*Admin, len(as))
	for _, a := range as {
		res[a.ID] = a
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
