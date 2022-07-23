package cmd

import (
	"context"
	"net/url"
	"sync"

	"github.com/and-period/furumaru/api/internal/messenger"
	messengerdb "github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/scheduler"
	messengersrv "github.com/and-period/furumaru/api/internal/messenger/service"
	"github.com/and-period/furumaru/api/internal/user"
	userdb "github.com/and-period/furumaru/api/internal/user/database"
	usersrv "github.com/and-period/furumaru/api/internal/user/service"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/and-period/furumaru/api/pkg/sqs"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"
)

type registry struct {
	appName   string
	env       string
	waitGroup *sync.WaitGroup
	job       scheduler.Scheduler
}

type params struct {
	config      *config
	logger      *zap.Logger
	waitGroup   *sync.WaitGroup
	aws         aws.Config
	secret      secret.Client
	producer    sqs.Producer
	adminWebURL *url.URL
	userWebURL  *url.URL
	dbHost      string
	dbPort      string
	dbUsername  string
	dbPassword  string
}

func newRegistry(ctx context.Context, conf *config, logger *zap.Logger) (*registry, error) {
	params := &params{
		config:    conf,
		logger:    logger,
		waitGroup: &sync.WaitGroup{},
	}

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

	// Amazon SQSの設定
	sqsParams := &sqs.Params{
		QueueURL: conf.SQSQueueURL,
	}
	params.producer = sqs.NewProducer(awscfg, sqsParams, sqs.WithDryRun(conf.SQSMockEnabled))

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

	// Databaseの設定
	dbClient, err := newDatabase("messengers", params)
	if err != nil {
		return nil, err
	}

	// Serviceの設定
	messengerService, err := newMessengerService(ctx, params)
	if err != nil {
		return nil, err
	}

	// Jobの設定
	dbParams := &messengerdb.Params{
		Database: dbClient,
	}
	jobParams := &scheduler.Params{
		WaitGroup: params.waitGroup,
		Database:  messengerdb.NewDatabase(dbParams),
		Messenger: messengerService,
	}
	return &registry{
		appName:   conf.AppName,
		env:       conf.Environment,
		waitGroup: params.waitGroup,
		job:       scheduler.NewScheduler(jobParams, scheduler.WithLogger(logger)),
	}, nil
}

func getSecret(ctx context.Context, p *params) error {
	// データベース認証情報の取得
	if p.config.DBSecretName == "" {
		p.dbHost = p.config.DBHost
		p.dbPort = p.config.DBPort
		p.dbUsername = p.config.DBUsername
		p.dbPassword = p.config.DBPassword
		return nil
	}
	secrets, err := p.secret.Get(ctx, p.config.DBSecretName)
	if err != nil {
		return err
	}
	p.dbHost = secrets["host"]
	p.dbPort = secrets["port"]
	p.dbUsername = secrets["username"]
	p.dbPassword = secrets["password"]
	return nil
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
	return database.NewClient(
		params,
		database.WithLogger(p.logger),
		database.WithTLS(p.config.DBEnabledTLS),
		database.WithTimeZone(p.config.DBTimeZone),
	)
}

func newMessengerService(ctx context.Context, p *params) (messenger.Service, error) {
	mysql, err := newDatabase("messengers", p)
	if err != nil {
		return nil, err
	}
	dbParams := &messengerdb.Params{
		Database: mysql,
	}
	user, err := newUserService(ctx, p)
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

func newUserService(ctx context.Context, p *params) (user.Service, error) {
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
	}
	return usersrv.NewService(params, usersrv.WithLogger(p.logger)), nil
}
