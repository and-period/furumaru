package cmd

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/store/broadcast/updater"
	storedb "github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/secret"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"
)

type registry struct {
	appName   string
	env       string
	waitGroup *sync.WaitGroup
	job       updater.Updater
}

type params struct {
	config     *config
	logger     *zap.Logger
	waitGroup  *sync.WaitGroup
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

	// AWS Secrets Managerの設定
	params.secret = secret.NewClient(awscfg)
	if err := getSecret(ctx, params); err != nil {
		return nil, err
	}

	// Databaseの設定
	dbClient, err := newDatabase("stores", params)
	if err != nil {
		return nil, err
	}

	// Jobの設定
	dbParams := &storedb.Params{
		Database: dbClient,
	}
	jobParams := &updater.Params{
		WaitGroup: params.waitGroup,
		Database:  storedb.NewDatabase(dbParams),
	}
	var job updater.Updater
	switch conf.RunType {
	case "CREATE":
		job = updater.NewCreator(jobParams, updater.WithLogger(logger))
	case "REMOVE":
		job = updater.NewRemover(jobParams, updater.WithLogger(logger))
	default:
		return nil, fmt.Errorf("cmd: unknown scheduler type. type=%s", conf.RunType)
	}

	return &registry{
		appName:   conf.AppName,
		env:       conf.Environment,
		waitGroup: params.waitGroup,
		job:       job,
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
	location, err := time.LoadLocation(p.config.DBTimeZone)
	if err != nil {
		return nil, err
	}
	return database.NewClient(
		params,
		database.WithLogger(p.logger),
		database.WithNow(p.now),
		database.WithTLS(p.config.DBEnabledTLS),
		database.WithLocation(location),
	)
}
