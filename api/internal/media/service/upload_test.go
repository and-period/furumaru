package service

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUploadEvent(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.GetUploadEventInput
		expect    *entity.UploadEvent
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.UploadEvent{Key: "dir/media.png"}).
					DoAndReturn(func(ctx context.Context, event *entity.UploadEvent) error {
						event.Status = entity.UploadStatusSucceeded
						event.FileGroup = "dir"
						event.FileType = "image/png"
						event.UploadURL = "http://example-tmp.s3.ap-northeast-1.amazonaws.com/dir/media.png?query=test"
						event.ReferenceURL = "http://example.s3.ap-northeast-1.amazonaws.com/dir/media.png"
						event.ExpiredAt = now.Add(12 * time.Hour)
						event.CreatedAt = now.Add(-2 * time.Hour)
						event.UpdatedAt = now.Add(-2 * time.Hour)
						return nil
					})
			},
			input: &media.GetUploadEventInput{
				Key: "dir/media.png",
			},
			expect: &entity.UploadEvent{
				Key:          "dir/media.png",
				Status:       entity.UploadStatusSucceeded,
				FileGroup:    "dir",
				FileType:     "image/png",
				UploadURL:    "http://example-tmp.s3.ap-northeast-1.amazonaws.com/dir/media.png?query=test",
				ReferenceURL: "http://example.s3.ap-northeast-1.amazonaws.com/dir/media.png",
				ExpiredAt:    now.Add(12 * time.Hour),
				CreatedAt:    now.Add(-2 * time.Hour),
				UpdatedAt:    now.Add(-2 * time.Hour),
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.GetUploadEventInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get upload event",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().Get(ctx, &entity.UploadEvent{Key: "dir/media.png"}).Return(assert.AnError)
			},
			input: &media.GetUploadEventInput{
				Key: "dir/media.png",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetUploadEvent(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestGetBroadcastArchiveMP4UploadURL(t *testing.T) {
	t.Parallel()
	broadcast := &entity.Broadcast{
		ID:         "broadcast-id",
		ScheduleID: "schdule-id",
		Status:     entity.BroadcastStatusDisabled,
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateBroadcastArchiveMP4UploadInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				path := fmt.Sprintf(entity.BroadcastArchiveMP4Path, "schedule-id")
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				generateUploadURLMocks(mocks, t, path, "mp4", "video/mp4", nil)
			},
			input: &media.GenerateBroadcastArchiveMP4UploadInput{
				GenerateUploadURLInput: media.GenerateUploadURLInput{
					FileType: "video/mp4",
				},
				ScheduleID: "schedule-id",
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &media.GenerateBroadcastArchiveMP4UploadInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get broadcast",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &media.GenerateBroadcastArchiveMP4UploadInput{
				GenerateUploadURLInput: media.GenerateUploadURLInput{
					FileType: "video/mp4",
				},
				ScheduleID: "schedule-id",
			},
			expect: exception.ErrInternal,
		},
		{
			name: "failed precondition",
			setup: func(ctx context.Context, mocks *mocks) {
				broadcast := &entity.Broadcast{Status: entity.BroadcastStatusActive}
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
			},
			input: &media.GenerateBroadcastArchiveMP4UploadInput{
				GenerateUploadURLInput: media.GenerateUploadURLInput{
					FileType: "video/mp4",
				},
				ScheduleID: "schedule-id",
			},
			expect: exception.ErrFailedPrecondition,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetBroadcastArchiveMP4UploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetBroadcastLiveMP4UploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.BroadcastLiveMP4Path, "mp4", "video/mp4", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "video/mp4",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetBroadcastLiveMP4UploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetVideoThumbnailUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.VideoThumbnailPath, "png", "image/png", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetVideoThumbnailUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetVideoFileUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.VideoMP4Path, "mp4", "video/mp4", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "video/mp4",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetVideoFileUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetCoordinatorThumbnailUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.CoordinatorThumbnailPath, "png", "image/png", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetCoordinatorThumbnailUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetCoordinatorHeaderUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.CoordinatorHeaderPath, "png", "image/png", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetCoordinatorHeaderUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetCoordinatorPromotionVideoUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.CoordinatorPromotionVideoPath, "mp4", "video/mp4", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "video/mp4",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetCoordinatorPromotionVideoUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetCoordinatorBonusVideoUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.CoordinatorBonusVideoPath, "mp4", "video/mp4", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "video/mp4",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetCoordinatorBonusVideoUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetProducerThumbnailUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.ProducerThumbnailPath, "png", "image/png", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetProducerThumbnailUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetProducerHeaderUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.ProducerHeaderPath, "png", "image/png", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetProducerHeaderUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetProducerPromotionVideoUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.ProducerPromotionVideoPath, "mp4", "video/mp4", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "video/mp4",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetProducerPromotionVideoUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetProducerBonusVideoUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.ProducerBonusVideoPath, "mp4", "video/mp4", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "video/mp4",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetProducerBonusVideoUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetUserThumbnailUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.UserThumbnailPath, "png", "image/png", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetUserThumbnailUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetProductMediaImageUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.ProductMediaImagePath, "png", "image/png", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetProductMediaImageUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetProductMediaVideoUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.ProductMediaVideoPath, "mp4", "video/mp4", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "video/mp4",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetProductMediaVideoUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetProductTypeIconUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.ProductTypeIconPath, "png", "image/png", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetProductTypeIconUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetScheduleThumbnailUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.ScheduleThumbnailPath, "png", "image/png", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetScheduleThumbnailUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetScheduleImageUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.ScheduleImagePath, "png", "image/png", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetScheduleImageUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetScheduleOpeningVideoUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.ScheduleOpeningVideoPath, "mp4", "video/mp4", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "video/mp4",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetScheduleOpeningVideoUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetExperienceMediaImageUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.ExperienceMediaImagePath, "png", "image/png", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetExperienceMediaImageUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetExperienceMediaVideoUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.ExperienceMediaVideoPath, "mp4", "video/mp4", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "video/mp4",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetExperienceMediaVideoUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetExperiencePromotionVideoUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.ExperiencePromotionVideoPath, "mp4", "video/mp4", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "video/mp4",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetExperiencePromotionVideoUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestGetSpotThumbnailUploadURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.GenerateUploadURLInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				generateUploadURLMocks(mocks, t, entity.SpotThumbnailPath, "png", "image/png", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.GetSpotThumbnailUploadURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func generateUploadURLMocks(mocks *mocks, t *testing.T, path, ext, contentType string, err error) {
	url := "http://example.com/media." + ext
	headers := map[string][]string{
		"Content-Type": {contentType},
	}
	mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any()).
		DoAndReturn(func(params *storage.GeneratePresignUploadURIParams) (string, map[string][]string, error) {
			assert.True(t, strings.HasPrefix(params.Key, path), params.Key)
			assert.True(t, strings.HasSuffix(params.Key, ext), params.Key)
			assert.Equal(t, params.ExpiresIn, 10*time.Minute)
			assert.Equal(t, params.ContentType, contentType)
			return url, headers, nil
		})
	mocks.cache.EXPECT().Insert(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, event *entity.UploadEvent) error {
			assert.Contains(t, event.FileType, ext)
			assert.Equal(t, event.UploadURL, url)
			return err
		})
}

func TestGenerateUploadURL(t *testing.T) {
	t.Parallel()
	now := time.Now()
	url := "http://example.com"
	headers := map[string][]string{
		"Content-Type": {"image/png"},
	}
	tests := []struct {
		name       string
		setup      func(ctx context.Context, mocks *mocks)
		input      *media.GenerateUploadURLInput
		regulation *entity.Regulation
		expect     *entity.UploadEvent
		expectErr  error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any()).
					DoAndReturn(func(params *storage.GeneratePresignUploadURIParams) (string, map[string][]string, error) {
						assert.True(t, strings.HasPrefix(params.Key, "coordinators/thumbnail/"), params.Key)
						assert.True(t, strings.HasSuffix(params.Key, ".png"), params.Key)
						assert.Equal(t, params.ExpiresIn, 10*time.Minute)
						assert.Equal(t, params.ContentType, "image/png")
						return url, headers, nil
					})
				mocks.cache.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			regulation: entity.CoordinatorThumbnailRegulation,
			expect: &entity.UploadEvent{
				Key:           "", // ignore
				Status:        entity.UploadStatusWaiting,
				FileGroup:     "coordinators/thumbnail",
				FileType:      "image/png",
				UploadURL:     "http://example.com",
				UploadHeaders: map[string][]string{"Content-Type": {"image/png"}},
				ExpiredAt:     now.Add(defaultUploadEventTTL),
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			expectErr: nil,
		},
		{
			name:       "invalid argument",
			setup:      func(ctx context.Context, mocks *mocks) {},
			input:      &media.GenerateUploadURLInput{},
			regulation: entity.CoordinatorThumbnailRegulation,
			expect:     nil,
			expectErr:  exception.ErrInvalidArgument,
		},
		{
			name:  "failed to get object key",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.GenerateUploadURLInput{
				FileType: "video/mp4",
			},
			regulation: entity.CoordinatorThumbnailRegulation,
			expect:     nil,
			expectErr:  exception.ErrInvalidArgument,
		},
		{
			name: "failed to generate presign upload uri",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any()).Return("", nil, assert.AnError)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			regulation: entity.CoordinatorThumbnailRegulation,
			expect:     nil,
			expectErr:  exception.ErrInternal,
		},
		{
			name: "failed to insert cache",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any()).Return(url, headers, nil)
				mocks.cache.EXPECT().Insert(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			regulation: entity.CoordinatorThumbnailRegulation,
			expect:     nil,
			expectErr:  exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.generateUploadURL(ctx, tt.input, tt.regulation)
			assert.ErrorIs(t, err, tt.expectErr)
			if err != nil {
				return
			}
			actual.Key = "" // ignore
			assert.Equal(t, tt.expect, actual)
		}, withNow(now)))
	}
}
