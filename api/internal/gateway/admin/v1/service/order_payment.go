package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

// PaymentType - 決済手段
type PaymentType int32

const (
	PaymentTypeUnknown PaymentType = 0
	PaymentTypeCash    PaymentType = 1 // 代引支払い
	PaymentTypeCard    PaymentType = 2 // クレジットカード払い
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
	id      string
	orderID string
}

func NewPaymentType(typ entity.PaymentType) PaymentType {
	switch typ {
	case entity.PaymentTypeCash:
		return PaymentTypeCash
	case entity.PaymentTypeCard:
		return PaymentTypeCard
	default:
		return PaymentTypeUnknown
	}
}

func (t PaymentType) Response() int32 {
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

func NewOrderPayment(payment *entity.OrderPayment, status entity.PaymentStatus) *OrderPayment {
	return &OrderPayment{
		OrderPayment: response.OrderPayment{
			TransactionID:  payment.TransactionID,
			PromotionID:    payment.PromotionID,
			PaymentID:      payment.PaymentID,
			PaymentType:    NewPaymentType(payment.PaymentType).Response(),
			Status:         NewPaymentStatus(status).Response(),
			Subtotal:       payment.Subtotal,
			Discount:       payment.Discount,
			ShippingCharge: payment.ShippingCharge,
			Tax:            payment.Tax,
			Total:          payment.Total,
			Lastname:       payment.Lastname,
			Firstname:      payment.Firstname,
			PostalCode:     payment.PostalCode,
			Prefecture:     payment.Prefecture,
			City:           payment.City,
			AddressLine1:   payment.AddressLine1,
			AddressLine2:   payment.AddressLine2,
			PhoneNumber:    payment.PhoneNumber,
		},
		id:      payment.ID,
		orderID: payment.OrderID,
	}
}

func (p *OrderPayment) Response() *response.OrderPayment {
	return &p.OrderPayment
}
