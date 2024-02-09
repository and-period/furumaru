package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/komoju"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListOrders(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 10, 10, 18, 30, 0, 0)
	params := &database.ListOrdersParams{
		CoordinatorID: "coordinator-id",
		Limit:         30,
		Offset:        0,
	}
	orders := entity.Orders{
		{
			ID:            "order-id",
			UserID:        "user-id",
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
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListOrdersInput
		expect      entity.Orders
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().List(gomock.Any(), params).Return(orders, nil)
				mocks.db.Order.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListOrdersInput{
				CoordinatorID: "coordinator-id",
				Limit:         30,
				Offset:        0,
			},
			expect:      orders,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &store.ListOrdersInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list orders",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Order.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListOrdersInput{
				CoordinatorID: "coordinator-id",
				Limit:         30,
				Offset:        0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count orders",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().List(gomock.Any(), params).Return(orders, nil)
				mocks.db.Order.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &store.ListOrdersInput{
				CoordinatorID: "coordinator-id",
				Limit:         30,
				Offset:        0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListOrders(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestGetOrder(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 10, 10, 18, 30, 0, 0)
	order := &entity.Order{
		ID:            "order-id",
		UserID:        "user-id",
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

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetOrderInput
		expect    *entity.Order
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
			},
			input: &store.GetOrderInput{
				OrderID: "order-id",
			},
			expect:    order,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetOrderInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get order",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(nil, assert.AnError)
			},
			input: &store.GetOrderInput{
				OrderID: "order-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetOrder(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestGetOrderByTransactionID(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 10, 10, 18, 30, 0, 0)
	order := &entity.Order{
		ID:            "order-id",
		UserID:        "user-id",
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

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetOrderByTransactionIDInput
		expect    *entity.Order
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().GetByTransactionID(ctx, "user-id", "transaction-id").Return(order, nil)
			},
			input: &store.GetOrderByTransactionIDInput{
				UserID:        "user-id",
				TransactionID: "transaction-id",
			},
			expect:    order,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetOrderByTransactionIDInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get order",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().GetByTransactionID(ctx, "user-id", "transaction-id").Return(nil, assert.AnError)
			},
			input: &store.GetOrderByTransactionIDInput{
				UserID:        "user-id",
				TransactionID: "transaction-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetOrderByTransactionID(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCaptureOrder(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 10, 10, 18, 30, 0, 0)
	order := &entity.Order{
		ID:            "order-id",
		UserID:        "user-id",
		PromotionID:   "",
		CoordinatorID: "coordinator-id",
		CreatedAt:     now,
		UpdatedAt:     now,
		OrderPayment: entity.OrderPayment{
			OrderID:           "order-id",
			AddressRevisionID: 1,
			TransactionID:     "transaction-id",
			PaymentID:         "payment-id",
			Status:            entity.PaymentStatusAuthorized,
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
	payment := &komoju.PaymentResponse{}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.CaptureOrderInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.komojuPayment.EXPECT().Capture(ctx, "payment-id").Return(payment, nil)
			},
			input: &store.CaptureOrderInput{
				OrderID: "order-id",
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.CaptureOrderInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get order",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(nil, assert.AnError)
			},
			input: &store.CaptureOrderInput{
				OrderID: "order-id",
			},
			expect: exception.ErrInternal,
		},
		{
			name: "failed to capture",
			setup: func(ctx context.Context, mocks *mocks) {
				order := &entity.Order{
					OrderPayment: entity.OrderPayment{Status: entity.PaymentStatusCaptured},
				}
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
			},
			input: &store.CaptureOrderInput{
				OrderID: "order-id",
			},
			expect: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to capture",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.komojuPayment.EXPECT().Capture(ctx, "payment-id").Return(nil, assert.AnError)
			},
			input: &store.CaptureOrderInput{
				OrderID: "order-id",
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.CaptureOrder(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestDraftOrder(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 10, 10, 18, 30, 0, 0)
	order := &entity.Order{
		ID:            "order-id",
		UserID:        "user-id",
		PromotionID:   "",
		CoordinatorID: "coordinator-id",
		CreatedAt:     now,
		UpdatedAt:     now,
		OrderPayment: entity.OrderPayment{
			OrderID:           "order-id",
			AddressRevisionID: 1,
			TransactionID:     "transaction-id",
			PaymentID:         "payment-id",
			Status:            entity.PaymentStatusCaptured,
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
	params := &database.DraftOrderParams{
		ShippingMessage: "購入ありがとうございます。",
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.DraftOrderInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.db.Order.EXPECT().Draft(ctx, "order-id", params).Return(nil)
			},
			input: &store.DraftOrderInput{
				OrderID:         "order-id",
				ShippingMessage: "購入ありがとうございます。",
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.DraftOrderInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get order",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(nil, assert.AnError)
			},
			input: &store.DraftOrderInput{
				OrderID:         "order-id",
				ShippingMessage: "購入ありがとうございます。",
			},
			expect: exception.ErrInternal,
		},
		{
			name: "failed to preservable",
			setup: func(ctx context.Context, mocks *mocks) {
				order := &entity.Order{}
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
			},
			input: &store.DraftOrderInput{
				OrderID:         "order-id",
				ShippingMessage: "購入ありがとうございます。",
			},
			expect: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to draft",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.db.Order.EXPECT().Draft(ctx, "order-id", params).Return(assert.AnError)
			},
			input: &store.DraftOrderInput{
				OrderID:         "order-id",
				ShippingMessage: "購入ありがとうございます。",
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DraftOrder(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestCompleteOrder(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 10, 10, 18, 30, 0, 0)
	order := &entity.Order{
		ID:            "order-id",
		UserID:        "user-id",
		PromotionID:   "",
		CoordinatorID: "coordinator-id",
		CreatedAt:     now,
		UpdatedAt:     now,
		OrderPayment: entity.OrderPayment{
			OrderID:           "order-id",
			AddressRevisionID: 1,
			TransactionID:     "transaction-id",
			PaymentID:         "payment-id",
			Status:            entity.PaymentStatusCaptured,
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
				Status:            entity.FulfillmentStatusFulfilled,
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
	params := &database.CompleteOrderParams{
		ShippingMessage: "購入ありがとうございます。",
		CompletedAt:     now,
	}
	messengerIn := &messenger.NotifyOrderShippedInput{
		OrderID: "order-id",
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.CompleteOrderInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.db.Order.EXPECT().Complete(ctx, "order-id", params).Return(nil)
				mocks.messenger.EXPECT().NotifyOrderShipped(gomock.Any(), messengerIn).Return(assert.AnError)
			},
			input: &store.CompleteOrderInput{
				OrderID:         "order-id",
				ShippingMessage: "購入ありがとうございます。",
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.CompleteOrderInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get order",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(nil, assert.AnError)
			},
			input: &store.CompleteOrderInput{
				OrderID:         "order-id",
				ShippingMessage: "購入ありがとうございます。",
			},
			expect: exception.ErrInternal,
		},
		{
			name: "failed to completable",
			setup: func(ctx context.Context, mocks *mocks) {
				order := &entity.Order{}
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
			},
			input: &store.CompleteOrderInput{
				OrderID:         "order-id",
				ShippingMessage: "購入ありがとうございます。",
			},
			expect: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to complete",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.db.Order.EXPECT().Complete(ctx, "order-id", params).Return(assert.AnError)
			},
			input: &store.CompleteOrderInput{
				OrderID:         "order-id",
				ShippingMessage: "購入ありがとうございます。",
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.CompleteOrder(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}, withNow(now)))
	}
}

func TestCancelOrder(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 10, 10, 18, 30, 0, 0)
	order := &entity.Order{
		ID:            "order-id",
		UserID:        "user-id",
		PromotionID:   "",
		CoordinatorID: "coordinator-id",
		CreatedAt:     now,
		UpdatedAt:     now,
		OrderPayment: entity.OrderPayment{
			OrderID:           "order-id",
			AddressRevisionID: 1,
			TransactionID:     "transaction-id",
			PaymentID:         "payment-id",
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
	payment := &komoju.PaymentResponse{}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.CancelOrderInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.komojuPayment.EXPECT().Cancel(ctx, "payment-id").Return(payment, nil)
			},
			input: &store.CancelOrderInput{
				OrderID: "order-id",
			},
			expect: nil,
		},
		{
			name:   "failed to capture",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.CancelOrderInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get order",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(nil, assert.AnError)
			},
			input: &store.CancelOrderInput{
				OrderID: "order-id",
			},
			expect: exception.ErrInternal,
		},
		{
			name: "failed to capture",
			setup: func(ctx context.Context, mocks *mocks) {
				order := &entity.Order{
					OrderPayment: entity.OrderPayment{Status: entity.PaymentStatusCaptured},
				}
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
			},
			input: &store.CancelOrderInput{
				OrderID: "order-id",
			},
			expect: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to capture",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.komojuPayment.EXPECT().Cancel(ctx, "payment-id").Return(nil, assert.AnError)
			},
			input: &store.CancelOrderInput{
				OrderID: "order-id",
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.CancelOrder(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestRefundOrder(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 10, 10, 18, 30, 0, 0)
	order := &entity.Order{
		ID:            "order-id",
		UserID:        "user-id",
		PromotionID:   "",
		CoordinatorID: "coordinator-id",
		CreatedAt:     now,
		UpdatedAt:     now,
		OrderPayment: entity.OrderPayment{
			OrderID:           "order-id",
			AddressRevisionID: 1,
			TransactionID:     "transaction-id",
			PaymentID:         "payment-id",
			Status:            entity.PaymentStatusCaptured,
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
				Status:            entity.FulfillmentStatusFulfilled,
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
	params := &komoju.RefundParams{
		PaymentID:   "payment-id",
		Amount:      1600,
		Description: "在庫が不足していたため。",
	}
	payment := &komoju.PaymentResponse{}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.RefundOrderInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.komojuPayment.EXPECT().Refund(ctx, params).Return(payment, nil)
			},
			input: &store.RefundOrderInput{
				OrderID:     "order-id",
				Description: "在庫が不足していたため。",
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.RefundOrderInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get order",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(nil, assert.AnError)
			},
			input: &store.RefundOrderInput{
				OrderID:     "order-id",
				Description: "在庫が不足していたため。",
			},
			expect: exception.ErrInternal,
		},
		{
			name: "failed to refundable",
			setup: func(ctx context.Context, mocks *mocks) {
				order := &entity.Order{}
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
			},
			input: &store.RefundOrderInput{
				OrderID:     "order-id",
				Description: "在庫が不足していたため。",
			},
			expect: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to refund",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.komojuPayment.EXPECT().Refund(ctx, params).Return(nil, assert.AnError)
			},
			input: &store.RefundOrderInput{
				OrderID:     "order-id",
				Description: "在庫が不足していたため。",
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.RefundOrder(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestUpdateOrderFulfillment(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 10, 10, 18, 30, 0, 0)
	order := &entity.Order{
		ID:            "order-id",
		UserID:        "user-id",
		PromotionID:   "",
		CoordinatorID: "coordinator-id",
		CreatedAt:     now,
		UpdatedAt:     now,
		OrderPayment: entity.OrderPayment{
			OrderID:           "order-id",
			AddressRevisionID: 1,
			TransactionID:     "transaction-id",
			PaymentID:         "payment-id",
			Status:            entity.PaymentStatusCaptured,
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
				Status:            entity.FulfillmentStatusFulfilled,
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
	params := &database.UpdateOrderFulfillmentParams{
		Status:          entity.FulfillmentStatusFulfilled,
		ShippingCarrier: entity.ShippingCarrierYamato,
		TrackingNumber:  "tracking-number",
		ShippedAt:       now,
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.UpdateOrderFulfillmentInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.db.Order.EXPECT().UpdateFulfillment(ctx, "order-id", "fulfillment-id", params).Return(nil)
			},
			input: &store.UpdateOrderFulfillmentInput{
				OrderID:         "order-id",
				FulfillmentID:   "fulfillment-id",
				ShippingCarrier: entity.ShippingCarrierYamato,
				TrackingNumber:  "tracking-number",
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.UpdateOrderFulfillmentInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get order",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(nil, assert.AnError)
			},
			input: &store.UpdateOrderFulfillmentInput{
				OrderID:         "order-id",
				FulfillmentID:   "fulfillment-id",
				ShippingCarrier: entity.ShippingCarrierYamato,
				TrackingNumber:  "tracking-number",
			},
			expect: exception.ErrInternal,
		},
		{
			name: "failed to completed",
			setup: func(ctx context.Context, mocks *mocks) {
				order := &entity.Order{Status: entity.OrderStatusCompleted}
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
			},
			input: &store.UpdateOrderFulfillmentInput{
				OrderID:         "order-id",
				FulfillmentID:   "fulfillment-id",
				ShippingCarrier: entity.ShippingCarrierYamato,
				TrackingNumber:  "tracking-number",
			},
			expect: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to update fulfillment",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.db.Order.EXPECT().UpdateFulfillment(ctx, "order-id", "fulfillment-id", params).Return(assert.AnError)
			},
			input: &store.UpdateOrderFulfillmentInput{
				OrderID:         "order-id",
				FulfillmentID:   "fulfillment-id",
				ShippingCarrier: entity.ShippingCarrierYamato,
				TrackingNumber:  "tracking-number",
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateOrderFulfillment(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}, withNow(now)))
	}
}

func TestAggregateOrders(t *testing.T) {
	t.Parallel()

	params := &database.AggregateOrdersParams{
		CoordinatorID: "coordinator-id",
		UserIDs:       []string{"user-id"},
	}
	orders := entity.AggregatedOrders{
		{
			UserID:     "user-id",
			OrderCount: 2,
			Subtotal:   6000,
			Discount:   1000,
			Total:      5000,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.AggregateOrdersInput
		expect    entity.AggregatedOrders
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Aggregate(ctx, params).Return(orders, nil)
			},
			input: &store.AggregateOrdersInput{
				CoordinatorID: "coordinator-id",
				UserIDs:       []string{"user-id"},
			},
			expect:    orders,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.AggregateOrdersInput{
				UserIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to aggregate",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Aggregate(ctx, params).Return(nil, assert.AnError)
			},
			input: &store.AggregateOrdersInput{
				CoordinatorID: "coordinator-id",
				UserIDs:       []string{"user-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.AggregateOrders(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestAggregateOrdersByPromotion(t *testing.T) {
	t.Parallel()

	params := &database.AggregateOrdersByPromotionParams{
		CoordinatorID: "coordinator-id",
		PromotionIDs:  []string{"promotion-id"},
	}
	orders := entity.AggregatedOrderPromotions{
		{
			PromotionID:   "promotion-id",
			OrderCount:    2,
			DiscountTotal: 1000,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.AggregateOrdersByPromotionInput
		expect    entity.AggregatedOrderPromotions
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().AggregateByPromotion(ctx, params).Return(orders, nil)
			},
			input: &store.AggregateOrdersByPromotionInput{
				CoordinatorID: "coordinator-id",
				PromotionIDs:  []string{"promotion-id"},
			},
			expect:    orders,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.AggregateOrdersByPromotionInput{
				PromotionIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to aggregate",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().AggregateByPromotion(ctx, params).Return(nil, assert.AnError)
			},
			input: &store.AggregateOrdersByPromotionInput{
				CoordinatorID: "coordinator-id",
				PromotionIDs:  []string{"promotion-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.AggregateOrdersByPromotion(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}
