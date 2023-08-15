package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/sfn"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type starter struct {
	now       func() time.Time
	logger    *zap.Logger
	waitGroup *sync.WaitGroup
	semaphore *semaphore.Weighted
	db        *database.Database
	sfn       sfn.StepFunction
}

func NewStarter(params *Params, opts ...Option) Scheduler {
	dopts := &options{
		logger:      zap.NewNop(),
		concurrency: 2,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &starter{
		now:       jst.Now,
		logger:    dopts.logger,
		waitGroup: params.WaitGroup,
		semaphore: semaphore.NewWeighted(dopts.concurrency),
		db:        params.Database,
		sfn:       params.StepFunction,
	}
}

func (s *starter) Lambda(ctx context.Context) (err error) {
	s.logger.Debug("Started Lambda function", zap.Time("now", s.now()))
	defer func() {
		s.logger.Debug("Finished Lambda function", zap.Time("now", s.now()), zap.Error(err))
	}()

	return s.run(ctx, s.now())
}

func (s *starter) Run(ctx context.Context, target time.Time) error {
	return s.run(ctx, target)
}

// run - ライブ配信のリソース作成と開始処理
func (s *starter) run(ctx context.Context, target time.Time) error {
	if err := s.startChannel(ctx, target); err != nil {
		return err
	}
	return s.createChannel(ctx, target)
}

// startChannel - ライブ配信を開始 (5分前)
func (s *starter) startChannel(ctx context.Context, target time.Time) error {
	params := &database.ListSchedulesParams{
		StartAtGte: target.Add(-5 * time.Minute), // マルシェ開催開始5分前〜
		EndAtLt:    target,                       // マルシェ開催終了前
	}
	schedules, err := s.db.Schedule.List(ctx, params)
	if err != nil {
		return err
	}
	if len(schedules) == 0 {
		return nil
	}

	eg, ectx := errgroup.WithContext(ctx)
	for i := range schedules {
		if err := s.semaphore.Acquire(ctx, 1); err != nil {
			return err
		}

		schedule := schedules[i]
		eg.Go(func() error {
			defer s.semaphore.Release(1)
			broadcast, err := s.db.Broadcast.GetByScheduleID(ectx, schedule.ID)
			if err != nil {
				return err
			}
			if broadcast.Status != entity.BroadcastStatusIdle {
				return nil // 停止中の場合のみ、起動処理を進める
			}
			// TODO: ライブ配信リソースの起動処理
			params := &database.UpdateBroadcastParams{
				Status: entity.BroadcastStatusActive,
			}
			return s.db.Broadcast.Update(ctx, broadcast.ID, params)
		})
	}
	return eg.Wait()
}

// createChannel - ライブ配信リソースの作成を開始 (30分前)
func (s *starter) createChannel(ctx context.Context, target time.Time) error {
	params := &database.ListSchedulesParams{
		StartAtGte: target.Add(-30 * time.Minute), // マルシェ開催開始30分前〜
		EndAtLt:    target,                        // マルシェ開催終了前
	}
	schedules, err := s.db.Schedule.List(ctx, params)
	if err != nil {
		return err
	}
	if len(schedules) == 0 {
		return nil
	}

	eg, ectx := errgroup.WithContext(ctx)
	for i := range schedules {
		if err := s.semaphore.Acquire(ctx, 1); err != nil {
			return err
		}

		schedule := schedules[i]
		eg.Go(func() error {
			defer s.semaphore.Release(1)
			broadcast, err := s.db.Broadcast.GetByScheduleID(ectx, schedule.ID)
			if err != nil {
				return err
			}
			if broadcast.Status != entity.BroadcastStatusDisabled {
				return nil // リソース未作成の場合のみ、作成処理を進める
			}
			payload := &CreatePayload{
				ScheduleID: broadcast.ScheduleID,
				ChannelInput: &CreateChannelPayload{
					Name:                   schedule.ID,
					StartTime:              schedule.StartAt.Format(time.RFC3339),
					InputLossImageSlateURI: schedule.ImageURL,
				},
				MP4Input: &CreateMp4Payload{
					OpeningVideoURL: schedule.OpeningVideoURL,
				},
				RtmpInput: &CreateRtmpPayload{
					StreamName: streamName,
				},
			}
			if err := s.sfn.StartExecution(ectx, payload); err != nil {
				return err
			}
			params := &database.UpdateBroadcastParams{
				Status: entity.BroadcastStatusWaiting,
			}
			return s.db.Broadcast.Update(ctx, broadcast.ID, params)
		})
	}
	return eg.Wait()
}
