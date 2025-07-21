package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/gin-gonic/gin"
)

func (h *handler) contactCategoryRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/contact-categories", h.authentication)

	r.GET("", h.ListContactCategories)
	r.GET("/:contactCategoryId")
}

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

	res := &response.ContactCategoriesResponse{
		ContactCategories: service.NewContactCategories(categories).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetContactCategory(ctx *gin.Context) {
	category, err := h.getContactCategory(ctx, util.GetParam(ctx, "contactCategoryId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ContactCategoryResponse{
		ContactCategory: category.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) multiGetContactCategories(
	ctx context.Context,
	categoryIDs []string,
) (service.ContactCategories, error) {
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

func (h *handler) getContactCategory(
	ctx context.Context,
	contactCategoryID string,
) (*service.ContactCategory, error) {
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
