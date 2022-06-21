package cmd

import (
	"context"
	"sync"

	"github.com/and-period/furumaru/api/internal/messenger/worker"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

type app struct {
	logger    *zap.Logger
	worker    worker.Worker
	waitGroup *sync.WaitGroup
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
	//nolint:errcheck
	defer app.logger.Sync()

	app.logger.Info("Started")
	switch conf.RunMethod {
	case "lambda":
		app.lambda()
	default:
		app.logger.Info("Not runnning...", zap.String("method", conf.RunMethod))
		return nil
	}
	app.waitGroup.Wait()
	app.logger.Info("Finished")
	return nil
}

func newApp(ctx context.Context, conf *config) (*app, error) {
	// 依存関係の解決
	reg, err := newRegistry(ctx, conf)
	if err != nil {
		return nil, err
	}

	return &app{
		logger:    reg.logger,
		worker:    reg.worker,
		waitGroup: reg.waitGroup,
	}, nil
}

func (a *app) lambda() {
	lambda.Start(a.worker.Lambda)
}
