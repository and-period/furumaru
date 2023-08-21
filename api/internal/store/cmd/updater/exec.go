package cmd

import (
	"context"
	"errors"
	"fmt"

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

	// Jobの起動
	logger.Info("Started")
	switch conf.RunMethod {
	case "lambda":
		switch conf.RunType {
		case "CREATE":
			lambda.StartWithOptions(reg.creator.Lambda, lambda.WithContext(ctx))
		case "REMOVE":
			return errors.New("cmd: not implemented")
		default:
			return fmt.Errorf("cmd: unknown scheduler type. type=%s", conf.RunType)
		}
	default:
		err = errors.New("not implemented")
	}

	defer logger.Info("Finished...")
	reg.waitGroup.Wait()
	return err
}
