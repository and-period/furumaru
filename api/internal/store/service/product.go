package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *service) ListProducts(ctx context.Context, in *store.ListProductsInput) (entity.Products, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	params := &database.ListProductsParams{
		Name:       in.Name,
		ProducerID: in.ProducerID,
		CreatedBy:  in.CoordinatorID,
		Limit:      int(in.Limit),
		Offset:     int(in.Offset),
	}
	products, err := s.db.Product.List(ctx, params)
	return products, exception.InternalError(err)
}

func (s *service) GetProduct(ctx context.Context, in *store.GetProductInput) (*entity.Product, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	product, err := s.db.Product.Get(ctx, in.ProductID)
	return product, exception.InternalError(err)
}

func (s *service) CreateProduct(ctx context.Context, in *store.CreateProductInput) (*entity.Product, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	multiMedia := make(entity.MultiProductMedia, len(in.Media))
	for i := range in.Media {
		multiMedia[i] = entity.NewProductMedia(in.Media[i].URL, in.Media[i].IsThumbnail)
	}
	params := &entity.NewProductParams{}
	return nil, exception.ErrNotImplemented
}

func (s *service) UpdateProduct(ctx context.Context, in *store.UpdateProductInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	return exception.ErrNotImplemented
}

func (s *service) DeleteProduct(ctx context.Context, in *store.DeleteProductInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	return exception.ErrNotImplemented
}
