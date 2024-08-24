package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/komoju"
	mock_media "github.com/and-period/furumaru/api/mock/media"
	mock_messenger "github.com/and-period/furumaru/api/mock/messenger"
	mock_dynamodb "github.com/and-period/furumaru/api/mock/pkg/dynamodb"
	mock_ivs "github.com/and-period/furumaru/api/mock/pkg/ivs"
	mock_postalcode "github.com/and-period/furumaru/api/mock/pkg/postalcode"
	mock_database "github.com/and-period/furumaru/api/mock/store/database"
	mock_komoju "github.com/and-period/furumaru/api/mock/store/komoju"
	mock_user "github.com/and-period/furumaru/api/mock/user"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/postalcode"
	govalidator "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type mocks struct {
	db            *dbMocks
	cache         *mock_dynamodb.MockClient
	user          *mock_user.MockService
	messenger     *mock_messenger.MockService
	media         *mock_media.MockService
	postalCode    *mock_postalcode.MockClient
	ivs           *mock_ivs.MockClient
	komojuPayment *mock_komoju.MockPayment
	komojuSession *mock_komoju.MockSession
}

type dbMocks struct {
	Category       *mock_database.MockCategory
	Experience     *mock_database.MockExperience
	ExperienceType *mock_database.MockExperienceType
	Order          *mock_database.MockOrder
	PaymentSystem  *mock_database.MockPaymentSystem
	Product        *mock_database.MockProduct
	ProductTag     *mock_database.MockProductTag
	ProductType    *mock_database.MockProductType
	Promotion      *mock_database.MockPromotion
	Shipping       *mock_database.MockShipping
	Schedule       *mock_database.MockSchedule
	Live           *mock_database.MockLive
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
		db:            newDBMocks(ctrl),
		cache:         mock_dynamodb.NewMockClient(ctrl),
		user:          mock_user.NewMockService(ctrl),
		messenger:     mock_messenger.NewMockService(ctrl),
		media:         mock_media.NewMockService(ctrl),
		postalCode:    mock_postalcode.NewMockClient(ctrl),
		ivs:           mock_ivs.NewMockClient(ctrl),
		komojuPayment: mock_komoju.NewMockPayment(ctrl),
		komojuSession: mock_komoju.NewMockSession(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		Category:       mock_database.NewMockCategory(ctrl),
		Experience:     mock_database.NewMockExperience(ctrl),
		ExperienceType: mock_database.NewMockExperienceType(ctrl),
		Order:          mock_database.NewMockOrder(ctrl),
		PaymentSystem:  mock_database.NewMockPaymentSystem(ctrl),
		Product:        mock_database.NewMockProduct(ctrl),
		ProductTag:     mock_database.NewMockProductTag(ctrl),
		ProductType:    mock_database.NewMockProductType(ctrl),
		Promotion:      mock_database.NewMockPromotion(ctrl),
		Shipping:       mock_database.NewMockShipping(ctrl),
		Schedule:       mock_database.NewMockSchedule(ctrl),
		Live:           mock_database.NewMockLive(ctrl),
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
			Category:       mocks.db.Category,
			Expericence:    mocks.db.Experience,
			ExperienceType: mocks.db.ExperienceType,
			Order:          mocks.db.Order,
			PaymentSystem:  mocks.db.PaymentSystem,
			Product:        mocks.db.Product,
			ProductTag:     mocks.db.ProductTag,
			ProductType:    mocks.db.ProductType,
			Promotion:      mocks.db.Promotion,
			Shipping:       mocks.db.Shipping,
			Schedule:       mocks.db.Schedule,
			Live:           mocks.db.Live,
		},
		Cache:      mocks.cache,
		User:       mocks.user,
		Messenger:  mocks.messenger,
		Media:      mocks.media,
		PostalCode: mocks.postalCode,
		Ivs:        mocks.ivs,
		Komoju: &komoju.Komoju{
			Payment: mocks.komojuPayment,
			Session: mocks.komojuSession,
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
			name:   "cache not found",
			err:    dynamodb.ErrNotFound,
			expect: exception.ErrNotFound,
		},
		{
			name:   "cache already exists",
			err:    dynamodb.ErrAlreadyExists,
			expect: exception.ErrAlreadyExists,
		},
		{
			name:   "cache resource exhausted",
			err:    dynamodb.ErrResourceExhausted,
			expect: exception.ErrResourceExhausted,
		},
		{
			name:   "cache canceled",
			err:    dynamodb.ErrCanceled,
			expect: exception.ErrCanceled,
		},
		{
			name:   "postal code invalid argument",
			err:    postalcode.ErrInvalidArgument,
			expect: exception.ErrInvalidArgument,
		},
		{
			name:   "postal code not found",
			err:    postalcode.ErrNotFound,
			expect: exception.ErrNotFound,
		},
		{
			name:   "postal code unavailable",
			err:    postalcode.ErrUnavailable,
			expect: exception.ErrUnavailable,
		},
		{
			name:   "postal code deadline exceeded",
			err:    postalcode.ErrTimeout,
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
