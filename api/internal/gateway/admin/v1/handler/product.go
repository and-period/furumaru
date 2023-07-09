package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) productRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("", h.ListProducts)
	arg.POST("", h.CreateProduct)
	arg.GET("/:productId", h.filterAccessProduct, h.GetProduct)
	arg.PATCH("/:productId", h.filterAccessProduct, h.UpdateProduct)
	arg.DELETE("/:productId", h.filterAccessProduct, h.DeleteProduct)
}

func (h *handler) filterAccessProduct(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			producers, err := h.getProducersByCoordinatorID(ctx, getAdminID(ctx))
			if err != nil {
				return false, err
			}
			product, err := h.getProduct(ctx, util.GetParam(ctx, "productId"))
			if err != nil {
				return false, err
			}
			return producers.Contains(product.ProducerID), nil
		},
		producer: func(ctx *gin.Context) (bool, error) {
			product, err := h.getProduct(ctx, util.GetParam(ctx, "productId"))
			if err != nil {
				return false, err
			}
			return currentAdmin(ctx, product.ProducerID), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.Next()
}

func (h *handler) ListProducts(ctx *gin.Context) {
	const (
		defaultLimit  = 20
		defaultOffset = 0
	)

	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	orders, err := h.newProductOrders(ctx)
	if err != nil {
		badRequest(ctx, err)
		return
	}

	in := &store.ListProductsInput{
		Name:       util.GetQuery(ctx, "name", ""),
		ProducerID: util.GetQuery(ctx, "producerId", ""),
		Limit:      limit,
		Offset:     offset,
		Orders:     orders,
	}
	if getRole(ctx) == service.AdminRoleCoordinator {
		producers, err := h.getProducersByCoordinatorID(ctx, getAdminID(ctx))
		if err != nil {
			httpError(ctx, err)
			return
		}
		in.ProducerIDs = producers.IDs()
	}
	products, total, err := h.store.ListProducts(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	if len(products) == 0 {
		res := &response.ProductsResponse{
			Products: []*response.Product{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	var (
		producers    service.Producers
		categories   service.Categories
		productTypes service.ProductTypes
		productTags  service.ProductTags
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		producers, err = h.multiGetProducers(ectx, products.ProducerIDs())
		return
	})
	eg.Go(func() (err error) {
		productTypes, err = h.multiGetProductTypes(ectx, products.ProductTypeIDs())
		if err != nil || len(productTypes) == 0 {
			return
		}
		categories, err = h.multiGetCategories(ectx, productTypes.CategoryIDs())
		return
	})
	eg.Go(func() (err error) {
		productTags, err = h.multiGetProductTags(ectx, products.ProductTagIDs())
		return
	})

	res := &response.ProductsResponse{
		Products:     service.NewProducts(products).Response(),
		Producers:    producers.Response(),
		Categories:   categories.Response(),
		ProductTypes: productTypes.Response(),
		ProductTags:  productTags.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newProductOrders(ctx *gin.Context) ([]*store.ListProductsOrder, error) {
	products := map[string]sentity.ProductOrderBy{
		"name":             sentity.ProductOrderByName,
		"public":           sentity.ProductOrderByPublic,
		"inventory":        sentity.ProductOrderByInventory,
		"price":            sentity.ProductOrderByPrice,
		"originPrefecture": sentity.ProductOrderByOriginPrefecture,
		"originCity":       sentity.ProductOrderByOriginCity,
		"createdAt":        sentity.ProductOrderByCreatedAt,
		"updatedAt":        sentity.ProductOrderByUpdatedAt,
	}
	params := util.GetOrders(ctx)
	res := make([]*store.ListProductsOrder, len(params))
	for i, p := range params {
		key, ok := products[p.Key]
		if !ok {
			return nil, fmt.Errorf("handler: unknown order key. key=%s: %w", p.Key, errInvalidOrderkey)
		}
		res[i] = &store.ListProductsOrder{
			Key:        key,
			OrderByASC: p.Direction == util.OrderByASC,
		}
	}
	return res, nil
}

func (h *handler) GetProduct(ctx *gin.Context) {
	product, err := h.getProduct(ctx, util.GetParam(ctx, "productId"))
	if err != nil {
		httpError(ctx, err)
		return
	}

	var (
		producer    *service.Producer
		category    *service.Category
		productType *service.ProductType
		productTags service.ProductTags
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		producer, err = h.getProducer(ectx, product.ProducerID)
		return
	})
	eg.Go(func() (err error) {
		productType, err = h.getProductType(ectx, product.ProductTypeID)
		if err != nil {
			return
		}
		category, err = h.getCategory(ectx, productType.ID)
		return
	})
	eg.Go(func() (err error) {
		productTags, err = h.multiGetProductTags(ectx, product.ProductTagIDs)
		return
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ProductResponse{
		Product:     product.Response(),
		Producer:    producer.Response(),
		Category:    category.Response(),
		ProductType: productType.Response(),
		ProductTags: productTags.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateProduct(ctx *gin.Context) {
	req := &request.CreateProductRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}
	if getRole(ctx).IsCoordinator() {
		producers, err := h.getProducersByCoordinatorID(ctx, getAdminID(ctx))
		if err != nil {
			httpError(ctx, err)
		}
		if !producers.Contains(req.ProducerID) {
			forbidden(ctx, errors.New("handler: not authorized this coordinator"))
			return
		}
	}

	var (
		producer    *service.Producer
		category    *service.Category
		productType *service.ProductType
		productTags service.ProductTags
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		producer, err = h.getProducer(ectx, req.ProducerID)
		return
	})
	eg.Go(func() (err error) {
		productType, err = h.getProductType(ectx, req.TypeID)
		if err != nil {
			return
		}
		category, err = h.getCategory(ectx, productType.ID)
		return
	})
	eg.Go(func() (err error) {
		productTags, err = h.multiGetProductTags(ectx, req.TagIDs)
		return
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		badRequest(ctx, err)
		return
	}
	if len(productTags) != len(req.TagIDs) {
		badRequest(ctx, errors.New("handler: unmatch product tags"))
		return
	}
	if err != nil {
		httpError(ctx, err)
		return
	}

	eg, ectx = errgroup.WithContext(ctx)
	productMedia := make([]*store.CreateProductMedia, len(req.Media))
	for i := range req.Media {
		i := i
		eg.Go(func() error {
			in := &media.UploadFileInput{
				URL: req.Media[i].URL,
			}
			url, err := h.media.UploadProductMedia(ectx, in)
			if err != nil {
				return err
			}
			productMedia[i] = &store.CreateProductMedia{
				URL:         url,
				IsThumbnail: req.Media[i].IsThumbnail,
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	weight, weightUnit := service.NewProductWeightFromRequest(req.Weight)
	in := &store.CreateProductInput{
		ProducerID:        req.ProducerID,
		TypeID:            req.TypeID,
		TagIDs:            req.TagIDs,
		Name:              req.Name,
		Description:       req.Description,
		Public:            req.Public,
		Inventory:         req.Inventory,
		Weight:            weight,
		WeightUnit:        weightUnit,
		Item:              1, // 1固定
		ItemUnit:          req.ItemUnit,
		ItemDescription:   req.ItemDescription,
		Media:             productMedia,
		Price:             req.Price,
		Cost:              req.Cost,
		ExpirationDate:    req.ExpirationDate,
		RecommendedPoints: h.newProductPoints(req.RecommendedPoint1, req.RecommendedPoint2, req.RecommendedPoint3),
		StorageMethodType: service.StorageMethodType(req.StorageMethodType).StoreEntity(),
		DeliveryType:      service.DeliveryType(req.DeliveryType).StoreEntity(),
		Box60Rate:         req.Box60Rate,
		Box80Rate:         req.Box80Rate,
		Box100Rate:        req.Box100Rate,
		OriginPrefecture:  codes.PrefectureValues[req.OriginPrefecture],
		OriginCity:        req.OriginCity,
		BusinessDays:      req.BusinessDays,
		StartAt:           jst.ParseFromUnix(req.StartAt),
		EndAt:             jst.ParseFromUnix(req.EndAt),
	}
	sproduct, err := h.store.CreateProduct(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ProductResponse{
		Product:     service.NewProduct(sproduct).Response(),
		Producer:    producer.Response(),
		Category:    category.Response(),
		ProductType: productType.Response(),
		ProductTags: productTags.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateProduct(ctx *gin.Context) {
	req := &request.UpdateProductRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	var productTags service.ProductTags
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		productType, err := h.getProductType(ectx, req.TypeID)
		if err != nil {
			return err
		}
		_, err = h.getCategory(ectx, productType.ID)
		return err
	})
	eg.Go(func() (err error) {
		productTags, err = h.multiGetProductTags(ectx, req.TagIDs)
		return
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		badRequest(ctx, err)
		return
	}
	if len(productTags) != len(req.TagIDs) {
		badRequest(ctx, errors.New("handler: unmatch product tags"))
		return
	}
	if err != nil {
		httpError(ctx, err)
		return
	}

	eg, ectx = errgroup.WithContext(ctx)
	productMedia := make([]*store.UpdateProductMedia, len(req.Media))
	for i := range req.Media {
		i := i
		eg.Go(func() error {
			in := &media.UploadFileInput{
				URL: req.Media[i].URL,
			}
			url, err := h.media.UploadProductMedia(ectx, in)
			if err != nil {
				return err
			}
			productMedia[i] = &store.UpdateProductMedia{
				URL:         url,
				IsThumbnail: req.Media[i].IsThumbnail,
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	weight, weightUnit := service.NewProductWeightFromRequest(req.Weight)
	in := &store.UpdateProductInput{
		ProductID:         util.GetParam(ctx, "productId"),
		ProducerID:        req.ProducerID,
		TypeID:            req.TypeID,
		TagIDs:            req.TagIDs,
		Name:              req.Name,
		Description:       req.Description,
		Public:            req.Public,
		Inventory:         req.Inventory,
		Weight:            weight,
		WeightUnit:        weightUnit,
		Item:              1, // 1固定
		ItemUnit:          req.ItemUnit,
		ItemDescription:   req.ItemDescription,
		Media:             productMedia,
		Price:             req.Price,
		Cost:              req.Cost,
		ExpirationDate:    req.ExpirationDate,
		RecommendedPoints: h.newProductPoints(req.RecommendedPoint1, req.RecommendedPoint2, req.RecommendedPoint3),
		StorageMethodType: service.StorageMethodType(req.StorageMethodType).StoreEntity(),
		DeliveryType:      service.DeliveryType(req.DeliveryType).StoreEntity(),
		Box60Rate:         req.Box60Rate,
		Box80Rate:         req.Box80Rate,
		Box100Rate:        req.Box100Rate,
		OriginPrefecture:  codes.PrefectureValues[req.OriginPrefecture],
		OriginCity:        req.OriginCity,
		BusinessDays:      req.BusinessDays,
		StartAt:           jst.ParseFromUnix(req.StartAt),
		EndAt:             jst.ParseFromUnix(req.EndAt),
	}
	if err := h.store.UpdateProduct(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) newProductPoints(points ...string) []string {
	res := make([]string, 0, len(points))
	for _, point := range points {
		if point == "" {
			continue
		}
		res = append(res, point)
	}
	return res
}

func (h *handler) DeleteProduct(ctx *gin.Context) {
	in := &store.DeleteProductInput{
		ProductID: util.GetParam(ctx, "productId"),
	}
	if err := h.store.DeleteProduct(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) multiGetProducts(ctx context.Context, productIDs []string) (service.Products, error) {
	if len(productIDs) == 0 {
		return service.Products{}, nil
	}
	in := &store.MultiGetProductsInput{
		ProductIDs: productIDs,
	}
	sproducts, err := h.store.MultiGetProducts(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProducts(sproducts), nil
}

func (h *handler) getProduct(ctx context.Context, productID string) (*service.Product, error) {
	in := &store.GetProductInput{
		ProductID: productID,
	}
	sproduct, err := h.store.GetProduct(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProduct(sproduct), nil
}
