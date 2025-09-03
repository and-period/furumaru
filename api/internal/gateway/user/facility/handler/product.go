package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Product
// @tag.description 商品関連
func (h *handler) productRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/products")

	r.GET("", h.ListProducts)
	r.GET("/:productId", h.GetProduct)
}

// @Summary     商品一覧取得
// @Description 商品の一覧を取得します。
// @Tags        Product
// @Router      /facilities/{facilityId}/products [get]
// @Param       facilityId path string true "施設ID"
// @Param       limit query int64 false "取得件数" default(20)
// @Param       offset query int64 false "取得開始位置" default(0)
// @Produce     json
// @Success     200 {object} response.ProductsResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
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
		Limit:            limit,
		Offset:           offset,
		OnlyPublished:    true,
		ExcludeOutOfSale: true,
		ExcludeDeleted:   true,
		ProducerID:       h.getProducerID(ctx),
		Orders: []*store.ListProductsOrder{
			// 売り切れでないもの順 && 公開日時が新しいもの順
			{Key: store.ListProductsOrderBySoldOut, OrderByASC: true},
			{Key: store.ListProductsOrderByStartAt, OrderByASC: false},
		},
	}
	products, total, err := h.store.ListProducts(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(products) == 0 {
		res := &response.ProductsResponse{
			Products:     []*response.Product{},
			Coordinators: []*response.Coordinator{},
			Producers:    []*response.Producer{},
			Categories:   []*response.Category{},
			ProductTypes: []*response.ProductType{},
			ProductTags:  []*response.ProductTag{},
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
		productRates service.ProductRates
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
	eg.Go(func() (err error) {
		productRates, err = h.aggregateProductRates(ectx, products.IDs()...)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	details := &service.ProductDetailsParams{
		Categories:   categories.Map(),
		ProductTypes: productTypes.Map(),
		ProductRates: productRates.MapByProductID(),
	}
	sproducts := service.NewProducts(products, details)

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

// @Summary     商品詳細取得
// @Description 商品の詳細を取得します。
// @Tags        Product
// @Router      /facilities/{facilityId}/products/{productId} [get]
// @Param       facilityId path string true "施設ID"
// @Param       productId path string true "商品ID"
// @Produce     json
// @Success     200 {object} response.ProductResponse
// @Failure     404 {object} util.ErrorResponse "商品が見つからない"
func (h *handler) GetProduct(ctx *gin.Context) {
	product, err := h.getProduct(ctx, h.getProducerID(ctx), util.GetParam(ctx, "productId"))
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

func (h *handler) multiGetProducts(ctx context.Context, producerID string, productIDs []string) (service.Products, error) {
	if len(productIDs) == 0 {
		return service.Products{}, nil
	}
	in := &store.MultiGetProductsInput{
		ProductIDs: productIDs,
	}
	products, err := h.store.MultiGetProducts(ctx, in)
	if err != nil || len(products) == 0 {
		return service.Products{}, err
	}
	products = products.FilterByProducerID(producerID).FilterByPublished()
	details, err := h.getProductDetails(ctx, products.IDs()...)
	if err != nil {
		return nil, err
	}
	return service.NewProducts(products, details), nil
}

func (h *handler) multiGetProductsByRevision(ctx context.Context, producerID string, revisionIDs []int64) (service.Products, error) {
	if len(revisionIDs) == 0 {
		return service.Products{}, nil
	}
	in := &store.MultiGetProductsByRevisionInput{
		ProductRevisionIDs: revisionIDs,
	}
	products, err := h.store.MultiGetProductsByRevision(ctx, in)
	if err != nil || len(products) == 0 {
		return service.Products{}, err
	}
	products = products.FilterByProducerID(producerID).FilterByPublished()
	details, err := h.getProductDetails(ctx, products.IDs()...)
	if err != nil {
		return nil, err
	}
	return service.NewProducts(products, details), nil
}

func (h *handler) getProduct(ctx context.Context, producerID, productID string) (*service.Product, error) {
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
	if product.ProducerID != producerID {
		return nil, exception.ErrNotFound
	}
	details, err := h.getProductDetails(ctx, productID)
	if err != nil {
		return nil, err
	}
	category := details.Categories[product.TypeID]
	rate := details.ProductRates[productID]
	return service.NewProduct(product, category, rate), nil
}

func (h *handler) getProductDetails(ctx context.Context, productIDs ...string) (*service.ProductDetailsParams, error) {
	var (
		categories   service.Categories
		productTypes service.ProductTypes
		productRates service.ProductRates
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		productTypes, err = h.multiGetProductTypes(ectx, productIDs)
		if err != nil {
			return
		}
		categories, err = h.multiGetCategories(ectx, productTypes.CategoryIDs())
		return
	})
	eg.Go(func() (err error) {
		productRates, err = h.aggregateProductRates(ectx, productIDs...)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	res := &service.ProductDetailsParams{
		Categories:   categories.Map(),
		ProductTypes: productTypes.Map(),
		ProductRates: productRates.MapByProductID(),
	}
	return res, nil
}
