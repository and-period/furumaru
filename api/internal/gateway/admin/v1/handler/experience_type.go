package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
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
	// TODO: 詳細の実装
	res := &response.ExperienceTypesResponse{
		ExperienceTypes: []*response.ExperienceType{},
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
