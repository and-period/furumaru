package entity

import (
	"time"

	"gorm.io/gorm"
)

// Store - 店舗情報
type Store struct {
	ID           int64          `gorm:"primaryKey;<-:create"`
	Name         string         `gorm:""`
	ThumbnailURL string         `gorm:""`
	CreatedAt    time.Time      `gorm:"<-:create"`
	UpdatedAt    time.Time      `gorm:""`
	DeletedAt    gorm.DeletedAt `gorm:"default:null"`
}

type Stores []*Store
