package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderItems_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		items OrderItems
	}{
		{
			name: "success",
			items: OrderItems{
				{FulfillmentID: "fulfillment-01", OrderID: "order-01"},
				{FulfillmentID: "fulfillment-02", OrderID: "order-01"},
			},
		},
		{
			name:  "empty",
			items: OrderItems{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var count int
			for range tt.items.All() {
				count++
			}
			assert.Equal(t, len(tt.items), count)
		})
	}
}

func TestOrderItems_IterGroupByFulfillmentID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		items      OrderItems
		expectKeys int
	}{
		{
			name: "success",
			items: OrderItems{
				{FulfillmentID: "fulfillment-01", OrderID: "order-01"},
				{FulfillmentID: "fulfillment-01", OrderID: "order-01"},
				{FulfillmentID: "fulfillment-02", OrderID: "order-02"},
			},
			expectKeys: 2,
		},
		{
			name:       "empty",
			items:      OrderItems{},
			expectKeys: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]OrderItems)
			for k, v := range tt.items.IterGroupByFulfillmentID() {
				result[k] = v
			}
			assert.Len(t, result, tt.expectKeys)
		})
	}
}

func TestOrderItems_IterGroupByOrderID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		items      OrderItems
		expectKeys int
	}{
		{
			name: "success",
			items: OrderItems{
				{FulfillmentID: "fulfillment-01", OrderID: "order-01"},
				{FulfillmentID: "fulfillment-02", OrderID: "order-01"},
				{FulfillmentID: "fulfillment-03", OrderID: "order-02"},
			},
			expectKeys: 2,
		},
		{
			name:       "empty",
			items:      OrderItems{},
			expectKeys: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]OrderItems)
			for k, v := range tt.items.IterGroupByOrderID() {
				result[k] = v
			}
			assert.Len(t, result, tt.expectKeys)
		})
	}
}
