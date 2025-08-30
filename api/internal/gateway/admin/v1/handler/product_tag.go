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

// @tag.name        ProductTag
// @tag.description 商品タグ関連
func (h *handler) productTagRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/product-tags", h.authentication)

	r.GET("", h.ListProductTags)
	r.POST("", h.CreateProductTag)
	r.PATCH("/:productTagId", h.UpdateProductTag)
	r.DELETE("/:productTagId", h.DeleteProductTag)
}

// @Summary     商品タグ一覧取得
// @Description 商品タグの一覧を取得します。名前でのフィルタリング、ソート順指定が可能です。
// @Tags        ProductTag
// @Router      /v1/product-tags [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Param       name query string false "商品タグ名(あいまい検索)" example("有機")
// @Param       orders query string false "ソート(name,-name)" example("-name")
// @Produce     json
// @Success     200 {object} response.ProductTagsResponse
func (h *handler) ListProductTags(ctx *gin.Context) {
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
	orders, err := h.newProductTagOrders(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.ListProductTagsInput{
		Name:   util.GetQuery(ctx, "name", ""),
		Limit:  limit,
		Offset: offset,
		Orders: orders,
	}
	productTags, total, err := h.store.ListProductTags(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ProductTagsResponse{
		ProductTags: service.NewProductTags(productTags).Response(),
		Total:       total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newProductTagOrders(ctx *gin.Context) ([]*store.ListProductTagsOrder, error) {
	productTags := map[string]store.ListProductTagsOrderKey{
		"name": store.ListProductTagsOrderByName,
	}
	params := util.GetOrders(ctx)
	res := make([]*store.ListProductTagsOrder, len(params))
	for i, p := range params {
		key, ok := productTags[p.Key]
		if !ok {
			return nil, fmt.Errorf("handler: unknown order key. key=%s: %w", p.Key, errInvalidOrderKey)
		}
		res[i] = &store.ListProductTagsOrder{
			Key:        key,
			OrderByASC: p.Direction == util.OrderByASC,
		}
	}
	return res, nil
}

// @Summary     商品タグ登録
// @Description 新しい商品タグを登録します。
// @Tags        ProductTag
// @Router      /v1/product-tags [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body request.CreateProductTagRequest true "商品タグ情報"
// @Produce     json
// @Success     200 {object} response.ProductTagResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     409 {object} util.ErrorResponse "すでに存在する商品タグ名"
func (h *handler) CreateProductTag(ctx *gin.Context) {
	req := &request.CreateProductTagRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.CreateProductTagInput{
		Name: req.Name,
	}
	productTag, err := h.store.CreateProductTag(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ProductTagResponse{
		ProductTag: service.NewProductTag(productTag).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     商品タグ更新
// @Description 商品タグの情報を更新します。
// @Tags        ProductTag
// @Router      /v1/product-tags/{productTagId} [patch]
// @Security    bearerauth
// @Param       productTagId path string true "商品タグID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body request.UpdateProductTagRequest true "商品タグ情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     404 {object} util.ErrorResponse "商品タグが存在しない"
// @Failure     409 {object} util.ErrorResponse "すでに存在する商品タグ名"
func (h *handler) UpdateProductTag(ctx *gin.Context) {
	req := &request.UpdateProductTagRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.UpdateProductTagInput{
		ProductTagID: util.GetParam(ctx, "productTagId"),
		Name:         req.Name,
	}
	if err := h.store.UpdateProductTag(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     商品タグ削除
// @Description 商品タグを削除します。
// @Tags        ProductTag
// @Router      /v1/product-tags/{productTagId} [delete]
// @Security    bearerauth
// @Param       productTagId path string true "商品タグID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     404 {object} util.ErrorResponse "商品タグが存在しない"
// @Failure     412 {object} util.ErrorResponse "商品側で紐づいているため削除不可"
func (h *handler) DeleteProductTag(ctx *gin.Context) {
	in := &store.DeleteProductTagInput{
		ProductTagID: util.GetParam(ctx, "productTagId"),
	}
	if err := h.store.DeleteProductTag(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) multiGetProductTags(ctx context.Context, productTagIDs []string) (service.ProductTags, error) {
	if len(productTagIDs) == 0 {
		return service.ProductTags{}, nil
	}
	in := &store.MultiGetProductTagsInput{
		ProductTagIDs: productTagIDs,
	}
	productTags, err := h.store.MultiGetProductTags(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProductTags(productTags), nil
}
