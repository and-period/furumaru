package updater

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media/broadcast/updater"
	mediadb "github.com/and-period/furumaru/api/internal/media/database/mysql"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/and-period/furumaru/api/pkg/secret"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type params struct {
	logger     *zap.Logger
	waitGroup  *sync.WaitGroup
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
		return err
	}

	// AWS Secrets Managerの設定
	params.secret = secret.NewClient(awscfg)
	if err := a.getSecret(ctx, params); err != nil {
		return err
	}

	// Loggerの設定
	logger, err := log.NewSentryLogger(params.sentryDsn, log.WithLogLevel(a.LogLevel), log.WithSentryLevel("error"))
	if err != nil {
		return err
	}
	params.logger = logger

	// Databaseの設定
	dbClient, err := a.newDatabase("stores", params)
	if err != nil {
		return err
	}

	// Jobの設定
	jobParams := &updater.Params{
		WaitGroup: params.waitGroup,
		Database:  mediadb.NewDatabase(dbClient),
	}
	switch a.RunType {
	case "START":
		a.updater = updater.NewStarter(jobParams, updater.WithLogger(params.logger))
	case "CLOSE":
		return errors.New("cmd: not implemented")
	default:
		return fmt.Errorf("cmd: unknown scheduler type. type=%s", a.RunType)
	}
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
	return mysql.NewClient(
		params,
		mysql.WithLogger(p.logger),
		mysql.WithNow(p.now),
		mysql.WithTLS(a.DBEnabledTLS),
		mysql.WithLocation(location),
	)
}
