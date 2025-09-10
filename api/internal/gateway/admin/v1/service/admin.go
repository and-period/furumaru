package service

import (
	"strings"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

// AdminType - 管理者ロール
type AdminType types.AdminType

// AdminStatus - 管理者ステータス
type AdminStatus types.AdminStatus

type Admin struct {
	types.Admin
}

type Admins []*Admin

func NewAdminType(role entity.AdminType) AdminType {
	switch role {
	case entity.AdminTypeAdministrator:
		return AdminType(types.AdminTypeAdministrator)
	case entity.AdminTypeCoordinator:
		return AdminType(types.AdminTypeCoordinator)
	case entity.AdminTypeProducer:
		return AdminType(types.AdminTypeProducer)
	default:
		return AdminType(types.AdminTypeUnknown)
	}
}

func (r AdminType) String() string {
	switch types.AdminType(r) {
	case types.AdminTypeAdministrator:
		return "administrator"
	case types.AdminTypeCoordinator:
		return "coordinator"
	case types.AdminTypeProducer:
		return "producer"
	default:
		return "unknown"
	}
}

func (r AdminType) IsCoordinator() bool {
	return types.AdminType(r) == types.AdminTypeCoordinator
}

func (r AdminType) Response() types.AdminType {
	return types.AdminType(r)
}

func NewAdminStatus(status entity.AdminStatus) AdminStatus {
	switch status {
	case entity.AdminStatusInvited:
		return AdminStatus(types.AdminStatusInvited)
	case entity.AdminStatusActivated:
		return AdminStatus(types.AdminStatusActivated)
	case entity.AdminStatusDeactivated:
		return AdminStatus(types.AdminStatusDeactivated)
	default:
		return AdminStatus(types.AdminStatusUnknown)
	}
}

func (s AdminStatus) Response() types.AdminStatus {
	return types.AdminStatus(s)
}

func NewAdmin(admin *entity.Admin) *Admin {
	return &Admin{
		Admin: types.Admin{
			ID:            admin.ID,
			Type:          NewAdminType(admin.Type).Response(),
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

func (a *Admin) Response() *types.Admin {
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

func (as Admins) Response() []*types.Admin {
	res := make([]*types.Admin, len(as))
	for i := range as {
		res[i] = as[i].Response()
	}
	return res
}
