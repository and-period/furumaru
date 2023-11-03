package admin

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"

	v1 "github.com/and-period/furumaru/api/internal/gateway/admin/v1/handler"
	khandler "github.com/and-period/furumaru/api/internal/gateway/komoju/handler"
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
	userdb "github.com/and-period/furumaru/api/internal/user/database/mysql"
	usersrv "github.com/and-period/furumaru/api/internal/user/service"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/and-period/furumaru/api/pkg/postalcode"
	"github.com/and-period/furumaru/api/pkg/rbac"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/and-period/furumaru/api/pkg/slack"
	"github.com/and-period/furumaru/api/pkg/sqs"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	sentrygo "github.com/getsentry/sentry-go"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rafaelhl/gorm-newrelic-telemetry-plugin/telemetry"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type params struct {
	logger               *zap.Logger
	waitGroup            *sync.WaitGroup
	enforcer             rbac.Enforcer
	aws                  aws.Config
	secret               secret.Client
	storage              storage.Bucket
	tmpStorage           storage.Bucket
	adminAuth            cognito.Client
	userAuth             cognito.Client
	messengerQueue       sqs.Producer
	mediaQueue           sqs.Producer
	slack                slack.Client
	newRelic             *newrelic.Application
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
		logger:    zap.NewNop(),
		now:       jst.Now,
		waitGroup: &sync.WaitGroup{},
		debugMode: a.LogLevel == "debug",
	}

	// Casbinの設定
	enforcer, err := rbac.NewEnforcer(a.RBACModelPath, a.RBACPolicyPath)
	if err != nil {
		return err
	}
	params.enforcer = enforcer

	// AWS SDKの設定
	awscfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(a.AWSRegion))
	if err != nil {
		return err
	}
	params.aws = awscfg

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
	adminAuthParams := &cognito.Params{
		UserPoolID:  a.CognitoAdminPoolID,
		AppClientID: a.CognitoAdminClientID,
	}
	params.adminAuth = cognito.NewClient(awscfg, adminAuthParams)
	userAuthParams := &cognito.Params{
		UserPoolID:  a.CognitoUserPoolID,
		AppClientID: a.CognitoUserClientID,
	}
	params.userAuth = cognito.NewClient(awscfg, userAuthParams, cognito.WithLogger(params.logger))

	// Amazon SQSの設定
	messengerSQSParams := &sqs.Params{
		QueueURL: a.SQSMessengerQueueURL,
	}
	params.messengerQueue = sqs.NewProducer(awscfg, messengerSQSParams, sqs.WithDryRun(a.SQSMockEnabled))
	mediaSQSParams := &sqs.Params{
		QueueURL: a.SQSMediaQueueURL,
	}
	params.mediaQueue = sqs.NewProducer(awscfg, mediaSQSParams, sqs.WithDryRun(a.SQSMockEnabled))

	// New Relicの設定
	if params.newRelicLicense != "" {
		newrelicApp, err := newrelic.NewApplication(
			newrelic.ConfigAppName(a.AppName),
			newrelic.ConfigLicense(params.newRelicLicense),
			newrelic.ConfigAppLogForwardingEnabled(true),
		)
		if err != nil {
			return err
		}
		params.newRelic = newrelicApp
	}

	// Sentryの設定
	if params.sentryDsn != "" {
		sentryOptions := sentrygo.ClientOptions{
			Dsn:   params.sentryDsn,
			Debug: params.debugMode,
		}
		if err := sentrygo.Init(sentryOptions); err != nil {
			return err
		}
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
	kpaymentParams := &kpayment.Params{
		Host:         a.KomojuHost,
		ClientID:     params.komojuClientID,
		ClientSecret: params.komojuClientPassword,
	}
	ksessionParams := &ksession.Params{
		Host:         a.KomojuHost,
		ClientID:     params.komojuClientID,
		ClientSecret: params.komojuClientPassword,
	}
	komojuOpts := []komoju.Option{
		komoju.WithLogger(logger),
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
		return err
	}
	params.adminWebURL = adminWebURL
	userWebURL, err := url.Parse(a.UserWebURL)
	if err != nil {
		return err
	}
	params.userWebURL = userWebURL

	// Serviceの設定
	mediaService, err := a.newMediaService(params)
	if err != nil {
		return err
	}
	messengerService, err := a.newMessengerService(params)
	if err != nil {
		return err
	}
	userService, err := a.newUserService(params, mediaService, messengerService)
	if err != nil {
		return err
	}
	storeService, err := a.newStoreService(params, userService, mediaService, messengerService)
	if err != nil {
		return err
	}

	// Handlerの設定
	v1Params := &v1.Params{
		WaitGroup: params.waitGroup,
		Enforcer:  enforcer,
		User:      userService,
		Store:     storeService,
		Messenger: messengerService,
		Media:     mediaService,
	}
	khandlerParams := &khandler.Params{
		WaitGroup: params.waitGroup,
	}
	a.logger = params.logger
	a.v1 = v1.NewHandler(v1Params, v1.WithLogger(params.logger))
	a.komoju = khandler.NewHandler(khandlerParams, khandler.WithLogger(params.logger))
	a.debugMode = params.debugMode
	a.waitGroup = params.waitGroup
	a.slack = params.slack
	a.newRelic = params.newRelic
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
		mysql.WithLogger(p.logger),
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

func (a *app) newMediaService(p *params) (media.Service, error) {
	mysql, err := a.newDatabase("media", p)
	if err != nil {
		return nil, err
	}
	params := &mediasrv.Params{
		WaitGroup: p.waitGroup,
		Database:  mediadb.NewDatabase(mysql),
		Storage:   p.storage,
		Tmp:       p.tmpStorage,
		Producer:  p.mediaQueue,
	}
	return mediasrv.NewService(params, mediasrv.WithLogger(p.logger))
}

func (a *app) newMessengerService(p *params) (messenger.Service, error) {
	mysql, err := a.newDatabase("messengers", p)
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
		Producer:    p.messengerQueue,
		AdminWebURL: p.adminWebURL,
		UserWebURL:  p.userWebURL,
		Database:    messengerdb.NewDatabase(mysql),
		User:        user,
		Store:       store,
	}
	return messengersrv.NewService(params, messengersrv.WithLogger(p.logger)), nil
}

func (a *app) newUserService(p *params, media media.Service, messenger messenger.Service) (user.Service, error) {
	mysql, err := a.newDatabase("users", p)
	if err != nil {
		return nil, err
	}
	store, err := a.newStoreService(p, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	params := &usersrv.Params{
		WaitGroup: p.waitGroup,
		Database:  userdb.NewDatabase(mysql),
		AdminAuth: p.adminAuth,
		UserAuth:  p.userAuth,
		Store:     store,
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
		User:       user,
		Messenger:  messenger,
		Media:      media,
		PostalCode: p.postalCode,
		Komoju:     p.komoju,
	}
	return storesrv.NewService(params, storesrv.WithLogger(p.logger)), nil
}
