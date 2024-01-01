// データベースに初期データの投入をします
//
//	usage: go run ./hack/database-seeds/main.go \
//	 -db-host='127.0.0.1' -db-port='3316' \
//	 -db-username='root' -db-password='12345678' \
//	 -src-dir='./hack/database/seeds/master'
package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/hack/database-seeds/common"
	"github.com/and-period/furumaru/api/hack/database-seeds/store"
	"github.com/and-period/furumaru/api/pkg/log"
	"go.uber.org/zap"
)

var (
	dbhost     string
	dbport     string
	dbusername string
	dbpassword string
	srcDir     string
)

func main() {
	start := time.Now()
	fmt.Println("Start..")
	if err := run(); err != nil {
		panic(err)
	}
	fmt.Printf("Done: %s\n", time.Since(start))
}

func run() error {
	flag.StringVar(&dbhost, "db-host", "mysql", "mysql server host")
	flag.StringVar(&dbport, "db-port", "3306", "mysql server port")
	flag.StringVar(&dbusername, "db-username", "root", "mysql auth username")
	flag.StringVar(&dbpassword, "db-password", "12345678", "mysql auth password")
	flag.StringVar(&srcDir, "src-dir", "./master", "source directory")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	logger, err := log.NewLogger(log.WithLogLevel("debug"))
	if err != nil {
		return err
	}
	defer logger.Sync() //nolint:errcheck

	params := &common.Params{
		Logger:     logger,
		DBHost:     dbhost,
		DBPort:     dbport,
		DBUsername: dbusername,
		DBPassword: dbpassword,
		SrcDir:     srcDir,
	}

	store, err := store.NewClient(params)
	if err != nil {
		logger.Error("Failed to create store client", zap.Error(err))
		return err
	}

	logger.Info("Database seeds will begin")
	if err := store.Execute(ctx); err != nil {
		return err
	}
	logger.Info("Database seeds has completed")
	return nil
}
