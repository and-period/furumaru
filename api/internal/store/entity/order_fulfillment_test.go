package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderFulfillments_MapByOrderID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		fulfillments OrderFulfillments
		expect       map[string]*OrderFulfillment
	}{
		{
			name: "success",
			fulfillments: OrderFulfillments{
				{
					ID:      "fulfillment-id01",
					OrderID: "order-id01",
				},
				{
					ID:      "fulfillment-id02",
					OrderID: "order-id02",
				},
			},
			expect: map[string]*OrderFulfillment{
				"order-id01": {
					ID:      "fulfillment-id01",
					OrderID: "order-id01",
				},
				"order-id02": {
					ID:      "fulfillment-id02",
					OrderID: "order-id02",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.fulfillments.MapByOrderID())
		})
	}
}
