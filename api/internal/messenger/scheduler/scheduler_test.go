package scheduler

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	mock_messenger "github.com/and-period/furumaru/api/mock/messenger"
	mock_database "github.com/and-period/furumaru/api/mock/messenger/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

type mocks struct {
	db        *dbMocks
	messenger *mock_messenger.MockService
}

type dbMocks struct {
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

type testCaller func(ctx context.Context, t *testing.T, scheduler *scheduler)

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		db:        newDBMocks(ctrl),
		messenger: mock_messenger.NewMockService(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		Message:         mock_database.NewMockMessage(ctrl),
		MessageTemplate: mock_database.NewMockMessageTemplate(ctrl),
		Notification:    mock_database.NewMockNotification(ctrl),
		ReceivedQueue:   mock_database.NewMockReceivedQueue(ctrl),
		ReportTemplate:  mock_database.NewMockReportTemplate(ctrl),
		Schedule:        mock_database.NewMockSchedule(ctrl),
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
			Message:         mocks.db.Message,
			MessageTemplate: mocks.db.MessageTemplate,
			Notification:    mocks.db.Notification,
			ReceivedQueue:   mocks.db.ReceivedQueue,
			ReportTemplate:  mocks.db.ReportTemplate,
			Schedule:        mocks.db.Schedule,
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
		ctx, cancel := context.WithCancel(t.Context())
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
	params := func(target time.Time) *database.ListSchedulesParams {
		return &database.ListSchedulesParams{
			Types: []entity.ScheduleType{
				entity.ScheduleTypeNotification,
				entity.ScheduleTypeStartLive,
				entity.ScheduleTypeReviewProductRequest,
				entity.ScheduleTypeReviewExperienceRequest,
			},
			Statuses: []entity.ScheduleStatus{
				entity.ScheduleStatusWaiting,
				entity.ScheduleStatusProcessing,
			},
			Since: jst.BeginningOfDay(target),
			Until: target,
		}
	}
	schedules := func(messageType entity.ScheduleType) entity.Schedules {
		return entity.Schedules{
			{
				MessageType: messageType,
				MessageID:   "message-id",
				Status:      entity.ScheduleStatusWaiting,
				Count:       0,
				SentAt:      now,
			},
		}
	}

	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		target time.Time
		expect error
	}{
		{
			name: "success notification",
			setup: func(ctx context.Context, mocks *mocks) {
				const messageType = entity.ScheduleTypeNotification
				schedules := schedules(messageType)
				mocks.db.Schedule.EXPECT().List(ctx, params(now)).Return(schedules, nil)
				mocks.messenger.EXPECT().NotifyNotification(gomock.Any(), gomock.Any()).Return(nil)
				mocks.db.Schedule.EXPECT().UpsertProcessing(gomock.Any(), schedules[0]).Return(nil)
				mocks.db.Schedule.EXPECT().UpdateDone(gomock.Any(), messageType, "message-id").Return(nil)
			},
			target: now,
			expect: nil,
		},
		{
			name: "success start live",
			setup: func(ctx context.Context, mocks *mocks) {
				const messageType = entity.ScheduleTypeStartLive
				schedules := schedules(messageType)
				mocks.db.Schedule.EXPECT().List(ctx, params(now)).Return(schedules, nil)
				mocks.messenger.EXPECT().NotifyStartLive(gomock.Any(), gomock.Any()).Return(nil)
				mocks.db.Schedule.EXPECT().UpsertProcessing(gomock.Any(), schedules[0]).Return(nil)
				mocks.db.Schedule.EXPECT().UpdateDone(gomock.Any(), messageType, "message-id").Return(nil)
			},
			target: now,
			expect: nil,
		},
		{
			name: "failed to list schedules",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().List(ctx, params(now)).Return(nil, assert.AnError)
			},
			target: now,
			expect: assert.AnError,
		},
		{
			name: "unknown message type",
			setup: func(ctx context.Context, mocks *mocks) {
				const messageType = entity.ScheduleTypeUnknown
				schedules := schedules(messageType)
				mocks.db.Schedule.EXPECT().List(ctx, params(now)).Return(schedules, nil)
			},
			target: now,
			expect: nil,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testScheduler(tt.setup, func(ctx context.Context, t *testing.T, scheduler *scheduler) {
			err := scheduler.Run(ctx, tt.target)
			assert.ErrorIs(t, err, tt.expect)
		}, withNow(now)))
	}
}

func TestScheduler_execute(t *testing.T) {
	t.Parallel()

	now := time.Now()
	schedule := &entity.Schedule{
		MessageType: entity.ScheduleTypeNotification,
		MessageID:   "notification-id",
		Status:      entity.ScheduleStatusWaiting,
		Count:       0,
		SentAt:      now,
	}

	tests := []struct {
		name     string
		setup    func(ctx context.Context, mocks *mocks)
		schedule *entity.Schedule
		execute  func(ctx context.Context, schedule *entity.Schedule) error
		expect   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().UpsertProcessing(ctx, schedule).Return(nil)
				mocks.db.Schedule.EXPECT().UpdateDone(ctx, entity.ScheduleTypeNotification, "notification-id").Return(nil)
			},
			schedule: schedule,
			execute: func(ctx context.Context, schedule *entity.Schedule) error {
				return nil
			},
			expect: nil,
		},
		{
			name:  "success non execute",
			setup: func(ctx context.Context, mocks *mocks) {},
			schedule: &entity.Schedule{
				MessageType: entity.ScheduleTypeNotification,
				MessageID:   "notification-id",
				Status:      entity.ScheduleStatusProcessing,
				Count:       1,
				SentAt:      now.Add(-time.Minute),
				CreatedAt:   now.Add(-time.Minute),
				UpdatedAt:   now.Add(-time.Minute),
			},
			execute: func(ctx context.Context, schedule *entity.Schedule) error {
				return nil
			},
			expect: nil,
		},
		{
			name: "success canceled",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().UpdateCancel(ctx, entity.ScheduleTypeNotification, "notification-id").Return(nil)
			},
			schedule: &entity.Schedule{
				MessageType: entity.ScheduleTypeNotification,
				MessageID:   "notification-id",
				Status:      entity.ScheduleStatusProcessing,
				Count:       2,
				SentAt:      now.Add(-time.Hour),
				CreatedAt:   now.Add(-time.Hour),
				UpdatedAt:   now.Add(-time.Hour),
			},
			execute: func(ctx context.Context, schedule *entity.Schedule) error {
				return nil
			},
			expect: nil,
		},
		{
			name: "failed to upsert processing",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().UpsertProcessing(ctx, schedule).Return(assert.AnError)
			},
			schedule: schedule,
			execute: func(ctx context.Context, schedule *entity.Schedule) error {
				return nil
			},
			expect: assert.AnError,
		},
		{
			name: "failed to notify notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().UpsertProcessing(ctx, schedule).Return(nil)
			},
			schedule: schedule,
			execute: func(ctx context.Context, schedule *entity.Schedule) error {
				return assert.AnError
			},
			expect: assert.AnError,
		},
		{
			name: "failed to update done",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().UpsertProcessing(ctx, schedule).Return(nil)
				mocks.db.Schedule.EXPECT().UpdateDone(ctx, entity.ScheduleTypeNotification, "notification-id").Return(assert.AnError)
			},
			schedule: schedule,
			execute: func(ctx context.Context, schedule *entity.Schedule) error {
				return nil
			},
			expect: assert.AnError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testScheduler(tt.setup, func(ctx context.Context, t *testing.T, scheduler *scheduler) {
			err := scheduler.execute(ctx, tt.schedule, tt.execute)
			assert.ErrorIs(t, err, tt.expect)
		}, withNow(now)))
	}
}
