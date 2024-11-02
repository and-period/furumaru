package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/and-period/furumaru/api/internal/messenger/worker"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type app struct {
	*cobra.Command
	logger                       *zap.Logger
	waitGroup                    *sync.WaitGroup
	worker                       worker.Worker
	AppName                      string `default:"messenger-worker" envconfig:"APP_NAME"`
	Environment                  string `default:"none"             envconfig:"ENV"`
	RunMethod                    string `default:"lambda"           envconfig:"RUN_METHOD"`
	LogPath                      string `default:""                 envconfig:"LOG_PATH"`
	LogLevel                     string `default:"info"             envconfig:"LOG_LEVEL"`
	DBTimeZone                   string `default:"Asia/Tokyo"       envconfig:"DB_TIMEZONE"`
	TiDBHost                     string `default:"127.0.0.1"        envconfig:"TIDB_HOST"`
	TiDBPort                     string `default:"4000"             envconfig:"TIDB_PORT"`
	TiDBUsername                 string `default:""                 envconfig:"TIDB_USERNAME"`
	TiDBPassword                 string `default:""                 envconfig:"TIDB_PASSWORD"`
	TiDBSecretName               string `default:""                 envconfig:"TIDB_SECRET_NAME"`
	SentryDsn                    string `default:""                 envconfig:"SENTRY_DSN"`
	SentrySecretName             string `default:""                 envconfig:"SENTRY_SECRET_NAME"`
	AWSRegion                    string `default:"ap-northeast-1"   envconfig:"AWS_REGION"`
	SendGridAPIKey               string `default:""                 envconfig:"SENDGRID_API_KEY"`
	SendGridTemplatePath         string `default:""                 envconfig:"SENDGRID_TEMPLATE_PATH"`
	SendGridSecretName           string `default:""                 envconfig:"SENDGRID_SECRET_NAME"`
	MailFromName                 string `default:""                 envconfig:"MAIL_FROM_NAME"`
	MailFromAddress              string `default:""                 envconfig:"MAIL_FROM_ADDRESS"`
	LINEChannelToken             string `default:""                 envconfig:"LINE_CHANNEL_TOKEN"`
	LINEChannelSecret            string `default:""                 envconfig:"LINE_CHANNEL_SECRET"`
	LINERoomID                   string `default:""                 envconfig:"LINE_ROOM_ID"`
	LINESecretName               string `default:""                 envconfig:"LINE_SECRET_NAME"`
	AdminFirebaseCredentialsJSON string `default:""                 envconfig:"ADMIN_FIREBASE_CREDENTIALS_JSON"`
	AdminFirebaseSecretName      string `default:""                 envconfig:"ADMIN_FIREBASE_SECRET_NAME"`
	UserFirebaseCredentialsJSON  string `default:""                 envconfig:"USER_FIREBASE_CREDENTIALS_JSON"`
	UserFirebaseSecretName       string `default:""                 envconfig:"USER_FIREBASE_SECRET_NAME"`
}

//nolint:revive
func NewApp() *app {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "messenger worker",
	}
	app := &app{Command: cmd}
	app.RunE = func(c *cobra.Command, args []string) error {
		return app.run()
	}
	return app
}

func (a *app) run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 環境変数の読み込み
	if err := envconfig.Process("", a); err != nil {
		return fmt.Errorf("worker: failed to load environment: %w", err)
	}

	// 依存関係の解決
	if err := a.inject(ctx); err != nil {
		return fmt.Errorf("worker: failed to new registry: %w", err)
	}
	defer a.logger.Sync() //nolint:errcheck

	// Workerの起動
	a.logger.Info("Started")
	switch a.RunMethod {
	case "lambda":
		lambda.StartWithOptions(a.worker.Lambda, lambda.WithContext(ctx))
	default:
		return errors.New("not implemented")
	}

	// Workerの停止
	a.logger.Info("Shutdown...")
	a.waitGroup.Wait()
	return nil
}
