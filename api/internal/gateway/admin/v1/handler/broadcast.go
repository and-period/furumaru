package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/gin-gonic/gin"
)

func (h *handler) broadcastRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/schedules/:scheduleId/broadcasts", h.authentication, h.filterAccessSchedule)

	r.GET("", h.GetBroadcast)
}

func (h *handler) GetBroadcast(ctx *gin.Context) {
	in := &media.GetBroadcastByScheduleIDInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
	}
	broadcast, err := h.media.GetBroadcastByScheduleID(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.BroadcastResponse{
		Broadcast: service.NewBroadcast(broadcast).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
