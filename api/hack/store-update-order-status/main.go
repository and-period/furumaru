// 注文履歴のステータスを更新します
//
//	usage: go run ./main.go \
//	 -db-host='127.0.0.1' -db-port='3316' \
//	 -db-username='root' -db-password='12345678'
package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/database/tidb"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/log"
	apmysql "github.com/and-period/furumaru/api/pkg/mysql"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	dbhost     string
	dbport     string
	dbusername string
	dbpassword string
)

type app struct {
	mysql  *apmysql.Client
	db     *database.Database
	logger *zap.Logger
}

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
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	logger, err := log.NewLogger(log.WithLogLevel("debug"))
	if err != nil {
		return err
	}
	defer logger.Sync() //nolint:errcheck

	params := &apmysql.Params{
		Socket:   "tcp",
		Host:     dbhost,
		Port:     dbport,
		Database: "stores",
		Username: dbusername,
		Password: dbpassword,
	}
	db, err := apmysql.NewClient(params, apmysql.WithLogger(logger))
	if err != nil {
		return err
	}

	app := &app{
		mysql:  db,
		db:     tidb.NewDatabase(db),
		logger: logger,
	}

	orders, err := app.db.Order.List(ctx, &database.ListOrdersParams{})
	if err != nil {
		app.logger.Error("Failed to list orders", zap.Error(err))
		return err
	}
	update := func(ctx context.Context, tx *gorm.DB, order *entity.Order) error {
		params := map[string]interface{}{
			"status": app.getOrderStatus(order),
		}
		stmt := tx.WithContext(ctx).Table("orders").Where("id = ?", order.ID)
		return stmt.Updates(params).Error
	}
	err = app.mysql.Transaction(ctx, func(tx *gorm.DB) error {
		for i := range orders {
			if err := update(ctx, tx, orders[i]); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		app.logger.Error("Failed to update order status", zap.Error(err))
		return err
	}
	app.logger.Info("Finished update order status")
	return nil
}

func (a *app) getOrderStatus(order *entity.Order) entity.OrderStatus {
	if order == nil {
		return entity.OrderStatusUnknown
	}
	switch order.OrderPayment.Status {
	case entity.PaymentStatusPending:
		return entity.OrderStatusUnpaid
	case entity.PaymentStatusAuthorized:
		return entity.OrderStatusWaiting
	case entity.PaymentStatusCaptured:
		if !order.Fulfilled() {
			return entity.OrderStatusPreparing
		}
		if order.CompletedAt.IsZero() {
			return entity.OrderStatusShipped
		}
		return entity.OrderStatusCompleted
	case entity.PaymentStatusCanceled:
		return entity.OrderStatusCanceled
	case entity.PaymentStatusRefunded:
		return entity.OrderStatusRefunded
	case entity.PaymentStatusFailed:
		return entity.OrderStatusFailed
	default:
		return entity.OrderStatusUnknown
	}
}
