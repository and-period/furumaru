package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListOrders(ctx context.Context, in *store.ListOrdersInput) (entity.Orders, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, exception.InternalError(err)
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
		return nil, 0, exception.InternalError(err)
	}
	return orders, total, nil
}

func (s *service) GetOrder(ctx context.Context, in *store.GetOrderInput) (*entity.Order, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	return order, exception.InternalError(err)
}

func (s *service) AggregatedOrders(ctx context.Context, in *store.AggregatedOrdersInput) (entity.AggregatedOrders, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	orders, err := s.db.Order.Aggregate(ctx, in.UserIDs)
	return orders, exception.InternalError(err)
}
