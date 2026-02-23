package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrders_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders Orders
	}{
		{
			name: "success",
			orders: Orders{
				{ID: "order-01", UserID: "user-01"},
				{ID: "order-02", UserID: "user-02"},
			},
		},
		{
			name:   "empty",
			orders: Orders{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, o := range tt.orders.All() {
				indices = append(indices, i)
				ids = append(ids, o.ID)
			}
			for i, o := range tt.orders {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, o.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.orders))
		})
	}
}

func TestOrders_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	orders := Orders{
		{ID: "order-01"},
		{ID: "order-02"},
		{ID: "order-03"},
	}
	var count int
	for range orders.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestOrders_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders Orders
	}{
		{
			name: "success",
			orders: Orders{
				{ID: "order-01", UserID: "user-01"},
				{ID: "order-02", UserID: "user-02"},
			},
		},
		{
			name:   "empty",
			orders: Orders{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Order)
			for k, v := range tt.orders.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.orders))
			for _, o := range tt.orders {
				assert.Contains(t, result, o.ID)
			}
		})
	}
}

func TestAggregatedUserOrders_All(t *testing.T) {
	t.Parallel()
	orders := AggregatedUserOrders{
		{UserID: "user-01", OrderCount: 5},
		{UserID: "user-02", OrderCount: 3},
	}
	var count int
	for range orders.All() {
		count++
	}
	assert.Equal(t, 2, count)
}

func TestAggregatedUserOrders_IterMap(t *testing.T) {
	t.Parallel()
	orders := AggregatedUserOrders{
		{UserID: "user-01", OrderCount: 5},
		{UserID: "user-02", OrderCount: 3},
	}
	result := make(map[string]*AggregatedUserOrder)
	for k, v := range orders.IterMap() {
		result[k] = v
	}
	assert.Len(t, result, 2)
	assert.Contains(t, result, "user-01")
	assert.Contains(t, result, "user-02")
}

func TestAggregatedOrderPayments_All(t *testing.T) {
	t.Parallel()
	payments := AggregatedOrderPayments{
		{PaymentMethodType: PaymentMethodTypeCreditCard, OrderCount: 10},
		{PaymentMethodType: PaymentMethodTypePayPay, OrderCount: 5},
	}
	var count int
	for range payments.All() {
		count++
	}
	assert.Equal(t, 2, count)
}

func TestAggregatedOrderPayments_IterMap(t *testing.T) {
	t.Parallel()
	payments := AggregatedOrderPayments{
		{PaymentMethodType: PaymentMethodTypeCreditCard, OrderCount: 10},
		{PaymentMethodType: PaymentMethodTypePayPay, OrderCount: 5},
	}
	result := make(map[PaymentMethodType]*AggregatedOrderPayment)
	for k, v := range payments.IterMap() {
		result[k] = v
	}
	assert.Len(t, result, 2)
	assert.Contains(t, result, PaymentMethodTypeCreditCard)
	assert.Contains(t, result, PaymentMethodTypePayPay)
}

func TestAggregatedOrderPromotions_All(t *testing.T) {
	t.Parallel()
	promotions := AggregatedOrderPromotions{
		{PromotionID: "promo-01", OrderCount: 5},
		{PromotionID: "promo-02", OrderCount: 3},
	}
	var count int
	for range promotions.All() {
		count++
	}
	assert.Equal(t, 2, count)
}

func TestAggregatedOrderPromotions_IterMap(t *testing.T) {
	t.Parallel()
	promotions := AggregatedOrderPromotions{
		{PromotionID: "promo-01", OrderCount: 5},
		{PromotionID: "promo-02", OrderCount: 3},
	}
	result := make(map[string]*AggregatedOrderPromotion)
	for k, v := range promotions.IterMap() {
		result[k] = v
	}
	assert.Len(t, result, 2)
	assert.Contains(t, result, "promo-01")
	assert.Contains(t, result, "promo-02")
}

func TestAggregatedPeriodOrders_All(t *testing.T) {
	t.Parallel()
	now := time.Now()
	orders := AggregatedPeriodOrders{
		{Period: now, OrderCount: 10},
		{Period: now.Add(24 * time.Hour), OrderCount: 5},
	}
	var count int
	for range orders.All() {
		count++
	}
	assert.Equal(t, 2, count)
}

func TestAggregatedPeriodOrders_IterMapByPeriod(t *testing.T) {
	t.Parallel()
	period1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	period2 := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	orders := AggregatedPeriodOrders{
		{Period: period1, OrderCount: 10},
		{Period: period2, OrderCount: 5},
	}
	result := make(map[time.Time]*AggregatedPeriodOrder)
	for k, v := range orders.IterMapByPeriod() {
		result[k] = v
	}
	assert.Len(t, result, 2)
	assert.Contains(t, result, period1)
	assert.Contains(t, result, period2)
}
