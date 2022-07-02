package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) productRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication())
	arg.GET("", h.ListProducts)
	arg.GET("/:productId", h.GetProduct)
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
	sproducts, err := h.store.ListProducts(ctx, in)
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
		// 生産者取得メソッドの実装
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
		// 品目取得メソッドの実装
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
