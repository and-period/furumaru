package updater

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/and-period/furumaru/api/internal/media/broadcast/updater"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type app struct {
	*cobra.Command
	logger           *zap.Logger
	waitGroup        *sync.WaitGroup
	updater          updater.Updater
	AppName          string `default:"media-scheduler" envconfig:"APP_NAME"`
	Environment      string `default:"none"            envconfig:"ENV"`
	RunMethod        string `default:"lambda"          envconfig:"RUN_METHOD"`
	RunType          string `default:""                envconfig:"RUN_TYPE"`
	LogPath          string `default:""                envconfig:"LOG_PATH"`
	LogLevel         string `default:"info"            envconfig:"LOG_LEVEL"`
	DBTimeZone       string `default:"Asia/Tokyo"      envconfig:"DB_TIMEZONE"`
	TiDBHost         string `default:"127.0.0.1"       envconfig:"TIDB_HOST"`
	TiDBPort         string `default:"4000"            envconfig:"TIDB_PORT"`
	TiDBUsername     string `default:""                envconfig:"TIDB_USERNAME"`
	TiDBPassword     string `default:""                envconfig:"TIDB_PASSWORD"`
	TiDBSecretName   string `default:""                envconfig:"TIDB_SECRET_NAME"`
	SentryDsn        string `default:""                envconfig:"SENTRY_DSN"`
	SentrySecretName string `default:""                envconfig:"SENTRY_SECRET_NAME"`
	AWSRegion        string `default:"ap-northeast-1"  envconfig:"AWS_REGION"`
	CDNURL           string `default:""                envconfig:"CDN_URL"`
}

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

	// 依存関係の解決
	if err := a.inject(ctx); err != nil {
		return fmt.Errorf("updater: failed to new registry: %w", err)
	}
	defer a.logger.Sync() //nolint:errcheck

	// Jobの起動
	a.logger.Info("Started")
	switch a.RunMethod {
	case "lambda":
		lambda.StartWithOptions(a.updater.Lambda, lambda.WithContext(ctx))
	default:
		return errors.New("not implemented")
	}

	defer a.logger.Info("Finished...")
	a.waitGroup.Wait()
	return nil
}
