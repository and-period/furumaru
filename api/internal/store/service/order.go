package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListOrders(ctx context.Context, in *store.ListOrdersInput) (entity.Orders, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListOrdersParams{
		CoordinatorID: in.CoordinatorID,
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

func (s *service) CaptureOrder(ctx context.Context, in *store.CaptureOrderInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	if err != nil {
		return internalError(err)
	}
	if !order.Capturable() {
		return fmt.Errorf("service: this order cannot be captured: %w", exception.ErrFailedPrecondition)
	}
	_, err = s.komoju.Payment.Capture(ctx, order.OrderPayment.PaymentID)
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

func (s *service) AggregateOrders(ctx context.Context, in *store.AggregateOrdersInput) (entity.AggregatedOrders, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	orders, err := s.db.Order.Aggregate(ctx, in.UserIDs)
	return orders, internalError(err)
}
