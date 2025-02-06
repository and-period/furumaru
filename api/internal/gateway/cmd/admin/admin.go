package admin

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	v1 "github.com/and-period/furumaru/api/internal/gateway/admin/v1/handler"
	khandler "github.com/and-period/furumaru/api/internal/gateway/komoju/handler"
	"github.com/and-period/furumaru/api/pkg/http"
	"github.com/and-period/furumaru/api/pkg/slack"
	"github.com/kelseyhightower/envconfig"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type app struct {
	*cobra.Command
	debugMode                         bool
	logger                            *zap.Logger
	waitGroup                         *sync.WaitGroup
	slack                             slack.Client
	newRelic                          *newrelic.Application
	v1                                v1.Handler
	komoju                            khandler.Handler
	AppName                           string   `default:"admin-gateway"  envconfig:"APP_NAME"`
	Environment                       string   `default:"none"           envconfig:"ENV"`
	Port                              int64    `default:"8080"           envconfig:"PORT"`
	MetricsPort                       int64    `default:"9090"           envconfig:"METRICS_PORT"`
	ShutdownDelaySec                  int64    `default:"20"             envconfig:"SHUTDOWN_DELAY_SEC"`
	LogPath                           string   `default:""               envconfig:"LOG_PATH"`
	LogLevel                          string   `default:"info"           envconfig:"LOG_LEVEL"`
	TraceSampleRate                   float64  `default:"0.0"            envconfig:"TRACE_SAMPLE_RATE"`
	DBTimeZone                        string   `default:"Asia/Tokyo"     envconfig:"DB_TIMEZONE"`
	TiDBHost                          string   `default:"127.0.0.1"      envconfig:"TIDB_HOST"`
	TiDBPort                          string   `default:"4000"           envconfig:"TIDB_PORT"`
	TiDBUsername                      string   `default:""               envconfig:"TIDB_USERNAME"`
	TiDBPassword                      string   `default:""               envconfig:"TIDB_PASSWORD"`
	TiDBSecretName                    string   `default:""               envconfig:"TIDB_SECRET_NAME"`
	GinMode                           string   `default:"release"        envconfig:"GIN_MODE"`
	NewRelicLicense                   string   `default:""               envconfig:"NEW_RELIC_LICENSE"`
	NewRelicSecretName                string   `default:""               envconfig:"NEW_RELIC_SECRET_NAME"`
	SentryDsn                         string   `default:""               envconfig:"SENTRY_DSN"`
	SentrySecretName                  string   `default:""               envconfig:"SENTRY_SECRET_NAME"`
	AWSRegion                         string   `default:"ap-northeast-1" envconfig:"AWS_REGION"`
	S3Bucket                          string   `default:""               envconfig:"S3_BUCKET"`
	S3TmpBucket                       string   `default:""               envconfig:"S3_TMP_BUCKET"`
	CognitoAdminPoolID                string   `default:""               envconfig:"COGNITO_ADMIN_POOL_ID"`
	CognitoAdminClientID              string   `default:""               envconfig:"COGNITO_ADMIN_CLIENT_ID"`
	CognitoAdminAuthDomain            string   `default:""               envconfig:"COGNITO_ADMIN_AUTH_DOMAIN"`
	CognitoAdminGoogleRedirectURL     string   `default:""               envconfig:"COGNITO_ADMIN_GOOGLE_REDIRECT_URL"`
	CognitoAdminLINERedirectURL       string   `default:""               envconfig:"COGNITO_ADMIN_LINE_REDIRECT_URL"`
	CognitoUserPoolID                 string   `default:""               envconfig:"COGNITO_USER_POOL_ID"`
	CognitoUserClientID               string   `default:""               envconfig:"COGNITO_USER_CLIENT_ID"`
	SQSMessengerQueueURL              string   `default:""               envconfig:"SQS_MESSENGER_QUEUE_URL"`
	SQSMediaQueueURL                  string   `default:""               envconfig:"SQS_MEDIA_QUEUE_URL"`
	SQSMockEnabled                    bool     `default:"false"          envconfig:"SQS_MOCK_ENABLED"`
	KomojuHost                        string   `default:""               envconfig:"KOMOJU_HOST"`
	KomojuClientID                    string   `default:""               envconfig:"KOMOJU_CLIENT_ID"`
	KomojuClientPassword              string   `default:""               envconfig:"KOMOJU_CLIENT_PASSWORD"`
	KomojuSecretName                  string   `default:""               envconfig:"KOMOJU_SECRET_NAME"`
	GoogleClientID                    string   `default:""               envconfig:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret                string   `default:""               envconfig:"GOOGLE_CLIENT_SECRET"`
	GoogleSecretName                  string   `default:""               envconfig:"GOOGLE_SECRET_NAME"`
	GoogleMapsPlatformAPIKey          string   `default:""               envconfig:"GOOGLE_MAPS_PLATFORM_API_KEY"`
	YoutubeAuthCallbackURL            string   `default:""               envconfig:"YOUTUBE_AUTH_CALLBACK_URL"`
	BatchMediaUpdateArchiveDefinition string   `default:""               envconfig:"BATCH_MEDIA_UPDATE_ARCHIVE_DEFINITION"`
	BatchMediaUpdateArchiveQueue      string   `default:""               envconfig:"BATCH_MEDIA_UPDATE_ARCHIVE_QUEUE"`
	AdminWebURL                       string   `default:""               envconfig:"ADMIN_WEB_URL"`
	UserWebURL                        string   `default:""               envconfig:"USER_WEB_URL"`
	AssetsURL                         string   `default:""               envconfig:"ASSETS_URL"`
	SlackAPIToken                     string   `default:""               envconfig:"SLACK_API_TOKEN"`
	SlackChannelID                    string   `default:""               envconfig:"SLACK_CHANNEL_ID"`
	SlackSecretName                   string   `default:""               envconfig:"SLACK_SECRET_NAME"`
	DefaultAdministratorGroupIDs      []string `default:""               envconfig:"DEFAULT_ADMINISTRATOR_GROUPS"`
	DefaultCoordinatorGroupIDs        []string `default:""               envconfig:"DEFAULT_COORDINATOR_GROUPS"`
	DefaultProducerGroupIDs           []string `default:""               envconfig:"DEFAULT_PRODUCER_GROUPS"`
}

//nolint:revive
func NewApp() *app {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "gateway admin",
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
		return fmt.Errorf("admin: failed to load environment: %w", err)
	}

	// 依存関係の解決
	if err := a.inject(ctx); err != nil {
		return fmt.Errorf("admin: failed to new registry: %w", err)
	}
	defer a.logger.Sync() //nolint:errcheck

	if err := a.v1.Setup(ctx); err != nil {
		return fmt.Errorf("admin: failed to setup http server: %w", err)
	}

	// HTTP Serverの設定
	rt := a.newRouter()
	hs := http.NewHTTPServer(rt, a.Port)

	// Metrics Serverの設定
	ms := http.NewMetricsServer(a.MetricsPort)

	// Serverの起動
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if err = ms.Serve(); err != nil {
			a.logger.Warn("Failed to run metrics server", zap.Error(err))
		}
		return
	})
	eg.Go(func() (err error) {
		if err = hs.Serve(); err != nil {
			a.logger.Warn("Failed to run http server", zap.Error(err))
		}
		return
	})
	eg.Go(func() (err error) {
		if err = a.v1.Sync(ectx); err != nil {
			a.logger.Warn("Failed to sync v1 handler", zap.Error(err))
		}
		return
	})
	a.logger.Info("Started server", zap.Int64("port", a.Port))
	defer func() {
		if r := recover(); r != nil {
			stackTrace := make([]byte, 1024)
			runtime.Stack(stackTrace, true)
			a.logger.Error("Occurred panic", zap.Any("value", r), zap.String("stackTrace", string(stackTrace)))
		}
	}()

	// シグナル検知設定
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-ectx.Done():
		a.logger.Warn("Done context", zap.Error(ectx.Err()))
	case signal := <-signalCh:
		a.logger.Info("Received signal", zap.Any("signal", signal))
		delay := time.Duration(a.ShutdownDelaySec) * time.Second
		a.logger.Info("Pre-shutdown", zap.Duration("delay", delay))
		time.Sleep(delay)
	}

	// Serverの停止
	a.logger.Info("Shutdown...")
	if err := hs.Stop(ectx); err != nil {
		a.logger.Error("Failed to stopeed http server", zap.Error(err))
		return err
	}
	if err := ms.Stop(ectx); err != nil {
		a.logger.Error("Failed to stopeed metrics server", zap.Error(err))
		return err
	}
	a.waitGroup.Wait()
	return eg.Wait()
}
