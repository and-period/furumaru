package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/gin-gonic/gin"
)

func (h *handler) experienceRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/experiences", h.authentication)

	r.GET("", h.ListExperiences)
	r.POST("", h.CreateExperience)
	r.GET("/:experienceId", h.GetExperience)
	r.PATCH("/:experienceId", h.UpdateExperience)
	r.DELETE("/:experienceId", h.DeleteExperience)
}

func (h *handler) ListExperiences(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.ExperiencesResponse{
		Experiences:     []*response.Experience{},
		Coordinators:    []*response.Coordinator{},
		Producers:       []*response.Producer{},
		ExperienceTypes: []*response.ExperienceType{},
		Total:           0,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetExperience(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.ExperienceResponse{
		Experience:     &response.Experience{},
		Coordinator:    &response.Coordinator{},
		Producer:       &response.Producer{},
		ExperienceType: &response.ExperienceType{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateExperience(ctx *gin.Context) {
	req := &request.CreateExperienceRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	// TODO: 詳細の実装
	res := &response.ExperienceResponse{
		Experience:     &response.Experience{},
		Coordinator:    &response.Coordinator{},
		Producer:       &response.Producer{},
		ExperienceType: &response.ExperienceType{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateExperience(ctx *gin.Context) {
	req := &request.UpdateExperienceRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	// TODO: 詳細の実装
	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteExperience(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.Status(http.StatusNoContent)
}
