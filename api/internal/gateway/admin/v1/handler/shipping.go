package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
)

func (h *handler) shippingRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/shippings", h.authentication)

	r.GET("", h.ListShippings)
	r.POST("", h.CreateShipping)
	r.GET("/:shippingId", h.GetShipping)
	r.PATCH("/:shippingId", h.UpdateShipping)
	r.DELETE("/:shippingId", h.DeleteShipping)
}

func (h *handler) ListShippings(ctx *gin.Context) {
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
	orders, err := h.newShippingOrders(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.ListShippingsInput{
		Name:   util.GetQuery(ctx, "name", ""),
		Limit:  limit,
		Offset: offset,
		Orders: orders,
	}
	if getRole(ctx) == service.AdminRoleCoordinator {
		in.CoordinatorID = getAdminID(ctx)
	}
	shippings, total, err := h.store.ListShippings(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(shippings) == 0 {
		res := &response.ShippingsResponse{
			Shippings: []*response.Shipping{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	coordinators, err := h.multiGetCoordinators(ctx, shippings.CoordinatorIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ShippingsResponse{
		Shippings:    service.NewShippings(shippings).Response(),
		Coordinators: coordinators.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newShippingOrders(ctx *gin.Context) ([]*store.ListShippingsOrder, error) {
	shippings := map[string]sentity.ShippingOrderBy{
		"name":            sentity.ShippingOrderByName,
		"hasFreeShipping": sentity.ShippingOrderByHasFreeShipping,
		"createdAt":       sentity.ShippingOrderByCreatedAt,
		"updatedAt":       sentity.ShippingOrderByUpdatedAt,
	}
	params := util.GetOrders(ctx)
	res := make([]*store.ListShippingsOrder, len(params))
	for i, p := range params {
		key, ok := shippings[p.Key]
		if !ok {
			return nil, fmt.Errorf("handler: unknown order key. key=%s: %w", p.Key, errInvalidOrderkey)
		}
		res[i] = &store.ListShippingsOrder{
			Key:        key,
			OrderByASC: p.Direction == util.OrderByASC,
		}
	}
	return res, nil
}

func (h *handler) GetShipping(ctx *gin.Context) {
	shipping, err := h.getShipping(ctx, util.GetParam(ctx, "shippingId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.ShippingResponse{
		Shipping: shipping.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateShipping(ctx *gin.Context) {
	req := &request.CreateShippingRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	if getRole(ctx) == service.AdminRoleCoordinator {
		if !currentAdmin(ctx, req.CoordinatorID) {
			h.forbidden(ctx, errors.New("handler: invalid coordinator id"))
			return
		}
	}

	in := &store.CreateShippingInput{
		CoordinatorID:      req.CoordinatorID,
		Name:               req.Name,
		IsDefault:          req.IsDefault,
		Box60Rates:         h.newShippingRatesForCreate(req.Box60Rates),
		Box60Refrigerated:  req.Box60Refrigerated,
		Box60Frozen:        req.Box60Frozen,
		Box80Rates:         h.newShippingRatesForCreate(req.Box80Rates),
		Box80Refrigerated:  req.Box80Refrigerated,
		Box80Frozen:        req.Box80Frozen,
		Box100Rates:        h.newShippingRatesForCreate(req.Box100Rates),
		Box100Refrigerated: req.Box100Refrigerated,
		Box100Frozen:       req.Box100Frozen,
		HasFreeShipping:    req.HasFreeShipping,
		FreeShippingRates:  req.FreeShippingRates,
	}
	sshipping, err := h.store.CreateShipping(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	coordinator, err := h.getCoordinator(ctx, sshipping.CoordinatorID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ShippingResponse{
		Shipping:    service.NewShipping(sshipping).Response(),
		Coordinator: coordinator.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newShippingRatesForCreate(in []*request.CreateShippingRate) []*store.CreateShippingRate {
	res := make([]*store.CreateShippingRate, len(in))
	for i := range in {
		res[i] = &store.CreateShippingRate{
			Name:            in[i].Name,
			Price:           in[i].Price,
			PrefectureCodes: in[i].PrefectureCodes,
		}
	}
	return res
}

func (h *handler) UpdateShipping(ctx *gin.Context) {
	req := &request.UpdateShippingRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.UpdateShippingInput{
		ShippingID:         util.GetParam(ctx, "shippingId"),
		Name:               req.Name,
		IsDefault:          req.IsDefault,
		Box60Rates:         h.newShippingRatesForUpdate(req.Box60Rates),
		Box60Refrigerated:  req.Box60Refrigerated,
		Box60Frozen:        req.Box60Frozen,
		Box80Rates:         h.newShippingRatesForUpdate(req.Box80Rates),
		Box80Refrigerated:  req.Box80Refrigerated,
		Box80Frozen:        req.Box80Frozen,
		Box100Rates:        h.newShippingRatesForUpdate(req.Box100Rates),
		Box100Refrigerated: req.Box100Refrigerated,
		Box100Frozen:       req.Box100Frozen,
		HasFreeShipping:    req.HasFreeShipping,
		FreeShippingRates:  req.FreeShippingRates,
	}
	if err := h.store.UpdateShipping(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) newShippingRatesForUpdate(in []*request.UpdateShippingRate) []*store.UpdateShippingRate {
	res := make([]*store.UpdateShippingRate, len(in))
	for i := range in {
		res[i] = &store.UpdateShippingRate{
			Name:            in[i].Name,
			Price:           in[i].Price,
			PrefectureCodes: in[i].PrefectureCodes,
		}
	}
	return res
}

func (h *handler) DeleteShipping(ctx *gin.Context) {
	in := &store.DeleteShippingInput{
		ShippingID: util.GetParam(ctx, "shippingId"),
	}
	if err := h.store.DeleteShipping(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) getShipping(ctx context.Context, shippingID string) (*service.Shipping, error) {
	in := &store.GetShippingInput{
		ShippingID: shippingID,
	}
	shipping, err := h.store.GetShipping(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewShipping(shipping), nil
}
