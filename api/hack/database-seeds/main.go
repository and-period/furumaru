// データベースに初期データの投入をします
//
//	usage: go run ./hack/database-seeds/main.go \
//	 -db-host='127.0.0.1' -db-port='3316' \
//	 -db-username='root' -db-password='12345678' \
//	 -src-dir='./hack/database-seeds/master'
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/and-period/furumaru/api/hack/database-seeds/common"
	"github.com/and-period/furumaru/api/hack/database-seeds/messenger"
	"github.com/and-period/furumaru/api/hack/database-seeds/store"
	"github.com/and-period/furumaru/api/hack/database-seeds/user"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/secret"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	logLevel     string
	awsRegion    string
	dbSecretName string
	dbTimeZone   string
	srcDir       string

	tidbHost     string
	tidbPort     string
	tidbUsername string
	tidbPassword string
)

type app struct {
	logger    *zap.Logger
	messenger common.Client
	store     common.Client
	user      common.Client
}

func main() {
	startedAt := jst.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	app, err := setup(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup: %v\n", err)
		os.Exit(1)
	}

	if err := app.run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run: %v\n", err)
		os.Exit(1)
	}

	endAt := jst.Now()

	const format = "2006-01-02 15:04:05"
	fmt.Printf("Done: %s -> %s (%s)\n",
		jst.Format(startedAt, format), jst.Format(endAt, format),
		endAt.Sub(startedAt).Truncate(time.Second).String(),
	)
}

func setup(ctx context.Context) (*app, error) {
	flag.StringVar(&awsRegion, "aws-region", "ap-northeast-1", "AWS region")
	flag.StringVar(&logLevel, "log-level", "debug", "log level")
	flag.StringVar(&dbSecretName, "db-secret-name", "", "AWS Secrets Manager secret name for TiDB")
	flag.StringVar(&dbTimeZone, "db-timezone", "Asia/Tokyo", "TiDB timezone")
	flag.StringVar(&srcDir, "src-dir", "./hack/database-seeds/master", "source directory")
	flag.Parse()

	if dbSecretName == "" {
		return nil, fmt.Errorf("db-secret-name is required")
	}

	// AWS SDKの設定
	awscfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(awsRegion))
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	// AWS Secrets Managerの設定
	secret := secret.NewClient(awscfg)

	secrets, err := secret.Get(ctx, dbSecretName)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}
	tidbHost = secrets["host"]
	tidbPort = secrets["port"]
	tidbUsername = secrets["username"]
	tidbPassword = secrets["password"]

	// Loggerの設定
	logger, err := log.NewLogger(log.WithLogLevel(logLevel))
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// モジュールごとのクライアント生成
	params := &common.Params{
		Logger:     logger,
		DBHost:     tidbHost,
		DBPort:     tidbPort,
		DBUsername: tidbUsername,
		DBPassword: tidbPassword,
		SrcDir:     srcDir,
	}
	messenger, err := messenger.NewClient(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create messenger client: %w", err)
	}
	store, err := store.NewClient(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create store client: %w", err)
	}
	user, err := user.NewClient(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user client: %w", err)
	}

	app := &app{
		logger:    logger,
		messenger: messenger,
		store:     store,
		user:      user,
	}
	return app, nil
}

func (a *app) run(ctx context.Context) error {
	a.logger.Info("Database seeds will begin")
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return a.user.Execute(ectx)
	})
	eg.Go(func() error {
		return a.store.Execute(ectx)
	})
	eg.Go(func() error {
		return a.messenger.Execute(ectx)
	})
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to execute database seeds: %w", err)
	}
	a.logger.Info("Database seeds has completed")
	return nil
}
