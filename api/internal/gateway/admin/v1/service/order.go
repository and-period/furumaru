package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
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
	OrderStatusWaiting   OrderStatus = 2 // 受注待ち
	OrderStatusPreparing OrderStatus = 3 // 発送準備中
	OrderStatusShipped   OrderStatus = 4 // 発送完了
	OrderStatusCompleted OrderStatus = 5 // 完了
	OrderStatusCanceled  OrderStatus = 6 // キャンセル
	OrderStatusRefunded  OrderStatus = 7 // 返金
	OrderStatusFailed    OrderStatus = 8 // 失敗
)

// OrderShippingType - 発送方法
type OrderShippingType int32

const (
	OrderShippingTypeUnknown  OrderShippingType = 0
	OrderShippingTypeNone     OrderShippingType = 1 // 発送なし
	OrderShippingTypeStandard OrderShippingType = 2 // 通常配送
	OrderShippingTypePickup   OrderShippingType = 3 // 店舗受取
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

func NewOrderShippingType(typ entity.OrderShippingType) OrderShippingType {
	switch typ {
	case entity.OrderShippingTypeNone:
		return OrderShippingTypeNone
	case entity.OrderShippingTypeStandard:
		return OrderShippingTypeStandard
	case entity.OrderShippingTypePickup:
		return OrderShippingTypePickup
	default:
		return OrderShippingTypeUnknown
	}
}

func (t OrderShippingType) Response() int32 {
	return int32(t)
}

func NewOrder(order *entity.Order, addresses map[int64]*Address, products map[int64]*Product, experiences map[int64]*Experience) *Order {
	return &Order{
		Order: response.Order{
			ID:              order.ID,
			UserID:          order.UserID,
			CoordinatorID:   order.CoordinatorID,
			PromotionID:     order.PromotionID,
			ManagementID:    order.ManagementID,
			ShippingMessage: order.ShippingMessage,
			Type:            NewOrderType(order.Type).Response(),
			Status:          NewOrderStatus(order.Status).Response(),
			ShippingType:    NewOrderShippingType(order.ShippingType).Response(),
			Metadata:        NewOrderMetadata(&order.OrderMetadata).Response(),
			Payment:         NewOrderPayment(&order.OrderPayment, addresses[order.OrderPayment.AddressRevisionID]).Response(),
			Refund:          NewOrderRefund(&order.OrderPayment).Response(),
			Fulfillments:    NewOrderFulfillments(order.OrderFulfillments, addresses).Response(),
			Items:           NewOrderItems(order.OrderItems, products).Response(),
			Experience:      NewOrderExperience(&order.OrderExperience, experiences[order.OrderExperience.ExperienceRevisionID]).Response(),
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

func NewOrders(orders entity.Orders, addresses map[int64]*Address, products map[int64]*Product, experiences map[int64]*Experience) Orders {
	res := make(Orders, len(orders))
	for i := range orders {
		res[i] = NewOrder(orders[i], addresses, products, experiences)
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
