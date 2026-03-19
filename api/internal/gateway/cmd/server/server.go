package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/cmd/admin"
	"github.com/and-period/furumaru/api/internal/gateway/cmd/user"
	"github.com/and-period/furumaru/api/pkg/http"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type app struct {
	*cobra.Command
	AppName         string `default:"gateway"        envconfig:"APP_NAME"`
	Environment     string `default:"none"           envconfig:"ENV"`
	AdminPort       int64  `default:"8080"           envconfig:"ADMIN_PORT"`
	UserPort        int64  `default:"8081"           envconfig:"USER_PORT"`
	MetricsPort     int64  `default:"9090"           envconfig:"METRICS_PORT"`
	ShutdownDelaySec int64 `default:"20"             envconfig:"SHUTDOWN_DELAY_SEC"`
	LogLevel        string `default:"info"           envconfig:"LOG_LEVEL"`
	SentryDsn       string `default:""               envconfig:"SENTRY_DSN"`
	SentrySecretName string `default:""              envconfig:"SENTRY_SECRET_NAME"`
}

func NewApp() *app {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "gateway server (admin + user)",
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
		return fmt.Errorf("server: failed to load environment: %w", err)
	}

	// ログの設定（1回だけ）
	logOpts := []log.Option{
		log.WithLogLevel(a.LogLevel),
		log.WithSentryDSN(a.SentryDsn),
		log.WithSentryServerName(a.AppName),
		log.WithSentryEnvironment(a.Environment),
		log.WithSentryLevel("error"),
	}
	logFlush, err := log.Start(ctx, logOpts...)
	if err != nil {
		return fmt.Errorf("server: failed to start logger: %w", err)
	}
	defer logFlush()

	// admin gateway のビルド
	adminResult, err := admin.NewApp().Build(ctx)
	if err != nil {
		return fmt.Errorf("server: failed to build admin gateway: %w", err)
	}

	// user gateway のビルド
	userResult, err := user.NewApp().Build(ctx)
	if err != nil {
		return fmt.Errorf("server: failed to build user gateway: %w", err)
	}

	// HTTP Serverの設定
	adminServer := http.NewHTTPServer(adminResult.Router, a.AdminPort)
	userServer := http.NewHTTPServer(userResult.Router, a.UserPort)

	// Metrics Serverの設定
	metricsServer := http.NewMetricsServer(a.MetricsPort)

	// Serverの起動
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if err = metricsServer.Serve(); err != nil {
			slog.Warn("Failed to run metrics server", log.Error(err))
		}
		return
	})
	eg.Go(func() (err error) {
		if err = adminServer.Serve(); err != nil {
			slog.Warn("Failed to run admin http server", log.Error(err))
		}
		return
	})
	eg.Go(func() (err error) {
		if err = userServer.Serve(); err != nil {
			slog.Warn("Failed to run user http server", log.Error(err))
		}
		return
	})
	// admin の Sync goroutine 起動
	for _, h := range adminResult.Handlers {
		h := h
		eg.Go(func() (err error) {
			if err = h.Sync(ectx); err != nil {
				slog.Warn("Failed to sync admin handler", log.Error(err))
			}
			return
		})
	}
	slog.Info("Started server",
		slog.Int64("adminPort", a.AdminPort),
		slog.Int64("userPort", a.UserPort),
		slog.Int64("metricsPort", a.MetricsPort),
	)
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
	if err := adminServer.Stop(ectx); err != nil {
		slog.Error("Failed to stop admin http server", log.Error(err))
		return err
	}
	if err := userServer.Stop(ectx); err != nil {
		slog.Error("Failed to stop user http server", log.Error(err))
		return err
	}
	if err := metricsServer.Stop(ectx); err != nil {
		slog.Error("Failed to stop metrics server", log.Error(err))
		return err
	}
	adminResult.WaitGroup.Wait()
	userResult.WaitGroup.Wait()
	return eg.Wait()
}
