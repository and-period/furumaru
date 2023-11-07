package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
)

// 支払いステータス
type PaymentStatus int32

const (
	PaymentStatusUnknown    PaymentStatus = 0
	PaymentStatusPending    PaymentStatus = 1 // 保留中
	PaymentStatusAuthorized PaymentStatus = 2 // 仮売上・オーソリ
	PaymentStatusCaptured   PaymentStatus = 3 // 実売上・キャプチャ
	PaymentStatusRefunded   PaymentStatus = 4 // 返金
	PaymentStatusFailed     PaymentStatus = 5 // 失敗/期限切れ
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

// 注文キャンセル種別
type RefundType int32

const (
	RefundTypeNone RefundType = 0 // キャンセルなし
)

// OrderPayment - 注文支払い情報
type OrderPayment struct {
	OrderID           string            `gorm:"primaryKey;<-:create"` // 注文履歴ID
	AddressRevisionID int64             `gorm:""`                     // 請求先情報ID
	Status            PaymentStatus     `gorm:""`                     // 決済状況
	TransactionID     string            `gorm:""`                     // 決済ID(決済代行システム用)
	MethodType        PaymentMethodType `gorm:""`                     // 決済手段種別
	Subtotal          int64             `gorm:""`                     // 購入金額
	Discount          int64             `gorm:""`                     // 割引金額
	ShippingFee       int64             `gorm:""`                     // 配送手数料
	Tax               int64             `gorm:""`                     // 消費税
	Total             int64             `gorm:""`                     // 合計金額
	RefundTotal       int64             `gorm:""`                     // 返金金額
	RefundType        RefundType        `gorm:""`                     // 注文キャンセル種別
	RefundReason      string            `gorm:""`                     // 注文キャンセル理由
	OrderedAt         time.Time         `gorm:"default:null"`         // 決済要求日時
	PaidAt            time.Time         `gorm:"default:null"`         // 決済承認日時(仮売上)
	CapturedAt        time.Time         `gorm:"default:null"`         // 決済確定日時(実売上)
	FailedAt          time.Time         `gorm:"default:null"`         // 決済失敗日時
	RefundedAt        time.Time         `gorm:"default:null"`         // 注文キャンセル日時(返金)
	CreatedAt         time.Time         `gorm:"<-:create"`            // 登録日時
	UpdatedAt         time.Time         `gorm:""`                     // 更新日時
}

type OrderPayments []*OrderPayment

func (p *OrderPayment) IsCanceled() bool {
	return p.Status == PaymentStatusRefunded
}

func (ps OrderPayments) AddressRevisionIDs() []int64 {
	return set.UniqBy(ps, func(p *OrderPayment) int64 {
		return p.AddressRevisionID
	})
}

func (ps OrderPayments) MapByOrderID() map[string]*OrderPayment {
	res := make(map[string]*OrderPayment, len(ps))
	for _, p := range ps {
		res[p.OrderID] = p
	}
	return res
}
