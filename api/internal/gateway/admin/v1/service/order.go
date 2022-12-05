package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
)

type Order struct {
	response.Order
	payment       *OrderPayment
	fulfillment   *OrderFulfillment
	refund        *OrderRefund
	items         OrderItems
	CoordinatorID string
}

type Orders []*Order

func NewOrder(order *entity.Order) *Order {
	return &Order{
		Order: response.Order{
			ID:          order.ID,
			UserID:      order.UserID,
			ScheduleID:  order.ScheduleID,
			OrderedAt:   jst.Unix(order.OrderedAt),
			PaidAt:      jst.Unix(order.ConfirmedAt),
			DeliveredAt: jst.Unix(order.DeliveredAt),
			CanceledAt:  jst.Unix(order.CanceledAt),
			CreatedAt:   order.CreatedAt.Unix(),
			UpdatedAt:   order.UpdatedAt.Unix(),
		},
		payment:       NewOrderPayment(&order.OrderPayment, order.PaymentStatus),
		fulfillment:   NewOrderFulfillment(&order.OrderFulfillment, order.FulfillmentStatus),
		refund:        NewOrderRefund(order),
		items:         NewOrderItems(order.OrderItems),
		CoordinatorID: order.CoordinatorID,
	}
}

func (o *Order) Fill(user *User, products map[string]*Product) {
	if user != nil {
		o.UserName = user.Name()
	}
	if o.items != nil {
		o.items.Fill(products)
	}
}

func (o *Order) ProductIDs() []string {
	return o.items.ProductIDs()
}

func (o *Order) Response() *response.Order {
	o.Payment = o.payment.Response()
	o.Fulfillment = o.fulfillment.Response()
	o.Refund = o.refund.Response()
	o.Items = o.items.Response()
	return &o.Order
}

func NewOrders(orders entity.Orders) Orders {
	res := make(Orders, len(orders))
	for i := range orders {
		res[i] = NewOrder(orders[i])
	}
	return res
}

func (os Orders) Fill(users map[string]*User, products map[string]*Product) {
	for i := range os {
		os[i].Fill(users[os[i].UserID], products)
	}
}

func (os Orders) UserIDs() []string {
	return set.UniqBy(os, func(o *Order) string {
		return o.UserID
	})
}

func (os Orders) ProductIDs() []string {
	res := set.NewEmpty[string](len(os))
	for i := range os {
		res.Add(os[i].ProductIDs()...)
	}
	return res.Slice()
}

func (os Orders) Response() []*response.Order {
	res := make([]*response.Order, len(os))
	for i := range os {
		res[i] = os[i].Response()
	}
	return res
}
