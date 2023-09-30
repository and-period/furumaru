package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
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
			ID:                "order-id",
			UserID:            "user-id",
			ScheduleID:        "schedule-id",
			PromotionID:       "",
			CoordinatorID:     "coordinator-id",
			PaymentStatus:     entity.PaymentStatusPending,
			FulfillmentStatus: entity.FulfillmentStatusUnfulfilled,
			RefundReason:      "",
			CreatedAt:         now,
			UpdatedAt:         now,
			OrderItems: []*entity.OrderItem{
				{
					OrderID:   "order-id",
					ProductID: "product-id01",
					Price:     100,
					Quantity:  1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					OrderID:   "order-id",
					ProductID: "product-id02",
					Price:     500,
					Quantity:  2,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			Payment: entity.Payment{
				OrderID:       "order-id",
				AddressID:     "address-id",
				TransactionID: "transaction-id",
				MethodType:    entity.PaymentMethodTypeCreditCard,
				Subtotal:      1100,
				Discount:      0,
				ShippingFee:   500,
				Tax:           160,
				Total:         1760,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			Fulfillment: entity.Fulfillment{
				OrderID:         "order-id",
				AddressID:       "address-id",
				TrackingNumber:  "",
				ShippingCarrier: entity.ShippingCarrierUnknown,
				ShippingMethod:  entity.DeliveryTypeNormal,
				BoxSize:         entity.ShippingSize60,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			Activities: []*entity.Activity{},
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
			expectErr:   exception.ErrUnknown,
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
			expectErr:   exception.ErrUnknown,
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
		ID:                "order-id",
		UserID:            "user-id",
		PaymentStatus:     entity.PaymentStatusPending,
		FulfillmentStatus: entity.FulfillmentStatusUnfulfilled,
		RefundReason:      "",
		CreatedAt:         now,
		UpdatedAt:         now,
		OrderItems: []*entity.OrderItem{
			{
				OrderID:   "order-id",
				ProductID: "product-id01",
				Price:     100,
				Quantity:  1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				OrderID:   "order-id",
				ProductID: "product-id02",
				Price:     500,
				Quantity:  2,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		Payment: entity.Payment{
			OrderID:       "order-id",
			AddressID:     "address-id",
			TransactionID: "transaction-id",
			MethodType:    entity.PaymentMethodTypeCreditCard,
			Subtotal:      1100,
			Discount:      0,
			ShippingFee:   500,
			Tax:           160,
			Total:         1760,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		Fulfillment: entity.Fulfillment{
			OrderID:         "order-id",
			AddressID:       "address-id",
			TrackingNumber:  "",
			ShippingCarrier: entity.ShippingCarrierUnknown,
			ShippingMethod:  entity.DeliveryTypeNormal,
			BoxSize:         entity.ShippingSize60,
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		Activities: []*entity.Activity{},
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
			expectErr: exception.ErrUnknown,
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

func TestAggregateOrders(t *testing.T) {
	t.Parallel()

	orders := entity.AggregatedOrders{
		{
			UserID:     "user-id",
			OrderCount: 2,
			Subtotal:   6000,
			Discount:   1000,
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
				mocks.db.Order.EXPECT().Aggregate(ctx, []string{"user-id"}).Return(orders, nil)
			},
			input: &store.AggregateOrdersInput{
				UserIDs: []string{"user-id"},
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
				mocks.db.Order.EXPECT().Aggregate(ctx, []string{"user-id"}).Return(nil, assert.AnError)
			},
			input: &store.AggregateOrdersInput{
				UserIDs: []string{"user-id"},
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
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
