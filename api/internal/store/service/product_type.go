package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListProductTypes(
	ctx context.Context, in *store.ListProductTypesInput,
) (entity.ProductTypes, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	params := &database.ListProductTypesParams{
		Name:       in.Name,
		CategoryID: in.CategoryID,
		Limit:      int(in.Limit),
		Offset:     int(in.Offset),
	}
	var (
		productTypes entity.ProductTypes
		total        int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		productTypes, err = s.db.ProductType.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.ProductType.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	return productTypes, total, nil
}

func (s *service) MultiGetProductTypes(
	ctx context.Context, in *store.MultiGetProductTypesInput,
) (entity.ProductTypes, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	productTypes, err := s.db.ProductType.MultiGet(ctx, in.ProductTypeIDs)
	return productTypes, exception.InternalError(err)
}

func (s *service) GetProductType(ctx context.Context, in *store.GetProductTypeInput) (*entity.ProductType, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	productType, err := s.db.ProductType.Get(ctx, in.ProductTypeID)
	return productType, exception.InternalError(err)
}

func (s *service) CreateProductType(
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

func (s *service) UpdateProductType(ctx context.Context, in *store.UpdateProductTypeInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.ProductType.Update(ctx, in.ProductTypeID, in.Name)
	return exception.InternalError(err)
}

func (s *service) DeleteProductType(ctx context.Context, in *store.DeleteProductTypeInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.ProductType.Delete(ctx, in.ProductTypeID)
	return exception.InternalError(err)
}
