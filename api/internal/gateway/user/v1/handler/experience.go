package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) experienceRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/experiences")

	r.GET("", h.ListExperiences)
	r.GET("/geolocation", h.ListExperiencesByGeolocation)
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
		ProducerID:      util.GetQuery(ctx, "producerId", ""),
	}
	experiences, total, err := h.store.ListExperiences(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res, err := h.newExperiencesResponse(ctx, experiences, total)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) ListExperiencesByGeolocation(ctx *gin.Context) {
	const defaultRadius = 20

	if _, ok := ctx.GetQuery("latitude"); !ok {
		h.badRequest(ctx, errors.New("handler: latitude is required"))
		return
	}
	if _, ok := ctx.GetQuery("longitude"); !ok {
		h.badRequest(ctx, errors.New("handler: longitude is required"))
		return
	}

	radius, err := util.GetQueryInt64(ctx, "radius", defaultRadius)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	latitude, err := util.GetQueryFloat64(ctx, "latitude", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	longitude, err := util.GetQueryFloat64(ctx, "longitude", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.ListExperiencesByGeolocationInput{
		Latitude:        latitude,
		Longitude:       longitude,
		Radius:          radius,
		OnlyPublished:   true,
		ExcludeFinished: true,
		ExcludeDeleted:  true,
		CoordinatorID:   util.GetQuery(ctx, "coordinatorId", ""),
		ProducerID:      util.GetQuery(ctx, "producerId", ""),
	}
	experiences, err := h.store.ListExperiencesByGeolocation(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res, err := h.newExperiencesResponse(ctx, experiences, int64(len(experiences)))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newExperiencesResponse(
	ctx context.Context, experiences entity.Experiences, total int64,
) (*response.ExperiencesResponse, error) {
	if len(experiences) == 0 {
		res := &response.ExperiencesResponse{
			Experiences:     []*response.Experience{},
			Coordinators:    []*response.Coordinator{},
			Producers:       []*response.Producer{},
			ExperienceTypes: []*response.ExperienceType{},
		}
		return res, nil
	}

	var (
		coordinators    service.Coordinators
		producers       service.Producers
		experienceTypes service.ExperienceTypes
		experienceRates service.ExperienceRates
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
	eg.Go(func() (err error) {
		experienceRates, err = h.aggregateExperienceRates(ectx, experiences.IDs()...)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	res := &response.ExperiencesResponse{
		Experiences:     service.NewExperiences(experiences, experienceRates.MapByExperienceID()).Response(),
		Coordinators:    coordinators.Response(),
		Producers:       producers.Response(),
		ExperienceTypes: experienceTypes.Response(),
		Total:           total,
	}
	return res, nil
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

func (h *handler) listExperiences(ctx context.Context, in *store.ListExperiencesInput) (service.Experiences, error) {
	experiences, _, err := h.store.ListExperiences(ctx, in)
	if err != nil || len(experiences) == 0 {
		return service.Experiences{}, err
	}
	experiences = experiences.FilterByPublished()
	rates, err := h.aggregateExperienceRates(ctx, experiences.IDs()...)
	if err != nil {
		return nil, err
	}
	return service.NewExperiences(experiences, rates.MapByExperienceID()), nil
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
	experiences = experiences.FilterByPublished()
	rates, err := h.aggregateExperienceRates(ctx, experiences.IDs()...)
	if err != nil {
		return nil, err
	}
	return service.NewExperiences(experiences, rates.MapByExperienceID()), nil
}

func (h *handler) multiGetExperiencesByRevision(ctx context.Context, revisionIDs []int64) (service.Experiences, error) {
	if len(revisionIDs) == 0 {
		return service.Experiences{}, nil
	}
	in := &store.MultiGetExperiencesByRevisionInput{
		ExperienceRevisionIDs: revisionIDs,
	}
	experiences, err := h.store.MultiGetExperiencesByRevision(ctx, in)
	if err != nil {
		return nil, err
	}
	experiences = experiences.FilterByPublished()
	rates, err := h.aggregateExperienceRates(ctx, experiences.IDs()...)
	if err != nil {
		return nil, err
	}
	return service.NewExperiences(experiences, rates.MapByExperienceID()), nil
}

func (h *handler) getExperience(ctx context.Context, experienceID string) (*service.Experience, error) {
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
	rates, err := h.aggregateExperienceRates(ctx, experienceID)
	if err != nil {
		return nil, err
	}
	return service.NewExperience(experience, rates.MapByExperienceID()[experience.ID]), nil
}
