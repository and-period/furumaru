package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type closer struct {
	now       func() time.Time
	logger    *zap.Logger
	waitGroup *sync.WaitGroup
	semaphore *semaphore.Weighted
	db        *database.Database
}

func NewCloser(params *Params, opts ...Option) Scheduler {
	dopts := &options{
		logger:      zap.NewNop(),
		concurrency: 2,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &closer{
		now:       jst.Now,
		logger:    dopts.logger,
		waitGroup: params.WaitGroup,
		semaphore: semaphore.NewWeighted(dopts.concurrency),
		db:        params.Database,
	}
}

func (c *closer) Lambda(ctx context.Context) (err error) {
	c.logger.Debug("Started Lambda function", zap.Time("now", c.now()))
	defer func() {
		c.logger.Debug("Finished Lambda function", zap.Time("now", c.now()), zap.Error(err))
	}()

	return c.run(ctx, c.now())
}

func (c *closer) Run(ctx context.Context, target time.Time) error {
	return c.run(ctx, target)
}

// run - ライブ配信のリソース削除と停止処理
func (c *closer) run(ctx context.Context, target time.Time) error {
	if err := c.removeChannel(ctx, target); err != nil {
		return err
	}
	return c.stopChannel(ctx, target)
}

// stopChannel - ライブ配信を停止 (1時間後)
func (c *closer) stopChannel(ctx context.Context, target time.Time) error {
	params := &database.ListSchedulesParams{
		EndAtGte: target.AddDate(0, 0, -1),   // 〜マルシェ開催終了1日後
		EndAtLt:  target.Add(-1 * time.Hour), // 〜マルシェ開催終了1時間後
	}
	schedules, err := c.db.Schedule.List(ctx, params)
	if err != nil {
		return err
	}
	if len(schedules) == 0 {
		return nil
	}

	eg, ectx := errgroup.WithContext(ctx)
	for i := range schedules {
		if err := c.semaphore.Acquire(ctx, 1); err != nil {
			return err
		}

		schedule := schedules[i]
		eg.Go(func() error {
			defer c.semaphore.Release(1)
			broadcast, err := c.db.Broadcast.GetByScheduleID(ectx, schedule.ID)
			if err != nil {
				return err
			}
			if broadcast.Status != entity.BroadcastStatusActive {
				return nil // 起動中の場合のみ、停止処理を進める
			}
			// TODO: ライブ配信リソースの停止処理
			params := &database.UpdateBroadcastParams{
				Status: entity.BroadcastStatusIdle,
			}
			return c.db.Broadcast.Update(ctx, broadcast.ID, params)
		})
	}
	return eg.Wait()
}

// removeChannel - ライブ配信を削除 (1時間後&&停止中)
func (c *closer) removeChannel(ctx context.Context, target time.Time) error {
	params := &database.ListSchedulesParams{
		EndAtGte: target.AddDate(0, 0, -1),   // 〜マルシェ開催終了1日後
		EndAtLt:  target.Add(-1 * time.Hour), // 〜マルシェ開催終了1時間後
	}
	schedules, err := c.db.Schedule.List(ctx, params)
	if err != nil {
		return err
	}
	if len(schedules) == 0 {
		return nil
	}

	eg, ectx := errgroup.WithContext(ctx)
	for i := range schedules {
		if err := c.semaphore.Acquire(ctx, 1); err != nil {
			return err
		}

		schedule := schedules[i]
		eg.Go(func() error {
			defer c.semaphore.Release(1)
			broadcast, err := c.db.Broadcast.GetByScheduleID(ectx, schedule.ID)
			if err != nil {
				return err
			}
			if broadcast.Status != entity.BroadcastStatusIdle {
				return nil // 停止中の場合のみ、削除処理を進める
			}
			// TODO: ライブ配信リソースの停止処理
			params := &database.UpdateBroadcastParams{
				Status: entity.BroadcastStatusWaiting,
			}
			return c.db.Broadcast.Update(ctx, broadcast.ID, params)
		})
	}
	return eg.Wait()
}
