package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/and-period/furumaru/api/pkg/http"
	"github.com/and-period/furumaru/api/pkg/log"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type app struct {
	logger    *zap.Logger
	server    http.Server
	metrics   http.Server
	waigGroup *sync.WaitGroup
}

func Exec() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := newConfig()
	if err != nil {
		return err
	}

	app, err := newApp(ctx, conf)
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
	app.waigGroup.Wait()
	return eg.Wait()
}

func newApp(ctx context.Context, conf *config) (*app, error) {
	// 依存関係の解決
	reg, err := newRegistry(ctx, conf)
	if err != nil {
		return nil, err
	}

	// HTTP Serverの設定
	lm, err := log.NewGinMiddleware(log.WithLogLevel(conf.LogLevel), log.WithOutput(conf.LogPath))
	if err != nil {
		return nil, err
	}

	rt := newRouter(reg, lm)
	hs := http.NewHTTPServer(rt, conf.Port)

	// Metrics Serverの設定
	ms := http.NewMetricsServer(conf.MetricsPort)

	return &app{
		logger:    reg.logger,
		server:    hs,
		metrics:   ms,
		waigGroup: reg.waitGroup,
	}, nil
}
