package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayments_MapByOrderID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		payments Payments
		expect   map[string]*Payment
	}{
		{
			name: "success",
			payments: Payments{
				{
					OrderID: "order-id01",
				},
				{
					OrderID: "order-id02",
				},
			},
			expect: map[string]*Payment{
				"order-id01": {
					OrderID: "order-id01",
				},
				"order-id02": {
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
