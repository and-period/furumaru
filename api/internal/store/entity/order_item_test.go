package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderItem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewOrderItemParams
		expect *OrderItem
	}{
		{
			name: "success",
			params: &NewOrderItemParams{
				OrderID:       "order-id",
				FulfillmentID: "fulfillment-id",
				Item: &CartItem{
					ProductID: "product-id",
					Quantity:  1,
				},
				Product: &Product{
					ID:   "product-id",
					Name: "じゃがいも",
					ProductRevision: ProductRevision{
						ID:        1,
						ProductID: "product-id01",
						Price:     500,
					},
				},
			},
			expect: &OrderItem{
				FulfillmentID:     "fulfillment-id",
				ProductRevisionID: 1,
				OrderID:           "order-id",
				Quantity:          1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewOrderItem(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestOrderItems(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		params    *NewOrderItemsParams
		expect    OrderItems
		expectErr error
	}{
		{
			name: "success",
			params: &NewOrderItemsParams{
				OrderID: "order-id",
				Fulfillment: &OrderFulfillment{
					ID:                "fulfillment-id",
					OrderID:           "order-id",
					AddressRevisionID: 1,
					Status:            FulfillmentStatusUnfulfilled,
					TrackingNumber:    "",
					ShippingCarrier:   ShippingCarrierUnknown,
					ShippingType:      ShippingTypeNormal,
					BoxNumber:         1,
					BoxSize:           ShippingSize60,
				},
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
			expect: OrderItems{
				{
					FulfillmentID:     "fulfillment-id",
					ProductRevisionID: 1,
					OrderID:           "order-id",
					Quantity:          1,
				},
				{
					FulfillmentID:     "fulfillment-id",
					ProductRevisionID: 2,
					OrderID:           "order-id",
					Quantity:          2,
				},
			},
			expectErr: nil,
		},
		{
			name: "not found",
			params: &NewOrderItemsParams{
				OrderID: "order-id",
				Fulfillment: &OrderFulfillment{
					ID:                "fulfillment-id",
					OrderID:           "order-id",
					AddressRevisionID: 1,
					Status:            FulfillmentStatusUnfulfilled,
					TrackingNumber:    "",
					ShippingCarrier:   ShippingCarrierUnknown,
					ShippingType:      ShippingTypeNormal,
					BoxNumber:         1,
					BoxSize:           ShippingSize60,
				},
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
				Products: nil,
			},
			expect:    nil,
			expectErr: errNotFoundProduct,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewOrderItems(tt.params)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestOrderItems_ProductRevisionIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		items  OrderItems
		expect []int64
	}{
		{
			name: "success",
			items: OrderItems{
				{
					FulfillmentID:     "fulfillment-id01",
					ProductRevisionID: 1,
					OrderID:           "order-id01",
				},
				{
					FulfillmentID:     "fulfillment-id01",
					ProductRevisionID: 2,
					OrderID:           "order-id01",
				},
				{
					FulfillmentID:     "fulfillment-id02",
					ProductRevisionID: 1,
					OrderID:           "order-id02",
				},
			},
			expect: []int64{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.items.ProductRevisionIDs())
		})
	}
}

func TestOrderItems_GroupByFulfillmentID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		items  OrderItems
		expect map[string]OrderItems
	}{
		{
			name: "success",
			items: OrderItems{
				{
					FulfillmentID:     "fulfillment-id01",
					ProductRevisionID: 1,
					OrderID:           "order-id01",
				},
				{
					FulfillmentID:     "fulfillment-id01",
					ProductRevisionID: 2,
					OrderID:           "order-id01",
				},
				{
					FulfillmentID:     "fulfillment-id02",
					ProductRevisionID: 1,
					OrderID:           "order-id02",
				},
			},
			expect: map[string]OrderItems{
				"fulfillment-id01": {
					{
						FulfillmentID:     "fulfillment-id01",
						ProductRevisionID: 1,
						OrderID:           "order-id01",
					},
					{
						FulfillmentID:     "fulfillment-id01",
						ProductRevisionID: 2,
						OrderID:           "order-id01",
					},
				},
				"fulfillment-id02": {
					{
						FulfillmentID:     "fulfillment-id02",
						ProductRevisionID: 1,
						OrderID:           "order-id02",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.items.GroupByFulfillmentID())
		})
	}
}

func TestOrderItems_GroupByOrderID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		items  OrderItems
		expect map[string]OrderItems
	}{
		{
			name: "success",
			items: OrderItems{
				{
					FulfillmentID:     "fulfillment-id01",
					ProductRevisionID: 1,
					OrderID:           "order-id01",
				},
				{
					FulfillmentID:     "fulfillment-id01",
					ProductRevisionID: 2,
					OrderID:           "order-id01",
				},
				{
					FulfillmentID:     "fulfillment-id02",
					ProductRevisionID: 1,
					OrderID:           "order-id02",
				},
			},
			expect: map[string]OrderItems{
				"order-id01": {
					{
						FulfillmentID:     "fulfillment-id01",
						ProductRevisionID: 1,
						OrderID:           "order-id01",
					},
					{
						FulfillmentID:     "fulfillment-id01",
						ProductRevisionID: 2,
						OrderID:           "order-id01",
					},
				},
				"order-id02": {
					{
						FulfillmentID:     "fulfillment-id02",
						ProductRevisionID: 1,
						OrderID:           "order-id02",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.items.GroupByOrderID())
		})
	}
}
