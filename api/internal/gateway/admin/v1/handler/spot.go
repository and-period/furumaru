package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Spot
// @tag.description スポット関連
func (h *handler) spotRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/spots", h.authentication)

	r.GET("", h.ListSpots)
	r.POST("", h.CreateSpot)
	r.GET("/:spotId", h.GetSpot)
	r.PATCH("/:spotId", h.filterAccessSpot, h.UpdateSpot)
	r.DELETE("/:spotId", h.filterAccessSpot, h.DeleteSpot)
	r.PATCH("/:spotId/approval", h.filterAccessSpot, h.ApproveSpot)
}

func (h *handler) filterAccessSpot(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			spot, err := h.getSpot(ctx, util.GetParam(ctx, "spotId"))
			if err != nil {
				return false, err
			}
			if spot.UserType != service.SpotUserTypeCoordinator {
				return false, nil
			}
			return spot.UserID == getAdminID(ctx), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Next()
}

// @Summary     スポット一覧取得
// @Description スポットの一覧を取得します。ページネーションと名前でのフィルタリングに対応しています。
// @Tags        Spot
// @Router      /v1/spots [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Param       name query string false "スポット名" example("春の公園")
// @Produce     json
// @Success     200 {object} types.SpotsResponse
func (h *handler) ListSpots(ctx *gin.Context) {
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

	in := &store.ListSpotsInput{
		Name:            util.GetQuery(ctx, "name", ""),
		ExcludeApproved: false,
		ExcludeDisabled: false,
		Limit:           limit,
		Offset:          offset,
	}
	spots, total, err := h.store.ListSpots(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(spots) == 0 {
		res := &types.SpotsResponse{
			Spots:        []*types.Spot{},
			SpotTypes:    []*types.SpotType{},
			Users:        []*types.User{},
			Coordinators: []*types.Coordinator{},
			Producers:    []*types.Producer{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}
	spotsMap := spots.GroupByUserType()

	var (
		spotTypes    service.SpotTypes
		users        service.Users
		coordinators service.Coordinators
		producers    service.Producers
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		spotTypes, err = h.multiGetSpotTypes(ectx, spots.TypeIDs())
		return
	})
	eg.Go(func() (err error) {
		spots := spotsMap[entity.SpotUserTypeUser]
		users, err = h.multiGetUsers(ectx, spots.UserIDs())
		return
	})
	eg.Go(func() (err error) {
		spots := spotsMap[entity.SpotUserTypeCoordinator]
		coordinators, err = h.multiGetCoordinators(ectx, spots.UserIDs())
		return
	})
	eg.Go(func() (err error) {
		spots := spotsMap[entity.SpotUserTypeProducer]
		producers, err = h.multiGetProducers(ectx, spots.UserIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.SpotsResponse{
		Spots:        service.NewSpots(spots).Response(),
		SpotTypes:    spotTypes.Response(),
		Users:        users.Response(),
		Coordinators: coordinators.Response(),
		Producers:    producers.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     スポット取得
// @Description 指定されたスポットの詳細情報を取得します。
// @Tags        Spot
// @Router      /v1/spots/{spotId} [get]
// @Security    bearerauth
// @Param       spotId path string true "スポットID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.SpotResponse
// @Failure     404 {object} util.ErrorResponse "スポットが存在しない"
func (h *handler) GetSpot(ctx *gin.Context) {
	spot, err := h.getSpot(ctx, util.GetParam(ctx, "spotId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	spotType, err := h.getSpotType(ctx, spot.TypeID)
	if err != nil && !errors.Is(err, exception.ErrNotFound) {
		h.httpError(ctx, err)
		return
	}

	res := &types.SpotResponse{
		Spot:     spot.Response(),
		SpotType: spotType.Response(),
	}

	switch spot.UserType {
	case service.SpotUserTypeUser:
		user, err := h.getUser(ctx, spot.UserID)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		res.User = user.Response()
	case service.SpotUserTypeCoordinator:
		coordinator, err := h.getCoordinator(ctx, spot.UserID)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		res.Coordinator = coordinator.Response()
	case service.SpotUserTypeProducer:
		producer, err := h.getProducer(ctx, spot.UserID)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		res.Producer = producer.Response()
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary     スポット登録
// @Description 新しいスポットを登録します。
// @Tags        Spot
// @Router      /v1/spots [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.CreateSpotRequest true "スポット情報"
// @Produce     json
// @Success     200 {object} types.SpotResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateSpot(ctx *gin.Context) {
	req := &types.CreateSpotRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	spotType, err := h.getSpotType(ctx, req.TypeID)
	if err != nil && !errors.Is(err, exception.ErrNotFound) {
		h.httpError(ctx, err)
		return
	}

	adminID := getAdminID(ctx)

	res := &types.SpotResponse{}
	switch getAdminType(ctx) {
	case service.AdminTypeCoordinator:
		coordinator, err := h.getCoordinator(ctx, adminID)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		res.Coordinator = coordinator.Response()
	case service.AdminTypeProducer:
		producer, err := h.getProducer(ctx, adminID)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		res.Producer = producer.Response()
	}

	in := &store.CreateSpotByAdminInput{
		TypeID:       req.TypeID,
		AdminID:      adminID,
		Name:         req.Name,
		Description:  req.Description,
		ThumbnailURL: req.ThumbnailURL,
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
	}
	spot, err := h.store.CreateSpotByAdmin(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res.Spot = service.NewSpot(spot).Response()
	res.SpotType = spotType.Response()
	ctx.JSON(http.StatusOK, res)
}

// @Summary     スポット更新
// @Description スポットの情報を更新します。
// @Tags        Spot
// @Router      /v1/spots/{spotId} [patch]
// @Security    bearerauth
// @Param       spotId path string true "スポットID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.UpdateSpotRequest true "スポット情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "スポットの更新権限がない"
// @Failure     404 {object} util.ErrorResponse "スポットが存在しない"
func (h *handler) UpdateSpot(ctx *gin.Context) {
	req := &types.UpdateSpotRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.UpdateSpotInput{
		SpotID:       util.GetParam(ctx, "spotId"),
		TypeID:       req.TypeID,
		Name:         req.Name,
		Description:  req.Description,
		ThumbnailURL: req.ThumbnailURL,
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
	}
	if err := h.store.UpdateSpot(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     スポット削除
// @Description スポットを削除します。
// @Tags        Spot
// @Router      /v1/spots/{spotId} [delete]
// @Security    bearerauth
// @Param       spotId path string true "スポットID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     403 {object} util.ErrorResponse "スポットの削除権限がない"
// @Failure     404 {object} util.ErrorResponse "スポットが存在しない"
func (h *handler) DeleteSpot(ctx *gin.Context) {
	in := &store.DeleteSpotInput{
		SpotID: util.GetParam(ctx, "spotId"),
	}
	if err := h.store.DeleteSpot(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     スポット承認
// @Description スポットの承認状態を更新します。
// @Tags        Spot
// @Router      /v1/spots/{spotId}/approval [patch]
// @Security    bearerauth
// @Param       spotId path string true "スポットID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.ApproveSpotRequest true "承認情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "スポットの承認権限がない"
// @Failure     404 {object} util.ErrorResponse "スポットが存在しない"
func (h *handler) ApproveSpot(ctx *gin.Context) {
	req := &types.ApproveSpotRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.ApproveSpotInput{
		SpotID:   util.GetParam(ctx, "spotId"),
		AdminID:  getAdminID(ctx),
		Approved: req.Approved,
	}
	if err := h.store.ApproveSpot(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) getSpot(ctx context.Context, spotID string) (*service.Spot, error) {
	in := &store.GetSpotInput{
		SpotID: spotID,
	}
	spot, err := h.store.GetSpot(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewSpot(spot), nil
}
