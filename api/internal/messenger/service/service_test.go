package service

import (
	"context"
	"errors"
	"net/url"
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

var (
	errmock        = errors.New("some error")
	adminWebURL, _ = url.Parse("htts://admin.and-period.jp")
	userWebURL, _  = url.Parse("htts://user.and-period.jp")
)

type mocks struct {
	db       *dbMocks
	producer *mock_sqs.MockProducer
	user     *mock_user.MockService
}

type dbMocks struct {
	Contact         *mock_database.MockContact
	Message         *mock_database.MockMessage
	MessageTemplate *mock_database.MockMessageTemplate
	Notification    *mock_database.MockNotification
	ReceivedQueue   *mock_database.MockReceivedQueue
	ReportTemplate  *mock_database.MockReportTemplate
	Schedule        *mock_database.MockSchedule
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
		db:       newDBMocks(ctrl),
		producer: mock_sqs.NewMockProducer(ctrl),
		user:     mock_user.NewMockService(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		Contact:         mock_database.NewMockContact(ctrl),
		Message:         mock_database.NewMockMessage(ctrl),
		MessageTemplate: mock_database.NewMockMessageTemplate(ctrl),
		Notification:    mock_database.NewMockNotification(ctrl),
		ReceivedQueue:   mock_database.NewMockReceivedQueue(ctrl),
		ReportTemplate:  mock_database.NewMockReportTemplate(ctrl),
		Schedule:        mock_database.NewMockSchedule(ctrl),
	}
}

func newService(mocks *mocks, opts ...testOption) *service {
	dopts := &testOptions{
		now: jst.Now,
	}
	for i := range opts {
		opts[i](dopts)
	}
	adminWebURL := func() *url.URL {
		url := *adminWebURL // copy
		return &url
	}
	userWebURL := func() *url.URL {
		url := *userWebURL // copy
		return &url
	}
	return &service{
		now:         dopts.now,
		logger:      zap.NewNop(),
		waitGroup:   &sync.WaitGroup{},
		validator:   validator.NewValidator(),
		producer:    mocks.producer,
		adminWebURL: adminWebURL,
		userWebURL:  userWebURL,
		db: &database.Database{
			Contact:         mocks.db.Contact,
			Message:         mocks.db.Message,
			MessageTemplate: mocks.db.MessageTemplate,
			Notification:    mocks.db.Notification,
			ReceivedQueue:   mocks.db.ReceivedQueue,
			ReportTemplate:  mocks.db.ReportTemplate,
			Schedule:        mocks.db.Schedule,
		},
		user: mocks.user,
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
