package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderPayments_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		payments OrderPayments
	}{
		{
			name: "success",
			payments: OrderPayments{
				{OrderID: "order-01", MethodType: PaymentMethodTypeCreditCard},
				{OrderID: "order-02", MethodType: PaymentMethodTypePayPay},
			},
		},
		{
			name:     "empty",
			payments: OrderPayments{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var count int
			for range tt.payments.All() {
				count++
			}
			assert.Equal(t, len(tt.payments), count)
		})
	}
}

func TestOrderPayments_IterMapByOrderID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		payments OrderPayments
	}{
		{
			name: "success",
			payments: OrderPayments{
				{OrderID: "order-01", MethodType: PaymentMethodTypeCreditCard},
				{OrderID: "order-02", MethodType: PaymentMethodTypePayPay},
			},
		},
		{
			name:     "empty",
			payments: OrderPayments{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*OrderPayment)
			for k, v := range tt.payments.IterMapByOrderID() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.payments))
			for _, p := range tt.payments {
				assert.Contains(t, result, p.OrderID)
			}
		})
	}
}
