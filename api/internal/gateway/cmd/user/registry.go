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
	awscredentials "github.com/aws/aws-sdk-go-v2/credentials"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type registry struct {
	v1        v1.APIV1Handler
	waitGroup *sync.WaitGroup
}

type serviceParams struct {
	waitGroup *sync.WaitGroup
	config    *config
	options   *options
	aws       aws.Config
	messenger messenger.MessengerService
}

type options struct {
	logger *zap.Logger
}

type option func(opts *options)

func withLogger(logger *zap.Logger) option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func newRegistry(ctx context.Context, conf *config, opts ...option) (*registry, error) {
	wg := &sync.WaitGroup{}

	// オプション設定の取得
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}

	// AWS SDKの設定
	awscreds := aws.NewCredentialsCache(
		awscredentials.NewStaticCredentialsProvider(conf.AWSAccessKey, conf.AWSSecretKey, ""),
	)
	awscfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(conf.AWSRegion),
		awsconfig.WithCredentialsProvider(awscreds),
	)
	if err != nil {
		return nil, err
	}

	// Amazon S3の設定
	storageParams := &storage.Params{
		Bucket: conf.S3BucketName,
	}

	// Serviceの設定
	srvParams := &serviceParams{
		waitGroup: wg,
		config:    conf,
		options:   dopts,
		aws:       awscfg,
	}
	messengerService, err := newMessengerService(ctx, srvParams)
	if err != nil {
		return nil, err
	}
	srvParams.messenger = messengerService

	userService, err := newUserService(ctx, srvParams)
	if err != nil {
		return nil, err
	}
	storeService, err := newStoreService(ctx, srvParams)
	if err != nil {
		return nil, err
	}

	// Handlerの設定
	v1Params := &v1.Params{
		WaitGroup:    &sync.WaitGroup{},
		Storage:      storage.NewBucket(awscfg, storageParams),
		UserService:  userService,
		StoreService: storeService,
	}

	return &registry{
		v1: v1.NewAPIV1Handler(v1Params, v1.WithLogger(dopts.logger)),
	}, nil
}

func newDatabase(params *database.Params, tls bool, timezone string, opts *options) (*database.Client, error) {
	return database.NewClient(
		params,
		database.WithLogger(opts.logger),
		database.WithTLS(tls),
		database.WithTimeZone(timezone),
	)
}

func newUserService(ctx context.Context, p *serviceParams) (user.UserService, error) {
	// MySQLの設定
	mysqlParams := &database.Params{
		Socket:   p.config.DBUserSocket,
		Host:     p.config.DBUserHost,
		Port:     p.config.DBUserPort,
		Database: "users",
		Username: p.config.DBUserUsername,
		Password: p.config.DBUserPassword,
	}
	mysql, err := newDatabase(mysqlParams, p.config.DBUserEnabledTLS, p.config.DBUserTimeZone, p.options)
	if err != nil {
		return nil, err
	}

	// Databaseの設定
	dbParams := &userdb.Params{
		Database: mysql,
	}

	// Amazon Cognitoの設定
	userAuthParams := &cognito.Params{
		UserPoolID:      p.config.CognitoUserPoolID,
		AppClientID:     p.config.CognitoUserClientID,
		AppClientSecret: p.config.CognitoUserClientSecret,
	}

	// User Serviceの設定
	params := &usersrv.Params{
		Database:         userdb.NewDatabase(dbParams),
		UserAuth:         cognito.NewClient(p.aws, userAuthParams),
		MessengerService: p.messenger,
		WaitGroup:        p.waitGroup,
	}
	return usersrv.NewUserService(
		params,
		usersrv.WithLogger(p.options.logger),
	), nil
}

func newStoreService(ctx context.Context, p *serviceParams) (store.StoreService, error) {
	// MySQLの設定
	mysqlParams := &database.Params{
		Socket:   p.config.DBStoreSocket,
		Host:     p.config.DBStoreHost,
		Port:     p.config.DBStorePort,
		Database: "stores",
		Username: p.config.DBStoreUsername,
		Password: p.config.DBStorePassword,
	}
	mysql, err := newDatabase(mysqlParams, p.config.DBStoreEnabledTLS, p.config.DBStoreTimeZone, p.options)
	if err != nil {
		return nil, err
	}

	// Databaseの設定
	dbParams := storedb.Params{
		Database: mysql,
	}

	// Store Serviceの設定
	params := &storesrv.Params{
		Database:  storedb.NewDatabase(&dbParams),
		WaitGroup: p.waitGroup,
	}
	return storesrv.NewStoreService(
		params,
		storesrv.WithLogger(p.options.logger),
	), nil
}

func newMessengerService(ctx context.Context, p *serviceParams) (messenger.MessengerService, error) {
	// Mailerの設定
	f, err := os.Open(p.config.SendGridTemplatePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var templateMap map[string]string
	d := yaml.NewDecoder(f)
	if err := d.Decode(&templateMap); err != nil {
		return nil, err
	}

	mailParams := &mailer.Params{
		APIKey:      p.config.SendGridAPIKey,
		FromName:    p.config.MailFromName,
		FromAddress: p.config.MailFromAddress,
		TemplateMap: templateMap,
	}

	userWebURL, err := url.Parse(p.config.UserWebURL)
	if err != nil {
		return nil, err
	}

	// Messenger Serviceの設定
	params := &messengersrv.Params{
		Mailer:      mailer.NewClient(mailParams, mailer.WithLogger(p.options.logger)),
		WaitGroup:   p.waitGroup,
		AdminWebURL: &url.URL{},
		UserWebURL:  userWebURL,
	}
	return messengersrv.NewMessengerService(
		params,
		messengersrv.WithLogger(p.options.logger),
	), nil
}
