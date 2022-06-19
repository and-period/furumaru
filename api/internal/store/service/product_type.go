package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *storeService) ListProductTypes(
	ctx context.Context, in *store.ListProductTypesInput,
) (entity.ProductTypes, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	params := &database.ListProductTypesParams{
		Name:       in.Name,
		CategoryID: in.CategoryID,
		Limit:      int(in.Limit),
		Offset:     int(in.Offset),
	}
	productTypes, err := s.db.ProductType.List(ctx, params)
	return productTypes, exception.InternalError(err)
}

func (s *storeService) CreateProductType(
	ctx context.Context, in *store.CreateProductTypeInput,
) (*entity.ProductType, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	productType := entity.NewProductType(in.Name, in.CategoryID)
	if err := s.db.ProductType.Create(ctx, productType); err != nil {
		return nil, exception.InternalError(err)
	}
	return productType, nil
}

func (s *storeService) UpdateProductType(ctx context.Context, in *store.UpdateProductTypeInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.ProductType.Update(ctx, in.ProductTypeID, in.Name)
	return exception.InternalError(err)
}

func (s *storeService) DeleteProductType(ctx context.Context, in *store.DeleteProductTypeInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.ProductType.Delete(ctx, in.ProductTypeID)
	return exception.InternalError(err)
}
