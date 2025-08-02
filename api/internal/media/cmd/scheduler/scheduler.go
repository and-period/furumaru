package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media/broadcast/scheduler"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

type app struct {
	*cobra.Command
	waitGroup               *sync.WaitGroup
	job                     scheduler.Scheduler
	AppName                 string `default:"media-scheduler" envconfig:"APP_NAME"`
	Environment             string `default:"none"            envconfig:"ENV"`
	RunMethod               string `default:"lambda"          envconfig:"RUN_METHOD"`
	RunType                 string `default:""                envconfig:"RUN_TYPE"`
	LogPath                 string `default:""                envconfig:"LOG_PATH"`
	LogLevel                string `default:"info"            envconfig:"LOG_LEVEL"`
	DBTimeZone              string `default:"Asia/Tokyo"      envconfig:"DB_TIMEZONE"`
	TiDBHost                string `default:"127.0.0.1"       envconfig:"TIDB_HOST"`
	TiDBPort                string `default:"4000"            envconfig:"TIDB_PORT"`
	TiDBUsername            string `default:""                envconfig:"TIDB_USERNAME"`
	TiDBPassword            string `default:""                envconfig:"TIDB_PASSWORD"`
	TiDBSecretName          string `default:""                envconfig:"TIDB_SECRET_NAME"`
	SentryDsn               string `default:""                envconfig:"SENTRY_DSN"`
	SentrySecretName        string `default:""                envconfig:"SENTRY_SECRET_NAME"`
	AWSRegion               string `default:"ap-northeast-1"  envconfig:"AWS_REGION"`
	TargetDatetime          string `default:""                envconfig:"TARGET_DATETIME"`
	StepFunctionARN         string `default:""                envconfig:"STEP_FUNCTION_ARN"`
	ArchiveBucketName       string `default:""                envconfig:"ARCHIVE_BUCKET_NAME"`
	MediaConvertEndpoint    string `default:""                envconfig:"MEDIA_CONVERT_ENDPOINT"`
	MediaConvertRoleARN     string `default:""                envconfig:"MEDIA_CONVERT_ROLE_ARN"`
	MediaConvertJobTemplate string `default:""                envconfig:"MEDIA_CONVERT_JOB_TEMPLATE"`
	CDNURL                  string `default:""                envconfig:"CDN_URL"`
}

func NewApp() *app {
	cmd := &cobra.Command{
		Use:   "scheduler",
		Short: "media scheduler",
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
		slog.ErrorContext(ctx, "Failed to parse target datetime", log.Error(err), slog.String("target", a.TargetDatetime))
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
