package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

// FulfillmentStatus - 配送状況
type FulfillmentStatus int32

const (
	FulfillmentStatusUnknown     FulfillmentStatus = 0
	FulfillmentStatusUnfulfilled FulfillmentStatus = 1 // 未発送
	FulfillmentStatusFulfilled   FulfillmentStatus = 2 // 発送済み
)

// ShippingCarrier - 配送会社
type ShippingCarrier int32

const (
	ShippingCarrierUnknown ShippingCarrier = 0
	ShippingCarrierYamato  ShippingCarrier = 1 // ヤマト運輸
	ShippingCarrierSagawa  ShippingCarrier = 2 // 佐川急便
)

// ShippingSize - 配送時の箱の大きさ
type ShippingSize int32

const (
	ShippingSizeUnknown ShippingSize = 0
	ShippingSize60      ShippingSize = 1 // 箱のサイズ:60
	ShippingSize80      ShippingSize = 2 // 箱のサイズ:80
	ShippingSize100     ShippingSize = 3 // 箱のサイズ:100
)

// ShippingType - 配送方法
type ShippingType int32

const (
	ShippingTypeUnknown ShippingType = 0
	ShippingTypeNormal  ShippingType = 1 // 常温・冷蔵便
	ShippingTypeFrozen  ShippingType = 2 // 冷凍便
)

type OrderFulfillment struct {
	response.OrderFulfillment
	orderID string
}

type OrderFulfillments []*OrderFulfillment

func NewFulfillmentStatus(status entity.FulfillmentStatus) FulfillmentStatus {
	switch status {
	case entity.FulfillmentStatusUnfulfilled:
		return FulfillmentStatusUnfulfilled
	case entity.FulfillmentStatusFulfilled:
		return FulfillmentStatusFulfilled
	default:
		return FulfillmentStatusUnknown
	}
}

func (s FulfillmentStatus) Response() int32 {
	return int32(s)
}

func NewShippingCarrier(carrier entity.ShippingCarrier) ShippingCarrier {
	switch carrier {
	case entity.ShippingCarrierYamato:
		return ShippingCarrierYamato
	case entity.ShippingCarrierSagawa:
		return ShippingCarrierSagawa
	default:
		return ShippingCarrierUnknown
	}
}

func (c ShippingCarrier) Response() int32 {
	return int32(c)
}

func NewShippingSize(size entity.ShippingSize) ShippingSize {
	switch size {
	case entity.ShippingSize60:
		return ShippingSize60
	case entity.ShippingSize80:
		return ShippingSize80
	case entity.ShippingSize100:
		return ShippingSize100
	default:
		return ShippingSizeUnknown
	}
}

func (s ShippingSize) Response() int32 {
	return int32(s)
}

func NewShippingType(typ entity.ShippingType) ShippingType {
	switch typ {
	case entity.ShippingTypeNormal:
		return ShippingTypeNormal
	case entity.ShippingTypeFrozen:
		return ShippingTypeFrozen
	default:
		return ShippingTypeUnknown
	}
}

func (t ShippingType) Response() int32 {
	return int32(t)
}

func NewOrderFulfillment(fulfillment *entity.OrderFulfillment, address *Address) *OrderFulfillment {
	return &OrderFulfillment{
		OrderFulfillment: response.OrderFulfillment{
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

func (f *OrderFulfillment) Response() *response.OrderFulfillment {
	return &f.OrderFulfillment
}

func NewOrderFulfillments(fulfillments entity.OrderFulfillments, addresses map[int64]*Address) OrderFulfillments {
	res := make(OrderFulfillments, len(fulfillments))
	for i, f := range fulfillments {
		res[i] = NewOrderFulfillment(f, addresses[f.AddressRevisionID])
	}
	return res
}

func (fs OrderFulfillments) Response() []*response.OrderFulfillment {
	res := make([]*response.OrderFulfillment, len(fs))
	for i := range fs {
		res[i] = fs[i].Response()
	}
	return res
}
