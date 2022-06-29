package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	mock_messenger "github.com/and-period/furumaru/api/mock/messenger"
	mock_database "github.com/and-period/furumaru/api/mock/store/database"
	mock_user "github.com/and-period/furumaru/api/mock/user"
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
	user      *mock_user.MockService
	messenger *mock_messenger.MockService
}

type dbMocks struct {
	Category    *mock_database.MockCategory
	Product     *mock_database.MockProduct
	ProductType *mock_database.MockProductType
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
		user:      mock_user.NewMockService(ctrl),
		messenger: mock_messenger.NewMockService(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		Category:    mock_database.NewMockCategory(ctrl),
		Product:     mock_database.NewMockProduct(ctrl),
		ProductType: mock_database.NewMockProductType(ctrl),
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
		sharedGroup: &singleflight.Group{},
		waitGroup:   &sync.WaitGroup{},
		validator:   validator.NewValidator(),
		db: &database.Database{
			Category:    mocks.db.Category,
			Product:     mocks.db.Product,
			ProductType: mocks.db.ProductType,
		},
		user:      mocks.user,
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

		srv := newService(mocks)
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
