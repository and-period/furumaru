package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	mock_media "github.com/and-period/furumaru/api/mock/media"
	mock_messenger "github.com/and-period/furumaru/api/mock/messenger"
	mock_cognito "github.com/and-period/furumaru/api/mock/pkg/cognito"
	mock_store "github.com/and-period/furumaru/api/mock/store"
	mock_database "github.com/and-period/furumaru/api/mock/user/database"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/jst"
	govalidator "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
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
	Address          *mock_database.MockAddress
	Admin            *mock_database.MockAdmin
	AdminGroup       *mock_database.MockAdminGroup
	AdminGroupRole   *mock_database.MockAdminGroupRole
	AdminGroupUser   *mock_database.MockAdminGroupUser
	AdminPolicy      *mock_database.MockAdminPolicy
	AdminRole        *mock_database.MockAdminRole
	AdminRolePolicy  *mock_database.MockAdminRolePolicy
	Administrator    *mock_database.MockAdministrator
	Coordinator      *mock_database.MockCoordinator
	Guest            *mock_database.MockGuest
	Member           *mock_database.MockMember
	Producer         *mock_database.MockProducer
	User             *mock_database.MockUser
	UserNotification *mock_database.MockUserNotification
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
		Address:          mock_database.NewMockAddress(ctrl),
		Admin:            mock_database.NewMockAdmin(ctrl),
		AdminGroup:       mock_database.NewMockAdminGroup(ctrl),
		AdminGroupRole:   mock_database.NewMockAdminGroupRole(ctrl),
		AdminGroupUser:   mock_database.NewMockAdminGroupUser(ctrl),
		AdminPolicy:      mock_database.NewMockAdminPolicy(ctrl),
		AdminRole:        mock_database.NewMockAdminRole(ctrl),
		AdminRolePolicy:  mock_database.NewMockAdminRolePolicy(ctrl),
		Administrator:    mock_database.NewMockAdministrator(ctrl),
		Coordinator:      mock_database.NewMockCoordinator(ctrl),
		Guest:            mock_database.NewMockGuest(ctrl),
		Member:           mock_database.NewMockMember(ctrl),
		Producer:         mock_database.NewMockProducer(ctrl),
		User:             mock_database.NewMockUser(ctrl),
		UserNotification: mock_database.NewMockUserNotification(ctrl),
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
			Address:          mocks.db.Address,
			Admin:            mocks.db.Admin,
			AdminGroup:       mocks.db.AdminGroup,
			AdminGroupRole:   mocks.db.AdminGroupRole,
			AdminGroupUser:   mocks.db.AdminGroupUser,
			AdminPolicy:      mocks.db.AdminPolicy,
			AdminRole:        mocks.db.AdminRole,
			AdminRolePolicy:  mocks.db.AdminRolePolicy,
			Administrator:    mocks.db.Administrator,
			Coordinator:      mocks.db.Coordinator,
			Guest:            mocks.db.Guest,
			Member:           mocks.db.Member,
			Producer:         mocks.db.Producer,
			User:             mocks.db.User,
			UserNotification: mocks.db.UserNotification,
		},
		AdminAuth: mocks.adminAuth,
		UserAuth:  mocks.userAuth,
		Store:     mocks.store,
		Messenger: mocks.messenger,
		Media:     mocks.media,
		DefaultAdminGroups: map[entity.AdminType][]string{
			entity.AdminTypeAdministrator: {"group-id"},
			entity.AdminTypeCoordinator:   {"group-id"},
			entity.AdminTypeProducer:      {"group-id"},
		},
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

func TestMain(m *testing.M) {
	opts := []goleak.Option{
		goleak.IgnoreTopFunction("go.opencensus.io/stats/view.(*worker).start"),
	}
	goleak.VerifyTestMain(m, opts...)
}

func TestService(t *testing.T) {
	t.Parallel()
	srv := NewService(&Params{}, WithLogger(zap.NewNop()))
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
			name:   "auth invalid argument",
			err:    cognito.ErrInvalidArgument,
			expect: exception.ErrInvalidArgument,
		},
		{
			name:   "auth unauthenticated",
			err:    cognito.ErrUnauthenticated,
			expect: exception.ErrUnauthenticated,
		},
		{
			name:   "auth not found",
			err:    cognito.ErrNotFound,
			expect: exception.ErrNotFound,
		},
		{
			name:   "auth already exists",
			err:    cognito.ErrAlreadyExists,
			expect: exception.ErrAlreadyExists,
		},
		{
			name:   "auth resource exhausted",
			err:    cognito.ErrResourceExhausted,
			expect: exception.ErrResourceExhausted,
		},
		{
			name:   "auth deadline exceeded",
			err:    cognito.ErrTimeout,
			expect: exception.ErrDeadlineExceeded,
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
