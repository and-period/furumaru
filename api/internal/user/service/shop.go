package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListShops(ctx context.Context, in *user.ListShopsInput) (entity.Shops, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListShopsParams{
		CoordinatorIDs: in.CoordinatorIDs,
		ProducerIDs:    in.ProducerIDs,
		Limit:          int(in.Limit),
		Offset:         int(in.Offset),
	}
	var (
		shops entity.Shops
		total int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		shops, err = s.db.Shop.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Shop.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return shops, total, nil
}

func (s *service) ListShopProducers(ctx context.Context, in *user.ListShopProducersInput) ([]string, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &database.ListShopProducersParams{
		ShopID: in.ShopID,
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	producerIDs, err := s.db.Shop.ListProducers(ctx, params)
	return producerIDs, internalError(err)
}

func (s *service) MultiGetShops(ctx context.Context, in *user.MultiGetShopsInput) (entity.Shops, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	shops, err := s.db.Shop.MultiGet(ctx, in.ShopIDs)
	return shops, internalError(err)
}

func (s *service) GetShop(ctx context.Context, in *user.GetShopInput) (*entity.Shop, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	shop, err := s.db.Shop.Get(ctx, in.ShopID)
	return shop, internalError(err)
}

func (s *service) GetShopByCoordinatorID(ctx context.Context, in *user.GetShopByCoordinatorIDInput) (*entity.Shop, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	shop, err := s.db.Shop.GetByCoordinatorID(ctx, in.CoordinatorID)
	return shop, internalError(err)
}

func (s *service) UpdateShop(ctx context.Context, in *user.UpdateShopInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	productTypesIn := &store.MultiGetProductTypesInput{
		ProductTypeIDs: in.ProductTypeIDs,
	}
	productTypes, err := s.store.MultiGetProductTypes(ctx, productTypesIn)
	if err != nil {
		return internalError(err)
	}
	if len(productTypes) != len(in.ProductTypeIDs) {
		return fmt.Errorf("service: contains invalid product type ids: %w", exception.ErrInvalidArgument)
	}
	params := &database.UpdateShopParams{
		Name:           in.Name,
		ProductTypeIDs: in.ProductTypeIDs,
		BusinessDays:   in.BusinessDays,
	}
	err = s.db.Shop.Update(ctx, in.ShopID, params)
	return internalError(err)
}

func (s *service) RelateShopProducer(ctx context.Context, in *user.RelateShopProducerInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Shop.RelateProducer(ctx, in.ShopID, in.ProducerID)
	return internalError(err)
}

func (s *service) UnrelateShopProducer(ctx context.Context, in *user.UnrelateShopProducerInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Shop.UnrelateProducer(ctx, in.ShopID, in.ProducerID)
	return internalError(err)
}
