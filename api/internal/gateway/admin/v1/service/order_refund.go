package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

// RefundType - 注文キャンセル種別
type RefundType int32

const (
	RefundTypeNone     RefundType = 0 // キャンセルなし
	RefundTypeCanceled RefundType = 1 // キャンセル
	RefundTypeRefunded RefundType = 2 // 返金
)

type OrderRefund struct {
	types.OrderRefund
	orderID string
}

type OrderRefunds []*OrderRefund

func NewRefundType(typ entity.RefundType) RefundType {
	switch typ {
	case entity.RefundTypeCanceled:
		return RefundTypeCanceled
	case entity.RefundTypeRefunded:
		return RefundTypeRefunded
	default:
		return RefundTypeNone
	}
}

func (t RefundType) Response() int32 {
	return int32(t)
}

func NewOrderRefund(payment *entity.OrderPayment) *OrderRefund {
	return &OrderRefund{
		OrderRefund: types.OrderRefund{
			Total:      payment.RefundTotal,
			Type:       NewRefundType(payment.RefundType).Response(),
			Reason:     payment.RefundReason,
			Canceled:   payment.IsCanceled(),
			CanceledAt: newOrderCanceledAt(payment),
		},
		orderID: payment.OrderID,
	}
}

func newOrderCanceledAt(payment *entity.OrderPayment) int64 {
	switch payment.RefundType {
	case entity.RefundTypeCanceled:
		return jst.Unix(payment.CanceledAt)
	case entity.RefundTypeRefunded:
		return jst.Unix(payment.RefundedAt)
	default:
		return 0
	}
}

func (r *OrderRefund) Response() *types.OrderRefund {
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
