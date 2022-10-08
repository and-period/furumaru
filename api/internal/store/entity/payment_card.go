package entity

import "time"

// PaymentCard - クレジットカード情報
type PaymentCard struct {
	ID           string    `gorm:"primaryKey;<-:create"` // 決済手段ID
	UserID       string    `gorm:""`                     // ユーザーID
	StripeUserID string    `gorm:""`                     // ユーザーID(Stripe用)
	IsDefault    bool      `gorm:""`                     // デフォルト支払い方法
	Brand        string    `gorm:""`                     // カード発行会社
	ExpYear      int64     `gorm:""`                     // 有効期限(年)
	ExpMonth     int64     `gorm:""`                     // 有効期限(月)
	Last4        int64     `gorm:""`                     // 下４桁の番号
	CreatedAt    time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt    time.Time `gorm:""`                     // 更新日時
}
