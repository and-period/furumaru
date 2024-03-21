package entity

import (
	"time"
)

// UserNotification - ユーザーの通知設定
type UserNotification struct {
	UserID    string    `gorm:"primaryKey;<-:create"` // ユーザーID
	Disabled  bool      `gorm:""`                     // 通知の停止
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type UserNotifications []*UserNotification

func NewUserNotification(userID string) *UserNotification {
	return &UserNotification{
		UserID:   userID,
		Disabled: false,
	}
}

func (n *UserNotification) Enabled() bool {
	if n == nil {
		return true
	}
	return !n.Disabled
}

func (ns UserNotifications) MapByUserID() map[string]*UserNotification {
	res := make(map[string]*UserNotification, len(ns))
	for _, n := range ns {
		res[n.UserID] = n
	}
	return res
}
