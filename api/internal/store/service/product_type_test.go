package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
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
		Orders: []*database.ListProductTypesOrder{
			{Key: database.ListProductTypesOrderByName, OrderByASC: true},
		},
	}
	productTypes := entity.ProductTypes{
		{
			ID:         "product-type-id",
			Name:       "じゃがいも",
			IconURL:    "https://and-period.jp/icon.png",
			CategoryID: "category-id",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListProductTypesInput
		expect      entity.ProductTypes
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().List(gomock.Any(), params).Return(productTypes, nil)
				mocks.db.ProductType.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListProductTypesInput{
				Name:       "じゃがいも",
				CategoryID: "category-id",
				Limit:      30,
				Offset:     0,
				Orders: []*store.ListProductTypesOrder{
					{Key: store.ListProductTypesOrderByName, OrderByASC: true},
				},
			},
			expect:      productTypes,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &store.ListProductTypesInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.ProductType.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListProductTypesInput{
				Name:       "じゃがいも",
				CategoryID: "category-id",
				Limit:      30,
				Offset:     0,
				Orders: []*store.ListProductTypesOrder{
					{Key: store.ListProductTypesOrderByName, OrderByASC: true},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().List(gomock.Any(), params).Return(productTypes, nil)
				mocks.db.ProductType.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &store.ListProductTypesInput{
				Name:       "じゃがいも",
				CategoryID: "category-id",
				Limit:      30,
				Offset:     0,
				Orders: []*store.ListProductTypesOrder{
					{Key: store.ListProductTypesOrderByName, OrderByASC: true},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListProductTypes(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestMultiGetProductTypes(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	productTypes := entity.ProductTypes{
		{
			ID:         "product-type-id",
			Name:       "じゃがいも",
			IconURL:    "https://and-period.jp/icon.png",
			CategoryID: "category-id",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetProductTypesInput
		expect    entity.ProductTypes
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().MultiGet(ctx, []string{"product-type-id"}).Return(productTypes, nil)
			},
			input: &store.MultiGetProductTypesInput{
				ProductTypeIDs: []string{"product-type-id"},
			},
			expect:    productTypes,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetProductTypesInput{
				ProductTypeIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().MultiGet(ctx, []string{"product-type-id"}).Return(nil, assert.AnError)
			},
			input: &store.MultiGetProductTypesInput{
				ProductTypeIDs: []string{"product-type-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetProductTypes(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetProductType(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	productType := &entity.ProductType{
		ID:         "product-type-id",
		Name:       "じゃがいも",
		IconURL:    "https://and-period.jp/icon.png",
		CategoryID: "category-id",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetProductTypeInput
		expect    *entity.ProductType
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().Get(ctx, "product-type-id").Return(productType, nil)
			},
			input: &store.GetProductTypeInput{
				ProductTypeID: "product-type-id",
			},
			expect:    productType,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetProductTypeInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get product type",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().Get(ctx, "product-type-id").Return(nil, assert.AnError)
			},
			input: &store.GetProductTypeInput{
				ProductTypeID: "product-type-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetProductType(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
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
							IconURL:    "https://and-period.jp/icon.png",
							CategoryID: "category-id",
						}
						assert.Equal(t, expect, productType)
						return nil
					})
			},
			input: &store.CreateProductTypeInput{
				Name:       "じゃがいも",
				IconURL:    "https://and-period.jp/icon.png",
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
				mocks.db.ProductType.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateProductTypeInput{
				Name:       "じゃがいも",
				IconURL:    "https://and-period.jp/icon.png",
				CategoryID: "category-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
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
				mocks.db.ProductType.EXPECT().Update(ctx, "product-type-id", "さつまいも", "https://tmp.and-period.jp/icon.png").Return(nil)
			},
			input: &store.UpdateProductTypeInput{
				ProductTypeID: "product-type-id",
				Name:          "さつまいも",
				IconURL:       "https://tmp.and-period.jp/icon.png",
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
				mocks.db.ProductType.EXPECT().Update(ctx, "product-type-id", "さつまいも", "https://and-period.jp/icon.png").Return(assert.AnError)
			},
			input: &store.UpdateProductTypeInput{
				ProductTypeID: "product-type-id",
				Name:          "さつまいも",
				IconURL:       "https://and-period.jp/icon.png",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateProductType(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteProductType(t *testing.T) {
	t.Parallel()
	params := &database.ListProductsParams{
		ProductTypeIDs: []string{"product-type-id"},
	}
	in := &user.RemoveCoordinatorProductTypeInput{
		ProductTypeID: "product-type-id",
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteProductTypeInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().Count(ctx, params).Return(int64(0), nil)
				mocks.db.Shop.EXPECT().RemoveProductType(ctx, "product-type-id").Return(nil)
				mocks.db.ProductType.EXPECT().Delete(ctx, "product-type-id").Return(nil)
				mocks.user.EXPECT().RemoveCoordinatorProductType(gomock.Any(), in).Return(assert.AnError)
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
			name: "failed to count",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().Count(ctx, params).Return(int64(0), assert.AnError)
			},
			input: &store.DeleteProductTypeInput{
				ProductTypeID: "product-type-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "associated with product",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().Count(ctx, params).Return(int64(1), nil)
			},
			input: &store.DeleteProductTypeInput{
				ProductTypeID: "product-type-id",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to remove from shop",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().Count(ctx, params).Return(int64(0), nil)
				mocks.db.Shop.EXPECT().RemoveProductType(ctx, "product-type-id").Return(assert.AnError)
			},
			input: &store.DeleteProductTypeInput{
				ProductTypeID: "product-type-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to delete",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().Count(ctx, params).Return(int64(0), nil)
				mocks.db.Shop.EXPECT().RemoveProductType(ctx, "product-type-id").Return(nil)
				mocks.db.ProductType.EXPECT().Delete(ctx, "product-type-id").Return(assert.AnError)
			},
			input: &store.DeleteProductTypeInput{
				ProductTypeID: "product-type-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteProductType(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
