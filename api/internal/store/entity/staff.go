package entity

import "time"

// StoreRole - 店舗スタッフ権限
type StoreRole int32

const (
	StoreRoleUnknown       StoreRole = 0
	StoreRoleAdministrator StoreRole = 1 // 管理者
	StoreRoleEditor        StoreRole = 2 // 編集者
	StoreRoleViewer        StoreRole = 3 // 閲覧者
)

// Staff - 店舗スタッフ情報
type Staff struct {
	StoreID   int64     `gorm:"primaryKey;<-:create"`
	UserID    string    `gorm:"primaryKey;<-:create"`
	Role      StoreRole `gorm:""`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time `gorm:""`
}

type Staffs []*Staff
