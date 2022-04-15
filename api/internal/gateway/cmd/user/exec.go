package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/and-period/marche/api/pkg/cors"
	"github.com/and-period/marche/api/pkg/http"
	"github.com/and-period/marche/api/pkg/log"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type app struct {
	logger  *zap.Logger
	server  http.Server
	metrics http.Server
}

func Exec() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := newConfig()
	if err != nil {
		return err
	}

	app, err := newApp(conf)
	if err != nil {
		return err
	}

	// Serverの起動
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		err = app.metrics.Serve()
		if err != nil {
			app.logger.Error("Failed to run metrics server", zap.Error(err))
		}
		return
	})
	eg.Go(func() (err error) {
		err = app.server.Serve()
		if err != nil {
			app.logger.Error("Failed to run gRPC server", zap.Error(err))
		}
		return
	})
	app.logger.Info("Started server", zap.Int64("port", conf.Port))

	// シグナル検知設定
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-ectx.Done():
		app.logger.Error("Done context", zap.Error(ectx.Err()))
	case <-signalCh:
		app.logger.Info("Received signal")
		delay := time.Duration(conf.ShutdownDelaySec) * time.Second
		app.logger.Info("Pre-shutdown", zap.String("delay", delay.String()))
		time.Sleep(delay)
	}

	app.logger.Info("Shutdown...")
	if err = app.server.Stop(ectx); err != nil {
		return err
	}
	if err = app.metrics.Stop(ectx); err != nil {
		return err
	}
	return eg.Wait()
}

func newApp(conf *config) (*app, error) {
	// Loggerの設定
	logger, err := log.NewLogger(log.WithLogLevel(conf.LogLevel), log.WithOutput(conf.LogPath))
	if err != nil {
		return nil, err
	}

	// 依存関係の解決
	reg, err := newRegistry(conf, withLogger(logger))
	if err != nil {
		return nil, err
	}

	// HTTP Serverの設定
	httpOpts := []gin.HandlerFunc{gin.Recovery(), gzip.Gzip(gzip.DefaultCompression)}

	cm := cors.NewGinMiddleware()
	httpOpts = append(httpOpts, cm)

	lm, err := log.NewGinMiddleware(log.WithLogLevel(conf.LogLevel), log.WithOutput(conf.LogPath))
	if err != nil {
		return nil, err
	}
	httpOpts = append(httpOpts, lm)

	rt := newRouter(reg, httpOpts...)
	hs := http.NewHTTPServer(rt, conf.Port)

	// Metrics Serverの設定
	ms := http.NewMetricsServer(conf.MetricsPort)

	return &app{
		logger:  logger,
		server:  hs,
		metrics: ms,
	}, nil
}
