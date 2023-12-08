package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/komoju"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListOrders(ctx context.Context, in *store.ListOrdersInput) (entity.Orders, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListOrdersParams{
		CoordinatorID: in.CoordinatorID,
		UserID:        in.UserID,
		Limit:         int(in.Limit),
		Offset:        int(in.Offset),
	}
	var (
		orders entity.Orders
		total  int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		orders, err = s.db.Order.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Order.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return orders, total, nil
}

func (s *service) GetOrder(ctx context.Context, in *store.GetOrderInput) (*entity.Order, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	return order, internalError(err)
}

func (s *service) GetOrderByTransactionID(ctx context.Context, in *store.GetOrderByTransactionIDInput) (*entity.Order, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	order, err := s.db.Order.GetByTransactionID(ctx, in.UserID, in.TransactionID)
	return order, internalError(err)
}

func (s *service) CaptureOrder(ctx context.Context, in *store.CaptureOrderInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	if err != nil {
		return internalError(err)
	}
	if !order.Capturable() {
		return fmt.Errorf("service: this order cannot be capture: %w", exception.ErrFailedPrecondition)
	}
	_, err = s.komoju.Payment.Capture(ctx, order.OrderPayment.PaymentID)
	return internalError(err)
}

func (s *service) DraftOrder(ctx context.Context, in *store.DraftOrderInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	if err != nil {
		return internalError(err)
	}
	if !order.Preservable() {
		return fmt.Errorf("service: this order cannot be save: %w", exception.ErrFailedPrecondition)
	}
	params := &database.DraftOrderParams{
		ShippingMessage: in.ShippingMessage,
	}
	err = s.db.Order.Draft(ctx, in.OrderID, params)
	return internalError(err)
}

func (s *service) CompleteOrder(ctx context.Context, in *store.CompleteOrderInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	if err != nil {
		return internalError(err)
	}
	if !order.Completable() {
		return fmt.Errorf("service: this order cannot be complete: %w", exception.ErrFailedPrecondition)
	}
	params := &database.CompleteOrderParams{
		ShippingMessage: in.ShippingMessage,
		CompletedAt:     s.now(),
	}
	err = s.db.Order.Complete(ctx, in.OrderID, params)
	return internalError(err)
}

func (s *service) CancelOrder(ctx context.Context, in *store.CancelOrderInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	if err != nil {
		return internalError(err)
	}
	if !order.Cancelable() {
		return fmt.Errorf("service: this order cannot be canceled: %w", exception.ErrFailedPrecondition)
	}
	_, err = s.komoju.Payment.Cancel(ctx, order.OrderPayment.PaymentID)
	return internalError(err)
}

func (s *service) RefundOrder(ctx context.Context, in *store.RefundOrderInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	if err != nil {
		return internalError(err)
	}
	if !order.Refundable() {
		return fmt.Errorf("service: this order cannot be refund: %w", exception.ErrFailedPrecondition)
	}
	params := &komoju.RefundParams{
		PaymentID:   order.OrderPayment.PaymentID,
		Amount:      order.OrderPayment.Total,
		Description: in.Description,
	}
	_, err = s.komoju.Payment.Refund(ctx, params)
	return internalError(err)
}

func (s *service) UpdateOrderFulfillment(ctx context.Context, in *store.UpdateOrderFulfillmentInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	if err != nil {
		return internalError(err)
	}
	if order.Completed() {
		return fmt.Errorf("service: this order is already completed: %w", exception.ErrFailedPrecondition)
	}
	params := &database.UpdateOrderFulfillmentParams{
		Status:          entity.FulfillmentStatusFulfilled,
		ShippingCarrier: in.ShippingCarrier,
		TrackingNumber:  in.TrackingNumber,
		ShippedAt:       s.now(),
	}
	err = s.db.Order.UpdateFulfillment(ctx, in.FulfillmentID, params)
	return internalError(err)
}

func (s *service) AggregateOrders(ctx context.Context, in *store.AggregateOrdersInput) (entity.AggregatedOrders, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &database.AggregateOrdersParams{
		CoordinatorID: in.CoordinatorID,
		UserIDs:       in.UserIDs,
	}
	orders, err := s.db.Order.Aggregate(ctx, params)
	return orders, internalError(err)
}
