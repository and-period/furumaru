package cmd

import (
	"context"
	"net/url"
	"sync"

	v1 "github.com/and-period/furumaru/api/internal/gateway/admin/v1/handler"
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
	"github.com/and-period/furumaru/api/pkg/line"
	"github.com/and-period/furumaru/api/pkg/rbac"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/and-period/furumaru/api/pkg/sqs"
	"github.com/and-period/furumaru/api/pkg/storage"
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
	waitGroup *sync.WaitGroup
	line      line.Client
	newRelic  *newrelic.Application
	v1        v1.Handler
}

type params struct {
	config          *config
	logger          *zap.Logger
	waitGroup       *sync.WaitGroup
	enforcer        rbac.Enforcer
	aws             aws.Config
	secret          secret.Client
	storage         storage.Bucket
	adminAuth       cognito.Client
	userAuth        cognito.Client
	producer        sqs.Producer
	line            line.Client
	newRelic        *newrelic.Application
	adminWebURL     *url.URL
	userWebURL      *url.URL
	dbHost          string
	dbPort          string
	dbUsername      string
	dbPassword      string
	lineToken       string
	lineSecret      string
	lineRoomID      string
	newRelicLicense string
}

//nolint:funlen
func newRegistry(ctx context.Context, conf *config, logger *zap.Logger) (*registry, error) {
	params := &params{
		config:    conf,
		logger:    logger,
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
	params.userAuth = cognito.NewClient(awscfg, userAuthParams)

	// Amazon SQSの設定
	sqsParams := &sqs.Params{
		QueueURL: conf.SQSQueueURL,
	}
	params.producer = sqs.NewProducer(awscfg, sqsParams, sqs.WithDryRun(conf.SQSMockEnabled))

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

	// LINEの設定
	if params.lineToken != "" {
		lineParams := &line.Params{
			Token:  params.lineToken,
			Secret: params.lineSecret,
			RoomID: params.lineRoomID,
		}
		linebot, err := line.NewClient(lineParams, line.WithLogger(logger))
		if err != nil {
			return nil, err
		}
		params.line = linebot
	}

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
	messengerService, err := newMessengerService(params)
	if err != nil {
		return nil, err
	}
	userService, err := newUserService(params, messengerService)
	if err != nil {
		return nil, err
	}
	storeService, err := newStoreService(params, userService, messengerService)
	if err != nil {
		return nil, err
	}

	// Handlerの設定
	v1Params := &v1.Params{
		WaitGroup: params.waitGroup,
		Enforcer:  enforcer,
		Storage:   storage.NewBucket(awscfg, storageParams),
		User:      userService,
		Store:     storeService,
		Messenger: messengerService,
	}
	return &registry{
		appName:   conf.AppName,
		env:       conf.Environment,
		waitGroup: params.waitGroup,
		line:      params.line,
		newRelic:  params.newRelic,
		v1:        v1.NewHandler(v1Params, v1.WithLogger(logger)),
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
		// LINE認証情報の取得
		if p.config.LINESecretName == "" {
			p.lineToken = p.config.LINEChannelToken
			p.lineSecret = p.config.LINEChannelSecret
			p.lineRoomID = p.config.LINERoomID
			return nil
		}
		secrets, err := p.secret.Get(ectx, p.config.LINESecretName)
		if err != nil {
			return err
		}
		p.lineToken = secrets["token"]
		p.lineSecret = secrets["secret"]
		p.lineRoomID = secrets["roomId"]
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
	cli, err := database.NewClient(
		params,
		database.WithLogger(p.logger),
		database.WithTLS(p.config.DBEnabledTLS),
		database.WithTimeZone(p.config.DBTimeZone),
	)
	if err != nil {
		return nil, err
	}
	if err := cli.DB.Use(telemetry.NewNrTracer(dbname, p.dbHost, string(newrelic.DatastoreMySQL))); err != nil {
		return nil, err
	}
	return cli, nil
}

func newMessengerService(p *params) (messenger.Service, error) {
	mysql, err := newDatabase("messengers", p)
	if err != nil {
		return nil, err
	}
	dbParams := &messengerdb.Params{
		Database: mysql,
	}
	user, err := newUserService(p, nil)
	if err != nil {
		return nil, err
	}
	params := &messengersrv.Params{
		WaitGroup:   p.waitGroup,
		Producer:    p.producer,
		AdminWebURL: p.adminWebURL,
		UserWebURL:  p.userWebURL,
		Database:    messengerdb.NewDatabase(dbParams),
		User:        user,
	}
	return messengersrv.NewService(params, messengersrv.WithLogger(p.logger)), nil
}

func newUserService(p *params, messenger messenger.Service) (user.Service, error) {
	mysql, err := newDatabase("users", p)
	if err != nil {
		return nil, err
	}
	dbParams := &userdb.Params{
		Database: mysql,
	}
	params := &usersrv.Params{
		WaitGroup: p.waitGroup,
		Database:  userdb.NewDatabase(dbParams),
		AdminAuth: p.adminAuth,
		UserAuth:  p.userAuth,
		Messenger: messenger,
	}
	return usersrv.NewService(params, usersrv.WithLogger(p.logger)), nil
}

func newStoreService(p *params, user user.Service, messenger messenger.Service) (store.Service, error) {
	mysql, err := newDatabase("stores", p)
	if err != nil {
		return nil, err
	}
	dbParams := &storedb.Params{
		Database: mysql,
	}
	params := &storesrv.Params{
		WaitGroup: p.waitGroup,
		Database:  storedb.NewDatabase(dbParams),
		User:      user,
		Messenger: messenger,
	}
	return storesrv.NewService(params, storesrv.WithLogger(p.logger)), nil
}
