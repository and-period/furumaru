package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	mock_messenger "github.com/and-period/furumaru/api/mock/messenger"
	mock_cognito "github.com/and-period/furumaru/api/mock/pkg/cognito"
	mock_storage "github.com/and-period/furumaru/api/mock/pkg/storage"
	mock_database "github.com/and-period/furumaru/api/mock/user/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/validator"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var errmock = errors.New("some error")

type mocks struct {
	storage   *mock_storage.MockBucket
	db        *dbMocks
	adminAuth *mock_cognito.MockClient
	userAuth  *mock_cognito.MockClient
	messenger *mock_messenger.MockMessengerService
}

type dbMocks struct {
	Admin *mock_database.MockAdmin
	User  *mock_database.MockUser
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

type testCaller func(ctx context.Context, t *testing.T, service *userService)

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		storage:   mock_storage.NewMockBucket(ctrl),
		db:        newDBMocks(ctrl),
		adminAuth: mock_cognito.NewMockClient(ctrl),
		userAuth:  mock_cognito.NewMockClient(ctrl),
		messenger: mock_messenger.NewMockMessengerService(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		Admin: mock_database.NewMockAdmin(ctrl),
		User:  mock_database.NewMockUser(ctrl),
	}
}

func newUserService(mocks *mocks, opts ...testOption) *userService {
	dopts := &testOptions{
		now: jst.Now,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &userService{
		now:         dopts.now,
		logger:      zap.NewNop(),
		sharedGroup: &singleflight.Group{},
		waitGroup:   &sync.WaitGroup{},
		validator:   validator.NewValidator(),
		storage:     mocks.storage,
		db: &database.Database{
			Admin: mocks.db.Admin,
			User:  mocks.db.User,
		},
		adminAuth: mocks.adminAuth,
		userAuth:  mocks.userAuth,
		messenger: mocks.messenger,
	}
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

		srv := newUserService(mocks)
		setup(ctx, mocks)

		testFunc(ctx, t, srv)
		srv.waitGroup.Wait()
	}
}

func TestUserService(t *testing.T) {
	t.Parallel()
	srv := NewUserService(&Params{}, WithLogger(zap.NewNop()))
	assert.NotNil(t, srv)
}
