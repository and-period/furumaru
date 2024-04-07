package resizer

import (
	"bytes"
	"context"
	"image"
	"io"
	"testing"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductMedia(t *testing.T) {
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
				url := "http://example.com/media/image_xxx.png"
				thumbnails := common.Images{
					{URL: url, Size: common.ImageSizeSmall},
					{URL: url, Size: common.ImageSizeMedium},
					{URL: url, Size: common.ImageSizeLarge},
				}
				mocks.storage.EXPECT().
					Download(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, url string) (io.Reader, error) {
						expect := map[string]io.Reader{
							"http://example.com/media/image.png": testImageFile(t),
							"http://example.com/media/video.mp4": testVideoFile(t),
						}
						file, ok := expect[url]
						require.True(t, ok)
						require.NotEmpty(t, file)
						return file, nil
					}).Times(2)
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(url, nil).Times(3)
				mocks.store.EXPECT().
					UpdateProductMedia(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, in *store.UpdateProductMediaInput) error {
						assert.Equal(t, "target-id", in.ProductID)
						assert.Len(t, in.Images, 1)
						assert.Equal(t, "http://example.com/media/image.png", in.Images[0].OriginURL)
						assert.ElementsMatch(t, thumbnails, in.Images[0].Images)
						return nil
					})
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProductMedia,
				URLs:     []string{"http://example.com/media/image.png", "http://example.com/media/video.mp4"},
			},
			expectErr: nil,
		},
		{
			name: "failed to empty url",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProductMedia,
				URLs:     []string{},
			},
			expectErr: errRequiredMediaURL,
		},
		{
			name: "failed to download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProductMedia,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: assert.AnError,
		},
		{
			name: "failed to resize images",
			setup: func(ctx context.Context, mocks *mocks) {
				file := &bytes.Buffer{}
				mocks.storage.EXPECT().Download(gomock.Any(), "http://example.com/media/image.png").Return(file, nil)
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProductMedia,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: image.ErrFormat,
		},
		{
			name: "failed to upload images",
			setup: func(ctx context.Context, mocks *mocks) {
				file := testImageFile(t)
				mocks.storage.EXPECT().Download(gomock.Any(), "http://example.com/media/image.png").Return(file, nil)
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", assert.AnError).AnyTimes()
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProductMedia,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: assert.AnError,
		},
		{
			name: "failed to update product media",
			setup: func(ctx context.Context, mocks *mocks) {
				file := testImageFile(t)
				url := "http://example.com/media/image_xxx.png"
				mocks.storage.EXPECT().Download(gomock.Any(), "http://example.com/media/image.png").Return(file, nil)
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(url, nil).Times(3)
				mocks.store.EXPECT().UpdateProductMedia(ctx, gomock.Any()).Return(assert.AnError)
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProductMedia,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testResizer(tt.setup, func(ctx context.Context, t *testing.T, resizer *resizer) {
			t.Parallel()
			err := resizer.productMedia(ctx, tt.payload)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
