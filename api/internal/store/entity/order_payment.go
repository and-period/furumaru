package entity

import "time"

// PaymentType - 決済手段
type PaymentType int32

const (
	PaymentTypeUnknown PaymentType = 0
	PaymentTypeCash    PaymentType = 1 // 代引支払い
	PaymentTypeCard    PaymentType = 2 // クレジットカード払い
)

// OrderPayment - 支払い情報
type OrderPayment struct {
	ID             string      `gorm:"primaryKey;<-:create"` // 決済ID
	TransactionID  string      `gorm:""`                     // 決済ID(Stripe用)
	OrderID        string      `gorm:""`                     // 注文履歴ID
	PromotionID    string      `gorm:"default:null"`         // プロモーションID
	PaymentID      string      `gorm:"default:null"`         // 決済手段ID
	PaymentType    PaymentType `gorm:""`                     // 決済手段
	Subtotal       int64       `gorm:""`                     // 購入金額
	Discount       int64       `gorm:""`                     // 割引金額
	ShippingCharge int64       `gorm:""`                     // 配送料金
	Tax            int64       `gorm:""`                     // 消費税
	Total          int64       `gorm:""`                     // 支払合計金額
	Lastname       string      `gorm:""`                     // 請求先情報 姓
	Firstname      string      `gorm:""`                     // 請求先情報 名
	PostalCode     string      `gorm:""`                     // 請求先情報 郵便番号
	Prefecture     string      `gorm:""`                     // 請求先情報 都道府県
	City           string      `gorm:""`                     // 請求先情報 市区町村
	AddressLine1   string      `gorm:""`                     // 請求先情報 町名・番地
	AddressLine2   string      `gorm:""`                     // 請求先情報 ビル名・号室など
	PhoneNumber    string      `gorm:""`                     // 請求先情報 電話番号
	CreatedAt      time.Time   `gorm:"<-:create"`            // 登録日時
	UpdatedAt      time.Time   `gorm:""`                     // 更新日時
}

type OrderPayments []*OrderPayment

func (ps OrderPayments) MapByOrderID() map[string]*OrderPayment {
	res := make(map[string]*OrderPayment, len(ps))
	for _, p := range ps {
		res[p.OrderID] = p
	}
	return res
}
