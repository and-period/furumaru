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

func TestListExperienceReviews(t *testing.T) {
	t.Parallel()

	now := time.Now()

	params := &database.ListExperienceReviewsParams{
		ExperienceID: "experience-id",
		UserID:       "user-id",
		Rates:        []int64{4, 5},
		Limit:        10,
	}
	reviews := entity.ExperienceReviews{
		{
			ID:           "review-id",
			ExperienceID: "experience-id",
			UserID:       "user-id",
			Rate:         5,
			Title:        "最高の体験",
			Comment:      "最高の体験でした。",
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListExperienceReviewsInput
		expect      entity.ExperienceReviews
		expectToken string
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReview.EXPECT().
					List(ctx, params).
					Return(reviews, "next-token", nil)
			},
			input: &store.ListExperienceReviewsInput{
				ExperienceID: "experience-id",
				UserID:       "user-id",
				Rates:        []int64{4, 5},
				Limit:        10,
			},
			expect:      reviews,
			expectToken: "next-token",
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &store.ListExperienceReviewsInput{},
			expect:      nil,
			expectToken: "",
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list experience reviews",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReview.EXPECT().List(ctx, params).Return(nil, "", assert.AnError)
			},
			input: &store.ListExperienceReviewsInput{
				ExperienceID: "experience-id",
				UserID:       "user-id",
				Rates:        []int64{4, 5},
				Limit:        10,
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
				actual, token, err := service.ListExperienceReviews(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Equal(t, tt.expectToken, token)
				assert.Equal(t, tt.expect, actual)
			}),
		)
	}
}

func TestGetExperienceReview(t *testing.T) {
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

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetExperienceReviewInput
		expect    *entity.ExperienceReview
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReview.EXPECT().Get(ctx, "review-id").Return(review, nil)
			},
			input: &store.GetExperienceReviewInput{
				ReviewID: "review-id",
			},
			expect:    review,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetExperienceReviewInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get experience review",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReview.EXPECT().Get(ctx, "review-id").Return(nil, assert.AnError)
			},
			input: &store.GetExperienceReviewInput{
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
				actual, err := service.GetExperienceReview(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Equal(t, tt.expect, actual)
			}),
		)
	}
}

func TestCreateExperienceReview(t *testing.T) {
	t.Parallel()

	now := time.Now()
	experience := &entity.Experience{
		ID:            "experience-id",
		CoordinatorID: "coordinator-id",
		ProducerID:    "producer-id",
		TypeID:        "experience-type-id",
		Title:         "じゃがいも収穫",
		Description:   "じゃがいもを収穫する体験です。",
		Public:        true,
		SoldOut:       false,
		Status:        entity.ExperienceStatusAccepting,
		ThumbnailURL:  "http://example.com/thumbnail.png",
		Media: []*entity.ExperienceMedia{
			{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
			{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
		},
		RecommendedPoints: []string{
			"じゃがいもを収穫する楽しさを体験できます。",
			"新鮮なじゃがいもを持ち帰ることができます。",
		},
		PromotionVideoURL:  "http://example.com/promotion.mp4",
		Duration:           60,
		Direction:          "彦根駅から徒歩10分",
		BusinessOpenTime:   "1000",
		BusinessCloseTime:  "1800",
		HostPostalCode:     "5220061",
		HostPrefecture:     "滋賀県",
		HostPrefectureCode: 25,
		HostCity:           "彦根市",
		HostAddressLine1:   "金亀町１−１",
		HostAddressLine2:   "",
		HostLongitude:      136.251739,
		HostLatitude:       35.276833,
		StartAt:            now.AddDate(0, 0, -1),
		EndAt:              now.AddDate(0, 0, 1),
		ExperienceRevision: entity.ExperienceRevision{
			ID:                    1,
			ExperienceID:          "experience-id",
			PriceAdult:            1000,
			PriceJuniorHighSchool: 800,
			PriceElementarySchool: 600,
			PricePreschool:        400,
			PriceSenior:           700,
			CreatedAt:             now,
			UpdatedAt:             now,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateExperienceReviewInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Experience.EXPECT().Get(ctx, "experience-id").Return(experience, nil)
				mocks.db.ExperienceReview.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, review *entity.ExperienceReview) error {
						expect := &entity.ExperienceReview{
							ID:           review.ID, // ignore
							ExperienceID: "experience-id",
							UserID:       "user-id",
							Rate:         5,
							Title:        "最高の体験",
							Comment:      "最高の体験でした。",
						}
						assert.Equal(t, expect, review)
						return nil
					})
			},
			input: &store.CreateExperienceReviewInput{
				ExperienceID: "experience-id",
				UserID:       "user-id",
				Rate:         5,
				Title:        "最高の体験",
				Comment:      "最高の体験でした。",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateExperienceReviewInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get experience",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Experience.EXPECT().Get(ctx, "experience-id").Return(nil, assert.AnError)
			},
			input: &store.CreateExperienceReviewInput{
				ExperienceID: "experience-id",
				UserID:       "user-id",
				Rate:         5,
				Title:        "最高の体験",
				Comment:      "最高の体験でした。",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create experience review",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Experience.EXPECT().Get(ctx, "experience-id").Return(experience, nil)
				mocks.db.ExperienceReview.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateExperienceReviewInput{
				ExperienceID: "experience-id",
				UserID:       "user-id",
				Rate:         5,
				Title:        "最高の体験",
				Comment:      "最高の体験でした。",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				_, err := service.CreateExperienceReview(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
			}),
		)
	}
}

func TestUpdateExperienceReview(t *testing.T) {
	t.Parallel()

	params := &database.UpdateExperienceReviewParams{
		Rate:    4,
		Title:   "良い体験",
		Comment: "良い体験でした。",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateExperienceReviewInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReview.EXPECT().Update(ctx, "review-id", params).Return(nil)
			},
			input: &store.UpdateExperienceReviewInput{
				ReviewID: "review-id",
				Rate:     4,
				Title:    "良い体験",
				Comment:  "良い体験でした。",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateExperienceReviewInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update experience review",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReview.EXPECT().
					Update(ctx, "review-id", params).
					Return(assert.AnError)
			},
			input: &store.UpdateExperienceReviewInput{
				ReviewID: "review-id",
				Rate:     4,
				Title:    "良い体験",
				Comment:  "良い体験でした。",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				err := service.UpdateExperienceReview(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
			}),
		)
	}
}

func TestDeleteExperienceReview(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteExperienceReviewInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReview.EXPECT().Delete(ctx, "review-id").Return(nil)
			},
			input: &store.DeleteExperienceReviewInput{
				ReviewID: "review-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteExperienceReviewInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete experience review",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReview.EXPECT().Delete(ctx, "review-id").Return(assert.AnError)
			},
			input: &store.DeleteExperienceReviewInput{
				ReviewID: "review-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				err := service.DeleteExperienceReview(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
			}),
		)
	}
}

func TestAggregateExperienceReviews(t *testing.T) {
	t.Parallel()

	params := &database.AggregateExperienceReviewsParams{
		ExperienceIDs: []string{"experience-id"},
	}
	reviews := entity.AggregatedExperienceReviews{
		{
			ExperienceID: "experience-id",
			Count:        4,
			Average:      2.5,
			Rate1:        2,
			Rate2:        0,
			Rate3:        1,
			Rate4:        0,
			Rate5:        1,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.AggregateExperienceReviewsInput
		expect    entity.AggregatedExperienceReviews
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReview.EXPECT().Aggregate(ctx, params).Return(reviews, nil)
			},
			input: &store.AggregateExperienceReviewsInput{
				ExperienceIDs: []string{"experience-id"},
			},
			expect:    reviews,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.AggregateExperienceReviewsInput{
				ExperienceIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to aggregate experience reviews",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceReview.EXPECT().
					Aggregate(ctx, params).
					Return(nil, assert.AnError)
			},
			input: &store.AggregateExperienceReviewsInput{
				ExperienceIDs: []string{"experience-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				actual, err := service.AggregateExperienceReviews(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Equal(t, tt.expect, actual)
			}),
		)
	}
}
