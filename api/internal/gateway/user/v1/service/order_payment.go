package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

// PaymentMethodType - 決済手段
type PaymentMethodType int32

const (
	PaymentMethodTypeUnknown      PaymentMethodType = 0
	PaymentMethodTypeCash         PaymentMethodType = 1  // 代引支払い
	PaymentMethodTypeCreditCard   PaymentMethodType = 2  // クレジットカード決済
	PaymentMethodTypeKonbini      PaymentMethodType = 3  // コンビニ決済
	PaymentMethodTypeBankTransfer PaymentMethodType = 4  // 銀行振込決済
	PaymentMethodTypePayPay       PaymentMethodType = 5  // QR決済（PayPay）
	PaymentMethodTypeLinePay      PaymentMethodType = 6  // QR決済（LINE Pay）
	PaymentMethodTypeMerpay       PaymentMethodType = 7  // QR決済（メルペイ）
	PaymentMethodTypeRakutenPay   PaymentMethodType = 8  // QR決済（楽天ペイ）
	PaymentMethodTypeAUPay        PaymentMethodType = 9  // QR決済（au PAY）
	PaymentMethodTypeFree         PaymentMethodType = 10 // 決済無し
	PaymentMethodTypePaidy        PaymentMethodType = 11 // ペイディ（Paidy）
	PaymentMethodTypePayEasy      PaymentMethodType = 12 // ペイジー（Pay-easy）
)

// PaymentStatus - 支払い状況
type PaymentStatus int32

const (
	PaymentStatusUnknown  PaymentStatus = 0
	PaymentStatusUnpaid   PaymentStatus = 1 // 未支払い
	PaymentStatusPaid     PaymentStatus = 2 // 支払い済み
	PaymentStatusCanceled PaymentStatus = 3 // キャンセル済み
	PaymentStatusFailed   PaymentStatus = 4 // 失敗
)

type OrderPayment struct {
	response.OrderPayment
}

type OrderPayments []*OrderPayment

func NewPaymentMethodType(typ entity.PaymentMethodType) PaymentMethodType {
	switch typ {
	case entity.PaymentMethodTypeCash:
		return PaymentMethodTypeCash
	case entity.PaymentMethodTypeCreditCard:
		return PaymentMethodTypeCreditCard
	case entity.PaymentMethodTypeKonbini:
		return PaymentMethodTypeKonbini
	case entity.PaymentMethodTypeBankTransfer:
		return PaymentMethodTypeBankTransfer
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
	case entity.PaymentMethodTypePaidy:
		return PaymentMethodTypePaidy
	case entity.PaymentMethodTypePayEasy:
		return PaymentMethodTypePayEasy
	case entity.PaymentMethodTypeNone:
		return PaymentMethodTypeFree
	default:
		return PaymentMethodTypeUnknown
	}
}

func (t PaymentMethodType) StoreEntity() entity.PaymentMethodType {
	switch t {
	case PaymentMethodTypeCash:
		return entity.PaymentMethodTypeCash
	case PaymentMethodTypeCreditCard:
		return entity.PaymentMethodTypeCreditCard
	case PaymentMethodTypeKonbini:
		return entity.PaymentMethodTypeKonbini
	case PaymentMethodTypeBankTransfer:
		return entity.PaymentMethodTypeBankTransfer
	case PaymentMethodTypePayPay:
		return entity.PaymentMethodTypePayPay
	case PaymentMethodTypeLinePay:
		return entity.PaymentMethodTypeLinePay
	case PaymentMethodTypeMerpay:
		return entity.PaymentMethodTypeMerpay
	case PaymentMethodTypeRakutenPay:
		return entity.PaymentMethodTypeRakutenPay
	case PaymentMethodTypeAUPay:
		return entity.PaymentMethodTypeAUPay
	case PaymentMethodTypePaidy:
		return entity.PaymentMethodTypePaidy
	case PaymentMethodTypePayEasy:
		return entity.PaymentMethodTypePayEasy
	case PaymentMethodTypeFree:
		return entity.PaymentMethodTypeNone
	default:
		return entity.PaymentMethodTypeUnknown
	}
}

func (t PaymentMethodType) Response() int32 {
	return int32(t)
}

func NewPaymentStatus(status entity.PaymentStatus) PaymentStatus {
	switch status {
	case entity.PaymentStatusPending:
		return PaymentStatusUnpaid
	case entity.PaymentStatusAuthorized, entity.PaymentStatusCaptured:
		return PaymentStatusPaid
	case entity.PaymentStatusCanceled, entity.PaymentStatusRefunded:
		return PaymentStatusCanceled
	case entity.PaymentStatusFailed, entity.PaymentStatusExpired:
		return PaymentStatusFailed
	default:
		return PaymentStatusUnknown
	}
}

func (s PaymentStatus) Response() int32 {
	return int32(s)
}

func NewOrderPayment(payment *entity.OrderPayment) *OrderPayment {
	return &OrderPayment{
		OrderPayment: response.OrderPayment{
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

func (p *OrderPayment) Response() *response.OrderPayment {
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
