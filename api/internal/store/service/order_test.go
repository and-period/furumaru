package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/komoju"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListOrders(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 10, 10, 18, 30, 0, 0)
	params := &database.ListOrdersParams{
		ShopID: "shop-id",
		Limit:  30,
		Offset: 0,
	}
	orders := entity.Orders{
		{
			ID:            "order-id",
			UserID:        "user-id",
			PromotionID:   "",
			ShopID:        "shop-id",
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
				ShopID: "shop-id",
				Limit:  30,
				Offset: 0,
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
				ShopID: "shop-id",
				Limit:  30,
				Offset: 0,
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
				ShopID: "shop-id",
				Limit:  30,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListOrders(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestListOrderUserIDs(t *testing.T) {
	t.Parallel()
	params := &database.ListOrdersParams{
		ShopID: "shop-id",
		Limit:  30,
		Offset: 0,
	}
	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListOrderUserIDsInput
		expect      []string
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().ListUserIDs(ctx, params).Return([]string{"user-id"}, int64(1), nil)
			},
			input: &store.ListOrderUserIDsInput{
				ShopID: "shop-id",
				Limit:  30,
				Offset: 0,
			},
			expect:      []string{"user-id"},
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &store.ListOrderUserIDsInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list user ids",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().ListUserIDs(ctx, params).Return(nil, int64(0), assert.AnError)
			},
			input: &store.ListOrderUserIDsInput{
				ShopID: "shop-id",
				Limit:  30,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListOrderUserIDs(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
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
		ShopID:        "shop-id",
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
		ShopID:        "shop-id",
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
		ShopID:        "shop-id",
		CoordinatorID: "coordinator-id",
		CreatedAt:     now,
		UpdatedAt:     now,
		Status:        entity.OrderStatusWaiting,
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
		ShopID:        "shop-id",
		CoordinatorID: "coordinator-id",
		CreatedAt:     now,
		UpdatedAt:     now,
		Status:        entity.OrderStatusPreparing,
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
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DraftOrder(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestCompleteProductOrder(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 10, 10, 18, 30, 0, 0)
	order := &entity.Order{
		ID:            "order-id",
		UserID:        "user-id",
		PromotionID:   "",
		ShopID:        "shop-id",
		CoordinatorID: "coordinator-id",
		CreatedAt:     now,
		UpdatedAt:     now,
		Status:        entity.OrderStatusShipped,
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
		input  *store.CompleteProductOrderInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.db.Order.EXPECT().Complete(ctx, "order-id", params).Return(nil)
				mocks.messenger.EXPECT().NotifyOrderShipped(gomock.Any(), messengerIn).Return(assert.AnError)
			},
			input: &store.CompleteProductOrderInput{
				OrderID:         "order-id",
				ShippingMessage: "購入ありがとうございます。",
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.CompleteProductOrderInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get order",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(nil, assert.AnError)
			},
			input: &store.CompleteProductOrderInput{
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
			input: &store.CompleteProductOrderInput{
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
			input: &store.CompleteProductOrderInput{
				OrderID:         "order-id",
				ShippingMessage: "購入ありがとうございます。",
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.CompleteProductOrder(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}, withNow(now)))
	}
}

func TestCompleteExperienceOrder(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 10, 10, 18, 30, 0, 0)
	order := &entity.Order{
		ID:            "order-id",
		UserID:        "user-id",
		PromotionID:   "",
		ShopID:        "shop-id",
		CoordinatorID: "coordinator-id",
		CreatedAt:     now,
		UpdatedAt:     now,
		Status:        entity.OrderStatusPreparing,
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
		OrderExperience: entity.OrderExperience{
			OrderID:               "order-id",
			ExperienceRevisionID:  1,
			AdultCount:            2,
			JuniorHighSchoolCount: 1,
			ElementarySchoolCount: 0,
			PreschoolCount:        0,
			SeniorCount:           0,
			Remarks:               entity.OrderExperienceRemarks{},
			CreatedAt:             now,
			UpdatedAt:             now,
		},
	}
	params := &database.CompleteOrderParams{
		CompletedAt: now,
	}
	messengerIn := &messenger.NotifyReviewRequestInput{
		OrderID: "order-id",
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.CompleteExperienceOrderInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.db.Order.EXPECT().Complete(ctx, "order-id", params).Return(nil)
				mocks.messenger.EXPECT().NotifyReviewRequest(gomock.Any(), messengerIn).Return(assert.AnError)
			},
			input: &store.CompleteExperienceOrderInput{
				OrderID: "order-id",
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.CompleteExperienceOrderInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get order",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(nil, assert.AnError)
			},
			input: &store.CompleteExperienceOrderInput{
				OrderID: "order-id",
			},
			expect: exception.ErrInternal,
		},
		{
			name: "failed to completable",
			setup: func(ctx context.Context, mocks *mocks) {
				order := &entity.Order{}
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
			},
			input: &store.CompleteExperienceOrderInput{
				OrderID: "order-id",
			},
			expect: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to complete",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Get(ctx, "order-id").Return(order, nil)
				mocks.db.Order.EXPECT().Complete(ctx, "order-id", params).Return(assert.AnError)
			},
			input: &store.CompleteExperienceOrderInput{
				OrderID: "order-id",
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.CompleteExperienceOrder(ctx, tt.input)
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
		ShopID:        "shop-id",
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
		ShopID:        "shop-id",
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
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.RefundOrder(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestUpdateOrderFulfillment(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 10, 10, 18, 30, 0, 0)
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
			name: "failed to update fulfillment",
			setup: func(ctx context.Context, mocks *mocks) {
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
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateOrderFulfillment(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}, withNow(now)))
	}
}

func TestAggregateOrders(t *testing.T) {
	t.Parallel()

	now := time.Now()
	params := &database.AggregateOrdersParams{
		ShopID:       "shop-id",
		CreatedAtGte: now.AddDate(0, 0, -7),
		CreatedAtLt:  now,
	}
	order := &entity.AggregatedOrder{
		OrderCount:    2,
		UserCount:     1,
		SalesTotal:    6000,
		DiscountTotal: 1000,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.AggregateOrdersInput
		expect    *entity.AggregatedOrder
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Aggregate(ctx, params).Return(order, nil)
			},
			input: &store.AggregateOrdersInput{
				ShopID:       "shop-id",
				CreatedAtGte: now.AddDate(0, 0, -7),
				CreatedAtLt:  now,
			},
			expect:    order,
			expectErr: nil,
		},
		{
			name: "failed to aggregate",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().Aggregate(ctx, params).Return(nil, assert.AnError)
			},
			input: &store.AggregateOrdersInput{
				ShopID:       "shop-id",
				CreatedAtGte: now.AddDate(0, 0, -7),
				CreatedAtLt:  now,
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.AggregateOrders(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestAggregateOrdersByUser(t *testing.T) {
	t.Parallel()

	params := &database.AggregateOrdersByUserParams{
		ShopID:  "shop-id",
		UserIDs: []string{"user-id"},
	}
	orders := entity.AggregatedUserOrders{
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
		input     *store.AggregateOrdersByUserInput
		expect    entity.AggregatedUserOrders
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().AggregateByUser(ctx, params).Return(orders, nil)
			},
			input: &store.AggregateOrdersByUserInput{
				ShopID:  "shop-id",
				UserIDs: []string{"user-id"},
			},
			expect:    orders,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.AggregateOrdersByUserInput{
				UserIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to aggregate",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().AggregateByUser(ctx, params).Return(nil, assert.AnError)
			},
			input: &store.AggregateOrdersByUserInput{
				ShopID:  "shop-id",
				UserIDs: []string{"user-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.AggregateOrdersByUser(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestAggregateOrdersByPaymentMethodType(t *testing.T) {
	t.Parallel()

	params := &database.AggregateOrdersByPaymentMethodTypeParams{
		ShopID:             "shop-id",
		PaymentMethodTypes: entity.AllPaymentMethodTypes,
	}
	orders := entity.AggregatedOrderPayments{
		{
			PaymentMethodType: entity.PaymentMethodTypeCreditCard,
			OrderCount:        2,
			UserCount:         1,
			SalesTotal:        6000,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.AggregateOrdersByPaymentMethodTypeInput
		expect    entity.AggregatedOrderPayments
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().AggregateByPaymentMethodType(ctx, params).Return(orders, nil)
			},
			input: &store.AggregateOrdersByPaymentMethodTypeInput{
				ShopID: "shop-id",
			},
			expect:    orders,
			expectErr: nil,
		},
		{
			name: "failed to aggregate",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().AggregateByPaymentMethodType(ctx, params).Return(nil, assert.AnError)
			},
			input: &store.AggregateOrdersByPaymentMethodTypeInput{
				ShopID: "shop-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.AggregateOrdersByPaymentMethodType(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestAggregateOrdersByPromotion(t *testing.T) {
	t.Parallel()

	params := &database.AggregateOrdersByPromotionParams{
		ShopID:       "shop-id",
		PromotionIDs: []string{"promotion-id"},
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
				ShopID:       "shop-id",
				PromotionIDs: []string{"promotion-id"},
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
				ShopID:       "shop-id",
				PromotionIDs: []string{"promotion-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.AggregateOrdersByPromotion(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestAggregateOrdersByPeriod(t *testing.T) {
	t.Parallel()

	now := time.Now()
	params := &database.AggregateOrdersByPeriodParams{
		ShopID:       "shop-id",
		PeriodType:   entity.AggregateOrderPeriodTypeDay,
		CreatedAtGte: now.AddDate(0, 0, -7),
		CreatedAtLt:  now,
	}
	orders := entity.AggregatedPeriodOrders{
		{
			Period:        now.Truncate(24 * time.Hour),
			OrderCount:    2,
			UserCount:     1,
			SalesTotal:    6000,
			DiscountTotal: 1000,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.AggregateOrdersByPeriodInput
		expect    entity.AggregatedPeriodOrders
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().AggregateByPeriod(ctx, params).Return(orders, nil)
			},
			input: &store.AggregateOrdersByPeriodInput{
				ShopID:       "shop-id",
				PeriodType:   entity.AggregateOrderPeriodTypeDay,
				CreatedAtGte: now.AddDate(0, 0, -7),
				CreatedAtLt:  now,
			},
			expect:    orders,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.AggregateOrdersByPeriodInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to aggregate",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().AggregateByPeriod(ctx, params).Return(nil, assert.AnError)
			},
			input: &store.AggregateOrdersByPeriodInput{
				ShopID:       "shop-id",
				PeriodType:   entity.AggregateOrderPeriodTypeDay,
				CreatedAtGte: now.AddDate(0, 0, -7),
				CreatedAtLt:  now,
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.AggregateOrdersByPeriod(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestExportOrders(t *testing.T) {
	t.Parallel()
	now := jst.Date(2024, 1, 23, 18, 30, 0, 0)
	ordersParams := &database.ListOrdersParams{
		ShopID:   "shop-id",
		Statuses: []entity.OrderStatus{entity.OrderStatusPreparing},
	}
	orders := entity.Orders{
		{
			ID:            "order-id",
			UserID:        "user-id",
			ShopID:        "shop-id",
			CoordinatorID: "coordinator-id",
			PromotionID:   "promotion-id",
			ManagementID:  1,
			Status:        entity.OrderStatusPreparing,
			OrderPayment: entity.OrderPayment{
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
				OrderedAt:         now,
				PaidAt:            now,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			OrderFulfillments: entity.OrderFulfillments{
				{
					ID:                "fulfillment-id",
					OrderID:           "order-id",
					AddressRevisionID: 1,
					TrackingNumber:    "",
					Status:            entity.FulfillmentStatusFulfilled,
					ShippingCarrier:   entity.ShippingCarrierUnknown,
					ShippingType:      entity.ShippingTypeNormal,
					BoxNumber:         1,
					BoxSize:           entity.ShippingSize60,
					CreatedAt:         now,
					UpdatedAt:         now,
				},
			},
			OrderItems: entity.OrderItems{
				{
					FulfillmentID:     "fulfillment-id",
					OrderID:           "order-id",
					ProductRevisionID: 1,
					Quantity:          1,
					CreatedAt:         now,
					UpdatedAt:         now,
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	addressesIn := &user.MultiGetAddressesByRevisionInput{
		AddressRevisionIDs: []int64{1},
	}
	addresses := uentity.Addresses{
		{
			ID:        "address-id",
			UserID:    "user-id",
			IsDefault: true,
			AddressRevision: uentity.AddressRevision{
				ID:             1,
				Lastname:       "&.",
				Firstname:      "購入者",
				LastnameKana:   "あんどどっと",
				FirstnameKana:  "こうにゅうしゃ",
				PostalCode:     "1000014",
				Prefecture:     "東京都",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "090-1234-1234",
			},
		},
	}
	products := entity.Products{
		{
			ID:              "product-id",
			TypeID:          "type-id",
			TagIDs:          []string{"tag-id"},
			ShopID:          "shop-id",
			CoordinatorID:   "coordinator-id",
			ProducerID:      "producer-id",
			Name:            "新鮮なじゃがいも",
			Description:     "新鮮なじゃがいもをお届けします。",
			Scope:           entity.ProductScopePublic,
			Inventory:       100,
			Weight:          100,
			WeightUnit:      entity.WeightUnitGram,
			Item:            1,
			ItemUnit:        "袋",
			ItemDescription: "1袋あたり100gのじゃがいも",
			Media: entity.MultiProductMedia{
				{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
				{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
			},
			ExpirationDate:    7,
			StorageMethodType: entity.StorageMethodTypeNormal,
			DeliveryType:      entity.DeliveryTypeNormal,
			Box60Rate:         50,
			Box80Rate:         40,
			Box100Rate:        30,
			OriginPrefecture:  "滋賀県",
			OriginCity:        "彦根市",
			ProductRevision: entity.ProductRevision{
				ID:        1,
				ProductID: "product-id",
				Price:     400,
				Cost:      300,
				CreatedAt: now,
				UpdatedAt: now,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.ExportOrdersInput
		expect    string
		expectErr error
	}{
		{
			name: "success general with body",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().List(ctx, ordersParams).Return(orders, nil)
				mocks.db.Product.EXPECT().MultiGetByRevision(gomock.Any(), []int64{1}).Return(products, nil)
				mocks.user.EXPECT().MultiGetAddressesByRevision(gomock.Any(), addressesIn).Return(addresses, nil)
			},
			input: &store.ExportOrdersInput{
				ShopID:          "shop-id",
				ShippingCarrier: entity.ShippingCarrierUnknown,
				EncodingType:    codes.CharacterEncodingTypeUTF8,
			},
			expect: "注文管理番号,ユーザーID,コーディネータID,お届け希望日,お届け希望時間帯,お届け先名,お届け先名（かな）,お届け先電話番号,お届け先郵便番号,お届け先都道府県,お届け先市区町村,お届け先町名・番地,お届け先ビル名・号室など,ご依頼主名,ご依頼主名（かな）,ご依頼主電話番号,ご依頼主郵便番号,ご依頼主都道府県,ご依頼主市区町村,ご依頼主町名・番地,ご依頼主ビル名・号室など,商品コード1,商品名1,商品1数量,商品コード2,商品名2,商品2数量,商品コード3,商品名3,商品3数量,決済手段,商品金額,割引金額,配送手数料,合計金額,注文日時,配送方法,箱のサイズ\n" +
				"order-id,user-id,coordinator-id,,,&. 購入者,あんどどっと こうにゅうしゃ,090-1234-1234,1000014,東京都,千代田区,永田町1-7-1,,&. 購入者,あんどどっと こうにゅうしゃ,090-1234-1234,1000014,東京都,千代田区,永田町1-7-1,,product-id,新鮮なじゃがいも,1,,,0,,,0,クレジットカード決済,1980,0,550,2530,2024-01-23 18:30:00,通常配送,60\n",
			expectErr: nil,
		},
		{
			name: "success general without body",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().List(ctx, ordersParams).Return(entity.Orders{}, nil)
			},
			input: &store.ExportOrdersInput{
				ShopID:          "shop-id",
				ShippingCarrier: entity.ShippingCarrierUnknown,
				EncodingType:    codes.CharacterEncodingTypeUTF8,
			},
			expect:    "注文管理番号,ユーザーID,コーディネータID,お届け希望日,お届け希望時間帯,お届け先名,お届け先名（かな）,お届け先電話番号,お届け先郵便番号,お届け先都道府県,お届け先市区町村,お届け先町名・番地,お届け先ビル名・号室など,ご依頼主名,ご依頼主名（かな）,ご依頼主電話番号,ご依頼主郵便番号,ご依頼主都道府県,ご依頼主市区町村,ご依頼主町名・番地,ご依頼主ビル名・号室など,商品コード1,商品名1,商品1数量,商品コード2,商品名2,商品2数量,商品コード3,商品名3,商品3数量,決済手段,商品金額,割引金額,配送手数料,合計金額,注文日時,配送方法,箱のサイズ\n",
			expectErr: nil,
		},
		{
			name: "success yamato with body",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().List(ctx, ordersParams).Return(orders, nil)
				mocks.db.Product.EXPECT().MultiGetByRevision(gomock.Any(), []int64{1}).Return(products, nil)
				mocks.user.EXPECT().MultiGetAddressesByRevision(gomock.Any(), addressesIn).Return(addresses, nil)
			},
			input: &store.ExportOrdersInput{
				ShopID:          "shop-id",
				ShippingCarrier: entity.ShippingCarrierYamato,
				EncodingType:    codes.CharacterEncodingTypeUTF8,
			},
			expect: "お客様管理番号,送り状種類,クール区分,伝票番号,出荷予定日,お届け予定日,配送時間帯,お届け先コード,お届け先電話番号,お届け先電話番号枝番,お届け先郵便番号,お届け先住所,お届け先アパートマンション名,お届け先会社・部門１,お届け先会社・部門２,お届け先名,お届け先名(ｶﾅ),敬称,ご依頼主コード,ご依頼主電話番号,ご依頼主電話番号枝番,ご依頼主郵便番号,ご依頼主住所,ご依頼主アパートマンション名,ご依頼主名,ご依頼主名(ｶﾅ),品名コード１,品名１,品名コード２,品名２,荷扱い１,荷扱い２,記事,ｺﾚｸﾄ代金引換額（税込),内消費税額等,止置き,営業所コード,発行枚数,個数口枠の印字,請求先顧客コード,請求先分類コード,運賃管理番号,クロネコwebコレクトデータ登録,クロネコwebコレクト加盟店番号,クロネコwebコレクト申込受付番号１,クロネコwebコレクト申込受付番号２,クロネコwebコレクト申込受付番号３,お届け予定ｅメール利用区分,お届け予定ｅメールe-mailアドレス,入力機種,お届け予定ｅメールメッセージ,お届け完了ｅメール利用区分,お届け完了ｅメールe-mailアドレス,お届け完了ｅメールメッセージ,クロネコ収納代行利用区分,予備,収納代行請求金額(税込),収納代行内消費税額等,収納代行請求先郵便番号,収納代行請求先住所,収納代行請求先住所（アパートマンション名）,収納代行請求先会社・部門名１,収納代行請求先会社・部門名２,収納代行請求先名(漢字),収納代行請求先名(カナ),収納代行問合せ先郵便番号,収納代行問合せ先住所,収納代行問合せ先住所（アパートマンション名）,収納代行問合せ先電話番号,収納代行管理番号,収納代行品名,収納代行備考,複数口くくりキー,検索キータイトル1,検索キー1,検索キータイトル2,検索キー2,検索キータイトル3,検索キー3,検索キータイトル4,検索キー4,検索キータイトル5,検索キー5,予備,予備,投函予定メール利用区分,投函予定メールe-mailアドレス,投函予定メールメッセージ,投函完了メール（お届け先宛）利用区分,投函予定メール（お届け先宛）e-mailアドレス,投函予定メール（お届け先宛）メッセージ,投函完了メール（ご依頼主宛）利用区分,投函予定メール（ご依頼主宛）e-mailアドレス,投函予定メール（ご依頼主宛）メッセージ\n" +
				"order-id,0,0,,,,,1,090-1234-1234,,1000014,東京都 千代田区 永田町1-7-1,,,,&. 購入者,あんどどっと こうにゅうしゃ,様,1,090-1234-1234,,1000014,東京都 千代田区 永田町1-7-1,,&. 購入者,あんどどっと こうにゅうしゃ,product-id,新鮮なじゃがいも,,,,,,0,0,1,,1,,,,,0,,,,,0,,,,0,,,0,,0,0,,,,,,,,,,,,,,,,ふるマルユーザーID,user-id,ふるマル注文履歴ID,order-id,,,,,,,,,0,,,0,,,,,\n",
			expectErr: nil,
		},
		{
			name: "success yamato without body",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().List(ctx, ordersParams).Return(entity.Orders{}, nil)
			},
			input: &store.ExportOrdersInput{
				ShopID:          "shop-id",
				ShippingCarrier: entity.ShippingCarrierYamato,
				EncodingType:    codes.CharacterEncodingTypeUTF8,
			},
			expect:    "お客様管理番号,送り状種類,クール区分,伝票番号,出荷予定日,お届け予定日,配送時間帯,お届け先コード,お届け先電話番号,お届け先電話番号枝番,お届け先郵便番号,お届け先住所,お届け先アパートマンション名,お届け先会社・部門１,お届け先会社・部門２,お届け先名,お届け先名(ｶﾅ),敬称,ご依頼主コード,ご依頼主電話番号,ご依頼主電話番号枝番,ご依頼主郵便番号,ご依頼主住所,ご依頼主アパートマンション名,ご依頼主名,ご依頼主名(ｶﾅ),品名コード１,品名１,品名コード２,品名２,荷扱い１,荷扱い２,記事,ｺﾚｸﾄ代金引換額（税込),内消費税額等,止置き,営業所コード,発行枚数,個数口枠の印字,請求先顧客コード,請求先分類コード,運賃管理番号,クロネコwebコレクトデータ登録,クロネコwebコレクト加盟店番号,クロネコwebコレクト申込受付番号１,クロネコwebコレクト申込受付番号２,クロネコwebコレクト申込受付番号３,お届け予定ｅメール利用区分,お届け予定ｅメールe-mailアドレス,入力機種,お届け予定ｅメールメッセージ,お届け完了ｅメール利用区分,お届け完了ｅメールe-mailアドレス,お届け完了ｅメールメッセージ,クロネコ収納代行利用区分,予備,収納代行請求金額(税込),収納代行内消費税額等,収納代行請求先郵便番号,収納代行請求先住所,収納代行請求先住所（アパートマンション名）,収納代行請求先会社・部門名１,収納代行請求先会社・部門名２,収納代行請求先名(漢字),収納代行請求先名(カナ),収納代行問合せ先郵便番号,収納代行問合せ先住所,収納代行問合せ先住所（アパートマンション名）,収納代行問合せ先電話番号,収納代行管理番号,収納代行品名,収納代行備考,複数口くくりキー,検索キータイトル1,検索キー1,検索キータイトル2,検索キー2,検索キータイトル3,検索キー3,検索キータイトル4,検索キー4,検索キータイトル5,検索キー5,予備,予備,投函予定メール利用区分,投函予定メールe-mailアドレス,投函予定メールメッセージ,投函完了メール（お届け先宛）利用区分,投函予定メール（お届け先宛）e-mailアドレス,投函予定メール（お届け先宛）メッセージ,投函完了メール（ご依頼主宛）利用区分,投函予定メール（ご依頼主宛）e-mailアドレス,投函予定メール（ご依頼主宛）メッセージ\n",
			expectErr: nil,
		},
		{
			name: "success sagawa with body",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().List(ctx, ordersParams).Return(orders, nil)
				mocks.db.Product.EXPECT().MultiGetByRevision(gomock.Any(), []int64{1}).Return(products, nil)
				mocks.user.EXPECT().MultiGetAddressesByRevision(gomock.Any(), addressesIn).Return(addresses, nil)
			},
			input: &store.ExportOrdersInput{
				ShopID:          "shop-id",
				ShippingCarrier: entity.ShippingCarrierSagawa,
				EncodingType:    codes.CharacterEncodingTypeUTF8,
			},
			expect: "お届け先コード取得区分,お届け先コード,お届け先電話番号,お届け先郵便番号,お届け先住所１,お届け先住所２,お届け先住所３,お届け先名称１,お届け先名称２,お客様管理番号,お客様コード,部署ご担当者コード,取得区分,部署ご担当者コード,部署ご担当者名称,荷送人電話番号,ご依頼主コード取得区分,ご依頼主コード,ご依頼主電話番号,ご依頼主郵便番号,ご依頼主住所１,ご依頼主住所２,ご依頼主名称１,ご依頼主名称２,荷姿,品名１,品名２,品名３,品名４,品名５,荷札荷姿,荷札品名1,荷札品名2,荷札品名3,荷札品名4,荷札品名5,荷札品名6,荷札品名7,荷札品名8,荷札品名9,荷札品名10,荷札品名11,出荷個数,スピード指定,クール便指定,配達日,配達指定時間帯,配達指定時間（時分）,代引金額,消費税,決済種別,保険金額,指定シール1,指定シール2,指定シール3,営業所受取,ＳＲＣ区分,営業所受取営業所コード,元着区分,メールアドレス,ご不在時連絡先,出荷日,お問い合せ送り状No.,出荷場印字区分,集約解除指定,編集1,編集2,編集3,編集4,編集5,編集6,編集7,編集8,編集9,編集10\n" +
				",,090-1234-1234,1000014,東京都,千代田区 永田町1-7-1,,&.,購入者,order-id,user-id,,,,,,,090-1234-1234,1000014,東京都 千代田区 永田町1-7-1,,&.,購入者,,新鮮なじゃがいも,,,,,,,,,,,,,,,,,1,,001,,,,0,0,,0,,,,,,,,,,,,,,,,,,,,,,,\n",
			expectErr: nil,
		},
		{
			name: "success sagawa without body",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().List(ctx, ordersParams).Return(entity.Orders{}, nil)
			},
			input: &store.ExportOrdersInput{
				ShopID:          "shop-id",
				ShippingCarrier: entity.ShippingCarrierSagawa,
				EncodingType:    codes.CharacterEncodingTypeUTF8,
			},
			expect:    "お届け先コード取得区分,お届け先コード,お届け先電話番号,お届け先郵便番号,お届け先住所１,お届け先住所２,お届け先住所３,お届け先名称１,お届け先名称２,お客様管理番号,お客様コード,部署ご担当者コード,取得区分,部署ご担当者コード,部署ご担当者名称,荷送人電話番号,ご依頼主コード取得区分,ご依頼主コード,ご依頼主電話番号,ご依頼主郵便番号,ご依頼主住所１,ご依頼主住所２,ご依頼主名称１,ご依頼主名称２,荷姿,品名１,品名２,品名３,品名４,品名５,荷札荷姿,荷札品名1,荷札品名2,荷札品名3,荷札品名4,荷札品名5,荷札品名6,荷札品名7,荷札品名8,荷札品名9,荷札品名10,荷札品名11,出荷個数,スピード指定,クール便指定,配達日,配達指定時間帯,配達指定時間（時分）,代引金額,消費税,決済種別,保険金額,指定シール1,指定シール2,指定シール3,営業所受取,ＳＲＣ区分,営業所受取営業所コード,元着区分,メールアドレス,ご不在時連絡先,出荷日,お問い合せ送り状No.,出荷場印字区分,集約解除指定,編集1,編集2,編集3,編集4,編集5,編集6,編集7,編集8,編集9,編集10\n",
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.ExportOrdersInput{
				ShippingCarrier: -1,
				EncodingType:    -1,
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list orders",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().List(ctx, ordersParams).Return(nil, assert.AnError)
			},
			input: &store.ExportOrdersInput{
				ShopID:          "shop-id",
				ShippingCarrier: entity.ShippingCarrierUnknown,
				EncodingType:    codes.CharacterEncodingTypeUTF8,
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to multi get products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().List(ctx, ordersParams).Return(orders, nil)
				mocks.db.Product.EXPECT().MultiGetByRevision(gomock.Any(), []int64{1}).Return(nil, assert.AnError)
				mocks.user.EXPECT().MultiGetAddressesByRevision(gomock.Any(), addressesIn).Return(addresses, nil)
			},
			input: &store.ExportOrdersInput{
				ShopID:          "shop-id",
				ShippingCarrier: entity.ShippingCarrierUnknown,
				EncodingType:    codes.CharacterEncodingTypeUTF8,
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to multi get addresses",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().List(ctx, ordersParams).Return(orders, nil)
				mocks.db.Product.EXPECT().MultiGetByRevision(gomock.Any(), []int64{1}).Return(products, nil)
				mocks.user.EXPECT().MultiGetAddressesByRevision(gomock.Any(), addressesIn).Return(nil, assert.AnError)
			},
			input: &store.ExportOrdersInput{
				ShopID:          "shop-id",
				ShippingCarrier: entity.ShippingCarrierUnknown,
				EncodingType:    codes.CharacterEncodingTypeUTF8,
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.ExportOrders(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, string(actual))
		}))
	}
}
