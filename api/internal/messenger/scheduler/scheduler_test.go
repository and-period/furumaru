package scheduler

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	mock_messenger "github.com/and-period/furumaru/api/mock/messenger"
	mock_database "github.com/and-period/furumaru/api/mock/messenger/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

var errmock = errors.New("some error")

type mocks struct {
	db        *dbMocks
	messenger *mock_messenger.MockService
}

type dbMocks struct {
	Contact        *mock_database.MockContact
	Notification   *mock_database.MockNotification
	ReceivedQueue  *mock_database.MockReceivedQueue
	ReportTemplate *mock_database.MockReportTemplate
	Schedule       *mock_database.MockSchedule
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

type testCaller func(ctx context.Context, t *testing.T, scheduler *scheduler)

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		db:        newDBMocks(ctrl),
		messenger: mock_messenger.NewMockService(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		Contact:        mock_database.NewMockContact(ctrl),
		Notification:   mock_database.NewMockNotification(ctrl),
		ReceivedQueue:  mock_database.NewMockReceivedQueue(ctrl),
		ReportTemplate: mock_database.NewMockReportTemplate(ctrl),
		Schedule:       mock_database.NewMockSchedule(ctrl),
	}
}

func newScheduler(mocks *mocks, opts ...testOption) *scheduler {
	dopts := &testOptions{
		now: jst.Now,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &scheduler{
		now:       dopts.now,
		logger:    zap.NewNop(),
		waitGroup: &sync.WaitGroup{},
		semaphore: semaphore.NewWeighted(1),
		db: &database.Database{
			Contact:        mocks.db.Contact,
			Notification:   mocks.db.Notification,
			ReceivedQueue:  mocks.db.ReceivedQueue,
			ReportTemplate: mocks.db.ReportTemplate,
			Schedule:       mocks.db.Schedule,
		},
		messenger: mocks.messenger,
	}
}

func testScheduler(
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

		w := newScheduler(mocks, opts...)
		setup(ctx, mocks)

		testFunc(ctx, t, w)
		w.waitGroup.Wait()
	}
}

func TestScheduler(t *testing.T) {
	t.Parallel()
	s := NewScheduler(&Params{}, WithLogger(zap.NewNop()), WithConcurrency(1))
	assert.NotNil(t, s)
}

func TestScheduler_Run(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		target    time.Time
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().List(gomock.Any(), gomock.Any()).Return(entity.Schedules{}, nil)
				mocks.db.Notification.EXPECT().List(gomock.Any(), gomock.Any()).Return(entity.Notifications{}, nil)
			},
			target:    now,
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testScheduler(tt.setup, func(ctx context.Context, t *testing.T, scheduler *scheduler) {
			err := scheduler.Run(ctx, tt.target)
			assert.ErrorIs(t, err, tt.expectErr)
		}, withNow(now)))
	}
}
