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

func TestUpsertExperienceReviewReaction(t *testing.T) {
	t.Parallel()

	now := time.Now()
	review := &entity.ExperienceReview{
		ID:           "review-id",
		ExperienceID: "experience-id",
		UserID:       "user-id",
		Rate:         5,
		Title:        "最高の体験",
		Comment:      "最高の体験でした。",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	reaction := &entity.ExperienceReviewReaction{
		ReviewID:     review.ID,
		UserID:       "user-id",
		ReactionType: entity.ExperienceReviewReactionTypeLike,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpsertExperienceReviewReactionInput
		expect    *entity.ExperienceReviewReaction
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReview.EXPECT().Get(ctx, review.ID).Return(review, nil)
				mocks.db.ExperienceReviewReaction.EXPECT().Upsert(ctx, reaction).Return(nil)
			},
			input: &store.UpsertExperienceReviewReactionInput{
				ReviewID:     review.ID,
				UserID:       "user-id",
				ReactionType: entity.ExperienceReviewReactionTypeLike,
			},
			expect:    reaction,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpsertExperienceReviewReactionInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get experience review",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReview.EXPECT().Get(ctx, review.ID).Return(nil, assert.AnError)
			},
			input: &store.UpsertExperienceReviewReactionInput{
				ReviewID:     review.ID,
				UserID:       "user-id",
				ReactionType: entity.ExperienceReviewReactionTypeLike,
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upsert experience review reaction",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReview.EXPECT().Get(ctx, review.ID).Return(review, nil)
				mocks.db.ExperienceReviewReaction.EXPECT().
					Upsert(ctx, reaction).
					Return(assert.AnError)
			},
			input: &store.UpsertExperienceReviewReactionInput{
				ReviewID:     review.ID,
				UserID:       "user-id",
				ReactionType: entity.ExperienceReviewReactionTypeLike,
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				actual, err := service.UpsertExperienceReviewReaction(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Equal(t, tt.expect, actual)
			}),
		)
	}
}

func TestDeleteExperienceReviewReaction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteExperienceReviewReactionInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReviewReaction.EXPECT().
					Delete(ctx, "preview-id", "user-id").
					Return(nil)
			},
			input: &store.DeleteExperienceReviewReactionInput{
				ReviewID: "preview-id",
				UserID:   "user-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteExperienceReviewReactionInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete experience review reaction",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReviewReaction.EXPECT().
					Delete(ctx, "preview-id", "user-id").
					Return(assert.AnError)
			},
			input: &store.DeleteExperienceReviewReactionInput{
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
				err := service.DeleteExperienceReviewReaction(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
			}),
		)
	}
}

func TestGetUserExperienceReviewReactions(t *testing.T) {
	t.Parallel()

	now := time.Now()
	reactions := entity.ExperienceReviewReactions{
		{
			ReviewID:     "review-id",
			UserID:       "user-id",
			ReactionType: entity.ExperienceReviewReactionTypeLike,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetUserExperienceReviewReactionsInput
		expect    entity.ExperienceReviewReactions
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReviewReaction.EXPECT().
					GetUserReactions(ctx, "experience-id", "user-id").
					Return(reactions, nil)
			},
			input: &store.GetUserExperienceReviewReactionsInput{
				ExperienceID: "experience-id",
				UserID:       "user-id",
			},
			expect:    reactions,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetUserExperienceReviewReactionsInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get user experience review reactions",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReviewReaction.EXPECT().
					GetUserReactions(ctx, "experience-id", "user-id").
					Return(nil, assert.AnError)
			},
			input: &store.GetUserExperienceReviewReactionsInput{
				ExperienceID: "experience-id",
				UserID:       "user-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				actual, err := service.GetUserExperienceReviewReactions(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Equal(t, tt.expect, actual)
			}),
		)
	}
}
