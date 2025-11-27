package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

// @tag.name        RelatedProducer
// @tag.description 関連生産者関連
func (h *handler) relatedProducerRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/coordinators/:coordinatorId/producers", h.authentication, h.filterAccessRelatedProducer)

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

// @Summary     関連生産者一覧取得
// @Description 指定されたコーディネーターに関連する生産者の一覧を取得します。
// @Tags        RelatedProducer
// @Router      /v1/coordinators/{coordinatorId}/producers [get]
// @Security    bearerauth
// @Param       coordinatorId path string true "コーディネーターID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Produce     json
// @Success     200 {object} types.ProducersResponse
// @Failure     403 {object} util.ErrorResponse "アクセス権限がない"
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

	coordinator, err := h.getCoordinator(ctx, util.GetParam(ctx, "coordinatorId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	in := &user.ListProducersInput{
		CoordinatorID: coordinator.ID,
		Limit:         limit,
		Offset:        offset,
	}
	producers, total, err := h.user.ListProducers(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.ProducersResponse{
		Producers:    service.NewProducers(producers).Response(),
		Coordinators: []*types.Coordinator{coordinator.Response()},
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}
