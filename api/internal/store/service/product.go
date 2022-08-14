package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListProducts(ctx context.Context, in *store.ListProductsInput) (entity.Products, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	orders := make([]*database.ListProductsOrder, len(in.Orders))
	for i := range in.Orders {
		orders[i] = &database.ListProductsOrder{
			Key:        in.Orders[i].Key,
			OrderByASC: in.Orders[i].OrderByASC,
		}
	}
	params := &database.ListProductsParams{
		Name:       in.Name,
		ProducerID: in.ProducerID,
		CreatedBy:  in.CoordinatorID,
		Limit:      int(in.Limit),
		Offset:     int(in.Offset),
		Orders:     orders,
	}
	var (
		products entity.Products
		total    int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		products, err = s.db.Product.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Product.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	return products, total, nil
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
	media := make(entity.MultiProductMedia, len(in.Media))
	for i := range in.Media {
		media[i] = entity.NewProductMedia(in.Media[i].URL, in.Media[i].IsThumbnail)
	}
	if err := media.Validate(); err != nil {
		return nil, fmt.Errorf("api: invalid media format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &user.GetCoordinatorInput{
			CoordinatorID: in.CoordinatorID,
		}
		_, err = s.user.GetCoordinator(ectx, in)
		return err
	})
	eg.Go(func() (err error) {
		in := &user.GetProducerInput{
			ProducerID: in.ProducerID,
		}
		_, err = s.user.GetProducer(ectx, in)
		return err
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		return nil, fmt.Errorf("api: invalid admin id: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, exception.InternalError(err)
	}
	params := &entity.NewProductParams{
		CoordinatorID:    in.CoordinatorID,
		ProducerID:       in.ProducerID,
		CategoryID:       in.CategoryID,
		TypeID:           in.TypeID,
		Name:             in.Name,
		Description:      in.Description,
		Public:           in.Public,
		Inventory:        in.Inventory,
		Weight:           in.Weight,
		WeightUnit:       in.WeightUnit,
		Item:             in.Item,
		ItemUnit:         in.ItemUnit,
		ItemDescription:  in.ItemDescription,
		Media:            media,
		Price:            in.Price,
		DeliveryType:     in.DeliveryType,
		Box60Rate:        in.Box60Rate,
		Box80Rate:        in.Box80Rate,
		Box100Rate:       in.Box100Rate,
		OriginPrefecture: in.OriginPrefecture,
		OriginCity:       in.OriginCity,
	}
	product := entity.NewProduct(params)
	if err := s.db.Product.Create(ctx, product); err != nil {
		return nil, exception.InternalError(err)
	}
	return product, nil
}

func (s *service) UpdateProduct(ctx context.Context, in *store.UpdateProductInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	media := make(entity.MultiProductMedia, len(in.Media))
	for i := range in.Media {
		media[i] = entity.NewProductMedia(in.Media[i].URL, in.Media[i].IsThumbnail)
	}
	if err := media.Validate(); err != nil {
		return fmt.Errorf("api: invalid media format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &user.GetCoordinatorInput{
			CoordinatorID: in.CoordinatorID,
		}
		_, err = s.user.GetCoordinator(ectx, in)
		return err
	})
	eg.Go(func() (err error) {
		in := &user.GetProducerInput{
			ProducerID: in.ProducerID,
		}
		_, err = s.user.GetProducer(ectx, in)
		return err
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		return fmt.Errorf("api: invalid admin id: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return exception.InternalError(err)
	}
	params := &database.UpdateProductParams{
		ProducerID:       in.ProducerID,
		CategoryID:       in.CategoryID,
		TypeID:           in.TypeID,
		Name:             in.Name,
		Description:      in.Description,
		Public:           in.Public,
		Inventory:        in.Inventory,
		Weight:           in.Weight,
		WeightUnit:       in.WeightUnit,
		Item:             in.Item,
		ItemUnit:         in.ItemUnit,
		ItemDescription:  in.ItemDescription,
		Media:            media,
		Price:            in.Price,
		DeliveryType:     in.DeliveryType,
		Box60Rate:        in.Box60Rate,
		Box80Rate:        in.Box80Rate,
		Box100Rate:       in.Box100Rate,
		OriginPrefecture: in.OriginPrefecture,
		OriginCity:       in.OriginCity,
		UpdatedBy:        in.CoordinatorID,
	}
	err = s.db.Product.Update(ctx, in.ProductID, params)
	return exception.InternalError(err)
}

func (s *service) DeleteProduct(ctx context.Context, in *store.DeleteProductInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.Product.Delete(ctx, in.ProductID)
	return exception.InternalError(err)
}
