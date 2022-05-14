package service

import (
	"context"

	"github.com/and-period/marche/api/internal/exception"
	"github.com/and-period/marche/api/internal/store"
	"github.com/and-period/marche/api/internal/store/database"
	"github.com/and-period/marche/api/internal/store/entity"
)

func (s *storeService) ListStores(ctx context.Context, in *store.ListStoresInput) (entity.Stores, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	params := &database.ListStoresParams{
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	stores, err := s.db.Store.List(ctx, params)
	return stores, exception.InternalError(err)
}

func (s *storeService) GetStore(ctx context.Context, in *store.GetStoreInput) (*entity.Store, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	store, err := s.db.Store.Get(ctx, in.StoreID)
	return store, exception.InternalError(err)
}

func (s *storeService) CreateStore(ctx context.Context, in *store.CreateStoreInput) (*entity.Store, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	store := entity.NewStore(in.Name)
	if err := s.db.Store.Create(ctx, store); err != nil {
		return nil, exception.InternalError(err)
	}
	return store, nil
}

func (s *storeService) UpdateStore(ctx context.Context, in *store.UpdateStoreInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.Store.Update(ctx, in.StoreID, in.Name, in.ThumbnailURL)
	return exception.InternalError(err)
}
