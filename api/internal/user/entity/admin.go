package entity

import (
	"time"

	"gorm.io/gorm"
)

// AdminRole - 管理者権限
type AdminRole int32

const (
	AdminRoleUnknown       AdminRole = 0
	AdminRoleAdministrator AdminRole = 1 // 管理者
	AdminRoleDeveloper     AdminRole = 2 // 開発者
	AdminRoleOperator      AdminRole = 3 // 運用者
)

// Admin - 管理者情報
type Admin struct {
	ID        string         `gorm:"primaryKey;<-:create"`
	CognitoID string         `gorm:""`
	Email     string         `gorm:"default:null"`
	Role      AdminRole      `gorm:""`
	CreatedAt time.Time      `gorm:"<-:create"`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"default:null"`
}
