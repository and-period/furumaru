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
			PromotionID: order.PromotionID,
			OrderedAt:   jst.Unix(order.OrderedAt),
			PaidAt:      jst.Unix(order.PaidAt),
			DeliveredAt: jst.Unix(order.ShippedAt),
			CanceledAt:  jst.Unix(order.RefundedAt),
			CreatedAt:   order.CreatedAt.Unix(),
			UpdatedAt:   order.UpdatedAt.Unix(),
		},
		payment:       NewOrderPayment(&order.Payment, order.PaymentStatus),
		fulfillment:   NewOrderFulfillment(&order.Fulfillment, order.FulfillmentStatus),
		refund:        NewOrderRefund(order),
		items:         NewOrderItems(order.OrderItems),
		CoordinatorID: order.CoordinatorID,
	}
}

func (o *Order) Fill(user *User, products map[string]*Product, addresses map[string]*Address) {
	if user != nil {
		o.UserName = user.Name()
	}
	if o.items != nil {
		o.items.Fill(products)
	}
	if a, ok := addresses[o.payment.AddressID]; ok {
		o.payment.Fill(a)
	}
	if a, ok := addresses[o.fulfillment.AddressID]; ok {
		o.fulfillment.Fill(a)
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

func (os Orders) Fill(users map[string]*User, products map[string]*Product, addresses map[string]*Address) {
	for i := range os {
		os[i].Fill(users[os[i].UserID], products, addresses)
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

func (os Orders) AddressIDs() []string {
	res := set.NewEmpty[string](len(os) * 2) // payment + fulfillment
	for i := range os {
		res.Add(os[i].payment.AddressID, os[i].fulfillment.AddressID)
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
