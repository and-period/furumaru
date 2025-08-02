package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/scheduler"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type app struct {
	*cobra.Command
	waitGroup        *sync.WaitGroup
	job              scheduler.Scheduler
	AppName          string `default:"messenger-scheduler" envconfig:"APP_NAME"`
	Environment      string `default:"none"                envconfig:"ENV"`
	RunMethod        string `default:"lambda"              envconfig:"RUN_METHOD"`
	LogPath          string `default:""                    envconfig:"LOG_PATH"`
	LogLevel         string `default:"info"                envconfig:"LOG_LEVEL"`
	DBTimeZone       string `default:"Asia/Tokyo"          envconfig:"DB_TIMEZONE"`
	TiDBHost         string `default:"127.0.0.1"           envconfig:"TIDB_HOST"`
	TiDBPort         string `default:"4000"                envconfig:"TIDB_PORT"`
	TiDBUsername     string `default:""                    envconfig:"TIDB_USERNAME"`
	TiDBPassword     string `default:""                    envconfig:"TIDB_PASSWORD"`
	TiDBSecretName   string `default:""                    envconfig:"TIDB_SECRET_NAME"`
	SentryDsn        string `default:""                    envconfig:"SENTRY_DSN"`
	SentrySecretName string `default:""                    envconfig:"SENTRY_SECRET_NAME"`
	AWSRegion        string `default:"ap-northeast-1"      envconfig:"AWS_REGION"`
	SQSQueueURL      string `default:""                    envconfig:"SQS_QUEUE_URL"`
	SQSMockEnabled   bool   `default:"false"               envconfig:"SQS_MOCK_ENABLED"`
	AminWebURL       string `default:""                    envconfig:"ADMIN_WEB_URL"`
	UserWebURL       string `default:""                    envconfig:"USER_WEB_URL"`
	TargetDatetime   string `default:""                    envconfig:"TARGET_DATETIME"`
}

func NewApp() *app {
	cmd := &cobra.Command{
		Use:   "scheduler",
		Short: "messenger scheduler",
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
		return fmt.Errorf("scheduler: failed to load environment: %w", err)
	}

	// ログの設定
	logOpts := []log.Option{
		log.WithLogLevel(a.LogLevel),
		log.WithSentryDSN(a.SentryDsn),
		log.WithSentryServerName(a.AppName),
		log.WithSentryEnvironment(a.Environment),
		log.WithSentryLevel("error"),
	}
	logFlush, err := log.Start(ctx, logOpts...)
	if err != nil {
		return fmt.Errorf("user: failed to start logger: %w", err)
	}
	defer logFlush()

	// 依存関係の解決
	if err := a.inject(ctx); err != nil {
		return fmt.Errorf("scheduler: failed to new registry: %w", err)
	}

	// Job実行に必要な引数の生成
	target, err := a.getTarget()
	if err != nil {
		slog.Error("Failed to parse target datetime", zap.Error(err), zap.String("target", a.TargetDatetime))
		return err
	}

	// Jobの起動
	slog.Info("Started")
	switch a.RunMethod {
	case "lambda":
		lambda.StartWithOptions(a.job.Lambda, lambda.WithContext(ctx))
	default:
		err = a.job.Run(ctx, target)
	}

	defer slog.Info("Finished...")
	a.waitGroup.Wait()
	return err
}

func (a *app) getTarget() (time.Time, error) {
	if a.TargetDatetime == "" {
		return jst.Now(), nil
	}
	return jst.Parse("2006-01-02 15:04:05", a.TargetDatetime)
}
