package entity

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

var errInvalidAdminRole = errors.New("entity: invalid admin role")

// AdminRole - 管理者権限
type AdminRole int32

const (
	AdminRoleUnknown       AdminRole = 0
	AdminRoleAdministrator AdminRole = 1 // 管理者
	AdminRoleProducer      AdminRole = 2 // 生産者
)

func (r AdminRole) Validate() error {
	switch r {
	case AdminRoleAdministrator, AdminRoleProducer:
		return nil
	default:
		return errInvalidAdminRole
	}
}

// Admin - 管理者情報
type Admin struct {
	ID            string         `gorm:"primaryKey;<-:create"`
	CognitoID     string         `gorm:""`
	Lastname      string         `gorm:""`
	Firstname     string         `gorm:""`
	LastnameKana  string         `gorm:""`
	FirstnameKana string         `gorm:""`
	Email         string         `gorm:"default:null"`
	ThumbnailURL  string         `gorm:""`
	Role          AdminRole      `gorm:""`
	CreatedAt     time.Time      `gorm:"<-:create"`
	UpdatedAt     time.Time      `gorm:""`
	DeletedAt     gorm.DeletedAt `gorm:"default:null"`
}

type Admins []*Admin

func NewAdmin(id, cognitoID, lastname, firstname, lastnameKana, firstnameKana, email string, role AdminRole) *Admin {
	return &Admin{
		ID:            id,
		CognitoID:     cognitoID,
		Lastname:      lastname,
		Firstname:     firstname,
		LastnameKana:  lastnameKana,
		FirstnameKana: firstnameKana,
		Email:         email,
		Role:          role,
	}
}

func (a *Admin) Name() string {
	return strings.Join([]string{a.Lastname, a.Firstname}, " ")
}

func (as Admins) Map() map[string]*Admin {
	res := make(map[string]*Admin, len(as))
	for _, a := range as {
		res[a.ID] = a
	}
	return res
}
