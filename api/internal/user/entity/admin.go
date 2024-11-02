package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

var errInvalidAdminRole = errors.New("entity: invalid admin role")

// AdminRole - 管理者権限
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

// Admin - 管理者共通情報
type Admin struct {
	ID            string         `gorm:"primaryKey;<-:create"` // 管理者ID
	CognitoID     string         `gorm:"default:null"`         // 管理者ID (Cognito用)
	Role          AdminRole      `gorm:"<-:create"`            // 管理者権限
	Status        AdminStatus    `gorm:"-"`                    // 管理者ステータス
	Lastname      string         `gorm:"default:null"`         // 姓
	Firstname     string         `gorm:"default:null"`         // 名
	LastnameKana  string         `gorm:"default:null"`         // 姓(かな)
	FirstnameKana string         `gorm:"default:null"`         // 名(かな)
	Email         string         `gorm:"default:null"`         // メールアドレス
	Device        string         `gorm:"default:null"`         // デバイストークン(Push通知用)
	FirstSignInAt time.Time      `gorm:"default:null"`         // 初回ログイン日時
	LastSignInAt  time.Time      `gorm:"default:null"`         // 最終ログイン日時
	CreatedAt     time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt     time.Time      `gorm:""`                     // 更新日時
	DeletedAt     gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type Admins []*Admin

type NewAdminParams struct {
	CognitoID     string
	Role          AdminRole
	Lastname      string
	Firstname     string
	LastnameKana  string
	FirstnameKana string
	Email         string
}

func NewAdminRole(role int32) (AdminRole, error) {
	res := AdminRole(role)
	if err := res.Validate(); err != nil {
		return AdminRoleUnknown, err
	}
	return res, nil
}

func (r AdminRole) Validate() error {
	switch r {
	case AdminRoleAdministrator, AdminRoleCoordinator, AdminRoleProducer:
		return nil
	default:
		return errInvalidAdminRole
	}
}

func NewAdmin(params *NewAdminParams) *Admin {
	return &Admin{
		ID:            uuid.Base58Encode(uuid.New()),
		CognitoID:     strings.ToLower(params.CognitoID), // Cognitoでは大文字小文字の区別がされず管理されているため
		Role:          params.Role,
		Lastname:      params.Lastname,
		Firstname:     params.Firstname,
		LastnameKana:  params.LastnameKana,
		FirstnameKana: params.FirstnameKana,
		Email:         params.Email,
	}
}

func (a *Admin) Name() string {
	return strings.TrimSpace(strings.Join([]string{a.Lastname, a.Firstname}, " "))
}

func (a *Admin) Fill() {
	if a.Role == AdminRoleProducer {
		// 生産者は認証機能を持たないため、一律無効状態にする
		a.Status = AdminStatusDeactivated
		return
	}
	switch {
	case !a.DeletedAt.Time.IsZero():
		a.Status = AdminStatusDeactivated
	case a.FirstSignInAt.IsZero():
		a.Status = AdminStatusInvited
	default:
		a.Status = AdminStatusActivated
	}
}

func (as Admins) Map() map[string]*Admin {
	res := make(map[string]*Admin, len(as))
	for _, a := range as {
		res[a.ID] = a
	}
	return res
}

func (as Admins) GroupByRole() map[AdminRole]Admins {
	const maxRoles = 4
	res := make(map[AdminRole]Admins, maxRoles)
	for _, a := range as {
		if _, ok := res[a.Role]; !ok {
			res[a.Role] = make(Admins, 0, len(as))
		}
		res[a.Role] = append(res[a.Role], a)
	}
	return res
}

func (as Admins) IDs() []string {
	res := make([]string, len(as))
	for i := range as {
		res[i] = as[i].ID
	}
	return res
}

func (as Admins) Devices() []string {
	set := set.NewEmpty[string](len(as))
	for i := range as {
		if as[i].Device == "" {
			continue
		}
		set.Add(as[i].Device)
	}
	return set.Slice()
}

func (as Admins) Fill() {
	for i := range as {
		as[i].Fill()
	}
}
