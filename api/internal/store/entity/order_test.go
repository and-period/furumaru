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

func TestOrders_IDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders Orders
		expect []string
	}{
		{
			name: "success",
			orders: Orders{
				{
					ID:               "order-id",
					OrderItems:       OrderItems{{ID: "item-id"}},
					OrderPayment:     OrderPayment{ID: "payment-id"},
					OrderFulfillment: OrderFulfillment{ID: "fulfillment-id"},
					OrderActivities:  OrderActivities{{ID: "activity-id"}},
				},
			},
			expect: []string{"order-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.orders.IDs())
		})
	}
}

func TestAggregatedOrders_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders AggregatedOrders
		expect map[string]*AggregatedOrder
	}{
		{
			name: "success",
			orders: AggregatedOrders{
				{
					UserID:     "user-id",
					OrderCount: 2,
					Subtotal:   3000,
					Discount:   0,
				},
			},
			expect: map[string]*AggregatedOrder{
				"user-id": {
					UserID:     "user-id",
					OrderCount: 2,
					Subtotal:   3000,
					Discount:   0,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.orders.Map())
		})
	}
}
