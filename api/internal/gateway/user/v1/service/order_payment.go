package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

// PaymentMethodType - 決済手段
type PaymentMethodType types.PaymentMethodType

// PaymentStatus - 支払い状況
type PaymentStatus types.PaymentStatus

type OrderPayment struct {
	types.OrderPayment
}

type OrderPayments []*OrderPayment

func NewPaymentMethodType(typ entity.PaymentMethodType) PaymentMethodType {
	switch typ {
	case entity.PaymentMethodTypeCash:
		return PaymentMethodType(types.PaymentMethodTypeCash)
	case entity.PaymentMethodTypeCreditCard:
		return PaymentMethodType(types.PaymentMethodTypeCreditCard)
	case entity.PaymentMethodTypeKonbini:
		return PaymentMethodType(types.PaymentMethodTypeKonbini)
	case entity.PaymentMethodTypeBankTransfer:
		return PaymentMethodType(types.PaymentMethodTypeBankTransfer)
	case entity.PaymentMethodTypePayPay:
		return PaymentMethodType(types.PaymentMethodTypePayPay)
	case entity.PaymentMethodTypeLinePay:
		return PaymentMethodType(types.PaymentMethodTypeLinePay)
	case entity.PaymentMethodTypeMerpay:
		return PaymentMethodType(types.PaymentMethodTypeMerpay)
	case entity.PaymentMethodTypeRakutenPay:
		return PaymentMethodType(types.PaymentMethodTypeRakutenPay)
	case entity.PaymentMethodTypeAUPay:
		return PaymentMethodType(types.PaymentMethodTypeAUPay)
	case entity.PaymentMethodTypePaidy:
		return PaymentMethodType(types.PaymentMethodTypePaidy)
	case entity.PaymentMethodTypePayEasy:
		return PaymentMethodType(types.PaymentMethodTypePayEasy)
	case entity.PaymentMethodTypeNone:
		return PaymentMethodType(types.PaymentMethodTypeFree)
	default:
		return PaymentMethodType(types.PaymentMethodTypeUnknown)
	}
}

func (t PaymentMethodType) StoreEntity() entity.PaymentMethodType {
	switch types.PaymentMethodType(t) {
	case types.PaymentMethodTypeCash:
		return entity.PaymentMethodTypeCash
	case types.PaymentMethodTypeCreditCard:
		return entity.PaymentMethodTypeCreditCard
	case types.PaymentMethodTypeKonbini:
		return entity.PaymentMethodTypeKonbini
	case types.PaymentMethodTypeBankTransfer:
		return entity.PaymentMethodTypeBankTransfer
	case types.PaymentMethodTypePayPay:
		return entity.PaymentMethodTypePayPay
	case types.PaymentMethodTypeLinePay:
		return entity.PaymentMethodTypeLinePay
	case types.PaymentMethodTypeMerpay:
		return entity.PaymentMethodTypeMerpay
	case types.PaymentMethodTypeRakutenPay:
		return entity.PaymentMethodTypeRakutenPay
	case types.PaymentMethodTypeAUPay:
		return entity.PaymentMethodTypeAUPay
	case types.PaymentMethodTypePaidy:
		return entity.PaymentMethodTypePaidy
	case types.PaymentMethodTypePayEasy:
		return entity.PaymentMethodTypePayEasy
	case types.PaymentMethodTypeFree:
		return entity.PaymentMethodTypeNone
	default:
		return entity.PaymentMethodTypeUnknown
	}
}

func (t PaymentMethodType) Response() types.PaymentMethodType {
	return types.PaymentMethodType(t)
}

func NewPaymentStatus(status entity.PaymentStatus) PaymentStatus {
	switch status {
	case entity.PaymentStatusPending:
		return PaymentStatus(types.PaymentStatusUnpaid)
	case entity.PaymentStatusAuthorized, entity.PaymentStatusCaptured:
		return PaymentStatus(types.PaymentStatusPaid)
	case entity.PaymentStatusCanceled, entity.PaymentStatusRefunded:
		return PaymentStatus(types.PaymentStatusCanceled)
	case entity.PaymentStatusFailed, entity.PaymentStatusExpired:
		return PaymentStatus(types.PaymentStatusFailed)
	default:
		return PaymentStatus(types.PaymentStatusUnknown)
	}
}

func (s PaymentStatus) Response() types.PaymentStatus {
	return types.PaymentStatus(s)
}

func NewOrderPayment(payment *entity.OrderPayment) *OrderPayment {
	return &OrderPayment{
		OrderPayment: types.OrderPayment{
			TransactionID: payment.TransactionID,
			MethodType:    NewPaymentMethodType(payment.MethodType).Response(),
			Status:        NewPaymentStatus(payment.Status).Response(),
			Subtotal:      payment.Subtotal,
			Discount:      payment.Discount,
			ShippingFee:   payment.ShippingFee,
			Total:         payment.Total,
			OrderedAt:     jst.Unix(payment.OrderedAt),
			PaidAt:        jst.Unix(payment.PaidAt),
		},
	}
}

func (p *OrderPayment) Response() *types.OrderPayment {
	if p == nil {
		return nil
	}
	return &p.OrderPayment
}

func NewOrderPayments(payments entity.OrderPayments) OrderPayments {
	res := make(OrderPayments, len(payments))
	for i := range payments {
		res[i] = NewOrderPayment(payments[i])
	}
	return res
}
