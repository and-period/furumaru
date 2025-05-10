package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

type UserOrder struct {
	response.UserOrder
}

type UserOrders []*UserOrder

func NewUserOrder(order *entity.Order) *UserOrder {
	return &UserOrder{
		UserOrder: response.UserOrder{
			OrderID:   order.ID,
			Status:    NewPaymentStatus(order.OrderPayment.Status).Response(),
			SubTotal:  order.Subtotal,
			Total:     order.Total,
			OrderedAt: jst.Unix(order.OrderedAt),
			PaidAt:    jst.Unix(order.PaidAt),
		},
	}
}

func (o *UserOrder) Response() *response.UserOrder {
	return &o.UserOrder
}

func NewUserOrders(orders entity.Orders) UserOrders {
	res := make(UserOrders, len(orders))
	for i := range orders {
		res[i] = NewUserOrder(orders[i])
	}
	return res
}

func (os UserOrders) Response() []*response.UserOrder {
	res := make([]*response.UserOrder, len(os))
	for i := range os {
		res[i] = os[i].Response()
	}
	return res
}
