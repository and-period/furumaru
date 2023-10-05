package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	mock_media "github.com/and-period/furumaru/api/mock/media"
	mock_messenger "github.com/and-period/furumaru/api/mock/messenger"
	mock_ivs "github.com/and-period/furumaru/api/mock/pkg/ivs"
	mock_postalcode "github.com/and-period/furumaru/api/mock/pkg/postalcode"
	mock_database "github.com/and-period/furumaru/api/mock/store/database"
	mock_user "github.com/and-period/furumaru/api/mock/user"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/postalcode"
	govalidator "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type mocks struct {
	db         *dbMocks
	user       *mock_user.MockService
	messenger  *mock_messenger.MockService
	media      *mock_media.MockService
	postalCode *mock_postalcode.MockClient
	ivs        *mock_ivs.MockClient
}

type dbMocks struct {
	Address     *mock_database.MockAddress
	Category    *mock_database.MockCategory
	Order       *mock_database.MockOrder
	Product     *mock_database.MockProduct
	ProductTag  *mock_database.MockProductTag
	ProductType *mock_database.MockProductType
	Promotion   *mock_database.MockPromotion
	Shipping    *mock_database.MockShipping
	Schedule    *mock_database.MockSchedule
	Live        *mock_database.MockLive
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
		db:         newDBMocks(ctrl),
		user:       mock_user.NewMockService(ctrl),
		messenger:  mock_messenger.NewMockService(ctrl),
		media:      mock_media.NewMockService(ctrl),
		postalCode: mock_postalcode.NewMockClient(ctrl),
		ivs:        mock_ivs.NewMockClient(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		Address:     mock_database.NewMockAddress(ctrl),
		Category:    mock_database.NewMockCategory(ctrl),
		Order:       mock_database.NewMockOrder(ctrl),
		Product:     mock_database.NewMockProduct(ctrl),
		ProductTag:  mock_database.NewMockProductTag(ctrl),
		ProductType: mock_database.NewMockProductType(ctrl),
		Promotion:   mock_database.NewMockPromotion(ctrl),
		Shipping:    mock_database.NewMockShipping(ctrl),
		Schedule:    mock_database.NewMockSchedule(ctrl),
		Live:        mock_database.NewMockLive(ctrl),
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
			Address:     mocks.db.Address,
			Category:    mocks.db.Category,
			Order:       mocks.db.Order,
			Product:     mocks.db.Product,
			ProductTag:  mocks.db.ProductTag,
			ProductType: mocks.db.ProductType,
			Promotion:   mocks.db.Promotion,
			Shipping:    mocks.db.Shipping,
			Schedule:    mocks.db.Schedule,
			Live:        mocks.db.Live,
		},
		User:       mocks.user,
		Messenger:  mocks.messenger,
		Media:      mocks.media,
		PostalCode: mocks.postalCode,
		Ivs:        mocks.ivs,
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
			expect: store.ErrInvalidArgument,
		},
		{
			name:   "database not found",
			err:    database.ErrNotFound,
			expect: store.ErrNotFound,
		},
		{
			name:   "database failed precondition",
			err:    database.ErrFailedPrecondition,
			expect: store.ErrFailedPrecondition,
		},
		{
			name:   "database already exists",
			err:    database.ErrAlreadyExists,
			expect: store.ErrAlreadyExists,
		},
		{
			name:   "database deadline exceeded",
			err:    database.ErrDeadlineExceeded,
			expect: store.ErrDeadlineExceeded,
		},
		{
			name:   "postal code invalid argument",
			err:    postalcode.ErrInvalidArgument,
			expect: store.ErrInvalidArgument,
		},
		{
			name:   "postal code not found",
			err:    postalcode.ErrNotFound,
			expect: store.ErrNotFound,
		},
		{
			name:   "postal code unavailable",
			err:    postalcode.ErrUnavailable,
			expect: store.ErrUnavailable,
		},
		{
			name:   "postal code deadline exceeded",
			err:    postalcode.ErrTimeout,
			expect: store.ErrDeadlineExceeded,
		},
		{
			name:   "context canceled",
			err:    context.Canceled,
			expect: store.ErrCanceled,
		},
		{
			name:   "context deadline exceeded",
			err:    context.DeadlineExceeded,
			expect: store.ErrDeadlineExceeded,
		},
		{
			name:   "other error",
			err:    assert.AnError,
			expect: store.ErrInternal,
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
