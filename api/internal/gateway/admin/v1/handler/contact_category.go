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
	arg := rg.Use(h.authentication)
	arg.GET("/:contactCategoryId")
}

func (h *handler) GetContactCategory(ctx *gin.Context) {
	category, err := h.getContactCategory(ctx, util.GetParam(ctx, "contactCategoryId"))
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ContactCategoryResponse{
		ContactCategory: category.Response(),
	}
	ctx.JSON(http.StatusOK, res)
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
