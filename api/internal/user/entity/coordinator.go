package entity

import (
	"time"

	"gorm.io/gorm"
)

// Coordinator - 仲介者情報
type Coordinator struct {
	ID               string         `gorm:"primaryKey;<-:create"` // 管理者ID
	Lastname         string         `gorm:""`                     // 姓
	Firstname        string         `gorm:""`                     // 名
	LastnameKana     string         `gorm:""`                     // 姓(かな)
	FirstnameKana    string         `gorm:""`                     // 名(かな)
	CompanyName      string         `gorm:""`                     // 会社名
	StoreName        string         `gorm:""`                     // 店舗名
	ThumbnailURL     string         `gorm:""`                     // サムネイルURL
	HeaderURL        string         `gorm:""`                     // ヘッダー画像URL
	TwitterAccount   string         `gorm:""`                     // SNS(Twitter)アカウント名
	InstagramAccount string         `gorm:""`                     // SNS(Instagram)アカウント名
	FacebookAccount  string         `gorm:""`                     // SNS(Facebook)アカウント名
	Email            string         `gorm:""`                     // メールアドレス
	PhoneNumber      string         `gorm:""`                     // 電話番号
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
