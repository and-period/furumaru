package worker

import (
	"context"
	"os"
	"sync"
	"time"

	firebase "firebase.google.com/go/v4"
	messengerdb "github.com/and-period/furumaru/api/internal/messenger/database/mysql"
	"github.com/and-period/furumaru/api/internal/messenger/worker"
	"github.com/and-period/furumaru/api/internal/user"
	userdb "github.com/and-period/furumaru/api/internal/user/database/mysql"
	usersrv "github.com/and-period/furumaru/api/internal/user/service"
	"github.com/and-period/furumaru/api/pkg/firebase/messaging"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/line"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/and-period/furumaru/api/pkg/secret"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rafaelhl/gorm-newrelic-telemetry-plugin/telemetry"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v2"
)

type params struct {
	logger            *zap.Logger
	waitGroup         *sync.WaitGroup
	mailer            mailer.Client
	line              line.Client
	messaging         messaging.Client
	aws               aws.Config
	firebase          *firebase.App
	secret            secret.Client
	now               func() time.Time
	dbHost            string
	dbPort            string
	dbUsername        string
	dbPassword        string
	sendGridAPIKey    string
	lineToken         string
	lineSecret        string
	lineRoomID        string
	googleCredentials []byte
}

func (a *app) inject(ctx context.Context, logger *zap.Logger) error {
	params := &params{
		logger:    logger,
		now:       jst.Now,
		waitGroup: &sync.WaitGroup{},
	}

	// AWS SDKの設定
	awscfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(a.AWSRegion))
	if err != nil {
		return err
	}
	params.aws = awscfg

	// AWS Secrets Managerの設定
	params.secret = secret.NewClient(awscfg)
	if err := a.getSecret(ctx, params); err != nil {
		return err
	}

	// Databaseの設定
	dbClient, err := a.newDatabase("messengers", params)
	if err != nil {
		return err
	}

	// メールテンプレートの設定
	f, err := os.Open(a.SendGridTemplatePath)
	if err != nil {
		return err
	}
	defer f.Close()
	var templateMap map[string]string
	d := yaml.NewDecoder(f)
	if err := d.Decode(&templateMap); err != nil {
		return err
	}

	// Mailerの設定
	mailParams := &mailer.Params{
		APIKey:      params.sendGridAPIKey,
		FromName:    a.MailFromName,
		FromAddress: a.MailFromAddress,
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
		return err
	}
	params.line = linebot

	// Firebaseの設定
	fbapp, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON(params.googleCredentials))
	if err != nil {
		return err
	}
	params.firebase = fbapp

	// Firebase Cloud Messagingの設定
	messaging, err := messaging.NewClient(ctx, fbapp, messaging.WithLogger(logger))
	if err != nil {
		return err
	}
	params.messaging = messaging

	// Serviceの設定
	userService, err := a.newUserService(params)
	if err != nil {
		return err
	}

	// Workerの設定
	workerParams := &worker.Params{
		WaitGroup: params.waitGroup,
		DB:        messengerdb.NewDatabase(dbClient),
		Mailer:    params.mailer,
		Line:      params.line,
		Messaging: params.messaging,
		User:      userService,
	}
	a.worker = worker.NewWorker(workerParams, worker.WithLogger(logger))
	a.waitGroup = params.waitGroup
	return nil
}

func (a *app) getSecret(ctx context.Context, p *params) error {
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// データベース認証情報の取得
		if a.DBSecretName == "" {
			p.dbHost = a.DBHost
			p.dbPort = a.DBPort
			p.dbUsername = a.DBUsername
			p.dbPassword = a.DBPassword
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.DBSecretName)
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
		if a.SendGridSecretName == "" {
			p.sendGridAPIKey = a.SendGridAPIKey
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.SendGridSecretName)
		if err != nil {
			return err
		}
		p.sendGridAPIKey = secrets["api_key"]
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
		// Google認証情報の取得
		if a.GoogleSecretName == "" {
			p.googleCredentials = []byte(a.GoogleCredentialsJSON)
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.GoogleSecretName)
		if err != nil {
			return err
		}
		p.googleCredentials = []byte(secrets["credentials"])
		return nil
	})
	return eg.Wait()
}

func (a *app) newDatabase(dbname string, p *params) (*mysql.Client, error) {
	params := &mysql.Params{
		Socket:   a.DBSocket,
		Host:     p.dbHost,
		Port:     p.dbPort,
		Database: dbname,
		Username: p.dbUsername,
		Password: p.dbPassword,
	}
	location, err := time.LoadLocation(a.DBTimeZone)
	if err != nil {
		return nil, err
	}
	cli, err := mysql.NewClient(
		params,
		mysql.WithLogger(p.logger),
		mysql.WithNow(p.now),
		mysql.WithTLS(a.DBEnabledTLS),
		mysql.WithLocation(location),
	)
	if err != nil {
		return nil, err
	}
	if err := cli.DB.Use(telemetry.NewNrTracer(dbname, p.dbHost, string(newrelic.DatastoreMySQL))); err != nil {
		return nil, err
	}
	return cli, nil
}

func (a *app) newUserService(p *params) (user.Service, error) {
	mysql, err := a.newDatabase("users", p)
	if err != nil {
		return nil, err
	}
	params := &usersrv.Params{
		WaitGroup: p.waitGroup,
		Database:  userdb.NewDatabase(mysql),
	}
	return usersrv.NewService(params, usersrv.WithLogger(p.logger)), nil
}
