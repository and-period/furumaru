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

func TestListVideoComments(t *testing.T) {
	t.Parallel()

	now := time.Now()
	params := &database.ListVideoCommentsParams{
		VideoID:      "video-id",
		WithDisabled: false,
		CreatedAtGte: now.Add(-time.Hour),
		CreatedAtLt:  now.Add(time.Hour),
		Limit:        20,
		NextToken:    "",
	}
	comments := entity.VideoComments{
		{
			ID:        "comment-id",
			VideoID:   "video-id",
			UserID:    "user-id",
			Content:   "面白かった",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *media.ListVideoCommentsInput
		expect      entity.VideoComments
		expectToken string
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.VideoComment.EXPECT().List(ctx, params).Return(comments, "", nil)
			},
			input: &media.ListVideoCommentsInput{
				VideoID:      "video-id",
				WithDisabled: false,
				CreatedAtGte: now.Add(-time.Hour),
				CreatedAtLt:  now.Add(time.Hour),
				Limit:        20,
				NextToken:    "",
			},
			expect:      comments,
			expectToken: "",
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &media.ListVideoCommentsInput{},
			expect:      nil,
			expectToken: "",
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list video comments",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.VideoComment.EXPECT().List(ctx, params).Return(nil, "", assert.AnError)
			},
			input: &media.ListVideoCommentsInput{
				VideoID:      "video-id",
				WithDisabled: false,
				CreatedAtGte: now.Add(-time.Hour),
				CreatedAtLt:  now.Add(time.Hour),
				Limit:        20,
				NextToken:    "",
			},
			expect:      nil,
			expectToken: "",
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, token, err := service.ListVideoComments(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.expectToken, token)
		}))
	}
}

func TestCreateVideoComment(t *testing.T) {
	t.Parallel()

	now := time.Now()
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

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.CreateVideoCommentInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().Get(ctx, "video-id").Return(video, nil)
				mocks.db.VideoComment.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, comment *entity.VideoComment) error {
						expect := &entity.VideoComment{
							ID:       comment.ID, // ignore
							VideoID:  "video-id",
							UserID:   "user-id",
							Content:  "面白かった",
							Disabled: false,
						}
						assert.Equal(t, expect, comment)
						return nil
					})
			},
			input: &media.CreateVideoCommentInput{
				VideoID: "video-id",
				UserID:  "user-id",
				Content: "面白かった",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.CreateVideoCommentInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get video",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().Get(ctx, "video-id").Return(nil, assert.AnError)
			},
			input: &media.CreateVideoCommentInput{
				VideoID: "video-id",
				UserID:  "user-id",
				Content: "面白かった",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "unpublished video",
			setup: func(ctx context.Context, mocks *mocks) {
				video := &entity.Video{Status: entity.VideoStatusPrivate}
				mocks.db.Video.EXPECT().Get(ctx, "video-id").Return(video, nil)
			},
			input: &media.CreateVideoCommentInput{
				VideoID: "video-id",
				UserID:  "user-id",
				Content: "面白かった",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to create video comment",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Video.EXPECT().Get(ctx, "video-id").Return(video, nil)
				mocks.db.VideoComment.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &media.CreateVideoCommentInput{
				VideoID: "video-id",
				UserID:  "user-id",
				Content: "面白かった",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateVideoComment(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateVideoComment(t *testing.T) {
	t.Parallel()

	params := &database.UpdateVideoCommentParams{
		Disabled: true,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UpdateVideoCommentInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.VideoComment.EXPECT().Update(ctx, "comment-id", params).Return(nil)
			},
			input: &media.UpdateVideoCommentInput{
				CommentID: "comment-id",
				Disabled:  true,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UpdateVideoCommentInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update video comment",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.VideoComment.EXPECT().Update(ctx, "comment-id", params).Return(assert.AnError)
			},
			input: &media.UpdateVideoCommentInput{
				CommentID: "comment-id",
				Disabled:  true,
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateVideoComment(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
