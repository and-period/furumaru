package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestScheduler_dispatchNotication(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 7, 19, 18, 30, 0, 0)
	schedulesParams := &database.ListSchedulesParams{
		Types: []entity.ScheduleType{entity.ScheduleTypeNotification},
		Since: jst.BeginningOfDay(now),
		Until: now,
	}
	notificationsParams := &database.ListNotificationsParams{
		Since:         jst.BeginningOfDay(now),
		Until:         now,
		OnlyPublished: true,
	}
	schedules := entity.Schedules{
		{
			MessageType: entity.ScheduleTypeNotification,
			MessageID:   "notification-id01",
			Status:      entity.ScheduleStatusProcessing,
			Count:       1,
			SentAt:      now,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	notifications := entity.Notifications{
		{ID: "notification-id01", PublishedAt: now},
		{ID: "notification-id02", PublishedAt: now},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		target    time.Time
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				schedule := &entity.Schedule{
					MessageType: entity.ScheduleTypeNotification,
					MessageID:   "notification-id02",
					Status:      entity.ScheduleStatusWaiting,
					Count:       0,
					SentAt:      now,
				}
				in := &messenger.NotifyNotificationInput{
					NotificationID: "notification-id02",
				}
				mocks.db.Schedule.EXPECT().List(gomock.Any(), schedulesParams).Return(schedules, nil)
				mocks.db.Notification.EXPECT().List(gomock.Any(), notificationsParams).Return(notifications, nil)
				mocks.db.Schedule.EXPECT().UpsertProcessing(gomock.Any(), schedule).Return(nil)
				mocks.messenger.EXPECT().NotifyNotification(gomock.Any(), in).Return(nil)
				mocks.db.Schedule.EXPECT().
					UpdateDone(gomock.Any(), entity.ScheduleTypeNotification, "notification-id02").
					Return(nil)
			},
			target:    now,
			expectErr: nil,
		},
		{
			name: "failed to list schedules",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().List(gomock.Any(), schedulesParams).Return(nil, assert.AnError)
				mocks.db.Notification.EXPECT().List(gomock.Any(), notificationsParams).Return(notifications, nil)
			},
			target:    now,
			expectErr: assert.AnError,
		},
		{
			name: "failed to list notifications",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().List(gomock.Any(), schedulesParams).Return(schedules, nil)
				mocks.db.Notification.EXPECT().List(gomock.Any(), notificationsParams).Return(nil, assert.AnError)
			},
			target:    now,
			expectErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testScheduler(tt.setup, func(ctx context.Context, t *testing.T, scheduler *scheduler) {
			err := scheduler.dispatchNotification(ctx, tt.target)
			assert.ErrorIs(t, err, tt.expectErr)
		}, withNow(now)))
	}
}

func TestScheduler_executeNotication(t *testing.T) {
	t.Parallel()

	now := time.Now()
	schedule := &entity.Schedule{
		MessageType: entity.ScheduleTypeNotification,
		MessageID:   "notification-id",
		Status:      entity.ScheduleStatusWaiting,
		Count:       0,
		SentAt:      now,
	}
	in := &messenger.NotifyNotificationInput{
		NotificationID: "notification-id",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		schedule  *entity.Schedule
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().UpsertProcessing(ctx, schedule).Return(nil)
				mocks.messenger.EXPECT().NotifyNotification(ctx, in).Return(nil)
				mocks.db.Schedule.EXPECT().UpdateDone(ctx, entity.ScheduleTypeNotification, "notification-id").Return(nil)
			},
			schedule:  schedule,
			expectErr: nil,
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
			expectErr: nil,
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
			expectErr: nil,
		},
		{
			name: "failed to upsert processing",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().UpsertProcessing(ctx, schedule).Return(assert.AnError)
			},
			schedule:  schedule,
			expectErr: assert.AnError,
		},
		{
			name: "failed to notify notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().UpsertProcessing(ctx, schedule).Return(nil)
				mocks.messenger.EXPECT().NotifyNotification(ctx, in).Return(assert.AnError)
			},
			schedule:  schedule,
			expectErr: assert.AnError,
		},
		{
			name: "failed to update done",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().UpsertProcessing(ctx, schedule).Return(nil)
				mocks.messenger.EXPECT().NotifyNotification(ctx, in).Return(nil)
				mocks.db.Schedule.EXPECT().UpdateDone(ctx, entity.ScheduleTypeNotification, "notification-id").Return(assert.AnError)
			},
			schedule:  schedule,
			expectErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testScheduler(tt.setup, func(ctx context.Context, t *testing.T, scheduler *scheduler) {
			err := scheduler.executeNotification(ctx, tt.schedule)
			assert.ErrorIs(t, err, tt.expectErr)
		}, withNow(now)))
	}
}
