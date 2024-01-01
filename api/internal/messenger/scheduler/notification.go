package scheduler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

func (s *scheduler) executeNotification(ctx context.Context, schedule *entity.Schedule) error {
	fn := func(ctx context.Context, schedule *entity.Schedule) error {
		in := &messenger.NotifyNotificationInput{
			NotificationID: schedule.MessageID,
		}
		return s.messenger.NotifyNotification(ctx, in)
	}
	return s.execute(ctx, schedule, fn)
}
