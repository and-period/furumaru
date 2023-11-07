package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		tt := tt
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
		tt := tt
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.items.GroupByOrderID())
		})
	}
}
