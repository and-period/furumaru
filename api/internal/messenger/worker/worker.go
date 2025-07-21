package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	const types = 4
	w.logger.Debug("Dispatch", zap.String("queueId", payload.QueueID), zap.Any("payload", payload))
	var mu sync.Mutex
	var errs error
	w.waitGroup.Add(types)
	go func() { // メール配信
		defer w.waitGroup.Done()
		if payload.Email == nil {
			return
		}
		err := w.execute(ctx, entity.NotifyTypeEmail, payload, w.multiSendMail)
		if err == nil {
			return
		}
		w.logger.Error(
			"Failed to multi send mail",
			zap.String("queueId", payload.QueueID),
			zap.Error(err),
		)
		mu.Lock()
		errs = errors.Join(errs, err)
		mu.Unlock()
	}()
	go func() { // メッセージ通知
		defer w.waitGroup.Done()
		if payload.Message == nil {
			return
		}
		err := w.execute(ctx, entity.NotifyTypeMessage, payload, w.createMessages)
		if err == nil {
			return
		}
		w.logger.Error(
			"Failed to create messages",
			zap.String("queueId", payload.QueueID),
			zap.Error(err),
		)
		mu.Lock()
		errs = errors.Join(errs, err)
		mu.Unlock()
	}()
	go func() { // プッシュ通知
		defer w.waitGroup.Done()
		if payload.Push == nil {
			return
		}
		err := w.execute(ctx, entity.NotifyTypePush, payload, w.multiSendPush)
		if err == nil {
			return
		}
		w.logger.Error(
			"Failed to multi send push",
			zap.String("queueId", payload.QueueID),
			zap.Error(err),
		)
		mu.Lock()
		errs = errors.Join(errs, err)
		mu.Unlock()
	}()
	go func() { // システムレポート
		defer w.waitGroup.Done()
		if payload.Report == nil {
			return
		}
		err := w.execute(ctx, entity.NotifyTypeReport, payload, w.sendReport)
		if err == nil {
			return
		}
		w.logger.Error(
			"Failed to send report",
			zap.String("queueId", payload.QueueID),
			zap.Error(err),
		)
		mu.Lock()
		errs = errors.Join(errs, err)
		mu.Unlock()
	}()
	w.waitGroup.Wait()
	return errs
}

func (w *worker) execute(
	ctx context.Context,
	notifyType entity.NotifyType,
	payload *entity.WorkerPayload,
	sendFn func(context.Context, *entity.WorkerPayload) error,
) error {
	queue, err := w.db.ReceivedQueue.Get(ctx, payload.QueueID, notifyType)
	if err != nil {
		return fmt.Errorf("worker: failed to get received queue: %w", err)
	}
	if queue.Done {
		w.logger.Info(
			"This queue is already done",
			zap.String("queueId", payload.QueueID),
			zap.Int32("notifyType", int32(notifyType)),
		)
		return nil
	}
	if err := sendFn(ctx, payload); err != nil {
		return fmt.Errorf("worker: failed to send function: %w", err)
	}
	if err := w.db.ReceivedQueue.UpdateDone(ctx, payload.QueueID, notifyType, true); err != nil {
		return fmt.Errorf("worker: failed to update done: %w", err)
	}
	return nil
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
