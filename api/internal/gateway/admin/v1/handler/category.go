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
	"github.com/gin-gonic/gin"
)

func (h *handler) categoryRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/categories", h.authentication)

	r.GET("", h.ListCategories)
	r.POST("", h.CreateCategory)
	r.PATCH("/:categoryId", h.UpdateCategory)
	r.DELETE("/:categoryId", h.DeleteCategory)
}

func (h *handler) ListCategories(ctx *gin.Context) {
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
	orders, err := h.newCategoryOrders(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.ListCategoriesInput{
		Name:   util.GetQuery(ctx, "name", ""),
		Limit:  limit,
		Offset: offset,
		Orders: orders,
	}
	categories, total, err := h.store.ListCategories(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.CategoriesResponse{
		Categories: service.NewCategories(categories).Response(),
		Total:      total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newCategoryOrders(ctx *gin.Context) ([]*store.ListCategoriesOrder, error) {
	categories := map[string]store.ListCategoriesOrderKey{
		"name": store.ListCategoriesOrderByName,
	}
	params := util.GetOrders(ctx)
	res := make([]*store.ListCategoriesOrder, len(params))
	for i, p := range params {
		key, ok := categories[p.Key]
		if !ok {
			return nil, fmt.Errorf("handler: unknown order key. key=%s: %w", p.Key, errInvalidOrderKey)
		}
		res[i] = &store.ListCategoriesOrder{
			Key:        key,
			OrderByASC: p.Direction == util.OrderByASC,
		}
	}
	return res, nil
}

func (h *handler) CreateCategory(ctx *gin.Context) {
	req := &request.CreateCategoryRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.CreateCategoryInput{
		Name: req.Name,
	}
	category, err := h.store.CreateCategory(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.CategoryResponse{
		Category: service.NewCategory(category).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateCategory(ctx *gin.Context) {
	req := &request.UpdateCategoryRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.UpdateCategoryInput{
		CategoryID: util.GetParam(ctx, "categoryId"),
		Name:       req.Name,
	}
	if err := h.store.UpdateCategory(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteCategory(ctx *gin.Context) {
	in := &store.DeleteCategoryInput{
		CategoryID: util.GetParam(ctx, "categoryId"),
	}
	if err := h.store.DeleteCategory(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) multiGetCategories(ctx context.Context, categoryIDs []string) (service.Categories, error) {
	if len(categoryIDs) == 0 {
		return service.Categories{}, nil
	}
	in := &store.MultiGetCategoriesInput{
		CategoryIDs: categoryIDs,
	}
	categories, err := h.store.MultiGetCategories(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewCategories(categories), nil
}

func (h *handler) getCategory(ctx context.Context, categoryID string) (*service.Category, error) {
	in := &store.GetCategoryInput{
		CategoryID: categoryID,
	}
	category, err := h.store.GetCategory(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewCategory(category), nil
}
