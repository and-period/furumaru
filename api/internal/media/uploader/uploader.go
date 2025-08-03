package uploader

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/url"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/aws/aws-lambda-go/events"
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
	waitGroup   *sync.WaitGroup
	cache       dynamodb.Client
	tmp         storage.Bucket
	storage     storage.Bucket
	concurrency int64
	storageURL  func() *url.URL
}

type options struct {
	concurrency int64
	storageURL  *url.URL
}

type Option func(*options)

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

func NewUploader(params *Params, opts ...Option) Uploader {
	defaultURL, _ := params.Storage.GetHost()
	dopts := &options{
		concurrency: 1,
		storageURL:  defaultURL,
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
		waitGroup:   params.WaitGroup,
		cache:       params.Cache,
		tmp:         params.Tmp,
		storage:     params.Storage,
		concurrency: dopts.concurrency,
		storageURL:  storageURL,
	}
}

func (u *uploader) Lambda(ctx context.Context, events events.SQSEvent) (err error) {
	slog.Debug("Started Lambda function", slog.Time("now", u.now()))
	defer func() {
		slog.Debug("Finished Lambda function", slog.Time("now", u.now()), log.Error(err))
	}()

	sm := semaphore.NewWeighted(u.concurrency)
	eg, ectx := errgroup.WithContext(ctx)
	for _, record := range events.Records {
		if err := sm.Acquire(ctx, 1); err != nil {
			return err
		}

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
		slog.Error("Failed to unmarshall sqs event", slog.Any("event", record), log.Error(err))
		return nil // リトライ不要なためnilで返す
	}
	for i := range payload.Records {
		err := u.run(ctx, &payload.Records[i])
		if err == nil {
			return nil
		}
		slog.Error("Failed to upload object", log.Error(err))
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
			slog.Error("Failed to update upload event", slog.Any("event", event), log.Error(err))
		}
	}()
	slog.Debug("Dispatch", slog.Any("record", record))
	// 画像のメタデータ取得
	key := record.S3.Object.URLDecodedKey
	metadata, err := u.tmp.GetMetadata(ctx, key)
	if err != nil {
		slog.Error("Failed to get metadata", slog.String("key", key), log.Error(err))
		return err
	}
	event.Key = key
	if err := u.cache.Get(ctx, event); err != nil {
		slog.Error("Failed to get upload event", slog.String("key", key), log.Error(err))
		return err
	}
	// バリデーション検証
	reg, err := event.Reguration()
	if err != nil {
		slog.Error("Unknown regulation", slog.String("key", key), log.Error(err))
		event.SetResult(false, "", u.now())
		return err
	}
	if err := reg.Validate(metadata.ContentType, metadata.ContentLength); err != nil {
		slog.Warn("Failed to check validation", log.Error(err))
		event.SetResult(false, "", u.now())
		return err
	}
	// 参照用S3バケットへコピーする
	md := u.newObjectMetadata(metadata.ContentType, reg.CacheTTL)
	if _, err := u.storage.Copy(ctx, u.tmp.GetBucketName(), key, key, md); err != nil {
		slog.Error("Failed to copy object", slog.String("key", key), log.Error(err))
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
	convertedKey, err := u.uploadConvertFile(ctx, event, reg)
	if err != nil {
		slog.Error("Failed to upload convert file", slog.String("key", key), log.Error(err))
		event.SetResult(false, "", u.now())
		return err
	}
	referenceURL.Path = convertedKey // 参照先としては変換後のファイルを指定
	event.SetResult(true, referenceURL.String(), u.now())
	return nil
}

func (u *uploader) newObjectMetadata(contentType string, ttl time.Duration) map[string]string {
	return map[string]string{
		"Content-Type":  contentType,
		"Cache-Control": "s-maxage=" + ttl.String(),
	}
}

func (u *uploader) isRetryable(err error) bool {
	return errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded)
}
