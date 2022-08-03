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
	eg, ectx := errgroup.WithContext(ctx)
	for _, scheduleType := range entity.ScheduleTypes {
		scheduleType := scheduleType
		eg.Go(func() (err error) {
			switch scheduleType {
			case entity.ScheduleTypeNotification:
				err = s.dispatchNotification(ectx, target)
			case entity.ScheduleTypeUnknown:
				return // 何もしない
			}
			return
		})
	}
	return eg.Wait()
}
