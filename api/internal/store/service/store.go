package service

import (
	"context"

	"github.com/and-period/marche/api/internal/store/database"
	"github.com/and-period/marche/api/internal/store/entity"
)

func (s *storeService) ListStores(ctx context.Context, in *ListStoresInput) (entity.Stores, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, storeError(err)
	}
	params := &database.ListStoresParams{
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	stores, err := s.db.Store.List(ctx, params)
	return stores, storeError(err)
}

func (s *storeService) GetStore(ctx context.Context, in *GetStoreInput) (*entity.Store, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, storeError(err)
	}
	store, err := s.db.Store.Get(ctx, in.StoreID)
	return store, storeError(err)
}

func (s *storeService) CreateStore(ctx context.Context, in *CreateStoreInput) (*entity.Store, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, storeError(err)
	}
	return nil, ErrNotImplemented
}

func (s *storeService) UpdateStore(ctx context.Context, in *UpdateStoreInput) error {
	if err := s.validator.Struct(in); err != nil {
		return storeError(err)
	}
	return ErrNotImplemented
}

func (s *storeService) UploadStoreThumbnail(ctx context.Context, in *UploadStoreThumbnailInput) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", storeError(err)
	}
	return "", ErrNotImplemented
}
