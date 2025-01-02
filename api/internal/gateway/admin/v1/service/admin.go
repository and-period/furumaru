package service

import (
	"strings"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

// AdminType - 管理者ロール
type AdminType int32

const (
	AdminTypeUnknown       AdminType = 0
	AdminTypeAdministrator AdminType = 1 // 管理者
	AdminTypeCoordinator   AdminType = 2 // コーディネータ
	AdminTypeProducer      AdminType = 3 // 生産者
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

func NewAdminType(role entity.AdminType) AdminType {
	switch role {
	case entity.AdminTypeAdministrator:
		return AdminTypeAdministrator
	case entity.AdminTypeCoordinator:
		return AdminTypeCoordinator
	case entity.AdminTypeProducer:
		return AdminTypeProducer
	default:
		return AdminTypeUnknown
	}
}

func (r AdminType) String() string {
	switch r {
	case AdminTypeAdministrator:
		return "administrator"
	case AdminTypeCoordinator:
		return "coordinator"
	case AdminTypeProducer:
		return "producer"
	default:
		return "unknown"
	}
}

func (r AdminType) IsCoordinator() bool {
	return r == AdminTypeCoordinator
}

func (r AdminType) Response() int32 {
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
