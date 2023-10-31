package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/gin-gonic/gin"
)

func (h *handler) scheduleRoutes(rg *gin.RouterGroup) {
	rg.GET("/:scheduleId", h.GetSchedule)
}

func (h *handler) GetSchedule(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.ScheduleResponse{
		Lives:     []*response.Live{},
		Producers: []*response.Producer{},
		Products:  []*response.Product{},
	}
	ctx.JSON(http.StatusOK, res)
}
