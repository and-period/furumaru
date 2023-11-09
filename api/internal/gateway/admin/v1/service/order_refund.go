package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

// RefundType - 注文キャンセル種別
type RefundType int32

const (
	RefundTypeNone RefundType = 0 // キャンセルなし
)

type OrderRefund struct {
	response.OrderRefund
	orderID string
}

type OrderRefunds []*OrderRefund

func NewRefundType(_ entity.RefundType) RefundType {
	return RefundTypeNone
}

func (t RefundType) Response() int32 {
	return int32(t)
}

func NewOrderRefund(payment *entity.OrderPayment) *OrderRefund {
	return &OrderRefund{
		OrderRefund: response.OrderRefund{
			Total:      payment.RefundTotal,
			Type:       NewRefundType(payment.RefundType).Response(),
			Reason:     payment.RefundReason,
			Canceled:   payment.IsCanceled(),
			CanceledAt: jst.Unix(payment.RefundedAt),
		},
		orderID: payment.OrderID,
	}
}

func (r *OrderRefund) Response() *response.OrderRefund {
	if r == nil {
		return nil
	}
	return &r.OrderRefund
}

func NewOrderRefunds(payments entity.OrderPayments) OrderRefunds {
	res := make(OrderRefunds, len(payments))
	for i := range payments {
		res[i] = NewOrderRefund(payments[i])
	}
	return res
}
