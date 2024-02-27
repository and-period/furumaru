package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListBroadcastComments(t *testing.T) {
	t.Parallel()

	now := time.Now()
	broadcast := &entity.Broadcast{
		ID:            "broadcast-id",
		ScheduleID:    "schedule-id",
		CoordinatorID: "coordinator-id",
		Status:        entity.BroadcastStatusIdle,
		InputURL:      "rtmp://127.0.0.1:1935/app/instance",
		OutputURL:     "http://example.com/index.m3u8",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	params := &database.ListBroadcastCommentsParams{
		BroadcastID:  "broadcast-id",
		WithDisabled: false,
		CreatedAtGte: now.Add(-time.Hour),
		CreatedAtLt:  now,
		Limit:        20,
		NextToken:    "next-token",
		Orders: []*database.ListBroadcastCommentsOrder{{
			Key:        entity.BroadcastCommentOrderByCreatedAt,
			OrderByASC: false,
		}},
	}
	comments := entity.BroadcastComments{
		{
			ID:          "comment-id",
			BroadcastID: "broadcast-id",
			UserID:      "user-id",
			Content:     "こんにちは",
			Disabled:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *media.ListBroadcastCommentsInput
		expect      entity.BroadcastComments
		expectToken string
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.db.BroadcastComment.EXPECT().List(ctx, params).Return(comments, "next-token", nil)
			},
			input: &media.ListBroadcastCommentsInput{
				ScheduleID:   "schedule-id",
				CreatedAtGte: now.Add(-time.Hour),
				CreatedAtLt:  now,
				Limit:        20,
				NextToken:    "next-token",
				Orders: []*media.ListBroadcastCommentsOrder{{
					Key:        entity.BroadcastCommentOrderByCreatedAt,
					OrderByASC: false,
				}},
			},
			expect: entity.BroadcastComments{
				{
					ID:          "comment-id",
					BroadcastID: "broadcast-id",
					UserID:      "user-id",
					Content:     "こんにちは",
					Disabled:    false,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			expectToken: "next-token",
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &media.ListBroadcastCommentsInput{},
			expect:      nil,
			expectToken: "",
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to get broadcast",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &media.ListBroadcastCommentsInput{
				ScheduleID:   "schedule-id",
				CreatedAtGte: now.Add(-time.Hour),
				CreatedAtLt:  now,
				Limit:        20,
				NextToken:    "next-token",
				Orders: []*media.ListBroadcastCommentsOrder{{
					Key:        entity.BroadcastCommentOrderByCreatedAt,
					OrderByASC: false,
				}},
			},
			expect:      nil,
			expectToken: "",
			expectErr:   exception.ErrInternal,
		},
		{
			name: "archive is fixed",
			setup: func(ctx context.Context, mocks *mocks) {
				broadcast := &entity.Broadcast{ArchiveFixed: true}
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
			},
			input: &media.ListBroadcastCommentsInput{
				ScheduleID:   "schedule-id",
				CreatedAtGte: now.Add(-time.Hour),
				CreatedAtLt:  now,
				Limit:        20,
				NextToken:    "next-token",
				Orders: []*media.ListBroadcastCommentsOrder{{
					Key:        entity.BroadcastCommentOrderByCreatedAt,
					OrderByASC: false,
				}},
			},
			expect:      entity.BroadcastComments{},
			expectToken: "",
			expectErr:   nil,
		},
		{
			name: "failed to list broadcast comments",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.db.BroadcastComment.EXPECT().List(ctx, params).Return(nil, "", assert.AnError)
			},
			input: &media.ListBroadcastCommentsInput{
				ScheduleID:   "schedule-id",
				CreatedAtGte: now.Add(-time.Hour),
				CreatedAtLt:  now,
				Limit:        20,
				NextToken:    "next-token",
				Orders: []*media.ListBroadcastCommentsOrder{{
					Key:        entity.BroadcastCommentOrderByCreatedAt,
					OrderByASC: false,
				}},
			},
			expect:      nil,
			expectToken: "",
			expectErr:   exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, token, err := service.ListBroadcastComments(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.expectToken, token)
		}))
	}
}

func TestCreateBroadcastComment(t *testing.T) {
	t.Parallel()

	now := time.Now()
	broadcast := &entity.Broadcast{
		ID:            "broadcast-id",
		ScheduleID:    "schedule-id",
		CoordinatorID: "coordinator-id",
		Status:        entity.BroadcastStatusIdle,
		InputURL:      "rtmp://127.0.0.1:1935/app/instance",
		OutputURL:     "http://example.com/index.m3u8",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.CreateBroadcastCommentInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.db.BroadcastComment.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, comment *entity.BroadcastComment) error {
						expect := &entity.BroadcastComment{
							ID:          comment.ID, // ignore
							BroadcastID: "broadcast-id",
							UserID:      "user-id",
							Content:     "こんにちは",
							Disabled:    false,
						}
						assert.Equal(t, expect, comment)
						return nil
					})
			},
			input: &media.CreateBroadcastCommentInput{
				ScheduleID: "schedule-id",
				UserID:     "user-id",
				Content:    "こんにちは",
			},
			expectErr: nil,
		},
		{
			name: "broadcast is disabled",
			setup: func(ctx context.Context, mocks *mocks) {
				broadcast := &entity.Broadcast{Status: entity.BroadcastStatusDisabled}
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
			},
			input: &media.CreateBroadcastCommentInput{
				ScheduleID: "schedule-id",
				UserID:     "user-id",
				Content:    "こんにちは",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to get broadcast",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &media.CreateBroadcastCommentInput{
				ScheduleID: "schedule-id",
				UserID:     "user-id",
				Content:    "こんにちは",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create broadcast comment",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.db.BroadcastComment.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &media.CreateBroadcastCommentInput{
				ScheduleID: "schedule-id",
				UserID:     "user-id",
				Content:    "こんにちは",
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateBroadcastComment(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestCreateBroadcastGuestComment(t *testing.T) {
	t.Parallel()

	now := time.Now()
	broadcast := &entity.Broadcast{
		ID:            "broadcast-id",
		ScheduleID:    "schedule-id",
		CoordinatorID: "coordinator-id",
		Status:        entity.BroadcastStatusIdle,
		InputURL:      "rtmp://127.0.0.1:1935/app/instance",
		OutputURL:     "http://example.com/index.m3u8",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.CreateBroadcastGuestCommentInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.db.BroadcastComment.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, comment *entity.BroadcastComment) error {
						expect := &entity.BroadcastComment{
							ID:          comment.ID, // ignore
							BroadcastID: "broadcast-id",
							UserID:      "",
							Content:     "こんにちは",
							Disabled:    false,
						}
						assert.Equal(t, expect, comment)
						return nil
					})
			},
			input: &media.CreateBroadcastGuestCommentInput{
				ScheduleID: "schedule-id",
				Content:    "こんにちは",
			},
			expectErr: nil,
		},
		{
			name: "broadcast is disabled",
			setup: func(ctx context.Context, mocks *mocks) {
				broadcast := &entity.Broadcast{Status: entity.BroadcastStatusDisabled}
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
			},
			input: &media.CreateBroadcastGuestCommentInput{
				ScheduleID: "schedule-id",
				Content:    "こんにちは",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to get broadcast",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &media.CreateBroadcastGuestCommentInput{
				ScheduleID: "schedule-id",
				Content:    "こんにちは",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create broadcast comment",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.db.BroadcastComment.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &media.CreateBroadcastGuestCommentInput{
				ScheduleID: "schedule-id",
				Content:    "こんにちは",
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateBroadcastGuestComment(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
