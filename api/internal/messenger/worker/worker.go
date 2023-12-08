package worker

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/firebase/messaging"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/line"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

var errUnknownUserType = errors.New("worker: unknown user type")

type Worker interface {
	Lambda(ctx context.Context, event events.SQSEvent) error
}

type Params struct {
	WaitGroup      *sync.WaitGroup
	Mailer         mailer.Client
	Line           line.Client
	AdminMessaging messaging.Client
	UserMessaging  messaging.Client
	DB             *database.Database
	User           user.Service
}

type worker struct {
	now            func() time.Time
	logger         *zap.Logger
	waitGroup      *sync.WaitGroup
	mailer         mailer.Client
	line           line.Client
	adminMessaging messaging.Client
	userMessagging messaging.Client
	db             *database.Database
	user           user.Service
	concurrency    int64
	maxRetries     int64
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
		now:            jst.Now,
		logger:         dopts.logger,
		waitGroup:      params.WaitGroup,
		mailer:         params.Mailer,
		line:           params.Line,
		adminMessaging: params.AdminMessaging,
		userMessagging: params.UserMessaging,
		db:             params.DB,
		user:           params.User,
		concurrency:    dopts.concurrency,
		maxRetries:     dopts.maxRetries,
	}
}

func (w *worker) Lambda(ctx context.Context, event events.SQSEvent) (err error) {
	w.logger.Debug("Started Lambda function", zap.Time("now", w.now()))
	defer func() {
		w.logger.Debug("Finished Lambda function", zap.Time("now", w.now()), zap.Error(err))
	}()

	sm := semaphore.NewWeighted(w.concurrency)
	eg, ectx := errgroup.WithContext(ctx)
	for _, record := range event.Records {
		if err := sm.Acquire(ctx, 1); err != nil {
			return err
		}
		record := record
		eg.Go(func() error {
			defer sm.Release(1)
			return w.dispatch(ectx, record)
		})
	}
	return eg.Wait()
}

func (w *worker) dispatch(ctx context.Context, record events.SQSMessage) error {
	payload := &entity.WorkerPayload{}
	if err := json.Unmarshal([]byte(record.Body), payload); err != nil {
		w.logger.Error("Failed to unmarshall sqs event", zap.Any("event", record), zap.Error(err))
		return nil // リトライ不要なためnilで返す
	}
	err := w.run(ctx, payload)
	if err == nil {
		return nil
	}
	w.logger.Error("Failed to send message", zap.Error(err))
	if w.isRetryable(err) {
		return err
	}
	return nil
}

func (w *worker) run(ctx context.Context, payload *entity.WorkerPayload) error {
	w.logger.Debug("Dispatch", zap.String("queueId", payload.QueueID), zap.Any("payload", payload))
	queue, err := w.db.ReceivedQueue.Get(ctx, payload.QueueID)
	if err != nil {
		return err
	}
	if queue.Done {
		w.logger.Info("This queue is already done", zap.String("queueId", payload.QueueID))
		return nil
	}
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// メール通知
		if payload.Email == nil {
			return nil
		}
		return w.multiSendMail(ectx, payload)
	})
	eg.Go(func() error {
		// プッシュ通知
		if payload.Push == nil {
			return nil
		}
		return w.multiSendPush(ectx, payload)
	})
	eg.Go(func() error {
		// メッセージ作成
		if payload.Message == nil {
			return nil
		}
		return w.messenger(ectx, payload)
	})
	eg.Go(func() error {
		// システムレポート
		if payload.Report == nil {
			return nil
		}
		return w.reporter(ectx, payload)
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	return w.db.ReceivedQueue.UpdateDone(ctx, payload.QueueID, true)
}

func (w *worker) isRetryable(err error) bool {
	return errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, mailer.ErrInternal) ||
		errors.Is(err, mailer.ErrUnavailable) ||
		errors.Is(err, mailer.ErrTimeout) ||
		errors.Is(err, messaging.ErrResourceExhausted) ||
		errors.Is(err, messaging.ErrInternal) ||
		errors.Is(err, messaging.ErrUnavailable) ||
		errors.Is(err, messaging.ErrTimeout) ||
		errors.Is(err, line.ErrInternal) ||
		errors.Is(err, line.ErrUnavailable) ||
		errors.Is(err, line.ErrResourceExhausted) ||
		errors.Is(err, line.ErrTimeout)
}
