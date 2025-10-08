package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Product
// @tag.description 商品関連
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
			product, err := h.getProduct(ctx, util.GetParam(ctx, "productId"))
			if err != nil {
				return false, err
			}
			return currentAdmin(ctx, product.CoordinatorID), nil
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

// @Summary     商品一覧取得
// @Description 商品の一覧を取得します。
// @Tags        Product
// @Router      /v1/products [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Param       name query string false "商品名(あいまい検索)" example("新じゃがいも")
// @Param       producerId query string false "生産者ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       orders query string false "ソート" example("-updatedAt")
// @Produce     json
// @Success     200 {object} types.ProductsResponse
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
		ShopID:         getShopID(ctx),
		Name:           util.GetQuery(ctx, "name", ""),
		ProducerID:     util.GetQuery(ctx, "producerId", ""),
		ExcludeDeleted: true,
		Limit:          limit,
		Offset:         offset,
		Orders:         orders,
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

func (h *handler) newProductOrders(ctx *gin.Context) ([]*store.ListProductsOrder, error) {
	products := map[string]store.ListProductsOrderKey{
		"name":             store.ListProductsOrderByName,
		"public":           store.ListProductsOrderByPublic,
		"inventory":        store.ListProductsOrderByInventory,
		"originPrefecture": store.ListProductsOrderByOriginPrefecture,
		"originCity":       store.ListProductsOrderByOriginCity,
		"createdAt":        store.ListProductsOrderByCreatedAt,
		"updatedAt":        store.ListProductsOrderByUpdatedAt,
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

// @Summary     商品取得
// @Description 指定された商品の詳細情報を取得します。
// @Tags        Product
// @Router      /v1/products/{productId} [get]
// @Security    bearerauth
// @Param       productId path string true "商品ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.ProductResponse
// @Failure     404 {object} util.ErrorResponse "商品が存在しない"
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

// @Summary     商品登録
// @Description 新しい商品を登録します。
// @Tags        Product
// @Router      /v1/products [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.CreateProductRequest true "商品情報"
// @Produce     json
// @Success     200 {object} types.ProductResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "商品の登録権限がない"
func (h *handler) CreateProduct(ctx *gin.Context) {
	req := &types.CreateProductRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	if getAdminType(ctx).IsCoordinator() {
		if req.CoordinatorID != getAdminID(ctx) {
			h.forbidden(ctx, errors.New("handler: not authorized this coordinator"))
			return
		}
	}

	shop, err := h.getShopByCoordinatorID(ctx, req.CoordinatorID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if !slices.Contains(shop.ProducerIDs, req.ProducerID) {
		h.forbidden(ctx, errors.New("handler: not authorized this coordinator"))
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
		coordinator, err = h.getCoordinator(ectx, req.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		producer, err = h.getProducer(ectx, req.ProducerID)
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
	err = eg.Wait()
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

	productScope := service.ProductScope(types.ProductScopePrivate)
	if req.Public {
		productScope = service.ProductScope(types.ProductScopePublic)
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
		ShopID:               shop.ID,
		CoordinatorID:        req.CoordinatorID,
		ProducerID:           req.ProducerID,
		TypeID:               req.TypeID,
		TagIDs:               req.TagIDs,
		Name:                 req.Name,
		Description:          req.Description,
		Scope:                productScope.StoreEntity(),
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

	res := &types.ProductResponse{
		Product:     service.NewProduct(sproduct).Response(),
		Coordinator: coordinator.Response(),
		Producer:    producer.Response(),
		Category:    category.Response(),
		ProductType: productType.Response(),
		ProductTags: productTags.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     商品更新
// @Description 商品の情報を更新します。
// @Tags        Product
// @Router      /v1/products/{productId} [patch]
// @Param       productId path string true "商品ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Security    bearerauth
// @Accept      json
// @Param       request body types.UpdateProductRequest true "商品情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "商品の更新権限がない"
// @Failure     404 {object} util.ErrorResponse "商品が存在しない"
func (h *handler) UpdateProduct(ctx *gin.Context) {
	req := &types.UpdateProductRequest{}
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

	productScope := service.ProductScope(types.ProductScopePrivate)
	if req.Public {
		productScope = service.ProductScope(types.ProductScopePublic)
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
		Scope:                productScope.StoreEntity(),
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

// @Summary     商品削除
// @Description 商品を削除します。
// @Tags        Product
// @Router      /v1/products/{productId} [delete]
// @Param       productId path string true "商品ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Security    bearerauth
// @Produce     json
// @Success     204
// @Failure     403 {object} util.ErrorResponse "商品の削除権限がない"
// @Failure     404 {object} util.ErrorResponse "商品が存在しない"
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
