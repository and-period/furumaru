package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestOrderFulfillment(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewOrderFulfillmentParams
		expect *OrderFulfillment
	}{
		{
			name: "success",
			params: &NewOrderFulfillmentParams{
				OrderID: "order-id",
				Address: &entity.Address{
					AddressRevision: entity.AddressRevision{
						ID:             1,
						AddressID:      "address-id",
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
					ID:        "address-id",
					UserID:    "user-id",
					IsDefault: false,
				},
				Basket: &CartBasket{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: []*CartItem{
						{
							ProductID: "product-id01",
							Quantity:  1,
						},
						{
							ProductID: "product-id02",
							Quantity:  2,
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			expect: &OrderFulfillment{
				OrderID:           "order-id",
				AddressRevisionID: 1,
				Status:            FulfillmentStatusUnfulfilled,
				TrackingNumber:    "",
				ShippingCarrier:   ShippingCarrierUnknown,
				ShippingType:      ShippingTypeNormal,
				BoxNumber:         1,
				BoxSize:           ShippingSize60,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewOrderFulfillment(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestOrderFulfillments(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		params             *NewOrderFulfillmentsParams
		expectFulfillments OrderFulfillments
		expectItems        OrderItems
		expectErr          error
	}{
		{
			name: "success",
			params: &NewOrderFulfillmentsParams{
				OrderID: "order-id",
				Address: &entity.Address{
					AddressRevision: entity.AddressRevision{
						ID:             1,
						AddressID:      "address-id",
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
					ID:        "address-id",
					UserID:    "user-id",
					IsDefault: false,
				},
				Baskets: CartBaskets{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: []*CartItem{
							{
								ProductID: "product-id01",
								Quantity:  1,
							},
							{
								ProductID: "product-id02",
								Quantity:  2,
							},
						},
						CoordinatorID: "coordinator-id",
					},
				},
				Products: map[string]*Product{
					"product-id01": {
						ID:   "product-id01",
						Name: "じゃがいも",
						ProductRevision: ProductRevision{
							ID:        1,
							ProductID: "product-id01",
							Price:     500,
						},
					},
					"product-id02": {
						ID:   "product-id02",
						Name: "人参",
						ProductRevision: ProductRevision{
							ID:        2,
							ProductID: "product-id02",
							Price:     1980,
						},
					},
				},
			},
			expectFulfillments: OrderFulfillments{
				{
					OrderID:           "order-id",
					AddressRevisionID: 1,
					Status:            FulfillmentStatusUnfulfilled,
					TrackingNumber:    "",
					ShippingCarrier:   ShippingCarrierUnknown,
					ShippingType:      ShippingTypeNormal,
					BoxNumber:         1,
					BoxSize:           ShippingSize60,
				},
			},
			expectItems: OrderItems{
				{
					ProductRevisionID: 1,
					OrderID:           "order-id",
					Quantity:          1,
				},
				{
					ProductRevisionID: 2,
					OrderID:           "order-id",
					Quantity:          2,
				},
			},
			expectErr: nil,
		},
		{
			name: "failed to create order items",
			params: &NewOrderFulfillmentsParams{
				OrderID: "order-id",
				Address: &entity.Address{
					AddressRevision: entity.AddressRevision{
						ID:             1,
						AddressID:      "address-id",
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
					ID:        "address-id",
					UserID:    "user-id",
					IsDefault: false,
				},
				Baskets: CartBaskets{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: []*CartItem{
							{
								ProductID: "product-id01",
								Quantity:  1,
							},
							{
								ProductID: "product-id02",
								Quantity:  2,
							},
						},
						CoordinatorID: "coordinator-id",
					},
				},
				Products: nil,
			},
			expectFulfillments: nil,
			expectItems:        nil,
			expectErr:          errNotFoundProduct,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fulfillments, items, err := NewOrderFulfillments(tt.params)
			for _, f := range fulfillments {
				f.ID = "" // ignore
			}
			for _, i := range items {
				i.FulfillmentID = "" // ignore
			}
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expectFulfillments, fulfillments)
			assert.ElementsMatch(t, tt.expectItems, items)
		})
	}
}

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
					ShippingType:      ShippingTypeNormal,
					BoxNumber:         1,
					BoxSize:           ShippingSize100,
				},
				{
					ID:                "fulfillment-id02",
					OrderID:           "order-id01",
					AddressRevisionID: 1,
					Status:            FulfillmentStatusUnfulfilled,
					ShippingCarrier:   ShippingCarrierYamato,
					ShippingType:      ShippingTypeNormal,
					BoxNumber:         2,
					BoxSize:           ShippingSize80,
				},
				{
					ID:                "fulfillment-id03",
					OrderID:           "order-id02",
					AddressRevisionID: 1,
					Status:            FulfillmentStatusUnfulfilled,
					ShippingCarrier:   ShippingCarrierYamato,
					ShippingType:      ShippingTypeFrozen,
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
					ShippingType:      ShippingTypeNormal,
					BoxNumber:         1,
					BoxSize:           ShippingSize100,
				},
				{
					ID:                "fulfillment-id02",
					OrderID:           "order-id01",
					AddressRevisionID: 1,
					Status:            FulfillmentStatusUnfulfilled,
					ShippingCarrier:   ShippingCarrierYamato,
					ShippingType:      ShippingTypeNormal,
					BoxNumber:         2,
					BoxSize:           ShippingSize80,
				},
				{
					ID:                "fulfillment-id03",
					OrderID:           "order-id02",
					AddressRevisionID: 1,
					Status:            FulfillmentStatusUnfulfilled,
					ShippingCarrier:   ShippingCarrierYamato,
					ShippingType:      ShippingTypeFrozen,
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
						ShippingType:      ShippingTypeNormal,
						BoxNumber:         1,
						BoxSize:           ShippingSize100,
					},
					{
						ID:                "fulfillment-id02",
						OrderID:           "order-id01",
						AddressRevisionID: 1,
						Status:            FulfillmentStatusUnfulfilled,
						ShippingCarrier:   ShippingCarrierYamato,
						ShippingType:      ShippingTypeNormal,
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
						ShippingType:      ShippingTypeFrozen,
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
