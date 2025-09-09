package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/gin-gonic/gin"
)

// @tag.name        ContactCategory
// @tag.description お問い合わせカテゴリ関連
func (h *handler) contactCategoryRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/contact-categories", h.authentication)

	r.GET("", h.ListContactCategories)
	r.GET("/:contactCategoryId", h.GetContactCategory)
}

// @Summary     お問い合わせカテゴリ一覧取得
// @Description お問い合わせカテゴリの一覧を取得します。
// @Tags        ContactCategory
// @Router      /v1/contact-categories [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Produce     json
// @Success     200 {object} types.ContactCategoriesResponse
func (h *handler) ListContactCategories(ctx *gin.Context) {
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

	in := &messenger.ListContactCategoriesInput{
		Limit:  limit,
		Offset: offset,
	}
	categories, err := h.messenger.ListContactCategories(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.ContactCategoriesResponse{
		ContactCategories: service.NewContactCategories(categories).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     お問い合わせカテゴリ取得
// @Description 指定されたお問い合わせカテゴリの詳細情報を取得します。
// @Tags        ContactCategory
// @Router      /v1/contact-categories/{contactCategoryId} [get]
// @Security    bearerauth
// @Param       contactCategoryId path string true "お問い合わせカテゴリID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.ContactCategoryResponse
// @Failure     404 {object} util.ErrorResponse "お問い合わせカテゴリが存在しない"
func (h *handler) GetContactCategory(ctx *gin.Context) {
	category, err := h.getContactCategory(ctx, util.GetParam(ctx, "contactCategoryId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.ContactCategoryResponse{
		ContactCategory: category.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) multiGetContactCategories(ctx context.Context, categoryIDs []string) (service.ContactCategories, error) {
	if len(categoryIDs) == 0 {
		return service.ContactCategories{}, nil
	}
	in := &messenger.MultiGetContactCategoriesInput{
		CategoryIDs: categoryIDs,
	}
	categories, err := h.messenger.MultiGetContactCategories(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewContactCategories(categories), nil
}

func (h *handler) getContactCategory(ctx context.Context, contactCategoryID string) (*service.ContactCategory, error) {
	in := &messenger.GetContactCategoryInput{
		CategoryID: contactCategoryID,
	}
	mcategory, err := h.messenger.GetContactCategory(ctx, in)
	if err != nil {
		return nil, err
	}
	category := service.NewContactCategory(mcategory)
	return category, nil
}
