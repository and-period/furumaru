package scheduler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

func (s *scheduler) executeReviewRequest(ctx context.Context, schedule *entity.Schedule) error {
	fn := func(ctx context.Context, schedule *entity.Schedule) error {
		in := &messenger.NotifyReviewRequestInput{
			OrderID: schedule.MessageID,
		}
		return s.messenger.NotifyReviewRequest(ctx, in)
	}
	return s.execute(ctx, schedule, fn)
}
