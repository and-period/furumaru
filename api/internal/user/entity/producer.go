package entity

import (
	"time"

	"gorm.io/gorm"
)

// Producer - 生産者情報
type Producer struct {
	Admin         `gorm:"-"`
	AdminID       string         `gorm:"primaryKey;<-:create"` // 管理者ID
	CoordinatorID string         `gorm:"default:null"`         // コーディネータID
	PhoneNumber   string         `gorm:""`                     // 電話番号
	StoreName     string         `gorm:""`                     // 店舗名
	ThumbnailURL  string         `gorm:""`                     // サムネイルURL
	HeaderURL     string         `gorm:""`                     // ヘッダー画像URL
	PostalCode    string         `gorm:""`                     // 郵便番号
	Prefecture    string         `gorm:""`                     // 都道府県
	City          string         `gorm:""`                     // 市区町村
	AddressLine1  string         `gorm:""`                     // 町名・番地
	AddressLine2  string         `gorm:""`                     // ビル名・号室など
	CreatedAt     time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt     time.Time      `gorm:""`                     // 更新日時
	DeletedAt     gorm.DeletedAt `gorm:"default:null"`         // 退会日時
}

type Producers []*Producer

type NewProducerParams struct {
	Admin        *Admin
	PhoneNumber  string
	StoreName    string
	ThumbnailURL string
	HeaderURL    string
	PostalCode   string
	Prefecture   string
	City         string
	AddressLine1 string
	AddressLine2 string
}

func NewProducer(params *NewProducerParams) *Producer {
	return &Producer{
		AdminID:      params.Admin.ID,
		PhoneNumber:  params.PhoneNumber,
		StoreName:    params.StoreName,
		ThumbnailURL: params.ThumbnailURL,
		HeaderURL:    params.HeaderURL,
		PostalCode:   params.PostalCode,
		Prefecture:   params.Prefecture,
		City:         params.City,
		AddressLine1: params.AddressLine1,
		AddressLine2: params.AddressLine2,
		Admin:        *params.Admin,
	}
}

func (p *Producer) Fill(admin *Admin) {
	p.Admin = *admin
}

func (ps Producers) IDs() []string {
	res := make([]string, len(ps))
	for i := range ps {
		res[i] = ps[i].AdminID
	}
	return res
}

func (ps Producers) Related() Producers {
	res := make(Producers, 0, len(ps))
	for _, p := range ps {
		if p.CoordinatorID == "" {
			continue
		}
		res = append(res, p)
	}
	return res
}

func (ps Producers) Unrelated() Producers {
	res := make(Producers, 0, len(ps))
	for _, p := range ps {
		if p.CoordinatorID != "" {
			continue
		}
		res = append(res, p)
	}
	return res
}
