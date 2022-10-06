package entity

import "time"

// OrderItem - 注文商品情報
type OrderItem struct {
	ID         string     `gorm:"primaryKey;<-:create"` // 注文商品ID
	OrderID    string     `gorm:""`                     // 注文履歴ID
	ProductID  string     `gorm:""`                     // 商品ID
	Price      int64      `gorm:""`                     // 購入価格
	Quantity   int64      `gorm:""`                     // 購入数量
	Weight     int64      `gorm:""`                     // 商品重量
	WeightUnit WeightUnit `gorm:""`                     // 商品重量単位
	CreatedAt  time.Time  `gorm:"<-:create"`            // 登録日時
	UpdatedAt  time.Time  `gorm:""`                     // 更新日時
}
