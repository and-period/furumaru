package entity

import (
	"slices"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/set"
)

// 支払いステータス
type PaymentStatus int32

const (
	PaymentStatusUnknown    PaymentStatus = 0
	PaymentStatusPending    PaymentStatus = 1 // 保留中・未支払い
	PaymentStatusAuthorized PaymentStatus = 2 // 仮売上・オーソリ
	PaymentStatusCaptured   PaymentStatus = 3 // 実売上・キャプチャ
	PaymentStatusCanceled   PaymentStatus = 4 // キャンセル
	PaymentStatusRefunded   PaymentStatus = 5 // 返金
	PaymentStatusFailed     PaymentStatus = 6 // 失敗
	PaymentStatusExpired    PaymentStatus = 7 // 期限切れ
)

var (
	PaymentSuccessStatuses  = []PaymentStatus{PaymentStatusAuthorized, PaymentStatusCaptured}
	PaymentFailedStatuses   = []PaymentStatus{PaymentStatusFailed, PaymentStatusExpired}
	PaymentRefundedStatuses = []PaymentStatus{PaymentStatusCanceled, PaymentStatusRefunded}
)

// PaymentMethodType - 決済手段
type PaymentMethodType int32

const (
	PaymentMethodTypeUnknown      PaymentMethodType = 0
	PaymentMethodTypeCash         PaymentMethodType = 1  // 代引支払い
	PaymentMethodTypeCreditCard   PaymentMethodType = 2  // クレジットカード決済
	PaymentMethodTypeKonbini      PaymentMethodType = 3  // コンビニ決済（セブン-イレブン、ローソン、ファミリーマート、ミニストップ、セイコーマート、デイリーヤマザキ）
	PaymentMethodTypeBankTransfer PaymentMethodType = 4  // 銀行振込決済
	PaymentMethodTypePayPay       PaymentMethodType = 5  // QR決済（PayPay）
	PaymentMethodTypeLinePay      PaymentMethodType = 6  // QR決済（LINE Pay）
	PaymentMethodTypeMerpay       PaymentMethodType = 7  // QR決済（メルペイ）
	PaymentMethodTypeRakutenPay   PaymentMethodType = 8  // QR決済（楽天ペイ）
	PaymentMethodTypeAUPay        PaymentMethodType = 9  // QR決済（au PAY）
	PaymentMethodTypeNone         PaymentMethodType = 10 // 決済なし
	PaymentMethodTypePaidy        PaymentMethodType = 11 // ペイディ（Paidy）
	PaymentMethodTypePayEasy      PaymentMethodType = 12 // ペイジー（Pay-easy）
)

var (
	// 即時決済対応の決済手段
	ImmediatePaymentMethodTypes = []PaymentMethodType{
		PaymentMethodTypeCreditCard,
		PaymentMethodTypePayPay,
		PaymentMethodTypeLinePay,
		PaymentMethodTypeMerpay,
		PaymentMethodTypeRakutenPay,
		PaymentMethodTypeAUPay,
		PaymentMethodTypePaidy,
	}
	// 後日決済対応の決済手段
	DeferredPaymentMethodTypes = []PaymentMethodType{
		PaymentMethodTypeKonbini,
		PaymentMethodTypeBankTransfer,
		PaymentMethodTypePayEasy,
	}
	// その他の決済手段
	OtherPaymentMethodTypes = []PaymentMethodType{
		PaymentMethodTypeCash,
		PaymentMethodTypeNone,
	}
	AllPaymentMethodTypes = slices.Concat(
		ImmediatePaymentMethodTypes,
		DeferredPaymentMethodTypes,
		[]PaymentMethodType{PaymentMethodTypeCash},
	)
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
	AddressRevisionID int64             `gorm:"default:null"`         // 請求先情報ID
	Status            PaymentStatus     `gorm:""`                     // 決済状況
	TransactionID     string            `gorm:""`                     // 決済ID(決済代行システム)
	PaymentID         string            `gorm:""`                     // 決済ID(決済代行システム)
	MethodType        PaymentMethodType `gorm:""`                     // 決済手段種別
	Subtotal          int64             `gorm:""`                     // 購入金額(税込)
	Discount          int64             `gorm:""`                     // 割引金額(税込)
	ShippingFee       int64             `gorm:""`                     // 配送手数料(税込)
	Tax               int64             `gorm:""`                     // 消費税(内税)
	Total             int64             `gorm:""`                     // 合計金額(税込)
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

type NewProductOrderPaymentParams struct {
	OrderID    string
	Pickup     bool
	MethodType PaymentMethodType
	Address    *entity.Address
	Baskets    CartBaskets
	Products   Products
	Shipping   *Shipping
	Promotion  *Promotion
}

type NewExperienceOrderPaymentParams struct {
	OrderID               string
	MethodType            PaymentMethodType
	Address               *entity.Address
	Experience            *Experience
	Promotion             *Promotion
	AdultCount            int64
	JuniorHighSchoolCount int64
	ElementarySchoolCount int64
	PreschoolCount        int64
	SeniorCount           int64
}

func (t PaymentMethodType) String() string {
	switch t {
	case PaymentMethodTypeCash:
		return "代引支払い"
	case PaymentMethodTypeCreditCard:
		return "クレジットカード決済"
	case PaymentMethodTypeKonbini:
		return "コンビニ決済"
	case PaymentMethodTypeBankTransfer:
		return "銀行振込決済"
	case PaymentMethodTypePayPay:
		return "QR決済（PayPay）"
	case PaymentMethodTypeLinePay:
		return "QR決済（LINE Pay）"
	case PaymentMethodTypeMerpay:
		return "QR決済（メルペイ）"
	case PaymentMethodTypeRakutenPay:
		return "QR決済（楽天ペイ）"
	case PaymentMethodTypeAUPay:
		return "QR決済（au PAY）"
	case PaymentMethodTypePaidy:
		return "ペイディ（Paidy）"
	case PaymentMethodTypePayEasy:
		return "ペイジー（Pay-easy）"
	case PaymentMethodTypeNone:
		return "決済なし"
	default:
		return ""
	}
}

func NewProductOrderPayment(params *NewProductOrderPaymentParams) (*OrderPayment, error) {
	var (
		addressRevisionID int64
		prefectureCode    int32
	)
	if params.Address != nil {
		addressRevisionID = params.Address.AddressRevision.ID
		prefectureCode = params.Address.PrefectureCode
		if err := codes.ValidatePrefectureValues(prefectureCode); err != nil {
			return nil, err
		}
	}
	sparams := &NewProductOrderPaymentSummaryParams{
		PrefectureCode: prefectureCode,
		Pickup:         params.Pickup,
		Baskets:        params.Baskets,
		Products:       params.Products,
		Shipping:       params.Shipping,
		Promotion:      params.Promotion,
	}
	summary, err := NewProductOrderPaymentSummary(sparams)
	if err != nil {
		return nil, err
	}
	return &OrderPayment{
		OrderID:           params.OrderID,
		AddressRevisionID: addressRevisionID,
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

func NewExperienceOrderPayment(params *NewExperienceOrderPaymentParams) (*OrderPayment, error) {
	var (
		addressRevisionID int64
		prefectureCode    int32
	)
	if params.Address != nil {
		addressRevisionID = params.Address.AddressRevision.ID
		prefectureCode = params.Address.PrefectureCode
		if err := codes.ValidatePrefectureValues(prefectureCode); err != nil {
			return nil, err
		}
	}
	sparams := &NewExperienceOrderPaymentSummaryParams{
		Experience:            params.Experience,
		Promotion:             params.Promotion,
		AdultCount:            params.AdultCount,
		JuniorHighSchoolCount: params.JuniorHighSchoolCount,
		ElementarySchoolCount: params.ElementarySchoolCount,
		PreschoolCount:        params.PreschoolCount,
		SeniorCount:           params.SeniorCount,
	}
	summary := NewExperienceOrderPaymentSummary(sparams)
	return &OrderPayment{
		OrderID:           params.OrderID,
		AddressRevisionID: addressRevisionID,
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
		p.Status == PaymentStatusFailed ||
		p.Status == PaymentStatusExpired
}

func (p *OrderPayment) IsCanceled() bool {
	return p.Status == PaymentStatusCanceled || p.Status == PaymentStatusRefunded
}

func (p *OrderPayment) IsImmediatePayment() bool {
	return slices.Contains(ImmediatePaymentMethodTypes, p.MethodType)
}

func (p *OrderPayment) IsDeferredPayment() bool {
	return slices.Contains(DeferredPaymentMethodTypes, p.MethodType)
}

func (p *OrderPayment) SetTransactionID(transactionID string, now time.Time) {
	p.TransactionID = transactionID
	p.OrderedAt = now
}

func (ps OrderPayments) AddressRevisionIDs() []int64 {
	res := set.NewEmpty[int64](len(ps))
	for _, p := range ps {
		if p.AddressRevisionID == 0 {
			continue
		}
		res.Add(p.AddressRevisionID)
	}
	return res.Slice()
}

func (ps OrderPayments) MapByOrderID() map[string]*OrderPayment {
	res := make(map[string]*OrderPayment, len(ps))
	for _, p := range ps {
		res[p.OrderID] = p
	}
	return res
}
