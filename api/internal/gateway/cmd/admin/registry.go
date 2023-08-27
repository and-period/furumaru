package cmd

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"

	v1 "github.com/and-period/furumaru/api/internal/gateway/admin/v1/handler"
	shandler "github.com/and-period/furumaru/api/internal/gateway/stripe/handler"
	"github.com/and-period/furumaru/api/internal/media"
	mediadb "github.com/and-period/furumaru/api/internal/media/database"
	mediasrv "github.com/and-period/furumaru/api/internal/media/service"
	"github.com/and-period/furumaru/api/internal/messenger"
	messengerdb "github.com/and-period/furumaru/api/internal/messenger/database"
	messengersrv "github.com/and-period/furumaru/api/internal/messenger/service"
	"github.com/and-period/furumaru/api/internal/store"
	storedb "github.com/and-period/furumaru/api/internal/store/database"
	storesrv "github.com/and-period/furumaru/api/internal/store/service"
	"github.com/and-period/furumaru/api/internal/user"
	userdb "github.com/and-period/furumaru/api/internal/user/database"
	usersrv "github.com/and-period/furumaru/api/internal/user/service"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/postalcode"
	"github.com/and-period/furumaru/api/pkg/rbac"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/and-period/furumaru/api/pkg/slack"
	"github.com/and-period/furumaru/api/pkg/sqs"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/and-period/furumaru/api/pkg/stripe"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rafaelhl/gorm-newrelic-telemetry-plugin/telemetry"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type registry struct {
	appName   string
	env       string
	debugMode bool
	waitGroup *sync.WaitGroup
	slack     slack.Client
	newRelic  *newrelic.Application
	v1        v1.Handler
	stripe    shandler.Handler
}

type params struct {
	config           *config
	logger           *zap.Logger
	waitGroup        *sync.WaitGroup
	enforcer         rbac.Enforcer
	aws              aws.Config
	secret           secret.Client
	storage          storage.Bucket
	tmpStorage       storage.Bucket
	adminAuth        cognito.Client
	userAuth         cognito.Client
	dynamodb         dynamodb.Client
	messengerQueue   sqs.Producer
	mediaQueue       sqs.Producer
	slack            slack.Client
	newRelic         *newrelic.Application
	receiver         stripe.Receiver
	adminWebURL      *url.URL
	userWebURL       *url.URL
	postalCode       postalcode.Client
	now              func() time.Time
	dbHost           string
	dbPort           string
	dbUsername       string
	dbPassword       string
	slackToken       string
	slackChannelID   string
	newRelicLicense  string
	stripeSecretKey  string
	stripeWebhookKey string
}

//nolint:funlen
func newRegistry(ctx context.Context, conf *config, logger *zap.Logger) (*registry, error) {
	params := &params{
		config:    conf,
		logger:    logger,
		now:       jst.Now,
		waitGroup: &sync.WaitGroup{},
	}

	// Casbinの設定
	enforcer, err := rbac.NewEnforcer(conf.RBACModelPath, conf.RBACPolicyPath)
	if err != nil {
		return nil, err
	}
	params.enforcer = enforcer

	// AWS SDKの設定
	awscfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(conf.AWSRegion))
	if err != nil {
		return nil, err
	}
	params.aws = awscfg

	// AWS Secrets Managerの設定
	params.secret = secret.NewClient(awscfg)
	if err := getSecret(ctx, params); err != nil {
		return nil, err
	}

	// Amazon S3の設定
	storageParams := &storage.Params{
		Bucket: conf.S3Bucket,
	}
	params.storage = storage.NewBucket(awscfg, storageParams)
	tmpStorageParams := &storage.Params{
		Bucket: conf.S3TmpBucket,
	}
	params.tmpStorage = storage.NewBucket(awscfg, tmpStorageParams, storage.WithLogger(params.logger))

	// Amazon Cognitoの設定
	adminAuthParams := &cognito.Params{
		UserPoolID:  conf.CognitoAdminPoolID,
		AppClientID: conf.CognitoAdminClientID,
	}
	params.adminAuth = cognito.NewClient(awscfg, adminAuthParams)
	userAuthParams := &cognito.Params{
		UserPoolID:  conf.CognitoUserPoolID,
		AppClientID: conf.CognitoUserClientID,
	}
	params.userAuth = cognito.NewClient(awscfg, userAuthParams, cognito.WithLogger(params.logger))

	// Amazon SQSの設定
	messengerSQSParams := &sqs.Params{
		QueueURL: conf.SQSMessengerQueueURL,
	}
	params.messengerQueue = sqs.NewProducer(awscfg, messengerSQSParams, sqs.WithDryRun(conf.SQSMockEnabled))
	mediaSQSParams := &sqs.Params{
		QueueURL: conf.SQSMediaQueueURL,
	}
	params.mediaQueue = sqs.NewProducer(awscfg, mediaSQSParams, sqs.WithDryRun(conf.SQSMockEnabled))

	// Amazon DynamoDBの設定
	dynamodbParams := &dynamodb.Params{
		TablePrefix: params.config.Environment,
	}
	params.dynamodb = dynamodb.NewClient(awscfg, dynamodbParams, dynamodb.WithLogger(params.logger))

	// New Relicの設定
	if params.newRelicLicense != "" {
		newrelicApp, err := newrelic.NewApplication(
			newrelic.ConfigAppName(conf.AppName),
			newrelic.ConfigLicense(params.newRelicLicense),
			newrelic.ConfigAppLogForwardingEnabled(true),
		)
		if err != nil {
			return nil, err
		}
		params.newRelic = newrelicApp
	}

	// Stripeの設定
	stripeParams := &stripe.Params{
		SecretKey:  params.stripeSecretKey,
		WebhookKey: params.stripeWebhookKey,
	}
	params.receiver = stripe.NewReceiver(stripeParams, stripe.WithLogger(logger))

	// Slackの設定
	if params.slackToken != "" {
		slackParams := &slack.Params{
			Token:     params.slackToken,
			ChannelID: params.slackChannelID,
		}
		params.slack = slack.NewClient(slackParams, slack.WithLogger(logger))
	}

	// PostalCodeの設定
	params.postalCode = postalcode.NewClient(&http.Client{}, postalcode.WithLogger(logger))

	// WebURLの設定
	adminWebURL, err := url.Parse(conf.AminWebURL)
	if err != nil {
		return nil, err
	}
	params.adminWebURL = adminWebURL
	userWebURL, err := url.Parse(conf.UserWebURL)
	if err != nil {
		return nil, err
	}
	params.userWebURL = userWebURL

	// Serviceの設定
	mediaService, err := newMediaService(params)
	if err != nil {
		return nil, err
	}
	messengerService, err := newMessengerService(params)
	if err != nil {
		return nil, err
	}
	userService, err := newUserService(params, mediaService, messengerService)
	if err != nil {
		return nil, err
	}
	storeService, err := newStoreService(params, userService, mediaService, messengerService)
	if err != nil {
		return nil, err
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
	shandlerParams := &shandler.Params{
		WaitGroup: params.waitGroup,
		Receiver:  params.receiver,
	}
	return &registry{
		appName:   conf.AppName,
		env:       conf.Environment,
		debugMode: conf.LogLevel == "debug",
		waitGroup: params.waitGroup,
		slack:     params.slack,
		newRelic:  params.newRelic,
		v1:        v1.NewHandler(v1Params, v1.WithLogger(logger)),
		stripe:    shandler.NewHandler(shandlerParams, shandler.WithLogger(logger)),
	}, nil
}

func getSecret(ctx context.Context, p *params) error {
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// データベース認証情報の取得
		if p.config.DBSecretName == "" {
			p.dbHost = p.config.DBHost
			p.dbPort = p.config.DBPort
			p.dbUsername = p.config.DBUsername
			p.dbPassword = p.config.DBPassword
			return nil
		}
		secrets, err := p.secret.Get(ectx, p.config.DBSecretName)
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
		if p.config.SlackSecretName == "" {
			p.slackToken = p.config.SlackAPIToken
			p.slackChannelID = p.config.SlackChannelID
			return nil
		}
		secrets, err := p.secret.Get(ectx, p.config.SlackSecretName)
		if err != nil {
			return err
		}
		p.slackToken = secrets["token"]
		p.slackChannelID = secrets["channelId"]
		return nil
	})
	eg.Go(func() error {
		// New Relic認証情報の取得
		if p.config.NewRelicSecretName == "" {
			p.newRelicLicense = p.config.NewRelicLicense
			return nil
		}
		secrets, err := p.secret.Get(ectx, p.config.NewRelicSecretName)
		if err != nil {
			return err
		}
		p.newRelicLicense = secrets["license"]
		return nil
	})
	eg.Go(func() error {
		// Stripe認証情報の取得
		if p.config.StripeSecretName == "" {
			p.stripeSecretKey = p.config.StripeSecretKey
			p.stripeWebhookKey = p.config.StripeWebhookKey
			return nil
		}
		secrets, err := p.secret.Get(ectx, p.config.StripeSecretName)
		if err != nil {
			return err
		}
		p.stripeSecretKey = secrets["secretKey"]
		p.stripeWebhookKey = secrets["webhookKey"]
		return nil
	})
	return eg.Wait()
}

func newDatabase(dbname string, p *params) (*database.Client, error) {
	params := &database.Params{
		Socket:   p.config.DBSocket,
		Host:     p.dbHost,
		Port:     p.dbPort,
		Database: dbname,
		Username: p.dbUsername,
		Password: p.dbPassword,
	}
	location, err := time.LoadLocation(p.config.DBTimeZone)
	if err != nil {
		return nil, err
	}
	cli, err := database.NewClient(
		params,
		database.WithLogger(p.logger),
		database.WithNow(p.now),
		database.WithTLS(p.config.DBEnabledTLS),
		database.WithLocation(location),
	)
	if err != nil {
		return nil, err
	}
	if err := cli.DB.Use(telemetry.NewNrTracer(dbname, p.dbHost, string(newrelic.DatastoreMySQL))); err != nil {
		return nil, err
	}
	return cli, nil
}

func newMediaService(p *params) (media.Service, error) {
	mysql, err := newDatabase("media", p)
	if err != nil {
		return nil, err
	}
	dbParams := &mediadb.Params{
		Database: mysql,
	}
	params := &mediasrv.Params{
		WaitGroup: p.waitGroup,
		Database:  mediadb.NewDatabase(dbParams),
		Storage:   p.storage,
		Tmp:       p.tmpStorage,
		Producer:  p.mediaQueue,
	}
	return mediasrv.NewService(params, mediasrv.WithLogger(p.logger))
}

func newMessengerService(p *params) (messenger.Service, error) {
	mysql, err := newDatabase("messengers", p)
	if err != nil {
		return nil, err
	}
	dbParams := &messengerdb.Params{
		Database: mysql,
	}
	user, err := newUserService(p, nil, nil)
	if err != nil {
		return nil, err
	}
	store, err := newStoreService(p, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	params := &messengersrv.Params{
		WaitGroup:   p.waitGroup,
		Producer:    p.messengerQueue,
		AdminWebURL: p.adminWebURL,
		UserWebURL:  p.userWebURL,
		Database:    messengerdb.NewDatabase(dbParams),
		User:        user,
		Store:       store,
	}
	return messengersrv.NewService(params, messengersrv.WithLogger(p.logger)), nil
}

func newUserService(p *params, media media.Service, messenger messenger.Service) (user.Service, error) {
	mysql, err := newDatabase("users", p)
	if err != nil {
		return nil, err
	}
	dbParams := &userdb.Params{
		Database: mysql,
	}
	store, err := newStoreService(p, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	params := &usersrv.Params{
		WaitGroup: p.waitGroup,
		Database:  userdb.NewDatabase(dbParams),
		AdminAuth: p.adminAuth,
		UserAuth:  p.userAuth,
		Store:     store,
		Messenger: messenger,
		Media:     media,
	}
	return usersrv.NewService(params, usersrv.WithLogger(p.logger)), nil
}

func newStoreService(
	p *params, user user.Service, media media.Service, messenger messenger.Service,
) (store.Service, error) {
	mysql, err := newDatabase("stores", p)
	if err != nil {
		return nil, err
	}
	dbParams := &storedb.Params{
		Database: mysql,
		DynamoDB: p.dynamodb,
	}
	params := &storesrv.Params{
		WaitGroup:  p.waitGroup,
		Database:   storedb.NewDatabase(dbParams),
		User:       user,
		Messenger:  messenger,
		Media:      media,
		PostalCode: p.postalCode,
	}
	return storesrv.NewService(params, storesrv.WithLogger(p.logger)), nil
}
