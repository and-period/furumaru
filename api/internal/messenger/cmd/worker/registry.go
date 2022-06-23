package cmd

import (
	"context"
	"os"
	"sync"

	"github.com/and-period/furumaru/api/internal/messenger/worker"
	"github.com/and-period/furumaru/api/internal/user"
	userdb "github.com/and-period/furumaru/api/internal/user/database"
	usersrv "github.com/and-period/furumaru/api/internal/user/service"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type registry struct {
	waitGroup *sync.WaitGroup
	worker    worker.Worker
}

type params struct {
	config    *config
	logger    *zap.Logger
	waitGroup *sync.WaitGroup
	mailer    mailer.Client
	user      user.Service
}

func newRegistry(ctx context.Context, conf *config, logger *zap.Logger) (*registry, error) {
	params := &params{
		config:    conf,
		logger:    logger,
		waitGroup: &sync.WaitGroup{},
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
		APIKey:      conf.SendGridAPIKey,
		FromName:    conf.MailFromName,
		FromAddress: conf.MailFromAddress,
		TemplateMap: templateMap,
	}
	params.mailer = mailer.NewClient(mailParams, mailer.WithLogger(logger))

	// Serviceの設定
	userService, err := newUserService(ctx, params)
	if err != nil {
		return nil, err
	}

	// Workerの設定
	workerParams := &worker.Params{
		WaitGroup: params.waitGroup,
		Mailer:    params.mailer,
		User:      userService,
	}
	return &registry{
		waitGroup: params.waitGroup,
		worker:    worker.NewWorker(workerParams, worker.WithLogger(logger)),
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

func newUserService(ctx context.Context, p *params) (user.Service, error) {
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
	}
	return usersrv.NewService(params, usersrv.WithLogger(p.logger)), nil
}
