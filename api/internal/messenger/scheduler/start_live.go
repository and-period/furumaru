package scheduler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

func (s *scheduler) executeStartLive(ctx context.Context, schedule *entity.Schedule) error {
	fn := func(ctx context.Context, schedule *entity.Schedule) error {
		in := &messenger.NotifyStartLiveInput{
			ScheduleID: schedule.MessageID,
		}
		return s.messenger.NotifyStartLive(ctx, in)
	}
	return s.execute(ctx, schedule, fn)
}
