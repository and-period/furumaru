package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

// PaymentMethodType - 決済手段
type PaymentMethodType int32

const (
	PaymentMethodTypeUnknown PaymentMethodType = 0
	PaymentMethodTypeCash    PaymentMethodType = 1 // 代引支払い
	PaymentMethodTypeCard    PaymentMethodType = 2 // クレジットカード払い
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
	case entity.PaymentMethodTypeCard:
		return PaymentMethodTypeCard
	default:
		return PaymentMethodTypeUnknown
	}
}

func (t PaymentMethodType) Response() int32 {
	return int32(t)
}

func NewPaymentStatus(status entity.PaymentStatus) PaymentStatus {
	switch status {
	case entity.PaymentStatusInitialized:
		return PaymentStatusUnpaid
	case entity.PaymentStatusPending:
		return PaymentStatusPending
	case entity.PaymentStatusAuthorized:
		return PaymentStatusAuthorized
	case entity.PaymentStatusCaptured:
		return PaymentStatusPaid
	case entity.PaymentStatusCanceled:
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
			MethodID:      payment.MethodID,
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
