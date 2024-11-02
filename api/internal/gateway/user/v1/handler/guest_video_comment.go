package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/gin-gonic/gin"
)

func (h *handler) guestVideoCommentRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/guests/videos/:videoId/comments")

	r.POST("", h.createVideoViewerLog, h.CreateGuestVideoComment)
}

func (h *handler) CreateGuestVideoComment(ctx *gin.Context) {
	req := &request.CreateGuestVideoCommentRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &media.CreateVideoGuestCommentInput{
		VideoID: util.GetParam(ctx, "videoId"),
		Content: req.Comment,
	}
	if _, err := h.media.CreateVideoGuestComment(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
