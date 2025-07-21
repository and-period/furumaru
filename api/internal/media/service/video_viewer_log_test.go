package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateVideoViewerLog(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 10, 20, 18, 30, 0, 0)
	video := &entity.Video{
		ID:            "video-id",
		CoordinatorID: "coordinator-id",
		ProductIDs:    []string{"product-id"},
		ExperienceIDs: []string{"experience-id"},
		Title:         "じゃがいも収穫",
		Description:   "じゃがいも収穫の説明",
		Status:        entity.VideoStatusPublished,
		ThumbnailURL:  "https://example.com/thumbnail.jpg",
		VideoURL:      "https://example.com/video.mp4",
		Public:        true,
		Limited:       false,
		VideoProducts: []*entity.VideoProduct{{
			VideoID:   "video-id",
			ProductID: "product-id",
			Priority:  1,
			CreatedAt: now,
			UpdatedAt: now,
		}},
		VideoExperiences: []*entity.VideoExperience{{
			VideoID:      "video-id",
			ExperienceID: "experience-id",
			Priority:     1,
			CreatedAt:    now,
			UpdatedAt:    now,
		}},
		PublishedAt: now.AddDate(0, 0, -1),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	log := &entity.VideoViewerLog{
		VideoID:   "video-id",
		SessionID: "session-id",
		UserID:    "user-id",
		UserAgent: "user-agent",
		ClientIP:  "127.0.0.1",
	}

	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.CreateVideoViewerLogInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().Get(ctx, "video-id").Return(video, nil)
				mocks.db.VideoViewerLog.EXPECT().Create(ctx, log).Return(nil)
			},
			input: &media.CreateVideoViewerLogInput{
				VideoID:   "video-id",
				SessionID: "session-id",
				UserID:    "user-id",
				UserAgent: "user-agent",
				ClientIP:  "127.0.0.1",
			},
			expect: nil,
		},
		{
			name: "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input:  &media.CreateVideoViewerLogInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get video",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().Get(ctx, "video-id").Return(nil, assert.AnError)
			},
			input: &media.CreateVideoViewerLogInput{
				VideoID:   "video-id",
				SessionID: "session-id",
				UserID:    "user-id",
				UserAgent: "user-agent",
				ClientIP:  "127.0.0.1",
			},
			expect: exception.ErrInternal,
		},
		{
			name: "video is not published",
			setup: func(ctx context.Context, mocks *mocks) {
				video := &entity.Video{Status: entity.VideoStatusPrivate}
				mocks.db.Video.EXPECT().Get(ctx, "video-id").Return(video, nil)
			},
			input: &media.CreateVideoViewerLogInput{
				VideoID:   "video-id",
				SessionID: "session-id",
				UserID:    "user-id",
				UserAgent: "user-agent",
				ClientIP:  "127.0.0.1",
			},
			expect: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to create broadacast viewer log",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().Get(ctx, "video-id").Return(video, nil)
				mocks.db.VideoViewerLog.EXPECT().Create(ctx, log).Return(assert.AnError)
			},
			input: &media.CreateVideoViewerLogInput{
				VideoID:   "video-id",
				SessionID: "session-id",
				UserID:    "user-id",
				UserAgent: "user-agent",
				ClientIP:  "127.0.0.1",
			},
			expect: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				err := service.CreateVideoViewerLog(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expect)
			}),
		)
	}
}

func TestAggregateVideoViewerLogs(t *testing.T) {
	t.Parallel()

	now := time.Now()
	aggregateParams := &database.AggregateVideoViewerLogsParams{
		VideoID:      "video-id",
		Interval:     entity.AggregateVideoViewerLogIntervalMinute,
		CreatedAtGte: now,
		CreatedAtLt:  now,
	}
	logs := entity.AggregatedVideoViewerLogs{
		{
			VideoID:    "video-id",
			ReportedAt: now,
			Total:      3,
		},
	}
	totalParams := &database.GetVideoTotalViewersParams{
		VideoID:      "video-id",
		CreatedAtGte: now,
		CreatedAtLt:  now,
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *media.AggregateVideoViewerLogsInput
		expect      entity.AggregatedVideoViewerLogs
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.VideoViewerLog.EXPECT().
					Aggregate(gomock.Any(), aggregateParams).
					Return(logs, nil)
				mocks.db.VideoViewerLog.EXPECT().
					GetTotal(gomock.Any(), totalParams).
					Return(int64(3), nil)
			},
			input: &media.AggregateVideoViewerLogsInput{
				VideoID:      "video-id",
				Interval:     entity.AggregateVideoViewerLogIntervalMinute,
				CreatedAtGte: now,
				CreatedAtLt:  now,
			},
			expect:      logs,
			expectTotal: 3,
			expectErr:   nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.AggregateVideoViewerLogsInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to aggregate viewer logs",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.VideoViewerLog.EXPECT().
					Aggregate(gomock.Any(), aggregateParams).
					Return(nil, assert.AnError)
				mocks.db.VideoViewerLog.EXPECT().
					GetTotal(gomock.Any(), totalParams).
					Return(int64(3), nil)
			},
			input: &media.AggregateVideoViewerLogsInput{
				VideoID:      "video-id",
				Interval:     entity.AggregateVideoViewerLogIntervalMinute,
				CreatedAtGte: now,
				CreatedAtLt:  now,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to aggregate viewer logs",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.VideoViewerLog.EXPECT().
					Aggregate(gomock.Any(), aggregateParams).
					Return(logs, nil)
				mocks.db.VideoViewerLog.EXPECT().
					GetTotal(gomock.Any(), totalParams).
					Return(int64(0), assert.AnError)
			},
			input: &media.AggregateVideoViewerLogsInput{
				VideoID:      "video-id",
				Interval:     entity.AggregateVideoViewerLogIntervalMinute,
				CreatedAtGte: now,
				CreatedAtLt:  now,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				actual, total, err := service.AggregateVideoViewerLogs(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Equal(t, tt.expect, actual)
				assert.Equal(t, tt.expectTotal, total)
			}),
		)
	}
}
