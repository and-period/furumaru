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
	StoreID   int64     `gorm:"primaryKey;<-:create"` // 店舗ID
	UserID    string    `gorm:"primaryKey;<-:create"` // スタッフID
	Role      StoreRole `gorm:""`                     // 権限
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type Staffs []*Staff

func (ss Staffs) UserIDs() []string {
	res := make([]string, len(ss))
	for i := range ss {
		res[i] = ss[i].UserID
	}
	return res
}
