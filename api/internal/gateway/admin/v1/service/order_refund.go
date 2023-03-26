package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

// OrderRefundType - 注文キャンセル理由
type OrderRefundType int32

const (
	OrderRefundTypeUnknown OrderRefundType = 0
)

type OrderRefund struct {
	response.OrderRefund
}

func NewOrderRefundType(_ entity.CancelType) OrderRefundType {
	// TODO: 必要な種別が決まったら詳細実装
	return OrderRefundTypeUnknown
}

func (t OrderRefundType) IsCanceled() bool {
	return t != OrderRefundTypeUnknown
}

func (t OrderRefundType) Response() int32 {
	return int32(t)
}

func NewOrderRefund(order *entity.Order) *OrderRefund {
	cancelType := NewOrderRefundType(order.CancelType)
	return &OrderRefund{
		OrderRefund: response.OrderRefund{
			Canceled: cancelType.IsCanceled(),
			Type:     cancelType.Response(),
			Reason:   order.CancelReason,
			Total:    order.RefundTotal,
		},
	}
}

func (r *OrderRefund) Response() *response.OrderRefund {
	return &r.OrderRefund
}
