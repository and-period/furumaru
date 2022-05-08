package entity

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Shop - 販売者情報
type Shop struct {
	ID            string         `gorm:"primaryKey;<-:create"`
	CognitoID     string         `gorm:""`
	Lastname      string         `gorm:""`
	Firstname     string         `gorm:""`
	LastnameKana  string         `gorm:""`
	FirstnameKana string         `gorm:""`
	Email         string         `gorm:"default:null"`
	CreatedAt     time.Time      `gorm:"<-:create"`
	UpdatedAt     time.Time      `gorm:""`
	DeletedAt     gorm.DeletedAt `gorm:"default:null"`
}

func (s *Shop) Name() string {
	return fmt.Sprintf("%s %s", s.Lastname, s.Firstname)
}

type Shops []*Shop

func (ss Shops) Map() map[string]*Shop {
	res := make(map[string]*Shop, len(ss))
	for _, s := range ss {
		res[s.ID] = s
	}
	return res
}
