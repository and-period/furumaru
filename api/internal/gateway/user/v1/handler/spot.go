package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/gin-gonic/gin"
)

func (h *handler) spotRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/spots")

	r.GET("", h.ListSpots)
	r.GET("/:spotId", h.GetSpot)
	r.POST("", h.CreateSpot)
	r.PATCH("/:spotId", h.UpdateSpot)
	r.DELETE("/:spotId", h.DeleteSpot)
}

func (h *handler) ListSpots(ctx *gin.Context) {
	const defaultRadius = 20

	_, err := util.GetQueryInt64(ctx, "radius", defaultRadius)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	// TODO: 詳細の実装
	res := &response.SpotsResponse{
		Spots: []*response.Spot{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetSpot(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.SpotResponse{}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateSpot(ctx *gin.Context) {
	req := &request.CreateSpotRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	// TODO: 詳細の実装
	res := &response.SpotResponse{}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateSpot(ctx *gin.Context) {
	req := &request.UpdateSpotRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	// TODO: 詳細の実装
	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteSpot(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.Status(http.StatusNoContent)
}
