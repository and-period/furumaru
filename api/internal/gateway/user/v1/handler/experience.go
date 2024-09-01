package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) experienceRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/experiences")

	r.GET("", h.ListExperiences)
	r.GET("/:experienceId", h.GetExperience)
}

func (h *handler) ListExperiences(ctx *gin.Context) {
	const (
		defaultLimit          = 20
		defaultOffset         = 0
		defaultPrefectureCode = 0
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
	prefectureCode, err := util.GetQueryInt32(ctx, "prefectureCode", defaultPrefectureCode)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.ListExperiencesInput{
		Name:            util.GetQuery(ctx, "name", ""),
		PrefectureCode:  prefectureCode,
		Limit:           limit,
		Offset:          offset,
		OnlyPublished:   true,
		ExcludeFinished: true,
		ExcludeDeleted:  true,
		CoordinatorID:   util.GetQuery(ctx, "coordinatorId", ""),
	}
	experiences, total, err := h.store.ListExperiences(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(experiences) == 0 {
		res := &response.ExperiencesResponse{
			Experiences:     []*response.Experience{},
			Coordinators:    []*response.Coordinator{},
			Producers:       []*response.Producer{},
			ExperienceTypes: []*response.ExperienceType{},
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
		coordinators, err = h.multiGetCoordinators(ectx, experiences.CoordinatorIDs())
		return
	})
	eg.Go(func() (err error) {
		producers, err = h.multiGetProducers(ectx, experiences.ProducerIDs())
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
	experience, err := h.getExperience(ctx, ctx.Param("experienceId"))
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
	return service.NewExperiences(experiences.FilterByPublished()), nil
}

func (h *handler) getExperience(ctx *gin.Context, experienceID string) (*service.Experience, error) {
	in := &store.GetExperienceInput{
		ExperienceID: experienceID,
	}
	experience, err := h.store.GetExperience(ctx, in)
	if err != nil {
		return nil, err
	}
	if !experience.Public {
		// 非公開のものは利用者側に表示しない
		return nil, exception.ErrNotFound
	}
	return service.NewExperience(experience), nil
}
