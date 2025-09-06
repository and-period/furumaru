package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
)

// OrderType - 注文種別
type OrderType int32

const (
	OrderTypeUnknown    OrderType = 0
	OrderTypeProduct    OrderType = 1 // 商品
	OrderTypeExperience OrderType = 2 // 体験
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

func NewOrderType(typ entity.OrderType) OrderType {
	switch typ {
	case entity.OrderTypeProduct:
		return OrderTypeProduct
	case entity.OrderTypeExperience:
		return OrderTypeExperience
	default:
		return OrderTypeUnknown
	}
}

func NewOrderTypeFromString(typ string) OrderType {
	switch typ {
	case "product":
		return OrderTypeProduct
	case "experience":
		return OrderTypeExperience
	default:
		return OrderTypeUnknown
	}
}

func (t OrderType) StoreEntity() entity.OrderType {
	switch t {
	case OrderTypeProduct:
		return entity.OrderTypeProduct
	case OrderTypeExperience:
		return entity.OrderTypeExperience
	default:
		return entity.OrderTypeUnknown
	}
}

func (t OrderType) Response() int32 {
	return int32(t)
}

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

func NewOrder(order *entity.Order, products map[int64]*Product) *Order {
	return &Order{
		Order: response.Order{
			ID:             order.ID,
			CoordinatorID:  order.CoordinatorID,
			PromotionID:    order.PromotionID,
			Type:           NewOrderType(order.Type).Response(),
			Status:         NewOrderStatus(order.Status).Response(),
			Payment:        NewOrderPayment(&order.OrderPayment).Response(),
			Refund:         NewOrderRefund(&order.OrderPayment).Response(),
			Items:          NewOrderItems(order.OrderItems, products).Response(),
			PickupAt:       jst.Unix(order.PickupAt),
			PickupLocation: order.PickupLocation,
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

func NewOrders(orders entity.Orders, products map[int64]*Product) Orders {
	res := make(Orders, len(orders))
	for i := range orders {
		res[i] = NewOrder(orders[i], products)
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
