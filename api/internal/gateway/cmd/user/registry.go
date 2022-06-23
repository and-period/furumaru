package cmd

import (
	"context"
	"net/url"
	"os"
	"sync"

	v1 "github.com/and-period/furumaru/api/internal/gateway/user/v1/handler"
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
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type registry struct {
	waitGroup *sync.WaitGroup
	v1        v1.Handler
}

type params struct {
	config     *config
	logger     *zap.Logger
	waitGroup  *sync.WaitGroup
	aws        aws.Config
	storage    storage.Bucket
	userAuth   cognito.Client
	mailer     mailer.Client
	userWebURL *url.URL
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

	// Amazon S3の設定
	storageParams := &storage.Params{
		Bucket: conf.S3Bucket,
	}
	params.storage = storage.NewBucket(awscfg, storageParams)

	// Amazon Cognitoの設定
	userAuthParams := &cognito.Params{
		UserPoolID:  conf.CognitoUserPoolID,
		AppClientID: conf.CognitoUserClientID,
	}
	params.userAuth = cognito.NewClient(awscfg, userAuthParams)

	// メールテンプレートの設定
	f, err := os.Open(conf.SendGridTemplatePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var templateMap map[string]string
	d := yaml.NewDecoder(f)
	if err := d.Decode(&templateMap); err != nil {
		return nil, err
	}

	// Mailerの設定
	mailParams := &mailer.Params{
		APIKey:      conf.SendGridAPIKey,
		FromName:    conf.MailFromName,
		FromAddress: conf.MailFromAddress,
		TemplateMap: templateMap,
	}
	params.mailer = mailer.NewClient(mailParams, mailer.WithLogger(logger))

	// WebURLの設定
	userWebURL, err := url.Parse(conf.UserWebURL)
	if err != nil {
		return nil, err
	}
	params.userWebURL = userWebURL

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
		Storage:   params.storage,
		User:      userService,
		Store:     storeService,
	}
	return &registry{
		waitGroup: params.waitGroup,
		v1:        v1.NewHandler(v1Params, v1.WithLogger(logger)),
	}, nil
}

func newDatabase(dbname string, conf *config, logger *zap.Logger) (*database.Client, error) {
	params := &database.Params{
		Socket:   conf.DBSocket,
		Host:     conf.DBHost,
		Port:     conf.DBPort,
		Database: dbname,
		Username: conf.DBUsername,
		Password: conf.DBPassword,
	}
	return database.NewClient(
		params,
		database.WithLogger(logger),
		database.WithTLS(conf.DBEnabledTLS),
		database.WithTimeZone(conf.DBTimeZone),
	)
}

func newMessengerService(ctx context.Context, p *params) messenger.Service {
	params := &messengersrv.Params{
		WaitGroup:   p.waitGroup,
		Mailer:      p.mailer,
		UserWebURL:  p.userWebURL,
		AdminWebURL: &url.URL{},
	}
	return messengersrv.NewService(params, messengersrv.WithLogger(p.logger))
}

func newUserService(ctx context.Context, p *params, messenger messenger.Service) (user.Service, error) {
	mysql, err := newDatabase("users", p.config, p.logger)
	if err != nil {
		return nil, err
	}
	dbParams := &userdb.Params{
		Database: mysql,
	}
	params := &usersrv.Params{
		WaitGroup: p.waitGroup,
		Database:  userdb.NewDatabase(dbParams),
		UserAuth:  p.userAuth,
		Messenger: messenger,
	}
	return usersrv.NewService(params, usersrv.WithLogger(p.logger)), nil
}

func newStoreService(
	ctx context.Context, p *params, user user.Service, messenger messenger.Service,
) (store.Service, error) {
	mysql, err := newDatabase("stores", p.config, p.logger)
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
