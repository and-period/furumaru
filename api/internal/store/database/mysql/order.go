package mysql

import (
	"context"
	"fmt"
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
	if p.UserID != "" {
		stmt = stmt.Where("user_id = ?", p.UserID)
	}
	stmt = stmt.Order("created_at DESC")
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
	return order, dbError(err)
}

func (o *order) GetByTransactionID(ctx context.Context, userID, transactionID string) (*entity.Order, error) {
	var order *entity.Order

	stmt := o.db.Statement(ctx, o.db.DB, orderTable, "orders.*").
		Joins("INNER JOIN order_payments ON orders.id = order_payments.order_id").
		Where("orders.user_id = ?", userID).
		Where("order_payments.transaction_id = ?", transactionID)

	if err := stmt.First(&order).Error; err != nil {
		return nil, dbError(err)
	}
	if err := o.fill(ctx, o.db.DB, order); err != nil {
		return nil, dbError(err)
	}
	return order, nil
}

func (o *order) Create(ctx context.Context, order *entity.Order) error {
	err := o.db.Transaction(ctx, func(tx *gorm.DB) error {
		count := func(stmt *gorm.DB) *gorm.DB {
			return stmt.Where("coordinator_id = ?", order.CoordinatorID)
		}
		total, err := o.db.Count(ctx, o.db.DB, &entity.Order{}, count)
		if err != nil {
			return err
		}

		now := o.now()
		order.ManagementID = total + 1
		order.CreatedAt, order.UpdatedAt = now, now
		order.OrderPayment.CreatedAt, order.OrderPayment.UpdatedAt = now, now
		for _, f := range order.OrderFulfillments {
			f.CreatedAt, f.UpdatedAt = now, now
		}
		for _, i := range order.OrderItems {
			i.CreatedAt, i.UpdatedAt = now, now
		}

		if err := tx.WithContext(ctx).Table(orderTable).Create(&order).Error; err != nil {
			return err
		}
		if err := tx.WithContext(ctx).Table(orderPaymentTable).Create(&order.OrderPayment).Error; err != nil {
			return err
		}
		if err := tx.WithContext(ctx).Table(orderFulfillmentTable).Create(&order.OrderFulfillments).Error; err != nil {
			return err
		}
		return tx.WithContext(ctx).Table(orderItemTable).Create(&order.OrderItems).Error
	})
	return dbError(err)
}

func (o *order) UpdatePayment(ctx context.Context, orderID string, params *database.UpdateOrderPaymentParams) error {
	err := o.db.Transaction(ctx, func(tx *gorm.DB) error {
		order, err := o.get(ctx, tx, orderID)
		if err != nil {
			return err
		}
		if order.IsCompleted() {
			return fmt.Errorf("mysql: this order is already completed: %w", database.ErrFailedPrecondition)
		}
		updatedAt := order.OrderPayment.UpdatedAt.Truncate(time.Second)
		if updatedAt.After(params.IssuedAt) {
			return fmt.Errorf("mysql: this event is older than the latest data: %w", database.ErrFailedPrecondition)
		}

		updates := map[string]interface{}{
			"status":     params.Status,
			"updated_at": o.now(),
		}
		switch params.Status {
		case entity.PaymentStatusAuthorized:
			updates["paid_at"] = params.IssuedAt
		case entity.PaymentStatusCaptured:
			updates["captured_at"] = params.IssuedAt
		case entity.PaymentStatusFailed:
			updates["failed_at"] = params.IssuedAt
		}
		if params.PaymentID != "" {
			updates["payment_id"] = params.PaymentID
		}

		stmt := tx.WithContext(ctx).
			Table(orderPaymentTable).
			Where("order_id = ?", orderID)
		if err := stmt.Updates(updates).Error; err != nil {
			return err
		}

		order.SetPaymentStatus(params.Status)
		return o.updateStatus(ctx, tx, order.ID, order.Status)
	})
	return dbError(err)
}

func (o *order) UpdateFulfillment(ctx context.Context, orderID, fulfillmentID string, params *database.UpdateOrderFulfillmentParams) error {
	err := o.db.Transaction(ctx, func(tx *gorm.DB) error {
		order, err := o.get(ctx, tx, orderID)
		if err != nil {
			return err
		}
		if order.IsCompleted() {
			return fmt.Errorf("mysql: this order is already completed: %w", database.ErrFailedPrecondition)
		}

		updates := map[string]interface{}{
			"status":     params.Status,
			"updated_at": o.now(),
		}
		switch params.Status {
		case entity.FulfillmentStatusFulfilled:
			updates["shipping_carrier"] = params.ShippingCarrier
			updates["tracking_number"] = params.TrackingNumber
			updates["shipped_at"] = params.ShippedAt
		case entity.FulfillmentStatusUnfulfilled:
			updates["shipping_carrier"] = entity.ShippingCarrierUnknown
			updates["tracking_number"] = nil
			updates["shipped_at"] = nil
		}
		stmt := tx.WithContext(ctx).
			Table(orderFulfillmentTable).
			Where("order_id = ?", orderID).
			Where("id = ?", fulfillmentID)
		if err := stmt.Updates(updates).Error; err != nil {
			return err
		}

		order.SetFulfillmentStatus(fulfillmentID, params.Status)
		return o.updateStatus(ctx, tx, order.ID, order.Status)
	})
	return dbError(err)
}

func (o *order) UpdateRefund(ctx context.Context, orderID string, params *database.UpdateOrderRefundParams) error {
	err := o.db.Transaction(ctx, func(tx *gorm.DB) error {
		order, err := o.get(ctx, tx, orderID)
		if err != nil {
			return err
		}
		updatedAt := order.OrderPayment.UpdatedAt.Truncate(time.Second)
		if updatedAt.After(params.IssuedAt) {
			return fmt.Errorf("mysql: this event is older than the latest data: %w", database.ErrFailedPrecondition)
		}

		updates := map[string]interface{}{
			"status":        params.Status,
			"refund_type":   params.RefundType,
			"refund_total":  params.RefundTotal,
			"refund_reason": params.RefundReason,
			"updated_at":    o.now(),
		}
		switch params.Status {
		case entity.PaymentStatusCanceled:
			updates["canceled_at"] = params.IssuedAt
		case entity.PaymentStatusRefunded:
			updates["refunded_at"] = params.IssuedAt
		}

		stmt := tx.WithContext(ctx).
			Table(orderPaymentTable).
			Where("order_id = ?", orderID)
		if err := stmt.Updates(updates).Error; err != nil {
			return err
		}

		order.SetPaymentStatus(params.Status)
		return o.updateStatus(ctx, tx, order.ID, order.Status)
	})
	return dbError(err)
}

func (o *order) Draft(ctx context.Context, orderID string, params *database.DraftOrderParams) error {
	updates := map[string]interface{}{
		"shipping_message": params.ShippingMessage,
		"updated_at":       o.now(),
	}
	err := o.db.DB.WithContext(ctx).
		Table(orderTable).
		Where("id = ?", orderID).
		Updates(updates).Error
	return dbError(err)
}

func (o *order) Complete(ctx context.Context, orderID string, params *database.CompleteOrderParams) error {
	now := o.now()
	updates := map[string]interface{}{
		"shipping_message": params.ShippingMessage,
		"status":           entity.OrderStatusCompleted,
		"completed_at":     now,
		"updated_at":       now,
	}
	err := o.db.DB.WithContext(ctx).
		Table(orderTable).
		Where("id = ?", orderID).
		Updates(updates).Error
	return dbError(err)
}

func (o *order) Aggregate(ctx context.Context, params *database.AggregateOrdersParams) (entity.AggregatedOrders, error) {
	var orders entity.AggregatedOrders

	fields := []string{
		"orders.user_id AS user_id",
		"COUNT(DISTINCT(orders.id)) AS order_count",
		"SUM(order_payments.subtotal) AS subtotal",
		"SUM(order_payments.discount) AS discount",
		"SUM(order_payments.total) AS total",
	}

	stmt := o.db.Statement(ctx, o.db.DB, orderTable, fields...).
		Joins("INNER JOIN order_payments ON order_payments.order_id = orders.id").
		Where("orders.user_id IN (?)", params.UserIDs).
		Where("order_payments.status IN (?)", entity.PaymentSuccessStatuses)
	if params.CoordinatorID != "" {
		stmt = stmt.Where("orders.coordinator_id = ?", params.CoordinatorID)
	}
	stmt = stmt.Group("orders.user_id")

	err := stmt.Scan(&orders).Error
	return orders, dbError(err)
}

func (o *order) AggregateByPromotion(
	ctx context.Context,
	params *database.AggregateOrdersByPromotionParams,
) (entity.AggregatedOrderPromotions, error) {
	var orders entity.AggregatedOrderPromotions

	fields := []string{
		"orders.promotion_id AS promotion_id",
		"COUNT(DISTINCT(orders.id)) AS order_count",
		"SUM(order_payments.discount) AS discount_total",
	}

	stmt := o.db.Statement(ctx, o.db.DB, orderTable, fields...).
		Joins("INNER JOIN order_payments ON order_payments.order_id = orders.id").
		Where("orders.promotion_id IN (?)", params.PromotionIDs).
		Where("order_payments.status IN (?)", entity.PaymentSuccessStatuses)
	if params.CoordinatorID != "" {
		stmt = stmt.Where("orders.coordinator_id = ?", params.CoordinatorID)
	}
	stmt = stmt.Group("orders.promotion_id")

	err := stmt.Scan(&orders).Error
	return orders, dbError(err)
}

func (o *order) get(ctx context.Context, tx *gorm.DB, orderID string, fields ...string) (*entity.Order, error) {
	var order *entity.Order

	stmt := o.db.Statement(ctx, tx, orderTable, fields...).
		Where("id = ?", orderID)

	if err := stmt.First(&order).Error; err != nil {
		return nil, err
	}
	if err := o.fill(ctx, o.db.DB, order); err != nil {
		return nil, err
	}
	return order, nil
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

func (o *order) updateStatus(ctx context.Context, tx *gorm.DB, orderID string, status entity.OrderStatus) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": o.now(),
	}
	stmt := tx.WithContext(ctx).
		Table(orderTable).
		Where("id = ?", orderID)
	return stmt.Updates(updates).Error
}
