package resizer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/backoff"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

var (
	errRequiredMediaURL       = errors.New("resizer: required media url")
	errInvalidFormat          = errors.New("resizer: invalid format")
	errUnsupportedImageSize   = errors.New("resizer: unsupported image size")
	errUnsupportedImageFormat = errors.New("resizer: unsupported image format")
	errUnknownFileType        = errors.New("resizer: unknown file type")
)

type Resizer interface {
	Lambda(ctx context.Context, event events.SQSEvent) error
}

type Params struct {
	WaitGroup *sync.WaitGroup
	Storage   storage.Bucket
	User      user.Service
	Store     store.Service
}

type resizer struct {
	now         func() time.Time
	logger      *zap.Logger
	waitGroup   *sync.WaitGroup
	storage     storage.Bucket
	user        user.Service
	store       store.Service
	concurrency int64
	maxRetries  int64
	ttl         time.Duration
}

type options struct {
	logger      *zap.Logger
	concurrency int64
	maxRetries  int64
	ttl         time.Duration
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

func WithMaxRetires(maxRetries int64) Option {
	return func(opts *options) {
		opts.maxRetries = maxRetries
	}
}

func WithCacheTTL(ttl time.Duration) Option {
	return func(opts *options) {
		opts.ttl = ttl
	}
}

func NewResizer(params *Params, opts ...Option) Resizer {
	dopts := &options{
		logger:      zap.NewNop(),
		concurrency: 1,
		maxRetries:  3,
		ttl:         5 * time.Minute,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &resizer{
		now:         jst.Now,
		logger:      dopts.logger,
		waitGroup:   params.WaitGroup,
		storage:     params.Storage,
		user:        params.User,
		store:       params.Store,
		concurrency: dopts.concurrency,
		ttl:         dopts.ttl,
	}
}

func (r *resizer) Lambda(ctx context.Context, event events.SQSEvent) (err error) {
	r.logger.Debug("Started Lambda function", zap.Time("now", r.now()))
	defer func() {
		r.logger.Debug("Finished Lambda function", zap.Time("now", r.now()), zap.Error(err))
	}()

	sm := semaphore.NewWeighted(r.concurrency)
	eg, ectx := errgroup.WithContext(ctx)
	for _, record := range event.Records {
		if err := sm.Acquire(ctx, 1); err != nil {
			return err
		}
		record := record
		eg.Go(func() error {
			defer sm.Release(1)
			return r.dispatch(ectx, record)
		})
	}
	return eg.Wait()
}

func (r *resizer) dispatch(ctx context.Context, record events.SQSMessage) error {
	payload := &entity.ResizerPayload{}
	if err := json.Unmarshal([]byte(record.Body), payload); err != nil {
		r.logger.Error("Failed to unmarshall sqs event", zap.Any("event", record), zap.Error(err))
		return nil // リトライ不要なためnilで返す
	}
	err := r.run(ctx, payload)
	if err == nil {
		return nil
	}
	r.logger.Error("Failed to send message", zap.Error(err))
	if r.isRetryable(err) {
		return err
	}
	return nil
}

func (r *resizer) run(ctx context.Context, payload *entity.ResizerPayload) error {
	r.logger.Debug("Dispatch", zap.Int32("fileType", int32(payload.FileType)), zap.String("targetId", payload.TargetID))
	if len(payload.URLs) == 0 {
		return fmt.Errorf("resizer: urls is length 0: %w", errInvalidFormat)
	}
	switch payload.FileType {
	case entity.FileTypeCoordinatorThumbnail:
		return r.coordinatorThumbnail(ctx, payload)
	case entity.FileTypeCoordinatorHeader:
		return r.coordinatorHeader(ctx, payload)
	case entity.FileTypeProducerThumbnail:
		return r.producerThumbnail(ctx, payload)
	case entity.FileTypeProducerHeader:
		return r.producerHeader(ctx, payload)
	case entity.FileTypeUserThumbnail:
		return r.userThumbnail(ctx, payload)
	case entity.FileTypeProductMedia:
		return r.productMedia(ctx, payload)
	case entity.FileTypeProductTypeIcon:
		return r.productTypeIcon(ctx, payload)
	case entity.FileTypeScheduleThumbnail:
		return r.scheduleThumbnail(ctx, payload)
	default:
		return fmt.Errorf("resizer: unknown file type. type=%d: %w", payload.FileType, errUnknownFileType)
	}
}

func (r *resizer) notify(ctx context.Context, payload *entity.ResizerPayload, fn func() error) error {
	retry := backoff.NewExponentialBackoff(r.maxRetries)
	err := backoff.Retry(ctx, retry, fn, backoff.WithRetryablel(r.isRetryable))
	if err != nil {
		r.logger.Error("Failed to notify resize action",
			zap.Int32("fileType", int32(payload.FileType)),
			zap.String("targetId", payload.TargetID),
			zap.Error(err))
		return err
	}
	return nil
}

func (r *resizer) isRetryable(err error) bool {
	return errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded)
}
