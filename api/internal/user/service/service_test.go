package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	mock_media "github.com/and-period/furumaru/api/mock/media"
	mock_messenger "github.com/and-period/furumaru/api/mock/messenger"
	mock_cognito "github.com/and-period/furumaru/api/mock/pkg/cognito"
	mock_store "github.com/and-period/furumaru/api/mock/store"
	mock_database "github.com/and-period/furumaru/api/mock/user/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type mocks struct {
	db        *dbMocks
	adminAuth *mock_cognito.MockClient
	userAuth  *mock_cognito.MockClient
	store     *mock_store.MockService
	messenger *mock_messenger.MockService
	media     *mock_media.MockService
}

type dbMocks struct {
	Admin         *mock_database.MockAdmin
	Administrator *mock_database.MockAdministrator
	Coordinator   *mock_database.MockCoordinator
	Member        *mock_database.MockMember
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
		store:     mock_store.NewMockService(ctrl),
		messenger: mock_messenger.NewMockService(ctrl),
		media:     mock_media.NewMockService(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		Admin:         mock_database.NewMockAdmin(ctrl),
		Administrator: mock_database.NewMockAdministrator(ctrl),
		Coordinator:   mock_database.NewMockCoordinator(ctrl),
		Member:        mock_database.NewMockMember(ctrl),
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
	params := &Params{
		WaitGroup: &sync.WaitGroup{},
		Database: &database.Database{
			Admin:         mocks.db.Admin,
			Administrator: mocks.db.Administrator,
			Coordinator:   mocks.db.Coordinator,
			Member:        mocks.db.Member,
			Producer:      mocks.db.Producer,
			User:          mocks.db.User,
		},
		AdminAuth: mocks.adminAuth,
		UserAuth:  mocks.userAuth,
		Store:     mocks.store,
		Messenger: mocks.messenger,
		Media:     mocks.media,
	}
	service := NewService(params).(*service)
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

func TestService(t *testing.T) {
	t.Parallel()
	srv := NewService(&Params{}, WithLogger(zap.NewNop()))
	assert.NotNil(t, srv)
}
