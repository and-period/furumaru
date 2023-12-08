package service

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadCoordinatorThumbnail(t *testing.T) {
	t.Parallel()
	file, _ := testImageFile(t)
	path := strings.Join([]string{entity.CoordinatorThumbnailPath, "calmato.png"}, "/")
	turl, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
	require.NoError(t, err)
	tpath := strings.TrimPrefix(turl.Path, "/")
	surl, err := url.Parse(strings.Join([]string{storageURL, path}, "/"))
	require.NoError(t, err)
	uurl, err := url.Parse(strings.Join([]string{unknownURL, path}, "/"))
	require.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UploadFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success to upload permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return(surl.String(), nil)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name: "success to download permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(file, nil)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UploadFileInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid url",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: unknownURL,
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "unknown storage host",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: uurl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to download temporary file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upload permanent file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to download permanent file when download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.UploadCoordinatorThumbnail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUploadCoordinatorHeader(t *testing.T) {
	t.Parallel()
	file, _ := testImageFile(t)
	path := strings.Join([]string{entity.CoordinatorHeaderPath, "calmato.png"}, "/")
	turl, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
	require.NoError(t, err)
	tpath := strings.TrimPrefix(turl.Path, "/")
	surl, err := url.Parse(strings.Join([]string{storageURL, path}, "/"))
	require.NoError(t, err)
	uurl, err := url.Parse(strings.Join([]string{unknownURL, path}, "/"))
	require.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UploadFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success to upload permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return(surl.String(), nil)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name: "success to download permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(file, nil)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UploadFileInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid url",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: unknownURL,
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "unknown storage host",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: uurl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to download temporary file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upload permanent file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to download permanent file when download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.UploadCoordinatorHeader(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUploadCoordinatorPromotionVideo(t *testing.T) {
	t.Parallel()
	file, _ := testImageFile(t)
	path := strings.Join([]string{entity.CoordinatorPromotionVideoPath, "calmato.mp4"}, "/")
	turl, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
	require.NoError(t, err)
	tpath := strings.TrimPrefix(turl.Path, "/")
	surl, err := url.Parse(strings.Join([]string{storageURL, path}, "/"))
	require.NoError(t, err)
	uurl, err := url.Parse(strings.Join([]string{unknownURL, path}, "/"))
	require.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UploadFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success to upload permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return(surl.String(), nil)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name: "success to download permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(file, nil)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UploadFileInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid url",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: unknownURL,
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "unknown storage host",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: uurl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to download temporary file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upload permanent file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to download permanent file when download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.UploadCoordinatorPromotionVideo(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUploadCoordinatorBonusVideo(t *testing.T) {
	t.Parallel()
	file, _ := testImageFile(t)
	path := strings.Join([]string{entity.CoordinatorBonusVideoPath, "calmato.mp4"}, "/")
	turl, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
	require.NoError(t, err)
	tpath := strings.TrimPrefix(turl.Path, "/")
	surl, err := url.Parse(strings.Join([]string{storageURL, path}, "/"))
	require.NoError(t, err)
	uurl, err := url.Parse(strings.Join([]string{unknownURL, path}, "/"))
	require.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UploadFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success to upload permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return(surl.String(), nil)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name: "success to download permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(file, nil)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UploadFileInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid url",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: unknownURL,
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "unknown storage host",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: uurl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to download temporary file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upload permanent file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to download permanent file when download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.UploadCoordinatorBonusVideo(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUploadProducerThumbnail(t *testing.T) {
	t.Parallel()
	file, _ := testImageFile(t)
	path := strings.Join([]string{entity.ProducerThumbnailPath, "calmato.png"}, "/")
	turl, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
	require.NoError(t, err)
	tpath := strings.TrimPrefix(turl.Path, "/")
	surl, err := url.Parse(strings.Join([]string{storageURL, path}, "/"))
	require.NoError(t, err)
	uurl, err := url.Parse(strings.Join([]string{unknownURL, path}, "/"))
	require.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UploadFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success to upload permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return(surl.String(), nil)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name: "success to download permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(file, nil)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UploadFileInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid url",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: unknownURL,
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "unknown storage host",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: uurl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to download temporary file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upload permanent file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to download permanent file when download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.UploadProducerThumbnail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUploadProducerHeader(t *testing.T) {
	t.Parallel()
	file, _ := testImageFile(t)
	path := strings.Join([]string{entity.ProducerHeaderPath, "calmato.png"}, "/")
	turl, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
	require.NoError(t, err)
	tpath := strings.TrimPrefix(turl.Path, "/")
	surl, err := url.Parse(strings.Join([]string{storageURL, path}, "/"))
	require.NoError(t, err)
	uurl, err := url.Parse(strings.Join([]string{unknownURL, path}, "/"))
	require.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UploadFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success to upload permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return(surl.String(), nil)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name: "success to download permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(file, nil)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UploadFileInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid url",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: unknownURL,
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "unknown storage host",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: uurl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to download temporary file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upload permanent file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to download permanent file when download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.UploadProducerHeader(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUploadProducerPromotionVideo(t *testing.T) {
	t.Parallel()
	file, _ := testImageFile(t)
	path := strings.Join([]string{entity.ProducerPromotionVideoPath, "calmato.mp4"}, "/")
	turl, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
	require.NoError(t, err)
	tpath := strings.TrimPrefix(turl.Path, "/")
	surl, err := url.Parse(strings.Join([]string{storageURL, path}, "/"))
	require.NoError(t, err)
	uurl, err := url.Parse(strings.Join([]string{unknownURL, path}, "/"))
	require.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UploadFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success to upload permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return(surl.String(), nil)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name: "success to download permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(file, nil)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UploadFileInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid url",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: unknownURL,
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "unknown storage host",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: uurl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to download temporary file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upload permanent file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to download permanent file when download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.UploadProducerPromotionVideo(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUploadProducerBonusVideo(t *testing.T) {
	t.Parallel()
	file, _ := testImageFile(t)
	path := strings.Join([]string{entity.ProducerBonusVideoPath, "calmato.mp4"}, "/")
	turl, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
	require.NoError(t, err)
	tpath := strings.TrimPrefix(turl.Path, "/")
	surl, err := url.Parse(strings.Join([]string{storageURL, path}, "/"))
	require.NoError(t, err)
	uurl, err := url.Parse(strings.Join([]string{unknownURL, path}, "/"))
	require.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UploadFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success to upload permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return(surl.String(), nil)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name: "success to download permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(file, nil)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UploadFileInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid url",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: unknownURL,
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "unknown storage host",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: uurl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to download temporary file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upload permanent file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to download permanent file when download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.UploadProducerBonusVideo(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUploadUserThumbnail(t *testing.T) {
	t.Parallel()
	file, _ := testImageFile(t)
	path := strings.Join([]string{entity.UserThumbnailPath, "calmato.png"}, "/")
	turl, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
	require.NoError(t, err)
	tpath := strings.TrimPrefix(turl.Path, "/")
	surl, err := url.Parse(strings.Join([]string{storageURL, path}, "/"))
	require.NoError(t, err)
	uurl, err := url.Parse(strings.Join([]string{unknownURL, path}, "/"))
	require.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UploadFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success to upload permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return(surl.String(), nil)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name: "success to download permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(file, nil)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UploadFileInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid url",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: unknownURL,
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "unknown storage host",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: uurl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to download temporary file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upload permanent file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to download permanent file when download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.UploadUserThumbnail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUploadProductMedia(t *testing.T) {
	t.Parallel()
	// 画像関連
	image, _ := testImageFile(t)
	ipath := strings.Join([]string{entity.ProductMediaImagePath, "calmato.png"}, "/")
	iturl, err := url.Parse(strings.Join([]string{tmpURL, ipath}, "/"))
	require.NoError(t, err)
	itpath := strings.TrimPrefix(iturl.Path, "/")
	isurl, err := url.Parse(strings.Join([]string{storageURL, ipath}, "/"))
	require.NoError(t, err)
	// 映像関連
	video, _ := testVideoFile(t)
	vpath := strings.Join([]string{entity.ProductMediaVideoPath, "calmato.png"}, "/")
	vturl, err := url.Parse(strings.Join([]string{tmpURL, vpath}, "/"))
	require.NoError(t, err)
	vtpath := strings.TrimPrefix(vturl.Path, "/")
	vsurl, err := url.Parse(strings.Join([]string{storageURL, vpath}, "/"))
	require.NoError(t, err)
	// その他
	uurl, err := url.Parse(strings.Join([]string{unknownURL, ipath}, "/"))
	require.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UploadFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success to upload permanent image file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, iturl.String()).Return(image, nil)
				mocks.storage.EXPECT().Upload(ctx, itpath, gomock.Any()).Return(isurl.String(), nil)
			},
			input: &media.UploadFileInput{
				URL: iturl.String(),
			},
			expect:    isurl.String(),
			expectErr: nil,
		},
		{
			name: "success to download permanent image file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, isurl.String()).Return(image, nil)
			},
			input: &media.UploadFileInput{
				URL: isurl.String(),
			},
			expect:    isurl.String(),
			expectErr: nil,
		},
		{
			name: "success to upload permanent video file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, vturl.String()).Return(video, nil)
				mocks.storage.EXPECT().Upload(ctx, vtpath, gomock.Any()).Return(vsurl.String(), nil)
			},
			input: &media.UploadFileInput{
				URL: vturl.String(),
			},
			expect:    vsurl.String(),
			expectErr: nil,
		},
		{
			name: "success to download permanent video file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, vsurl.String()).Return(video, nil)
			},
			input: &media.UploadFileInput{
				URL: vsurl.String(),
			},
			expect:    vsurl.String(),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UploadFileInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid url",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: unknownURL,
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "unknown storage host",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: uurl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to download temporary file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, iturl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: iturl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upload permanent file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, iturl.String()).Return(image, nil)
				mocks.storage.EXPECT().Upload(ctx, itpath, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: iturl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to download permanent file when download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, isurl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: isurl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.UploadProductMedia(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUploadProductTypeIcon(t *testing.T) {
	t.Parallel()
	file, _ := testImageFile(t)
	path := strings.Join([]string{entity.ProductTypeIconPath, "calmato.png"}, "/")
	turl, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
	require.NoError(t, err)
	tpath := strings.TrimPrefix(turl.Path, "/")
	surl, err := url.Parse(strings.Join([]string{storageURL, path}, "/"))
	require.NoError(t, err)
	uurl, err := url.Parse(strings.Join([]string{unknownURL, path}, "/"))
	require.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UploadFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success to upload permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return(surl.String(), nil)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name: "success to download permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(file, nil)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UploadFileInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid url",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: unknownURL,
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "unknown storage host",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: uurl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to download temporary file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upload permanent file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to download permanent file when download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.UploadProductTypeIcon(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUploadScheduleThumbnail(t *testing.T) {
	t.Parallel()
	file, _ := testImageFile(t)
	path := strings.Join([]string{entity.ScheduleThumbnailPath, "calmato.png"}, "/")
	turl, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
	require.NoError(t, err)
	tpath := strings.TrimPrefix(turl.Path, "/")
	surl, err := url.Parse(strings.Join([]string{storageURL, path}, "/"))
	require.NoError(t, err)
	uurl, err := url.Parse(strings.Join([]string{unknownURL, path}, "/"))
	require.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UploadFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success to upload permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return(surl.String(), nil)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name: "success to download permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(file, nil)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UploadFileInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid url",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: unknownURL,
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "unknown storage host",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: uurl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to download temporary file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upload permanent file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to download permanent file when download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.UploadScheduleThumbnail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUploadScheduleImage(t *testing.T) {
	t.Parallel()
	file, _ := testImageFile(t)
	path := strings.Join([]string{entity.ScheduleImagePath, "calmato.png"}, "/")
	turl, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
	require.NoError(t, err)
	tpath := strings.TrimPrefix(turl.Path, "/")
	surl, err := url.Parse(strings.Join([]string{storageURL, path}, "/"))
	require.NoError(t, err)
	uurl, err := url.Parse(strings.Join([]string{unknownURL, path}, "/"))
	require.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UploadFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success to upload permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return(surl.String(), nil)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name: "success to download permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(file, nil)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UploadFileInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid url",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: unknownURL,
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "unknown storage host",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: uurl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to download temporary file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upload permanent file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to download permanent file when download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.UploadScheduleImage(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUploadScheduleOpeningVideo(t *testing.T) {
	t.Parallel()
	file, _ := testImageFile(t)
	path := strings.Join([]string{entity.ScheduleOpeningVideoPath, "calmato.mp4"}, "/")
	turl, err := url.Parse(strings.Join([]string{tmpURL, path}, "/"))
	require.NoError(t, err)
	tpath := strings.TrimPrefix(turl.Path, "/")
	surl, err := url.Parse(strings.Join([]string{storageURL, path}, "/"))
	require.NoError(t, err)
	uurl, err := url.Parse(strings.Join([]string{unknownURL, path}, "/"))
	require.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UploadFileInput
		expect    string
		expectErr error
	}{
		{
			name: "success to upload permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return(surl.String(), nil)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name: "success to download permanent file",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(file, nil)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    surl.String(),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UploadFileInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid url",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: unknownURL,
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "unknown storage host",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &media.UploadFileInput{
				URL: uurl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to download temporary file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upload permanent file when upload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.tmp.EXPECT().Download(ctx, turl.String()).Return(file, nil)
				mocks.storage.EXPECT().Upload(ctx, tpath, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: turl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to download permanent file when download",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, surl.String()).Return(nil, assert.AnError)
			},
			input: &media.UploadFileInput{
				URL: surl.String(),
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.UploadScheduleOpeningVideo(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}
