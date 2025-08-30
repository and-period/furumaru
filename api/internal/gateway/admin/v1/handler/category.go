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

// @tag.name        Category
// @tag.description 商品種別関連
func (h *handler) categoryRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/categories", h.authentication)

	r.GET("", h.ListCategories)
	r.POST("", h.CreateCategory)
	r.PATCH("/:categoryId", h.UpdateCategory)
	r.DELETE("/:categoryId", h.DeleteCategory)
}

// @Summary     商品種別一覧取得
// @Description 商品種別の一覧を取得します。
// @Tags        Category
// @Router      /v1/categories [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Param       name query string false "商品種別名(あいまい検索)(32文字以内)" example("野菜")
// @Param       orders query string false "ソート(name,-name)" example("-name")
// @Produce     json
// @Success     200 {object} response.CategoriesResponse
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

// @Summary     商品種別登録
// @Description 新しい商品種別を登録します。
// @Tags        Category
// @Router      /v1/categories [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body request.CreateCategoryRequest true "商品種別情報"
// @Produce     json
// @Success     200 {object} response.CategoryResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     409 {object} util.ErrorResponse "すでに存在する商品種別名"
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

// @Summary     商品種別更新
// @Description 商品種別の情報を更新します。
// @Tags        Category
// @Router      /v1/categories/{categoryId} [patch]
// @Security    bearerauth
// @Param       categoryId path string true "商品種別ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body request.UpdateCategoryRequest true "商品種別情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     404 {object} util.ErrorResponse "商品種別が存在しない"
// @Failure     409 {object} util.ErrorResponse "すでに存在する商品種別名"
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

// @Summary     商品種別削除
// @Description 商品種別を削除します。
// @Tags        Category
// @Router      /v1/categories/{categoryId} [delete]
// @Security    bearerauth
// @Param       categoryId path string true "商品種別ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     404 {object} util.ErrorResponse "商品種別が存在しない"
// @Failure     412 {object} util.ErrorResponse "品目側で紐づいているため削除不可"
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
