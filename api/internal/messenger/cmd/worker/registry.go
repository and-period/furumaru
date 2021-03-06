package cmd

import (
	"context"
	"os"
	"sync"

	messengerdb "github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/worker"
	"github.com/and-period/furumaru/api/internal/user"
	userdb "github.com/and-period/furumaru/api/internal/user/database"
	usersrv "github.com/and-period/furumaru/api/internal/user/service"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/line"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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
	line           line.Client
	aws            aws.Config
	secret         secret.Client
	dbHost         string
	dbPort         string
	dbUsername     string
	dbPassword     string
	sendGridAPIKey string
	lineToken      string
	lineSecret     string
	lineRoomID     string
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

	// Databaseの設定
	mysql, err := newDatabase("messengers", params)
	if err != nil {
		return nil, err
	}
	dbParams := &messengerdb.Params{
		Database: mysql,
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

	// LINEの設定
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

	// Serviceの設定
	userService, err := newUserService(ctx, params)
	if err != nil {
		return nil, err
	}

	// Workerの設定
	workerParams := &worker.Params{
		WaitGroup: params.waitGroup,
		DB:        messengerdb.NewDatabase(dbParams),
		Mailer:    params.mailer,
		Line:      params.line,
		User:      userService,
	}
	return &registry{
		waitGroup: params.waitGroup,
		worker:    worker.NewWorker(workerParams, worker.WithLogger(logger)),
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
		// SendGrid認証情報の取得
		if p.config.SendGridSecretName == "" {
			p.sendGridAPIKey = p.config.SendGridAPIKey
			return nil
		}
		secrets, err := p.secret.Get(ectx, p.config.SendGridSecretName)
		if err != nil {
			return err
		}
		p.sendGridAPIKey = secrets["api_key"]
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
