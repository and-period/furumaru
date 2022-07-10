package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	mock_database "github.com/and-period/furumaru/api/mock/messenger/database"
	mock_sqs "github.com/and-period/furumaru/api/mock/pkg/sqs"
	mock_user "github.com/and-period/furumaru/api/mock/user"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/validator"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var errmock = errors.New("some error")

type mocks struct {
	db       *dbMocks
	producer *mock_sqs.MockProducer
	user     *mock_user.MockService
}

type dbMocks struct {
	Notification *mock_database.MockNotification
	Contact      *mock_database.MockContact
}

type testOptions struct {
	now func() time.Time
}

type testOption func(opts *testOptions)

type testCaller func(ctx context.Context, t *testing.T, service *service)

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		db:       newDBMocks(ctrl),
		producer: mock_sqs.NewMockProducer(ctrl),
		user:     mock_user.NewMockService(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		Notification: mock_database.NewMockNotification(ctrl),
		Contact:      mock_database.NewMockContact(ctrl),
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
		now:       dopts.now,
		logger:    zap.NewNop(),
		waitGroup: &sync.WaitGroup{},
		validator: validator.NewValidator(),
		db: &database.Database{
			Notification: mocks.db.Notification,
			Contact:      mocks.db.Contact,
		},
		producer: mocks.producer,
		user:     mocks.user,
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
