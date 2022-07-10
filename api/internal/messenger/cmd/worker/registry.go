package cmd

import (
	"context"
	"net/url"
	"os"
	"sync"

	"github.com/and-period/furumaru/api/internal/messenger/worker"
	"github.com/and-period/furumaru/api/internal/user"
	userdb "github.com/and-period/furumaru/api/internal/user/database"
	usersrv "github.com/and-period/furumaru/api/internal/user/service"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type registry struct {
	waitGroup *sync.WaitGroup
	worker    worker.Worker
}

type params struct {
	config         *config
	logger         *zap.Logger
	waitGroup      *sync.WaitGroup
	mailer         mailer.Client
	aws            aws.Config
	secret         secret.Client
	adminWebURL    *url.URL
	userWebURL     *url.URL
	dbHost         string
	dbPort         string
	dbUsername     string
	dbPassword     string
	sendGridAPIKey string
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
		APIKey:      params.sendGridAPIKey,
		FromName:    conf.MailFromName,
		FromAddress: conf.MailFromAddress,
		TemplateMap: templateMap,
	}
	params.mailer = mailer.NewClient(mailParams, mailer.WithLogger(logger))

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
	userService, err := newUserService(ctx, params)
	if err != nil {
		return nil, err
	}

	// Workerの設定
	workerParams := &worker.Params{
		WaitGroup:   params.waitGroup,
		Mailer:      params.mailer,
		AdminWebURL: params.adminWebURL,
		UserWebURL:  params.userWebURL,
		User:        userService,
	}
	return &registry{
		waitGroup: params.waitGroup,
		worker:    worker.NewWorker(workerParams, worker.WithLogger(logger)),
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
	// SendGrid認証情報の取得
	if p.config.SendGridSecretName == "" {
		p.sendGridAPIKey = p.config.SendGridAPIKey
	} else {
		secrets, err := p.secret.Get(ctx, p.config.SendGridSecretName)
		if err != nil {
			return err
		}
		p.sendGridAPIKey = secrets["api_key"]
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
