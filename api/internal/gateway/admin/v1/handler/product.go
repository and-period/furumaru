package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) productRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication())
	arg.GET("", h.ListProducts)
	arg.POST("", h.CreateProduct)
	arg.GET("/:productId", h.GetProduct)
	arg.PATCH("/:productId", h.UpdateProduct)
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

	in := &store.ListProductsInput{
		Name:          util.GetQuery(ctx, "name", ""),
		CoordinatorID: util.GetQuery(ctx, "coordinatorId", ""),
		ProducerID:    util.GetQuery(ctx, "producerId", ""),
		Limit:         limit,
		Offset:        offset,
	}
	sproducts, total, err := h.store.ListProducts(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	products := service.NewProducts(sproducts)

	var (
		producers  uentity.Producers
		categories sentity.Categories
		types      sentity.ProductTypes
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &user.MultiGetProducersInput{
			ProducerIDs: products.ProducerIDs(),
		}
		producers, err = h.user.MultiGetProducers(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &store.MultiGetCategoriesInput{
			CategoryIDs: products.CategoryIDs(),
		}
		categories, err = h.store.MultiGetCategories(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &store.MultiGetProductTypesInput{
			ProductTypeIDs: products.ProductTypeIDs(),
		}
		types, err = h.store.MultiGetProductTypes(ectx, in)
		return
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ProductsResponse{
		Products:     products.Response(),
		ProductTypes: service.NewProductTypes(types).Response(),
		Categories:   service.NewCategories(categories).Response(),
		Producers:    service.NewProducers(producers).Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetProduct(ctx *gin.Context) {
	in := &store.GetProductInput{
		ProductID: util.GetParam(ctx, "productId"),
	}
	product, err := h.store.GetProduct(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ProductResponse{
		Product: service.NewProduct(product).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateProduct(ctx *gin.Context) {
	req := &request.CreateProductRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	media := make([]*store.CreateProductMedia, len(req.Media))
	for i := range req.Media {
		media[i] = &store.CreateProductMedia{
			URL:         req.Media[i].URL,
			IsThumbnail: req.Media[i].IsThumbnail,
		}
	}
	weight, weightUnit := service.NewProductWeightFromRequest(req.Weight)

	in := &store.CreateProductInput{
		CoordinatorID:    getAdminID(ctx),
		ProducerID:       req.ProducerID,
		CategoryID:       req.CategoryID,
		TypeID:           req.TypeID,
		Name:             req.Name,
		Description:      req.Description,
		Public:           req.Public,
		Inventory:        req.Inventory,
		Weight:           weight,
		WeightUnit:       weightUnit,
		Item:             1, // 1固定
		ItemUnit:         req.ItemUnit,
		ItemDescription:  req.ItemDescription,
		Media:            media,
		Price:            req.Price,
		DeliveryType:     service.DeliveryType(req.DeliveryType).StoreEntity(),
		Box60Rate:        req.Box60Rate,
		Box80Rate:        req.Box80Rate,
		Box100Rate:       req.Box100Rate,
		OriginPrefecture: req.OriginPrefecture,
		OriginCity:       req.OriginCity,
	}
	product, err := h.store.CreateProduct(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ProductResponse{
		Product: service.NewProduct(product).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateProduct(ctx *gin.Context) {
	req := &request.UpdateProductRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	media := make([]*store.UpdateProductMedia, len(req.Media))
	for i := range req.Media {
		media[i] = &store.UpdateProductMedia{
			URL:         req.Media[i].URL,
			IsThumbnail: req.Media[i].IsThumbnail,
		}
	}
	weight, weightUnit := service.NewProductWeightFromRequest(req.Weight)

	in := &store.UpdateProductInput{
		ProductID:        util.GetParam(ctx, "productId"),
		CoordinatorID:    getAdminID(ctx),
		ProducerID:       req.ProducerID,
		CategoryID:       req.CategoryID,
		TypeID:           req.TypeID,
		Name:             req.Name,
		Description:      req.Description,
		Public:           req.Public,
		Inventory:        req.Inventory,
		Weight:           weight,
		WeightUnit:       weightUnit,
		Item:             1, // 1固定
		ItemUnit:         req.ItemUnit,
		ItemDescription:  req.ItemDescription,
		Media:            media,
		Price:            req.Price,
		DeliveryType:     service.DeliveryType(req.DeliveryType).StoreEntity(),
		Box60Rate:        req.Box60Rate,
		Box80Rate:        req.Box80Rate,
		Box100Rate:       req.Box100Rate,
		OriginPrefecture: req.OriginPrefecture,
		OriginCity:       req.OriginCity,
	}
	if err := h.store.UpdateProduct(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
