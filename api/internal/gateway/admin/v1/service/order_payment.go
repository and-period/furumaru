package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
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
	PaymentStatusUnpaid     PaymentStatus = 1 // 未支払い
	PaymentStatusAuthorized PaymentStatus = 2 // オーソリ済み
	PaymentStatusPaid       PaymentStatus = 3 // 支払い済み
	PaymentStatusCanceled   PaymentStatus = 4 // キャンセル済み
	PaymentStatusFailed     PaymentStatus = 5 // 失敗
)

type OrderPayment struct {
	response.OrderPayment
	orderID string
}

type OrderPayments []*OrderPayment

func NewPaymentMethodType(typ sentity.PaymentMethodType) PaymentMethodType {
	switch typ {
	case sentity.PaymentMethodTypeCash:
		return PaymentMethodTypeCash
	case sentity.PaymentMethodTypeCreditCard:
		return PaymentMethodTypeCreditCard
	case sentity.PaymentMethodTypeKonbini:
		return PaymentMethodTypeKonbini
	case sentity.PaymentMethodTypeBankTranser:
		return PaymentMethodTypeBankTranser
	case sentity.PaymentMethodTypePayPay:
		return PaymentMethodTypePayPay
	case sentity.PaymentMethodTypeLinePay:
		return PaymentMethodTypeLinePay
	case sentity.PaymentMethodTypeMerpay:
		return PaymentMethodTypeMerpay
	case sentity.PaymentMethodTypeRakutenPay:
		return PaymentMethodTypeRakutenPay
	case sentity.PaymentMethodTypeAUPay:
		return PaymentMethodTypeAUPay
	default:
		return PaymentMethodTypeUnknown
	}
}

func (t PaymentMethodType) Response() int32 {
	return int32(t)
}

func NewPaymentStatus(status sentity.PaymentStatus) PaymentStatus {
	switch status {
	case sentity.PaymentStatusPending:
		return PaymentStatusUnpaid
	case sentity.PaymentStatusAuthorized:
		return PaymentStatusAuthorized
	case sentity.PaymentStatusCaptured:
		return PaymentStatusPaid
	case sentity.PaymentStatusCanceled, sentity.PaymentStatusRefunded:
		return PaymentStatusCanceled
	case sentity.PaymentStatusFailed:
		return PaymentStatusFailed
	default:
		return PaymentStatusUnknown
	}
}

func (s PaymentStatus) Response() int32 {
	return int32(s)
}

func NewOrderPayment(payment *sentity.OrderPayment, address *Address) *OrderPayment {
	return &OrderPayment{
		OrderPayment: response.OrderPayment{
			TransactionID: payment.TransactionID,
			MethodType:    NewPaymentMethodType(payment.MethodType).Response(),
			Status:        NewPaymentStatus(payment.Status).Response(),
			Subtotal:      payment.Subtotal,
			Discount:      payment.Discount,
			ShippingFee:   payment.ShippingFee,
			Tax:           payment.Tax,
			Total:         payment.Total,
			OrderedAt:     jst.Unix(payment.OrderedAt),
			PaidAt:        jst.Unix(payment.PaidAt),
			Address:       address.Response(),
		},
		orderID: payment.OrderID,
	}
}

func (p *OrderPayment) Response() *response.OrderPayment {
	if p == nil {
		return nil
	}
	return &p.OrderPayment
}

func NewOrderPayments(payments sentity.OrderPayments, addresses map[int64]*Address) OrderPayments {
	res := make(OrderPayments, len(payments))
	for i, p := range payments {
		res[i] = NewOrderPayment(p, addresses[p.AddressRevisionID])
	}
	return res
}
