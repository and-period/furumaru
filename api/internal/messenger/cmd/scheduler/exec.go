package cmd

import (
	"context"

	"github.com/and-period/furumaru/api/pkg/jst"
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

	// Job実行に必要な引数の生成
	target := jst.Now()
	if conf.TargetDatetime != "" {
		target, err = jst.Parse("2006-01-02 15:04:05", conf.TargetDatetime)
		if err != nil {
			logger.Error("Failed to parse target datetime", zap.Error(err), zap.String("target", conf.TargetDatetime))
			return err
		}
	}

	// Jobの起動
	logger.Info("Started")
	switch conf.RunMethod {
	case "lambda":
		logger.Info("Started Lambda function")
		lambda.Start(reg.job.Lambda)
	default:
		logger.Info("Started manual function")
		err = reg.job.Run(ctx, target)
	}

	defer logger.Info("Finished...")
	reg.waitGroup.Wait()
	return nil
}
