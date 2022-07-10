package cmd

import (
	"context"
	"sync"

	v1 "github.com/and-period/furumaru/api/internal/gateway/admin/v1/handler"
	"github.com/and-period/furumaru/api/internal/messenger"
	messengersrv "github.com/and-period/furumaru/api/internal/messenger/service"
	"github.com/and-period/furumaru/api/internal/store"
	storedb "github.com/and-period/furumaru/api/internal/store/database"
	storesrv "github.com/and-period/furumaru/api/internal/store/service"
	"github.com/and-period/furumaru/api/internal/user"
	userdb "github.com/and-period/furumaru/api/internal/user/database"
	usersrv "github.com/and-period/furumaru/api/internal/user/service"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/rbac"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/and-period/furumaru/api/pkg/sqs"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"
)

type registry struct {
	waitGroup *sync.WaitGroup
	v1        v1.Handler
}

type params struct {
	config     *config
	logger     *zap.Logger
	waitGroup  *sync.WaitGroup
	enforcer   rbac.Enforcer
	aws        aws.Config
	secret     secret.Client
	storage    storage.Bucket
	adminAuth  cognito.Client
	userAuth   cognito.Client
	producer   sqs.Producer
	dbHost     string
	dbPort     string
	dbUsername string
	dbPassword string
}

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

	// Serviceの設定
	messengerService := newMessengerService(ctx, params)
	userService, err := newUserService(ctx, params, messengerService)
	if err != nil {
		return nil, err
	}
	storeService, err := newStoreService(ctx, params, userService, messengerService)
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
	}
	return &registry{
		waitGroup: params.waitGroup,
		v1:        v1.NewHandler(v1Params, v1.WithLogger(logger)),
	}, nil
}

func getSecret(ctx context.Context, p *params) error {
	// データベース認証情報の取得
	if p.config.DBSecretName == "" {
		p.dbHost = p.config.DBHost
		p.dbPort = p.config.DBPort
		p.dbUsername = p.config.DBUsername
		p.dbPassword = p.config.DBPassword
	} else {
		secrets, err := p.secret.Get(ctx, p.config.DBSecretName)
		if err != nil {
			return err
		}
		p.dbHost = secrets["host"]
		p.dbPort = secrets["port"]
		p.dbUsername = secrets["username"]
		p.dbPassword = secrets["password"]
	}
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

func newMessengerService(ctx context.Context, p *params) messenger.Service {
	params := &messengersrv.Params{
		WaitGroup: p.waitGroup,
		Producer:  p.producer,
	}
	return messengersrv.NewService(params, messengersrv.WithLogger(p.logger))
}

func newUserService(ctx context.Context, p *params, messenger messenger.Service) (user.Service, error) {
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

func newStoreService(
	ctx context.Context, p *params, user user.Service, messenger messenger.Service,
) (store.Service, error) {
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
