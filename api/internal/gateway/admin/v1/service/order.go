package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
)

// OrderStatus - 注文ステータス
type OrderStatus int32

const (
	OrderStatusUnknown   OrderStatus = 0
	OrderStatusUnpaid    OrderStatus = 1 // 支払い待ち
	OrderStatusWaiting   OrderStatus = 2 // 受注待ち
	OrderStatusPreparing OrderStatus = 3 // 発送準備中
	OrderStatusShipped   OrderStatus = 4 // 発送完了
	OrderStatusCompleted OrderStatus = 5 // 完了
	OrderStatusCanceled  OrderStatus = 6 // キャンセル
	OrderStatusRefunded  OrderStatus = 7 // 返金
	OrderStatusFailed    OrderStatus = 8 // 失敗
)

type Order struct {
	response.Order
}

type Orders []*Order

func NewOrderStatus(status entity.OrderStatus) OrderStatus {
	switch status {
	case entity.OrderStatusUnpaid:
		return OrderStatusUnpaid
	case entity.OrderStatusWaiting:
		return OrderStatusWaiting
	case entity.OrderStatusPreparing:
		return OrderStatusPreparing
	case entity.OrderStatusShipped:
		return OrderStatusShipped
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
	return &Order{
		Order: response.Order{
			ID:              order.ID,
			UserID:          order.UserID,
			CoordinatorID:   order.CoordinatorID,
			PromotionID:     order.PromotionID,
			ManagementID:    order.ManagementID,
			ShippingMessage: order.ShippingMessage,
			Status:          NewOrderStatus(order.Status).Response(),
			Payment:         NewOrderPayment(&order.OrderPayment, addresses[order.OrderPayment.AddressRevisionID]).Response(),
			Refund:          NewOrderRefund(&order.OrderPayment).Response(),
			Fulfillments:    NewOrderFulfillments(order.OrderFulfillments, addresses).Response(),
			Items:           NewOrderItems(order.OrderItems, products).Response(),
			CreatedAt:       jst.Unix(order.CreatedAt),
			UpdatedAt:       jst.Unix(order.UpdatedAt),
			CompletedAt:     jst.Unix(order.CompletedAt),
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
