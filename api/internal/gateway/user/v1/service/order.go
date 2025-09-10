package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/set"
)

// OrderType - 注文種別
type OrderType types.OrderType

// OrderStatus - 注文ステータス
type OrderStatus types.OrderStatus

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
	case entity.OrderStatusWaiting, entity.OrderStatusPreparing, entity.OrderStatusShipped:
		return OrderStatus(types.OrderStatusPreparing)
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

func NewOrder(order *entity.Order, addresses map[int64]*Address, products map[int64]*Product, experiences map[int64]*Experience) *Order {
	var (
		billingAddress, shippingAddress *Address
		experience                      *Experience
	)
	if address, ok := addresses[order.AddressRevisionID]; ok {
		billingAddress = address
	}
	if exp, ok := experiences[order.ExperienceRevisionID]; ok {
		experience = exp
	}
	if len(order.OrderFulfillments) > 0 {
		// 現状すべての配送先が同一になっているため
		if address, ok := addresses[order.OrderFulfillments[0].AddressRevisionID]; ok {
			shippingAddress = address
		}
	}
	return &Order{
		Order: types.Order{
			ID:              order.ID,
			CoordinatorID:   order.CoordinatorID,
			PromotionID:     order.PromotionID,
			Type:            NewOrderType(order.Type).Response(),
			Status:          NewOrderStatus(order.Status).Response(),
			Payment:         NewOrderPayment(&order.OrderPayment).Response(),
			Refund:          NewOrderRefund(&order.OrderPayment).Response(),
			Fulfillments:    NewOrderFulfillments(order.OrderFulfillments).Response(),
			Items:           NewOrderItems(order.OrderItems, products).Response(),
			Experience:      NewOrderExperience(&order.OrderExperience, experience).Response(),
			BillingAddress:  billingAddress.Response(),
			ShippingAddress: shippingAddress.Response(),
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
