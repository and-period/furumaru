package service

import (
	"context"
	"io"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).
					DoAndReturn(func(key string, expiresIn time.Duration) (string, error) {
						assert.True(t, strings.HasPrefix(key, entity.CoordinatorThumbnailPath), key)
						assert.True(t, strings.HasSuffix(key, ".png"), key)
						return "http://example.com/image.png", nil
					})
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
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).
					DoAndReturn(func(key string, expiresIn time.Duration) (string, error) {
						assert.True(t, strings.HasPrefix(key, entity.CoordinatorHeaderPath), key)
						assert.True(t, strings.HasSuffix(key, ".png"), key)
						return "http://example.com/image.png", nil
					})
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
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).
					DoAndReturn(func(key string, expiresIn time.Duration) (string, error) {
						assert.True(t, strings.HasPrefix(key, entity.CoordinatorPromotionVideoPath), key)
						assert.True(t, strings.HasSuffix(key, ".mp4"), key)
						return "http://example.com/video.mp4", nil
					})
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
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).
					DoAndReturn(func(key string, expiresIn time.Duration) (string, error) {
						assert.True(t, strings.HasPrefix(key, entity.CoordinatorBonusVideoPath), key)
						assert.True(t, strings.HasSuffix(key, ".mp4"), key)
						return "http://example.com/video.mp4", nil
					})
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
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).
					DoAndReturn(func(key string, expiresIn time.Duration) (string, error) {
						assert.True(t, strings.HasPrefix(key, entity.ProducerThumbnailPath), key)
						assert.True(t, strings.HasSuffix(key, ".png"), key)
						return "http://example.com/image.png", nil
					})
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
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).
					DoAndReturn(func(key string, expiresIn time.Duration) (string, error) {
						assert.True(t, strings.HasPrefix(key, entity.ProducerHeaderPath), key)
						assert.True(t, strings.HasSuffix(key, ".png"), key)
						return "http://example.com/image.png", nil
					})
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
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).
					DoAndReturn(func(key string, expiresIn time.Duration) (string, error) {
						assert.True(t, strings.HasPrefix(key, entity.ProducerPromotionVideoPath), key)
						assert.True(t, strings.HasSuffix(key, ".mp4"), key)
						return "http://example.com/video.mp4", nil
					})
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
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).
					DoAndReturn(func(key string, expiresIn time.Duration) (string, error) {
						assert.True(t, strings.HasPrefix(key, entity.ProducerBonusVideoPath), key)
						assert.True(t, strings.HasSuffix(key, ".mp4"), key)
						return "http://example.com/video.mp4", nil
					})
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
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).
					DoAndReturn(func(key string, expiresIn time.Duration) (string, error) {
						assert.True(t, strings.HasPrefix(key, entity.UserThumbnailPath), key)
						assert.True(t, strings.HasSuffix(key, ".png"), key)
						return "http://example.com/image.png", nil
					})
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
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).
					DoAndReturn(func(key string, expiresIn time.Duration) (string, error) {
						assert.True(t, strings.HasPrefix(key, entity.ProductMediaImagePath), key)
						assert.True(t, strings.HasSuffix(key, ".png"), key)
						return "http://example.com/image.png", nil
					})
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
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).
					DoAndReturn(func(key string, expiresIn time.Duration) (string, error) {
						assert.True(t, strings.HasPrefix(key, entity.ProductMediaVideoPath), key)
						assert.True(t, strings.HasSuffix(key, ".mp4"), key)
						return "http://example.com/video.mp4", nil
					})
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
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).
					DoAndReturn(func(key string, expiresIn time.Duration) (string, error) {
						assert.True(t, strings.HasPrefix(key, entity.ProductTypeIconPath), key)
						assert.True(t, strings.HasSuffix(key, ".png"), key)
						return "http://example.com/image.png", nil
					})
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
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).
					DoAndReturn(func(key string, expiresIn time.Duration) (string, error) {
						assert.True(t, strings.HasPrefix(key, entity.ScheduleThumbnailPath), key)
						assert.True(t, strings.HasSuffix(key, ".png"), key)
						return "http://example.com/image.png", nil
					})
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
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).
					DoAndReturn(func(key string, expiresIn time.Duration) (string, error) {
						assert.True(t, strings.HasPrefix(key, entity.ScheduleImagePath), key)
						assert.True(t, strings.HasSuffix(key, ".png"), key)
						return "http://example.com/image.png", nil
					})
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
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).
					DoAndReturn(func(key string, expiresIn time.Duration) (string, error) {
						assert.True(t, strings.HasPrefix(key, entity.ScheduleOpeningVideoPath), key)
						assert.True(t, strings.HasSuffix(key, ".mp4"), key)
						return "http://example.com/video.mp4", nil
					})
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

func TestGenerateUploadURL(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(ctx context.Context, mocks *mocks)
		input      *media.GenerateUploadURLInput
		regulation *entity.Regulation
		expect     string
		expectErr  error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).Return("http://example.com", nil)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			regulation: entity.CoordinatorThumbnailRegulation,
			expect:     "http://example.com",
			expectErr:  nil,
		},
		{
			name:       "invalid argument",
			setup:      func(ctx context.Context, mocks *mocks) {},
			input:      &media.GenerateUploadURLInput{},
			regulation: entity.CoordinatorThumbnailRegulation,
			expect:     "",
			expectErr:  exception.ErrInvalidArgument,
		},
		{
			name:  "failed to get object key",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.GenerateUploadURLInput{
				FileType: "video/mp4",
			},
			regulation: entity.CoordinatorThumbnailRegulation,
			expect:     "",
			expectErr:  exception.ErrInvalidArgument,
		},
		{
			name: "failed to generate presign upload uri",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().GeneratePresignUploadURI(gomock.Any(), 10*time.Minute).Return("", assert.AnError)
			},
			input: &media.GenerateUploadURLInput{
				FileType: "image/png",
			},
			regulation: entity.CoordinatorThumbnailRegulation,
			expect:     "",
			expectErr:  exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.generateUploadURL(tt.input, tt.regulation)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestGenerateCoordinatorThumbnail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     func() *media.GenerateFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, path string, file io.Reader) (string, error) {
						assert.True(t, strings.HasPrefix(path, entity.CoordinatorThumbnailPath), path)
						u, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
						require.NoError(t, err)
						return u.String(), nil
					})
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    strings.Join([]string{tmpURL, entity.CoordinatorThumbnailPath}, "/"),
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: func() *media.GenerateFileInput {
				return &media.GenerateFileInput{}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid file regulation",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				header.Size = 10<<20 + 1
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GenerateCoordinatorThumbnail(ctx, tt.input())
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Contains(t, actual, tt.expect)
		}))
	}
}

func TestGenerateCoordinatorHeader(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     func() *media.GenerateFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, path string, file io.Reader) (string, error) {
						assert.True(t, strings.HasPrefix(path, entity.CoordinatorHeaderPath), path)
						u, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
						require.NoError(t, err)
						return u.String(), nil
					})
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    strings.Join([]string{tmpURL, entity.CoordinatorHeaderPath}, "/"),
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: func() *media.GenerateFileInput {
				return &media.GenerateFileInput{}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid file regulation",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				header.Size = 10<<20 + 1
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GenerateCoordinatorHeader(ctx, tt.input())
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Contains(t, actual, tt.expect)
		}))
	}
}

func TestGenerateCoordinatorPromotionVideo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     func() *media.GenerateFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, path string, file io.Reader) (string, error) {
						assert.True(t, strings.HasPrefix(path, entity.CoordinatorPromotionVideoPath), path)
						u, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
						require.NoError(t, err)
						return u.String(), nil
					})
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    strings.Join([]string{tmpURL, entity.CoordinatorPromotionVideoPath}, "/"),
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: func() *media.GenerateFileInput {
				return &media.GenerateFileInput{}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid file regulation",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				header.Size = 200<<20 + 1
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GenerateCoordinatorPromotionVideo(ctx, tt.input())
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Contains(t, actual, tt.expect)
		}))
	}
}

func TestGenerateCoordinatorBonusVideo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     func() *media.GenerateFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, path string, file io.Reader) (string, error) {
						assert.True(t, strings.HasPrefix(path, entity.CoordinatorBonusVideoPath), path)
						u, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
						require.NoError(t, err)
						return u.String(), nil
					})
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    strings.Join([]string{tmpURL, entity.CoordinatorBonusVideoPath}, "/"),
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: func() *media.GenerateFileInput {
				return &media.GenerateFileInput{}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid file regulation",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				header.Size = 200<<20 + 1
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GenerateCoordinatorBonusVideo(ctx, tt.input())
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Contains(t, actual, tt.expect)
		}))
	}
}

func TestGenerateProducerThumbnail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     func() *media.GenerateFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, path string, file io.Reader) (string, error) {
						assert.True(t, strings.HasPrefix(path, entity.ProducerThumbnailPath), path)
						u, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
						require.NoError(t, err)
						return u.String(), nil
					})
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    strings.Join([]string{tmpURL, entity.ProducerThumbnailPath}, "/"),
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: func() *media.GenerateFileInput {
				return &media.GenerateFileInput{}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid file regulation",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				header.Size = 10<<20 + 1
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GenerateProducerThumbnail(ctx, tt.input())
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Contains(t, actual, tt.expect)
		}))
	}
}

func TestGenerateProducerHeader(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     func() *media.GenerateFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, path string, file io.Reader) (string, error) {
						assert.True(t, strings.HasPrefix(path, entity.ProducerHeaderPath), path)
						u, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
						require.NoError(t, err)
						return u.String(), nil
					})
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    strings.Join([]string{tmpURL, entity.ProducerHeaderPath}, "/"),
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: func() *media.GenerateFileInput {
				return &media.GenerateFileInput{}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid file regulation",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				header.Size = 10<<20 + 1
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GenerateProducerHeader(ctx, tt.input())
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Contains(t, actual, tt.expect)
		}))
	}
}

func TestGenerateProducerPromotionVideo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     func() *media.GenerateFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, path string, file io.Reader) (string, error) {
						assert.True(t, strings.HasPrefix(path, entity.ProducerPromotionVideoPath), path)
						u, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
						require.NoError(t, err)
						return u.String(), nil
					})
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    strings.Join([]string{tmpURL, entity.ProducerPromotionVideoPath}, "/"),
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: func() *media.GenerateFileInput {
				return &media.GenerateFileInput{}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid file regulation",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				header.Size = 200<<20 + 1
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GenerateProducerPromotionVideo(ctx, tt.input())
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Contains(t, actual, tt.expect)
		}))
	}
}

func TestGenerateProducerBonusVideo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     func() *media.GenerateFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, path string, file io.Reader) (string, error) {
						assert.True(t, strings.HasPrefix(path, entity.ProducerBonusVideoPath), path)
						u, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
						require.NoError(t, err)
						return u.String(), nil
					})
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    strings.Join([]string{tmpURL, entity.ProducerBonusVideoPath}, "/"),
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: func() *media.GenerateFileInput {
				return &media.GenerateFileInput{}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid file regulation",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				header.Size = 200<<20 + 1
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GenerateProducerBonusVideo(ctx, tt.input())
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Contains(t, actual, tt.expect)
		}))
	}
}

func TestGenerateUserThumbnail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     func() *media.GenerateFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, path string, file io.Reader) (string, error) {
						assert.True(t, strings.HasPrefix(path, entity.UserThumbnailPath), path)
						u, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
						require.NoError(t, err)
						return u.String(), nil
					})
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    strings.Join([]string{tmpURL, entity.UserThumbnailPath}, "/"),
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: func() *media.GenerateFileInput {
				return &media.GenerateFileInput{}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid file regulation",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				header.Size = 10<<20 + 1
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GenerateUserThumbnail(ctx, tt.input())
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Contains(t, actual, tt.expect)
		}))
	}
}

func TestGenerateProductMediaImage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     func() *media.GenerateFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, path string, file io.Reader) (string, error) {
						assert.True(t, strings.HasPrefix(path, entity.ProductMediaImagePath), path)
						u, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
						require.NoError(t, err)
						return u.String(), nil
					})
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    strings.Join([]string{tmpURL, entity.ProductMediaImagePath}, "/"),
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: func() *media.GenerateFileInput {
				return &media.GenerateFileInput{}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid file regulation",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				header.Size = 10<<20 + 1
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GenerateProductMediaImage(ctx, tt.input())
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Contains(t, actual, tt.expect)
		}))
	}
}

func TestGenerateProductMediaVideo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     func() *media.GenerateFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, path string, file io.Reader) (string, error) {
						assert.True(t, strings.HasPrefix(path, entity.ProductMediaVideoPath), path)
						u, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
						require.NoError(t, err)
						return u.String(), nil
					})
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    strings.Join([]string{tmpURL, entity.ProductMediaVideoPath}, "/"),
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: func() *media.GenerateFileInput {
				return &media.GenerateFileInput{}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid file regulation",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				header.Size = 200<<20 + 1
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GenerateProductMediaVideo(ctx, tt.input())
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Contains(t, actual, tt.expect)
		}))
	}
}

func TestGenerateProductTypeIcon(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     func() *media.GenerateFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, path string, file io.Reader) (string, error) {
						assert.True(t, strings.HasPrefix(path, entity.ProductTypeIconPath), path)
						u, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
						require.NoError(t, err)
						return u.String(), nil
					})
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    strings.Join([]string{tmpURL, entity.ProductTypeIconPath}, "/"),
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: func() *media.GenerateFileInput {
				return &media.GenerateFileInput{}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid file regulation",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				header.Size = 10<<20 + 1
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GenerateProductTypeIcon(ctx, tt.input())
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Contains(t, actual, tt.expect)
		}))
	}
}

func TestGenerateScheduleThumbnail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     func() *media.GenerateFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, path string, file io.Reader) (string, error) {
						assert.True(t, strings.HasPrefix(path, entity.ScheduleThumbnailPath), path)
						u, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
						require.NoError(t, err)
						return u.String(), nil
					})
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    strings.Join([]string{tmpURL, entity.ScheduleThumbnailPath}, "/"),
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: func() *media.GenerateFileInput {
				return &media.GenerateFileInput{}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid file regulation",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				header.Size = 10<<20 + 1
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GenerateScheduleThumbnail(ctx, tt.input())
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Contains(t, actual, tt.expect)
		}))
	}
}

func TestGenerateScheduleImage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     func() *media.GenerateFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, path string, file io.Reader) (string, error) {
						assert.True(t, strings.HasPrefix(path, entity.ScheduleImagePath), path)
						u, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
						require.NoError(t, err)
						return u.String(), nil
					})
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    strings.Join([]string{tmpURL, entity.ScheduleImagePath}, "/"),
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: func() *media.GenerateFileInput {
				return &media.GenerateFileInput{}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid file regulation",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				header.Size = 10<<20 + 1
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			input: func() *media.GenerateFileInput {
				file, header := testImageFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GenerateScheduleImage(ctx, tt.input())
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Contains(t, actual, tt.expect)
		}))
	}
}

func TestGenerateScheduleOpeningVideo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     func() *media.GenerateFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, path string, file io.Reader) (string, error) {
						assert.True(t, strings.HasPrefix(path, entity.ScheduleOpeningVideoPath), path)
						u, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
						require.NoError(t, err)
						return u.String(), nil
					})
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    strings.Join([]string{tmpURL, entity.ScheduleOpeningVideoPath}, "/"),
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: func() *media.GenerateFileInput {
				return &media.GenerateFileInput{}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid file regulation",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				header.Size = 200<<20 + 1
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			input: func() *media.GenerateFileInput {
				file, header := testVideoFile(t)
				return &media.GenerateFileInput{File: file, Header: header}
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GenerateScheduleOpeningVideo(ctx, tt.input())
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Contains(t, actual, tt.expect)
		}))
	}
}
