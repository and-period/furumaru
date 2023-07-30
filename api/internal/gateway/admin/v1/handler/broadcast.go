package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
)

func (h *handler) broadcastRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication, h.filterAccessSchedule)
	arg.GET("", h.GetBroadcast)
}

func (h *handler) GetBroadcast(ctx *gin.Context) {
	in := &store.GetBroadcastByScheduleIDInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
	}
	broadcast, err := h.store.GetBroadcastByScheduleID(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	res := &response.BroadcastResponse{
		Broadcast: service.NewBroadcast(broadcast).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
