package scheduler

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger"
	messengerdb "github.com/and-period/furumaru/api/internal/messenger/database/tidb"
	"github.com/and-period/furumaru/api/internal/messenger/scheduler"
	messengersrv "github.com/and-period/furumaru/api/internal/messenger/service"
	"github.com/and-period/furumaru/api/internal/store"
	storedb "github.com/and-period/furumaru/api/internal/store/database/tidb"
	storesrv "github.com/and-period/furumaru/api/internal/store/service"
	"github.com/and-period/furumaru/api/internal/user"
	userdb "github.com/and-period/furumaru/api/internal/user/database/tidb"
	usersrv "github.com/and-period/furumaru/api/internal/user/service"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/and-period/furumaru/api/pkg/sqs"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rafaelhl/gorm-newrelic-telemetry-plugin/telemetry"
	"golang.org/x/sync/errgroup"
)

type params struct {
	waitGroup    *sync.WaitGroup
	aws          aws.Config
	secret       secret.Client
	producer     sqs.Producer
	adminWebURL  *url.URL
	userWebURL   *url.URL
	now          func() time.Time
	tidbHost     string
	tidbPort     string
	tidbUsername string
	tidbPassword string
	sentryDsn    string
}

func (a *app) inject(ctx context.Context) error {
	params := &params{
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

	// Amazon SQSの設定
	sqsParams := &sqs.Params{
		QueueURL: a.SQSQueueURL,
	}
	params.producer = sqs.NewProducer(awscfg, sqsParams, sqs.WithDryRun(a.SQSMockEnabled))

	// WebURLの設定
	adminWebURL, err := url.Parse(a.AminWebURL)
	if err != nil {
		return fmt.Errorf("cmd: failed to parse admin web url: %w", err)
	}
	params.adminWebURL = adminWebURL
	userWebURL, err := url.Parse(a.UserWebURL)
	if err != nil {
		return fmt.Errorf("cmd: failed to parse user web url: %w", err)
	}
	params.userWebURL = userWebURL

	// Databaseの設定
	dbClient, err := a.newTiDB("messengers", params)
	if err != nil {
		return fmt.Errorf("cmd: failed to create database client: %w", err)
	}

	// Serviceの設定
	messengerService, err := a.newMessengerService(params)
	if err != nil {
		return fmt.Errorf("cmd: failed to create messenger service: %w", err)
	}

	// Jobの設定
	jobParams := &scheduler.Params{
		WaitGroup: params.waitGroup,
		Database:  messengerdb.NewDatabase(dbClient),
		Messenger: messengerService,
	}
	a.job = scheduler.NewScheduler(jobParams)
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
	cli, err := mysql.NewTiDBClient(
		params,
		mysql.WithNow(p.now),
		mysql.WithLocation(location),
	)
	if err != nil {
		return nil, err
	}
	if err := cli.DB.Use(telemetry.NewNrTracer(dbname, p.tidbHost, string(newrelic.DatastoreMySQL))); err != nil {
		return nil, err
	}
	return cli, nil
}

func (a *app) newMessengerService(p *params) (messenger.Service, error) {
	mysql, err := a.newTiDB("messengers", p)
	if err != nil {
		return nil, err
	}
	user, err := a.newUserService(p)
	if err != nil {
		return nil, err
	}
	store, err := a.newStoreService(p)
	if err != nil {
		return nil, err
	}
	params := &messengersrv.Params{
		WaitGroup:   p.waitGroup,
		Producer:    p.producer,
		AdminWebURL: p.adminWebURL,
		UserWebURL:  p.userWebURL,
		Database:    messengerdb.NewDatabase(mysql),
		User:        user,
		Store:       store,
	}
	return messengersrv.NewService(params), nil
}

func (a *app) newUserService(p *params) (user.Service, error) {
	mysql, err := a.newTiDB("users", p)
	if err != nil {
		return nil, err
	}
	params := &usersrv.Params{
		WaitGroup: p.waitGroup,
		Database:  userdb.NewDatabase(mysql),
	}
	return usersrv.NewService(params), nil
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
	return storesrv.NewService(params), nil
}
