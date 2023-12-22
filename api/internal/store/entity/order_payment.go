package entity

import (
	"errors"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/store/komoju"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/set"
)

var errNotFoundAddress = errors.New("entity: not found address")

// 支払いステータス
type PaymentStatus int32

const (
	PaymentStatusUnknown    PaymentStatus = 0
	PaymentStatusPending    PaymentStatus = 1 // 保留中・未支払い
	PaymentStatusAuthorized PaymentStatus = 2 // 仮売上・オーソリ
	PaymentStatusCaptured   PaymentStatus = 3 // 実売上・キャプチャ
	PaymentStatusCanceled   PaymentStatus = 4 // キャンセル
	PaymentStatusRefunded   PaymentStatus = 5 // 返金
	PaymentStatusFailed     PaymentStatus = 6 // 失敗/期限切れ
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
	RefundTypeNone     RefundType = 0 // キャンセルなし
	RefundTypeCanceled RefundType = 1 // キャンセル
	RefundTypeRefunded RefundType = 2 // 返金
)

// OrderPayment - 注文支払い情報
type OrderPayment struct {
	OrderID           string            `gorm:"primaryKey;<-:create"` // 注文履歴ID
	AddressRevisionID int64             `gorm:""`                     // 請求先情報ID
	Status            PaymentStatus     `gorm:""`                     // 決済状況
	TransactionID     string            `gorm:""`                     // 決済ID(決済代行システム)
	PaymentID         string            `gorm:""`                     // 決済ID(決済代行システム)
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
	CanceledAt        time.Time         `gorm:"default:null"`         // 注文キャンセル日時（実売上前）
	RefundedAt        time.Time         `gorm:"default:null"`         // 注文キャンセル日時（実売上後）
	CreatedAt         time.Time         `gorm:"<-:create"`            // 登録日時
	UpdatedAt         time.Time         `gorm:""`                     // 更新日時
}

type OrderPayments []*OrderPayment

type NewOrderPaymentParams struct {
	OrderID    string
	MethodType PaymentMethodType
	Address    *entity.Address
	Baskets    CartBaskets
	Products   Products
	Shipping   *Shipping
	Promotion  *Promotion
}

func NewPaymentStatus(status komoju.PaymentStatus) PaymentStatus {
	switch status {
	case komoju.PaymentStatusPending:
		return PaymentStatusPending
	case komoju.PaymentStatusAuthorized:
		return PaymentStatusAuthorized
	case komoju.PaymentStatusCaptured:
		return PaymentStatusCaptured
	case komoju.PaymentStatusRefunded:
		return PaymentStatusRefunded
	case komoju.PaymentStatusCancelled, komoju.PaymentStatusExpired:
		return PaymentStatusFailed
	default:
		return PaymentStatusUnknown
	}
}

func NewKomojuPaymentTypes(methodType PaymentMethodType) []komoju.PaymentType {
	switch methodType {
	case PaymentMethodTypeCash:
		// 未対応
		return []komoju.PaymentType{}
	case PaymentMethodTypeCreditCard:
		return []komoju.PaymentType{komoju.PaymentTypeCreditCard}
	case PaymentMethodTypeKonbini:
		return []komoju.PaymentType{komoju.PaymentTypeKonbini}
	case PaymentMethodTypeBankTranser:
		return []komoju.PaymentType{komoju.PaymentTypeBankTransfer}
	case PaymentMethodTypePayPay:
		return []komoju.PaymentType{komoju.PaymentTypePayPay}
	case PaymentMethodTypeLinePay:
		return []komoju.PaymentType{komoju.PaymentTypeLinePay}
	case PaymentMethodTypeMerpay:
		return []komoju.PaymentType{komoju.PaymentTypeMerpay}
	case PaymentMethodTypeRakutenPay:
		return []komoju.PaymentType{komoju.PaymentTypeRakutenPay}
	case PaymentMethodTypeAUPay:
		return []komoju.PaymentType{komoju.PaymentTypeAUPay}
	default:
		return []komoju.PaymentType{}
	}
}

func NewOrderPayment(params *NewOrderPaymentParams) (*OrderPayment, error) {
	if params.Address == nil {
		return nil, errNotFoundAddress
	}
	if err := codes.ValidatePrefectureValues(params.Address.PrefectureCode); err != nil {
		return nil, err
	}
	sparams := &NewOrderPaymentSummaryParams{
		PrefectureCode: params.Address.PrefectureCode,
		Baskets:        params.Baskets,
		Products:       params.Products,
		Shipping:       params.Shipping,
		Promotion:      params.Promotion,
	}
	summary, err := NewOrderPaymentSummary(sparams)
	if err != nil {
		return nil, err
	}
	return &OrderPayment{
		OrderID:           params.OrderID,
		AddressRevisionID: params.Address.AddressRevision.ID,
		Status:            PaymentStatusPending,
		TransactionID:     "",
		MethodType:        params.MethodType,
		Subtotal:          summary.Subtotal,
		Discount:          summary.Discount,
		ShippingFee:       summary.ShippingFee,
		Tax:               summary.Tax,
		Total:             summary.Total,
	}, nil
}

func (p *OrderPayment) IsCompleted() bool {
	return p.Status == PaymentStatusCaptured ||
		p.Status == PaymentStatusCanceled ||
		p.Status == PaymentStatusRefunded ||
		p.Status == PaymentStatusFailed
}

func (p *OrderPayment) IsCanceled() bool {
	return p.Status == PaymentStatusCanceled || p.Status == PaymentStatusRefunded
}

func (p *OrderPayment) SetTransactionID(transactionID string, now time.Time) {
	p.TransactionID = transactionID
	p.OrderedAt = now
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
