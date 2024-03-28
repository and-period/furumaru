package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/set"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListProducts(ctx context.Context, in *store.ListProductsInput) (entity.Products, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	orders := make([]*database.ListProductsOrder, len(in.Orders))
	for i := range in.Orders {
		orders[i] = &database.ListProductsOrder{
			Key:        in.Orders[i].Key,
			OrderByASC: in.Orders[i].OrderByASC,
		}
	}
	params := &database.ListProductsParams{
		Name:           in.Name,
		CoordinatorID:  in.CoordinatorID,
		ProducerID:     in.ProducerID,
		ProducerIDs:    in.ProducerIDs,
		OnlyPublished:  in.OnlyPublished,
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
	if errors.Is(err, exception.ErrNotFound) {
		return nil, fmt.Errorf("api: invalid admin id: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewProductParams{
		CoordinatorID:        in.CoordinatorID,
		ProducerID:           in.ProducerID,
		TypeID:               in.TypeID,
		TagIDs:               in.TagIDs,
		Name:                 in.Name,
		Description:          in.Description,
		Public:               in.Public,
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
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		s.resizeProduct(context.Background(), product.ID, product.Media)
	}()
	return product, nil
}

func (s *service) UpdateProduct(ctx context.Context, in *store.UpdateProductInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	if _, err := codes.ToPrefectureJapanese(in.OriginPrefectureCode); err != nil {
		return fmt.Errorf("service: invalid prefecture: %w: %s", exception.ErrInvalidArgument, err.Error())
	}
	product, err := s.db.Product.Get(ctx, in.ProductID)
	if err != nil {
		return internalError(err)
	}
	currentMedia := product.Media.MapByURL()
	media := make(entity.MultiProductMedia, len(in.Media))
	for i, m := range in.Media {
		media[i] = entity.NewProductMedia(m.URL, m.IsThumbnail)
		if images, ok := currentMedia[m.URL]; ok {
			media[i].SetImages(images.Images)
		}
	}
	if err := media.Validate(); err != nil {
		return fmt.Errorf("api: invalid media format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	params := &database.UpdateProductParams{
		TypeID:               in.TypeID,
		TagIDs:               in.TagIDs,
		Name:                 in.Name,
		Description:          in.Description,
		Public:               in.Public,
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
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		set := set.NewEmpty[string](len(product.Media))
		for _, m := range product.Media {
			set.Add(m.URL)
		}
		ms := make(entity.MultiProductMedia, 0, len(in.Media))
		for i := range media {
			if set.Contains(media[i].URL) {
				continue
			}
			ms = append(ms, media[i])
		}
		s.resizeProduct(context.Background(), product.ID, ms)
	}()
	return nil
}

func (s *service) UpdateProductMedia(ctx context.Context, in *store.UpdateProductMediaInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	resized := make(map[string]common.Images, len(in.Images))
	for _, image := range in.Images {
		resized[image.OriginURL] = image.Images
	}
	setFn := func(media entity.MultiProductMedia) (exists bool) {
		for i := range media {
			images, ok := resized[media[i].URL]
			if !ok {
				continue
			}
			exists = true
			media[i].SetImages(images)
		}
		return
	}
	err := s.db.Product.UpdateMedia(ctx, in.ProductID, setFn)
	return internalError(err)
}

func (s *service) DeleteProduct(ctx context.Context, in *store.DeleteProductInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Product.Delete(ctx, in.ProductID)
	return internalError(err)
}

func (s *service) resizeProduct(ctx context.Context, productID string, ms entity.MultiProductMedia) {
	if len(ms) == 0 {
		return
	}
	urls := make([]string, len(ms))
	for i := range ms {
		urls[i] = ms[i].URL
	}
	in := &media.ResizeFileInput{
		TargetID: productID,
		URLs:     urls,
	}
	if err := s.media.ResizeProductMedia(ctx, in); err != nil {
		s.logger.Error("Failed to resize product media",
			zap.String("productId", productID), zap.Error(err),
		)
	}
}
