package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media/database"
	mock_database "github.com/and-period/furumaru/api/mock/media/database"
	mock_dynamodb "github.com/and-period/furumaru/api/mock/pkg/dynamodb"
	mock_medialive "github.com/and-period/furumaru/api/mock/pkg/medialive"
	mock_sqs "github.com/and-period/furumaru/api/mock/pkg/sqs"
	mock_storage "github.com/and-period/furumaru/api/mock/pkg/storage"
	mock_store "github.com/and-period/furumaru/api/mock/store"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/storage"
	govalidator "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

var (
	tmpURL     = "http://tmp.and-period.jp"
	storageURL = "http://and-period.jp"
	unknownURL = "http://example.com"
)

type mocks struct {
	db       *dbMocks
	cache    *mock_dynamodb.MockClient
	store    *mock_store.MockService
	tmp      *mock_storage.MockBucket
	storage  *mock_storage.MockBucket
	producer *mock_sqs.MockProducer
	media    *mock_medialive.MockMediaLive
}

type dbMocks struct {
	Broadcast *mock_database.MockBroadcast
}

type testOptions struct {
	now func() time.Time
}

type testOption func(opts *testOptions)

func withNow(now time.Time) testOption {
	return func(opts *testOptions) {
		opts.now = func() time.Time {
			return now
		}
	}
}

type testCaller func(ctx context.Context, t *testing.T, service *service)

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		db:       newDBMocks(ctrl),
		cache:    mock_dynamodb.NewMockClient(ctrl),
		store:    mock_store.NewMockService(ctrl),
		tmp:      mock_storage.NewMockBucket(ctrl),
		storage:  mock_storage.NewMockBucket(ctrl),
		producer: mock_sqs.NewMockProducer(ctrl),
		media:    mock_medialive.NewMockMediaLive(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		Broadcast: mock_database.NewMockBroadcast(ctrl),
	}
}

func newService(mocks *mocks, opts ...testOption) *service {
	dopts := &testOptions{
		now: jst.Now,
	}
	for i := range opts {
		opts[i](dopts)
	}
	params := &Params{
		WaitGroup: &sync.WaitGroup{},
		Database: &database.Database{
			Broadcast: mocks.db.Broadcast,
		},
		Cache:     mocks.cache,
		Store:     mocks.store,
		Tmp:       mocks.tmp,
		Storage:   mocks.storage,
		Producer:  mocks.producer,
		MediaLive: mocks.media,
	}
	tmpHost, _ := url.Parse(tmpURL)
	storageHost, _ := url.Parse(storageURL)
	mocks.tmp.EXPECT().GetHost().Return(tmpHost, nil)
	mocks.storage.EXPECT().GetHost().Return(storageHost, nil)
	srv, _ := NewService(params)
	service := srv.(*service)
	service.now = func() time.Time {
		return dopts.now()
	}
	return service
}

func testService(
	setup func(ctx context.Context, mocks *mocks),
	testFunc testCaller,
	opts ...testOption,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mocks := newMocks(ctrl)

		srv := newService(mocks, opts...)
		setup(ctx, mocks)

		testFunc(ctx, t, srv)
		srv.waitGroup.Wait()
	}
}

func testImageFile(t *testing.T) (io.Reader, *multipart.FileHeader) {
	const filename, format = "and-period.png", "image"

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	defer writer.Close()

	filepath := testFilePath(t, filename)
	file, err := os.Open(filepath)
	require.NoError(t, err, err)

	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, format, filename))
	header.Set("Content-Type", "multipart/form-data")
	part := &multipart.FileHeader{
		Filename: filepath,
		Header:   header,
		Size:     3 << 20, // 3MB
	}

	return file, part
}

func testVideoFile(t *testing.T) (io.Reader, *multipart.FileHeader) {
	const filename, format = "and-period.mp4", "video"

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	defer writer.Close()

	filepath := testFilePath(t, filename)
	file, err := os.Open(filepath)
	require.NoError(t, err, err)

	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, format, filename))
	header.Set("Content-Type", "multipart/form-data")
	part := &multipart.FileHeader{
		Filename: filepath,
		Header:   header,
		Size:     10 << 20, // 10MB
	}

	return file, part
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

func TestService(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mocks := newMocks(ctrl)
	params := &Params{Storage: mocks.storage, Tmp: mocks.tmp}
	surl, err := url.Parse(storageURL)
	require.NoError(t, err)
	turl, err := url.Parse(tmpURL)
	require.NoError(t, err)
	mocks.storage.EXPECT().GetHost().Return(surl, nil)
	mocks.tmp.EXPECT().GetHost().Return(turl, nil)
	srv, err := NewService(params, WithLogger(zap.NewNop()))
	assert.NoError(t, err)
	assert.NotNil(t, srv)
}

func TestInternalError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "not error",
			err:    nil,
			expect: nil,
		},
		{
			name:   "validation error",
			err:    govalidator.ValidationErrors{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name:   "database not found",
			err:    database.ErrNotFound,
			expect: exception.ErrNotFound,
		},
		{
			name:   "database failed precondition",
			err:    database.ErrFailedPrecondition,
			expect: exception.ErrFailedPrecondition,
		},
		{
			name:   "database already exists",
			err:    database.ErrAlreadyExists,
			expect: exception.ErrAlreadyExists,
		},
		{
			name:   "database deadline exceeded",
			err:    database.ErrDeadlineExceeded,
			expect: exception.ErrDeadlineExceeded,
		},
		{
			name:   "storage invalid argument",
			err:    storage.ErrInvalidURL,
			expect: exception.ErrInvalidArgument,
		},
		{
			name:   "storage not found",
			err:    storage.ErrNotFound,
			expect: exception.ErrNotFound,
		},
		{
			name:   "context canceled",
			err:    context.Canceled,
			expect: exception.ErrCanceled,
		},
		{
			name:   "context deadline exceeded",
			err:    context.DeadlineExceeded,
			expect: exception.ErrDeadlineExceeded,
		},
		{
			name:   "other error",
			err:    assert.AnError,
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := internalError(tt.err)
			assert.ErrorIs(t, actual, tt.expect)
		})
	}
}
