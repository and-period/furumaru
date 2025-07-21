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

func TestUpsertProductReviewReaction(t *testing.T) {
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
	reaction := &entity.ProductReviewReaction{
		ReviewID:     review.ID,
		UserID:       "user-id",
		ReactionType: entity.ProductReviewReactionTypeLike,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpsertProductReviewReactionInput
		expect    *entity.ProductReviewReaction
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReview.EXPECT().Get(ctx, review.ID).Return(review, nil)
				mocks.db.ProductReviewReaction.EXPECT().Upsert(ctx, reaction).Return(nil)
			},
			input: &store.UpsertProductReviewReactionInput{
				ReviewID:     review.ID,
				UserID:       "user-id",
				ReactionType: entity.ProductReviewReactionTypeLike,
			},
			expect:    reaction,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpsertProductReviewReactionInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get product review",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReview.EXPECT().Get(ctx, review.ID).Return(nil, assert.AnError)
			},
			input: &store.UpsertProductReviewReactionInput{
				ReviewID:     review.ID,
				UserID:       "user-id",
				ReactionType: entity.ProductReviewReactionTypeLike,
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upsert product review reaction",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReview.EXPECT().Get(ctx, review.ID).Return(review, nil)
				mocks.db.ProductReviewReaction.EXPECT().Upsert(ctx, reaction).Return(assert.AnError)
			},
			input: &store.UpsertProductReviewReactionInput{
				ReviewID:     review.ID,
				UserID:       "user-id",
				ReactionType: entity.ProductReviewReactionTypeLike,
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				actual, err := service.UpsertProductReviewReaction(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Equal(t, tt.expect, actual)
			}),
		)
	}
}

func TestDeleteProductReviewReaction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteProductReviewReactionInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReviewReaction.EXPECT().
					Delete(ctx, "preview-id", "user-id").
					Return(nil)
			},
			input: &store.DeleteProductReviewReactionInput{
				ReviewID: "preview-id",
				UserID:   "user-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteProductReviewReactionInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete product review reaction",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReviewReaction.EXPECT().
					Delete(ctx, "preview-id", "user-id").
					Return(assert.AnError)
			},
			input: &store.DeleteProductReviewReactionInput{
				ReviewID: "preview-id",
				UserID:   "user-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				err := service.DeleteProductReviewReaction(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
			}),
		)
	}
}

func TestGetUserProductReviewReactions(t *testing.T) {
	t.Parallel()

	now := time.Now()
	reactions := entity.ProductReviewReactions{
		{
			ReviewID:     "review-id",
			UserID:       "user-id",
			ReactionType: entity.ProductReviewReactionTypeLike,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetUserProductReviewReactionsInput
		expect    entity.ProductReviewReactions
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReviewReaction.EXPECT().
					GetUserReactions(ctx, "product-id", "user-id").
					Return(reactions, nil)
			},
			input: &store.GetUserProductReviewReactionsInput{
				ProductID: "product-id",
				UserID:    "user-id",
			},
			expect:    reactions,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetUserProductReviewReactionsInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get user product review reactions",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductReviewReaction.EXPECT().
					GetUserReactions(ctx, "product-id", "user-id").
					Return(nil, assert.AnError)
			},
			input: &store.GetUserProductReviewReactionsInput{
				ProductID: "product-id",
				UserID:    "user-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				actual, err := service.GetUserProductReviewReactions(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Equal(t, tt.expect, actual)
			}),
		)
	}
}
