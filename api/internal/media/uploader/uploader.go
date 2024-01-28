package uploader

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Uploader interface {
	Lambda(ctx context.Context, event events.SQSEvent) error
}

type Params struct {
	WaitGroup *sync.WaitGroup
	Cache     dynamodb.Client
	Tmp       storage.Bucket
	Storage   storage.Bucket
}

type uploader struct {
	now         func() time.Time
	logger      *zap.Logger
	waitGroup   *sync.WaitGroup
	cache       dynamodb.Client
	tmp         storage.Bucket
	storage     storage.Bucket
	concurrency int64
	cacheTTL    time.Duration
}

type options struct {
	logger      *zap.Logger
	concurrency int64
	cacheTTL    time.Duration
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

func WithCacheTTL(ttl time.Duration) Option {
	return func(opts *options) {
		opts.cacheTTL = ttl
	}
}

func NewUploader(params *Params, opts ...Option) Uploader {
	dopts := &options{
		logger:      zap.NewNop(),
		concurrency: 1,
		cacheTTL:    6 * time.Hour,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &uploader{
		now:         jst.Now,
		logger:      dopts.logger,
		waitGroup:   params.WaitGroup,
		cache:       params.Cache,
		tmp:         params.Tmp,
		storage:     params.Storage,
		concurrency: dopts.concurrency,
		cacheTTL:    dopts.cacheTTL,
	}
}

func (u *uploader) Lambda(ctx context.Context, events events.SQSEvent) (err error) {
	u.logger.Debug("Started Lambda function", zap.Time("now", u.now()))
	defer func() {
		u.logger.Debug("Finished Lambda function", zap.Time("now", u.now()), zap.Error(err))
	}()

	sm := semaphore.NewWeighted(u.concurrency)
	eg, ectx := errgroup.WithContext(ctx)
	for _, record := range events.Records {
		if err := sm.Acquire(ctx, 1); err != nil {
			return err
		}

		record := record
		eg.Go(func() error {
			defer sm.Release(1)
			return u.dispatch(ectx, record)
		})
	}
	return eg.Wait()
}

func (u *uploader) dispatch(ctx context.Context, record events.SQSMessage) error {
	payload := &events.S3EventRecord{}
	if err := json.Unmarshal([]byte(record.Body), payload); err != nil {
		u.logger.Error("Failed to unmarshall sqs event", zap.Any("event", record), zap.Error(err))
		return nil // リトライ不要なためnilで返す
	}
	err := u.run(ctx, payload)
	if err == nil {
		return nil
	}
	u.logger.Error("Failed to upload object", zap.Error(err))
	if u.isRetryable(err) {
		return err
	}
	return nil
}

func (u *uploader) run(ctx context.Context, record *events.S3EventRecord) error {
	u.logger.Debug("Dispatch", zap.Any("record", record))
	key := record.S3.Object.URLDecodedKey
	metadata, err := u.tmp.GetMetadata(ctx, key)
	if err != nil {
		u.logger.Error("Failed to get metadata", zap.String("key", key), zap.Error(err))
		return err
	}
	reg, err := entity.FindByObjectKey(key)
	if err != nil {
		u.logger.Error("Unknown regulation", zap.String("key", key), zap.Error(err))
		return err
	}
	if err := reg.Validate(metadata.ContentType, metadata.ContentLength); err != nil {
		u.logger.Warn("Failed to check validation", zap.Error(err))
		return err
	}
	_, err = u.storage.Copy(ctx, u.tmp.GenerateS3URI(key), key)
	if err != nil {
		u.logger.Error("Failed to copy object", zap.String("key", key), zap.Error(err))
		return err
	}
	// TODO: 結果をキャッシュに保存
	return nil
}

func (u *uploader) isRetryable(err error) bool {
	return errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded)
}
