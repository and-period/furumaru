package entity

import (
	"time"
)

// Guest - ゲスト情報
type Guest struct {
	UserID      string    `gorm:"primaryKey;<-:create"` // ユーザーID
	Email       string    `gorm:""`                     // メールアドレス
	PhoneNumber string    `gorm:""`                     // 電話番号
	CreatedAt   time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt   time.Time `gorm:""`                     // 更新日時
}

type Guests []*Guest

func (gs Guests) Map() map[string]*Guest {
	res := make(map[string]*Guest, len(gs))
	for _, g := range gs {
		res[g.UserID] = g
	}
	return res
}
