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

func (s *service) ListShippings(ctx context.Context, in *store.ListShippingsInput) (entity.Shippings, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	orders := make([]*database.ListShippingsOrder, len(in.Orders))
	for i := range in.Orders {
		orders[i] = &database.ListShippingsOrder{
			Key:        in.Orders[i].Key,
			OrderByASC: in.Orders[i].OrderByASC,
		}
	}
	params := &database.ListShippingsParams{
		CoordinatorID: in.CoordinatorID,
		Name:          in.Name,
		Limit:         int(in.Limit),
		Offset:        int(in.Offset),
		Orders:        orders,
	}
	var (
		shippings entity.Shippings
		total     int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		shippings, err = s.db.Shipping.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Shipping.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	return shippings, total, nil
}

func (s *service) MultiGetShippings(
	ctx context.Context, in *store.MultiGetShippingsInput,
) (entity.Shippings, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	shippings, err := s.db.Shipping.MultiGet(ctx, in.ShippingIDs)
	return shippings, exception.InternalError(err)
}

func (s *service) GetShipping(ctx context.Context, in *store.GetShippingInput) (*entity.Shipping, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	shipping, err := s.db.Shipping.Get(ctx, in.ShippingID)
	return shipping, exception.InternalError(err)
}

func (s *service) CreateShipping(ctx context.Context, in *store.CreateShippingInput) (*entity.Shipping, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	box60Rates, err := s.newShippingRatesFromCreate(in.Box60Rates)
	if err != nil {
		return nil, fmt.Errorf("api: invalid box 60 rates format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	box80Rates, err := s.newShippingRatesFromCreate(in.Box80Rates)
	if err != nil {
		return nil, fmt.Errorf("api: invalid box 80 rates format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	box100Rates, err := s.newShippingRatesFromCreate(in.Box100Rates)
	if err != nil {
		return nil, fmt.Errorf("api: invalid box 100 rates format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	params := &entity.NewShippingParams{
		CoordinatorID:      in.CoordinatorID,
		Name:               in.Name,
		IsDefault:          in.IsDefault,
		Box60Rates:         box60Rates,
		Box60Refrigerated:  in.Box60Refrigerated,
		Box60Frozen:        in.Box60Frozen,
		Box80Rates:         box80Rates,
		Box80Refrigerated:  in.Box80Refrigerated,
		Box80Frozen:        in.Box80Frozen,
		Box100Rates:        box100Rates,
		Box100Refrigerated: in.Box100Refrigerated,
		Box100Frozen:       in.Box100Frozen,
		HasFreeShipping:    in.HasFreeShipping,
		FreeShippingRates:  in.FreeShippingRates,
	}
	shipping := entity.NewShipping(params)
	if err := s.db.Shipping.Create(ctx, shipping); err != nil {
		return nil, exception.InternalError(err)
	}
	return shipping, nil
}

func (s *service) newShippingRatesFromCreate(in []*store.CreateShippingRate) (entity.ShippingRates, error) {
	rates := make(entity.ShippingRates, len(in))
	for i := range in {
		rates[i] = entity.NewShippingRate(int64(i+1), in[i].Name, in[i].Price, in[i].Prefectures)
	}
	if err := rates.Validate(); err != nil {
		return nil, err
	}
	return rates, nil
}

func (s *service) UpdateShipping(ctx context.Context, in *store.UpdateShippingInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	box60Rates, err := s.newShippingRatesFromUpdate(in.Box60Rates)
	if err != nil {
		return fmt.Errorf("api: invalid box 60 rates format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	box80Rates, err := s.newShippingRatesFromUpdate(in.Box80Rates)
	if err != nil {
		return fmt.Errorf("api: invalid box 80 rates format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	box100Rates, err := s.newShippingRatesFromUpdate(in.Box100Rates)
	if err != nil {
		return fmt.Errorf("api: invalid box 100 rates format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	params := &database.UpdateShippingParams{
		Name:               in.Name,
		IsDefault:          in.IsDefault,
		Box60Rates:         box60Rates,
		Box60Refrigerated:  in.Box60Refrigerated,
		Box60Frozen:        in.Box60Frozen,
		Box80Rates:         box80Rates,
		Box80Refrigerated:  in.Box80Refrigerated,
		Box80Frozen:        in.Box80Frozen,
		Box100Rates:        box100Rates,
		Box100Refrigerated: in.Box100Refrigerated,
		Box100Frozen:       in.Box100Frozen,
		HasFreeShipping:    in.HasFreeShipping,
		FreeShippingRates:  in.FreeShippingRates,
	}
	err = s.db.Shipping.Update(ctx, in.ShippingID, params)
	return exception.InternalError(err)
}

func (s *service) newShippingRatesFromUpdate(in []*store.UpdateShippingRate) (entity.ShippingRates, error) {
	rates := make(entity.ShippingRates, len(in))
	for i := range in {
		rates[i] = entity.NewShippingRate(int64(i+1), in[i].Name, in[i].Price, in[i].Prefectures)
	}
	if err := rates.Validate(); err != nil {
		return nil, err
	}
	return rates, nil
}

func (s *service) DeleteShipping(ctx context.Context, in *store.DeleteShippingInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.Shipping.Delete(ctx, in.ShippingID)
	return exception.InternalError(err)
}
