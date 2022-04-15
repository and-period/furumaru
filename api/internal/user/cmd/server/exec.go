package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	apgrpc "github.com/and-period/marche/api/pkg/grpc"
	aphttp "github.com/and-period/marche/api/pkg/http"
	"github.com/and-period/marche/api/pkg/log"
	"github.com/and-period/marche/api/proto/user"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type app struct {
	logger  *zap.Logger
	server  apgrpc.Server
	metrics aphttp.Server
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
	if err = app.metrics.Stop(ectx); err != nil {
		return err
	}
	app.server.Stop()
	return eg.Wait()
}

func newApp(ctx context.Context, conf *config) (*app, error) {
	// Loggerの設定
	logger, err := log.NewLogger(log.WithLogLevel(conf.LogLevel), log.WithOutput(conf.LogPath))
	if err != nil {
		return nil, err
	}

	// 依存関係の解決
	reg, err := newRegistry(ctx, conf, withLogger(logger))
	if err != nil {
		return nil, err
	}

	// gRPC Serverの設定
	gRPCOpts := apgrpc.NewGRPCOptions(apgrpc.WithLogger(logger))

	s := grpc.NewServer(gRPCOpts...)
	user.RegisterUserServiceServer(s, reg.userServer)

	gs, err := apgrpc.NewGRPCServer(s, conf.Port)
	if err != nil {
		return nil, err
	}

	// Metrics Serverの設定
	ms := aphttp.NewMetricsServer(conf.MetricsPort)

	return &app{
		logger:  logger,
		server:  gs,
		metrics: ms,
	}, nil
}
