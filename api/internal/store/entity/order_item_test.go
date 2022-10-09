package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
					ID:      "item-id01",
					OrderID: "order-id01",
				},
				{
					ID:      "item-id02",
					OrderID: "order-id01",
				},
				{
					ID:      "item-id03",
					OrderID: "order-id02",
				},
			},
			expect: map[string]OrderItems{
				"order-id01": {
					{
						ID:      "item-id01",
						OrderID: "order-id01",
					},
					{
						ID:      "item-id02",
						OrderID: "order-id01",
					},
				},
				"order-id02": {
					{
						ID:      "item-id03",
						OrderID: "order-id02",
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
