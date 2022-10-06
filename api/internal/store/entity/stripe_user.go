package entity

import "time"

// StripeUser - Stripe顧客情報
type StripeUser struct {
	ID        string    `gorm:"primaryKey;<-:create"` // ユーザーID(Stripe用)
	UserID    string    `gorm:""`                     // ユーザーID
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}
