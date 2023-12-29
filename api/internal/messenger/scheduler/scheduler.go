package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Scheduler interface {
	Run(ctx context.Context, target time.Time) error
	Lambda(ctx context.Context) error
}

type Params struct {
	WaitGroup *sync.WaitGroup
	Database  *database.Database
	Messenger messenger.Service
}

type scheduler struct {
	now       func() time.Time
	logger    *zap.Logger
	waitGroup *sync.WaitGroup
	semaphore *semaphore.Weighted
	db        *database.Database
	messenger messenger.Service
}

type options struct {
	logger      *zap.Logger
	concurrency int64
}

type Option func(*options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithConcurrency(concurrency int64) Option {
	return func(opts *options) {
		opts.concurrency = concurrency
	}
}

func NewScheduler(params *Params, opts ...Option) Scheduler {
	dopts := &options{
		logger:      zap.NewNop(),
		concurrency: 2,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &scheduler{
		now:       jst.Now,
		logger:    dopts.logger,
		waitGroup: params.WaitGroup,
		semaphore: semaphore.NewWeighted(dopts.concurrency),
		db:        params.Database,
		messenger: params.Messenger,
	}
}

func (s *scheduler) Lambda(ctx context.Context) (err error) {
	s.logger.Debug("Started Lambda function", zap.Time("now", s.now()))
	defer func() {
		s.logger.Debug("Finished Lambda function", zap.Time("now", s.now()), zap.Error(err))
	}()

	return s.run(ctx, s.now())
}

func (s *scheduler) Run(ctx context.Context, target time.Time) error {
	return s.run(ctx, target)
}

func (s *scheduler) run(ctx context.Context, target time.Time) error {
	params := &database.ListSchedulesParams{
		Types:    entity.ScheduleTypes,
		Statuses: []entity.ScheduleStatus{entity.ScheduleStatusWaiting, entity.ScheduleStatusProcessing},
		Since:    jst.BeginningOfDay(target),
		Until:    target,
	}
	schedules, err := s.db.Schedule.List(ctx, params)
	if err != nil {
		s.logger.Error("Failed to list schedules", zap.Error(err))
		return err
	}

	eg, ectx := errgroup.WithContext(ctx)
	for i := range schedules {
		if err := s.semaphore.Acquire(ctx, 1); err != nil {
			return err
		}

		schedule := schedules[i]
		eg.Go(func() error {
			defer s.semaphore.Release(1)
			return s.dispatch(ectx, schedule)
		})
	}
	return eg.Wait()
}

func (s *scheduler) dispatch(ctx context.Context, schedule *entity.Schedule) error {
	switch schedule.MessageType {
	case entity.ScheduleTypeNotification:
		return s.executeNotification(ctx, schedule)
	case entity.ScheduleTypeStartLive:
		return s.executeStartLive(ctx, schedule)
	default:
		s.logger.Warn("Received unknown message type", zap.Any("schedule", schedule))
		return nil // 何もしない
	}
}

func (s *scheduler) execute(ctx context.Context, schedule *entity.Schedule, fn func(context.Context, *entity.Schedule) error) error {
	now := s.now()
	// 通知前処理
	if schedule.ShouldCancel(now) {
		return s.db.Schedule.UpdateCancel(ctx, schedule.MessageType, schedule.MessageID)
	}
	if !schedule.Executable(now) {
		return nil
	}
	if err := s.db.Schedule.UpsertProcessing(ctx, schedule); err != nil {
		return err
	}
	// 通知処理
	if err := fn(ctx, schedule); err != nil {
		return err
	}
	// 通知後処理
	return s.db.Schedule.UpdateDone(ctx, schedule.MessageType, schedule.MessageID)
}
