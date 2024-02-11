package uploader

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/and-period/furumaru/api/internal/media/uploader"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type app struct {
	*cobra.Command
	logger           *zap.Logger
	waitGroup        *sync.WaitGroup
	uploader         uploader.Uploader
	AppName          string `envconfig:"APP_NAME" default:"media-uploader"`
	Environment      string `envconfig:"ENV" default:"none"`
	RunMethod        string `envconfig:"RUN_METHOD" default:"lambda"`
	RunType          string `envconfig:"RUN_TYPE" default:""`
	LogPath          string `envconfig:"LOG_PATH" default:""`
	LogLevel         string `envconfig:"LOG_LEVEL" default:"info"`
	SentryDsn        string `envconfig:"SENTRY_DSN" default:""`
	SentrySecretName string `envconfig:"SENTRY_SECRET_NAME" default:""`
	AWSRegion        string `envconfig:"AWS_REGION" default:"ap-northeast-1"`
	S3Bucket         string `envconfig:"S3_BUCKET" default:""`
	S3TmpBucket      string `envconfig:"S3_TMP_BUCKET" default:""`
	CDNURL           string `envconfig:"CDN_URL" default:""`
}

//nolint:revive
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

	// 依存関係の解決
	if err := a.inject(ctx); err != nil {
		return fmt.Errorf("updater: failed to new registry: %w", err)
	}
	defer a.logger.Sync() //nolint:errcheck

	// Jobの起動
	a.logger.Info("Started")
	switch a.RunMethod {
	case "lambda":
		lambda.StartWithOptions(a.uploader.Lambda, lambda.WithContext(ctx))
	default:
		return errors.New("not implemented")
	}

	defer a.logger.Info("Finished...")
	a.waitGroup.Wait()
	return nil
}
