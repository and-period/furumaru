package entity

import (
	"time"

	"gorm.io/gorm"
)

// PaymentMethodType - 決済手段
type PaymentMethodType int32

const (
	PaymentMethodTypeUnknown     PaymentMethodType = 0
	PaymentMethodTypeCash        PaymentMethodType = 1 // 代引支払い
	PaymentMethodTypeCreditCard  PaymentMethodType = 2 // クレジットカード決済
	PaymentMethodTypeKonbini     PaymentMethodType = 3 // コンビニ決済
	PaymentMethodTypeBankTranser PaymentMethodType = 4 // 銀行振込決済
	PaymentMethodTypePayPay      PaymentMethodType = 5 // QR決済（PayPay）
	PaymentMethodTypeLinePay     PaymentMethodType = 6 // QR決済（LINE Pay）
	PaymentMethodTypeMerpay      PaymentMethodType = 7 // QR決済（メルペイ）
	PaymentMethodTypeRakutenPay  PaymentMethodType = 8 // QR決済（楽天ペイ）
	PaymentMethodTypeAUPay       PaymentMethodType = 9 // QR決済（au PAY）
)

// Payment - 注文支払い情報
type Payment struct {
	OrderID       string            `gorm:"primaryKey;<-:create"` // 注文履歴ID
	AddressID     string            `gorm:""`                     // 請求先情報ID
	TransactionID string            `gorm:""`                     // 決済ID(決済代行システム用)
	MethodType    PaymentMethodType `gorm:""`                     // 決済手段種別
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
