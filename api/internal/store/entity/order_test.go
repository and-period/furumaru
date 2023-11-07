package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrder_Fill(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name         string
		order        *Order
		payment      *OrderPayment
		fulfillments OrderFulfillments
		items        OrderItems
		expect       *Order
	}{
		{
			name: "success",
			order: &Order{
				ID:            "order-id",
				UserID:        "user-id",
				CoordinatorID: "coordinator-id",
				PromotionID:   "promotion-id",
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			payment:      &OrderPayment{OrderID: "order-id"},
			fulfillments: OrderFulfillments{{OrderID: "order-id"}},
			items:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
			expect: &Order{
				ID:                "order-id",
				UserID:            "user-id",
				CoordinatorID:     "coordinator-id",
				PromotionID:       "promotion-id",
				OrderPayment:      OrderPayment{OrderID: "order-id"},
				OrderFulfillments: OrderFulfillments{{OrderID: "order-id"}},
				OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.order.Fill(tt.payment, tt.fulfillments, tt.items)
			assert.Equal(t, tt.expect, tt.order)
		})
	}
}

func TestOrders_IsCanceled(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		order  *Order
		expect bool
	}{
		{
			name: "canceled",
			order: &Order{
				ID:                "order-id",
				OrderPayment:      OrderPayment{OrderID: "order-id", Status: PaymentStatusRefunded},
				OrderFulfillments: OrderFulfillments{{OrderID: "order-id"}},
				OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
			},
			expect: true,
		},
		{
			name: "not canceled",
			order: &Order{
				ID:                "order-id",
				OrderPayment:      OrderPayment{OrderID: "order-id", Status: PaymentStatusPending},
				OrderFulfillments: OrderFulfillments{{OrderID: "order-id"}},
				OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
			},
			expect: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.order.IsCanceled())
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
					ID:                "order-id",
					OrderPayment:      OrderPayment{OrderID: "order-id"},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id"}},
					OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
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

func TestOrders_UserIDs(t *testing.T) {
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
					ID:                "order-id",
					UserID:            "user-id",
					CoordinatorID:     "coordinator-id",
					OrderPayment:      OrderPayment{OrderID: "order-id"},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id"}},
					OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
				},
			},
			expect: []string{"user-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.orders.UserIDs())
		})
	}
}

func TestOrders_CoordinatorID(t *testing.T) {
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
					ID:                "order-id",
					UserID:            "user-id",
					CoordinatorID:     "coordinator-id",
					OrderPayment:      OrderPayment{OrderID: "order-id"},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id"}},
					OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
				},
			},
			expect: []string{"coordinator-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.orders.CoordinatorIDs())
		})
	}
}

func TestOrders_PromotionIDs(t *testing.T) {
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
					ID:                "order-id01",
					UserID:            "user-id",
					CoordinatorID:     "coordinator-id",
					PromotionID:       "promotion-id",
					OrderPayment:      OrderPayment{OrderID: "order-id"},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id"}},
					OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
				},
				{
					ID:                "order-id02",
					UserID:            "user-id",
					CoordinatorID:     "coordinator-id",
					OrderPayment:      OrderPayment{OrderID: "order-id"},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id"}},
					OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
				},
			},
			expect: []string{"promotion-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.orders.PromotionIDs())
		})
	}
}

func TestOrders_AddressRevisionIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders Orders
		expect []int64
	}{
		{
			name: "success",
			orders: Orders{
				{
					ID:                "order-id",
					UserID:            "user-id",
					OrderPayment:      OrderPayment{OrderID: "order-id", AddressRevisionID: 1},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id", AddressRevisionID: 2}},
					OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
				},
			},
			expect: []int64{1, 2},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.orders.AddressRevisionIDs())
		})
	}
}

func TestOrders_ProductRevisionIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders Orders
		expect []int64
	}{
		{
			name: "success",
			orders: Orders{
				{
					ID:                "order-id",
					UserID:            "user-id",
					OrderPayment:      OrderPayment{OrderID: "order-id", AddressRevisionID: 1},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id", AddressRevisionID: 2}},
					OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
				},
			},
			expect: []int64{1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.orders.ProductRevisionIDs())
		})
	}
}

func TestOrders_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		orders       Orders
		payments     map[string]*OrderPayment
		fulfillments map[string]OrderFulfillments
		items        map[string]OrderItems
		expect       Orders
	}{
		{
			name: "success",
			orders: Orders{
				{ID: "order-id01"},
				{ID: "order-id02"},
				{ID: "order-id03"},
			},
			payments: map[string]*OrderPayment{
				"order-id01": {OrderID: "order-id01"},
				"order-id02": {OrderID: "order-id02"},
			},
			fulfillments: map[string]OrderFulfillments{
				"order-id01": {{OrderID: "order-id01"}},
				"order-id02": {{OrderID: "order-id02"}},
			},
			items: map[string]OrderItems{
				"order-id01": {{OrderID: "order-id01", ProductRevisionID: 1}},
				"order-id02": {{OrderID: "order-id02", ProductRevisionID: 1}},
			},
			expect: Orders{
				{
					ID:                "order-id01",
					OrderPayment:      OrderPayment{OrderID: "order-id01"},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id01"}},
					OrderItems:        OrderItems{{OrderID: "order-id01", ProductRevisionID: 1}},
				},
				{
					ID:                "order-id02",
					OrderPayment:      OrderPayment{OrderID: "order-id02"},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id02"}},
					OrderItems:        OrderItems{{OrderID: "order-id02", ProductRevisionID: 1}},
				},
				{
					ID:           "order-id03",
					OrderPayment: OrderPayment{},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.orders.Fill(tt.payments, tt.fulfillments, tt.items)
			assert.Equal(t, tt.expect, tt.orders)
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
