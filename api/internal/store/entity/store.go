package entity

import (
	"time"

	"gorm.io/gorm"
)

// Store - 店舗情報
type Store struct {
	ID           int64          `gorm:"primaryKey;<-:create"` // 店舗ID
	Name         string         `gorm:""`                     // 店舗名
	ThumbnailURL string         `gorm:""`                     // サムネイルURL
	CreatedAt    time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt    time.Time      `gorm:""`                     // 更新日時
	DeletedAt    gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type Stores []*Store

func NewStore(name string) *Store {
	return &Store{
		Name: name,
	}
}
