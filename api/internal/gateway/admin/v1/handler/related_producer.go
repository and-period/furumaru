package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

func (h *handler) relatedProducerRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication, h.filterAccessRelatedProducer)
	arg.GET("", h.ListRelatedProducers)
	arg.POST("", h.RelateProducers)
	arg.DELETE("/:producerId", h.filterAccessProducer, h.UnrelateProducer)
}

func (h *handler) filterAccessRelatedProducer(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			coordinatorID := util.GetParam(ctx, "coordinatorId")
			return currentAdmin(ctx, coordinatorID), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		httpError(ctx, err)
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
		badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.ListProducersInput{
		CoordinatorID: util.GetParam(ctx, "coordinatorId"),
		Limit:         limit,
		Offset:        offset,
	}
	producers, total, err := h.user.ListProducers(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	coordinators, err := h.multiGetCoordinators(ctx, producers.CoordinatorIDs())
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ProducersResponse{
		Producers:    service.NewProducers(producers).Response(),
		Coordinators: coordinators.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) RelateProducers(ctx *gin.Context) {
	req := &request.RelateProducersRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	coordinator, err := h.getCoordinator(ctx, util.GetParam(ctx, "coordinatorId"))
	if err != nil {
		httpError(ctx, err)
		return
	}

	in := &user.RelateProducersInput{
		CoordinatorID: coordinator.ID,
		ProducerIDs:   req.ProducerIDs,
	}
	if err := h.user.RelateProducers(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) UnrelateProducer(ctx *gin.Context) {
	if _, err := h.getCoordinator(ctx, util.GetParam(ctx, "coordinatorId")); err != nil {
		httpError(ctx, err)
		return
	}

	in := &user.UnrelateProducerInput{
		ProducerID: util.GetParam(ctx, "producerId"),
	}
	if err := h.user.UnrelateProducer(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
