package scheduler

import (
	"context"
	"fmt"
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
	waitGroup      *sync.WaitGroup
	job            scheduler.Scheduler
	AppName        string `envconfig:"APP_NAME" default:"messenger-scheduler"`
	Environment    string `envconfig:"ENV" default:"none"`
	RunMethod      string `envconfig:"RUN_METHOD" default:"lambda"`
	LogPath        string `envconfig:"LOG_PATH" default:""`
	LogLevel       string `envconfig:"LOG_LEVEL" default:"info"`
	DBSocket       string `envconfig:"DB_SOCKET" default:"tcp"`
	DBHost         string `envconfig:"DB_HOST" default:"127.0.0.1"`
	DBPort         string `envconfig:"DB_PORT" default:"3306"`
	DBUsername     string `envconfig:"DB_USERNAME" default:"root"`
	DBPassword     string `envconfig:"DB_PASSWORD" default:""`
	DBTimeZone     string `envconfig:"DB_TIMEZONE" default:"Asia/Tokyo"`
	DBEnabledTLS   bool   `envconfig:"DB_ENABLED_TLS" default:"false"`
	DBSecretName   string `envconfig:"DB_SECRET_NAME" default:""`
	AWSRegion      string `envconfig:"AWS_REGION" default:"ap-northeast-1"`
	SQSQueueURL    string `envconfig:"SQS_QUEUE_URL" default:""`
	SQSMockEnabled bool   `envconfig:"SQS_MOCK_ENABLED" default:"false"`
	AminWebURL     string `envconfig:"ADMIN_WEB_URL" default:""`
	UserWebURL     string `envconfig:"USER_WEB_URL" default:""`
	TargetDatetime string `envconfig:"TARGET_DATETIME" default:""`
}

//nolint:revive
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

	// Loggerの設定
	logger, err := log.NewLogger(log.WithLogLevel(a.LogLevel), log.WithOutput(a.LogPath))
	if err != nil {
		return fmt.Errorf("scheduler: failed to new logger: %w", err)
	}
	defer logger.Sync() //nolint:errcheck

	// 依存関係の解決
	if err := a.inject(ctx, logger); err != nil {
		logger.Error("Failed to new registry", zap.Error(err))
		return err
	}

	// Job実行に必要な引数の生成
	target, err := a.getTarget()
	if err != nil {
		logger.Error("Failed to parse target datetime", zap.Error(err), zap.String("target", a.TargetDatetime))
		return err
	}

	// Jobの起動
	logger.Info("Started")
	switch a.RunMethod {
	case "lambda":
		lambda.StartWithOptions(a.job.Lambda, lambda.WithContext(ctx))
	default:
		err = a.job.Run(ctx, target)
	}

	defer logger.Info("Finished...")
	a.waitGroup.Wait()
	return err
}

func (a *app) getTarget() (time.Time, error) {
	if a.TargetDatetime == "" {
		return jst.Now(), nil
	}
	return jst.Parse("2006-01-02 15:04:05", a.TargetDatetime)
}
