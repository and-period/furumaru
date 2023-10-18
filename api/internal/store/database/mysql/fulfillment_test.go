package mysql

import (
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
)

func testFulfillment(orderID, addressID string, now time.Time) *entity.Fulfillment {
	return &entity.Fulfillment{
		OrderID:         orderID,
		AddressID:       addressID,
		TrackingNumber:  "",
		ShippingCarrier: entity.ShippingCarrierUnknown,
		ShippingMethod:  entity.DeliveryTypeNormal,
		BoxSize:         entity.ShippingSize60,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
