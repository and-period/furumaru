package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/gin-gonic/gin"
)

func (h *handler) uploadRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/upload", h.authentication)

	r.GET("/state", h.GetUploadState)
	r.POST("/users/thumbnail", h.CreateUserThumbnailURL)
}

func (h *handler) GetUploadState(ctx *gin.Context) {
	in := &media.GetUploadEventInput{
		Key: util.GetQuery(ctx, "key", ""),
	}
	event, err := h.media.GetUploadEvent(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.UploadStateResponse{
		URL:    event.ReferenceURL,
		Status: service.NewUploadStatus(event.Status).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateUserThumbnailURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetUserThumbnailUploadURL)
}

func (h *handler) getUploadURL(ctx *gin.Context, fn func(context.Context, *media.GenerateUploadURLInput) (*entity.UploadEvent, error)) {
	req := &request.GetUploadURLRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &media.GenerateUploadURLInput{
		FileType: req.FileType,
	}
	event, err := fn(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.UploadURLResponse{
		Key: event.Key,
		URL: event.UploadURL,
	}
	ctx.JSON(http.StatusOK, res)
}
