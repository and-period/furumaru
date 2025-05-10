package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/stretchr/testify/assert"
)

func TestScheduler_executeReviewRequest(t *testing.T) {
	t.Parallel()

	now := time.Now()
	schedule := &entity.Schedule{
		MessageType: entity.ScheduleTypeReviewProductRequest,
		MessageID:   "order-id",
		Status:      entity.ScheduleStatusWaiting,
		Count:       0,
		SentAt:      now,
	}
	in := &messenger.NotifyReviewRequestInput{
		OrderID: "order-id",
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
				mocks.messenger.EXPECT().NotifyReviewRequest(ctx, in).Return(nil)
				mocks.db.Schedule.EXPECT().UpdateDone(ctx, entity.ScheduleTypeReviewProductRequest, "order-id").Return(nil)
			},
			schedule:  schedule,
			expectErr: nil,
		},
		{
			name: "failed to notify review request",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Schedule.EXPECT().UpsertProcessing(ctx, schedule).Return(nil)
				mocks.messenger.EXPECT().NotifyReviewRequest(ctx, in).Return(assert.AnError)
			},
			schedule:  schedule,
			expectErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testScheduler(tt.setup, func(ctx context.Context, t *testing.T, scheduler *scheduler) {
			err := scheduler.executeReviewRequest(ctx, tt.schedule)
			assert.ErrorIs(t, err, tt.expectErr)
		}, withNow(now)))
	}
}
