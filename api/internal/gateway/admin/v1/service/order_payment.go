package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
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

// PaymentStatus - 支払い状況
type PaymentStatus int32

const (
	PaymentStatusUnknown    PaymentStatus = 0
	PaymentStatusUnpaid     PaymentStatus = 1 // 未払い
	PaymentStatusPending    PaymentStatus = 2 // 保留中
	PaymentStatusAuthorized PaymentStatus = 3 // オーソリ済み
	PaymentStatusPaid       PaymentStatus = 4 // 支払い済み
	PaymentStatusRefunded   PaymentStatus = 5 // 返金済み
	PaymentStatusExpired    PaymentStatus = 6 // 期限切れ
)

type OrderPayment struct {
	response.OrderPayment
	orderID string
}

func NewPaymentMethodType(typ entity.PaymentMethodType) PaymentMethodType {
	switch typ {
	case entity.PaymentMethodTypeCash:
		return PaymentMethodTypeCash
	case entity.PaymentMethodTypeCreditCard:
		return PaymentMethodTypeCreditCard
	case entity.PaymentMethodTypeKonbini:
		return PaymentMethodTypeKonbini
	case entity.PaymentMethodTypeBankTranser:
		return PaymentMethodTypeBankTranser
	case entity.PaymentMethodTypePayPay:
		return PaymentMethodTypePayPay
	case entity.PaymentMethodTypeLinePay:
		return PaymentMethodTypeLinePay
	case entity.PaymentMethodTypeMerpay:
		return PaymentMethodTypeMerpay
	case entity.PaymentMethodTypeRakutenPay:
		return PaymentMethodTypeRakutenPay
	case entity.PaymentMethodTypeAUPay:
		return PaymentMethodTypeAUPay
	default:
		return PaymentMethodTypeUnknown
	}
}

func (t PaymentMethodType) Response() int32 {
	return int32(t)
}

func NewPaymentStatus(status entity.PaymentStatus) PaymentStatus {
	switch status {
	case entity.PaymentStatusPending:
		return PaymentStatusPending
	case entity.PaymentStatusAuthorized:
		return PaymentStatusAuthorized
	case entity.PaymentStatusCaptured:
		return PaymentStatusPaid
	case entity.PaymentStatusRefunded:
		return PaymentStatusRefunded
	case entity.PaymentStatusFailed:
		return PaymentStatusExpired
	default:
		return PaymentStatusUnknown
	}
}

func (s PaymentStatus) Response() int32 {
	return int32(s)
}

func NewOrderPayment(payment *entity.Payment, status entity.PaymentStatus) *OrderPayment {
	return &OrderPayment{
		OrderPayment: response.OrderPayment{
			TransactionID: payment.TransactionID,
			MethodType:    NewPaymentMethodType(payment.MethodType).Response(),
			Status:        NewPaymentStatus(status).Response(),
			Subtotal:      payment.Subtotal,
			Discount:      payment.Discount,
			ShippingFee:   payment.ShippingFee,
			Tax:           payment.Tax,
			Total:         payment.Total,
			AddressID:     payment.AddressID,
		},
		orderID: payment.OrderID,
	}
}

func (p *OrderPayment) Fill(address *Address) {
	p.Address = address.Response()
}

func (p *OrderPayment) Response() *response.OrderPayment {
	return &p.OrderPayment
}
