package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
)

func (h *handler) productTypeRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/categories/:categoryId/product-types", h.authentication)

	r.GET("", h.ListProductTypes)
	r.POST("", h.CreateProductType)
	r.PATCH("/:productTypeId", h.UpdateProductType)
	r.DELETE("/:productTypeId", h.DeleteProductType)

	rg.GET("/categories/-/product-types", h.authentication, h.ListProductTypes)
}

func (h *handler) ListProductTypes(ctx *gin.Context) {
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
	orders, err := h.newProductTypeOrders(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	typesIn := &store.ListProductTypesInput{
		Name:       util.GetQuery(ctx, "name", ""),
		CategoryID: util.GetParam(ctx, "categoryId"),
		Limit:      limit,
		Offset:     offset,
		Orders:     orders,
	}
	productTypes, total, err := h.store.ListProductTypes(ctx, typesIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(productTypes) == 0 {
		res := &response.ProductTypesResponse{
			ProductTypes: []*response.ProductType{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	categories, err := h.multiGetCategories(ctx, productTypes.CategoryIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ProductTypesResponse{
		ProductTypes: service.NewProductTypes(productTypes).Response(),
		Categories:   categories.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newProductTypeOrders(ctx *gin.Context) ([]*store.ListProductTypesOrder, error) {
	types := map[string]sentity.ProductTypeOrderBy{
		"name": sentity.ProductTypeOrderByName,
	}
	params := util.GetOrders(ctx)
	res := make([]*store.ListProductTypesOrder, len(params))
	for i, p := range params {
		key, ok := types[p.Key]
		if !ok {
			return nil, fmt.Errorf("handler: unknown order key. key=%s: %w", p.Key, errInvalidOrderkey)
		}
		res[i] = &store.ListProductTypesOrder{
			Key:        key,
			OrderByASC: p.Direction == util.OrderByASC,
		}
	}
	return res, nil
}

func (h *handler) CreateProductType(ctx *gin.Context) {
	req := &request.CreateProductTypeRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	category, err := h.getCategory(ctx, util.GetParam(ctx, "categoryId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	typeIn := &store.CreateProductTypeInput{
		CategoryID: category.ID,
		Name:       req.Name,
		IconURL:    req.IconURL,
	}
	sproductType, err := h.store.CreateProductType(ctx, typeIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	productType := service.NewProductType(sproductType)

	res := &response.ProductTypeResponse{
		ProductType: productType.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateProductType(ctx *gin.Context) {
	req := &request.UpdateProductTypeRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.UpdateProductTypeInput{
		ProductTypeID: util.GetParam(ctx, "productTypeId"),
		Name:          req.Name,
		IconURL:       req.IconURL,
	}
	if err := h.store.UpdateProductType(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) DeleteProductType(ctx *gin.Context) {
	in := &store.DeleteProductTypeInput{
		ProductTypeID: util.GetParam(ctx, "productTypeId"),
	}
	if err := h.store.DeleteProductType(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) multiGetProductTypes(ctx context.Context, productTypeIDs []string) (service.ProductTypes, error) {
	if len(productTypeIDs) == 0 {
		return service.ProductTypes{}, nil
	}
	in := &store.MultiGetProductTypesInput{
		ProductTypeIDs: productTypeIDs,
	}
	sproductTypes, err := h.store.MultiGetProductTypes(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProductTypes(sproductTypes), nil
}

func (h *handler) getProductType(ctx context.Context, productTypeID string) (*service.ProductType, error) {
	in := &store.GetProductTypeInput{
		ProductTypeID: productTypeID,
	}
	sproductType, err := h.store.GetProductType(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProductType(sproductType), nil
}
