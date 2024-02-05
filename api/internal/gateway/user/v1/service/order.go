package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/set"
)

// OrderStatus - 注文ステータス
type OrderStatus int32

const (
	OrderStatusUnknown   OrderStatus = 0
	OrderStatusUnpaid    OrderStatus = 1 // 支払い待ち
	OrderStatusPreparing OrderStatus = 2 // 発送対応中
	OrderStatusCompleted OrderStatus = 3 // 完了
	OrderStatusCanceled  OrderStatus = 4 // キャンセル
	OrderStatusRefunded  OrderStatus = 5 // 返金
	OrderStatusFailed    OrderStatus = 6 // 失敗
)

type Order struct {
	response.Order
}

type Orders []*Order

func NewOrderStatus(status entity.OrderStatus) OrderStatus {
	switch status {
	case entity.OrderStatusUnpaid:
		return OrderStatusUnpaid
	case entity.OrderStatusWaiting, entity.OrderStatusPreparing, entity.OrderStatusShipped:
		return OrderStatusPreparing
	case entity.OrderStatusCompleted:
		return OrderStatusCompleted
	case entity.OrderStatusCanceled:
		return OrderStatusCanceled
	case entity.OrderStatusRefunded:
		return OrderStatusRefunded
	case entity.OrderStatusFailed:
		return OrderStatusFailed
	default:
		return OrderStatusUnknown
	}
}

func (s OrderStatus) Response() int32 {
	return int32(s)
}

func NewOrder(order *entity.Order, addresses map[int64]*Address, products map[int64]*Product) *Order {
	var billingAddress, shippingAddress *Address
	if address, ok := addresses[order.OrderPayment.AddressRevisionID]; ok {
		billingAddress = address
	}
	if len(order.OrderFulfillments) > 0 {
		// 現状すべての配送先が同一になっているため
		if address, ok := addresses[order.OrderFulfillments[0].AddressRevisionID]; ok {
			shippingAddress = address
		}
	}
	return &Order{
		Order: response.Order{
			ID:              order.ID,
			CoordinatorID:   order.CoordinatorID,
			PromotionID:     order.PromotionID,
			Status:          NewOrderStatus(order.Status).Response(),
			Payment:         NewOrderPayment(&order.OrderPayment).Response(),
			Refund:          NewOrderRefund(&order.OrderPayment).Response(),
			Fulfillments:    NewOrderFulfillments(order.OrderFulfillments).Response(),
			Items:           NewOrderItems(order.OrderItems, products).Response(),
			BillingAddress:  billingAddress.Response(),
			ShippingAddress: shippingAddress.Response(),
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
