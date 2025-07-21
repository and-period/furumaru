package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

func (h *handler) relatedProducerRoutes(rg *gin.RouterGroup) {
	r := rg.Group(
		"/coordinators/:coordinatorId/producers",
		h.authentication,
		h.filterAccessRelatedProducer,
	)

	r.GET("", h.ListRelatedProducers)
}

func (h *handler) filterAccessRelatedProducer(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			coordinatorID := util.GetParam(ctx, "coordinatorId")
			return currentAdmin(ctx, coordinatorID), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Next()
}

func (h *handler) ListRelatedProducers(ctx *gin.Context) {
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

	in := &user.ListProducersInput{
		CoordinatorID: util.GetParam(ctx, "coordinatorId"),
		Limit:         limit,
		Offset:        offset,
	}
	producers, total, err := h.user.ListProducers(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	coordinators, err := h.multiGetCoordinators(ctx, producers.CoordinatorIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ProducersResponse{
		Producers:    service.NewProducers(producers).Response(),
		Coordinators: coordinators.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}
