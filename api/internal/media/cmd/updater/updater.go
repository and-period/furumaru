package updater

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/and-period/furumaru/api/internal/media/broadcast/updater"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type app struct {
	*cobra.Command
	waitGroup    *sync.WaitGroup
	updater      updater.Updater
	AppName      string `envconfig:"APP_NAME" default:"media-scheduler"`
	Environment  string `envconfig:"ENV" default:"none"`
	RunMethod    string `envconfig:"RUN_METHOD" default:"lambda"`
	RunType      string `envconfig:"RUN_TYPE" default:""`
	LogPath      string `envconfig:"LOG_PATH" default:""`
	LogLevel     string `envconfig:"LOG_LEVEL" default:"info"`
	DBSocket     string `envconfig:"DB_SOCKET" default:"tcp"`
	DBHost       string `envconfig:"DB_HOST" default:"127.0.0.1"`
	DBPort       string `envconfig:"DB_PORT" default:"3306"`
	DBUsername   string `envconfig:"DB_USERNAME" default:"root"`
	DBPassword   string `envconfig:"DB_PASSWORD" default:""`
	DBTimeZone   string `envconfig:"DB_TIMEZONE" default:"Asia/Tokyo"`
	DBEnabledTLS bool   `envconfig:"DB_ENABLED_TLS" default:"false"`
	DBSecretName string `envconfig:"DB_SECRET_NAME" default:""`
	AWSRegion    string `envconfig:"AWS_REGION" default:"ap-northeast-1"`
}

//nolint:revive
func NewApp() *app {
	cmd := &cobra.Command{
		Use:   "updater",
		Short: "media updater",
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
		return fmt.Errorf("updater: failed to load environment: %w", err)
	}

	// Loggerの設定
	logger, err := log.NewLogger(log.WithLogLevel(a.LogLevel), log.WithOutput(a.LogPath))
	if err != nil {
		return fmt.Errorf("updater: failed to new logger: %w", err)
	}
	defer logger.Sync() //nolint:errcheck

	// 依存関係の解決
	if err := a.inject(ctx, logger); err != nil {
		logger.Error("Failed to new registry", zap.Error(err))
		return err
	}

	// Jobの起動
	logger.Info("Started")
	switch a.RunMethod {
	case "lambda":
		lambda.StartWithOptions(a.updater.Lambda, lambda.WithContext(ctx))
	default:
		err = errors.New("not implemented")
	}

	defer logger.Info("Finished...")
	a.waitGroup.Wait()
	return err
}