package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
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
	r.GET("/merchant-feed", h.GetMerchantCenterFeed)
}

// @Summary     商品一覧取得
// @Description 商品の一覧を取得します。
// @Tags        Product
// @Router      /products [get]
// @Param       limit query int64 false "取得上限数(max:200)" default(20)
// @Param       offset query int64 false "取得開始位置(min:0)" default(0)
// @Param       coordinatorId query string false "コーディネータID"
// @Produce     json
// @Success     200 {object} types.ProductsResponse
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
	var shopID string
	if coordinatorID := util.GetQuery(ctx, "coordinatorId", ""); coordinatorID != "" {
		coordinator, err := h.getCoordinator(ctx, coordinatorID)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		shopID = coordinator.ShopID
	}

	in := &store.ListProductsInput{
		Limit:            limit,
		Offset:           offset,
		OnlyPublished:    true,
		ExcludeOutOfSale: true,
		ExcludeDeleted:   true,
		ShopID:           shopID,
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
		res := &types.ProductsResponse{
			Products:     []*types.Product{},
			Coordinators: []*types.Coordinator{},
			Producers:    []*types.Producer{},
			Categories:   []*types.Category{},
			ProductTypes: []*types.ProductType{},
			ProductTags:  []*types.ProductTag{},
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

	res := &types.ProductsResponse{
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
// @Description 商品の詳細情報を取得します。
// @Tags        Product
// @Router      /products/{productId} [get]
// @Param       productId path string true "商品ID"
// @Produce     json
// @Success     200 {object} types.ProductResponse
// @Failure     404 {object} util.ErrorResponse "商品が見つからない"
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

	res := &types.ProductResponse{
		Product:     product.Response(),
		Coordinator: coordinator.Response(),
		Producer:    producer.Response(),
		Category:    category.Response(),
		ProductType: productType.Response(),
		ProductTags: productTags.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) listProducts(ctx context.Context, in *store.ListProductsInput) (service.Products, error) {
	products, _, err := h.store.ListProducts(ctx, in)
	if err != nil || len(products) == 0 {
		return service.Products{}, err
	}
	products = products.FilterByPublished()
	details, err := h.getProductDetails(ctx, products.IDs()...)
	if err != nil {
		return nil, err
	}
	return service.NewProducts(products, details), nil
}

func (h *handler) multiGetProducts(ctx context.Context, productIDs []string) (service.Products, error) {
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
	products = products.FilterByPublished()
	details, err := h.getProductDetails(ctx, products.IDs()...)
	if err != nil {
		return nil, err
	}
	return service.NewProducts(products, details), nil
}

func (h *handler) multiGetProductsByRevision(ctx context.Context, revisionIDs []int64) (service.Products, error) {
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
	products = products.FilterByPublished()
	details, err := h.getProductDetails(ctx, products.IDs()...)
	if err != nil {
		return nil, err
	}
	return service.NewProducts(products, details), nil
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

// @Summary     Merchant Centerフィード取得
// @Description Google Merchant Center用の商品フィードをXML形式で取得します。
// @Tags        Product
// @Router      /products/merchant-feed [get]
// @Produce     xml
// @Success     200 {string} string "XML形式の商品フィード"
func (h *handler) GetMerchantCenterFeed(ctx *gin.Context) {
	const (
		title       = "ふるマル - 全国ふるさとマルシェ"
		description = "地域・地方の特産品を扱うECマーケットプレイス"
		version     = "2.0"
		xmlns       = "http://base.google.com/ns/1.0"
		contentType = "application/xml; charset=utf-8"
	)

	in := &store.ListProductsInput{
		NoLimit:          true,
		OnlyPublished:    true,
		ExcludeOutOfSale: true,
		ExcludeDeleted:   true,
		Orders: []*store.ListProductsOrder{
			{Key: store.ListProductsOrderByUpdatedAt, OrderByASC: false},
		},
	}
	products, _, err := h.store.ListProducts(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.MerchantCenterFeedResponse{
		Version: version,
		Xmlns:   xmlns,
		Channel: &types.MerchantCenterChannel{
			Title:       title,
			Link:        h.userWebURL().String(),
			Description: description,
			Items:       []*types.MerchantCenterItem{},
		},
	}
	if len(products) == 0 {
		ctx.Header("Content-Type", contentType)
		ctx.XML(http.StatusOK, res)
		return
	}

	var (
		details      *service.ProductDetailsParams
		coordinators service.Coordinators
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		details, err = h.getProductDetails(ectx, products.IDs()...)
		return
	})
	eg.Go(func() (err error) {
		coordinators, err = h.multiGetCoordinators(ectx, products.CoordinatorIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	params := &service.NewMerchantCenterItemsParams{
		Products:     service.NewProducts(products, details),
		Coordinators: coordinators.Map(),
		Details:      details,
		WebURL:       h.userWebURL,
	}
	res.Channel.Items = service.NewMerchantCenterItems(params).Response()

	ctx.Header("Content-Type", contentType)
	ctx.XML(http.StatusOK, res)
}
