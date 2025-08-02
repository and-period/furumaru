package uploader

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/and-period/furumaru/api/internal/media/uploader"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

type app struct {
	*cobra.Command
	waitGroup        *sync.WaitGroup
	uploader         uploader.Uploader
	AppName          string `default:"media-uploader" envconfig:"APP_NAME"`
	Environment      string `default:"none"           envconfig:"ENV"`
	RunMethod        string `default:"lambda"         envconfig:"RUN_METHOD"`
	RunType          string `default:""               envconfig:"RUN_TYPE"`
	LogPath          string `default:""               envconfig:"LOG_PATH"`
	LogLevel         string `default:"info"           envconfig:"LOG_LEVEL"`
	SentryDsn        string `default:""               envconfig:"SENTRY_DSN"`
	SentrySecretName string `default:""               envconfig:"SENTRY_SECRET_NAME"`
	AWSRegion        string `default:"ap-northeast-1" envconfig:"AWS_REGION"`
	S3Bucket         string `default:""               envconfig:"S3_BUCKET"`
	S3TmpBucket      string `default:""               envconfig:"S3_TMP_BUCKET"`
	CDNURL           string `default:""               envconfig:"CDN_URL"`
}

func NewApp() *app {
	cmd := &cobra.Command{
		Use:   "uploader",
		Short: "media uploader",
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
		return fmt.Errorf("updater: failed to new registry: %w", err)
	}

	// Jobの起動
	slog.Info("Started")
	switch a.RunMethod {
	case "lambda":
		lambda.StartWithOptions(a.uploader.Lambda, lambda.WithContext(ctx))
	default:
		return errors.New("not implemented")
	}

	defer slog.Info("Finished...")
	a.waitGroup.Wait()
	return nil
}
