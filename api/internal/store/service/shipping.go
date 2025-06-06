package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListShippingsByShopID(
	ctx context.Context, in *store.ListShippingsByShopIDInput,
) (entity.Shippings, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListShippingsParams{
		ShopID: in.ShopID,
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	var (
		shippings entity.Shippings
		total     int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var err error
		shippings, err = s.db.Shipping.List(ectx, params)
		return err
	})
	eg.Go(func() error {
		var err error
		total, err = s.db.Shipping.Count(ectx, params)
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return shippings, total, nil
}

func (s *service) ListShippingsByShopIDs(ctx context.Context, in *store.ListShippingsByShopIDsInput) (entity.Shippings, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	shippings, err := s.db.Shipping.ListByShopIDs(ctx, in.ShopIDs)
	return shippings, internalError(err)
}

func (s *service) ListShippingsByCoordinatorIDs(
	ctx context.Context, in *store.ListShippingsByCoordinatorIDsInput,
) (entity.Shippings, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	shippings, err := s.db.Shipping.ListByCoordinatorIDs(ctx, in.CoordinatorIDs)
	return shippings, internalError(err)
}

func (s *service) MultiGetShippingsByRevision(
	ctx context.Context, in *store.MultiGetShippingsByRevisionInput,
) (entity.Shippings, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	shippings, err := s.db.Shipping.MultiGetByRevision(ctx, in.ShippingRevisionIDs)
	return shippings, internalError(err)
}

func (s *service) GetShipping(ctx context.Context, in *store.GetShippingInput) (*entity.Shipping, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	shipping, err := s.db.Shipping.Get(ctx, in.ShippingID)
	return shipping, internalError(err)
}

func (s *service) GetDefaultShipping(ctx context.Context, in *store.GetDefaultShippingInput) (*entity.Shipping, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	shipping, err := s.db.Shipping.GetDefault(ctx)
	return shipping, internalError(err)
}

func (s *service) GetShippingByShopID(ctx context.Context, in *store.GetShippingByShopIDInput) (*entity.Shipping, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	shipping, err := s.db.Shipping.GetByShopID(ctx, in.ShopID)
	return shipping, internalError(err)
}

func (s *service) GetShippingByCoordinatorID(
	ctx context.Context, in *store.GetShippingByCoordinatorIDInput,
) (*entity.Shipping, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	shipping, err := s.db.Shipping.GetByCoordinatorID(ctx, in.CoordinatorID)
	return shipping, internalError(err)
}

func (s *service) CreateShipping(ctx context.Context, in *store.CreateShippingInput) (*entity.Shipping, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
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

	var inUse bool
	// 現在の配送設定を取得
	_, err = s.db.Shipping.GetByCoordinatorID(ctx, in.CoordinatorID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return nil, internalError(err)
	}
	// 配送設定が存在しない場合、今回作成するものを使用中のものとして作成する
	if errors.Is(err, database.ErrNotFound) {
		inUse = true
	}

	params := &entity.NewShippingParams{
		ShopID:            in.ShopID,
		CoordinatorID:     in.CoordinatorID,
		Name:              in.Name,
		Box60Rates:        box60Rates,
		Box60Frozen:       in.Box60Frozen,
		Box80Rates:        box80Rates,
		Box80Frozen:       in.Box80Frozen,
		Box100Rates:       box100Rates,
		Box100Frozen:      in.Box100Frozen,
		HasFreeShipping:   in.HasFreeShipping,
		FreeShippingRates: in.FreeShippingRates,
		InUse:             inUse,
	}
	shipping := entity.NewShipping(params)
	if err := s.db.Shipping.Create(ctx, shipping); err != nil {
		return nil, internalError(err)
	}
	return shipping, nil
}

func (s *service) UpdateShipping(ctx context.Context, in *store.UpdateShippingInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
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
		Name:              in.Name,
		Box60Rates:        box60Rates,
		Box60Frozen:       in.Box60Frozen,
		Box80Rates:        box80Rates,
		Box80Frozen:       in.Box80Frozen,
		Box100Rates:       box100Rates,
		Box100Frozen:      in.Box100Frozen,
		HasFreeShipping:   in.HasFreeShipping,
		FreeShippingRates: in.FreeShippingRates,
	}
	err = s.db.Shipping.Update(ctx, in.ShippingID, params)
	return internalError(err)
}

func (s *service) UpdateShippingInUse(ctx context.Context, in *store.UpdateShippingInUseInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Shipping.UpdateInUse(ctx, in.ShopID, in.ShippingID)
	return internalError(err)
}

func (s *service) UpdateDefaultShipping(ctx context.Context, in *store.UpdateDefaultShippingInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	box60Rates, err := s.newShippingRatesFromUpdateDefault(in.Box60Rates)
	if err != nil {
		return fmt.Errorf("api: invalid box 60 rates format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	box80Rates, err := s.newShippingRatesFromUpdateDefault(in.Box80Rates)
	if err != nil {
		return fmt.Errorf("api: invalid box 80 rates format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	box100Rates, err := s.newShippingRatesFromUpdateDefault(in.Box100Rates)
	if err != nil {
		return fmt.Errorf("api: invalid box 100 rates format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	params := &database.UpdateShippingParams{
		Box60Rates:        box60Rates,
		Box60Frozen:       in.Box60Frozen,
		Box80Rates:        box80Rates,
		Box80Frozen:       in.Box80Frozen,
		Box100Rates:       box100Rates,
		Box100Frozen:      in.Box100Frozen,
		HasFreeShipping:   in.HasFreeShipping,
		FreeShippingRates: in.FreeShippingRates,
	}
	err = s.db.Shipping.Update(ctx, entity.DefaultShippingID, params)
	return internalError(err)
}

// Deprecated
func (s *service) UpsertShipping(ctx context.Context, in *store.UpsertShippingInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	box60Rates, err := s.newShippingRatesFromUpsert(in.Box60Rates)
	if err != nil {
		return fmt.Errorf("api: invalid box 60 rates format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	box80Rates, err := s.newShippingRatesFromUpsert(in.Box80Rates)
	if err != nil {
		return fmt.Errorf("api: invalid box 80 rates format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	box100Rates, err := s.newShippingRatesFromUpsert(in.Box100Rates)
	if err != nil {
		return fmt.Errorf("api: invalid box 100 rates format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	shipping, err := s.db.Shipping.GetByCoordinatorID(ctx, in.CoordinatorID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return internalError(err)
	}
	if errors.Is(err, database.ErrNotFound) {
		// create
		params := &entity.NewShippingParams{
			ShopID:            in.ShopID,
			CoordinatorID:     in.CoordinatorID,
			Box60Rates:        box60Rates,
			Box60Frozen:       in.Box60Frozen,
			Box80Rates:        box80Rates,
			Box80Frozen:       in.Box80Frozen,
			Box100Rates:       box100Rates,
			Box100Frozen:      in.Box100Frozen,
			HasFreeShipping:   in.HasFreeShipping,
			FreeShippingRates: in.FreeShippingRates,
			InUse:             true,
		}
		shipping = entity.NewShipping(params)
		err = s.db.Shipping.Create(ctx, shipping)
	} else {
		// update
		params := &database.UpdateShippingParams{
			Box60Rates:        box60Rates,
			Box60Frozen:       in.Box60Frozen,
			Box80Rates:        box80Rates,
			Box80Frozen:       in.Box80Frozen,
			Box100Rates:       box100Rates,
			Box100Frozen:      in.Box100Frozen,
			HasFreeShipping:   in.HasFreeShipping,
			FreeShippingRates: in.FreeShippingRates,
		}
		err = s.db.Shipping.Update(ctx, shipping.ID, params)
	}
	return internalError(err)
}

func (s *service) DeleteShipping(ctx context.Context, in *store.DeleteShippingInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Shipping.Delete(ctx, in.ShippingID)
	return internalError(err)
}

func (s *service) newShippingRatesFromUpdateDefault(in []*store.UpdateDefaultShippingRate) (entity.ShippingRates, error) {
	rates := make(entity.ShippingRates, len(in))
	for i := range in {
		rates[i] = entity.NewShippingRate(int64(i+1), in[i].Name, in[i].Price, in[i].PrefectureCodes)
	}
	if err := rates.Validate(); err != nil {
		return nil, err
	}
	return rates, nil
}

func (s *service) newShippingRatesFromCreate(in []*store.CreateShippingRate) (entity.ShippingRates, error) {
	rates := make(entity.ShippingRates, len(in))
	for i := range in {
		rates[i] = entity.NewShippingRate(int64(i+1), in[i].Name, in[i].Price, in[i].PrefectureCodes)
	}
	if err := rates.Validate(); err != nil {
		return nil, err
	}
	return rates, nil
}

func (s *service) newShippingRatesFromUpdate(in []*store.UpdateShippingRate) (entity.ShippingRates, error) {
	rates := make(entity.ShippingRates, len(in))
	for i := range in {
		rates[i] = entity.NewShippingRate(int64(i+1), in[i].Name, in[i].Price, in[i].PrefectureCodes)
	}
	if err := rates.Validate(); err != nil {
		return nil, err
	}
	return rates, nil
}

func (s *service) newShippingRatesFromUpsert(in []*store.UpsertShippingRate) (entity.ShippingRates, error) {
	rates := make(entity.ShippingRates, len(in))
	for i := range in {
		rates[i] = entity.NewShippingRate(int64(i+1), in[i].Name, in[i].Price, in[i].PrefectureCodes)
	}
	if err := rates.Validate(); err != nil {
		return nil, err
	}
	return rates, nil
}
