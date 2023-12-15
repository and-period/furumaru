package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) productRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/products")

	r.GET("", h.ListProducts)
	r.GET("/:productId", h.GetProduct)
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

	in := &store.ListProductsInput{
		Limit:         limit,
		Offset:        offset,
		OnlyPublished: true,
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
	return service.NewProducts(products.FilterByPublished()), nil
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
	return service.NewProducts(products.FilterByPublished()), nil
}

func (h *handler) getProduct(ctx context.Context, productID string) (*service.Product, error) {
	in := &store.GetProductInput{
		ProductID: productID,
	}
	product, err := h.store.GetProduct(ctx, in)
	if err != nil {
		return nil, err
	}
	if !product.Public {
		// 非公開のものは利用者側に表示しない
		return nil, exception.ErrNotFound
	}
	return service.NewProduct(product), nil
}
