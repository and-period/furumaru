package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderFulfillments_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		fulfillments OrderFulfillments
	}{
		{
			name: "success",
			fulfillments: OrderFulfillments{
				{ID: "fulfillment-01", OrderID: "order-01"},
				{ID: "fulfillment-02", OrderID: "order-01"},
			},
		},
		{
			name:         "empty",
			fulfillments: OrderFulfillments{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var count int
			for range tt.fulfillments.All() {
				count++
			}
			assert.Equal(t, len(tt.fulfillments), count)
		})
	}
}

func TestOrderFulfillments_IterGroupByOrderID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		fulfillments OrderFulfillments
		expectKeys   int
	}{
		{
			name: "success",
			fulfillments: OrderFulfillments{
				{ID: "fulfillment-01", OrderID: "order-01"},
				{ID: "fulfillment-02", OrderID: "order-01"},
				{ID: "fulfillment-03", OrderID: "order-02"},
			},
			expectKeys: 2,
		},
		{
			name:         "empty",
			fulfillments: OrderFulfillments{},
			expectKeys:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]OrderFulfillments)
			for k, v := range tt.fulfillments.IterGroupByOrderID() {
				result[k] = v
			}
			assert.Len(t, result, tt.expectKeys)
		})
	}
}
