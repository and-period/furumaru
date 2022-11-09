package cmd

import (
	"context"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media/resizer"
	"github.com/and-period/furumaru/api/internal/user"
	userdb "github.com/and-period/furumaru/api/internal/user/database"
	usersrv "github.com/and-period/furumaru/api/internal/user/service"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rafaelhl/gorm-newrelic-telemetry-plugin/telemetry"
	"go.uber.org/zap"
)

type registry struct {
	appName   string
	env       string
	waitGroup *sync.WaitGroup
	resizer   resizer.Resizer
}

type params struct {
	config     *config
	logger     *zap.Logger
	waitGroup  *sync.WaitGroup
	aws        aws.Config
	storage    storage.Bucket
	secret     secret.Client
	now        func() time.Time
	dbHost     string
	dbPort     string
	dbUsername string
	dbPassword string
}

func newRegistry(ctx context.Context, conf *config, logger *zap.Logger) (*registry, error) {
	params := &params{
		config:    conf,
		logger:    logger,
		now:       jst.Now,
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

	// Amazon S3の設定
	storageParams := &storage.Params{
		Bucket: conf.S3Bucket,
	}
	params.storage = storage.NewBucket(awscfg, storageParams)

	// Serviceの設定
	userService, err := newUserService(params)
	if err != nil {
		return nil, err
	}

	// Resizerの設定
	resizerParams := &resizer.Params{
		WaitGroup: params.waitGroup,
		Storage:   params.storage,
		User:      userService,
	}
	return &registry{
		appName:   conf.AppName,
		env:       conf.Environment,
		waitGroup: params.waitGroup,
		resizer:   resizer.NewResizer(resizerParams, resizer.WithLogger(logger)),
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
	cli, err := database.NewClient(
		params,
		database.WithLogger(p.logger),
		database.WithNow(p.now),
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

func newUserService(p *params) (user.Service, error) {
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
