package worker

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Worker interface {
	Lambda(ctx context.Context, event events.SQSEvent) error
}

type Params struct {
	WaitGroup *sync.WaitGroup
	Mailer    mailer.Client
	User      user.Service
}

type worker struct {
	now         func() time.Time
	logger      *zap.Logger
	waitGroup   *sync.WaitGroup
	mailer      mailer.Client
	user        user.Service
	concurrency int64
	maxRetries  int64
}

type options struct {
	logger      *zap.Logger
	concurrency int64
	maxRetries  int64
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

func WithMaxRetries(maxRetries int64) Option {
	return func(opts *options) {
		opts.maxRetries = maxRetries
	}
}

func NewWorker(params *Params, opts ...Option) Worker {
	dopts := &options{
		logger:      zap.NewNop(),
		concurrency: 1,
		maxRetries:  3,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &worker{
		now:         jst.Now,
		logger:      dopts.logger,
		waitGroup:   params.WaitGroup,
		mailer:      params.Mailer,
		user:        params.User,
		concurrency: dopts.concurrency,
		maxRetries:  dopts.maxRetries,
	}
}

func (w *worker) Lambda(ctx context.Context, event events.SQSEvent) error {
	sm := semaphore.NewWeighted(w.concurrency)
	eg, ectx := errgroup.WithContext(ctx)
	for _, record := range event.Records {
		if err := sm.Acquire(ctx, 1); err != nil {
			return err
		}
		record := record
		eg.Go(func() error {
			defer sm.Release(1)
			payload := &messenger.WorkerPayload{}
			if err := json.Unmarshal([]byte(record.Body), payload); err != nil {
				w.logger.Error("Failed to unmarshall sqs event", zap.Any("event", event), zap.Error(err))
				return nil // リトライ不要なためnilで返す
			}
			err := w.dispatch(ectx, record.MessageId, payload)
			if err == nil {
				return nil
			}
			w.logger.Error("Failed to dispatch", zap.Error(err))
			if w.retryable(err) {
				return err
			}
			return nil
		})
	}
	return eg.Wait()
}

func (w *worker) dispatch(ctx context.Context, queueID string, payload *messenger.WorkerPayload) error {
	w.logger.Debug("Dispatch", zap.String("queueId", queueID), zap.Any("payload", payload))
	if payload.Email != nil {
		if err := w.sendInfoMail(ctx, payload); err != nil {
			w.logger.Error("Failed to send email", zap.Error(err))
			return err
		}
	}
	// TODO: プッシュ通知
	return nil
}

func (w *worker) retryable(err error) bool {
	return errors.Is(err, mailer.ErrTimeout) ||
		errors.Is(err, mailer.ErrUnavailable) ||
		errors.Is(err, mailer.ErrInternal)
}
