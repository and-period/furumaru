package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

const orderTable = "orders"

type order struct {
	db  *database.Client
	now func() time.Time
}

func NewOrder(db *database.Client) Order {
	return &order{
		db:  db,
		now: jst.Now,
	}
}

func (o *order) List(ctx context.Context, params *ListOrdersParams, fields ...string) (entity.Orders, error) {
	var orders entity.Orders

	stmt := o.db.Statement(ctx, o.db.DB, orderTable, fields...)
	stmt = params.stmt(stmt)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	if err := stmt.Find(&orders).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	if err := o.fill(ctx, o.db.DB, orders...); err != nil {
		return nil, exception.InternalError(err)
	}
	return orders, nil
}

func (o *order) Count(ctx context.Context, params *ListOrdersParams) (int64, error) {
	total, err := o.db.Count(ctx, o.db.DB, &entity.Order{}, nil)
	return total, exception.InternalError(err)
}

func (o *order) Get(ctx context.Context, orderID string, fields ...string) (*entity.Order, error) {
	order, err := o.get(ctx, o.db.DB, orderID, fields...)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	if err := o.fill(ctx, o.db.DB, order); err != nil {
		return nil, exception.InternalError(err)
	}
	return order, nil
}

func (o *order) Aggregate(ctx context.Context, userIDs []string) (entity.AggregatedOrders, error) {
	var orders entity.AggregatedOrders

	fields := []string{
		"orders.user_id AS user_id",
		"COUNT(DISTINCT(orders.id)) AS order_count",
		"SUM(payments.subtotal) AS subtotal",
		"SUM(payments.discount) AS discount",
	}

	stmt := o.db.Statement(ctx, o.db.DB, orderTable, fields...).
		Joins("INNER JOIN payments ON payments.order_id = orders.id").
		Where("orders.user_id IN (?)", userIDs).
		Group("orders.user_id")

	err := stmt.Scan(&orders).Error
	return orders, exception.InternalError(err)
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
		payments     entity.Payments
		fulfillments entity.Fulfillments
		activities   entity.Activities
		items        entity.OrderItems
	)

	ids := entity.Orders(orders).IDs()
	if len(ids) == 0 {
		return nil
	}

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		stmt := o.db.Statement(ectx, tx, paymentTable).Where("order_id IN (?)", ids)
		return stmt.Find(&payments).Error
	})
	eg.Go(func() error {
		stmt := o.db.Statement(ectx, tx, fulfillmentTable).Where("order_id IN (?)", ids)
		return stmt.Find(&fulfillments).Error
	})
	eg.Go(func() error {
		stmt := o.db.Statement(ectx, tx, activityTable).Where("order_id IN (?)", ids)
		return stmt.Find(&activities).Error
	})
	eg.Go(func() error {
		stmt := o.db.Statement(ectx, tx, orderItemTable).Where("order_id IN (?)", ids)
		return stmt.Find(&items).Error
	})
	if err := eg.Wait(); err != nil {
		return err
	}

	paymentMap := payments.MapByOrderID()
	fulfillmentMap := fulfillments.MapByOrderID()
	activitiesMap := activities.GroupByOrderID()
	itemsMap := items.GroupByOrderID()

	for i, o := range orders {
		payment, ok := paymentMap[o.ID]
		if !ok {
			payment = &entity.Payment{}
		}
		fulfillment, ok := fulfillmentMap[o.ID]
		if !ok {
			fulfillment = &entity.Fulfillment{}
		}
		orders[i].Fill(payment, fulfillment, activitiesMap[o.ID], itemsMap[o.ID])
	}
	return nil
}
