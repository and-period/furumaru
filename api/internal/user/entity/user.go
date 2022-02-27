package entity

import (
	"time"

	"gorm.io/gorm"
)

type ProviderType int32

const (
	ProviderTypeUnknown ProviderType = 0
	ProviderTypeEmail   ProviderType = 1
	ProviderTypeOAuth   ProviderType = 2
)

type User struct {
	ID           string         `gorm:"primaryKey;<-:create"`
	CognitoID    string         `gorm:""`
	ProviderType ProviderType   `gorm:""`
	Email        string         `gorm:"default:null"`
	PhoneNumber  string         `gorm:"default:null"`
	CreatedAt    time.Time      `gorm:"<-:create"`
	UpdatedAt    time.Time      `gorm:""`
	VerifiedAt   time.Time      `gorm:"default:null"`
	DeletedAt    gorm.DeletedAt `gorm:"default:null"`
}
