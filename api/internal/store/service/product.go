package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListProducts(ctx context.Context, in *store.ListProductsInput) (entity.Products, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	orders, err := s.newListProductsOrders(in.Orders)
	if err != nil {
		return nil, 0, fmt.Errorf("service: invalid list products orders: err=%s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	params := &database.ListProductsParams{
		Name:           in.Name,
		ShopID:         in.ShopID,
		ProducerID:     in.ProducerID,
		ProducerIDs:    in.ProducerIDs,
		Scopes:         in.Scopes,
		ExcludeDeleted: in.ExcludeDeleted,
		Limit:          int(in.Limit),
		Offset:         int(in.Offset),
		Orders:         orders,
	}
	if in.ExcludeOutOfSale {
		params.EndAtGte = s.now()
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
		return nil, 0, internalError(err)
	}
	return products, total, nil
}

func (s *service) newListProductsOrders(in []*store.ListProductsOrder) ([]*database.ListProductsOrder, error) {
	res := make([]*database.ListProductsOrder, len(in))
	for i := range in {
		var key database.ListProductsOrderKey
		switch in[i].Key {
		case store.ListProductsOrderByName:
			key = database.ListProductsOrderByName
		case store.ListProductsOrderBySoldOut:
			key = database.ListProductsOrderBySoldOut
		case store.ListProductsOrderByPublic:
			key = database.ListProductsOrderByPublic
		case store.ListProductsOrderByInventory:
			key = database.ListProductsOrderByInventory
		case store.ListProductsOrderByOriginPrefecture:
			key = database.ListProductsOrderByOriginPrefecture
		case store.ListProductsOrderByOriginCity:
			key = database.ListProductsOrderByOriginCity
		case store.ListProductsOrderByStartAt:
			key = database.ListProductsOrderByStartAt
		case store.ListProductsOrderByCreatedAt:
			key = database.ListProductsOrderByCreatedAt
		case store.ListProductsOrderByUpdatedAt:
			key = database.ListProductsOrderByUpdatedAt
		default:
			return nil, errors.New("service: invalid order key")
		}
		res[i] = &database.ListProductsOrder{
			Key:        key,
			OrderByASC: in[i].OrderByASC,
		}
	}
	return res, nil
}

func (s *service) MultiGetProducts(
	ctx context.Context, in *store.MultiGetProductsInput,
) (entity.Products, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	products, err := s.db.Product.MultiGet(ctx, in.ProductIDs)
	return products, internalError(err)
}

func (s *service) MultiGetProductsByRevision(
	ctx context.Context, in *store.MultiGetProductsByRevisionInput,
) (entity.Products, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	products, err := s.db.Product.MultiGetByRevision(ctx, in.ProductRevisionIDs)
	return products, internalError(err)
}

func (s *service) GetProduct(ctx context.Context, in *store.GetProductInput) (*entity.Product, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	product, err := s.db.Product.Get(ctx, in.ProductID)
	return product, internalError(err)
}

func (s *service) CreateProduct(ctx context.Context, in *store.CreateProductInput) (*entity.Product, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
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
		_, err = s.db.Shop.Get(ectx, in.ShopID)
		return
	})
	eg.Go(func() (err error) {
		in := &user.GetCoordinatorInput{
			CoordinatorID: in.CoordinatorID,
		}
		_, err = s.user.GetCoordinator(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &user.GetProducerInput{
			ProducerID: in.ProducerID,
		}
		_, err = s.user.GetProducer(ectx, in)
		return
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) || errors.Is(err, database.ErrNotFound) {
		return nil, fmt.Errorf("api: invalid admin id: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewProductParams{
		ShopID:               in.ShopID,
		CoordinatorID:        in.CoordinatorID,
		ProducerID:           in.ProducerID,
		TypeID:               in.TypeID,
		TagIDs:               in.TagIDs,
		Name:                 in.Name,
		Description:          in.Description,
		Scope:                in.Scope,
		Inventory:            in.Inventory,
		Weight:               in.Weight,
		WeightUnit:           in.WeightUnit,
		Item:                 in.Item,
		ItemUnit:             in.ItemUnit,
		ItemDescription:      in.ItemDescription,
		Media:                media,
		Price:                in.Price,
		Cost:                 in.Cost,
		ExpirationDate:       in.ExpirationDate,
		RecommendedPoints:    in.RecommendedPoints,
		StorageMethodType:    in.StorageMethodType,
		DeliveryType:         in.DeliveryType,
		Box60Rate:            in.Box60Rate,
		Box80Rate:            in.Box80Rate,
		Box100Rate:           in.Box100Rate,
		OriginPrefectureCode: in.OriginPrefectureCode,
		OriginCity:           in.OriginCity,
		StartAt:              in.StartAt,
		EndAt:                in.EndAt,
	}
	product, err := entity.NewProduct(params)
	if err != nil {
		return nil, fmt.Errorf("service: failed to new product: %w: %s", exception.ErrInvalidArgument, err.Error())
	}
	if err := s.db.Product.Create(ctx, product); err != nil {
		return nil, internalError(err)
	}
	return product, nil
}

func (s *service) UpdateProduct(ctx context.Context, in *store.UpdateProductInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	if _, err := codes.ToPrefectureJapanese(in.OriginPrefectureCode); err != nil {
		return fmt.Errorf("service: invalid prefecture: %w: %s", exception.ErrInvalidArgument, err.Error())
	}
	media := make(entity.MultiProductMedia, len(in.Media))
	for i, m := range in.Media {
		media[i] = entity.NewProductMedia(m.URL, m.IsThumbnail)
	}
	if err := media.Validate(); err != nil {
		return fmt.Errorf("api: invalid media format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	params := &database.UpdateProductParams{
		TypeID:               in.TypeID,
		TagIDs:               in.TagIDs,
		Name:                 in.Name,
		Description:          in.Description,
		Scope:                in.Scope,
		Inventory:            in.Inventory,
		Weight:               in.Weight,
		WeightUnit:           in.WeightUnit,
		Item:                 in.Item,
		ItemUnit:             in.ItemUnit,
		ItemDescription:      in.ItemDescription,
		Media:                media,
		Price:                in.Price,
		Cost:                 in.Cost,
		ExpirationDate:       in.ExpirationDate,
		RecommendedPoints:    in.RecommendedPoints,
		StorageMethodType:    in.StorageMethodType,
		DeliveryType:         in.DeliveryType,
		Box60Rate:            in.Box60Rate,
		Box80Rate:            in.Box80Rate,
		Box100Rate:           in.Box100Rate,
		OriginPrefectureCode: in.OriginPrefectureCode,
		OriginCity:           in.OriginCity,
		StartAt:              in.StartAt,
		EndAt:                in.EndAt,
	}
	if err := s.db.Product.Update(ctx, in.ProductID, params); err != nil {
		return internalError(err)
	}
	return nil
}

func (s *service) DeleteProduct(ctx context.Context, in *store.DeleteProductInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	videosIn := &media.ListProductVideosInput{
		ProductID: in.ProductID,
	}
	videos, err := s.media.ListProductVideos(ctx, videosIn)
	if err != nil {
		return internalError(err)
	}
	if len(videos) > 0 {
		return fmt.Errorf("service: product has videos: %w", exception.ErrFailedPrecondition)
	}
	err = s.db.Product.Delete(ctx, in.ProductID)
	return internalError(err)
}
