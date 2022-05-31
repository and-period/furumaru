package cmd

import (
	"context"
	"net/url"
	"os"
	"sync"

	v1 "github.com/and-period/furumaru/api/internal/gateway/admin/v1/handler"
	"github.com/and-period/furumaru/api/internal/messenger"
	messengersrv "github.com/and-period/furumaru/api/internal/messenger/service"
	"github.com/and-period/furumaru/api/internal/user"
	userdb "github.com/and-period/furumaru/api/internal/user/database"
	usersrv "github.com/and-period/furumaru/api/internal/user/service"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/and-period/furumaru/api/pkg/rbac"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awscredentials "github.com/aws/aws-sdk-go-v2/credentials"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type registry struct {
	v1        v1.APIV1Handler
	logger    *zap.Logger
	waitGroup *sync.WaitGroup
}

type serviceParams struct {
	waitGroup *sync.WaitGroup
	logger    *zap.Logger
	config    *config
	aws       aws.Config
	messenger messenger.MessengerService
}

func newRegistry(ctx context.Context, conf *config) (*registry, error) {
	wg := &sync.WaitGroup{}

	// Loggerの設定
	logger, err := log.NewLogger(log.WithLogLevel(conf.LogLevel), log.WithOutput(conf.LogPath))
	if err != nil {
		return nil, err
	}

	// Casbinの設定
	enforcer, err := rbac.NewEnforcer(conf.RBACModelPath, conf.RBACPolicyPath)
	if err != nil {
		return nil, err
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
		logger:    logger,
		config:    conf,
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

	// Handlerの設定
	v1Params := &v1.Params{
		WaitGroup:   wg,
		Enforcer:    enforcer,
		Storage:     storage.NewBucket(awscfg, storageParams),
		UserService: userService,
	}

	return &registry{
		v1:        v1.NewAPIV1Handler(v1Params, v1.WithLogger(logger)),
		logger:    logger,
		waitGroup: wg,
	}, nil
}

func newDatabase(params *database.Params, tls bool, timezone string, logger *zap.Logger) (*database.Client, error) {
	return database.NewClient(
		params,
		database.WithLogger(logger),
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
	mysql, err := newDatabase(mysqlParams, p.config.DBUserEnabledTLS, p.config.DBUserTimeZone, p.logger)
	if err != nil {
		return nil, err
	}

	// Databaseの設定
	dbParams := &userdb.Params{
		Database: mysql,
	}

	// Amazon Cognitoの設定
	adminAuthParams := &cognito.Params{
		UserPoolID:      p.config.CognitoAdminPoolID,
		AppClientID:     p.config.CognitoAdminClientID,
		AppClientSecret: p.config.CognitoAdminClientSecret,
	}
	userAuthParams := &cognito.Params{
		UserPoolID:      p.config.CognitoUserPoolID,
		AppClientID:     p.config.CognitoUserClientID,
		AppClientSecret: p.config.CognitoUserClientSecret,
	}

	// User Serviceの設定
	params := &usersrv.Params{
		Database:         userdb.NewDatabase(dbParams),
		AdminAuth:        cognito.NewClient(p.aws, adminAuthParams),
		UserAuth:         cognito.NewClient(p.aws, userAuthParams),
		MessengerService: p.messenger,
		WaitGroup:        p.waitGroup,
	}
	return usersrv.NewUserService(
		params,
		usersrv.WithLogger(p.logger),
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

	adminWebURL, err := url.Parse(p.config.AminWebURL)
	if err != nil {
		return nil, err
	}
	userWebURL, err := url.Parse(p.config.UserWebURL)
	if err != nil {
		return nil, err
	}

	// Messenger Serviceの設定
	params := &messengersrv.Params{
		Mailer:      mailer.NewClient(mailParams, mailer.WithLogger(p.logger)),
		WaitGroup:   p.waitGroup,
		AdminWebURL: adminWebURL,
		UserWebURL:  userWebURL,
	}
	return messengersrv.NewMessengerService(
		params,
		messengersrv.WithLogger(p.logger),
	), nil
}
