package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type OrderRefund struct {
	response.OrderRefund
}

func NewOrderRefund(order *entity.Order) *OrderRefund {
	return &OrderRefund{
		OrderRefund: response.OrderRefund{
			Canceled: order.IsCanceled(),
			Reason:   order.RefundReason,
			Total:    order.RefundTotal,
		},
	}
}

func (r *OrderRefund) Response() *response.OrderRefund {
	return &r.OrderRefund
}
