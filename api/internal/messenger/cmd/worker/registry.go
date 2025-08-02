package worker

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	firebase "firebase.google.com/go/v4"
	messengerdb "github.com/and-period/furumaru/api/internal/messenger/database/tidb"
	"github.com/and-period/furumaru/api/internal/messenger/worker"
	"github.com/and-period/furumaru/api/internal/user"
	userdb "github.com/and-period/furumaru/api/internal/user/database/tidb"
	usersrv "github.com/and-period/furumaru/api/internal/user/service"
	"github.com/and-period/furumaru/api/pkg/firebase/messaging"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/line"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/and-period/furumaru/api/pkg/secret"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rafaelhl/gorm-newrelic-telemetry-plugin/telemetry"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/option"
	yaml "gopkg.in/yaml.v2"
)

type params struct {
	logger                   *zap.Logger
	waitGroup                *sync.WaitGroup
	mailer                   mailer.Client
	line                     line.Client
	adminMessaging           messaging.Client
	userMessaging            messaging.Client
	secret                   secret.Client
	now                      func() time.Time
	tidbHost                 string
	tidbPort                 string
	tidbUsername             string
	tidbPassword             string
	sentryDsn                string
	sendGridAPIKey           string
	sendGridTemplateMap      map[string]string
	lineToken                string
	lineSecret               string
	lineRoomID               string
	adminFirebaseCredentials []byte
	userFirebaseCredentials  []byte
}

func (a *app) inject(ctx context.Context) error {
	params := &params{
		logger:    zap.NewNop(),
		now:       jst.Now,
		waitGroup: &sync.WaitGroup{},
	}

	// AWS SDKの設定
	awscfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(a.AWSRegion))
	if err != nil {
		return fmt.Errorf("cmd: failed to load aws config: %w", err)
	}

	// AWS Secrets Managerの設定
	params.secret = secret.NewClient(awscfg)
	if err := a.getSecret(ctx, params); err != nil {
		return fmt.Errorf("cmd: failed to get secret: %w", err)
	}

	// Loggerの設定
	logger, err := log.NewSentryLogger(params.sentryDsn,
		log.WithLogLevel(a.LogLevel),
		log.WithSentryServerName(a.AppName),
		log.WithSentryEnvironment(a.Environment),
		log.WithSentryLevel("error"),
	)
	if err != nil {
		return fmt.Errorf("cmd: failed to create sentry logger: %w", err)
	}
	params.logger = logger

	// Databaseの設定
	dbClient, err := a.newTiDB("messengers", params)
	if err != nil {
		return fmt.Errorf("cmd: failed to create database client: %w", err)
	}

	// メールテンプレートの設定
	if params.sendGridTemplateMap == nil {
		f, err := os.Open(a.SendGridTemplatePath)
		if err != nil {
			return fmt.Errorf("cmd: failed to open sendgrid template file: %w", err)
		}
		defer f.Close() //nolint:errcheck
		var templateMap map[string]string
		d := yaml.NewDecoder(f)
		if err := d.Decode(&templateMap); err != nil {
			return fmt.Errorf("cmd: failed to decode sendgrid template yaml: %w", err)
		}
		params.sendGridTemplateMap = templateMap
	}

	// Mailerの設定
	mailParams := &mailer.Params{
		APIKey:      params.sendGridAPIKey,
		FromName:    a.MailFromName,
		FromAddress: a.MailFromAddress,
		TemplateMap: params.sendGridTemplateMap,
	}
	params.mailer = mailer.NewClient(mailParams, mailer.WithLogger(params.logger))

	// LINEの設定
	lineParams := &line.Params{
		Token:  params.lineToken,
		Secret: params.lineSecret,
		RoomID: params.lineRoomID,
	}
	linebot, err := line.NewClient(lineParams, line.WithLogger(params.logger))
	if err != nil {
		return fmt.Errorf("cmd: failed to create line client: %w", err)
	}
	params.line = linebot

	// Firebaseの設定（管理者用）
	afbapp, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON(params.adminFirebaseCredentials))
	if err != nil {
		return fmt.Errorf("cmd: failed to create firebase client for admin: %w", err)
	}

	// Firebase Cloud Messagingの設定（管理者用）
	amessaging, err := messaging.NewClient(ctx, afbapp, messaging.WithLogger(params.logger))
	if err != nil {
		return fmt.Errorf("cmd: failed to create firebase messaging client for admin: %w", err)
	}
	params.adminMessaging = amessaging

	// Firebaseの設定（利用者用）
	ufbapp, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON(params.userFirebaseCredentials))
	if err != nil {
		return fmt.Errorf("cmd: failed to create firebase client for user: %w", err)
	}

	// Firebase Cloud Messagingの設定（利用者用）
	umessaging, err := messaging.NewClient(ctx, ufbapp, messaging.WithLogger(params.logger))
	if err != nil {
		return fmt.Errorf("cmd: failed to create firebase messaging client for user: %w", err)
	}
	params.userMessaging = umessaging

	// Serviceの設定
	userService, err := a.newUserService(params)
	if err != nil {
		return fmt.Errorf("cmd: failed to create user service: %w", err)
	}

	// Workerの設定
	workerParams := &worker.Params{
		WaitGroup:      params.waitGroup,
		DB:             messengerdb.NewDatabase(dbClient),
		Mailer:         params.mailer,
		Line:           params.line,
		AdminMessaging: params.adminMessaging,
		UserMessaging:  params.userMessaging,
		User:           userService,
	}
	a.worker = worker.NewWorker(workerParams, worker.WithLogger(params.logger))
	a.waitGroup = params.waitGroup
	return nil
}

func (a *app) getSecret(ctx context.Context, p *params) error {
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// データベース（TiDB）認証情報の取得
		if a.TiDBSecretName == "" {
			p.tidbHost = a.TiDBHost
			p.tidbPort = a.TiDBPort
			p.tidbUsername = a.TiDBUsername
			p.tidbPassword = a.TiDBPassword
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.TiDBSecretName)
		if err != nil {
			return err
		}
		p.tidbHost = secrets["host"]
		p.tidbPort = secrets["port"]
		p.tidbUsername = secrets["username"]
		p.tidbPassword = secrets["password"]
		return nil
	})
	eg.Go(func() error {
		// Sentry認証情報の取得
		if a.SentrySecretName == "" {
			p.sentryDsn = a.SentryDsn
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.SentrySecretName)
		if err != nil {
			return err
		}
		p.sentryDsn = secrets["dsn"]
		return nil
	})
	eg.Go(func() error {
		// SendGrid認証情報の取得
		if a.SendGridAPIKeySecretName == "" {
			p.sendGridAPIKey = a.SendGridAPIKey
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.SendGridAPIKeySecretName)
		if err != nil {
			return err
		}
		p.sendGridAPIKey = secrets["api_key"]
		return nil
	})
	eg.Go(func() error {
		// SendGridテンプレート情報の取得
		if a.SendGridTemplateSecretName == "" {
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.SendGridTemplateSecretName)
		if err != nil {
			return err
		}
		p.sendGridTemplateMap = secrets
		return nil
	})
	eg.Go(func() error {
		// LINE認証情報の取得
		if a.LINESecretName == "" {
			p.lineToken = a.LINEChannelToken
			p.lineSecret = a.LINEChannelSecret
			p.lineRoomID = a.LINERoomID
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.LINESecretName)
		if err != nil {
			return err
		}
		p.lineToken = secrets["token"]
		p.lineSecret = secrets["secret"]
		p.lineRoomID = secrets["roomId"]
		return nil
	})
	eg.Go(func() error {
		// Firebase認証情報の取得（管理者用）
		if a.AdminFirebaseSecretName == "" {
			p.adminFirebaseCredentials = []byte(a.AdminFirebaseCredentialsJSON)
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.AdminFirebaseSecretName)
		if err != nil {
			return err
		}
		p.adminFirebaseCredentials = []byte(secrets["credentials"])
		return nil
	})
	eg.Go(func() error {
		// Firebase認証情報の取得（利用者用）
		if a.UserFirebaseSecretName == "" {
			p.userFirebaseCredentials = []byte(a.UserFirebaseCredentialsJSON)
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.UserFirebaseSecretName)
		if err != nil {
			return err
		}
		p.userFirebaseCredentials = []byte(secrets["credentials"])
		return nil
	})
	return eg.Wait()
}

func (a *app) newTiDB(dbname string, p *params) (*mysql.Client, error) {
	params := &mysql.Params{
		Host:     p.tidbHost,
		Port:     p.tidbPort,
		Database: dbname,
		Username: p.tidbUsername,
		Password: p.tidbPassword,
	}
	location, err := time.LoadLocation(a.DBTimeZone)
	if err != nil {
		return nil, err
	}
	cli, err := mysql.NewTiDBClient(
		params,
		mysql.WithNow(p.now),
		mysql.WithLocation(location),
	)
	if err != nil {
		return nil, err
	}
	if err := cli.DB.Use(telemetry.NewNrTracer(dbname, p.tidbHost, string(newrelic.DatastoreMySQL))); err != nil {
		return nil, err
	}
	return cli, nil
}

func (a *app) newUserService(p *params) (user.Service, error) {
	mysql, err := a.newTiDB("users", p)
	if err != nil {
		return nil, err
	}
	params := &usersrv.Params{
		WaitGroup: p.waitGroup,
		Database:  userdb.NewDatabase(mysql),
	}
	return usersrv.NewService(params, usersrv.WithLogger(p.logger)), nil
}
