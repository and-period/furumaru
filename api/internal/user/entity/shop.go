package entity

import (
	"time"

	"gorm.io/gorm"
)

// Shop - 販売者情報
type Shop struct {
	ID        string         `gorm:"primaryKey;<-:create"`
	CognitoID string         `gorm:""`
	Name      string         `gorm:""`
	Email     string         `gorm:"default:null"`
	CreatedAt time.Time      `gorm:"<-:create"`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"default:null"`
}

type Shops []*Shop

func (ss Shops) Map() map[string]*Shop {
	res := make(map[string]*Shop, len(ss))
	for _, s := range ss {
		res[s.ID] = s
	}
	return res
}
