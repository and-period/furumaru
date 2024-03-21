package entity

import "time"

// UserNotification - ユーザーの通知設定
type UserNotification struct {
	UserID        string    `gorm:"primaryKey;<-:create"` // ユーザーID
	EmailDisabled bool      `gorm:""`                     // メール通知の停止
	CreatedAt     time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt     time.Time `gorm:""`                     // 更新日時
}

func NewUserNotification(userID string) *UserNotification {
	return &UserNotification{
		UserID:        userID,
		EmailDisabled: false,
	}
}
