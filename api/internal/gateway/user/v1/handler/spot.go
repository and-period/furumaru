package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Spot
// @tag.description スポット関連
func (h *handler) spotRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/spots")

	r.GET("", h.ListSpots)
	r.GET("/:spotId", h.GetSpot)
	r.POST("", h.authentication, h.CreateSpot)
	r.PATCH("/:spotId", h.authentication, h.UpdateSpot)
	r.DELETE("/:spotId", h.authentication, h.DeleteSpot)
}

// @Summary     スポット一覧取得
// @Description 指定された位置情報周辺のスポット一覧を取得します。
// @Tags        Spot
// @Router      /spots [get]
// @Param       latitude query number true "緯度"
// @Param       longitude query number true "経度"
// @Param       radius query int64 false "検索半径（km）" default(20)
// @Produce     json
// @Success     200 {object} response.SpotsResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) ListSpots(ctx *gin.Context) {
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

	in := &store.ListSpotsByGeolocationInput{
		Latitude:        latitude,
		Longitude:       longitude,
		Radius:          radius,
		ExcludeDisabled: true,
	}
	sspots, err := h.store.ListSpotsByGeolocation(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(sspots) == 0 {
		res := &response.SpotsResponse{
			Spots:     []*response.Spot{},
			SpotTypes: []*response.SpotType{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	spotsMap := sspots.GroupByUserType()

	var (
		mu        sync.Mutex
		spotTypes service.SpotTypes
	)
	spots := make(service.Spots, 0, len(sspots))

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		spotTypes, err = h.multiGetSpotTypes(ectx, sspots.TypeIDs())
		return
	})
	eg.Go(func() error {
		sspots := spotsMap[sentity.SpotUserTypeUser]
		users, err := h.multiGetUsers(ectx, sspots.UserIDs())
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		spots = append(spots, service.NewSpotsByUser(sspots, users.Map())...)
		return nil
	})
	eg.Go(func() error {
		sspots := spotsMap[sentity.SpotUserTypeCoordinator]
		coordinators, err := h.multiGetCoordinators(ectx, sspots.UserIDs())
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		spots = append(spots, service.NewSpotsByCoordinator(sspots, coordinators.Map())...)
		return nil
	})
	eg.Go(func() error {
		sspots := spotsMap[sentity.SpotUserTypeProducer]
		producers, err := h.multiGetProducers(ectx, sspots.UserIDs())
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		spots = append(spots, service.NewSpotsByProducer(sspots, producers.Map())...)
		return nil
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.SpotsResponse{
		Spots:     spots.Response(),
		SpotTypes: spotTypes.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     スポット詳細取得
// @Description 指定されたIDのスポット詳細を取得します。
// @Tags        Spot
// @Router      /spots/{spotId} [get]
// @Param       spotId path string true "スポットID"
// @Produce     json
// @Success     200 {object} response.SpotResponse
// @Failure     404 {object} util.ErrorResponse "スポットが見つかりません"
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
	res := &response.SpotResponse{
		Spot:     spot.Response(),
		SpotType: spotType.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     スポット登録
// @Description 新しいスポットを登録します。
// @Tags        Spot
// @Router      /spots [post]
// @Security    bearerauth
// @Accept      json
// @Produce     json
// @Param       body body request.CreateSpotRequest true "スポット情報"
// @Success     200 {object} response.SpotResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) CreateSpot(ctx *gin.Context) {
	req := &request.CreateSpotRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	user, err := h.getMember(ctx, h.getUserID(ctx))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	spotType, err := h.getSpotType(ctx, req.TypeID)
	if err != nil && !errors.Is(err, exception.ErrNotFound) {
		h.httpError(ctx, err)
		return
	}
	in := &store.CreateSpotByUserInput{
		TypeID:       req.TypeID,
		UserID:       user.ID,
		Name:         req.Name,
		Description:  req.Description,
		ThumbnailURL: req.ThumbnailURL,
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
	}
	spot, err := h.store.CreateSpotByUser(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.SpotResponse{
		Spot:     service.NewSpotByUser(spot, user).Response(),
		SpotType: spotType.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     スポット更新
// @Description 指定されたIDのスポット情報を更新します。
// @Tags        Spot
// @Router      /spots/{spotId} [patch]
// @Security    bearerauth
// @Accept      json
// @Param       spotId path string true "スポットID"
// @Param       body body request.UpdateSpotRequest true "更新するスポット情報"
// @Success     204 "更新成功"
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     403 {object} util.ErrorResponse "アクセス権限がありません"
// @Failure     404 {object} util.ErrorResponse "スポットが見つかりません"
func (h *handler) UpdateSpot(ctx *gin.Context) {
	req := &request.UpdateSpotRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	spot, err := h.getSpot(ctx, util.GetParam(ctx, "spotId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if spot.UserID != h.getUserID(ctx) {
		h.forbidden(ctx, errors.New("handler: user is not owner"))
		return
	}
	in := &store.UpdateSpotInput{
		SpotID:       spot.ID,
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
// @Description 指定されたIDのスポットを削除します。
// @Tags        Spot
// @Router      /spots/{spotId} [delete]
// @Security    bearerauth
// @Param       spotId path string true "スポットID"
// @Success     204 "削除成功"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     403 {object} util.ErrorResponse "アクセス権限がありません"
// @Failure     404 {object} util.ErrorResponse "スポットが見つかりません"
func (h *handler) DeleteSpot(ctx *gin.Context) {
	spot, err := h.getSpot(ctx, util.GetParam(ctx, "spotId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if spot.UserID != h.getUserID(ctx) {
		h.forbidden(ctx, errors.New("handler: user is not owner"))
		return
	}
	in := &store.DeleteSpotInput{
		SpotID: spot.ID,
	}
	if err := h.store.DeleteSpot(ctx, in); err != nil {
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
	if !spot.Approved {
		return nil, fmt.Errorf("spot is not approved: %w", exception.ErrNotFound)
	}

	switch spot.UserType {
	case sentity.SpotUserTypeUser:
		user, err := h.getMember(ctx, spot.UserID)
		if err != nil {
			return nil, err
		}
		return service.NewSpotByUser(spot, user), nil
	case sentity.SpotUserTypeCoordinator:
		coordinator, err := h.getCoordinator(ctx, spot.UserID)
		if err != nil {
			return nil, err
		}
		return service.NewSpotByCoordinator(spot, coordinator), nil
	case sentity.SpotUserTypeProducer:
		producer, err := h.getProducer(ctx, spot.UserID)
		if err != nil {
			return nil, err
		}
		return service.NewSpotByProducer(spot, producer), nil
	default:
		return nil, fmt.Errorf("unknown user type: %w", exception.ErrNotFound)
	}
}
