package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/gin-gonic/gin"
)

func (h *handler) topRoutes(rg *gin.RouterGroup) {
	rg.GET("/common", h.TopCommon)
}

func (h *handler) TopCommon(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.TopCommonResponse{
		Lives:    []*response.TopCommonLive{},
		Archives: []*response.TopCommonArchive{},
	}
	ctx.JSON(http.StatusOK, res)
}
