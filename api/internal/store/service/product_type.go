package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListProductTypes(
	ctx context.Context, in *store.ListProductTypesInput,
) (entity.ProductTypes, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	orders := make([]*database.ListProductTypesOrder, len(in.Orders))
	for i := range in.Orders {
		orders[i] = &database.ListProductTypesOrder{
			Key:        in.Orders[i].Key,
			OrderByASC: in.Orders[i].OrderByASC,
		}
	}
	params := &database.ListProductTypesParams{
		Name:       in.Name,
		CategoryID: in.CategoryID,
		Limit:      int(in.Limit),
		Offset:     int(in.Offset),
		Orders:     orders,
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
		return nil, 0, internalError(err)
	}
	return productTypes, total, nil
}

func (s *service) MultiGetProductTypes(
	ctx context.Context, in *store.MultiGetProductTypesInput,
) (entity.ProductTypes, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	productTypes, err := s.db.ProductType.MultiGet(ctx, in.ProductTypeIDs)
	return productTypes, internalError(err)
}

func (s *service) GetProductType(ctx context.Context, in *store.GetProductTypeInput) (*entity.ProductType, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	productType, err := s.db.ProductType.Get(ctx, in.ProductTypeID)
	return productType, internalError(err)
}

func (s *service) CreateProductType(
	ctx context.Context, in *store.CreateProductTypeInput,
) (*entity.ProductType, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewProductTypeParams{
		CategoryID: in.CategoryID,
		Name:       in.Name,
		IconURL:    in.IconURL,
	}
	productType := entity.NewProductType(params)
	if err := s.db.ProductType.Create(ctx, productType); err != nil {
		return nil, internalError(err)
	}
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		s.resizeProductType(context.Background(), productType.ID, productType.IconURL)
	}()
	return productType, nil
}

func (s *service) UpdateProductType(ctx context.Context, in *store.UpdateProductTypeInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	productType, err := s.db.ProductType.Get(ctx, in.ProductTypeID)
	if err != nil {
		return internalError(err)
	}
	if err := s.db.ProductType.Update(ctx, in.ProductTypeID, in.Name, in.IconURL); err != nil {
		return internalError(err)
	}
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		var iconURL string
		if productType.IconURL != in.IconURL {
			iconURL = in.IconURL
		}
		s.resizeProductType(context.Background(), productType.ID, iconURL)
	}()
	return nil
}

func (s *service) UpdateProductTypeIcons(ctx context.Context, in *store.UpdateProductTypeIconsInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.ProductType.UpdateIcons(ctx, in.ProductTypeID, in.Icons)
	return internalError(err)
}

func (s *service) DeleteProductType(ctx context.Context, in *store.DeleteProductTypeInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &database.ListProductsParams{
		ProductTypeIDs: []string{in.ProductTypeID},
	}
	total, err := s.db.Product.Count(ctx, params)
	if err != nil {
		return internalError(err)
	}
	if total > 0 {
		return fmt.Errorf("service: associated with product: %w", exception.ErrFailedPrecondition)
	}
	if err := s.db.ProductType.Delete(ctx, in.ProductTypeID); err != nil {
		return internalError(err)
	}
	s.waitGroup.Add(1)
	go func(productTypeID string) {
		defer s.waitGroup.Done()
		in := &user.RemoveCoordinatorProductTypeInput{
			ProductTypeID: productTypeID,
		}
		if err := s.user.RemoveCoordinatorProductType(context.Background(), in); err != nil {
			s.logger.Error("Failed to remove product type in coordinators",
				zap.String("productTypeId", productTypeID), zap.Error(err))
		}
	}(in.ProductTypeID)
	return nil
}

func (s *service) resizeProductType(ctx context.Context, productTypeID, iconURL string) {
	if iconURL == "" {
		return
	}
	in := &media.ResizeFileInput{
		TargetID: productTypeID,
		URLs:     []string{iconURL},
	}
	if err := s.media.ResizeProductTypeIcon(ctx, in); err != nil {
		s.logger.Error("Failed to resize product type icon",
			zap.String("productTypeId", productTypeID), zap.Error(err),
		)
	}
}
