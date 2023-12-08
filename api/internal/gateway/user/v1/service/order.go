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

func NewOrderStatus(order *entity.Order) OrderStatus {
	if order == nil {
		return OrderStatusUnknown
	}
	switch order.OrderPayment.Status {
	case entity.PaymentStatusPending:
		return OrderStatusUnpaid
	case entity.PaymentStatusAuthorized, entity.PaymentStatusCaptured:
		if order.CompletedAt.IsZero() {
			return OrderStatusPreparing
		}
		return OrderStatusCompleted
	case entity.PaymentStatusCanceled:
		return OrderStatusCanceled
	case entity.PaymentStatusRefunded:
		return OrderStatusRefunded
	case entity.PaymentStatusFailed:
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
			ID:            order.ID,
			CoordinatorID: order.CoordinatorID,
			PromotionID:   order.PromotionID,
			Status:        NewOrderStatus(order).Response(),
			Payment:       NewOrderPayment(&order.OrderPayment, addresses[order.OrderPayment.AddressRevisionID]).Response(),
			Refund:        NewOrderRefund(&order.OrderPayment).Response(),
			Fulfillments:  NewOrderFulfillments(order.OrderFulfillments, addresses).Response(),
			Items:         NewOrderItems(order.OrderItems, products).Response(),
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
