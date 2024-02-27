package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/gin-gonic/gin"
)

func (h *handler) guestLiveCommentRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/guests/schedules/:scheduleId/comments")

	r.POST("", h.createBroadcastViewerLog, h.CreateGuestLiveComment)
}

func (h *handler) CreateGuestLiveComment(ctx *gin.Context) {
	req := &request.CreateGuestLiveCommentRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	schedule, err := h.getSchedule(ctx, util.GetParam(ctx, "scheduleId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	in := &media.CreateBroadcastGuestCommentInput{
		ScheduleID: schedule.ID,
		Content:    req.Comment,
	}
	if _, err := h.media.CreateBroadcastGuestComment(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
