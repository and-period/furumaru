package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/stretchr/testify/assert"
)

func TestListProducts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.ListProductsInput
		expect    entity.Products
		expectErr error
	}{
		{
			name:  "not implemented",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.ListProductsInput{
				Name:          "みかん",
				CoordinatorID: "",
				ProducerID:    "",
				Limit:         30,
				Offset:        0,
			},
			expect:    nil,
			expectErr: exception.ErrNotImplemented,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.ListProductsInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.ListProducts(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetProductInput
		expect    *entity.Product
		expectErr error
	}{
		{
			name:  "not implemented",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.GetProductInput{
				ProductID: "product-id",
			},
			expect:    nil,
			expectErr: exception.ErrNotImplemented,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetProductInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetProduct(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateProduct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateProductInput
		expect    *entity.Product
		expectErr error
	}{
		{
			name:  "not implemented",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.CreateProductInput{
				CoordinatorID: "coordinator-id",
				ProducerID:    "producer-id",
				CategoryID:    "category-id",
				TypeID:        "product-type-id",
				Name:          "新鮮なじゃがいも", Description: "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*store.CreateProductMedia{
					{
						URL:         "https://and-period.jp/thumbnail.png",
						IsThumbnail: true,
					},
				},
				Price:            400,
				DeliveryType:     entity.DeliveryTypeNormal,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
			},
			expect:    nil,
			expectErr: exception.ErrNotImplemented,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateProductInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.CreateProduct(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUpdateProduct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateProductInput
		expect    entity.Products
		expectErr error
	}{
		{
			name:  "not implemented",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.UpdateProductInput{
				ProductID:       "product-id",
				CoordinatorID:   "coordinator-id",
				CategoryID:      "category-id",
				TypeID:          "product-type-id",
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*store.CreateProductMedia{
					{
						URL:         "https://and-period.jp/thumbnail.png",
						IsThumbnail: true,
					},
				},
				Price:            400,
				DeliveryType:     entity.DeliveryTypeNormal,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
			},
			expect:    nil,
			expectErr: exception.ErrNotImplemented,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateProductInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateProduct(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteProduct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteProductInput
		expect    entity.Products
		expectErr error
	}{
		{
			name:  "not implemented",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.DeleteProductInput{
				ProductID: "product-id",
			},
			expect:    nil,
			expectErr: exception.ErrNotImplemented,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteProductInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteProduct(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
