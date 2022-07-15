package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/gin-gonic/gin"
)

func (h *handler) contactRoutes(rg *gin.RouterGroup) {
	rg.POST("", h.CreateContact)
}

func (h *handler) CreateContact(ctx *gin.Context) {
	req := &request.CreateContactRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &messenger.CreateContactInput{
		Title:       req.Title,
		Content:     req.Content,
		Username:    req.Username,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}
	if _, err := h.messenger.CreateContact(ctx, in); err != nil {
		httpError(ctx, err)
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
