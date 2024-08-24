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
	"golang.org/x/sync/errgroup"
)

func (h *handler) experienceRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/experiences", h.authentication)

	r.GET("", h.ListExperiences)
	r.POST("", h.CreateExperience)
	r.GET("/:experienceId", h.filterAccessExperience, h.GetExperience)
	r.PATCH("/:experienceId", h.filterAccessExperience, h.UpdateExperience)
	r.DELETE("/:experienceId", h.filterAccessExperience, h.DeleteExperience)
}

func (h *handler) filterAccessExperience(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			producers, err := h.getProducersByCoordinatorID(ctx, getAdminID(ctx))
			if err != nil {
				return false, err
			}
			experience, err := h.getExperience(ctx, util.GetParam(ctx, "experienceId"))
			if err != nil {
				return false, err
			}
			return producers.Contains(experience.ProducerID), nil
		},
		producer: func(ctx *gin.Context) (bool, error) {
			experience, err := h.getExperience(ctx, util.GetParam(ctx, "experienceId"))
			if err != nil {
				return false, err
			}
			return experience.ProducerID == getAdminID(ctx), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Next()
}

func (h *handler) ListExperiences(ctx *gin.Context) {
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

	in := &store.ListExperiencesInput{
		Name:       util.GetQuery(ctx, "name", ""),
		ProducerID: util.GetQuery(ctx, "producerId", ""),
		Limit:      limit,
		Offset:     offset,
		NoLimit:    false,
	}
	if getRole(ctx) == service.AdminRoleCoordinator {
		producers, err := h.getProducersByCoordinatorID(ctx, getAdminID(ctx))
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		// 生産者が紐づかない場合、体験が存在しないためアーリーリターンする
		if len(producers) == 0 {
			res := &response.ExperiencesResponse{
				Experiences: []*response.Experience{},
			}
			ctx.JSON(http.StatusOK, res)
			return
		}
	}
	experiences, total, err := h.store.ListExperiences(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(experiences) == 0 {
		res := &response.ExperiencesResponse{
			Experiences: []*response.Experience{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	var (
		coordinators    service.Coordinators
		producers       service.Producers
		experienceTypes service.ExperienceTypes
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinators, err = h.multiGetCoordinators(ctx, experiences.CoordinatorIDs())
		return
	})
	eg.Go(func() (err error) {
		producers, err = h.multiGetProducers(ctx, experiences.ProducerIDs())
		return
	})
	eg.Go(func() (err error) {
		experienceTypes, err = h.multiGetExperienceTypes(ectx, experiences.ExperienceTypeIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ExperiencesResponse{
		Experiences:     service.NewExperiences(experiences).Response(),
		Coordinators:    coordinators.Response(),
		Producers:       producers.Response(),
		ExperienceTypes: experienceTypes.Response(),
		Total:           total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetExperience(ctx *gin.Context) {
	experience, err := h.getExperience(ctx, util.GetParam(ctx, "experienceId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		coordinator    *service.Coordinator
		producer       *service.Producer
		experienceType *service.ExperienceType
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinator(ectx, experience.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		producer, err = h.getProducer(ectx, experience.ProducerID)
		return
	})
	eg.Go(func() (err error) {
		experienceType, err = h.getExperienceType(ectx, experience.ExperienceTypeID)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ExperienceResponse{
		Experience:     experience.Response(),
		Coordinator:    coordinator.Response(),
		Producer:       producer.Response(),
		ExperienceType: experienceType.Response(),
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

func (h *handler) multiGetExperiences(ctx context.Context, experienceIDs []string) (service.Experiences, error) {
	if len(experienceIDs) == 0 {
		return service.Experiences{}, nil
	}
	in := &store.MultiGetExperiencesInput{
		ExperienceIDs: experienceIDs,
	}
	experiences, err := h.store.MultiGetExperiences(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewExperiences(experiences), nil
}

func (h *handler) getExperience(ctx context.Context, experienceID string) (*service.Experience, error) {
	in := &store.GetExperienceInput{
		ExperienceID: experienceID,
	}
	experience, err := h.store.GetExperience(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewExperience(experience), nil
}
