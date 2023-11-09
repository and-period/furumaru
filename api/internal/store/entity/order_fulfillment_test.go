package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFulfillments_AddressRevisionIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		fulfillments OrderFulfillments
		expect       []int64
	}{
		{
			name: "success",
			fulfillments: OrderFulfillments{
				{
					ID:                "fulfillment-id01",
					OrderID:           "order-id01",
					AddressRevisionID: 1,
					Status:            FulfillmentStatusUnfulfilled,
					ShippingCarrier:   ShippingCarrierYamato,
					ShippingMethod:    DeliveryTypeNormal,
					BoxNumber:         1,
					BoxSize:           ShippingSize100,
				},
				{
					ID:                "fulfillment-id02",
					OrderID:           "order-id01",
					AddressRevisionID: 1,
					Status:            FulfillmentStatusUnfulfilled,
					ShippingCarrier:   ShippingCarrierYamato,
					ShippingMethod:    DeliveryTypeNormal,
					BoxNumber:         2,
					BoxSize:           ShippingSize80,
				},
				{
					ID:                "fulfillment-id03",
					OrderID:           "order-id02",
					AddressRevisionID: 1,
					Status:            FulfillmentStatusUnfulfilled,
					ShippingCarrier:   ShippingCarrierYamato,
					ShippingMethod:    DeliveryTypeFrozen,
					BoxNumber:         1,
					BoxSize:           ShippingSize80,
				},
			},
			expect: []int64{1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.fulfillments.AddressRevisionIDs())
		})
	}
}

func TestFulfillments_GroupByOrderID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		fulfillments OrderFulfillments
		expect       map[string]OrderFulfillments
	}{
		{
			name: "success",
			fulfillments: OrderFulfillments{
				{
					ID:                "fulfillment-id01",
					OrderID:           "order-id01",
					AddressRevisionID: 1,
					Status:            FulfillmentStatusUnfulfilled,
					ShippingCarrier:   ShippingCarrierYamato,
					ShippingMethod:    DeliveryTypeNormal,
					BoxNumber:         1,
					BoxSize:           ShippingSize100,
				},
				{
					ID:                "fulfillment-id02",
					OrderID:           "order-id01",
					AddressRevisionID: 1,
					Status:            FulfillmentStatusUnfulfilled,
					ShippingCarrier:   ShippingCarrierYamato,
					ShippingMethod:    DeliveryTypeNormal,
					BoxNumber:         2,
					BoxSize:           ShippingSize80,
				},
				{
					ID:                "fulfillment-id03",
					OrderID:           "order-id02",
					AddressRevisionID: 1,
					Status:            FulfillmentStatusUnfulfilled,
					ShippingCarrier:   ShippingCarrierYamato,
					ShippingMethod:    DeliveryTypeFrozen,
					BoxNumber:         1,
					BoxSize:           ShippingSize80,
				},
			},
			expect: map[string]OrderFulfillments{
				"order-id01": {
					{
						ID:                "fulfillment-id01",
						OrderID:           "order-id01",
						AddressRevisionID: 1,
						Status:            FulfillmentStatusUnfulfilled,
						ShippingCarrier:   ShippingCarrierYamato,
						ShippingMethod:    DeliveryTypeNormal,
						BoxNumber:         1,
						BoxSize:           ShippingSize100,
					},
					{
						ID:                "fulfillment-id02",
						OrderID:           "order-id01",
						AddressRevisionID: 1,
						Status:            FulfillmentStatusUnfulfilled,
						ShippingCarrier:   ShippingCarrierYamato,
						ShippingMethod:    DeliveryTypeNormal,
						BoxNumber:         2,
						BoxSize:           ShippingSize80,
					},
				},
				"order-id02": {
					{
						ID:                "fulfillment-id03",
						OrderID:           "order-id02",
						AddressRevisionID: 1,
						Status:            FulfillmentStatusUnfulfilled,
						ShippingCarrier:   ShippingCarrierYamato,
						ShippingMethod:    DeliveryTypeFrozen,
						BoxNumber:         1,
						BoxSize:           ShippingSize80,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.fulfillments.GroupByOrderID())
		})
	}
}
