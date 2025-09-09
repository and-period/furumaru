package handler

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Experience
// @tag.description 体験関連
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

// @Summary     体験一覧取得
// @Description 体験の一覧を取得します。店舗、生産者、名前でのフィルタリングが可能です。
// @Tags        Experience
// @Router      /v1/experiences [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Param       name query string false "体験名(あいまい検索)" example("農業体験")
// @Param       producerId query string false "生産者ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.ExperiencesResponse
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
		res := &types.ExperiencesResponse{
			Experiences:     []*types.Experience{},
			Coordinators:    []*types.Coordinator{},
			Producers:       []*types.Producer{},
			ExperienceTypes: []*types.ExperienceType{},
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

	res := &types.ExperiencesResponse{
		Experiences:     service.NewExperiences(experiences).Response(),
		Coordinators:    coordinators.Response(),
		Producers:       producers.Response(),
		ExperienceTypes: experienceTypes.Response(),
		Total:           total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     体験取得
// @Description 指定された体験の詳細情報を取得します。
// @Tags        Experience
// @Router      /v1/experiences/{experienceId} [get]
// @Security    bearerauth
// @Param       experienceId path string true "体験ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.ExperienceResponse
// @Failure     403 {object} util.ErrorResponse "体験の参照権限がない"
// @Failure     404 {object} util.ErrorResponse "体験が存在しない"
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

	res := &types.ExperienceResponse{
		Experience:     experience.Response(),
		Coordinator:    coordinator.Response(),
		Producer:       producer.Response(),
		ExperienceType: experienceType.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     体験登録
// @Description 新しい体験を登録します。コーディネーターは管理店舗の生産者の体験のみ登録可能です。
// @Tags        Experience
// @Router      /v1/experiences [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.CreateExperienceRequest true "体験情報"
// @Produce     json
// @Success     200 {object} types.ExperienceResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "体験の登録権限がない"
func (h *handler) CreateExperience(ctx *gin.Context) {
	req := &types.CreateExperienceRequest{}
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

	res := &types.ExperienceResponse{
		Experience:     service.NewExperience(experience).Response(),
		Coordinator:    coordinator.Response(),
		Producer:       producer.Response(),
		ExperienceType: experienceType.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     体験更新
// @Description 体験の情報を更新します。
// @Tags        Experience
// @Router      /v1/experiences/{experienceId} [patch]
// @Security    bearerauth
// @Param       experienceId path string true "体験ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.UpdateExperienceRequest true "体験情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "体験の更新権限がない"
// @Failure     404 {object} util.ErrorResponse "体験が存在しない"
func (h *handler) UpdateExperience(ctx *gin.Context) {
	req := &types.UpdateExperienceRequest{}
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

// @Summary     体験削除
// @Description 体験を削除します。
// @Tags        Experience
// @Router      /v1/experiences/{experienceId} [delete]
// @Security    bearerauth
// @Param       experienceId path string true "体験ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     403 {object} util.ErrorResponse "体験の削除権限がない"
// @Failure     404 {object} util.ErrorResponse "体験が存在しない"
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
