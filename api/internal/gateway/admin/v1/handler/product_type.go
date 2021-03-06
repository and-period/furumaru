package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *handler) productTypeRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication())
	arg.GET("", h.ListProductTypes)
	arg.POST("", h.CreateProductType)
	arg.PATCH("/:productTypeId", h.UpdateProductType)
	arg.DELETE("/:productTypeId", h.DeleteProductType)
}

func (h *handler) ListProductTypes(ctx *gin.Context) {
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

	typesIn := &store.ListProductTypesInput{
		Name:       util.GetQuery(ctx, "name", ""),
		CategoryID: util.GetParam(ctx, "categoryId"),
		Limit:      limit,
		Offset:     offset,
	}
	stypes, total, err := h.store.ListProductTypes(ctx, typesIn)
	if err != nil {
		httpError(ctx, err)
		return
	}
	productTypes := service.NewProductTypes(stypes)

	categoriesIn := &store.MultiGetCategoriesInput{
		CategoryIDs: productTypes.CategoryIDs(),
	}
	categories, err := h.store.MultiGetCategories(ctx, categoriesIn)
	if err != nil {
		httpError(ctx, err)
		return
	}

	// TODO: 後から実装
	h.logger.Debug("TODO", zap.Any("categories", categories))

	res := &response.ProductTypesResponse{
		ProductTypes: productTypes.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateProductType(ctx *gin.Context) {
	req := &request.CreateProductTypeRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &store.CreateProductTypeInput{
		CategoryID: util.GetParam(ctx, "categoryId"),
		Name:       req.Name,
	}
	productType, err := h.store.CreateProductType(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ProductTypeResponse{
		ProductType: service.NewProductType(productType).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateProductType(ctx *gin.Context) {
	req := &request.UpdateProductTypeRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &store.UpdateProductTypeInput{
		ProductTypeID: util.GetParam(ctx, "productTypeId"),
		Name:          req.Name,
	}
	if err := h.store.UpdateProductType(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) DeleteProductType(ctx *gin.Context) {
	in := &store.DeleteProductTypeInput{
		ProductTypeID: util.GetParam(ctx, "productTypeId"),
	}
	if err := h.store.DeleteProductType(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
