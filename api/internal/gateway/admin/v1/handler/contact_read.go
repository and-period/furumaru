package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/gin-gonic/gin"
)

func (h *handler) contactReadRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/contact-reads", h.authentication)

	r.POST("", h.CreateContactRead)
}

func (h *handler) CreateContactRead(ctx *gin.Context) {
	req := &request.CreateContactReadRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &messenger.CreateContactReadInput{
		ContactID: req.ContactID,
		UserID:    req.UserID,
		UserType:  entity.ContactUserType(req.UserType),
	}
	scontactRead, err := h.messenger.CreateContactRead(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	contactRead := service.NewContactRead(scontactRead)

	res := &response.ContactReadResponse{
		ContactRead: contactRead.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
