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
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type registry struct {
	worker    worker.Worker
	logger    *zap.Logger
	waitGroup *sync.WaitGroup
}

type serviceParams struct {
	waitGroup *sync.WaitGroup
	logger    *zap.Logger
	config    *config
	user      user.UserService
}

func newRegistry(ctx context.Context, conf *config) (*registry, error) {
	wg := &sync.WaitGroup{}

	// Loggerの設定
	logger, err := log.NewLogger(log.WithLogLevel(conf.LogLevel), log.WithOutput(conf.LogPath))
	if err != nil {
		return nil, err
	}

	// Mailerの設定
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

	mailParams := &mailer.Params{
		APIKey:      conf.SendGridAPIKey,
		FromName:    conf.MailFromName,
		FromAddress: conf.MailFromAddress,
		TemplateMap: templateMap,
	}

	// Serviceの設定
	srvParams := &serviceParams{
		waitGroup: wg,
		logger:    logger,
		config:    conf,
	}

	userService, err := newUserService(ctx, srvParams)
	if err != nil {
		return nil, err
	}
	srvParams.user = userService

	// Workerの設定
	workerParams := &worker.Params{
		WaitGroup:   wg,
		Mailer:      mailer.NewClient(mailParams, mailer.WithLogger(logger)),
		UserService: userService,
	}

	return &registry{
		worker:    worker.NewWorker(workerParams, worker.WithLogger(logger)),
		logger:    logger,
		waitGroup: wg,
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

func newUserService(ctx context.Context, p *serviceParams) (user.UserService, error) {
	// MySQLの設定
	mysql, err := newDatabase("users", p.config, p.logger)
	if err != nil {
		return nil, err
	}

	// Databaseの設定
	dbParams := &userdb.Params{
		Database: mysql,
	}

	// User Serviceの設定
	params := &usersrv.Params{
		Database:  userdb.NewDatabase(dbParams),
		WaitGroup: p.waitGroup,
	}
	return usersrv.NewUserService(
		params,
		usersrv.WithLogger(p.logger),
	), nil
}
