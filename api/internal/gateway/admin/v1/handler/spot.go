package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/gin-gonic/gin"
)

func (h *handler) spotRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/spots")

	r.GET("", h.ListSpots)
	r.POST("", h.CreateSpot)
	r.GET("/:spotId", h.GetSpot)
	r.PATCH("/:spotId", h.filterAccessSpot, h.UpdateSpot)
	r.DELETE("/:spotId", h.filterAccessSpot, h.DeleteSpot)
	r.PATCH("/:spotId/approval", h.filterAccessSpot, h.ApproveSpot)
}

func (h *handler) filterAccessSpot(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			spot, err := h.getSpot(ctx, util.GetParam(ctx, "spotId"))
			if err != nil {
				return false, err
			}
			if service.NewSpotUserTypeFromInt32(spot.UserType) != service.SpotUserTypeAdmin {
				return false, nil
			}
			return spot.UserID == getAdminID(ctx), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Next()
}

func (h *handler) ListSpots(ctx *gin.Context) {
	const (
		defaultLimit  = 20
		defaultOffset = 0
	)

	_, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	_, err = util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	// TODO: 詳細の実装
	res := &response.SpotsResponse{
		Spots:  []*response.Spot{},
		Users:  []*response.User{},
		Admins: []*response.Admin{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetSpot(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.SpotResponse{
		Admins: []*response.Admin{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateSpot(ctx *gin.Context) {
	req := &request.CreateSpotRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	// TODO: 詳細の実装
	res := &response.SpotResponse{
		Admins: []*response.Admin{},
	}
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

func (h *handler) ApproveSpot(ctx *gin.Context) {
	req := &request.ApproveSpotRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	// TODO: 詳細の実装
	ctx.Status(http.StatusNoContent)
}

func (h *handler) getSpot(ctx context.Context, spotID string) (*service.Spot, error) {
	return &service.Spot{}, nil
}
