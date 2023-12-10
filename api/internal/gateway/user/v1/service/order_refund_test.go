package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestRefundType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		typ    entity.RefundType
		expect RefundType
	}{
		{
			name:   "canceled",
			typ:    entity.RefundTypeCanceled,
			expect: RefundTypeCanceled,
		},
		{
			name:   "refunded",
			typ:    entity.RefundTypeRefunded,
			expect: RefundTypeRefunded,
		},
		{
			name:   "none",
			typ:    entity.RefundTypeNone,
			expect: RefundTypeNone,
		},
	}
	for _, tt := range tests {
		tt := tt
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
			typ:    RefundTypeNone,
			expect: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
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
		expect *OrderRefund
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
				Tax:               253,
				Total:             2783,
				RefundTotal:       0,
				RefundType:        entity.RefundTypeNone,
				RefundReason:      "",
				OrderedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
				PaidAt:            jst.Date(2022, 1, 1, 0, 0, 0, 0),
				RefundedAt:        time.Time{},
				CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &OrderRefund{
				OrderRefund: response.OrderRefund{
					Total:      0,
					Type:       RefundTypeNone.Response(),
					Reason:     "",
					Canceled:   false,
					CanceledAt: 0,
				},
				orderID: "order-id",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
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
		expect *response.OrderRefund
	}{
		{
			name: "success",
			refund: &OrderRefund{
				OrderRefund: response.OrderRefund{
					Total:      0,
					Type:       RefundTypeNone.Response(),
					Reason:     "",
					Canceled:   false,
					CanceledAt: 0,
				},
			},
			expect: &response.OrderRefund{
				Total:      0,
				Type:       RefundTypeNone.Response(),
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
		tt := tt
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
					Tax:               253,
					Total:             2783,
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
					OrderRefund: response.OrderRefund{
						Total:      0,
						Type:       RefundTypeNone.Response(),
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewOrderRefunds(tt.orders))
		})
	}
}
