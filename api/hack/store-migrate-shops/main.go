// コーディネータの情報を元に店舗を作成する
//
//	usage: go run ./main.go \
//	 -tidb-host='localhost' \
//	 -tidb-username='root' \
//	 -tidb-password='password'
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	sdb "github.com/and-period/furumaru/api/internal/store/database"
	stidb "github.com/and-period/furumaru/api/internal/store/database/tidb"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	udb "github.com/and-period/furumaru/api/internal/user/database"
	utidb "github.com/and-period/furumaru/api/internal/user/database/tidb"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"go.uber.org/zap"
)

var (
	tidbHost     string
	tidbPort     string
	tidbUsername string
	tidbPassword string
)

type app struct {
	logger *zap.Logger
	store  *sdb.Database
	user   *udb.Database
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

func setup(_ context.Context) (*app, error) {
	flag.StringVar(&tidbHost, "tidb-host", "localhost", "TiDB host")
	flag.StringVar(&tidbPort, "tidb-port", "4000", "TiDB port")
	flag.StringVar(&tidbUsername, "tidb-username", "root", "TiDB username")
	flag.StringVar(&tidbPassword, "tidb-password", "", "TiDB password")
	flag.Parse()

	// Loggerの設定
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// TiDBの設定
	storedb, err := mysql.NewTiDBClient(&mysql.Params{
		Host:     tidbHost,
		Port:     tidbPort,
		Database: "stores",
		Username: tidbUsername,
		Password: tidbPassword,
	}, mysql.WithNow(jst.Now), mysql.WithLocation(jst.Location()))
	if err != nil {
		return nil, fmt.Errorf("failed to create tidb client for store: %w", err)
	}

	userdb, err := mysql.NewTiDBClient(&mysql.Params{
		Host:     tidbHost,
		Port:     tidbPort,
		Database: "users",
		Username: tidbUsername,
		Password: tidbPassword,
	}, mysql.WithNow(jst.Now), mysql.WithLocation(jst.Location()))
	if err != nil {
		return nil, fmt.Errorf("failed to create tidb client for user: %w", err)
	}

	app := &app{
		logger: logger,
		store:  stidb.NewDatabase(storedb),
		user:   utidb.NewDatabase(userdb),
	}
	return app, nil
}

func (a *app) run(ctx context.Context) error {
	a.logger.Info("start migration")

	cparams := &udb.ListCoordinatorsParams{}
	coordinators, err := a.user.Coordinator.List(ctx, cparams)
	if err != nil {
		return fmt.Errorf("failed to list coordinators: %w", err)
	}
	if len(coordinators) == 0 {
		return nil
	}
	a.logger.Info("coordinators", zap.Int("count", len(coordinators)))

	for _, coordinator := range coordinators {
		shop, err := a.store.Shop.GetByCoordinatorID(ctx, coordinator.ID)
		if err != nil && !errors.Is(err, sdb.ErrNotFound) {
			return fmt.Errorf("failed to get shop by coordinator id: %w", err)
		}

		if shop == nil {
			params := &sentity.ShopParams{
				CoordinatorID:  coordinator.ID,
				ProductTypeIDs: coordinator.ProductTypeIDs,
				Name:           coordinator.MarcheName,
			}
			shop = sentity.NewShop(params)
			if err := a.store.Shop.Create(ctx, shop); err != nil {
				return fmt.Errorf("failed to create shop: %w", err)
			}
		}

		pparams := &udb.ListProducersParams{
			CoordinatorID: coordinator.ID,
		}
		producers, err := a.user.Producer.List(ctx, pparams)
		if err != nil {
			return fmt.Errorf("failed to list producers: %w", err)
		}
		if len(producers) == 0 {
			continue
		}

		for _, producer := range producers {
			if err := a.store.Shop.RelateProducer(ctx, shop.ID, producer.ID); err != nil {
				return fmt.Errorf("failed to relate producer: %w", err)
			}
		}
	}

	a.logger.Info("end migration")
	return nil
}
