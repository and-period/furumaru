package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
)

func (h *handler) experienceTypeRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/experience-types", h.authentication)

	r.GET("", h.ListExperienceTypes)
	r.POST("", h.CreateExperienceType)
	r.PATCH("/:experienceTypeId", h.UpdateExperienceType)
	r.DELETE("/:experienceTypeId", h.DeleteExperienceType)
}

func (h *handler) ListExperienceTypes(ctx *gin.Context) {
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

	in := &store.ListExperienceTypesInput{
		Name:   util.GetQuery(ctx, "name", ""),
		Limit:  limit,
		Offset: offset,
	}
	types, total, err := h.store.ListExperienceTypes(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ExperienceTypesResponse{
		ExperienceTypes: service.NewExperienceTypes(types).Response(),
		Total:           total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateExperienceType(ctx *gin.Context) {
	req := &request.CreateExperienceTypeRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	// TODO: 詳細の実装
	res := &response.ExperienceTypeResponse{
		ExperienceType: &response.ExperienceType{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateExperienceType(ctx *gin.Context) {
	req := &request.UpdateExperienceTypeRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	// TODO: 詳細の実装
	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteExperienceType(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.Status(http.StatusNoContent)
}

func (h *handler) multiGetExperienceTypes(ctx context.Context, experienceTypeIDs []string) (service.ExperienceTypes, error) {
	if len(experienceTypeIDs) == 0 {
		return service.ExperienceTypes{}, nil
	}
	in := &store.MultiGetExperienceTypesInput{
		ExperienceTypeIDs: experienceTypeIDs,
	}
	types, err := h.store.MultiGetExperienceTypes(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewExperienceTypes(types), nil
}
