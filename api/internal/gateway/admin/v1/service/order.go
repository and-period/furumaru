package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
)

// OrderType - 注文種別
type OrderType types.OrderType

// OrderStatus - 注文ステータス
type OrderStatus types.OrderStatus

// OrderShippingType - 発送方法
type OrderShippingType types.OrderShippingType

type Order struct {
	types.Order
}

type Orders []*Order

func NewOrderType(typ entity.OrderType) OrderType {
	switch typ {
	case entity.OrderTypeProduct:
		return OrderType(types.OrderTypeProduct)
	case entity.OrderTypeExperience:
		return OrderType(types.OrderTypeExperience)
	default:
		return OrderType(types.OrderTypeUnknown)
	}
}

func NewOrderTypeFromString(typ string) OrderType {
	switch typ {
	case "product":
		return OrderType(types.OrderTypeProduct)
	case "experience":
		return OrderType(types.OrderTypeExperience)
	default:
		return OrderType(types.OrderTypeUnknown)
	}
}

func (t OrderType) StoreEntity() entity.OrderType {
	switch types.OrderType(t) {
	case types.OrderTypeProduct:
		return entity.OrderTypeProduct
	case types.OrderTypeExperience:
		return entity.OrderTypeExperience
	default:
		return entity.OrderTypeUnknown
	}
}

func (t OrderType) Response() types.OrderType {
	return types.OrderType(t)
}

func NewOrderStatus(status entity.OrderStatus) OrderStatus {
	switch status {
	case entity.OrderStatusUnpaid:
		return OrderStatus(types.OrderStatusUnpaid)
	case entity.OrderStatusWaiting:
		return OrderStatus(types.OrderStatusWaiting)
	case entity.OrderStatusPreparing:
		return OrderStatus(types.OrderStatusPreparing)
	case entity.OrderStatusShipped:
		return OrderStatus(types.OrderStatusShipped)
	case entity.OrderStatusCompleted:
		return OrderStatus(types.OrderStatusCompleted)
	case entity.OrderStatusCanceled:
		return OrderStatus(types.OrderStatusCanceled)
	case entity.OrderStatusRefunded:
		return OrderStatus(types.OrderStatusRefunded)
	case entity.OrderStatusFailed:
		return OrderStatus(types.OrderStatusFailed)
	default:
		return OrderStatus(types.OrderStatusUnknown)
	}
}

func (s OrderStatus) Response() types.OrderStatus {
	return types.OrderStatus(s)
}

func NewOrderShippingType(typ entity.OrderShippingType) OrderShippingType {
	switch typ {
	case entity.OrderShippingTypeNone:
		return OrderShippingType(types.OrderShippingTypeNone)
	case entity.OrderShippingTypeStandard:
		return OrderShippingType(types.OrderShippingTypeStandard)
	case entity.OrderShippingTypePickup:
		return OrderShippingType(types.OrderShippingTypePickup)
	default:
		return OrderShippingType(types.OrderShippingTypeUnknown)
	}
}

func (t OrderShippingType) Response() types.OrderShippingType {
	return types.OrderShippingType(t)
}

func NewOrder(order *entity.Order, addresses map[int64]*Address, products map[int64]*Product, experiences map[int64]*Experience) *Order {
	return &Order{
		Order: types.Order{
			ID:              order.ID,
			UserID:          order.UserID,
			CoordinatorID:   order.CoordinatorID,
			PromotionID:     order.PromotionID,
			ManagementID:    order.ManagementID,
			ShippingMessage: order.ShippingMessage,
			Type:            NewOrderType(order.Type).Response(),
			Status:          NewOrderStatus(order.Status).Response(),
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
	return set.UniqBy(o.Items, func(i *types.OrderItem) string {
		return i.ProductID
	})
}

func (o *Order) Response() *types.Order {
	return &o.Order
}

func NewOrders(orders entity.Orders, addresses map[int64]*Address, products map[int64]*Product, experiences map[int64]*Experience) Orders {
	res := make(Orders, len(orders))
	for i := range orders {
		res[i] = NewOrder(orders[i], addresses, products, experiences)
	}
	return res
}

func (os Orders) Response() []*types.Order {
	res := make([]*types.Order, len(os))
	for i := range os {
		res[i] = os[i].Response()
	}
	return res
}
