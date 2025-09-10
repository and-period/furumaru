package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestRefundType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		typ    entity.RefundType
		expect types.RefundType
	}{
		{
			name:   "canceled",
			typ:    entity.RefundTypeCanceled,
			expect: types.RefundTypeCanceled,
		},
		{
			name:   "refunded",
			typ:    entity.RefundTypeRefunded,
			expect: types.RefundTypeRefunded,
		},
		{
			name:   "none",
			typ:    entity.RefundTypeNone,
			expect: types.RefundTypeNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewRefundType(tt.typ))
		})
	}
}

func TestRefundType_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		typ    RefundType
		expect int32
	}{
		{
			name:   "none",
			typ:    RefundType(types.RefundTypeNone),
			expect: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.typ.Response())
		})
	}
}

func TestOrderRefund(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		order  *entity.OrderPayment
		expect *types.OrderRefund
	}{
		{
			name: "success",
			order: &entity.OrderPayment{
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
			expect: &types.OrderRefund{
				Total:      0,
				Type:       NewRefundType(entity.RefundTypeNone).Response(),
				Reason:     "",
				Canceled:   false,
				CanceledAt: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewOrderRefund(tt.order))
		})
	}
}

func TestOrderRefund_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		refund *OrderRefund
		expect *types.OrderRefund
	}{
		{
			name: "success",
			refund: &OrderRefund{
				OrderRefund: types.OrderRefund{
				Total:      0,
				Type:       NewRefundType(entity.RefundTypeNone).Response(),
				Reason:     "",
				Canceled:   false,
				CanceledAt: 0,
				},
			},
			expect: &types.OrderRefund{
				Total:      0,
				Type:       NewRefundType(entity.RefundTypeNone).Response(),
				Reason:     "",
				Canceled:   false,
				CanceledAt: 0,
			},
		},
		{
			name:   "empty",
			refund: nil,
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.refund.Response())
		})
	}
}

func TestOrderRefunds(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders entity.OrderPayments
		expect OrderRefunds
	}{
		{
			name: "success",
			orders: entity.OrderPayments{
				{
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
			},
			expect: OrderRefunds{
				{
					OrderRefund: types.OrderRefund{
						Total:      0,
						Type:       NewRefundType(entity.RefundTypeNone).Response(),
						Reason:     "",
						Canceled:   false,
						CanceledAt: 0,
					},
					orderID: "order-id",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewOrderRefunds(tt.orders))
		})
	}
}
