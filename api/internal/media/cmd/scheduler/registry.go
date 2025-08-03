package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media/broadcast/scheduler"
	mediadb "github.com/and-period/furumaru/api/internal/media/database/tidb"
	"github.com/and-period/furumaru/api/internal/store"
	storedb "github.com/and-period/furumaru/api/internal/store/database/tidb"
	storesrv "github.com/and-period/furumaru/api/internal/store/service"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/mediaconvert"
	"github.com/and-period/furumaru/api/pkg/medialive"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/and-period/furumaru/api/pkg/sfn"
	"github.com/and-period/furumaru/api/pkg/storage"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type params struct {
	logger       *zap.Logger
	waitGroup    *sync.WaitGroup
	secret       secret.Client
	now          func() time.Time
	tidbHost     string
	tidbPort     string
	tidbUsername string
	tidbPassword string
	sentryDsn    string
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

	// AWS Step Functionsの設定
	sfnParams := &sfn.Params{
		StateMachineARN: a.StepFunctionARN,
	}
	sfnClient := sfn.NewStepFunction(awscfg, sfnParams)

	// AWS Media Liveの設定
	mediaLiveClient := medialive.NewMediaLive(awscfg)

	// AWS Media Convertの設定
	mediaConvertParams := mediaconvert.Params{
		Endpoint: a.MediaConvertEndpoint,
		RoleARN:  a.MediaConvertRoleARN,
	}
	mediaConvertClient := mediaconvert.NewMediaConvert(awscfg, &mediaConvertParams)

	// Amazon S3の設定
	storageParams := &storage.Params{
		Bucket: a.ArchiveBucketName,
	}
	storageClient := storage.NewBucket(awscfg, storageParams)

	// Databaseの設定
	dbClient, err := a.newTiDB("media", params)
	if err != nil {
		return fmt.Errorf("cmd: failed to create database client: %w", err)
	}

	// Serviceの設定
	storeService, err := a.newStoreService(params)
	if err != nil {
		return fmt.Errorf("cmd: failed to create store service: %w", err)
	}

	// Jobの設定
	jobParams := &scheduler.Params{
		WaitGroup:          params.waitGroup,
		Database:           mediadb.NewDatabase(dbClient),
		Storage:            storageClient,
		Store:              storeService,
		StepFunction:       sfnClient,
		MediaLive:          mediaLiveClient,
		MediaConvert:       mediaConvertClient,
		Environment:        a.Environment,
		ArchiveBucketName:  a.ArchiveBucketName,
		ConvertJobTemplate: a.MediaConvertJobTemplate,
	}
	switch a.RunType {
	case "START":
		a.job = scheduler.NewStarter(jobParams, scheduler.WithLogger(params.logger))
	case "CLOSE":
		a.job = scheduler.NewCloser(jobParams, scheduler.WithLogger(params.logger), scheduler.WithStorageURL(a.CDNURL))
	default:
		return fmt.Errorf("cmd: unknown scheduler type. type=%s", a.RunType)
	}
	a.waitGroup = params.waitGroup
	return nil
}

func (a *app) getSecret(ctx context.Context, p *params) error {
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// データベース（TiDB）認証情報の取得
		if a.TiDBSecretName == "" {
			p.tidbHost = a.TiDBHost
			p.tidbPort = a.TiDBPort
			p.tidbUsername = a.TiDBUsername
			p.tidbPassword = a.TiDBPassword
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.TiDBSecretName)
		if err != nil {
			return err
		}
		p.tidbHost = secrets["host"]
		p.tidbPort = secrets["port"]
		p.tidbUsername = secrets["username"]
		p.tidbPassword = secrets["password"]
		return nil
	})
	eg.Go(func() error {
		// Sentry認証情報の取得
		if a.SentrySecretName == "" {
			p.sentryDsn = a.SentryDsn
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.SentrySecretName)
		if err != nil {
			return err
		}
		p.sentryDsn = secrets["dsn"]
		return nil
	})
	return eg.Wait()
}

func (a *app) newTiDB(dbname string, p *params) (*mysql.Client, error) {
	params := &mysql.Params{
		Host:     p.tidbHost,
		Port:     p.tidbPort,
		Database: dbname,
		Username: p.tidbUsername,
		Password: p.tidbPassword,
	}
	location, err := time.LoadLocation(a.DBTimeZone)
	if err != nil {
		return nil, err
	}
	return mysql.NewTiDBClient(
		params,
		mysql.WithNow(p.now),
		mysql.WithLocation(location),
	)
}

func (a *app) newStoreService(p *params) (store.Service, error) {
	mysql, err := a.newTiDB("stores", p)
	if err != nil {
		return nil, err
	}
	params := &storesrv.Params{
		WaitGroup: p.waitGroup,
		Database:  storedb.NewDatabase(mysql),
	}
	return storesrv.NewService(params, storesrv.WithLogger(p.logger)), nil
}
