package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
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
	types.Order
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
