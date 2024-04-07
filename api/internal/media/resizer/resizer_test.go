package resizer

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/and-period/furumaru/api/internal/media/entity"
	mock_storage "github.com/and-period/furumaru/api/mock/pkg/storage"
	mock_store "github.com/and-period/furumaru/api/mock/store"
	mock_user "github.com/and-period/furumaru/api/mock/user"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type mocks struct {
	storage *mock_storage.MockBucket
	user    *mock_user.MockService
	store   *mock_store.MockService
}

type testCaller func(ctx context.Context, t *testing.T, resizer *resizer)

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		storage: mock_storage.NewMockBucket(ctrl),
		user:    mock_user.NewMockService(ctrl),
		store:   mock_store.NewMockService(ctrl),
	}
}

func newResizer(mocks *mocks) *resizer {
	params := &Params{
		WaitGroup: &sync.WaitGroup{},
		Storage:   mocks.storage,
		User:      mocks.user,
		Store:     mocks.store,
	}
	resizer := NewResizer(params).(*resizer)
	resizer.concurrency = 1
	return resizer
}

func testResizer(
	setup func(ctx context.Context, mocks *mocks),
	testFunc testCaller,
) func(t *testing.T) {
	return func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mocks := newMocks(ctrl)

		r := newResizer(mocks)
		setup(ctx, mocks)

		testFunc(ctx, t, r)
		r.waitGroup.Wait()
	}
}

func testImageFile(t *testing.T) io.Reader {
	const filename = "and-period.png"

	filepath := testFilePath(t, filename)
	file, err := os.Open(filepath)
	require.NoError(t, err, err)

	return file
}

func testVideoFile(t *testing.T) io.Reader {
	const filename, format = "and-period.mp4", "video"

	filepath := testFilePath(t, filename)
	file, err := os.Open(filepath)
	require.NoError(t, err, err)

	return file
}

func testFilePath(t *testing.T, filename string) string {
	dir, err := os.Getwd()
	assert.NoError(t, err)

	strs := strings.Split(dir, "api/internal")
	if len(strs) == 0 {
		t.Fatal("test: invalid file path")
	}
	return filepath.Join(strs[0], "/api/tmp", filename)
}

func TestResizer(t *testing.T) {
	t.Parallel()
	w := NewResizer(&Params{},
		WithLogger(zap.NewNop()),
		WithConcurrency(1),
		WithConcurrency(3),
		WithMaxRetires(3),
	)
	assert.NotNil(t, w)
}

func TestResizer_Dispatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		record    events.SQSMessage
		expectErr error
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			record: events.SQSMessage{
				Body: `{"id":"", "fileType":0, "urls":[]}`,
			},
			expectErr: nil,
		},
		{
			name:      "failed to unmarshall sqs event",
			setup:     func(ctx context.Context, mocks *mocks) {},
			record:    events.SQSMessage{},
			expectErr: nil,
		},
		{
			name: "failed to run with retry",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, gomock.Any()).Return(nil, context.Canceled)
			},
			record: events.SQSMessage{
				Body: `{"id":"", "fileType":1, "urls":["http://example.com/media/image.png"]}`,
			},
			expectErr: context.Canceled,
		},
		{
			name: "failed to run without retry",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.storage.EXPECT().Download(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			record: events.SQSMessage{
				Body: `{"id":"", "fileType":1, "urls":["http://example.com/media/image.png"]}`,
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testResizer(tt.setup, func(ctx context.Context, t *testing.T, resizer *resizer) {
			t.Parallel()
			err := resizer.dispatch(ctx, tt.record)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestResizer_Run(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		payload   *entity.ResizerPayload
		expectErr error
	}{
		{
			name: "received coordinator thumbnail event",
			setup: func(ctx context.Context, mocks *mocks) {
				file := testImageFile(t)
				mocks.storage.EXPECT().Download(ctx, gomock.Any()).Return(file, nil)
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()
				mocks.user.EXPECT().UpdateCoordinatorThumbnails(ctx, gomock.Any()).Return(nil)
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeCoordinatorThumbnail,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: nil,
		},
		{
			name: "received coordinator header event",
			setup: func(ctx context.Context, mocks *mocks) {
				file := testImageFile(t)
				mocks.storage.EXPECT().Download(ctx, gomock.Any()).Return(file, nil)
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()
				mocks.user.EXPECT().UpdateCoordinatorHeaders(ctx, gomock.Any()).Return(nil)
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeCoordinatorHeader,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: nil,
		},
		{
			name: "received producer thumbnail event",
			setup: func(ctx context.Context, mocks *mocks) {
				file := testImageFile(t)
				mocks.storage.EXPECT().Download(ctx, gomock.Any()).Return(file, nil)
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()
				mocks.user.EXPECT().UpdateProducerThumbnails(ctx, gomock.Any()).Return(nil)
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProducerThumbnail,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: nil,
		},
		{
			name: "received producer header event",
			setup: func(ctx context.Context, mocks *mocks) {
				file := testImageFile(t)
				mocks.storage.EXPECT().Download(ctx, gomock.Any()).Return(file, nil)
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()
				mocks.user.EXPECT().UpdateProducerHeaders(ctx, gomock.Any()).Return(nil)
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProducerHeader,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: nil,
		},
		{
			name: "received product media event",
			setup: func(ctx context.Context, mocks *mocks) {
				file := testImageFile(t)
				mocks.storage.EXPECT().Download(gomock.Any(), gomock.Any()).Return(file, nil)
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()
				mocks.store.EXPECT().UpdateProductMedia(ctx, gomock.Any()).Return(nil)
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProductMedia,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: nil,
		},
		{
			name: "received product type icon event",
			setup: func(ctx context.Context, mocks *mocks) {
				file := testImageFile(t)
				mocks.storage.EXPECT().Download(ctx, gomock.Any()).Return(file, nil)
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()
				mocks.store.EXPECT().UpdateProductTypeIcons(ctx, gomock.Any()).Return(nil)
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeProductTypeIcon,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: nil,
		},
		{
			name: "received schedume thumbnail event",
			setup: func(ctx context.Context, mocks *mocks) {
				file := testImageFile(t)
				mocks.storage.EXPECT().Download(ctx, gomock.Any()).Return(file, nil)
				mocks.storage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()
				mocks.store.EXPECT().UpdateScheduleThumbnails(ctx, gomock.Any()).Return(nil)
			},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeScheduleThumbnail,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: nil,
		},
		{
			name:  "failed to empty urls",
			setup: func(ctx context.Context, mocks *mocks) {},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeCoordinatorThumbnail,
				URLs:     []string{},
			},
			expectErr: errInvalidFormat,
		},
		{
			name:  "failed to invalid file type",
			setup: func(ctx context.Context, mocks *mocks) {},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeUnknown,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			expectErr: errUnknownFileType,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testResizer(tt.setup, func(ctx context.Context, t *testing.T, resizer *resizer) {
			t.Parallel()
			err := resizer.run(ctx, tt.payload)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestResizer_Notify(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		payload   *entity.ResizerPayload
		fn        func() error
		expectErr error
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeCoordinatorThumbnail,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			fn:        func() error { return nil },
			expectErr: nil,
		},
		{
			name:  "failed to function",
			setup: func(ctx context.Context, mocks *mocks) {},
			payload: &entity.ResizerPayload{
				TargetID: "target-id",
				FileType: entity.FileTypeCoordinatorThumbnail,
				URLs:     []string{"http://example.com/media/image.png"},
			},
			fn:        func() error { return assert.AnError },
			expectErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testResizer(tt.setup, func(ctx context.Context, t *testing.T, resizer *resizer) {
			t.Parallel()
			err := resizer.notify(ctx, tt.payload, tt.fn)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
