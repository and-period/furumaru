package user

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	v1 "github.com/and-period/furumaru/api/internal/gateway/user/v1/handler"
	"github.com/and-period/furumaru/api/internal/media"
	mediadb "github.com/and-period/furumaru/api/internal/media/database/mysql"
	mediasrv "github.com/and-period/furumaru/api/internal/media/service"
	"github.com/and-period/furumaru/api/internal/messenger"
	messengerdb "github.com/and-period/furumaru/api/internal/messenger/database/mysql"
	messengersrv "github.com/and-period/furumaru/api/internal/messenger/service"
	"github.com/and-period/furumaru/api/internal/store"
	storedb "github.com/and-period/furumaru/api/internal/store/database/mysql"
	"github.com/and-period/furumaru/api/internal/store/komoju"
	kpayment "github.com/and-period/furumaru/api/internal/store/komoju/payment"
	ksession "github.com/and-period/furumaru/api/internal/store/komoju/session"
	storesrv "github.com/and-period/furumaru/api/internal/store/service"
	"github.com/and-period/furumaru/api/internal/user"
	userdb "github.com/and-period/furumaru/api/internal/user/database/tidb"
	usersrv "github.com/and-period/furumaru/api/internal/user/service"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/and-period/furumaru/api/pkg/postalcode"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/and-period/furumaru/api/pkg/sentry"
	"github.com/and-period/furumaru/api/pkg/slack"
	"github.com/and-period/furumaru/api/pkg/sqs"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rafaelhl/gorm-newrelic-telemetry-plugin/telemetry"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type params struct {
	serviceName          string
	logger               *zap.Logger
	waitGroup            *sync.WaitGroup
	aws                  aws.Config
	secret               secret.Client
	storage              storage.Bucket
	tmpStorage           storage.Bucket
	userAuth             cognito.Client
	cache                dynamodb.Client
	producer             sqs.Producer
	slack                slack.Client
	newRelic             *newrelic.Application
	sentry               sentry.Client
	komoju               *komoju.Komoju
	adminWebURL          *url.URL
	userWebURL           *url.URL
	postalCode           postalcode.Client
	now                  func() time.Time
	debugMode            bool
	dbHost               string
	dbPort               string
	dbUsername           string
	dbPassword           string
	tidbHost             string
	tidbPort             string
	tidbUsername         string
	tidbPassword         string
	slackToken           string
	slackChannelID       string
	newRelicLicense      string
	sentryDsn            string
	komojuClientID       string
	komojuClientPassword string
}

//nolint:funlen
func (a *app) inject(ctx context.Context) error {
	params := &params{
		serviceName: fmt.Sprintf("%s-%s", a.AppName, a.Environment),
		logger:      zap.NewNop(),
		now:         jst.Now,
		waitGroup:   &sync.WaitGroup{},
		debugMode:   a.LogLevel == "debug",
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
	tmpStorageParams := &storage.Params{
		Bucket: a.S3TmpBucket,
	}
	params.tmpStorage = storage.NewBucket(awscfg, tmpStorageParams, storage.WithLogger(params.logger))

	// Amazon Cognitoの設定
	userAuthParams := &cognito.Params{
		UserPoolID:  a.CognitoUserPoolID,
		AppClientID: a.CognitoUserClientID,
	}
	params.userAuth = cognito.NewClient(awscfg, userAuthParams)

	// Amazon SQSの設定
	sqsParams := &sqs.Params{
		QueueURL: a.SQSQueueURL,
	}
	params.producer = sqs.NewProducer(awscfg, sqsParams, sqs.WithDryRun(a.SQSMockEnabled))

	// Amazon DynamoDBの設定
	dbParams := &dynamodb.Params{
		TablePrefix: "furumaru",
		TableSuffix: a.Environment,
	}
	params.cache = dynamodb.NewClient(awscfg, dbParams, dynamodb.WithLogger(params.logger))

	// New Relicの設定
	if params.newRelicLicense != "" {
		labels := map[string]string{
			"app":     "furumaru",
			"env":     a.Environment,
			"service": params.serviceName,
			"type":    "backend",
		}
		newrelicApp, err := newrelic.NewApplication(
			newrelic.ConfigAppName(params.serviceName),
			newrelic.ConfigLicense(params.newRelicLicense),
			newrelic.ConfigAppLogMetricsEnabled(true),
			newrelic.ConfigAppLogForwardingEnabled(true),
			newrelic.ConfigCustomInsightsEventsEnabled(true),
			newrelic.ConfigAppLogEnabled(true),
			newrelic.ConfigAppLogForwardingEnabled(true),
			func(cfg *newrelic.Config) {
				cfg.HostDisplayName = params.serviceName
				cfg.Labels = labels
			},
		)
		if err != nil {
			return fmt.Errorf("cmd: failed to create newrelic client: %w", err)
		}
		params.newRelic = newrelicApp
	}

	// Sentryの設定
	if params.sentryDsn != "" {
		sentryApp, err := sentry.NewClient(
			sentry.WithServerName(a.AppName),
			sentry.WithEnvironment(a.Environment),
			sentry.WithDSN(params.sentryDsn),
			sentry.WithBind(true),
			sentry.WithTracesSampleRate(a.TraceSampleRate),
			sentry.WithProfilesSampleRate(a.ProfileSampleRate),
		)
		if err != nil {
			return fmt.Errorf("cmd: failed to create sentry client: %w", err)
		}
		params.sentry = sentryApp
	}

	// Slackの設定
	if params.slackToken != "" {
		slackParams := &slack.Params{
			Token:     params.slackToken,
			ChannelID: params.slackChannelID,
		}
		params.slack = slack.NewClient(slackParams, slack.WithLogger(params.logger))
	}

	// KOMOJUの設定
	captureMode := komoju.CaptureModeManual
	if a.CheckoutAutoCaptured {
		captureMode = komoju.CaptureModeAuto
	}
	kpaymentParams := &kpayment.Params{
		Host:         a.KomojuHost,
		ClientID:     params.komojuClientID,
		ClientSecret: params.komojuClientPassword,
	}
	ksessionParams := &ksession.Params{
		Host:         a.KomojuHost,
		ClientID:     params.komojuClientID,
		ClientSecret: params.komojuClientPassword,
		CaptureMode:  captureMode,
	}
	komojuOpts := []komoju.Option{
		komoju.WithLogger(params.logger),
		komoju.WithDebugMode(params.debugMode),
	}
	komojuParams := &komoju.Params{
		Payment: kpayment.NewClient(&http.Client{}, kpaymentParams, komojuOpts...),
		Session: ksession.NewClient(&http.Client{}, ksessionParams, komojuOpts...),
	}
	params.komoju = komoju.NewKomoju(komojuParams)

	// PostalCodeの設定
	params.postalCode = postalcode.NewClient(&http.Client{}, postalcode.WithLogger(params.logger))

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

	// Serviceの設定
	mediaService, err := a.newMediaService(params)
	if err != nil {
		return fmt.Errorf("cmd: failed to create media service: %w", err)
	}
	messengerService, err := a.newMessengerService(params)
	if err != nil {
		return fmt.Errorf("cmd: failed to create messenger service: %w", err)
	}
	userService, err := a.newUserService(params, mediaService, messengerService)
	if err != nil {
		return fmt.Errorf("cmd: failed to create user service: %w", err)
	}
	storeService, err := a.newStoreService(params, userService, mediaService, messengerService)
	if err != nil {
		return fmt.Errorf("cmd: failed to create store service: %w", err)
	}

	// Handlerの設定
	v1Params := &v1.Params{
		WaitGroup: params.waitGroup,
		User:      userService,
		Store:     storeService,
		Messenger: messengerService,
		Media:     mediaService,
	}
	a.v1 = v1.NewHandler(v1Params,
		v1.WithEnvironment(a.Environment),
		v1.WithLogger(params.logger),
		v1.WithSentry(params.sentry),
	)
	a.logger = params.logger
	a.debugMode = params.debugMode
	a.waitGroup = params.waitGroup
	a.slack = params.slack
	a.newRelic = params.newRelic
	return nil
}

func (a *app) getSecret(ctx context.Context, p *params) error {
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// データベース（MySQL）認証情報の取得
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
		// Slack認証情報の取得
		if a.SlackSecretName == "" {
			p.slackToken = a.SlackAPIToken
			p.slackChannelID = a.SlackChannelID
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.SlackSecretName)
		if err != nil {
			return err
		}
		p.slackToken = secrets["token"]
		p.slackChannelID = secrets["channelId"]
		return nil
	})
	eg.Go(func() error {
		// New Relic認証情報の取得
		if a.NewRelicSecretName == "" {
			p.newRelicLicense = a.NewRelicLicense
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.NewRelicSecretName)
		if err != nil {
			return err
		}
		p.newRelicLicense = secrets["license"]
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
	eg.Go(func() error {
		// KOMOJU接続情報の取得
		if a.KomojuSecretName == "" {
			p.komojuClientID = a.KomojuClientID
			p.komojuClientPassword = a.KomojuClientPassword
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.KomojuSecretName)
		if err != nil {
			return err
		}
		p.komojuClientID = secrets["clientId"]
		p.komojuClientPassword = secrets["clientPassword"]
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
		mysql.WithLogger(p.logger), // TODO: 動作検証が終わり次第削除
	)
	if err != nil {
		return nil, err
	}
	if err := cli.DB.Use(telemetry.NewNrTracer(dbname, p.tidbHost, string(newrelic.DatastoreMySQL))); err != nil {
		return nil, err
	}
	return cli, nil
}

func (a *app) newMediaService(p *params) (media.Service, error) {
	mysql, err := a.newDatabase("media", p)
	if err != nil {
		return nil, err
	}
	params := &mediasrv.Params{
		WaitGroup: p.waitGroup,
		Database:  mediadb.NewDatabase(mysql),
		Cache:     p.cache,
		Storage:   p.storage,
		Tmp:       p.tmpStorage,
	}
	return mediasrv.NewService(params, mediasrv.WithLogger(p.logger))
}

func (a *app) newMessengerService(p *params) (messenger.Service, error) {
	db, err := a.newDatabase("messengers", p)
	if err != nil {
		return nil, err
	}
	user, err := a.newUserService(p, nil, nil)
	if err != nil {
		return nil, err
	}
	store, err := a.newStoreService(p, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	params := &messengersrv.Params{
		WaitGroup:   p.waitGroup,
		Producer:    p.producer,
		AdminWebURL: p.adminWebURL,
		UserWebURL:  p.userWebURL,
		Database:    messengerdb.NewDatabase(db),
		User:        user,
		Store:       store,
	}
	return messengersrv.NewService(params, messengersrv.WithLogger(p.logger)), nil
}

func (a *app) newUserService(p *params, media media.Service, messenger messenger.Service) (user.Service, error) {
	mysql, err := a.newTiDB("users", p)
	if err != nil {
		return nil, err
	}
	params := &usersrv.Params{
		WaitGroup: p.waitGroup,
		Database:  userdb.NewDatabase(mysql),
		UserAuth:  p.userAuth,
		Messenger: messenger,
		Media:     media,
	}
	return usersrv.NewService(params, usersrv.WithLogger(p.logger)), nil
}

func (a *app) newStoreService(
	p *params, user user.Service, media media.Service, messenger messenger.Service,
) (store.Service, error) {
	mysql, err := a.newDatabase("stores", p)
	if err != nil {
		return nil, err
	}
	params := &storesrv.Params{
		WaitGroup:  p.waitGroup,
		Database:   storedb.NewDatabase(mysql),
		Cache:      p.cache,
		User:       user,
		Messenger:  messenger,
		Media:      media,
		PostalCode: p.postalCode,
		Komoju:     p.komoju,
	}
	return storesrv.NewService(params, storesrv.WithLogger(p.logger)), nil
}
