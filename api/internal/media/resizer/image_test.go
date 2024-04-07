package resizer

import (
	"bytes"
	"context"
	"image"
	"io"
	"testing"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUploadImages(t *testing.T) {
	t.Parallel()

	img := testImageFile(t)

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		originURL string
		images    map[common.ImageSize]io.Reader
		expect    common.Images
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, path string, body io.Reader) (string, error) {
						expect := []string{
							"media/image_240.png",
							"media/image_675.png",
							"media/image_900.png",
						}
						assert.Contains(t, expect, path)
						return "http://example.com/media/image_xxx.png", nil
					}).Times(3)
			},
			originURL: "http://example.com/media/image.png",
			images: map[common.ImageSize]io.Reader{
				common.ImageSizeSmall:  img,
				common.ImageSizeMedium: img,
				common.ImageSizeLarge:  img,
			},
			expect: common.Images{
				{URL: "http://example.com/media/image_xxx.png", Size: common.ImageSizeSmall},
				{URL: "http://example.com/media/image_xxx.png", Size: common.ImageSizeMedium},
				{URL: "http://example.com/media/image_xxx.png", Size: common.ImageSizeLarge},
			},
			expectErr: nil,
		},
		{
			name:      "success for empty",
			setup:     func(ctx context.Context, mocks *mocks) {},
			originURL: "http://example.com/media/image.png",
			images:    map[common.ImageSize]io.Reader{},
			expect:    common.Images{},
			expectErr: nil,
		},
		{
			name:      "success for empty",
			setup:     func(ctx context.Context, mocks *mocks) {},
			originURL: "http://example.com/media/image.png",
			images:    map[common.ImageSize]io.Reader{},
			expect:    common.Images{},
			expectErr: nil,
		},
		{
			name:      "failed to generate file path",
			setup:     func(ctx context.Context, mocks *mocks) {},
			originURL: "http://example.com/media/image.png",
			images: map[common.ImageSize]io.Reader{
				common.ImageSizeUnknown: img,
			},
			expect:    nil,
			expectErr: errUnsupportedImageSize,
		},
		{
			name: "failed to upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", assert.AnError).AnyTimes()
			},
			originURL: "http://example.com/media/image.png",
			images: map[common.ImageSize]io.Reader{
				common.ImageSizeSmall:  img,
				common.ImageSizeMedium: img,
				common.ImageSizeLarge:  img,
			},
			expect:    nil,
			expectErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testResizer(tt.setup, func(ctx context.Context, t *testing.T, resizer *resizer) {
			t.Parallel()
			actual, err := resizer.uploadImages(ctx, tt.originURL, tt.images)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGenerateFilePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		originURL string
		size      common.ImageSize
		expect    string
		expectErr error
	}{
		{
			name:      "success to image size small",
			setup:     func(ctx context.Context, mocks *mocks) {},
			originURL: "http://example.com/images/image.png",
			size:      common.ImageSizeSmall,
			expect:    "images/image_240.png",
			expectErr: nil,
		},
		{
			name:      "success to image size medium",
			setup:     func(ctx context.Context, mocks *mocks) {},
			originURL: "http://example.com/images/image.png",
			size:      common.ImageSizeMedium,
			expect:    "images/image_675.png",
			expectErr: nil,
		},
		{
			name:      "success to image size large",
			setup:     func(ctx context.Context, mocks *mocks) {},
			originURL: "http://example.com/images/image.png",
			size:      common.ImageSizeLarge,
			expect:    "images/image_900.png",
			expectErr: nil,
		},
		{
			name:      "failed to generage url for unknown image size",
			setup:     func(ctx context.Context, mocks *mocks) {},
			originURL: "images/image.png",
			size:      common.ImageSizeUnknown,
			expect:    "",
			expectErr: errUnsupportedImageSize,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testResizer(tt.setup, func(ctx context.Context, t *testing.T, resizer *resizer) {
			actual, err := resizer.generateFilePath(tt.originURL, tt.size)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestResizeImages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		originURL string
		file      io.Reader
		expectLen int
		expectErr error
	}{
		{
			name:      "success",
			setup:     func(ctx context.Context, mocks *mocks) {},
			originURL: "http://example.com/media/image.png",
			file:      testImageFile(t),
			expectLen: 3,
			expectErr: nil,
		},
		{
			name:      "success for no need to resize",
			setup:     func(ctx context.Context, mocks *mocks) {},
			originURL: "http://example.com/media/video.mp4",
			file:      testVideoFile(t),
			expectLen: 0,
			expectErr: nil,
		},
		{
			name:      "success for source image is empty",
			setup:     func(ctx context.Context, mocks *mocks) {},
			originURL: "http://example.com/media/image.png",
			file:      nil,
			expectLen: 0,
			expectErr: nil,
		},
		{
			name:      "failed to resize image",
			setup:     func(ctx context.Context, mocks *mocks) {},
			originURL: "http://example.com/media/image.png",
			file:      &bytes.Buffer{},
			expectLen: 0,
			expectErr: image.ErrFormat,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testResizer(tt.setup, func(ctx context.Context, t *testing.T, resizer *resizer) {
			actual, err := resizer.resizeImages(tt.originURL, tt.file)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Len(t, actual, tt.expectLen)
		}))
	}
}

func TestResizeImage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		file      io.Reader
		size      common.ImageSize
		expect    string
		expectErr error
	}{
		{
			name:      "success to small",
			setup:     func(ctx context.Context, mocks *mocks) {},
			file:      testImageFile(t),
			size:      common.ImageSizeSmall,
			expect:    "and-period_240.png",
			expectErr: nil,
		},
		{
			name:      "success to medium",
			setup:     func(ctx context.Context, mocks *mocks) {},
			file:      testImageFile(t),
			size:      common.ImageSizeMedium,
			expect:    "and-period_675.png",
			expectErr: nil,
		},
		{
			name:      "success to large",
			setup:     func(ctx context.Context, mocks *mocks) {},
			file:      testImageFile(t),
			size:      common.ImageSizeMedium,
			expect:    "and-period_900.png",
			expectErr: nil,
		},
		{
			name:      "failed to decode",
			setup:     func(ctx context.Context, mocks *mocks) {},
			file:      &bytes.Buffer{},
			size:      common.ImageSizeSmall,
			expect:    "",
			expectErr: image.ErrFormat,
		},
		{
			name:      "failed to invalid size",
			setup:     func(ctx context.Context, mocks *mocks) {},
			file:      testImageFile(t),
			size:      common.ImageSizeUnknown,
			expect:    "",
			expectErr: errUnsupportedImageSize,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testResizer(tt.setup, func(ctx context.Context, t *testing.T, resizer *resizer) {
			_, err := resizer.resizeImage(tt.file, tt.size)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
