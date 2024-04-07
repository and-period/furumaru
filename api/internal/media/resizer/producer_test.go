package resizer

import (
	"bytes"
	"context"
	"image"
	"testing"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestProducerThumbnail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		payload   *entity.ResizerPayload
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				file := testImageFile(t)
				url := "http://example.com/media/image_xxx.png"
				thumbnails := common.Images{
					{URL: url, Size: common.ImageSizeSmall},
					{URL: url, Size: common.ImageSizeMedium},
					{URL: url, Size: common.ImageSizeLarge},
				}
				mocks.storage.EXPECT().Download(ctx, "http://example.com/media/image.png").Return(file, nil)
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(url, nil).Times(3)
				mocks.user.EXPECT().
					UpdateProducerThumbnails(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, in *user.UpdateProducerThumbnailsInput) error {
						assert.Equal(t, "target-id", in.ProducerID)
						assert.ElementsMatch(t, thumbnails, in.Thumbnails)
						return nil
					})
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProducerThumbnail,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: nil,
		},
		{
			name: "failed to empty url",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProducerThumbnail,
				URLs:     []string{},
			},
			expectErr: errRequiredMediaURL,
		},
		{
			name: "failed to download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, "http://example.com/media/image.png").Return(nil, assert.AnError)
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProducerThumbnail,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: assert.AnError,
		},
		{
			name: "failed to resize images",
			setup: func(ctx context.Context, mocks *mocks) {
				file := &bytes.Buffer{}
				mocks.storage.EXPECT().Download(ctx, "http://example.com/media/image.png").Return(file, nil)
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProducerThumbnail,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: image.ErrFormat,
		},
		{
			name: "failed to upload images",
			setup: func(ctx context.Context, mocks *mocks) {
				file := testImageFile(t)
				mocks.storage.EXPECT().Download(ctx, "http://example.com/media/image.png").Return(file, nil)
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", assert.AnError).AnyTimes()
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProducerThumbnail,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: assert.AnError,
		},
		{
			name: "failed to update producer thumbnails",
			setup: func(ctx context.Context, mocks *mocks) {
				file := testImageFile(t)
				url := "http://example.com/media/image_xxx.png"
				mocks.storage.EXPECT().Download(ctx, "http://example.com/media/image.png").Return(file, nil)
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(url, nil).Times(3)
				mocks.user.EXPECT().UpdateProducerThumbnails(ctx, gomock.Any()).Return(assert.AnError)
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProducerThumbnail,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testResizer(tt.setup, func(ctx context.Context, t *testing.T, resizer *resizer) {
			t.Parallel()
			err := resizer.producerThumbnail(ctx, tt.payload)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestProducerHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		payload   *entity.ResizerPayload
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				file := testImageFile(t)
				url := "http://example.com/media/image_xxx.png"
				headers := common.Images{
					{URL: url, Size: common.ImageSizeSmall},
					{URL: url, Size: common.ImageSizeMedium},
					{URL: url, Size: common.ImageSizeLarge},
				}
				mocks.storage.EXPECT().Download(ctx, "http://example.com/media/image.png").Return(file, nil)
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(url, nil).Times(3)
				mocks.user.EXPECT().
					UpdateProducerHeaders(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, in *user.UpdateProducerHeadersInput) error {
						assert.Equal(t, "target-id", in.ProducerID)
						assert.ElementsMatch(t, headers, in.Headers)
						return nil
					})
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProducerHeader,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: nil,
		},
		{
			name:  "failed to empty url",
			setup: func(ctx context.Context, mocks *mocks) {},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProducerHeader,
				URLs:     []string{},
			},
			expectErr: errRequiredMediaURL,
		},
		{
			name: "failed to download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, "http://example.com/media/image.png").Return(nil, assert.AnError)
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProducerHeader,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: assert.AnError,
		},
		{
			name: "failed to resize images",
			setup: func(ctx context.Context, mocks *mocks) {
				file := &bytes.Buffer{}
				mocks.storage.EXPECT().Download(ctx, "http://example.com/media/image.png").Return(file, nil)
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProducerHeader,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: image.ErrFormat,
		},
		{
			name: "failed to upload images",
			setup: func(ctx context.Context, mocks *mocks) {
				file := testImageFile(t)
				mocks.storage.EXPECT().Download(ctx, "http://example.com/media/image.png").Return(file, nil)
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", assert.AnError).AnyTimes()
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProducerHeader,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: assert.AnError,
		},
		{
			name: "failed to update producer headers",
			setup: func(ctx context.Context, mocks *mocks) {
				file := testImageFile(t)
				url := "http://example.com/media/image_xxx.png"
				mocks.storage.EXPECT().Download(ctx, "http://example.com/media/image.png").Return(file, nil)
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(url, nil).Times(3)
				mocks.user.EXPECT().UpdateProducerHeaders(ctx, gomock.Any()).Return(assert.AnError)
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProducerHeader,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testResizer(tt.setup, func(ctx context.Context, t *testing.T, resizer *resizer) {
			t.Parallel()
			err := resizer.producerHeader(ctx, tt.payload)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
