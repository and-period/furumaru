package resizer

import (
	"context"
	"errors"

	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

func (a *app) run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 環境変数の読み込み
	conf, err := newConfig()
	if err != nil {
		return err
	}

	// Loggerの設定
	logger, err := log.NewLogger(log.WithLogLevel(conf.LogLevel), log.WithOutput(conf.LogPath))
	if err != nil {
		return err
	}
	defer logger.Sync() //nolint:errcheck

	// 依存関係の解決
	reg, err := newRegistry(ctx, conf, logger)
	if err != nil {
		logger.Error("Failed to new registry", zap.Error(err))
		return err
	}

	// Workerの起動
	logger.Info("Started")
	switch conf.RunMethod {
	case "lambda":
		lambda.StartWithOptions(reg.resizer.Lambda, lambda.WithContext(ctx))
	default:
		err = errors.New("not implemented")
	}

	// Workerの停止
	logger.Info("Shutdown...")
	reg.waitGroup.Wait()
	return err
}
