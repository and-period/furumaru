package resizer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media/resizer"
	"github.com/and-period/furumaru/api/internal/store"
	storedb "github.com/and-period/furumaru/api/internal/store/database/mysql"
	storesrv "github.com/and-period/furumaru/api/internal/store/service"
	"github.com/and-period/furumaru/api/internal/user"
	userdb "github.com/and-period/furumaru/api/internal/user/database/mysql"
	usersrv "github.com/and-period/furumaru/api/internal/user/service"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rafaelhl/gorm-newrelic-telemetry-plugin/telemetry"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type params struct {
	logger     *zap.Logger
	waitGroup  *sync.WaitGroup
	aws        aws.Config
	storage    storage.Bucket
	secret     secret.Client
	now        func() time.Time
	dbHost     string
	dbPort     string
	dbUsername string
	dbPassword string
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
	params.aws = awscfg

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

	// Serviceの設定
	userService, err := a.newUserService(params)
	if err != nil {
		return fmt.Errorf("cmd: failed to create user service: %w", err)
	}
	storeService, err := a.newStoreService(params)
	if err != nil {
		return fmt.Errorf("cmd: failed to create store service: %w", err)
	}

	// Resizerの設定
	resizerParams := &resizer.Params{
		WaitGroup: params.waitGroup,
		Storage:   params.storage,
		User:      userService,
		Store:     storeService,
	}
	a.resizer = resizer.NewResizer(resizerParams, resizer.WithLogger(params.logger), resizer.WithCacheTTL(a.CDNCacheTTL))
	a.logger = params.logger
	a.waitGroup = params.waitGroup
	return nil
}

func (a *app) getSecret(ctx context.Context, p *params) error {
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// データベース認証情報の取得
		if a.DBSecretName == "" {
			p.dbHost = a.DBHost
			p.dbPort = a.DBPort
			p.dbUsername = a.DBUsername
			p.dbPassword = a.DBPassword
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.DBSecretName)
		if err != nil {
			return err
		}
		p.dbHost = secrets["host"]
		p.dbPort = secrets["port"]
		p.dbUsername = secrets["username"]
		p.dbPassword = secrets["password"]
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

func (a *app) newDatabase(dbname string, p *params) (*mysql.Client, error) {
	params := &mysql.Params{
		Socket:   a.DBSocket,
		Host:     p.dbHost,
		Port:     p.dbPort,
		Database: dbname,
		Username: p.dbUsername,
		Password: p.dbPassword,
	}
	location, err := time.LoadLocation(a.DBTimeZone)
	if err != nil {
		return nil, err
	}
	cli, err := mysql.NewClient(
		params,
		mysql.WithNow(p.now),
		mysql.WithTLS(a.DBEnabledTLS),
		mysql.WithLocation(location),
	)
	if err != nil {
		return nil, err
	}
	if err := cli.DB.Use(telemetry.NewNrTracer(dbname, p.dbHost, string(newrelic.DatastoreMySQL))); err != nil {
		return nil, err
	}
	return cli, nil
}

func (a *app) newUserService(p *params) (user.Service, error) {
	mysql, err := a.newDatabase("users", p)
	if err != nil {
		return nil, err
	}
	params := &usersrv.Params{
		WaitGroup: p.waitGroup,
		Database:  userdb.NewDatabase(mysql),
	}
	return usersrv.NewService(params, usersrv.WithLogger(p.logger)), nil
}

func (a *app) newStoreService(p *params) (store.Service, error) {
	mysql, err := a.newDatabase("stores", p)
	if err != nil {
		return nil, err
	}
	params := &storesrv.Params{
		WaitGroup: p.waitGroup,
		Database:  storedb.NewDatabase(mysql),
	}
	return storesrv.NewService(params, storesrv.WithLogger(p.logger)), nil
}
