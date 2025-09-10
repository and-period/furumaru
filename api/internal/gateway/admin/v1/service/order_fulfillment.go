package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

// FulfillmentStatus - 配送状況
type FulfillmentStatus types.FulfillmentStatus

// ShippingCarrier - 配送会社
type ShippingCarrier types.ShippingCarrier

// ShippingSize - 配送時の箱の大きさ
type ShippingSize types.ShippingSize

// ShippingType - 配送方法
type ShippingType types.ShippingType

type OrderFulfillment struct {
	types.OrderFulfillment
	orderID string
}

type OrderFulfillments []*OrderFulfillment

func NewFulfillmentStatus(status entity.FulfillmentStatus) FulfillmentStatus {
	switch status {
	case entity.FulfillmentStatusUnfulfilled:
		return FulfillmentStatus(types.FulfillmentStatusUnfulfilled)
	case entity.FulfillmentStatusFulfilled:
		return FulfillmentStatus(types.FulfillmentStatusFulfilled)
	default:
		return FulfillmentStatus(types.FulfillmentStatusUnknown)
	}
}

func (s FulfillmentStatus) Response() types.FulfillmentStatus {
	return types.FulfillmentStatus(s)
}

func NewShippingCarrier(carrier entity.ShippingCarrier) ShippingCarrier {
	switch carrier {
	case entity.ShippingCarrierYamato:
		return ShippingCarrier(types.ShippingCarrierYamato)
	case entity.ShippingCarrierSagawa:
		return ShippingCarrier(types.ShippingCarrierSagawa)
	default:
		return ShippingCarrier(types.ShippingCarrierUnknown)
	}
}

func (c ShippingCarrier) Response() types.ShippingCarrier {
	return types.ShippingCarrier(c)
}

func NewShippingSize(size entity.ShippingSize) ShippingSize {
	switch size {
	case entity.ShippingSize60:
		return ShippingSize(types.ShippingSize60)
	case entity.ShippingSize80:
		return ShippingSize(types.ShippingSize80)
	case entity.ShippingSize100:
		return ShippingSize(types.ShippingSize100)
	default:
		return ShippingSize(types.ShippingSizeUnknown)
	}
}

func (s ShippingSize) Response() types.ShippingSize {
	return types.ShippingSize(s)
}

func NewShippingType(typ entity.ShippingType) ShippingType {
	switch typ {
	case entity.ShippingTypeNormal:
		return ShippingType(types.ShippingTypeNormal)
	case entity.ShippingTypeFrozen:
		return ShippingType(types.ShippingTypeFrozen)
	case entity.ShippingTypePickup:
		return ShippingType(types.ShippingTypePickup)
	default:
		return ShippingType(types.ShippingTypeUnknown)
	}
}

func (t ShippingType) Response() types.ShippingType {
	return types.ShippingType(t)
}

func NewOrderFulfillment(fulfillment *entity.OrderFulfillment, address *Address) *OrderFulfillment {
	return &OrderFulfillment{
		OrderFulfillment: types.OrderFulfillment{
			FulfillmentID:   fulfillment.ID,
			TrackingNumber:  fulfillment.TrackingNumber,
			Status:          NewFulfillmentStatus(fulfillment.Status).Response(),
			ShippingCarrier: NewShippingCarrier(fulfillment.ShippingCarrier).Response(),
			ShippingType:    NewShippingType(fulfillment.ShippingType).Response(),
			BoxNumber:       fulfillment.BoxNumber,
			BoxSize:         NewShippingSize(fulfillment.BoxSize).Response(),
			BoxRate:         fulfillment.BoxRate,
			ShippedAt:       jst.Unix(fulfillment.ShippedAt),
			Address:         address.Response(),
		},
		orderID: fulfillment.OrderID,
	}
}

func (f *OrderFulfillment) Response() *types.OrderFulfillment {
	return &f.OrderFulfillment
}

func NewOrderFulfillments(fulfillments entity.OrderFulfillments, addresses map[int64]*Address) OrderFulfillments {
	res := make(OrderFulfillments, len(fulfillments))
	for i, f := range fulfillments {
		res[i] = NewOrderFulfillment(f, addresses[f.AddressRevisionID])
	}
	return res
}

func (fs OrderFulfillments) Response() []*types.OrderFulfillment {
	res := make([]*types.OrderFulfillment, len(fs))
	for i := range fs {
		res[i] = fs[i].Response()
	}
	return res
}
