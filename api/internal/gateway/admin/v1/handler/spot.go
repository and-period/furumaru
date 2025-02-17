package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

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
		res := &response.SpotsResponse{
			Spots:        []*response.Spot{},
			SpotTypes:    []*response.SpotType{},
			Users:        []*response.User{},
			Coordinators: []*response.Coordinator{},
			Producers:    []*response.Producer{},
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

	res := &response.SpotsResponse{
		Spots:        service.NewSpots(spots).Response(),
		SpotTypes:    spotTypes.Response(),
		Users:        users.Response(),
		Coordinators: coordinators.Response(),
		Producers:    producers.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

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

func (h *handler) CreateSpot(ctx *gin.Context) {
	req := &request.CreateSpotRequest{}
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

	res := &response.SpotResponse{}
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

func (h *handler) UpdateSpot(ctx *gin.Context) {
	req := &request.UpdateSpotRequest{}
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

func (h *handler) ApproveSpot(ctx *gin.Context) {
	req := &request.ApproveSpotRequest{}
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
