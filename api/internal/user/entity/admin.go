package entity

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

var errInvalidAdminRole = errors.New("entity: invalid admin role")

// AdminRole - 管理者権限
type AdminRole int32

const (
	AdminRoleUnknown       AdminRole = 0
	AdminRoleAdministrator AdminRole = 1 // 管理者
	AdminRoleCoordinator   AdminRole = 2 // 仲介者
	AdminRoleProducer      AdminRole = 3 // 生産者
)

// Admin - 管理者情報
type Admin struct {
	ID            string         `gorm:"primaryKey;<-:create"` // 管理者ID
	CognitoID     string         `gorm:""`                     // 管理者ID (Cognito用)
	Lastname      string         `gorm:""`                     // 姓
	Firstname     string         `gorm:""`                     // 名
	LastnameKana  string         `gorm:""`                     // 姓(かな)
	FirstnameKana string         `gorm:""`                     // 名(かな)
	StoreName     string         `gorm:""`                     // 店舗名
	ThumbnailURL  string         `gorm:""`                     // サムネイルURL
	Email         string         `gorm:""`                     // メールアドレス
	PhoneNumber   string         `gorm:""`                     // 電話番号
	PostalCode    string         `gorm:""`                     // 郵便番号
	Prefecture    string         `gorm:""`                     // 都道府県
	City          string         `gorm:""`                     // 市区町村
	AddressLine1  string         `gorm:""`                     // 町名・番地
	AddressLine2  string         `gorm:""`                     // ビル名・号室など
	Role          AdminRole      `gorm:""`                     // 権限
	CreatedAt     time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt     time.Time      `gorm:""`                     // 更新日時
	DeletedAt     gorm.DeletedAt `gorm:"default:null"`         // 退会日時
}

type Admins []*Admin

type NewAdministratorParams struct {
	ID            string
	CognitoID     string
	Lastname      string
	Firstname     string
	LastnameKana  string
	FirstnameKana string
	Email         string
	PhoneNumber   string
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
	case AdminRoleAdministrator, AdminRoleProducer:
		return nil
	default:
		return errInvalidAdminRole
	}
}

func NewAdminRoles(roles []int32) ([]AdminRole, error) {
	res := make([]AdminRole, len(roles))
	for i := range roles {
		role, err := NewAdminRole(roles[i])
		if err != nil {
			return nil, err
		}
		res[i] = role
	}
	return res, nil
}

func NewAdministrator(params *NewAdministratorParams) *Admin {
	return &Admin{
		ID:            params.ID,
		CognitoID:     params.CognitoID,
		Lastname:      params.Lastname,
		Firstname:     params.Firstname,
		LastnameKana:  params.LastnameKana,
		FirstnameKana: params.FirstnameKana,
		Email:         params.Email,
		PhoneNumber:   params.PhoneNumber,
		Role:          AdminRoleAdministrator,
	}
}

func (a *Admin) Name() string {
	return strings.TrimSpace(strings.Join([]string{a.Lastname, a.Firstname}, " "))
}

func (as Admins) Map() map[string]*Admin {
	res := make(map[string]*Admin, len(as))
	for _, a := range as {
		res[a.ID] = a
	}
	return res
}
