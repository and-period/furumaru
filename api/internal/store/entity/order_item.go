package entity

import "time"

// OrderItem - 注文商品情報
type OrderItem struct {
	OrderID   string    `gorm:"primaryKey;<-:create"` // 注文履歴ID
	ProductID string    `gorm:"primaryKey;<-:create"` // 商品ID
	Price     int64     `gorm:""`                     // 購入価格
	Quantity  int64     `gorm:""`                     // 購入数量
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type OrderItems []*OrderItem

func (is OrderItems) GroupByOrderID() map[string]OrderItems {
	res := make(map[string]OrderItems, len(is))
	for _, i := range is {
		if _, ok := res[i.OrderID]; !ok {
			res[i.OrderID] = make(OrderItems, 0, len(is))
		}
		res[i.OrderID] = append(res[i.OrderID], i)
	}
	return res
}
