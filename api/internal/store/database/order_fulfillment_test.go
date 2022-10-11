package database

import (
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
)

func testOrderFulfillment(id, orderID, shippingID string, now time.Time) *entity.OrderFulfillment {
	return &entity.OrderFulfillment{
		ID:              id,
		OrderID:         orderID,
		ShippingID:      shippingID,
		TrackingNumber:  "",
		ShippingCarrier: entity.ShippingCarrierUnknown,
		ShippingMethod:  entity.DeliveryTypeNormal,
		BoxSize:         entity.ShippingSize60,
		BoxCount:        1,
		WeightTotal:     1000,
		Lastname:        "&.",
		Firstname:       "スタッフ",
		PostalCode:      "1000014",
		Prefecture:      "東京都",
		City:            "千代田区",
		AddressLine1:    "永田町1-7-1",
		AddressLine2:    "",
		PhoneNumber:     "+819012345678",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
