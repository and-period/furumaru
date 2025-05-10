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

func TestListProductTags(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	params := &database.ListProductTagsParams{
		Name:   "野菜",
		Limit:  30,
		Offset: 0,
		Orders: []*database.ListProductTagsOrder{
			{Key: database.ListProductTagsOrderByName, OrderByASC: true},
		},
	}
	productTags := entity.ProductTags{
		{
			ID:        "product-tag-id",
			Name:      "野菜",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListProductTagsInput
		expect      entity.ProductTags
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductTag.EXPECT().List(gomock.Any(), params).Return(productTags, nil)
				mocks.db.ProductTag.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListProductTagsInput{
				Name:   "野菜",
				Limit:  30,
				Offset: 0,
				Orders: []*store.ListProductTagsOrder{
					{Key: store.ListProductTagsOrderByName, OrderByASC: true},
				},
			},
			expect:      productTags,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &store.ListProductTagsInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductTag.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.ProductTag.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListProductTagsInput{
				Name:   "野菜",
				Limit:  30,
				Offset: 0,
				Orders: []*store.ListProductTagsOrder{
					{Key: store.ListProductTagsOrderByName, OrderByASC: true},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductTag.EXPECT().List(gomock.Any(), params).Return(productTags, nil)
				mocks.db.ProductTag.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &store.ListProductTagsInput{
				Name:   "野菜",
				Limit:  30,
				Offset: 0,
				Orders: []*store.ListProductTagsOrder{
					{Key: store.ListProductTagsOrderByName, OrderByASC: true},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListProductTags(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestMultiGetProductTags(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	productTags := entity.ProductTags{
		{
			ID:        "product-tag-id",
			Name:      "野菜",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetProductTagsInput
		expect    entity.ProductTags
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductTag.EXPECT().MultiGet(ctx, []string{"product-tag-id"}).Return(productTags, nil)
			},
			input: &store.MultiGetProductTagsInput{
				ProductTagIDs: []string{"product-tag-id"},
			},
			expect:    productTags,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetProductTagsInput{
				ProductTagIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductTag.EXPECT().MultiGet(ctx, []string{"product-tag-id"}).Return(nil, assert.AnError)
			},
			input: &store.MultiGetProductTagsInput{
				ProductTagIDs: []string{"product-tag-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetProductTags(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetProductTag(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	productTag := &entity.ProductTag{
		ID:        "product-tag-id",
		Name:      "野菜",
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetProductTagInput
		expect    *entity.ProductTag
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductTag.EXPECT().Get(ctx, "product-tag-id").Return(productTag, nil)
			},
			input: &store.GetProductTagInput{
				ProductTagID: "product-tag-id",
			},
			expect:    productTag,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetProductTagInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductTag.EXPECT().Get(ctx, "product-tag-id").Return(nil, assert.AnError)
			},
			input: &store.GetProductTagInput{
				ProductTagID: "product-tag-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetProductTag(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateProductTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateProductTagInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductTag.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, productTag *entity.ProductTag) error {
						expect := &entity.ProductTag{
							ID:   productTag.ID, // ignore
							Name: "野菜",
						}
						assert.Equal(t, expect, productTag)
						return nil
					})
			},
			input: &store.CreateProductTagInput{
				Name: "野菜",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateProductTagInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductTag.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateProductTagInput{
				Name: "野菜",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateProductTag(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateProductTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateProductTagInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductTag.EXPECT().Update(ctx, "product-tag-id", "野菜").Return(nil)
			},
			input: &store.UpdateProductTagInput{
				ProductTagID: "product-tag-id",
				Name:         "野菜",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateProductTagInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductTag.EXPECT().Update(ctx, "product-tag-id", "野菜").Return(assert.AnError)
			},
			input: &store.UpdateProductTagInput{
				ProductTagID: "product-tag-id",
				Name:         "野菜",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateProductTag(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteProductTag(t *testing.T) {
	t.Parallel()

	params := &database.ListProductsParams{
		ProductTagID: "product-tag-id",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteProductTagInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().Count(ctx, params).Return(int64(0), nil)
				mocks.db.ProductTag.EXPECT().Delete(ctx, "product-tag-id").Return(nil)
			},
			input: &store.DeleteProductTagInput{
				ProductTagID: "product-tag-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteProductTagInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to count",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().Count(ctx, params).Return(int64(0), assert.AnError)
			},
			input: &store.DeleteProductTagInput{
				ProductTagID: "product-tag-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "associated with product",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().Count(ctx, params).Return(int64(1), nil)
			},
			input: &store.DeleteProductTagInput{
				ProductTagID: "product-tag-id",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to delete",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().Count(ctx, params).Return(int64(0), nil)
				mocks.db.ProductTag.EXPECT().Delete(ctx, "product-tag-id").Return(assert.AnError)
			},
			input: &store.DeleteProductTagInput{
				ProductTagID: "product-tag-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteProductTag(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
