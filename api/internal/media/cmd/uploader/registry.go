package uploader

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media/uploader"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/and-period/furumaru/api/pkg/storage"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"
)

type params struct {
	logger     *zap.Logger
	waitGroup  *sync.WaitGroup
	secret     secret.Client
	storage    storage.Bucket
	tmpStorage storage.Bucket
	cache      dynamodb.Client
	now        func() time.Time
	sentryDsn  string
}

func (a *app) inject(ctx context.Context) error {
	params := &params{
		logger:    zap.NewNop(),
		now:       jst.Now,
		waitGroup: &sync.WaitGroup{},
	}

	// AWS SDKの設定
	awscfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(a.AWSRegion))
	if err != nil {
		return fmt.Errorf("cmd: failed to load aws config: %w", err)
	}

	// AWS Secrets Managerの設定
	params.secret = secret.NewClient(awscfg)
	if err := a.getSecret(ctx, params); err != nil {
		return fmt.Errorf("cmd: failed to get secret: %w", err)
	}

	// Loggerの設定
	logger, err := log.NewSentryLogger(params.sentryDsn,
		log.WithLogLevel(a.LogLevel),
		log.WithSentryServerName(a.AppName),
		log.WithSentryEnvironment(a.Environment),
		log.WithSentryLevel("error"),
	)
	if err != nil {
		return fmt.Errorf("cmd: failed to create sentry logger: %w", err)
	}
	params.logger = logger

	// Amazon S3の設定
	storageParams := &storage.Params{
		Bucket: a.S3Bucket,
	}
	params.storage = storage.NewBucket(awscfg, storageParams)
	tmpStorageParams := &storage.Params{
		Bucket: a.S3TmpBucket,
	}
	params.tmpStorage = storage.NewBucket(awscfg, tmpStorageParams, storage.WithLogger(params.logger))

	// Amazon DynamoDBの設定
	dbParams := &dynamodb.Params{
		TablePrefix: "furumaru",
		TableSuffix: a.Environment,
	}
	params.cache = dynamodb.NewClient(awscfg, dbParams, dynamodb.WithLogger(params.logger))

	// Uploaderの設定
	uploaderParams := &uploader.Params{
		WaitGroup: params.waitGroup,
		Storage:   params.storage,
		Tmp:       params.tmpStorage,
		Cache:     params.cache,
	}
	a.uploader = uploader.NewUploader(uploaderParams, uploader.WithLogger(params.logger), uploader.WithCacheDomain(a.CDNDomain))
	a.logger = params.logger
	a.waitGroup = params.waitGroup
	return nil
}

func (a *app) getSecret(ctx context.Context, p *params) error {
	// Sentry認証情報の取得
	if a.SentrySecretName == "" {
		p.sentryDsn = a.SentryDsn
		return nil
	}
	secrets, err := p.secret.Get(ctx, a.SentrySecretName)
	if err != nil {
		return err
	}
	p.sentryDsn = secrets["dsn"]
	return nil
}
