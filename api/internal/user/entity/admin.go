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
	Type          AdminType      `gorm:"<-:create"`            // 管理者種別
	GroupIDs      []string       `gorm:"-"`                    // 管理者グループID一覧
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
	Type          AdminType
	Lastname      string
	Firstname     string
	LastnameKana  string
	FirstnameKana string
	Email         string
}

func NewAdminType(typ int32) (AdminType, error) {
	res := AdminType(typ)
	if err := res.Validate(); err != nil {
		return AdminTypeUnknown, err
	}
	return res, nil
}

func (r AdminType) Validate() error {
	switch r {
	case AdminTypeAdministrator, AdminTypeCoordinator, AdminTypeProducer:
		return nil
	default:
		return errInvalidAdminRole
	}
}

func NewAdmin(params *NewAdminParams) *Admin {
	return &Admin{
		ID:            uuid.Base58Encode(uuid.New()),
		CognitoID:     strings.ToLower(params.CognitoID), // Cognitoでは大文字小文字の区別がされず管理されているため
		Type:          params.Type,
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

func (a *Admin) Fill(groups AdminGroupUsers) (err error) {
	a.SetStatus()
	a.GroupIDs = groups.GroupIDs()
	return
}

func (a *Admin) SetStatus() {
	switch {
	case a.Type == AdminTypeProducer:
		// 生産者は認証機能を持たないため、一律無効状態にする
		a.Status = AdminStatusDeactivated
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

func (as Admins) GroupByType() map[AdminType]Admins {
	const maxRoles = 4
	res := make(map[AdminType]Admins, maxRoles)
	for _, a := range as {
		if _, ok := res[a.Type]; !ok {
			res[a.Type] = make(Admins, 0, len(as))
		}
		res[a.Type] = append(res[a.Type], a)
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

func (as Admins) Fill(groups map[string]AdminGroupUsers) error {
	for _, a := range as {
		if err := a.Fill(groups[a.ID]); err != nil {
			return err
		}
	}
	return nil
}
