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

func (s *service) ListProductTags(
	ctx context.Context, in *store.ListProductTagsInput,
) (entity.ProductTags, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	orders := make([]*database.ListProductTagsOrder, len(in.Orders))
	for i := range in.Orders {
		orders[i] = &database.ListProductTagsOrder{
			Key:        in.Orders[i].Key,
			OrderByASC: in.Orders[i].OrderByASC,
		}
	}
	params := &database.ListProductTagsParams{
		Name:   in.Name,
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
		Orders: orders,
	}
	var (
		productTags entity.ProductTags
		total       int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		productTags, err = s.db.ProductTag.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.ProductTag.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return productTags, total, nil
}

func (s *service) MultiGetProductTags(
	ctx context.Context, in *store.MultiGetProductTagsInput,
) (entity.ProductTags, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	productTags, err := s.db.ProductTag.MultiGet(ctx, in.ProductTagIDs)
	return productTags, internalError(err)
}

func (s *service) GetProductTag(ctx context.Context, in *store.GetProductTagInput) (*entity.ProductTag, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	productTag, err := s.db.ProductTag.Get(ctx, in.ProductTagID)
	return productTag, internalError(err)
}

func (s *service) CreateProductTag(ctx context.Context, in *store.CreateProductTagInput) (*entity.ProductTag, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	productTag := entity.NewProductTag(in.Name)
	if err := s.db.ProductTag.Create(ctx, productTag); err != nil {
		return nil, internalError(err)
	}
	return productTag, nil
}

func (s *service) UpdateProductTag(ctx context.Context, in *store.UpdateProductTagInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.ProductTag.Update(ctx, in.ProductTagID, in.Name)
	return internalError(err)
}

func (s *service) DeleteProductTag(ctx context.Context, in *store.DeleteProductTagInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &database.ListProductsParams{
		ProductTagID: in.ProductTagID,
	}
	total, err := s.db.Product.Count(ctx, params)
	if err != nil {
		return internalError(err)
	}
	if total > 0 {
		return fmt.Errorf("service: associated with product: %w", exception.ErrFailedPrecondition)
	}
	err = s.db.ProductTag.Delete(ctx, in.ProductTagID)
	return internalError(err)
}
