package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/stretchr/testify/assert"
)

func TestShippingSize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.ShippingSize
		expect ShippingSize
	}{
		{
			name:   "size 60",
			status: entity.ShippingSize60,
			expect: ShippingSize60,
		},
		{
			name:   "size 80",
			status: entity.ShippingSize80,
			expect: ShippingSize80,
		},
		{
			name:   "size 100",
			status: entity.ShippingSize100,
			expect: ShippingSize100,
		},
		{
			name:   "unknown",
			status: entity.ShippingSizeUnknown,
			expect: ShippingSizeUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewShippingSize(tt.status))
		})
	}
}

func TestShippingSize_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status ShippingSize
		expect int32
	}{
		{
			name:   "success",
			status: ShippingSize60,
			expect: 1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.status.Response())
		})
	}
}

func TestShippingType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.ShippingType
		expect ShippingType
	}{
		{
			name:   "normal",
			status: entity.ShippingTypeNormal,
			expect: ShippingTypeNormal,
		},
		{
			name:   "frozen",
			status: entity.ShippingTypeFrozen,
			expect: ShippingTypeFrozen,
		},
		{
			name:   "unknown",
			status: entity.ShippingTypeUnknown,
			expect: ShippingTypeUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewShippingType(tt.status))
		})
	}
}

func TestShippingType_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status ShippingType
		expect int32
	}{
		{
			name:   "success",
			status: ShippingTypeNormal,
			expect: 1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.status.Response())
		})
	}
}

func TestPaymentMethodType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		methodType entity.PaymentMethodType
		expect     PaymentMethodType
	}{
		{
			name:       "cash",
			methodType: entity.PaymentMethodTypeCash,
			expect:     PaymentMethodTypeCash,
		},
		{
			name:       "credit card",
			methodType: entity.PaymentMethodTypeCreditCard,
			expect:     PaymentMethodTypeCreditCard,
		},
		{
			name:       "konbini",
			methodType: entity.PaymentMethodTypeKonbini,
			expect:     PaymentMethodTypeKonbini,
		},
		{
			name:       "bank transfer",
			methodType: entity.PaymentMethodTypeBankTranser,
			expect:     PaymentMethodTypeBankTranser,
		},
		{
			name:       "paypay",
			methodType: entity.PaymentMethodTypePayPay,
			expect:     PaymentMethodTypePayPay,
		},
		{
			name:       "line pay",
			methodType: entity.PaymentMethodTypeLinePay,
			expect:     PaymentMethodTypeLinePay,
		},
		{
			name:       "merpay",
			methodType: entity.PaymentMethodTypeMerpay,
			expect:     PaymentMethodTypeMerpay,
		},
		{
			name:       "rakuten pay",
			methodType: entity.PaymentMethodTypeRakutenPay,
			expect:     PaymentMethodTypeRakutenPay,
		},
		{
			name:       "au pay",
			methodType: entity.PaymentMethodTypeAUPay,
			expect:     PaymentMethodTypeAUPay,
		},
		{
			name:       "unknown",
			methodType: entity.PaymentMethodTypeUnknown,
			expect:     PaymentMethodTypeUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewPaymentMethodType(tt.methodType))
		})
	}
}

func TestPaymentMethodType_StoreEntity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		methodType PaymentMethodType
		expect     entity.PaymentMethodType
	}{
		{
			name:       "cash",
			methodType: PaymentMethodTypeCash,
			expect:     entity.PaymentMethodTypeCash,
		},
		{
			name:       "credit card",
			methodType: PaymentMethodTypeCreditCard,
			expect:     entity.PaymentMethodTypeCreditCard,
		},
		{
			name:       "konbini",
			methodType: PaymentMethodTypeKonbini,
			expect:     entity.PaymentMethodTypeKonbini,
		},
		{
			name:       "bank transfer",
			methodType: PaymentMethodTypeBankTranser,
			expect:     entity.PaymentMethodTypeBankTranser,
		},
		{
			name:       "paypay",
			methodType: PaymentMethodTypePayPay,
			expect:     entity.PaymentMethodTypePayPay,
		},
		{
			name:       "line pay",
			methodType: PaymentMethodTypeLinePay,
			expect:     entity.PaymentMethodTypeLinePay,
		},
		{
			name:       "merpay",
			methodType: PaymentMethodTypeMerpay,
			expect:     entity.PaymentMethodTypeMerpay,
		},
		{
			name:       "rakuten pay",
			methodType: PaymentMethodTypeRakutenPay,
			expect:     entity.PaymentMethodTypeRakutenPay,
		},
		{
			name:       "au pay",
			methodType: PaymentMethodTypeAUPay,
			expect:     entity.PaymentMethodTypeAUPay,
		},
		{
			name:       "unknown",
			methodType: PaymentMethodTypeUnknown,
			expect:     entity.PaymentMethodTypeUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.methodType.StoreEntity())
		})
	}
}

func TestPaymentMethodType_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		PaymentMethodType PaymentMethodType
		expect            int32
	}{
		{
			name:              "success",
			PaymentMethodType: PaymentMethodTypeCreditCard,
			expect:            2,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.PaymentMethodType.Response())
		})
	}
}
