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
	db        *dbMocks
	adminAuth *mock_cognito.MockClient
	userAuth  *mock_cognito.MockClient
	messenger *mock_messenger.MockService
}

type dbMocks struct {
	AdminAuth     *mock_database.MockAdminAuth
	Administrator *mock_database.MockAdministrator
	Coordinator   *mock_database.MockCoordinator
	Producer      *mock_database.MockProducer
	User          *mock_database.MockUser
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
		db:        newDBMocks(ctrl),
		adminAuth: mock_cognito.NewMockClient(ctrl),
		userAuth:  mock_cognito.NewMockClient(ctrl),
		messenger: mock_messenger.NewMockService(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		AdminAuth:     mock_database.NewMockAdminAuth(ctrl),
		Administrator: mock_database.NewMockAdministrator(ctrl),
		Coordinator:   mock_database.NewMockCoordinator(ctrl),
		Producer:      mock_database.NewMockProducer(ctrl),
		User:          mock_database.NewMockUser(ctrl),
	}
}

func newService(mocks *mocks, opts ...testOption) *service {
	dopts := &testOptions{
		now: jst.Now,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &service{
		now:         dopts.now,
		logger:      zap.NewNop(),
		waitGroup:   &sync.WaitGroup{},
		sharedGroup: &singleflight.Group{},
		validator:   validator.NewValidator(),
		db: &database.Database{
			AdminAuth:     mocks.db.AdminAuth,
			Administrator: mocks.db.Administrator,
			Coordinator:   mocks.db.Coordinator,
			Producer:      mocks.db.Producer,
			User:          mocks.db.User,
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

		srv := newService(mocks, opts...)
		setup(ctx, mocks)

		testFunc(ctx, t, srv)
		srv.waitGroup.Wait()
	}
}

func TestService(t *testing.T) {
	t.Parallel()
	srv := NewService(&Params{}, WithLogger(zap.NewNop()))
	assert.NotNil(t, srv)
}
