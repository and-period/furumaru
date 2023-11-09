package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

const (
	orderTable            = "orders"
	orderFulfillmentTable = "order_fulfillments"
	orderItemTable        = "order_items"
	orderPaymentTable     = "order_payments"
)

type order struct {
	db  *mysql.Client
	now func() time.Time
}

func newOrder(db *mysql.Client) database.Order {
	return &order{
		db:  db,
		now: jst.Now,
	}
}

type listOrdersParams database.ListOrdersParams

func (p listOrdersParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.CoordinatorID != "" {
		stmt = stmt.Where("coordinator_id = ?", p.CoordinatorID)
	}
	return stmt
}

func (p listOrdersParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (o *order) List(ctx context.Context, params *database.ListOrdersParams, fields ...string) (entity.Orders, error) {
	var orders entity.Orders

	p := listOrdersParams(*params)

	stmt := o.db.Statement(ctx, o.db.DB, orderTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&orders).Error; err != nil {
		return nil, dbError(err)
	}
	if err := o.fill(ctx, o.db.DB, orders...); err != nil {
		return nil, dbError(err)
	}
	return orders, nil
}

func (o *order) Count(ctx context.Context, params *database.ListOrdersParams) (int64, error) {
	p := listOrdersParams(*params)

	total, err := o.db.Count(ctx, o.db.DB, &entity.Order{}, p.stmt)
	return total, dbError(err)
}

func (o *order) Get(ctx context.Context, orderID string, fields ...string) (*entity.Order, error) {
	order, err := o.get(ctx, o.db.DB, orderID, fields...)
	if err != nil {
		return nil, dbError(err)
	}
	if err := o.fill(ctx, o.db.DB, order); err != nil {
		return nil, dbError(err)
	}
	return order, nil
}

func (o *order) Aggregate(ctx context.Context, userIDs []string) (entity.AggregatedOrders, error) {
	var orders entity.AggregatedOrders

	fields := []string{
		"orders.user_id AS user_id",
		"COUNT(DISTINCT(orders.id)) AS order_count",
		"SUM(order_payments.subtotal) AS subtotal",
		"SUM(order_payments.discount) AS discount",
	}

	stmt := o.db.Statement(ctx, o.db.DB, orderTable, fields...).
		Joins("INNER JOIN order_payments ON order_payments.order_id = orders.id").
		Where("orders.user_id IN (?)", userIDs).
		Group("orders.user_id")

	err := stmt.Scan(&orders).Error
	return orders, dbError(err)
}

func (o *order) get(ctx context.Context, tx *gorm.DB, orderID string, fields ...string) (*entity.Order, error) {
	var order *entity.Order

	err := o.db.Statement(ctx, tx, orderTable, fields...).
		Where("id = ?", orderID).
		First(&order).Error
	return order, err
}

func (o *order) fill(ctx context.Context, tx *gorm.DB, orders ...*entity.Order) error {
	var (
		payments     entity.OrderPayments
		fulfillments entity.OrderFulfillments
		items        entity.OrderItems
	)

	ids := entity.Orders(orders).IDs()
	if len(ids) == 0 {
		return nil
	}

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		stmt := o.db.Statement(ectx, tx, orderPaymentTable).Where("order_id IN (?)", ids)
		return stmt.Find(&payments).Error
	})
	eg.Go(func() error {
		stmt := o.db.Statement(ectx, tx, orderFulfillmentTable).Where("order_id IN (?)", ids)
		return stmt.Find(&fulfillments).Error
	})
	eg.Go(func() error {
		stmt := o.db.Statement(ectx, tx, orderItemTable).Where("order_id IN (?)", ids)
		return stmt.Find(&items).Error
	})
	if err := eg.Wait(); err != nil {
		return err
	}

	entity.Orders(orders).Fill(payments.MapByOrderID(), fulfillments.GroupByOrderID(), items.GroupByOrderID())
	return nil
}
