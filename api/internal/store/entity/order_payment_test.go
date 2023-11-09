package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderPayments_AddressRevisionIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		payments OrderPayments
		expect   []int64
	}{
		{
			name: "success",
			payments: OrderPayments{
				{
					OrderID:           "order-id01",
					AddressRevisionID: 1,
					TransactionID:     "transaction-id01",
					Status:            PaymentStatusCaptured,
					MethodType:        PaymentMethodTypeCreditCard,
					Subtotal:          1980,
					Discount:          0,
					ShippingFee:       550,
					Tax:               253,
					Total:             2783,
				},
				{
					OrderID:           "order-id02",
					AddressRevisionID: 1,
					TransactionID:     "transaction-id02",
					Status:            PaymentStatusCaptured,
					MethodType:        PaymentMethodTypeCreditCard,
					Subtotal:          3000,
					Discount:          0,
					ShippingFee:       0,
					Tax:               300,
					Total:             3300,
				},
			},
			expect: []int64{1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.payments.AddressRevisionIDs())
		})
	}
}

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
					OrderID:           "order-id01",
					AddressRevisionID: 1,
					TransactionID:     "transaction-id01",
					Status:            PaymentStatusCaptured,
					MethodType:        PaymentMethodTypeCreditCard,
					Subtotal:          1980,
					Discount:          0,
					ShippingFee:       550,
					Tax:               253,
					Total:             2783,
				},
				{
					OrderID:           "order-id02",
					AddressRevisionID: 1,
					TransactionID:     "transaction-id02",
					Status:            PaymentStatusCaptured,
					MethodType:        PaymentMethodTypeCreditCard,
					Subtotal:          3000,
					Discount:          0,
					ShippingFee:       0,
					Tax:               300,
					Total:             3300,
				},
			},
			expect: map[string]*OrderPayment{
				"order-id01": {
					OrderID:           "order-id01",
					AddressRevisionID: 1,
					TransactionID:     "transaction-id01",
					Status:            PaymentStatusCaptured,
					MethodType:        PaymentMethodTypeCreditCard,
					Subtotal:          1980,
					Discount:          0,
					ShippingFee:       550,
					Tax:               253,
					Total:             2783,
				},
				"order-id02": {
					OrderID:           "order-id02",
					AddressRevisionID: 1,
					TransactionID:     "transaction-id02",
					Status:            PaymentStatusCaptured,
					MethodType:        PaymentMethodTypeCreditCard,
					Subtotal:          3000,
					Discount:          0,
					ShippingFee:       0,
					Tax:               300,
					Total:             3300,
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
