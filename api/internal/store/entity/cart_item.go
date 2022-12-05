package entity

import "time"

// CartItem - カート内の商品情報
type CartItem struct {
	CartID    string    `gorm:"primaryKey;<-:create"` // カートID
	ProductID string    `gorm:"primaryKey;<-:create"` // 商品ID
	Quantity  int64     `gorm:""`                     // 数量
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type CartItems []*CartItem
