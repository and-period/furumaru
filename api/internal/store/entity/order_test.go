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
		payment     *Payment
		fulfillment *Fulfillment
		activities  Activities
		items       OrderItems
		expect      *Order
	}{
		{
			name:        "success",
			order:       &Order{},
			payment:     &Payment{OrderID: "order-id"},
			fulfillment: &Fulfillment{OrderID: "order-id"},
			activities:  Activities{{OrderID: "order-id", ID: "activity-id"}},
			items:       OrderItems{{OrderID: "order-id", ProductID: "item-id"}},
			expect: &Order{
				Payment:     Payment{OrderID: "order-id"},
				Fulfillment: Fulfillment{OrderID: "order-id"},
				Activities:  Activities{{OrderID: "order-id", ID: "activity-id"}},
				OrderItems:  OrderItems{{OrderID: "order-id", ProductID: "item-id"}},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.order.Fill(tt.payment, tt.fulfillment, tt.activities, tt.items)
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
					ID:          "order-id",
					Payment:     Payment{OrderID: "order-id"},
					Fulfillment: Fulfillment{OrderID: "order-id"},
					Activities:  Activities{{OrderID: "order-id", ID: "activity-id"}},
					OrderItems:  OrderItems{{OrderID: "order-id", ProductID: "item-id"}},
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
