package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
)

func (h *handler) spotTypeRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/spot-types")

	r.GET("", h.ListSpotTypes)
}

func (h *handler) ListSpotTypes(ctx *gin.Context) {
	const (
		defaultLimit  = 20
		defaultOffset = 0
	)

	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.ListSpotTypesInput{
		Name:   util.GetQuery(ctx, "name", ""),
		Limit:  limit,
		Offset: offset,
	}
	types, total, err := h.store.ListSpotTypes(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.SpotTypesResponse{
		SpotTypes: service.NewSpotTypes(types).Response(),
		Total:     total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) multiGetSpotTypes(
	ctx context.Context,
	spotTypeIDs []string,
) (service.SpotTypes, error) {
	if len(spotTypeIDs) == 0 {
		return service.SpotTypes{}, nil
	}
	in := &store.MultiGetSpotTypesInput{
		SpotTypeIDs: spotTypeIDs,
	}
	spotTypes, err := h.store.MultiGetSpotTypes(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewSpotTypes(spotTypes), nil
}

func (h *handler) getSpotType(ctx context.Context, spotTypeID string) (*service.SpotType, error) {
	in := &store.GetSpotTypeInput{
		SpotTypeID: spotTypeID,
	}
	spotType, err := h.store.GetSpotType(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewSpotType(spotType), nil
}
