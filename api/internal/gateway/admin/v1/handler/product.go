package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) productRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/products", h.authentication)

	r.GET("", h.ListProducts)
	r.POST("", h.CreateProduct)
	r.GET("/:productId", h.filterAccessProduct, h.GetProduct)
	r.PATCH("/:productId", h.filterAccessProduct, h.UpdateProduct)
	r.DELETE("/:productId", h.filterAccessProduct, h.DeleteProduct)
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
		h.httpError(ctx, err)
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
		h.badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	orders, err := h.newProductOrders(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.ListProductsInput{
		Name:           util.GetQuery(ctx, "name", ""),
		ProducerID:     util.GetQuery(ctx, "producerId", ""),
		ExcludeDeleted: true,
		Limit:          limit,
		Offset:         offset,
		Orders:         orders,
	}
	if getRole(ctx) == service.AdminRoleCoordinator {
		producers, err := h.getProducersByCoordinatorID(ctx, getAdminID(ctx))
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		in.ProducerIDs = producers.IDs()
	}
	products, total, err := h.store.ListProducts(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
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
		coordinators service.Coordinators
		producers    service.Producers
		categories   service.Categories
		productTypes service.ProductTypes
		productTags  service.ProductTags
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinators, err = h.multiGetCoordinators(ectx, products.CoordinatorIDs())
		return
	})
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
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	sproducts := service.NewProducts(products)
	sproducts.Fill(productTypes.Map(), categories.Map())

	res := &response.ProductsResponse{
		Products:     sproducts.Response(),
		Coordinators: coordinators.Response(),
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
			return nil, fmt.Errorf("handler: unknown order key. key=%s: %w", p.Key, errInvalidOrderKey)
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
		h.httpError(ctx, err)
		return
	}

	var (
		coordinator *service.Coordinator
		producer    *service.Producer
		category    *service.Category
		productType *service.ProductType
		productTags service.ProductTags
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinator(ectx, product.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		producer, err = h.getProducer(ectx, product.ProducerID)
		return
	})
	eg.Go(func() (err error) {
		productType, err = h.getProductType(ectx, product.ProductTypeID)
		if err != nil {
			return
		}
		category, err = h.getCategory(ectx, productType.CategoryID)
		return
	})
	eg.Go(func() (err error) {
		productTags, err = h.multiGetProductTags(ectx, product.ProductTagIDs)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	product.Fill(category)

	res := &response.ProductResponse{
		Product:     product.Response(),
		Coordinator: coordinator.Response(),
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
		h.badRequest(ctx, err)
		return
	}
	if getRole(ctx).IsCoordinator() {
		if req.CoordinatorID != getAdminID(ctx) {
			h.forbidden(ctx, errors.New("handler: not authorized this coordinator"))
			return
		}
		producers, err := h.getProducersByCoordinatorID(ctx, getAdminID(ctx))
		if err != nil {
			h.httpError(ctx, err)
		}
		if !producers.Contains(req.ProducerID) {
			h.forbidden(ctx, errors.New("handler: not authorized this coordinator"))
			return
		}
	}

	var (
		coordinator *service.Coordinator
		producer    *service.Producer
		category    *service.Category
		productType *service.ProductType
		productTags service.ProductTags
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinator(ectx, req.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		producer, err = h.getProducer(ectx, req.ProducerID)
		if err != nil {
			return
		}
		if producer.CoordinatorID != req.CoordinatorID {
			return fmt.Errorf("handler: unmatch coordinator id: %w", exception.ErrInvalidArgument)
		}
		return
	})
	eg.Go(func() (err error) {
		productType, err = h.getProductType(ectx, req.TypeID)
		if err != nil {
			return
		}
		category, err = h.getCategory(ectx, productType.CategoryID)
		return
	})
	eg.Go(func() (err error) {
		productTags, err = h.multiGetProductTags(ectx, req.TagIDs)
		return
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		h.badRequest(ctx, err)
		return
	}
	if len(productTags) != len(req.TagIDs) {
		h.badRequest(ctx, errors.New("handler: unmatch product tags"))
		return
	}
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	productMedia := make([]*store.CreateProductMedia, len(req.Media))
	for i := range req.Media {
		productMedia[i] = &store.CreateProductMedia{
			URL:         req.Media[i].URL,
			IsThumbnail: req.Media[i].IsThumbnail,
		}
	}
	weight, weightUnit := service.NewProductWeightFromRequest(req.Weight)
	in := &store.CreateProductInput{
		CoordinatorID:        req.CoordinatorID,
		ProducerID:           req.ProducerID,
		TypeID:               req.TypeID,
		TagIDs:               req.TagIDs,
		Name:                 req.Name,
		Description:          req.Description,
		Public:               req.Public,
		Inventory:            req.Inventory,
		Weight:               weight,
		WeightUnit:           weightUnit,
		Item:                 1, // 1固定
		ItemUnit:             req.ItemUnit,
		ItemDescription:      req.ItemDescription,
		Media:                productMedia,
		Price:                req.Price,
		Cost:                 req.Cost,
		ExpirationDate:       req.ExpirationDate,
		RecommendedPoints:    h.newProductPoints(req.RecommendedPoint1, req.RecommendedPoint2, req.RecommendedPoint3),
		StorageMethodType:    service.StorageMethodType(req.StorageMethodType).StoreEntity(),
		DeliveryType:         service.DeliveryType(req.DeliveryType).StoreEntity(),
		Box60Rate:            req.Box60Rate,
		Box80Rate:            req.Box80Rate,
		Box100Rate:           req.Box100Rate,
		OriginPrefectureCode: req.OriginPrefectureCode,
		OriginCity:           req.OriginCity,
		StartAt:              jst.ParseFromUnix(req.StartAt),
		EndAt:                jst.ParseFromUnix(req.EndAt),
	}
	sproduct, err := h.store.CreateProduct(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ProductResponse{
		Product:     service.NewProduct(sproduct).Response(),
		Coordinator: coordinator.Response(),
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
		h.badRequest(ctx, err)
		return
	}

	var productTags service.ProductTags
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		productType, err := h.getProductType(ectx, req.TypeID)
		if err != nil {
			return err
		}
		_, err = h.getCategory(ectx, productType.CategoryID)
		return err
	})
	eg.Go(func() (err error) {
		productTags, err = h.multiGetProductTags(ectx, req.TagIDs)
		return
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		h.badRequest(ctx, err)
		return
	}
	if len(productTags) != len(req.TagIDs) {
		h.badRequest(ctx, errors.New("handler: unmatch product tags"))
		return
	}
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	productMedia := make([]*store.UpdateProductMedia, len(req.Media))
	for i := range req.Media {
		productMedia[i] = &store.UpdateProductMedia{
			URL:         req.Media[i].URL,
			IsThumbnail: req.Media[i].IsThumbnail,
		}
	}
	weight, weightUnit := service.NewProductWeightFromRequest(req.Weight)
	in := &store.UpdateProductInput{
		ProductID:            util.GetParam(ctx, "productId"),
		TypeID:               req.TypeID,
		TagIDs:               req.TagIDs,
		Name:                 req.Name,
		Description:          req.Description,
		Public:               req.Public,
		Inventory:            req.Inventory,
		Weight:               weight,
		WeightUnit:           weightUnit,
		Item:                 1, // 1固定
		ItemUnit:             req.ItemUnit,
		ItemDescription:      req.ItemDescription,
		Media:                productMedia,
		Price:                req.Price,
		Cost:                 req.Cost,
		ExpirationDate:       req.ExpirationDate,
		RecommendedPoints:    h.newProductPoints(req.RecommendedPoint1, req.RecommendedPoint2, req.RecommendedPoint3),
		StorageMethodType:    service.StorageMethodType(req.StorageMethodType).StoreEntity(),
		DeliveryType:         service.DeliveryType(req.DeliveryType).StoreEntity(),
		Box60Rate:            req.Box60Rate,
		Box80Rate:            req.Box80Rate,
		Box100Rate:           req.Box100Rate,
		OriginPrefectureCode: req.OriginPrefectureCode,
		OriginCity:           req.OriginCity,
		StartAt:              jst.ParseFromUnix(req.StartAt),
		EndAt:                jst.ParseFromUnix(req.EndAt),
	}
	if err := h.store.UpdateProduct(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
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
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) multiGetProducts(ctx context.Context, productIDs []string) (service.Products, error) {
	if len(productIDs) == 0 {
		return service.Products{}, nil
	}
	in := &store.MultiGetProductsInput{
		ProductIDs: productIDs,
	}
	products, err := h.store.MultiGetProducts(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProducts(products), nil
}

func (h *handler) multiGetProductsByRevision(ctx context.Context, revisionIDs []int64) (service.Products, error) {
	if len(revisionIDs) == 0 {
		return service.Products{}, nil
	}
	in := &store.MultiGetProductsByRevisionInput{
		ProductRevisionIDs: revisionIDs,
	}
	products, err := h.store.MultiGetProductsByRevision(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProducts(products), nil
}

func (h *handler) getProduct(ctx context.Context, productID string) (*service.Product, error) {
	in := &store.GetProductInput{
		ProductID: productID,
	}
	product, err := h.store.GetProduct(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProduct(product), nil
}
