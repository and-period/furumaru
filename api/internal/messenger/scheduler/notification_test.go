package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/stretchr/testify/assert"
)

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
			name: "failed to notify notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().UpsertProcessing(ctx, schedule).Return(nil)
				mocks.messenger.EXPECT().NotifyNotification(ctx, in).Return(assert.AnError)
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
