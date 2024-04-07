package uploader

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
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
	storageURL  func() *url.URL
	ttl         time.Duration
}

type options struct {
	logger      *zap.Logger
	concurrency int64
	storageURL  *url.URL
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

func WithStorageURL(storageURL string) Option {
	return func(opts *options) {
		url, err := url.Parse(storageURL)
		if err != nil {
			return
		}
		opts.storageURL = url
	}
}

func WithCacheTTL(ttl time.Duration) Option {
	return func(opts *options) {
		opts.ttl = ttl
	}
}

func NewUploader(params *Params, opts ...Option) Uploader {
	defaultURL, _ := params.Storage.GetHost()
	dopts := &options{
		logger:      zap.NewNop(),
		concurrency: 1,
		storageURL:  defaultURL,
		ttl:         5 * time.Minute,
	}
	for i := range opts {
		opts[i](dopts)
	}
	storageURL := func() *url.URL {
		url := *dopts.storageURL // copy
		return &url
	}
	return &uploader{
		now:         jst.Now,
		logger:      dopts.logger,
		waitGroup:   params.WaitGroup,
		cache:       params.Cache,
		tmp:         params.Tmp,
		storage:     params.Storage,
		concurrency: dopts.concurrency,
		storageURL:  storageURL,
		ttl:         dopts.ttl,
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
	payload := &events.S3Event{}
	if err := json.Unmarshal([]byte(record.Body), payload); err != nil {
		u.logger.Error("Failed to unmarshall sqs event", zap.Any("event", record), zap.Error(err))
		return nil // リトライ不要なためnilで返す
	}
	for i := range payload.Records {
		err := u.run(ctx, &payload.Records[i])
		if err == nil {
			return nil
		}
		u.logger.Error("Failed to upload object", zap.Error(err))
		if u.isRetryable(err) {
			return err
		}
	}
	return nil
}

func (u *uploader) run(ctx context.Context, record *events.S3EventRecord) error {
	event := &entity.UploadEvent{}
	defer func() {
		if event.Status == entity.UploadStatusUnknown {
			return
		}
		if err := u.cache.Insert(ctx, event); err != nil {
			u.logger.Error("Failed to update upload event", zap.Any("event", event), zap.Error(err))
		}
	}()
	u.logger.Debug("Dispatch", zap.Any("record", record))
	// 画像のメタデータ取得
	key := record.S3.Object.URLDecodedKey
	metadata, err := u.tmp.GetMetadata(ctx, key)
	if err != nil {
		u.logger.Error("Failed to get metadata", zap.String("key", key), zap.Error(err))
		return err
	}
	event.Key = key
	if err := u.cache.Get(ctx, event); err != nil {
		u.logger.Error("Failed to get upload event", zap.String("key", key), zap.Error(err))
		return err
	}
	// バリデーション検証
	reg, err := event.Reguration()
	if err != nil {
		u.logger.Error("Unknown regulation", zap.String("key", key), zap.Error(err))
		event.SetResult(false, "", u.now())
		return err
	}
	if err := reg.Validate(metadata.ContentType, metadata.ContentLength); err != nil {
		u.logger.Warn("Failed to check validation", zap.Error(err))
		event.SetResult(false, "", u.now())
		return err
	}
	// 参照用S3バケットへコピーする
	md := u.newObjectMetadata(metadata.ContentType)
	if _, err := u.storage.Copy(ctx, u.tmp.GetBucketName(), key, key, md); err != nil {
		u.logger.Error("Failed to copy object", zap.String("key", key), zap.Error(err))
		event.SetResult(false, "", u.now())
		return err
	}
	// 結果の保存
	referenceURL := u.storageURL()
	referenceURL.Path = key
	if !reg.ShouldConvert(metadata.ContentType) {
		event.SetResult(true, referenceURL.String(), u.now())
		return nil
	}
	// ファイル変換が必要な場合は変換して保存
	convertedKey, err := u.uploadConvetFile(ctx, event, reg)
	if err != nil {
		u.logger.Error("Failed to convert file", zap.String("key", key), zap.Error(err))
		event.SetResult(false, "", u.now())
		return err
	}
	referenceURL.Path = convertedKey // 参照先としては変換後のファイルを指定
	event.SetResult(true, referenceURL.String(), u.now())
	return nil
}

func (u *uploader) newObjectMetadata(contentType string) map[string]string {
	return map[string]string{
		"Content-Type":  contentType,
		"Cache-Control": "max-age=" + u.ttl.String(),
	}
}

func (u *uploader) isRetryable(err error) bool {
	return errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded)
}
