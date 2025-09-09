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
			methodType: entity.PaymentMethodTypeBankTransfer,
			expect:     PaymentMethodTypeBankTransfer,
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
			name:       "paidy",
			methodType: entity.PaymentMethodTypePaidy,
			expect:     PaymentMethodTypePaidy,
		},
		{
			name:       "pay easy",
			methodType: entity.PaymentMethodTypePayEasy,
			expect:     PaymentMethodTypePayEasy,
		},
		{
			name:       "free",
			methodType: entity.PaymentMethodTypeNone,
			expect:     PaymentMethodTypeFree,
		},
		{
			name:       "unknown",
			methodType: entity.PaymentMethodTypeUnknown,
			expect:     PaymentMethodTypeUnknown,
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
			methodType: PaymentMethodTypeBankTransfer,
			expect:     entity.PaymentMethodTypeBankTransfer,
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
			name:       "paidy",
			methodType: PaymentMethodTypePaidy,
			expect:     entity.PaymentMethodTypePaidy,
		},
		{
			name:       "pay easy",
			methodType: PaymentMethodTypePayEasy,
			expect:     entity.PaymentMethodTypePayEasy,
		},
		{
			name:       "none",
			methodType: PaymentMethodTypeFree,
			expect:     entity.PaymentMethodTypeNone,
		},
		{
			name:       "unknown",
			methodType: PaymentMethodTypeUnknown,
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.PaymentMethodType.Response())
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
			expect: PaymentStatusUnpaid,
		},
		{
			name:   "paid",
			status: entity.PaymentStatusCaptured,
			expect: PaymentStatusPaid,
		},
		{
			name:   "canceled",
			status: entity.PaymentStatusCanceled,
			expect: PaymentStatusCanceled,
		},
		{
			name:   "refunded",
			status: entity.PaymentStatusRefunded,
			expect: PaymentStatusCanceled,
		},
		{
			name:   "expired",
			status: entity.PaymentStatusFailed,
			expect: PaymentStatusFailed,
		},
		{
			name:   "unknown",
			status: entity.PaymentStatusUnknown,
			expect: PaymentStatusUnknown,
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
			status: PaymentStatusPaid,
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
					MethodType:    PaymentMethodTypeCreditCard.Response(),
					Status:        PaymentStatusPaid.Response(),
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
					MethodType:    PaymentMethodTypeCreditCard.Response(),
					Status:        PaymentStatusPaid.Response(),
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
				MethodType:    PaymentMethodTypeCreditCard.Response(),
				Status:        PaymentStatusPaid.Response(),
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
