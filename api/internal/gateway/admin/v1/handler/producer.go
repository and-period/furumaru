package handler

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/gin-gonic/gin"
)

func (h *handler) producerRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/producers", h.authentication)

	r.GET("", h.ListProducers)
	r.POST("", h.CreateProducer)
	r.GET("/:producerId", h.filterAccessProducer, h.GetProducer)
	r.PATCH("/:producerId", h.filterAccessProducer, h.UpdateProducer)
	r.DELETE("/:producerId", h.filterAccessProducer, h.DeleteProducer)
}

func (h *handler) filterAccessProducer(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			shop, err := h.getShop(ctx, getShopID(ctx))
			if err != nil {
				return false, err
			}
			return slices.Contains(shop.ProducerIDs, util.GetParam(ctx, "producerId")), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Next()
}

func (h *handler) ListProducers(ctx *gin.Context) {
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

	var (
		producers service.Producers
		total     int64
	)
	if getAdminType(ctx) == service.AdminTypeCoordinator {
		in := &store.ListShopProducersInput{
			ShopID: getShopID(ctx),
			Limit:  limit,
			Offset: offset,
		}
		producerIDs, err := h.store.ListShopProducers(ctx, in)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		producers, err = h.multiGetProducers(ctx, producerIDs)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		total = int64(len(producerIDs))
	} else {
		in := &user.ListProducersInput{
			Name:   util.GetQuery(ctx, "username", ""),
			Limit:  limit,
			Offset: offset,
		}
		var ps entity.Producers
		ps, total, err = h.user.ListProducers(ctx, in)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		producers = service.NewProducers(ps)
	}

	shops, err := h.listShopsByProducerIDs(ctx, producers.IDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	coordinators, err := h.multiGetCoordinators(ctx, shops.CoordinatorIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ProducersResponse{
		Producers:    producers.Response(),
		Shops:        shops.Response(),
		Coordinators: coordinators.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetProducer(ctx *gin.Context) {
	producer, err := h.getProducer(ctx, util.GetParam(ctx, "producerId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	shops, err := h.listShopsByProducerIDs(ctx, []string{producer.ID})
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ProducerResponse{
		Producer: producer.Response(),
		Shops:    shops.Response(),
	}
	if len(shops) == 0 {
		ctx.JSON(http.StatusOK, res)
		return
	}

	coordinators, err := h.multiGetCoordinators(ctx, shops.CoordinatorIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res.Coordinators = coordinators.Response()

	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateProducer(ctx *gin.Context) {
	req := &request.CreateProducerRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	if getAdminType(ctx) == service.AdminTypeCoordinator {
		if !currentAdmin(ctx, req.CoordinatorID) {
			h.forbidden(ctx, errors.New("handler: invalid coordinator id"))
			return
		}
	}

	in := &user.CreateProducerInput{
		CoordinatorID:     req.CoordinatorID,
		Lastname:          req.Lastname,
		Firstname:         req.Firstname,
		LastnameKana:      req.LastnameKana,
		FirstnameKana:     req.FirstnameKana,
		Username:          req.Username,
		Profile:           req.Profile,
		ThumbnailURL:      req.ThumbnailURL,
		HeaderURL:         req.HeaderURL,
		PromotionVideoURL: req.PromotionVideoURL,
		BonusVideoURL:     req.BonusVideoURL,
		InstagramID:       req.InstagramID,
		FacebookID:        req.FacebookID,
		Email:             req.Email,
		PhoneNumber:       req.PhoneNumber,
		PostalCode:        req.PostalCode,
		PrefectureCode:    req.PrefectureCode,
		City:              req.City,
		AddressLine1:      req.AddressLine1,
		AddressLine2:      req.AddressLine2,
	}
	producer, err := h.user.CreateProducer(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	coordinator, err := h.getCoordinator(ctx, req.CoordinatorID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	shop, err := h.getShopByCoordinatorID(ctx, coordinator.ID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ProducerResponse{
		Producer:     service.NewProducer(producer).Response(),
		Shops:        []*response.Shop{shop.Response()},
		Coordinators: []*response.Coordinator{coordinator.Response()},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateProducer(ctx *gin.Context) {
	req := &request.UpdateProducerRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.UpdateProducerInput{
		ProducerID:        util.GetParam(ctx, "producerId"),
		Lastname:          req.Lastname,
		Firstname:         req.Firstname,
		LastnameKana:      req.LastnameKana,
		FirstnameKana:     req.FirstnameKana,
		Username:          req.Username,
		Profile:           req.Profile,
		ThumbnailURL:      req.ThumbnailURL,
		HeaderURL:         req.HeaderURL,
		PromotionVideoURL: req.PromotionVideoURL,
		BonusVideoURL:     req.BonusVideoURL,
		InstagramID:       req.InstagramID,
		FacebookID:        req.FacebookID,
		Email:             req.Email,
		PhoneNumber:       req.PhoneNumber,
		PostalCode:        req.PostalCode,
		PrefectureCode:    req.PrefectureCode,
		City:              req.City,
		AddressLine1:      req.AddressLine1,
		AddressLine2:      req.AddressLine2,
	}
	if err := h.user.UpdateProducer(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteProducer(ctx *gin.Context) {
	in := &user.DeleteProducerInput{
		ProducerID: util.GetParam(ctx, "producerId"),
	}
	if err := h.user.DeleteProducer(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) multiGetProducers(ctx context.Context, producerIDs []string) (service.Producers, error) {
	if len(producerIDs) == 0 {
		return service.Producers{}, nil
	}
	in := &user.MultiGetProducersInput{
		ProducerIDs: producerIDs,
	}
	producers, err := h.user.MultiGetProducers(ctx, in)
	if err != nil {
		return nil, err
	}
	if len(producers) == 0 {
		return service.Producers{}, nil
	}
	return service.NewProducers(producers), nil
}

func (h *handler) getProducer(ctx context.Context, producerID string) (*service.Producer, error) {
	in := &user.GetProducerInput{
		ProducerID: producerID,
	}
	producer, err := h.user.GetProducer(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProducer(producer), nil
}
