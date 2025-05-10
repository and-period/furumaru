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
	"gorm.io/gorm"
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
	db     *mysql.Client
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
		db:     storedb,
	}
	return app, nil
}

func (a *app) run(ctx context.Context) error {
	a.logger.Info("start migrate shops")
	if err := a.createShops(ctx); err != nil {
		return fmt.Errorf("failed to create shops: %w", err)
	}
	a.logger.Info("finish migrate shops")

	a.logger.Info("start fill experiences")
	if err := a.fillExperiences(ctx); err != nil {
		return fmt.Errorf("failed to fill experiences: %w", err)
	}
	a.logger.Info("finish fill experiences")

	a.logger.Info("start fill orders")
	if err := a.fillOrders(ctx); err != nil {
		return fmt.Errorf("failed to fill orders: %w", err)
	}
	a.logger.Info("finish fill orders")

	a.logger.Info("start fill products")
	if err := a.fillProducts(ctx); err != nil {
		return fmt.Errorf("failed to fill products: %w", err)
	}
	a.logger.Info("finish fill products")

	a.logger.Info("start fill schedules")
	if err := a.fillSchedules(ctx); err != nil {
		return fmt.Errorf("failed to fill schedules: %w", err)
	}
	a.logger.Info("finish fill schedules")

	a.logger.Info("start fill shippings")
	if err := a.fillShippings(ctx); err != nil {
		return fmt.Errorf("failed to fill shippings: %w", err)
	}
	a.logger.Info("finish fill shippings")

	return nil
}

func (a *app) createShops(ctx context.Context) error {
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
				BusinessDays:   coordinator.BusinessDays,
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
	return nil
}

func (a *app) fillExperiences(ctx context.Context) error {
	eparams := &sdb.ListExperiencesParams{}
	experiences, err := a.store.Experience.List(ctx, eparams)
	if err != nil {
		return fmt.Errorf("failed to list experiences: %w", err)
	}
	if len(experiences) == 0 {
		return nil
	}
	a.logger.Info("experiences", zap.Int("count", len(experiences)))

	err = a.db.Transaction(ctx, func(tx *gorm.DB) error {
		sparams := &sdb.ListShopsParams{
			CoordinatorIDs: experiences.CoordinatorIDs(),
		}
		shops, err := a.store.Shop.List(ctx, sparams)
		if err != nil {
			return fmt.Errorf("failed to list shops: %w", err)
		}
		shopMap := shops.MapByCoordinatorID()

		for _, experience := range experiences {
			shop, ok := shopMap[experience.CoordinatorID]
			if !ok {
				return fmt.Errorf("shop not found: %s", experience.ID)
			}

			updates := map[string]interface{}{
				"shop_id":    shop.ID,
				"updated_at": jst.Now(),
			}
			if err := tx.WithContext(ctx).Table("experiences").Where("id = ?", experience.ID).Updates(updates).Error; err != nil {
				return fmt.Errorf("failed to update experience: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update experiences: %w", err)
	}
	return nil
}

func (a *app) fillOrders(ctx context.Context) error {
	oparams := &sdb.ListOrdersParams{}
	orders, err := a.store.Order.List(ctx, oparams)
	if err != nil {
		return fmt.Errorf("failed to list orders: %w", err)
	}
	if len(orders) == 0 {
		return nil
	}
	a.logger.Info("orders", zap.Int("count", len(orders)))

	err = a.db.Transaction(ctx, func(tx *gorm.DB) error {
		sparams := &sdb.ListShopsParams{
			CoordinatorIDs: orders.CoordinatorIDs(),
		}
		shops, err := a.store.Shop.List(ctx, sparams)
		if err != nil {
			return fmt.Errorf("failed to list shops: %w", err)
		}
		shopMap := shops.MapByCoordinatorID()

		for _, order := range orders {
			shop, ok := shopMap[order.CoordinatorID]
			if !ok {
				return fmt.Errorf("shop not found: %s", order.ID)
			}

			updates := map[string]interface{}{
				"shop_id":    shop.ID,
				"updated_at": jst.Now(),
			}
			if err := tx.WithContext(ctx).Table("orders").Where("id = ?", order.ID).Updates(updates).Error; err != nil {
				return fmt.Errorf("failed to update order: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update orders: %w", err)
	}
	return nil
}

func (a *app) fillProducts(ctx context.Context) error {
	pparams := &sdb.ListProductsParams{}
	products, err := a.store.Product.List(ctx, pparams)
	if err != nil {
		return fmt.Errorf("failed to list products: %w", err)
	}
	if len(products) == 0 {
		return nil
	}
	a.logger.Info("products", zap.Int("count", len(products)))

	err = a.db.Transaction(ctx, func(tx *gorm.DB) error {
		sparams := &sdb.ListShopsParams{
			CoordinatorIDs: products.CoordinatorIDs(),
		}
		shops, err := a.store.Shop.List(ctx, sparams)
		if err != nil {
			return fmt.Errorf("failed to list shops: %w", err)
		}
		shopMap := shops.MapByCoordinatorID()

		for _, product := range products {
			shop, ok := shopMap[product.CoordinatorID]
			if !ok {
				return fmt.Errorf("shop not found: %s", product.ID)
			}

			updates := map[string]interface{}{
				"shop_id":    shop.ID,
				"updated_at": jst.Now(),
			}
			if err := tx.WithContext(ctx).Table("products").Where("id = ?", product.ID).Updates(updates).Error; err != nil {
				return fmt.Errorf("failed to update product: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update products: %w", err)
	}
	return nil
}

func (a *app) fillSchedules(ctx context.Context) error {
	sparams := &sdb.ListSchedulesParams{}
	schedules, err := a.store.Schedule.List(ctx, sparams)
	if err != nil {
		return fmt.Errorf("failed to list schedules: %w", err)
	}
	if len(schedules) == 0 {
		return nil
	}
	a.logger.Info("schedules", zap.Int("count", len(schedules)))

	err = a.db.Transaction(ctx, func(tx *gorm.DB) error {
		sparams := &sdb.ListShopsParams{
			CoordinatorIDs: schedules.CoordinatorIDs(),
		}
		shops, err := a.store.Shop.List(ctx, sparams)
		if err != nil {
			return fmt.Errorf("failed to list shops: %w", err)
		}
		shopMap := shops.MapByCoordinatorID()

		for _, schedule := range schedules {
			shop, ok := shopMap[schedule.CoordinatorID]
			if !ok {
				return fmt.Errorf("shop not found: %s", schedule.ID)
			}

			updates := map[string]interface{}{
				"shop_id":    shop.ID,
				"updated_at": jst.Now(),
			}
			if err := tx.WithContext(ctx).Table("schedules").Where("id = ?", schedule.ID).Updates(updates).Error; err != nil {
				return fmt.Errorf("failed to update schedule: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update schedules: %w", err)
	}
	return nil
}

func (a *app) fillShippings(ctx context.Context) error {
	var shippings sentity.Shippings
	if err := a.db.DB.WithContext(ctx).Table("shippings").Find(&shippings).Error; err != nil {
		return fmt.Errorf("failed to list shippings: %w", err)
	}
	if len(shippings) == 0 {
		return nil
	}
	a.logger.Info("shippings", zap.Int("count", len(shippings)))

	err := a.db.Transaction(ctx, func(tx *gorm.DB) error {
		sparams := &sdb.ListShopsParams{
			CoordinatorIDs: shippings.CoordinatorIDs(),
		}
		shops, err := a.store.Shop.List(ctx, sparams)
		if err != nil {
			return fmt.Errorf("failed to list shops: %w", err)
		}
		shopMap := shops.MapByCoordinatorID()

		for _, shipping := range shippings {
			if shipping.ID == sentity.DefaultShippingID {
				continue
			}

			shop, ok := shopMap[shipping.CoordinatorID]
			if !ok {
				return fmt.Errorf("shop not found: %s", shipping.ID)
			}

			updates := map[string]interface{}{
				"shop_id":    shop.ID,
				"updated_at": jst.Now(),
			}
			if err := tx.WithContext(ctx).Table("shippings").Where("id = ?", shipping.ID).Updates(updates).Error; err != nil {
				return fmt.Errorf("failed to update shipping: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update shippings: %w", err)
	}
	return nil
}
