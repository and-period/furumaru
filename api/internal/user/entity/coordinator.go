package entity

import (
	"time"

	"gorm.io/gorm"
)

// Coordinator - コーディネータ情報
type Coordinator struct {
	Admin            `gorm:"-"`
	AdminID          string         `gorm:"primaryKey;<-:create"` // 管理者ID
	PhoneNumber      string         `gorm:""`                     // 電話番号
	CompanyName      string         `gorm:""`                     // 会社名
	StoreName        string         `gorm:""`                     // 店舗名
	ThumbnailURL     string         `gorm:""`                     // サムネイルURL
	HeaderURL        string         `gorm:""`                     // ヘッダー画像URL
	TwitterAccount   string         `gorm:""`                     // SNS(Twitter)アカウント名
	InstagramAccount string         `gorm:""`                     // SNS(Instagram)アカウント名
	FacebookAccount  string         `gorm:""`                     // SNS(Facebook)アカウント名
	PostalCode       string         `gorm:""`                     // 郵便番号
	Prefecture       string         `gorm:""`                     // 都道府県
	City             string         `gorm:""`                     // 市区町村
	AddressLine1     string         `gorm:""`                     // 町名・番地
	AddressLine2     string         `gorm:""`                     // ビル名・号室など
	CreatedAt        time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt        time.Time      `gorm:""`                     // 更新日時
	DeletedAt        gorm.DeletedAt `gorm:"default:null"`         // 退会日時
}

type Coordinators []*Coordinator

type NewCoordinatorParams struct {
	Admin            *Admin
	PhoneNumber      string
	CompanyName      string
	StoreName        string
	ThumbnailURL     string
	HeaderURL        string
	TwitterAccount   string
	InstagramAccount string
	FacebookAccount  string
	PostalCode       string
	Prefecture       string
	City             string
	AddressLine1     string
	AddressLine2     string
}

func NewCoordinator(params *NewCoordinatorParams) *Coordinator {
	return &Coordinator{
		AdminID:          params.Admin.ID,
		PhoneNumber:      params.PhoneNumber,
		CompanyName:      params.CompanyName,
		StoreName:        params.StoreName,
		ThumbnailURL:     params.ThumbnailURL,
		HeaderURL:        params.HeaderURL,
		TwitterAccount:   params.TwitterAccount,
		InstagramAccount: params.InstagramAccount,
		FacebookAccount:  params.FacebookAccount,
		PostalCode:       params.PostalCode,
		Prefecture:       params.Prefecture,
		City:             params.City,
		AddressLine1:     params.AddressLine1,
		AddressLine2:     params.AddressLine2,
		Admin:            *params.Admin,
	}
}

func (c *Coordinator) Fill(admin *Admin) {
	c.Admin = *admin
}

func (cs Coordinators) IDs() []string {
	res := make([]string, len(cs))
	for i := range cs {
		res[i] = cs[i].AdminID
	}
	return res
}
