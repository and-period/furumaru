package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestPaymentMethodType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		methodType entity.PaymentMethodType
		expect     types.PaymentMethodType
	}{
		{
			name:       "cash",
			methodType: entity.PaymentMethodTypeCash,
			expect:     types.PaymentMethodTypeCash,
		},
		{
			name:       "credit card",
			methodType: entity.PaymentMethodTypeCreditCard,
			expect:     types.PaymentMethodTypeCreditCard,
		},
		{
			name:       "konbini",
			methodType: entity.PaymentMethodTypeKonbini,
			expect:     types.PaymentMethodTypeKonbini,
		},
		{
			name:       "bank transfer",
			methodType: entity.PaymentMethodTypeBankTransfer,
			expect:     types.PaymentMethodTypeBankTransfer,
		},
		{
			name:       "paypay",
			methodType: entity.PaymentMethodTypePayPay,
			expect:     types.PaymentMethodTypePayPay,
		},
		{
			name:       "line pay",
			methodType: entity.PaymentMethodTypeLinePay,
			expect:     types.PaymentMethodTypeLinePay,
		},
		{
			name:       "merpay",
			methodType: entity.PaymentMethodTypeMerpay,
			expect:     types.PaymentMethodTypeMerpay,
		},
		{
			name:       "rakuten pay",
			methodType: entity.PaymentMethodTypeRakutenPay,
			expect:     types.PaymentMethodTypeRakutenPay,
		},
		{
			name:       "au pay",
			methodType: entity.PaymentMethodTypeAUPay,
			expect:     types.PaymentMethodTypeAUPay,
		},
		{
			name:       "paidy",
			methodType: entity.PaymentMethodTypePaidy,
			expect:     types.PaymentMethodTypePaidy,
		},
		{
			name:       "pay easy",
			methodType: entity.PaymentMethodTypePayEasy,
			expect:     types.PaymentMethodTypePayEasy,
		},
		{
			name:       "free",
			methodType: entity.PaymentMethodTypeNone,
			expect:     types.PaymentMethodTypeFree,
		},
		{
			name:       "unknown",
			methodType: entity.PaymentMethodTypeUnknown,
			expect:     types.PaymentMethodTypeUnknown,
		},
	}
	for _, tt := range tests {
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
			methodType: PaymentMethodType(types.PaymentMethodTypeCash),
			expect:     entity.PaymentMethodTypeCash,
		},
		{
			name:       "credit card",
			methodType: PaymentMethodType(types.PaymentMethodTypeCreditCard),
			expect:     entity.PaymentMethodTypeCreditCard,
		},
		{
			name:       "konbini",
			methodType: PaymentMethodType(types.PaymentMethodTypeKonbini),
			expect:     entity.PaymentMethodTypeKonbini,
		},
		{
			name:       "bank transfer",
			methodType: PaymentMethodType(types.PaymentMethodTypeBankTransfer),
			expect:     entity.PaymentMethodTypeBankTransfer,
		},
		{
			name:       "paypay",
			methodType: PaymentMethodType(types.PaymentMethodTypePayPay),
			expect:     entity.PaymentMethodTypePayPay,
		},
		{
			name:       "line pay",
			methodType: PaymentMethodType(types.PaymentMethodTypeLinePay),
			expect:     entity.PaymentMethodTypeLinePay,
		},
		{
			name:       "merpay",
			methodType: PaymentMethodType(types.PaymentMethodTypeMerpay),
			expect:     entity.PaymentMethodTypeMerpay,
		},
		{
			name:       "rakuten pay",
			methodType: PaymentMethodType(types.PaymentMethodTypeRakutenPay),
			expect:     entity.PaymentMethodTypeRakutenPay,
		},
		{
			name:       "au pay",
			methodType: PaymentMethodType(types.PaymentMethodTypeAUPay),
			expect:     entity.PaymentMethodTypeAUPay,
		},
		{
			name:       "paidy",
			methodType: PaymentMethodType(types.PaymentMethodTypePaidy),
			expect:     entity.PaymentMethodTypePaidy,
		},
		{
			name:       "pay easy",
			methodType: PaymentMethodType(types.PaymentMethodTypePayEasy),
			expect:     entity.PaymentMethodTypePayEasy,
		},
		{
			name:       "none",
			methodType: PaymentMethodType(types.PaymentMethodTypeFree),
			expect:     entity.PaymentMethodTypeNone,
		},
		{
			name:       "unknown",
			methodType: PaymentMethodType(types.PaymentMethodTypeUnknown),
			expect:     entity.PaymentMethodTypeUnknown,
		},
	}
	for _, tt := range tests {
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
		methodType        PaymentMethodType
		expect            int32
	}{
		{
			name:              "success",
			methodType: PaymentMethodType(types.PaymentMethodTypeCreditCard),
			expect:            2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.methodType.Response())
		})
	}
}

func TestPaymentStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.PaymentStatus
		expect PaymentStatus
	}{
		{
			name:   "pending",
			status: entity.PaymentStatusPending,
			expect: PaymentStatus(types.PaymentStatusUnpaid),
		},
		{
			name:   "paid",
			status: entity.PaymentStatusCaptured,
			expect: PaymentStatus(types.PaymentStatusPaid),
		},
		{
			name:   "canceled",
			status: entity.PaymentStatusCanceled,
			expect: PaymentStatus(types.PaymentStatusCanceled),
		},
		{
			name:   "refunded",
			status: entity.PaymentStatusRefunded,
			expect: PaymentStatus(types.PaymentStatusCanceled),
		},
		{
			name:   "expired",
			status: entity.PaymentStatusFailed,
			expect: PaymentStatus(types.PaymentStatusFailed),
		},
		{
			name:   "unknown",
			status: entity.PaymentStatusUnknown,
			expect: PaymentStatus(types.PaymentStatusUnknown),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewPaymentStatus(tt.status))
		})
	}
}

func TestPaymentStatus_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status PaymentStatus
		expect int32
	}{
		{
			name:   "success",
			status: PaymentStatus(types.PaymentStatusPaid),
			expect: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.status.Response())
		})
	}
}

func TestOrderPayment(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payment *entity.OrderPayment
		expect  *OrderPayment
	}{
		{
			name: "success",
			payment: &entity.OrderPayment{
				OrderID:           "order-id",
				AddressRevisionID: 1,
				TransactionID:     "transaction-id",
				Status:            entity.PaymentStatusCaptured,
				MethodType:        entity.PaymentMethodTypeCreditCard,
				Subtotal:          1980,
				Discount:          0,
				ShippingFee:       550,
				Tax:               230,
				Total:             2530,
				RefundTotal:       0,
				RefundType:        entity.RefundTypeNone,
				RefundReason:      "",
				OrderedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
				PaidAt:            jst.Date(2022, 1, 1, 0, 0, 0, 0),
				RefundedAt:        time.Time{},
				CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &OrderPayment{
				OrderPayment: types.OrderPayment{
					TransactionID: "transaction-id",
					MethodType:    NewPaymentMethodType(entity.PaymentMethodTypeCreditCard).Response(),
					Status:        NewPaymentStatus(entity.PaymentStatusCaptured).Response(),
					Subtotal:      1980,
					Discount:      0,
					ShippingFee:   550,
					Total:         2530,
					OrderedAt:     1640962800,
					PaidAt:        1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewOrderPayment(tt.payment))
		})
	}
}

func TestOrderPayment_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payment *OrderPayment
		expect  *types.OrderPayment
	}{
		{
			name: "success",
			payment: &OrderPayment{
				OrderPayment: types.OrderPayment{
				TransactionID: "transaction-id",
				MethodType:    NewPaymentMethodType(entity.PaymentMethodTypeCreditCard).Response(),
				Status:        NewPaymentStatus(entity.PaymentStatusCaptured).Response(),
				Subtotal:      1100,
				Discount:      0,
				ShippingFee:   500,
				Total:         1600,
				OrderedAt:     1640962800,
				PaidAt:        1640962800,
				},
			},
			expect: &types.OrderPayment{
				TransactionID: "transaction-id",
				MethodType:    NewPaymentMethodType(entity.PaymentMethodTypeCreditCard).Response(),
				Status:        NewPaymentStatus(entity.PaymentStatusCaptured).Response(),
				Subtotal:      1100,
				Discount:      0,
				ShippingFee:   500,
				Total:         1600,
				OrderedAt:     1640962800,
				PaidAt:        1640962800,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.payment.Response())
		})
	}
}
