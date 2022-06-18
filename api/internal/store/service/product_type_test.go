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

func TestListProductTypes(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	params := &database.ListProductTypesParams{
		Name:       "じゃがいも",
		CategoryID: "category-id",
		Limit:      30,
		Offset:     0,
	}
	productTypes := entity.ProductTypes{
		{
			ID:         "product-type-id",
			Name:       "じゃがいも",
			CategoryID: "category-id",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.ListProductTypesInput
		expect    entity.ProductTypes
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().List(ctx, params).Return(productTypes, nil)
			},
			input: &store.ListProductTypesInput{
				Name:       "じゃがいも",
				CategoryID: "category-id",
				Limit:      30,
				Offset:     0,
			},
			expect:    productTypes,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.ListProductTypesInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().List(ctx, params).Return(nil, errmock)
			},
			input: &store.ListProductTypesInput{
				Name:       "じゃがいも",
				CategoryID: "category-id",
				Limit:      30,
				Offset:     0,
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *storeService) {
			actual, err := service.ListProductTypes(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestCreateProductType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateProductTypeInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, productType *entity.ProductType) error {
						expect := &entity.ProductType{
							ID:         productType.ID, // ignore
							Name:       "じゃがいも",
							CategoryID: "category-id",
						}
						assert.Equal(t, expect, productType)
						return nil
					})
			},
			input: &store.CreateProductTypeInput{
				Name:       "じゃがいも",
				CategoryID: "category-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateProductTypeInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().Create(ctx, gomock.Any()).Return(errmock)
			},
			input: &store.CreateProductTypeInput{
				Name:       "じゃがいも",
				CategoryID: "category-id",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *storeService) {
			_, err := service.CreateProductType(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateProductType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateProductTypeInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().Update(ctx, "product-type-id", "さつまいも").Return(nil)
			},
			input: &store.UpdateProductTypeInput{
				ProductTypeID: "product-type-id",
				Name:          "さつまいも",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateProductTypeInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().Update(ctx, "product-type-id", "さつまいも").Return(errmock)
			},
			input: &store.UpdateProductTypeInput{
				ProductTypeID: "product-type-id",
				Name:          "さつまいも",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *storeService) {
			err := service.UpdateProductType(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteProductType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteProductTypeInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().Delete(ctx, "product-type-id").Return(nil)
			},
			input: &store.DeleteProductTypeInput{
				ProductTypeID: "product-type-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteProductTypeInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().Delete(ctx, "product-type-id").Return(errmock)
			},
			input: &store.DeleteProductTypeInput{
				ProductTypeID: "product-type-id",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *storeService) {
			err := service.DeleteProductType(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
