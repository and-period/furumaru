package tidb

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
	orderExperienceTable  = "order_experiences"
	orderMetadataTable    = "order_metadata"
)

type order struct {
	db  *mysql.Client
	now func() time.Time
}

func NewOrder(db *mysql.Client) database.Order {
	return &order{
		db:  db,
		now: jst.Now,
	}
}

type listOrdersParams database.ListOrdersParams

func (p listOrdersParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.ShopID != "" {
		stmt = stmt.Where("shop_id = ?", p.ShopID)
	}
	if p.UserID != "" {
		stmt = stmt.Where("user_id = ?", p.UserID)
	}
	if len(p.Types) > 0 {
		stmt = stmt.Where("type IN (?)", p.Types)
	}
	if len(p.Statuses) > 0 {
		stmt = stmt.Where("status IN (?)", p.Statuses)
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
	stmt = stmt.Order("created_at DESC")

	if err := stmt.Find(&orders).Error; err != nil {
		return nil, dbError(err)
	}
	if err := o.fill(ctx, o.db.DB, orders...); err != nil {
		return nil, dbError(err)
	}
	return orders, nil
}

func (o *order) ListUserIDs(ctx context.Context, params *database.ListOrdersParams) ([]string, int64, error) {
	var userIDs []string
	var total int64

	p := listOrdersParams(*params)

	stmt := o.db.Statement(ctx, o.db.DB, orderTable, "DISTINCT(user_id)")
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)
	if err := stmt.Find(&userIDs).Error; err != nil {
		return nil, 0, dbError(err)
	}

	stmt = o.db.Statement(ctx, o.db.DB, orderTable, "COUNT(DISTINCT(user_id))")
	stmt = p.stmt(stmt)
	if err := stmt.Count(&total).Error; err != nil {
		return nil, 0, dbError(err)
	}

	return userIDs, total, nil
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

func (o *order) GetByTransactionIDWithSessionID(ctx context.Context, sessionID, transactionID string) (*entity.Order, error) {
	var order *entity.Order

	stmt := o.db.Statement(ctx, o.db.DB, orderTable, "orders.*").
		Joins("INNER JOIN order_payments ON orders.id = order_payments.order_id").
		Where("orders.session_id = ?", sessionID).
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
		order.OrderMetadata.CreatedAt, order.OrderMetadata.UpdatedAt = now, now

		if err := tx.WithContext(ctx).Table(orderTable).Create(&order).Error; err != nil {
			return err
		}
		if err := tx.WithContext(ctx).Table(orderPaymentTable).Create(&order.OrderPayment).Error; err != nil {
			return err
		}
		if err := tx.WithContext(ctx).Table(orderMetadataTable).Create(&order.OrderMetadata).Error; err != nil {
			return err
		}

		switch order.Type {
		case entity.OrderTypeProduct:
			for _, f := range order.OrderFulfillments {
				f.CreatedAt, f.UpdatedAt = now, now
			}
			for _, i := range order.OrderItems {
				i.CreatedAt, i.UpdatedAt = now, now
			}
			if err := tx.WithContext(ctx).Table(orderFulfillmentTable).Create(&order.OrderFulfillments).Error; err != nil {
				return err
			}
			if err := tx.WithContext(ctx).Table(orderItemTable).Create(&order.OrderItems).Error; err != nil {
				return err
			}
		case entity.OrderTypeExperience:
			order.OrderExperience.CreatedAt, order.OrderExperience.UpdatedAt = now, now
			internal := newInternalOrderExperience(&order.OrderExperience)
			if err := tx.WithContext(ctx).Table(orderExperienceTable).Create(&internal).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return dbError(err)
}

func (o *order) UpdateAuthorized(ctx context.Context, orderID string, params *database.UpdateOrderAuthorizedParams) error {
	p := &updateOrderPaymentParams{
		orderID:  orderID,
		status:   entity.PaymentStatusAuthorized,
		issuedAt: params.IssuedAt,
		validate: func(order *entity.Order) error {
			if order.Completed() {
				return fmt.Errorf("tidb: this order is already completed: %w", database.ErrFailedPrecondition)
			}
			return nil
		},
		updates: map[string]interface{}{
			"status":     entity.PaymentStatusAuthorized,
			"paid_at":    params.IssuedAt,
			"updated_at": o.now(),
		},
	}
	if params.PaymentID != "" {
		p.updates["payment_id"] = params.PaymentID
	}
	return o.updatePayment(ctx, p)
}

func (o *order) UpdateCaptured(ctx context.Context, orderID string, params *database.UpdateOrderCapturedParams) error {
	p := &updateOrderPaymentParams{
		orderID:  orderID,
		status:   entity.PaymentStatusCaptured,
		issuedAt: params.IssuedAt,
		validate: o.noopUpdatePaymentValidate,
		updates: map[string]interface{}{
			"status":      entity.PaymentStatusCaptured,
			"captured_at": params.IssuedAt,
			"updated_at":  o.now(),
		},
	}
	if params.PaymentID != "" {
		p.updates["payment_id"] = params.PaymentID
	}
	return o.updatePayment(ctx, p)
}

func (o *order) UpdateFailed(ctx context.Context, orderID string, params *database.UpdateOrderFailedParams) error {
	p := &updateOrderPaymentParams{
		orderID:  orderID,
		status:   params.Status,
		issuedAt: params.IssuedAt,
		validate: o.noopUpdatePaymentValidate,
		updates: map[string]interface{}{
			"status":     params.Status,
			"failed_at":  params.IssuedAt,
			"updated_at": o.now(),
		},
	}
	if params.PaymentID != "" {
		p.updates["payment_id"] = params.PaymentID
	}
	return o.updatePayment(ctx, p)
}

func (o *order) UpdateRefunded(ctx context.Context, orderID string, params *database.UpdateOrderRefundedParams) error {
	p := &updateOrderPaymentParams{
		orderID:  orderID,
		status:   params.Status,
		issuedAt: params.IssuedAt,
		validate: o.noopUpdatePaymentValidate,
		updates: map[string]interface{}{
			"status":        params.Status,
			"refund_type":   params.RefundType,
			"refund_total":  params.RefundTotal,
			"refund_reason": params.RefundReason,
			"updated_at":    o.now(),
		},
	}
	switch params.Status {
	case entity.PaymentStatusCanceled:
		p.updates["canceled_at"] = params.IssuedAt
	case entity.PaymentStatusRefunded:
		p.updates["refunded_at"] = params.IssuedAt
	}
	return o.updatePayment(ctx, p)
}

func (o *order) UpdateFulfillment(ctx context.Context, orderID, fulfillmentID string, params *database.UpdateOrderFulfillmentParams) error {
	err := o.db.Transaction(ctx, func(tx *gorm.DB) error {
		order, err := o.get(ctx, tx, orderID)
		if err != nil {
			return err
		}
		if order.Completed() {
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

func (o *order) Aggregate(ctx context.Context, params *database.AggregateOrdersParams) (*entity.AggregatedOrder, error) {
	var orders entity.AggregatedOrder

	fields := []string{
		"COUNT(DISTINCT(orders.id)) AS order_count",
		"COUNT(DISTINCT(orders.user_id)) AS user_count",
		"SUM(order_payments.subtotal) AS sales_total",
		"SUM(order_payments.discount) AS discount_total",
	}

	stmt := o.db.Statement(ctx, o.db.DB, orderTable, fields...).
		Joins("INNER JOIN order_payments ON order_payments.order_id = orders.id").
		Where("order_payments.status IN (?)", entity.PaymentSuccessStatuses).
		Where("orders.created_at >= ?", params.CreatedAtGte).
		Where("orders.created_at < ?", params.CreatedAtLt)
	if params.ShopID != "" {
		stmt = stmt.Where("orders.shop_id = ?", params.ShopID)
	}

	err := stmt.Scan(&orders).Error
	return &orders, dbError(err)
}

func (o *order) AggregateByUser(ctx context.Context, params *database.AggregateOrdersByUserParams) (entity.AggregatedUserOrders, error) {
	var orders entity.AggregatedUserOrders

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
	if params.ShopID != "" {
		stmt = stmt.Where("orders.shop_id = ?", params.ShopID)
	}
	stmt = stmt.Group("orders.user_id")

	err := stmt.Scan(&orders).Error
	return orders, dbError(err)
}

func (o *order) AggregateByPaymentMethodType(
	ctx context.Context,
	params *database.AggregateOrdersByPaymentMethodTypeParams,
) (entity.AggregatedOrderPayments, error) {
	var payments entity.AggregatedOrderPayments

	fields := []string{
		"order_payments.method_type AS payment_method_type",
		"COUNT(DISTINCT(orders.id)) AS order_count",
		"COUNT(DISTINCT(orders.user_id)) AS user_count",
		"SUM(order_payments.subtotal) AS sales_total",
	}

	stmt := o.db.Statement(ctx, o.db.DB, orderTable, fields...).
		Joins("INNER JOIN order_payments ON order_payments.order_id = orders.id").
		Where("order_payments.status IN (?)", params.PaymentMethodTypes).
		Where("order_payments.status IN (?)", entity.PaymentSuccessStatuses)
	if params.ShopID != "" {
		stmt = stmt.Where("orders.shop_id = ?", params.ShopID)
	}
	if !params.CreatedAtGte.IsZero() {
		stmt = stmt.Where("orders.created_at >= ?", params.CreatedAtGte)
	}
	if !params.CreatedAtLt.IsZero() {
		stmt = stmt.Where("orders.created_at < ?", params.CreatedAtLt)
	}
	stmt = stmt.Group("order_payments.method_type")

	err := stmt.Scan(&payments).Error
	return payments, dbError(err)
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
	if params.ShopID != "" {
		stmt = stmt.Where("orders.shop_id = ?", params.ShopID)
	}
	stmt = stmt.Group("orders.promotion_id")

	err := stmt.Scan(&orders).Error
	return orders, dbError(err)
}

func (o *order) AggregateByPeriod(
	ctx context.Context,
	params *database.AggregateOrdersByPeriodParams,
) (entity.AggregatedPeriodOrders, error) {
	var internal internalAggregatedPeriodOrders

	var period string
	switch params.PeriodType {
	case entity.AggregateOrderPeriodTypeDay:
		period = "DATE_FORMAT(orders.created_at, '%Y-%m-%d')" // 日付
	case entity.AggregateOrderPeriodTypeWeek:
		period = "DATE_FORMAT(SUBDATE(orders.created_at, WEEKDAY(orders.created_at)+1), '%Y-%m-%d')" // 週のはじめ（日曜日）
	case entity.AggregateOrderPeriodTypeMonth:
		period = "DATE_FORMAT(orders.created_at, '%Y-%m-01')" // 月のはじめ
	default:
		return nil, fmt.Errorf("tidb: invalid period type: %w", database.ErrInvalidArgument)
	}

	fields := []string{
		fmt.Sprintf("%s AS period", period),
		"COUNT(DISTINCT(orders.id)) AS order_count",
		"COUNT(DISTINCT(orders.user_id)) AS user_count",
		"SUM(order_payments.subtotal) AS sales_total",
		"SUM(order_payments.discount) AS discount_total",
	}

	stmt := o.db.Statement(ctx, o.db.DB, orderTable, fields...).
		Joins("INNER JOIN order_payments ON order_payments.order_id = orders.id").
		Where("order_payments.status IN (?)", entity.PaymentSuccessStatuses).
		Where("orders.created_at >= ?", params.CreatedAtGte).
		Where("orders.created_at < ?", params.CreatedAtLt)
	if params.ShopID != "" {
		stmt = stmt.Where("orders.shop_id = ?", params.ShopID)
	}
	stmt = stmt.Group("period").Order("period ASC")

	if err := stmt.Scan(&internal).Error; err != nil {
		return nil, dbError(err)
	}
	return internal.entities(), nil
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
		experiences  entity.OrderExperiences
		metadata     entity.MultiOrderMetadata
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
	eg.Go(func() error {
		var internal internalOrderExperiences
		stmt := o.db.Statement(ectx, tx, orderExperienceTable).Where("order_id IN (?)", ids)
		if err := stmt.Find(&internal).Error; err != nil {
			return err
		}
		experiences = internal.entities()
		return nil
	})
	eg.Go(func() (err error) {
		stmt := o.db.Statement(ectx, tx, orderMetadataTable).Where("order_id IN (?)", ids)
		return stmt.Find(&metadata).Error
	})
	if err := eg.Wait(); err != nil {
		return err
	}

	entity.Orders(orders).Fill(
		payments.MapByOrderID(),
		fulfillments.GroupByOrderID(),
		items.GroupByOrderID(),
		experiences.MapByOrderID(),
		metadata.MapByOrderID(),
	)
	return nil
}

type updateOrderPaymentParams struct {
	orderID  string
	status   entity.PaymentStatus
	issuedAt time.Time
	validate func(order *entity.Order) error
	updates  map[string]interface{}
}

func (o *order) noopUpdatePaymentValidate(order *entity.Order) error {
	return nil
}

func (o *order) updatePayment(ctx context.Context, params *updateOrderPaymentParams) error {
	err := o.db.Transaction(ctx, func(tx *gorm.DB) error {
		order, err := o.get(ctx, tx, params.orderID)
		if err != nil {
			return err
		}
		updatedAt := order.OrderPayment.UpdatedAt.Truncate(time.Second)
		if updatedAt.After(params.issuedAt) {
			return fmt.Errorf("tidb: this refunded event is older than the latest data: %w", database.ErrFailedPrecondition)
		}
		if err := params.validate(order); err != nil {
			return err
		}

		stmt := tx.WithContext(ctx).
			Table(orderPaymentTable).
			Where("order_id = ?", params.orderID)
		if err := stmt.Updates(params.updates).Error; err != nil {
			return err
		}

		order.SetPaymentStatus(params.status)
		return o.updateStatus(ctx, tx, order.ID, order.Status)
	})
	return dbError(err)
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

type internalOrderExperience struct {
	entity.OrderExperience `gorm:"embedded"`
	RemarksJSON            mysql.JSONColumn[entity.OrderExperienceRemarks] `gorm:"default:null;column:remarks"` // 備考(JSON)
}

type internalOrderExperiences []*internalOrderExperience

func newInternalOrderExperience(experience *entity.OrderExperience) *internalOrderExperience {
	return &internalOrderExperience{
		OrderExperience: *experience,
		RemarksJSON:     mysql.NewJSONColumn(experience.Remarks),
	}
}

func (e *internalOrderExperience) entity() *entity.OrderExperience {
	exp := e.OrderExperience
	exp.Remarks = e.RemarksJSON.Val
	return &exp
}



func (es internalOrderExperiences) entities() entity.OrderExperiences {
	res := make(entity.OrderExperiences, len(es))
	for i := range es {
		res[i] = es[i].entity()
	}
	return res
}

type internalAggregatedPeriodOrder struct {
	Period        string
	OrderCount    int64
	UserCount     int64
	SalesTotal    int64
	DiscountTotal int64
}

type internalAggregatedPeriodOrders []*internalAggregatedPeriodOrder

func (o *internalAggregatedPeriodOrder) entity() *entity.AggregatedPeriodOrder {
	period, _ := jst.Parse("2006-01-02", o.Period)
	return &entity.AggregatedPeriodOrder{
		Period:        period,
		OrderCount:    o.OrderCount,
		UserCount:     o.UserCount,
		SalesTotal:    o.SalesTotal,
		DiscountTotal: o.DiscountTotal,
	}
}

func (o internalAggregatedPeriodOrders) entities() entity.AggregatedPeriodOrders {
	res := make(entity.AggregatedPeriodOrders, len(o))
	for i := range o {
		res[i] = o[i].entity()
	}
	return res
}
