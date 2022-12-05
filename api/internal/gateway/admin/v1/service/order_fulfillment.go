package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
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

type OrderFulfillment struct {
	response.OrderFulfillment
	orderID string
}

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

func NewOrderFulfillment(fulfillment *entity.Fulfillment, status entity.FulfillmentStatus) *OrderFulfillment {
	return &OrderFulfillment{
		OrderFulfillment: response.OrderFulfillment{
			TrackingNumber:  fulfillment.TrackingNumber,
			Status:          NewFulfillmentStatus(status).Response(),
			ShippingCarrier: NewShippingCarrier(fulfillment.ShippingCarrier).Response(),
			ShippingMethod:  NewDeliveryType(fulfillment.ShippingMethod).Response(),
			BoxSize:         NewShippingSize(fulfillment.BoxSize).Response(),
			AddressID:       fulfillment.AddressID,
		},
		orderID: fulfillment.OrderID,
	}
}

func (f *OrderFulfillment) Fill(address *Address) {
	f.Address = address.Response()
}

func (f *OrderFulfillment) Response() *response.OrderFulfillment {
	return &f.OrderFulfillment
}
