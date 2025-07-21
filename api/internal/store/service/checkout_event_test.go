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
	"github.com/and-period/furumaru/api/internal/store/komoju"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNotifyPaymentAuthorized(t *testing.T) {
	t.Parallel()
	now := time.Now()
	order := func(methodType entity.PaymentMethodType) *entity.Order {
		return &entity.Order{
			ID:            "order-id",
			UserID:        "user-id",
			SessionID:     "session-id",
			PromotionID:   "",
			CoordinatorID: "coordinator-id",
			CreatedAt:     now,
			UpdatedAt:     now,
			OrderPayment: entity.OrderPayment{
				OrderID:           "order-id",
				AddressRevisionID: 1,
				TransactionID:     "transaction-id",
				Status:            entity.PaymentStatusPending,
				MethodType:        methodType,
				Subtotal:          1100,
				Discount:          0,
				ShippingFee:       500,
				Tax:               145,
				Total:             1600,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			OrderFulfillments: entity.OrderFulfillments{
				{
					ID:                "fulfillment-id",
					OrderID:           "order-id",
					AddressRevisionID: 1,
					Status:            entity.FulfillmentStatusUnfulfilled,
					TrackingNumber:    "",
					ShippingCarrier:   entity.ShippingCarrierUnknown,
					ShippingType:      entity.ShippingTypeNormal,
					BoxNumber:         1,
					BoxSize:           entity.ShippingSize60,
					CreatedAt:         now,
					UpdatedAt:         now,
				},
			},
			OrderItems: []*entity.OrderItem{
				{
					FulfillmentID:     "fufillment-id",
					ProductRevisionID: 1,
					OrderID:           "order-id",
					Quantity:          1,
					CreatedAt:         now,
					UpdatedAt:         now,
				},
				{
					FulfillmentID:     "fufillment-id",
					ProductRevisionID: 2,
					OrderID:           "order-id",
					Quantity:          2,
					CreatedAt:         now,
					UpdatedAt:         now,
				},
			},
		}
	}
	params := &database.UpdateOrderAuthorizedParams{
		PaymentID: "payment-id",
		IssuedAt:  now,
	}
	payment := &komoju.PaymentResponse{
		PaymentInfo: &komoju.PaymentInfo{
			ID:     "payment-id",
			Status: komoju.PaymentStatusAuthorized,
		},
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.NotifyPaymentAuthorizedInput
		expect error
	}{
		{
			name: "success when immediate payment",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().
					Get(ctx, "order-id").
					Return(order(entity.PaymentMethodTypeCreditCard), nil)
				mocks.db.Order.EXPECT().UpdateAuthorized(ctx, "order-id", params).Return(nil)
				mocks.komojuPayment.EXPECT().Show(ctx, "payment-id").Return(payment, nil)
				mocks.komojuPayment.EXPECT().
					Capture(ctx, "payment-id").
					Return(&komoju.PaymentResponse{}, nil)
				mocks.cache.EXPECT().
					Get(gomock.Any(), &entity.Cart{SessionID: "session-id"}).
					Return(assert.AnError)
			},
			input: &store.NotifyPaymentAuthorizedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					IssuedAt:  now,
					Status:    entity.PaymentStatusAuthorized,
				},
			},
			expect: nil,
		},
		{
			name: "success when not immediate payment",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().
					Get(ctx, "order-id").
					Return(order(entity.PaymentMethodTypeKonbini), nil)
			},
			input: &store.NotifyPaymentAuthorizedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					IssuedAt:  now,
					Status:    entity.PaymentStatusAuthorized,
				},
			},
			expect: nil,
		},
		{
			name: "success when already captured",
			setup: func(ctx context.Context, mocks *mocks) {
				payment := &komoju.PaymentResponse{
					PaymentInfo: &komoju.PaymentInfo{
						ID:     "payment-id",
						Status: komoju.PaymentStatusCaptured,
					},
				}
				mocks.db.Order.EXPECT().
					Get(ctx, "order-id").
					Return(order(entity.PaymentMethodTypeCreditCard), nil)
				mocks.db.Order.EXPECT().
					UpdateAuthorized(ctx, "order-id", params).
					Return(database.ErrFailedPrecondition)
				mocks.komojuPayment.EXPECT().Show(ctx, "payment-id").Return(payment, nil)
			},
			input: &store.NotifyPaymentAuthorizedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					IssuedAt:  now,
					Status:    entity.PaymentStatusAuthorized,
				},
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.NotifyPaymentAuthorizedInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid status",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.NotifyPaymentAuthorizedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					IssuedAt:  now,
					Status:    entity.PaymentStatusCaptured,
				},
			},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get order",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(nil, assert.AnError)
			},
			input: &store.NotifyPaymentAuthorizedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					IssuedAt:  now,
					Status:    entity.PaymentStatusAuthorized,
				},
			},
			expect: exception.ErrInternal,
		},
		{
			name: "failed to update payment status",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().
					Get(ctx, "order-id").
					Return(order(entity.PaymentMethodTypeCreditCard), nil)
				mocks.db.Order.EXPECT().
					UpdateAuthorized(ctx, "order-id", params).
					Return(assert.AnError)
			},
			input: &store.NotifyPaymentAuthorizedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					IssuedAt:  now,
					Status:    entity.PaymentStatusAuthorized,
				},
			},
			expect: exception.ErrInternal,
		},
		{
			name: "failed to show payment",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().
					Get(ctx, "order-id").
					Return(order(entity.PaymentMethodTypeCreditCard), nil)
				mocks.db.Order.EXPECT().
					UpdateAuthorized(ctx, "order-id", params).
					Return(database.ErrFailedPrecondition)
				mocks.komojuPayment.EXPECT().Show(ctx, "payment-id").Return(nil, assert.AnError)
			},
			input: &store.NotifyPaymentAuthorizedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					IssuedAt:  now,
					Status:    entity.PaymentStatusAuthorized,
				},
			},
			expect: exception.ErrInternal,
		},
		{
			name: "failed to capture",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().
					Get(ctx, "order-id").
					Return(order(entity.PaymentMethodTypeCreditCard), nil)
				mocks.db.Order.EXPECT().
					UpdateAuthorized(ctx, "order-id", params).
					Return(database.ErrFailedPrecondition)
				mocks.komojuPayment.EXPECT().Show(ctx, "payment-id").Return(payment, nil)
				mocks.komojuPayment.EXPECT().Capture(ctx, "payment-id").Return(nil, assert.AnError)
			},
			input: &store.NotifyPaymentAuthorizedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					IssuedAt:  now,
					Status:    entity.PaymentStatusAuthorized,
				},
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				err := service.NotifyPaymentAuthorized(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expect)
			}),
		)
	}
}

func TestNotifyPaymentCaptured(t *testing.T) {
	t.Parallel()
	now := time.Now()
	order := &entity.Order{
		ID:            "order-id",
		UserID:        "user-id",
		SessionID:     "session-id",
		PromotionID:   "",
		CoordinatorID: "coordinator-id",
		CreatedAt:     now,
		UpdatedAt:     now,
		OrderPayment: entity.OrderPayment{
			OrderID:           "order-id",
			AddressRevisionID: 1,
			TransactionID:     "transaction-id",
			Status:            entity.PaymentStatusPending,
			MethodType:        entity.PaymentMethodTypeCreditCard,
			Subtotal:          1100,
			Discount:          0,
			ShippingFee:       500,
			Tax:               145,
			Total:             1600,
			CreatedAt:         now,
			UpdatedAt:         now,
		},
		OrderFulfillments: entity.OrderFulfillments{
			{
				ID:                "fulfillment-id",
				OrderID:           "order-id",
				AddressRevisionID: 1,
				Status:            entity.FulfillmentStatusUnfulfilled,
				TrackingNumber:    "",
				ShippingCarrier:   entity.ShippingCarrierUnknown,
				ShippingType:      entity.ShippingTypeNormal,
				BoxNumber:         1,
				BoxSize:           entity.ShippingSize60,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
		OrderItems: []*entity.OrderItem{
			{
				FulfillmentID:     "fufillment-id",
				ProductRevisionID: 1,
				OrderID:           "order-id",
				Quantity:          1,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			{
				FulfillmentID:     "fufillment-id",
				ProductRevisionID: 2,
				OrderID:           "order-id",
				Quantity:          2,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
	}
	params := &database.UpdateOrderCapturedParams{
		PaymentID: "payment-id",
		IssuedAt:  now,
	}
	in := &messenger.NotifyOrderCapturedInput{
		OrderID: "order-id",
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.NotifyPaymentCapturedInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.db.Order.EXPECT().UpdateCaptured(ctx, "order-id", params).Return(nil)
				mocks.messenger.EXPECT().
					NotifyOrderCaptured(gomock.Any(), in).
					Return(assert.AnError)
			},
			input: &store.NotifyPaymentCapturedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					Status:    entity.PaymentStatusCaptured,
					IssuedAt:  now,
				},
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.NotifyPaymentCapturedInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid status",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.NotifyPaymentCapturedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					Status:    entity.PaymentStatusAuthorized,
					IssuedAt:  now,
				},
			},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get order",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(nil, assert.AnError)
			},
			input: &store.NotifyPaymentCapturedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					Status:    entity.PaymentStatusCaptured,
					IssuedAt:  now,
				},
			},
			expect: exception.ErrInternal,
		},
		{
			name: "already updated",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.db.Order.EXPECT().
					UpdateCaptured(ctx, "order-id", params).
					Return(database.ErrFailedPrecondition)
			},
			input: &store.NotifyPaymentCapturedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					Status:    entity.PaymentStatusCaptured,
					IssuedAt:  now,
				},
			},
			expect: nil,
		},
		{
			name: "failed to update payment status",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.db.Order.EXPECT().
					UpdateCaptured(ctx, "order-id", params).
					Return(assert.AnError)
			},
			input: &store.NotifyPaymentCapturedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					Status:    entity.PaymentStatusCaptured,
					IssuedAt:  now,
				},
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				err := service.NotifyPaymentCaptured(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expect)
			}),
		)
	}
}

func TestNotifyPaymentFailed(t *testing.T) {
	t.Parallel()
	now := time.Now()
	order := &entity.Order{
		ID:            "order-id",
		UserID:        "user-id",
		SessionID:     "session-id",
		PromotionID:   "",
		CoordinatorID: "coordinator-id",
		CreatedAt:     now,
		UpdatedAt:     now,
		OrderPayment: entity.OrderPayment{
			OrderID:           "order-id",
			AddressRevisionID: 1,
			TransactionID:     "transaction-id",
			Status:            entity.PaymentStatusPending,
			MethodType:        entity.PaymentMethodTypeCreditCard,
			Subtotal:          1100,
			Discount:          0,
			ShippingFee:       500,
			Tax:               145,
			Total:             1600,
			CreatedAt:         now,
			UpdatedAt:         now,
		},
		OrderFulfillments: entity.OrderFulfillments{
			{
				ID:                "fulfillment-id",
				OrderID:           "order-id",
				AddressRevisionID: 1,
				Status:            entity.FulfillmentStatusUnfulfilled,
				TrackingNumber:    "",
				ShippingCarrier:   entity.ShippingCarrierUnknown,
				ShippingType:      entity.ShippingTypeNormal,
				BoxNumber:         1,
				BoxSize:           entity.ShippingSize60,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
		OrderItems: []*entity.OrderItem{
			{
				FulfillmentID:     "fufillment-id",
				ProductRevisionID: 1,
				OrderID:           "order-id",
				Quantity:          1,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			{
				FulfillmentID:     "fufillment-id",
				ProductRevisionID: 2,
				OrderID:           "order-id",
				Quantity:          2,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
	}
	params := &database.UpdateOrderFailedParams{
		PaymentID: "payment-id",
		Status:    entity.PaymentStatusFailed,
		IssuedAt:  now,
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.NotifyPaymentFailedInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.db.Order.EXPECT().UpdateFailed(ctx, "order-id", params).Return(nil)
				mocks.db.Product.EXPECT().
					DecreaseInventory(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(2)
			},
			input: &store.NotifyPaymentFailedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					Status:    entity.PaymentStatusFailed,
					IssuedAt:  now,
				},
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.NotifyPaymentFailedInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid status",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.NotifyPaymentFailedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					Status:    entity.PaymentStatusAuthorized,
					IssuedAt:  now,
				},
			},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get order",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(nil, assert.AnError)
			},
			input: &store.NotifyPaymentFailedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					Status:    entity.PaymentStatusFailed,
					IssuedAt:  now,
				},
			},
			expect: exception.ErrInternal,
		},
		{
			name: "failed to update payment status",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.db.Order.EXPECT().UpdateFailed(ctx, "order-id", params).Return(assert.AnError)
			},
			input: &store.NotifyPaymentFailedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					Status:    entity.PaymentStatusFailed,
					IssuedAt:  now,
				},
			},
			expect: exception.ErrInternal,
		},
		{
			name: "already updated",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.db.Order.EXPECT().
					UpdateFailed(ctx, "order-id", params).
					Return(database.ErrFailedPrecondition)
			},
			input: &store.NotifyPaymentFailedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					Status:    entity.PaymentStatusFailed,
					IssuedAt:  now,
				},
			},
			expect: nil,
		},
		{
			name: "failed to increase product inventory",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.db.Order.EXPECT().UpdateFailed(ctx, "order-id", params).Return(nil)
				mocks.db.Product.EXPECT().
					DecreaseInventory(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(assert.AnError).
					MinTimes(1)
			},
			input: &store.NotifyPaymentFailedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:   "order-id",
					PaymentID: "payment-id",
					Status:    entity.PaymentStatusFailed,
					IssuedAt:  now,
				},
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				err := service.NotifyPaymentFailed(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expect)
			}),
		)
	}
}

func TestNotifyPaymentRefunded(t *testing.T) {
	t.Parallel()
	now := time.Now()
	params := &database.UpdateOrderRefundedParams{
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
				mocks.db.Order.EXPECT().UpdateRefunded(ctx, "order-id", params).Return(nil)
			},
			input: &store.NotifyPaymentRefundedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:  "order-id",
					Status:   entity.PaymentStatusRefunded,
					IssuedAt: now,
				},
				Type:   entity.RefundTypeRefunded,
				Total:  1980,
				Reason: "在庫不足のため。",
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
				mocks.db.Order.EXPECT().
					UpdateRefunded(ctx, "order-id", params).
					Return(assert.AnError)
			},
			input: &store.NotifyPaymentRefundedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:  "order-id",
					Status:   entity.PaymentStatusRefunded,
					IssuedAt: now,
				},
				Type:   entity.RefundTypeRefunded,
				Total:  1980,
				Reason: "在庫不足のため。",
			},
			expect: exception.ErrInternal,
		},
		{
			name: "already updated",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().
					UpdateRefunded(ctx, "order-id", params).
					Return(database.ErrFailedPrecondition)
			},
			input: &store.NotifyPaymentRefundedInput{
				NotifyPaymentPayload: store.NotifyPaymentPayload{
					OrderID:  "order-id",
					Status:   entity.PaymentStatusRefunded,
					IssuedAt: now,
				},
				Type:   entity.RefundTypeRefunded,
				Total:  1980,
				Reason: "在庫不足のため。",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				err := service.NotifyPaymentRefunded(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expect)
			}),
		)
	}
}
