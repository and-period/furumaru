package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderPayments_MapByOrderID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		payments OrderPayments
		expect   map[string]*OrderPayment
	}{
		{
			name: "success",
			payments: OrderPayments{
				{
					ID:      "payment-id01",
					OrderID: "order-id01",
				},
				{
					ID:      "payment-id02",
					OrderID: "order-id02",
				},
			},
			expect: map[string]*OrderPayment{
				"order-id01": {
					ID:      "payment-id01",
					OrderID: "order-id01",
				},
				"order-id02": {
					ID:      "payment-id02",
					OrderID: "order-id02",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.payments.MapByOrderID())
		})
	}
}
