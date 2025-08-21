package user

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway"
	"github.com/and-period/furumaru/api/pkg/http"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/slack"
	"github.com/kelseyhightower/envconfig"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type app struct {
	*cobra.Command
	debugMode                    bool
	waitGroup                    *sync.WaitGroup
	slack                        slack.Client
	newRelic                     *newrelic.Application
	v1                           gateway.Handler
	facility                     gateway.Handler
	AppName                      string  `default:"user-gateway"   envconfig:"APP_NAME"`
	Environment                  string  `default:"none"           envconfig:"ENV"`
	Port                         int64   `default:"8080"           envconfig:"PORT"`
	MetricsPort                  int64   `default:"9090"           envconfig:"METRICS_PORT"`
	ShutdownDelaySec             int64   `default:"20"             envconfig:"SHUTDOWN_DELAY_SEC"`
	LogPath                      string  `default:""               envconfig:"LOG_PATH"`
	LogLevel                     string  `default:"info"           envconfig:"LOG_LEVEL"`
	TraceSampleRate              float64 `default:"0.0"            envconfig:"TRACE_SAMPLE_RATE"`
	DBTimeZone                   string  `default:"Asia/Tokyo"     envconfig:"DB_TIMEZONE"`
	TiDBHost                     string  `default:"127.0.0.1"      envconfig:"TIDB_HOST"`
	TiDBPort                     string  `default:"4000"           envconfig:"TIDB_PORT"`
	TiDBUsername                 string  `default:""               envconfig:"TIDB_USERNAME"`
	TiDBPassword                 string  `default:""               envconfig:"TIDB_PASSWORD"`
	TiDBSecretName               string  `default:""               envconfig:"TIDB_SECRET_NAME"`
	GinMode                      string  `default:"release"        envconfig:"GIN_MODE"`
	NewRelicLicense              string  `default:""               envconfig:"NEW_RELIC_LICENSE"`
	NewRelicSecretName           string  `default:""               envconfig:"NEW_RELIC_SECRET_NAME"`
	SentryDsn                    string  `default:""               envconfig:"SENTRY_DSN"`
	SentrySecretName             string  `default:""               envconfig:"SENTRY_SECRET_NAME"`
	AWSRegion                    string  `default:"ap-northeast-1" envconfig:"AWS_REGION"`
	S3Bucket                     string  `default:""               envconfig:"S3_BUCKET"`
	S3TmpBucket                  string  `default:""               envconfig:"S3_TMP_BUCKET"`
	CognitoUserPoolID            string  `default:""               envconfig:"COGNITO_USER_POOL_ID"`
	CognitoUserAuthDomain        string  `default:""               envconfig:"COGNITO_USER_AUTH_DOMAIN"`
	CognitoUserGoogleRedirectURL string  `default:""               envconfig:"COGNITO_USER_GOOGLE_REDIRECT_URL"`
	CognitoUserLINERedirectURL   string  `default:""               envconfig:"COGNITO_USER_LINE_REDIRECT_URL"`
	CognitoUserClientID          string  `default:""               envconfig:"COGNITO_USER_CLIENT_ID"`
	SQSQueueURL                  string  `default:""               envconfig:"SQS_QUEUE_URL"`
	SQSMockEnabled               bool    `default:"false"          envconfig:"SQS_MOCK_ENABLED"`
	KomojuHost                   string  `default:""               envconfig:"KOMOJU_HOST"`
	KomojuClientID               string  `default:""               envconfig:"KOMOJU_CLIENT_ID"`
	KomojuClientPassword         string  `default:""               envconfig:"KOMOJU_CLIENT_PASSWORD"`
	KomojuSecretName             string  `default:""               envconfig:"KOMOJU_SECRET_NAME"`
	GoogleSecretName             string  `default:""               envconfig:"GOOGLE_SECRET_NAME"`
	GoogleMapsPlatformAPIKey     string  `default:""               envconfig:"GOOGLE_MAPS_PLATFORM_API_KEY"`
	CheckoutAutoCaptured         bool    `default:"false"          envconfig:"CHECKOUT_AUTO_CAPTURED"`
	SlackAPIToken                string  `default:""               envconfig:"SLACK_API_TOKEN"`
	SlackChannelID               string  `default:""               envconfig:"SLACK_CHANNEL_ID"`
	SlackSecretName              string  `default:""               envconfig:"SLACK_SECRET_NAME"`
	CookieBaseDomain             string  `default:""               envconfig:"COOKIE_BASE_DOMAIN"`
	AminWebURL                   string  `default:""               envconfig:"ADMIN_WEB_URL"`
	UserWebURL                   string  `default:""               envconfig:"USER_WEB_URL"`
}

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
		return fmt.Errorf("user: failed to new registry: %w", err)
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
			slog.Warn("Failed to run metrics server", log.Error(err))
		}
		return
	})
	eg.Go(func() (err error) {
		if err = hs.Serve(); err != nil {
			slog.Warn("Failed to run http server", log.Error(err))
		}
		return
	})
	slog.Info("Started server", slog.Int64("port", a.Port))
	defer func() {
		if r := recover(); r != nil {
			stackTrace := make([]byte, 1024)
			runtime.Stack(stackTrace, true)
			slog.Error("Occurred panic", slog.Any("value", r), slog.String("stackTrace", string(stackTrace)))
		}
	}()

	// シグナル検知設定
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-ectx.Done():
		slog.Warn("Done context", log.Error(ectx.Err()))
	case signal := <-signalCh:
		slog.Info("Received signal", slog.Any("signal", signal))
		delay := time.Duration(a.ShutdownDelaySec) * time.Second
		slog.Info("Pre-shutdown", slog.Duration("delay", delay))
		time.Sleep(delay)
	}

	// Serverの停止
	slog.Info("Shutdown...")
	if err := hs.Stop(ectx); err != nil {
		slog.Error("Failed to stopeed http server", log.Error(err))
		return err
	}
	if err := ms.Stop(ectx); err != nil {
		slog.Error("Failed to stopeed metrics server", log.Error(err))
		return err
	}
	a.waitGroup.Wait()
	return eg.Wait()
}
