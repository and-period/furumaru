package handler

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/pkg/jst"
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
			experience, err := h.getExperience(ctx, util.GetParam(ctx, "experienceId"))
			if err != nil {
				return false, err
			}
			return currentAdmin(ctx, experience.CoordinatorID), nil
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
		ShopID:         getShopID(ctx),
		Name:           util.GetQuery(ctx, "name", ""),
		ProducerID:     util.GetQuery(ctx, "producerId", ""),
		ExcludeDeleted: true,
		Limit:          limit,
		Offset:         offset,
		NoLimit:        false,
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
	if getAdminType(ctx).IsCoordinator() {
		if req.CoordinatorID != getAdminID(ctx) {
			h.forbidden(ctx, errors.New("handler: not allowed to create experience"))
			return
		}
		shop, err := h.getShop(ctx, getShopID(ctx))
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		if !slices.Contains(shop.ProducerIDs, req.ProducerID) {
			h.forbidden(ctx, errors.New("handler: not allowed to create experience"))
			return
		}
	}

	var (
		shop           *service.Shop
		coordinator    *service.Coordinator
		producer       *service.Producer
		experienceType *service.ExperienceType
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		shop, err = h.getShopByCoordinatorID(ectx, req.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinator(ectx, req.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		producer, err = h.getProducer(ectx, req.ProducerID)
		return
	})
	eg.Go(func() (err error) {
		experienceType, err = h.getExperienceType(ectx, req.TypeID)
		return
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		h.badRequest(ctx, err)
		return
	}
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	media := make([]*store.CreateExperienceMedia, len(req.Media))
	for i := range req.Media {
		media[i] = &store.CreateExperienceMedia{
			URL:         req.Media[i].URL,
			IsThumbnail: req.Media[i].IsThumbnail,
		}
	}
	in := &store.CreateExperienceInput{
		ShopID:                shop.ID,
		CoordinatorID:         req.CoordinatorID,
		ProducerID:            req.ProducerID,
		TypeID:                req.TypeID,
		Title:                 req.Title,
		Description:           req.Description,
		Public:                req.Public,
		SoldOut:               req.SoldOut,
		Media:                 media,
		PriceAdult:            req.PriceAdult,
		PriceJuniorHighSchool: req.PriceJuniorHighSchool,
		PriceElementarySchool: req.PriceElementarySchool,
		PricePreschool:        req.PricePreschool,
		PriceSenior:           req.PriceSenior,
		RecommendedPoints:     h.newExperiencePoints(req.RecommendedPoint1, req.RecommendedPoint2, req.RecommendedPoint3),
		PromotionVideoURL:     req.PromotionVideoURL,
		Duration:              req.Duration,
		Direction:             req.Direction,
		BusinessOpenTime:      req.BusinessOpenTime,
		BusinessCloseTime:     req.BusinessCloseTime,
		HostPostalCode:        req.HostPostalCode,
		HostPrefectureCode:    req.HostPrefectureCode,
		HostCity:              req.HostCity,
		HostAddressLine1:      req.HostAddressLine1,
		HostAddressLine2:      req.HostAddressLine2,
		StartAt:               jst.ParseFromUnix(req.StartAt),
		EndAt:                 jst.ParseFromUnix(req.EndAt),
	}
	experience, err := h.store.CreateExperience(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ExperienceResponse{
		Experience:     service.NewExperience(experience).Response(),
		Coordinator:    coordinator.Response(),
		Producer:       producer.Response(),
		ExperienceType: experienceType.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateExperience(ctx *gin.Context) {
	req := &request.UpdateExperienceRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	_, err := h.getExperienceType(ctx, req.TypeID)
	if errors.Is(err, exception.ErrNotFound) {
		h.badRequest(ctx, err)
		return
	}
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	media := make([]*store.UpdateExperienceMedia, len(req.Media))
	for i := range req.Media {
		media[i] = &store.UpdateExperienceMedia{
			URL:         req.Media[i].URL,
			IsThumbnail: req.Media[i].IsThumbnail,
		}
	}
	in := &store.UpdateExperienceInput{
		ExperienceID:          util.GetParam(ctx, "experienceId"),
		TypeID:                req.TypeID,
		Title:                 req.Title,
		Description:           req.Description,
		Public:                req.Public,
		SoldOut:               req.SoldOut,
		Media:                 media,
		PriceAdult:            req.PriceAdult,
		PriceJuniorHighSchool: req.PriceJuniorHighSchool,
		PriceElementarySchool: req.PriceElementarySchool,
		PricePreschool:        req.PricePreschool,
		PriceSenior:           req.PriceSenior,
		RecommendedPoints:     h.newExperiencePoints(req.RecommendedPoint1, req.RecommendedPoint2, req.RecommendedPoint3),
		PromotionVideoURL:     req.PromotionVideoURL,
		Duration:              req.Duration,
		Direction:             req.Direction,
		BusinessOpenTime:      req.BusinessOpenTime,
		BusinessCloseTime:     req.BusinessCloseTime,
		HostPostalCode:        req.HostPostalCode,
		HostPrefectureCode:    req.HostPrefectureCode,
		HostCity:              req.HostCity,
		HostAddressLine1:      req.HostAddressLine1,
		HostAddressLine2:      req.HostAddressLine2,
		StartAt:               jst.ParseFromUnix(req.StartAt),
		EndAt:                 jst.ParseFromUnix(req.EndAt),
	}
	if err := h.store.UpdateExperience(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) newExperiencePoints(points ...string) []string {
	res := make([]string, 0, len(points))
	for _, point := range points {
		if point == "" {
			continue
		}
		res = append(res, point)
	}
	return res
}

func (h *handler) DeleteExperience(ctx *gin.Context) {
	in := &store.DeleteExperienceInput{
		ExperienceID: util.GetParam(ctx, "experienceId"),
	}
	if err := h.store.DeleteExperience(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
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
