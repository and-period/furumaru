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

// @tag.name        ProductType
// @tag.description 品目関連
func (h *handler) productTypeRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/categories/:categoryId/product-types", h.authentication)

	r.GET("", h.ListProductTypes)
	r.POST("", h.CreateProductType)
	r.PATCH("/:productTypeId", h.UpdateProductType)
	r.DELETE("/:productTypeId", h.DeleteProductType)

	rg.GET("/categories/-/product-types", h.authentication, h.ListProductTypes)
}

// @Summary     品目一覧取得
// @Description 品目の一覧を取得します。商品種別ID省略時は全品目を取得します。
// @Tags        ProductType
// @Router      /v1/categories/-/product-types [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Param       name query string false "品目名(あいまい検索)" example("いも")
// @Param       orders query string false "ソート(name,-name)" example("-name")
// @Produce     json
// @Success     200 {object} response.ProductTypesResponse

// @Summary     品目一覧取得
// @Description 品目の一覧を取得します。商品種別ID指定時はその種別の品目のみ取得します。
// @Tags        ProductType
// @Router      /v1/categories/{categoryId}/product-types [get]
// @Security    bearerauth
// @Param       categoryId path string true "商品種別ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Param       name query string false "品目名(あいまい検索)" example("いも")
// @Param       orders query string false "ソート(name,-name)" example("-name")
// @Produce     json
// @Success     200 {object} response.ProductTypesResponse
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
			Categories:   []*response.Category{},
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
	types := map[string]store.ListProductTypesOrderKey{
		"name": store.ListProductTypesOrderByName,
	}
	params := util.GetOrders(ctx)
	res := make([]*store.ListProductTypesOrder, len(params))
	for i, p := range params {
		key, ok := types[p.Key]
		if !ok {
			return nil, fmt.Errorf("handler: unknown order key. key=%s: %w", p.Key, errInvalidOrderKey)
		}
		res[i] = &store.ListProductTypesOrder{
			Key:        key,
			OrderByASC: p.Direction == util.OrderByASC,
		}
	}
	return res, nil
}

// @Summary     品目登録
// @Description 新しい品目を登録します。
// @Tags        ProductType
// @Router      /v1/categories/{categoryId}/product-types [post]
// @Security    bearerauth
// @Param       categoryId path string true "商品種別ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body request.CreateProductTypeRequest true "品目情報"
// @Produce     json
// @Success     200 {object} response.ProductTypeResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     409 {object} util.ErrorResponse "すでに存在する品目名"
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

// @Summary     品目更新
// @Description 品目の情報を更新します。
// @Tags        ProductType
// @Router      /v1/categories/{categoryId}/product-types/{productTypeId} [patch]
// @Security    bearerauth
// @Param       categoryId path string true "商品種別ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       productTypeId path string true "品目ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body request.UpdateProductTypeRequest true "品目情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     404 {object} util.ErrorResponse "品目が存在しない"
// @Failure     409 {object} util.ErrorResponse "すでに存在する品目名"
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

	ctx.Status(http.StatusNoContent)
}

// @Summary     品目削除
// @Description 品目を削除します。
// @Tags        ProductType
// @Router      /v1/categories/{categoryId}/product-types/{productTypeId} [delete]
// @Security    bearerauth
// @Param       categoryId path string true "商品種別ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       productTypeId path string true "品目ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     404 {object} util.ErrorResponse "品目が存在しない"
// @Failure     412 {object} util.ErrorResponse "商品側で紐づいているため削除不可"
func (h *handler) DeleteProductType(ctx *gin.Context) {
	in := &store.DeleteProductTypeInput{
		ProductTypeID: util.GetParam(ctx, "productTypeId"),
	}
	if err := h.store.DeleteProductType(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
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
