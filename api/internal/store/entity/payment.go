package entity

import (
	"time"

	"gorm.io/gorm"
)

// PaymentMethodType - 決済手段
type PaymentMethodType int32

const (
	PaymentMethodTypeUnknown PaymentMethodType = 0
	PaymentMethodTypeCash    PaymentMethodType = 1 // 代引支払い
	PaymentMethodTypeCard    PaymentMethodType = 2 // クレジットカード払い
)

// Payment - 注文支払い情報
type Payment struct {
	OrderID       string            `gorm:"primaryKey;<-:create"` // 注文履歴ID
	AddressID     string            `gorm:""`                     // 請求先情報ID
	TransactionID string            `gorm:""`                     // 決済ID(Stripe用)
	MethodType    PaymentMethodType `gorm:""`                     // 決済手段種別
	MethodID      string            `gorm:"default:null"`         // 決済手段ID
	Subtotal      int64             `gorm:""`                     // 購入金額
	Discount      int64             `gorm:""`                     // 割引金額
	ShippingFee   int64             `gorm:""`                     // 配送手数料
	Tax           int64             `gorm:""`                     // 消費税
	Total         int64             `gorm:""`                     // 合計金額
	RefundTotal   int64             `gorm:""`                     // 返金金額
	CreatedAt     time.Time         `gorm:"<-:create"`            // 登録日時
	UpdatedAt     time.Time         `gorm:""`                     // 更新日時
	DeletedAt     gorm.DeletedAt    `gorm:"default:null"`         // 削除日時
}

type Payments []*Payment

func (ps Payments) MapByOrderID() map[string]*Payment {
	res := make(map[string]*Payment, len(ps))
	for _, p := range ps {
		res[p.OrderID] = p
	}
	return res
}
