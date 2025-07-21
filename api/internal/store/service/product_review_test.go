package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListProductReviews(t *testing.T) {
	t.Parallel()

	now := time.Now()

	params := &database.ListProductReviewsParams{
		ProductID: "product-id",
		UserID:    "user-id",
		Rates:     []int64{4, 5},
		Limit:     10,
	}
	reviews := entity.ProductReviews{
		{
			ID:        "review-id",
			ProductID: "product-id",
			UserID:    "user-id",
			Rate:      5,
			Title:     "最高の商品",
			Comment:   "最高の商品でした。",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListProductReviewsInput
		expect      entity.ProductReviews
		expectToken string
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReview.EXPECT().List(ctx, params).Return(reviews, "next-token", nil)
			},
			input: &store.ListProductReviewsInput{
				ProductID: "product-id",
				UserID:    "user-id",
				Rates:     []int64{4, 5},
				Limit:     10,
			},
			expect:      reviews,
			expectToken: "next-token",
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &store.ListProductReviewsInput{},
			expect:      nil,
			expectToken: "",
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list product reviews",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReview.EXPECT().List(ctx, params).Return(nil, "", assert.AnError)
			},
			input: &store.ListProductReviewsInput{
				ProductID: "product-id",
				UserID:    "user-id",
				Rates:     []int64{4, 5},
				Limit:     10,
			},
			expect:      nil,
			expectToken: "",
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				actual, token, err := service.ListProductReviews(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Equal(t, tt.expectToken, token)
				assert.Equal(t, tt.expect, actual)
			}),
		)
	}
}

func TestGetProductReview(t *testing.T) {
	t.Parallel()

	now := time.Now()

	review := &entity.ProductReview{
		ID:        "review-id",
		ProductID: "product-id",
		UserID:    "user-id",
		Rate:      5,
		Title:     "最高の商品",
		Comment:   "最高の商品でした。",
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetProductReviewInput
		expect    *entity.ProductReview
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReview.EXPECT().Get(ctx, "review-id").Return(review, nil)
			},
			input: &store.GetProductReviewInput{
				ReviewID: "review-id",
			},
			expect:    review,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetProductReviewInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get product review",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReview.EXPECT().Get(ctx, "review-id").Return(nil, assert.AnError)
			},
			input: &store.GetProductReviewInput{
				ReviewID: "review-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				actual, err := service.GetProductReview(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Equal(t, tt.expect, actual)
			}),
		)
	}
}

func TestCreateProductReview(t *testing.T) {
	t.Parallel()

	now := time.Now()
	product := &entity.Product{
		ID:              "product-id",
		TypeID:          "type-id",
		TagIDs:          []string{"tag-id"},
		CoordinatorID:   "coordinator-id",
		ProducerID:      "producer-id",
		Name:            "新鮮なじゃがいも",
		Description:     "新鮮なじゃがいもをお届けします。",
		Public:          true,
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
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateProductReviewInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().Get(ctx, "product-id").Return(product, nil)
				mocks.db.ProductReview.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, review *entity.ProductReview) error {
						expect := &entity.ProductReview{
							ID:        review.ID, // ignore
							ProductID: "product-id",
							UserID:    "user-id",
							Rate:      5,
							Title:     "最高の商品",
							Comment:   "最高の商品でした。",
						}
						assert.Equal(t, expect, review)
						return nil
					})
			},
			input: &store.CreateProductReviewInput{
				ProductID: "product-id",
				UserID:    "user-id",
				Rate:      5,
				Title:     "最高の商品",
				Comment:   "最高の商品でした。",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateProductReviewInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get product",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().Get(ctx, "product-id").Return(nil, assert.AnError)
			},
			input: &store.CreateProductReviewInput{
				ProductID: "product-id",
				UserID:    "user-id",
				Rate:      5,
				Title:     "最高の商品",
				Comment:   "最高の商品でした。",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create product review",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().Get(ctx, "product-id").Return(product, nil)
				mocks.db.ProductReview.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateProductReviewInput{
				ProductID: "product-id",
				UserID:    "user-id",
				Rate:      5,
				Title:     "最高の商品",
				Comment:   "最高の商品でした。",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				_, err := service.CreateProductReview(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
			}),
		)
	}
}

func TestUpdateProductReview(t *testing.T) {
	t.Parallel()

	params := &database.UpdateProductReviewParams{
		Rate:    4,
		Title:   "良い商品",
		Comment: "良い商品でした。",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateProductReviewInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReview.EXPECT().Update(ctx, "review-id", params).Return(nil)
			},
			input: &store.UpdateProductReviewInput{
				ReviewID: "review-id",
				Rate:     4,
				Title:    "良い商品",
				Comment:  "良い商品でした。",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateProductReviewInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update product review",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReview.EXPECT().
					Update(ctx, "review-id", params).
					Return(assert.AnError)
			},
			input: &store.UpdateProductReviewInput{
				ReviewID: "review-id",
				Rate:     4,
				Title:    "良い商品",
				Comment:  "良い商品でした。",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				err := service.UpdateProductReview(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
			}),
		)
	}
}

func TestDeleteProductReview(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteProductReviewInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReview.EXPECT().Delete(ctx, "review-id").Return(nil)
			},
			input: &store.DeleteProductReviewInput{
				ReviewID: "review-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteProductReviewInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete product review",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReview.EXPECT().Delete(ctx, "review-id").Return(assert.AnError)
			},
			input: &store.DeleteProductReviewInput{
				ReviewID: "review-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				err := service.DeleteProductReview(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
			}),
		)
	}
}

func TestAggregateProductReviews(t *testing.T) {
	t.Parallel()

	params := &database.AggregateProductReviewsParams{
		ProductIDs: []string{"product-id"},
	}
	reviews := entity.AggregatedProductReviews{
		{
			ProductID: "product-id",
			Count:     4,
			Average:   2.5,
			Rate1:     2,
			Rate2:     0,
			Rate3:     1,
			Rate4:     0,
			Rate5:     1,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.AggregateProductReviewsInput
		expect    entity.AggregatedProductReviews
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReview.EXPECT().Aggregate(ctx, params).Return(reviews, nil)
			},
			input: &store.AggregateProductReviewsInput{
				ProductIDs: []string{"product-id"},
			},
			expect:    reviews,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.AggregateProductReviewsInput{
				ProductIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to aggregate product reviews",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReview.EXPECT().Aggregate(ctx, params).Return(nil, assert.AnError)
			},
			input: &store.AggregateProductReviewsInput{
				ProductIDs: []string{"product-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				actual, err := service.AggregateProductReviews(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Equal(t, tt.expect, actual)
			}),
		)
	}
}
