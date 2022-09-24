package entity

import (
	"strings"
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

type AdministratorOrderBy string

const (
	AdministratorOrderByLastname    AdministratorOrderBy = "lastname"
	AdministratorOrderByFirstname   AdministratorOrderBy = "firstname"
	AdministratorOrderByEmail       AdministratorOrderBy = "email"
	AdministratorOrderByPhoneNumber AdministratorOrderBy = "phone_number"
)

// Administrator - システム管理者情報
type Administrator struct {
	ID            string         `gorm:"primaryKey;<-:create"` // Deprecated: 管理者ID
	AdminID       string         `gorm:""`                     // 管理者ID
	Lastname      string         `gorm:""`                     // Deprecated: 姓
	Firstname     string         `gorm:""`                     // Deprecated: 名
	LastnameKana  string         `gorm:""`                     // Deprecated: 姓(かな)
	FirstnameKana string         `gorm:""`                     // Deprecated: 名(かな)
	Email         string         `gorm:""`                     // Deprecated: メールアドレス
	PhoneNumber   string         `gorm:""`                     // 電話番号
	CreatedAt     time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt     time.Time      `gorm:""`                     // 更新日時
	DeletedAt     gorm.DeletedAt `gorm:"default:null"`         // 退会日時
}

type Administrators []*Administrator

type NewAdministratorParams struct {
	Lastname      string
	Firstname     string
	LastnameKana  string
	FirstnameKana string
	Email         string
	PhoneNumber   string
}

func NewAdministrator(params *NewAdministratorParams) *Administrator {
	return &Administrator{
		ID:            uuid.Base58Encode(uuid.New()),
		Lastname:      params.Lastname,
		Firstname:     params.Firstname,
		LastnameKana:  params.LastnameKana,
		FirstnameKana: params.FirstnameKana,
		Email:         params.Email,
		PhoneNumber:   params.PhoneNumber,
	}
}

func (a *Administrator) Name() string {
	return strings.TrimSpace(strings.Join([]string{a.Lastname, a.Firstname}, " "))
}

func (as Administrators) IDs() []string {
	res := make([]string, len(as))
	for i := range as {
		res[i] = as[i].ID
	}
	return res
}
