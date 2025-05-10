package service

import (
	"context"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media/database"
	mock_database "github.com/and-period/furumaru/api/mock/media/database"
	mock_batch "github.com/and-period/furumaru/api/mock/pkg/batch"
	mock_dynamodb "github.com/and-period/furumaru/api/mock/pkg/dynamodb"
	mock_medialive "github.com/and-period/furumaru/api/mock/pkg/medialive"
	mock_sqs "github.com/and-period/furumaru/api/mock/pkg/sqs"
	mock_storage "github.com/and-period/furumaru/api/mock/pkg/storage"
	mock_youtube "github.com/and-period/furumaru/api/mock/pkg/youtube"
	mock_store "github.com/and-period/furumaru/api/mock/store"
	mock_user "github.com/and-period/furumaru/api/mock/user"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"github.com/and-period/furumaru/api/pkg/youtube"
	govalidator "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

var (
	tmpURL     = "http://tmp.and-period.jp"
	storageURL = "http://and-period.jp"
)

type mocks struct {
	db             *dbMocks
	cache          *mock_dynamodb.MockClient
	store          *mock_store.MockService
	tmp            *mock_storage.MockBucket
	storage        *mock_storage.MockBucket
	producer       *mock_sqs.MockProducer
	batch          *mock_batch.MockClient
	user           *mock_user.MockService
	media          *mock_medialive.MockMediaLive
	youtube        *mock_youtube.MockYoutube
	youtubeService *mock_youtube.MockService
	youtubeAuth    *mock_youtube.MockAuth
}

type dbMocks struct {
	Broadcast          *mock_database.MockBroadcast
	BroadcastComment   *mock_database.MockBroadcastComment
	BroadcastViewerLog *mock_database.MockBroadcastViewerLog
	Video              *mock_database.MockVideo
	VideoComment       *mock_database.MockVideoComment
	VideoViewerLog     *mock_database.MockVideoViewerLog
}

type testOptions struct {
	now  func() time.Time
	uuid func() string
}

type testOption func(opts *testOptions)

func withNow(now time.Time) testOption {
	return func(opts *testOptions) {
		opts.now = func() time.Time {
			return now
		}
	}
}

func withUUID(uuid string) testOption {
	return func(opts *testOptions) {
		opts.uuid = func() string {
			return uuid
		}
	}
}

type testCaller func(ctx context.Context, t *testing.T, service *service)

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		db:             newDBMocks(ctrl),
		cache:          mock_dynamodb.NewMockClient(ctrl),
		store:          mock_store.NewMockService(ctrl),
		tmp:            mock_storage.NewMockBucket(ctrl),
		storage:        mock_storage.NewMockBucket(ctrl),
		producer:       mock_sqs.NewMockProducer(ctrl),
		batch:          mock_batch.NewMockClient(ctrl),
		user:           mock_user.NewMockService(ctrl),
		media:          mock_medialive.NewMockMediaLive(ctrl),
		youtube:        mock_youtube.NewMockYoutube(ctrl),
		youtubeService: mock_youtube.NewMockService(ctrl),
		youtubeAuth:    mock_youtube.NewMockAuth(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		Broadcast:          mock_database.NewMockBroadcast(ctrl),
		BroadcastComment:   mock_database.NewMockBroadcastComment(ctrl),
		BroadcastViewerLog: mock_database.NewMockBroadcastViewerLog(ctrl),
		Video:              mock_database.NewMockVideo(ctrl),
		VideoComment:       mock_database.NewMockVideoComment(ctrl),
		VideoViewerLog:     mock_database.NewMockVideoViewerLog(ctrl),
	}
}

func newService(mocks *mocks, opts ...testOption) *service {
	dopts := &testOptions{
		now:  jst.Now,
		uuid: uuid.New,
	}
	for i := range opts {
		opts[i](dopts)
	}
	params := &Params{
		WaitGroup: &sync.WaitGroup{},
		Database: &database.Database{
			Broadcast:          mocks.db.Broadcast,
			BroadcastComment:   mocks.db.BroadcastComment,
			BroadcastViewerLog: mocks.db.BroadcastViewerLog,
			Video:              mocks.db.Video,
			VideoComment:       mocks.db.VideoComment,
			VideoViewerLog:     mocks.db.VideoViewerLog,
		},
		Cache:                        mocks.cache,
		User:                         mocks.user,
		Store:                        mocks.store,
		Tmp:                          mocks.tmp,
		Storage:                      mocks.storage,
		Producer:                     mocks.producer,
		Batch:                        mocks.batch,
		MediaLive:                    mocks.media,
		Youtube:                      mocks.youtube,
		BatchUpdateArchiveDefinition: "batch-update-archive-definition",
		BatchUpdateArchiveQueue:      "batch-update-archive-queue",
		BatchUpdateArchiveCommand: func(broadcastID string) []string {
			return []string{"batch-update-archive-command", broadcastID}
		},
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
	service.generateID = func() string {
		return dopts.uuid()
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
		ctx, cancel := context.WithCancel(t.Context())
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

func TestMain(m *testing.M) {
	opts := []goleak.Option{
		goleak.IgnoreTopFunction("go.opencensus.io/stats/view.(*worker).start"),
	}
	goleak.VerifyTestMain(m, opts...)
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
			name:   "youtube invalid argument",
			err:    youtube.ErrBadRequest,
			expect: exception.ErrInvalidArgument,
		},
		{
			name:   "youtube unauthorized",
			err:    youtube.ErrUnauthorized,
			expect: exception.ErrUnauthenticated,
		},
		{
			name:   "youtube forbidden",
			err:    youtube.ErrForbidden,
			expect: exception.ErrForbidden,
		},
		{
			name:   "youtube too many requests",
			err:    youtube.ErrTooManyRequests,
			expect: exception.ErrResourceExhausted,
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := internalError(tt.err)
			assert.ErrorIs(t, actual, tt.expect)
		})
	}
}
