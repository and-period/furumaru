package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrder_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		order       *Order
		items       OrderItems
		payment     *OrderPayment
		fulfillment *OrderFulfillment
		activities  OrderActivities
		expect      *Order
	}{
		{
			name:        "success",
			order:       &Order{},
			items:       OrderItems{{ID: "item-id"}},
			payment:     &OrderPayment{ID: "payment-id"},
			fulfillment: &OrderFulfillment{ID: "fulfillment-id"},
			activities:  OrderActivities{{ID: "activity-id"}},
			expect: &Order{
				OrderItems:       OrderItems{{ID: "item-id"}},
				OrderPayment:     OrderPayment{ID: "payment-id"},
				OrderFulfillment: OrderFulfillment{ID: "fulfillment-id"},
				OrderActivities:  OrderActivities{{ID: "activity-id"}},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.order.Fill(tt.items, tt.payment, tt.fulfillment, tt.activities)
			assert.Equal(t, tt.expect, tt.order)
		})
	}
}
