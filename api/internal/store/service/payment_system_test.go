package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/stretchr/testify/assert"
)

func TestMultiGetPaymentSystems(t *testing.T) {
	t.Parallel()
	now := time.Now()
	systems := entity.PaymentSystems{
		{
			MethodType: entity.PaymentMethodTypeCreditCard,
			Status:     entity.PaymentSystemStatusInUse,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetPaymentSystemsInput
		expect    entity.PaymentSystems
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				types := []entity.PaymentMethodType{entity.PaymentMethodTypeCreditCard}
				mocks.db.PaymentSystem.EXPECT().MultiGet(ctx, types).Return(systems, nil)
			},
			input: &store.MultiGetPaymentSystemsInput{
				MethodTypes: []entity.PaymentMethodType{entity.PaymentMethodTypeCreditCard},
			},
			expect:    systems,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetPaymentSystemsInput{
				MethodTypes: []entity.PaymentMethodType{entity.PaymentMethodTypeUnknown},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multi get payment systems",
			setup: func(ctx context.Context, mocks *mocks) {
				types := []entity.PaymentMethodType{entity.PaymentMethodTypeCreditCard}
				mocks.db.PaymentSystem.EXPECT().MultiGet(ctx, types).Return(nil, assert.AnError)
			},
			input: &store.MultiGetPaymentSystemsInput{
				MethodTypes: []entity.PaymentMethodType{entity.PaymentMethodTypeCreditCard},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetPaymentSystems(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestGetPaymentSystem(t *testing.T) {
	t.Parallel()
	now := time.Now()
	system := &entity.PaymentSystem{
		MethodType: entity.PaymentMethodTypeCreditCard,
		Status:     entity.PaymentSystemStatusInUse,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetPaymentSystemInput
		expect    *entity.PaymentSystem
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.PaymentSystem.EXPECT().Get(ctx, entity.PaymentMethodTypeCreditCard).Return(system, nil)
			},
			input: &store.GetPaymentSystemInput{
				MethodType: entity.PaymentMethodTypeCreditCard,
			},
			expect:    system,
			expectErr: nil,
		},
		{
			name: "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input:     &store.GetPaymentSystemInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get payment system",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.PaymentSystem.EXPECT().Get(ctx, entity.PaymentMethodTypeCreditCard).Return(nil, assert.AnError)
			},
			input: &store.GetPaymentSystemInput{
				MethodType: entity.PaymentMethodTypeCreditCard,
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetPaymentSystem(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUpdatePaymentSystem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdatePaymentStatusInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.PaymentSystem.EXPECT().
					Update(ctx, entity.PaymentMethodTypeCreditCard, entity.PaymentSystemStatusOutage).
					Return(nil)
			},
			input: &store.UpdatePaymentStatusInput{
				MethodType: entity.PaymentMethodTypeCreditCard,
				Status:     entity.PaymentSystemStatusOutage,
			},
			expectErr: nil,
		},
		{
			name: "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input:     &store.UpdatePaymentStatusInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update payment system",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.PaymentSystem.EXPECT().
					Update(ctx, entity.PaymentMethodTypeCreditCard, entity.PaymentSystemStatusOutage).
					Return(assert.AnError)
			},
			input: &store.UpdatePaymentStatusInput{
				MethodType: entity.PaymentMethodTypeCreditCard,
				Status:     entity.PaymentSystemStatusOutage,
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdatePaymentSystem(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
