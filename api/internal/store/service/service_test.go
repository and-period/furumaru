package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	mock_media "github.com/and-period/furumaru/api/mock/media"
	mock_messenger "github.com/and-period/furumaru/api/mock/messenger"
	mock_postalcode "github.com/and-period/furumaru/api/mock/pkg/postalcode"
	mock_database "github.com/and-period/furumaru/api/mock/store/database"
	mock_user "github.com/and-period/furumaru/api/mock/user"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mocks struct {
	db         *dbMocks
	user       *mock_user.MockService
	messenger  *mock_messenger.MockService
	media      *mock_media.MockService
	postalCode *mock_postalcode.MockClient
}

type dbMocks struct {
	Category    *mock_database.MockCategory
	Order       *mock_database.MockOrder
	Product     *mock_database.MockProduct
	ProductType *mock_database.MockProductType
	Promotion   *mock_database.MockPromotion
	Shipping    *mock_database.MockShipping
	Schedule    *mock_database.MockSchedule
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
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		Category:    mock_database.NewMockCategory(ctrl),
		Order:       mock_database.NewMockOrder(ctrl),
		Product:     mock_database.NewMockProduct(ctrl),
		ProductType: mock_database.NewMockProductType(ctrl),
		Promotion:   mock_database.NewMockPromotion(ctrl),
		Shipping:    mock_database.NewMockShipping(ctrl),
		Schedule:    mock_database.NewMockSchedule(ctrl),
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
			Category:    mocks.db.Category,
			Order:       mocks.db.Order,
			Product:     mocks.db.Product,
			ProductType: mocks.db.ProductType,
			Promotion:   mocks.db.Promotion,
			Shipping:    mocks.db.Shipping,
			Schedule:    mocks.db.Schedule,
		},
		User:       mocks.user,
		Messenger:  mocks.messenger,
		Media:      mocks.media,
		PostalCode: mocks.postalCode,
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
