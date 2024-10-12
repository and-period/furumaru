package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/stretchr/testify/assert"
)

func TestNotifyPaymentCompleted(t *testing.T) {
	t.Parallel()
	now := time.Now()
	params := &database.UpdateOrderPaymentParams{
		Status:    entity.PaymentStatusCaptured,
		PaymentID: "payment-id",
		IssuedAt:  now,
	}
	in := &messenger.NotifyOrderAuthorizedInput{
		OrderID: "order-id",
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.NotifyPaymentCompletedInput
		expect error
	}{
		{
			name: "success captured",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().UpdatePayment(ctx, "order-id", params).Return(nil)
				mocks.messenger.EXPECT().NotifyOrderAuthorized(ctx, in).Return(nil)
			},
			input: &store.NotifyPaymentCompletedInput{
				OrderID:   "order-id",
				PaymentID: "payment-id",
				Status:    entity.PaymentStatusCaptured,
				IssuedAt:  now,
			},
			expect: nil,
		},
		{
			name: "success canceled",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &database.UpdateOrderPaymentParams{
					Status:    entity.PaymentStatusCanceled,
					PaymentID: "payment-id",
					IssuedAt:  now,
				}
				mocks.db.Order.EXPECT().UpdatePayment(ctx, "order-id", params).Return(nil)
			},
			input: &store.NotifyPaymentCompletedInput{
				OrderID:   "order-id",
				PaymentID: "payment-id",
				Status:    entity.PaymentStatusCanceled,
				IssuedAt:  now,
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.NotifyPaymentCompletedInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update payment status",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().UpdatePayment(ctx, "order-id", params).Return(assert.AnError)
			},
			input: &store.NotifyPaymentCompletedInput{
				OrderID:   "order-id",
				PaymentID: "payment-id",
				Status:    entity.PaymentStatusCaptured,
				IssuedAt:  now,
			},
			expect: exception.ErrInternal,
		},
		{
			name: "already updated",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().UpdatePayment(ctx, "order-id", params).Return(database.ErrFailedPrecondition)
			},
			input: &store.NotifyPaymentCompletedInput{
				OrderID:   "order-id",
				PaymentID: "payment-id",
				Status:    entity.PaymentStatusCaptured,
				IssuedAt:  now,
			},
			expect: nil,
		},
		{
			name: "failed to notify order authorized",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().UpdatePayment(ctx, "order-id", params).Return(nil)
				mocks.messenger.EXPECT().NotifyOrderAuthorized(ctx, in).Return(assert.AnError)
			},
			input: &store.NotifyPaymentCompletedInput{
				OrderID:   "order-id",
				PaymentID: "payment-id",
				Status:    entity.PaymentStatusCaptured,
				IssuedAt:  now,
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.NotifyPaymentCompleted(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestNotifyPaymentRefunded(t *testing.T) {
	t.Parallel()
	now := time.Now()
	params := &database.UpdateOrderRefundParams{
		Status:       entity.PaymentStatusRefunded,
		RefundType:   entity.RefundTypeRefunded,
		RefundTotal:  1980,
		RefundReason: "在庫不足のため。",
		IssuedAt:     now,
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.NotifyPaymentRefundedInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().UpdateRefund(ctx, "order-id", params).Return(nil)
			},
			input: &store.NotifyPaymentRefundedInput{
				OrderID:  "order-id",
				Status:   entity.PaymentStatusRefunded,
				Type:     entity.RefundTypeRefunded,
				Total:    1980,
				Reason:   "在庫不足のため。",
				IssuedAt: now,
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.NotifyPaymentRefundedInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update payment status",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().UpdateRefund(ctx, "order-id", params).Return(assert.AnError)
			},
			input: &store.NotifyPaymentRefundedInput{
				OrderID:  "order-id",
				Status:   entity.PaymentStatusRefunded,
				Type:     entity.RefundTypeRefunded,
				Total:    1980,
				Reason:   "在庫不足のため。",
				IssuedAt: now,
			},
			expect: exception.ErrInternal,
		},
		{
			name: "already updated",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().UpdateRefund(ctx, "order-id", params).Return(database.ErrFailedPrecondition)
			},
			input: &store.NotifyPaymentRefundedInput{
				OrderID:  "order-id",
				Status:   entity.PaymentStatusRefunded,
				Type:     entity.RefundTypeRefunded,
				Total:    1980,
				Reason:   "在庫不足のため。",
				IssuedAt: now,
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.NotifyPaymentRefunded(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}
