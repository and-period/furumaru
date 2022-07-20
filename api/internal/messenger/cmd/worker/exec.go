package cmd

import (
	"context"

	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

func Exec() error {
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

	txn := reg.newRelic.StartTransaction(reg.appName)
	defer txn.End()

	// Workerの起動
	logger.Info("Started")
	switch conf.RunMethod {
	case "lambda":
		logger.Info("Started Lambda function")
		lambda.StartWithOptions(reg.worker.Lambda, lambda.WithContext(ctx))
	default:
		logger.Warn("Not runnning...", zap.String("method", conf.RunMethod))
		return nil
	}

	// Workerの停止
	logger.Info("Shutdown...")
	reg.waitGroup.Wait()
	return nil
}
