package entity

import (
	"strings"
	"time"
)

// Guest - ゲスト情報
type Guest struct {
	UserID        string    `gorm:"primaryKey;<-:create"` // ユーザーID
	Lastname      string    `gorm:""`                     // 姓
	Firstname     string    `gorm:""`                     // 名
	LastnameKana  string    `gorm:""`                     // 姓（かな）
	FirstnameKana string    `gorm:""`                     // 名（かな）
	Email         string    `gorm:""`                     // メールアドレス
	CreatedAt     time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt     time.Time `gorm:""`                     // 更新日時
}

type Guests []*Guest

func (g *Guest) Name() string {
	str := strings.Join([]string{g.Lastname, g.Firstname}, " ")
	return strings.TrimSpace(str)
}

func (gs Guests) Map() map[string]*Guest {
	res := make(map[string]*Guest, len(gs))
	for _, g := range gs {
		res[g.UserID] = g
	}
	return res
}
