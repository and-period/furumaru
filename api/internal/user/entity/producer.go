package entity

import (
	"strings"
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

// Producer - 生産者情報
type Producer struct {
	ID            string         `gorm:"primaryKey;<-:create"` // 管理者ID
	Lastname      string         `gorm:""`                     // 姓
	Firstname     string         `gorm:""`                     // 名
	LastnameKana  string         `gorm:""`                     // 姓(かな)
	FirstnameKana string         `gorm:""`                     // 名(かな)
	StoreName     string         `gorm:""`                     // 店舗名
	ThumbnailURL  string         `gorm:""`                     // サムネイルURL
	HeaderURL     string         `gorm:""`                     // ヘッダー画像URL
	Email         string         `gorm:""`                     // メールアドレス
	PhoneNumber   string         `gorm:""`                     // 電話番号
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
	Lastname      string
	Firstname     string
	LastnameKana  string
	FirstnameKana string
	StoreName     string
	ThumbnailURL  string
	HeaderURL     string
	Email         string
	PhoneNumber   string
	PostalCode    string
	Prefecture    string
	City          string
	AddressLine1  string
	AddressLine2  string
}

func NewProducer(params *NewProducerParams) *Producer {
	return &Producer{
		ID:            uuid.Base58Encode(uuid.New()),
		Lastname:      params.Lastname,
		Firstname:     params.Firstname,
		LastnameKana:  params.LastnameKana,
		FirstnameKana: params.FirstnameKana,
		StoreName:     params.StoreName,
		ThumbnailURL:  params.ThumbnailURL,
		HeaderURL:     params.HeaderURL,
		Email:         params.Email,
		PhoneNumber:   params.PhoneNumber,
		PostalCode:    params.PostalCode,
		Prefecture:    params.Prefecture,
		City:          params.City,
		AddressLine1:  params.AddressLine1,
		AddressLine2:  params.AddressLine2,
	}
}

func (p *Producer) Name() string {
	return strings.TrimSpace(strings.Join([]string{p.Lastname, p.Firstname}, " "))
}

func (ps Producers) IDs() []string {
	res := make([]string, len(ps))
	for i := range ps {
		res[i] = ps[i].ID
	}
	return res
}
