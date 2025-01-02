package entity

import (
	"time"

	"gorm.io/gorm"
)

// AdminType - 管理者種別
type AdminType int32

const (
	AdminTypeUnknown       AdminType = 0
	AdminTypeAdministrator AdminType = 1 // 管理者
	AdminTypeCoordinator   AdminType = 2 // コーディネータ
	AdminTypeProducer      AdminType = 3 // 生産者
)

// AdminGroup - 管理者グループ情報
type AdminGroup struct {
	ID             string         `gorm:"primaryKey;<-:create"` // 管理者グループID
	Type           AdminType      `gorm:"<-:create"`            // 管理者グループ種別
	Name           string         `gorm:""`                     // 管理者グループ名
	Description    string         `gorm:""`                     // 説明
	CreatedAdminID string         `gorm:""`                     // 登録者ID
	UpdatedAdminID string         `gorm:""`                     // 更新者ID
	CreatedAt      time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt      time.Time      `gorm:""`                     // 更新日時
	DeletedAt      gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

// AdminGroupRole - 管理者グループと管理者権限の紐付け情報
type AdminGroupRole struct {
	AdminGroupID   string         `gorm:"primaryKey;<-:create"` // 管理者グループID
	RoleID         string         `gorm:"primaryKey;<-:create"` // 管理者権限ID
	CreatedAdminID string         `gorm:""`                     // 登録者ID
	UpdatedAdminID string         `gorm:""`                     // 更新者ID
	ExpiredAt      time.Time      `gorm:"default:null"`         // 有効期限
	CreatedAt      time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt      time.Time      `gorm:""`                     // 更新日時
	DeletedAt      gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

// AdminRole - 管理者権限情報
type AdminRole struct {
	ID             string         `gorm:"primaryKey;<-:create"` // 管理者権限ID
	Name           string         `gorm:""`                     // 管理者権限名
	Note           string         `gorm:""`                     // 備考
	CreatedAdminID string         `gorm:""`                     // 登録者ID
	UpdatedAdminID string         `gorm:""`                     // 更新者ID
	CreatedAt      time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt      time.Time      `gorm:""`                     // 更新日時
	DeletedAt      gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

// AdminPolicy - 管理者ポリシー情報
type AdminPolicy struct {
	ID          string    `gorm:"primaryKey;<-:create"` // 管理者ポリシーID
	Name        string    `gorm:""`                     // 管理者ポリシー名
	Description string    `gorm:""`                     // 説明
	Path        string    `gorm:""`                     // マッチパターン - Path
	Method      string    `gorm:""`                     // マッチパターン - Method
	CreatedAt   time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt   time.Time `gorm:""`                     // 更新日時
}

// AdminRolePolicy - 管理者ロールと管理者ポリシーの紐付け情報
type AdminRolePolicy struct {
	RoleID         string    `gorm:"primaryKey;<-:create"` // 管理者権限ID
	PolicyID       string    `gorm:"primaryKey;<-:create"` // 管理者ポリシーID
	CreatedAdminID string    `gorm:""`                     // 登録者ID
	UpdatedAdminID string    `gorm:""`                     // 更新者ID
	CreatedAt      time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt      time.Time `gorm:""`                     // 更新日時
}
