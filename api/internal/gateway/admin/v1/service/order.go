package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
)

type Order struct {
	response.Order
}

type Orders []*Order

func NewOrder(order *entity.Order, addresses map[int64]*Address, products map[int64]*Product) *Order {
	return &Order{
		Order: response.Order{
			ID:            order.ID,
			UserID:        order.UserID,
			CoordinatorID: order.CoordinatorID,
			PromotionID:   order.PromotionID,
			Payment:       NewOrderPayment(&order.OrderPayment, addresses[order.OrderPayment.AddressRevisionID]).Response(),
			Refund:        NewOrderRefund(&order.OrderPayment).Response(),
			Fulfillments:  NewOrderFulfillments(order.OrderFulfillments, addresses).Response(),
			Items:         NewOrderItems(order.OrderItems, products).Response(),
			CreatedAt:     jst.Unix(order.CreatedAt),
			UpdatedAt:     jst.Unix(order.UpdatedAt),
		},
	}
}

func (o *Order) ProductIDs() []string {
	return set.UniqBy(o.Items, func(i *response.OrderItem) string {
		return i.ProductID
	})
}

func (o *Order) Response() *response.Order {
	return &o.Order
}

func NewOrders(orders entity.Orders, addresses map[int64]*Address, products map[int64]*Product) Orders {
	res := make(Orders, len(orders))
	for i := range orders {
		res[i] = NewOrder(orders[i], addresses, products)
	}
	return res
}

func (os Orders) Response() []*response.Order {
	res := make([]*response.Order, len(os))
	for i := range os {
		res[i] = os[i].Response()
	}
	return res
}
