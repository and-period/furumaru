package entity

import (
	"time"

	"gorm.io/gorm"
)

// Shop - 店舗情報
type Shop struct {
	ID            string         `gorm:"primaryKey;<-:create"` // 店舗ID
	CoordinatorID string         `gorm:""`                     // コーディネータID
	ProducerIDs   []string       `gorm:"-"`                    // 生産者ID一覧
	Name          string         `gorm:""`                     // 店舗名
	Activated     bool           `gorm:""`                     // 有効フラグ
	CreatedAt     time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt     time.Time      `gorm:""`                     // 更新日時
	DeletedAt     gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type Shops []Shop
