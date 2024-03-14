package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/stretchr/testify/assert"
)

func TestReserveStartLive(t *testing.T) {
	t.Parallel()
	now := time.Now()
	scheduleIn := &store.GetScheduleInput{
		ScheduleID: "schedule-id",
	}
	live := &sentity.Schedule{
		ID:      "schedule-id",
		Status:  sentity.ScheduleStatusWaiting,
		StartAt: now.Add(time.Hour),
		EndAt:   now.Add(2 * time.Hour),
	}
	schedule := &entity.Schedule{
		MessageType: entity.ScheduleTypeStartLive,
		MessageID:   "schedule-id",
		Status:      entity.ScheduleStatusWaiting,
		Count:       0,
		SentAt:      now.Add(time.Hour - (10 * time.Minute)),
		Deadline:    now.Add(time.Hour),
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *messenger.ReserveStartLiveInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetSchedule(ctx, scheduleIn).Return(live, nil)
				mocks.db.Schedule.EXPECT().Upsert(ctx, schedule).Return(nil)
			},
			input: &messenger.ReserveStartLiveInput{
				ScheduleID: "schedule-id",
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &messenger.ReserveStartLiveInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "not found live schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetSchedule(ctx, scheduleIn).Return(nil, exception.ErrNotFound)
			},
			input: &messenger.ReserveStartLiveInput{
				ScheduleID: "schedule-id",
			},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get live schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetSchedule(ctx, scheduleIn).Return(nil, assert.AnError)
			},
			input: &messenger.ReserveStartLiveInput{
				ScheduleID: "schedule-id",
			},
			expect: exception.ErrInternal,
		},
		{
			name: "failed to upsert schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetSchedule(ctx, scheduleIn).Return(live, nil)
				mocks.db.Schedule.EXPECT().Upsert(ctx, schedule).Return(assert.AnError)
			},
			input: &messenger.ReserveStartLiveInput{
				ScheduleID: "schedule-id",
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ReserveStartLive(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestReserveNotification(t *testing.T) {
	t.Parallel()
	now := time.Now()
	notification := &entity.Notification{
		ID:          "notification-id",
		PublishedAt: now,
	}
	schedule := &entity.Schedule{
		MessageType: entity.ScheduleTypeNotification,
		MessageID:   "notification-id",
		Status:      entity.ScheduleStatusWaiting,
		Count:       0,
		SentAt:      now,
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *messenger.ReserveNotificationInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Notification.EXPECT().Get(ctx, "notification-id").Return(notification, nil)
				mocks.db.Schedule.EXPECT().Upsert(ctx, schedule).Return(nil)
			},
			input: &messenger.ReserveNotificationInput{
				NotificationID: "notification-id",
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &messenger.ReserveNotificationInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "not found notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Notification.EXPECT().Get(ctx, "notification-id").Return(nil, database.ErrNotFound)
			},
			input: &messenger.ReserveNotificationInput{
				NotificationID: "notification-id",
			},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Notification.EXPECT().Get(ctx, "notification-id").Return(nil, assert.AnError)
			},
			input: &messenger.ReserveNotificationInput{
				NotificationID: "notification-id",
			},
			expect: exception.ErrInternal,
		},
		{
			name: "failed to upsert schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Notification.EXPECT().Get(ctx, "notification-id").Return(notification, nil)
				mocks.db.Schedule.EXPECT().Upsert(ctx, schedule).Return(assert.AnError)
			},
			input: &messenger.ReserveNotificationInput{
				NotificationID: "notification-id",
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ReserveNotification(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}
