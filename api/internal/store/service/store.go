package service

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/and-period/marche/api/internal/store/database"
	"github.com/and-period/marche/api/internal/store/entity"
	"github.com/and-period/marche/api/pkg/uuid"
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
	store := entity.NewStore(in.Name)
	if err := s.db.Store.Create(ctx, store); err != nil {
		return nil, storeError(err)
	}
	return store, nil
}

func (s *storeService) UpdateStore(ctx context.Context, in *UpdateStoreInput) error {
	if err := s.validator.Struct(in); err != nil {
		return storeError(err)
	}
	err := s.db.Store.Update(ctx, in.StoreID, in.Name, in.ThumbnailURL)
	return storeError(err)
}

func (s *storeService) UploadStoreThumbnail(ctx context.Context, in *UploadStoreThumbnailInput) (string, error) {
	const format = "stores/%s/thumbnails/%s"
	if err := s.validator.Struct(in); err != nil {
		return "", storeError(err)
	}
	var b bytes.Buffer
	if _, err := b.Write(in.Image); err != nil {
		return "", storeError(err)
	}
	path := fmt.Sprintf(format, strconv.FormatInt(in.StoreID, 10), uuid.Base58Encode(uuid.New()))
	thumbnailURL, err := s.storage.Upload(ctx, path, &b)
	if err != nil {
		return "", storeError(err)
	}
	return thumbnailURL, nil
}
