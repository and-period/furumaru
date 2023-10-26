package user

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	v1 "github.com/and-period/furumaru/api/internal/gateway/user/v1/handler"
	"github.com/and-period/furumaru/api/pkg/http"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/slack"
	"github.com/kelseyhightower/envconfig"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type app struct {
	*cobra.Command
	debugMode            bool
	waitGroup            *sync.WaitGroup
	slack                slack.Client
	newRelic             *newrelic.Application
	v1                   v1.Handler
	AppName              string `envconfig:"APP_NAME" default:"user-gateway"`
	Environment          string `envconfig:"ENV" default:"none"`
	Port                 int64  `envconfig:"PORT" default:"8080"`
	MetricsPort          int64  `envconfig:"METRICS_PORT" default:"9090"`
	ShutdownDelaySec     int64  `envconfig:"SHUTDOWN_DELAY_SEC" default:"20"`
	LogPath              string `envconfig:"LOG_PATH" default:""`
	LogLevel             string `envconfig:"LOG_LEVEL" default:"info"`
	DBSocket             string `envconfig:"DB_SOCKET" default:"tcp"`
	DBHost               string `envconfig:"DB_HOST" default:"127.0.0.1"`
	DBPort               string `envconfig:"DB_PORT" default:"3306"`
	DBUsername           string `envconfig:"DB_USERNAME" default:"root"`
	DBPassword           string `envconfig:"DB_PASSWORD" default:""`
	DBTimeZone           string `envconfig:"DB_TIMEZONE" default:"Asia/Tokyo"`
	DBEnabledTLS         bool   `envconfig:"DB_ENABLED_TLS" default:"false"`
	DBSecretName         string `envconfig:"DB_SECRET_NAME" default:""`
	NewRelicLicense      string `envconfig:"NEW_RELIC_LICENSE" default:""`
	NewRelicSecretName   string `envconfig:"NEW_RELIC_SECRET_NAME" default:""`
	AWSRegion            string `envconfig:"AWS_REGION" default:"ap-northeast-1"`
	S3Bucket             string `envconfig:"S3_BUCKET" default:""`
	S3TmpBucket          string `envconfig:"S3_TMP_BUCKET" default:""`
	CognitoUserPoolID    string `envconfig:"COGNITO_USER_POOL_ID" default:""`
	CognitoUserClientID  string `envconfig:"COGNITO_USER_CLIENT_ID" default:""`
	SQSQueueURL          string `envconfig:"SQS_QUEUE_URL" default:""`
	SQSMockEnabled       bool   `envconfig:"SQS_MOCK_ENABLED" default:"false"`
	KomojuHost           string `envconfig:"KOMOJU_HOST" default:""`
	KomojuClientID       string `encvonfig:"KOMOJU_CLIENT_ID" default:""`
	KomojuClientPassword string `encvonfig:"KOMOJU_CLIENT_PASSWORD" default:""`
	KomojuSecretName     string `encvonfig:"KOMOJU_SECRET_NAME" default:""`
	SlackAPIToken        string `envconfig:"SLACK_API_TOKEN" default:""`
	SlackChannelID       string `envconfig:"SLACK_CHANNEL_ID" default:""`
	SlackSecretName      string `envconfig:"SLACK_SECRET_NAME" default:""`
	AminWebURL           string `envconfig:"ADMIN_WEB_URL" default:""`
	UserWebURL           string `envconfig:"USER_WEB_URL" default:""`
}

//nolint:revive
func NewApp() *app {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "gateway user",
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
		return fmt.Errorf("user: failed to load environment: %w", err)
	}

	// Loggerの設定
	logger, err := log.NewLogger(log.WithLogLevel(a.LogLevel), log.WithOutput(a.LogPath))
	if err != nil {
		return fmt.Errorf("user: failed to new logger: %w", err)
	}
	defer logger.Sync() //nolint:errcheck

	// 依存関係の解決
	if err := a.inject(ctx, logger); err != nil {
		logger.Error("Failed to new registry", zap.Error(err))
		return err
	}

	// HTTP Serverの設定
	rt := a.newRouter(logger)
	hs := http.NewHTTPServer(rt, a.Port)

	// Metrics Serverの設定
	ms := http.NewMetricsServer(a.MetricsPort)

	// Serverの起動
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if err = ms.Serve(); err != nil {
			logger.Error("Failed to run metrics server", zap.Error(err))
		}
		return
	})
	eg.Go(func() (err error) {
		if err = hs.Serve(); err != nil {
			logger.Error("Failed to run http server", zap.Error(err))
		}
		return
	})
	logger.Info("Started server", zap.Int64("port", a.Port))

	// シグナル検知設定
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-ectx.Done():
		logger.Error("Done context", zap.Error(ectx.Err()))
	case signal := <-signalCh:
		logger.Info("Received signal", zap.Any("signal", signal))
		delay := time.Duration(a.ShutdownDelaySec) * time.Second
		logger.Info("Pre-shutdown", zap.Duration("delay", delay))
		time.Sleep(delay)
	}

	// Serverの停止
	logger.Info("Shutdown...")
	if err = hs.Stop(ectx); err != nil {
		logger.Error("Failed to stopeed http server", zap.Error(err))
		return err
	}
	if err = ms.Stop(ectx); err != nil {
		logger.Error("Failed to stopeed metrics server", zap.Error(err))
		return err
	}
	a.waitGroup.Wait()
	return eg.Wait()
}
