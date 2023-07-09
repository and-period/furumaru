package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
)

func (h *handler) shippingRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("", h.ListShippings)
	arg.POST("", h.CreateShipping)
	arg.GET("/:shippingId", h.GetShipping)
	arg.PATCH("/:shippingId", h.UpdateShipping)
	arg.DELETE("/:shippingId", h.DeleteShipping)
}

func (h *handler) ListShippings(ctx *gin.Context) {
	const (
		defaultLimit  = 20
		defaultOffset = 0
	)

	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	orders, err := h.newShippingOrders(ctx)
	if err != nil {
		badRequest(ctx, err)
		return
	}

	in := &store.ListShippingsInput{
		Name:   util.GetQuery(ctx, "name", ""),
		Limit:  limit,
		Offset: offset,
		Orders: orders,
	}
	sshippings, total, err := h.store.ListShippings(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	shippings, err := service.NewShippings(sshippings)
	if err != nil {
		httpError(ctx, err)
		return
	}
	res := &response.ShippingsResponse{
		Shippings: shippings.Response(),
		Total:     total,
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
		httpError(ctx, err)
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
		badRequest(ctx, err)
		return
	}
	box60Rates, err := h.newShippingRatesForCreate(req.Box60Rates)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	box80Rates, err := h.newShippingRatesForCreate(req.Box80Rates)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	box100Rates, err := h.newShippingRatesForCreate(req.Box100Rates)
	if err != nil {
		badRequest(ctx, err)
		return
	}

	in := &store.CreateShippingInput{
		Name:               req.Name,
		Box60Rates:         box60Rates,
		Box60Refrigerated:  req.Box60Refrigerated,
		Box60Frozen:        req.Box60Frozen,
		Box80Rates:         box80Rates,
		Box80Refrigerated:  req.Box80Refrigerated,
		Box80Frozen:        req.Box80Frozen,
		Box100Rates:        box100Rates,
		Box100Refrigerated: req.Box100Refrigerated,
		Box100Frozen:       req.Box100Frozen,
		HasFreeShipping:    req.HasFreeShipping,
		FreeShippingRates:  req.FreeShippingRates,
	}
	sshipping, err := h.store.CreateShipping(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	shipping, err := service.NewShipping(sshipping)
	if err != nil {
		httpError(ctx, err)
		return
	}
	res := &response.ShippingResponse{
		Shipping: shipping.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newShippingRatesForCreate(in []*request.CreateShippingRate) ([]*store.CreateShippingRate, error) {
	res := make([]*store.CreateShippingRate, len(in))
	for i := range in {
		prefectures, err := codes.ToPrefectureValues(in[i].Prefectures...)
		if err != nil {
			return nil, err
		}
		res[i] = &store.CreateShippingRate{
			Name:        in[i].Name,
			Price:       in[i].Price,
			Prefectures: prefectures,
		}
	}
	return res, nil
}

func (h *handler) UpdateShipping(ctx *gin.Context) {
	req := &request.UpdateShippingRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}
	box60Rates, err := h.newShippingRatesForUpdate(req.Box60Rates)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	box80Rates, err := h.newShippingRatesForUpdate(req.Box80Rates)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	box100Rates, err := h.newShippingRatesForUpdate(req.Box100Rates)
	if err != nil {
		badRequest(ctx, err)
		return
	}

	in := &store.UpdateShippingInput{
		ShippingID:         util.GetParam(ctx, "shippingId"),
		Name:               req.Name,
		Box60Rates:         box60Rates,
		Box60Refrigerated:  req.Box60Refrigerated,
		Box60Frozen:        req.Box60Frozen,
		Box80Rates:         box80Rates,
		Box80Refrigerated:  req.Box80Refrigerated,
		Box80Frozen:        req.Box80Frozen,
		Box100Rates:        box100Rates,
		Box100Refrigerated: req.Box100Refrigerated,
		Box100Frozen:       req.Box100Frozen,
		HasFreeShipping:    req.HasFreeShipping,
		FreeShippingRates:  req.FreeShippingRates,
	}
	if err := h.store.UpdateShipping(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) newShippingRatesForUpdate(in []*request.UpdateShippingRate) ([]*store.UpdateShippingRate, error) {
	res := make([]*store.UpdateShippingRate, len(in))
	for i := range in {
		prefectures, err := codes.ToPrefectureValues(in[i].Prefectures...)
		if err != nil {
			return nil, err
		}
		res[i] = &store.UpdateShippingRate{
			Name:        in[i].Name,
			Price:       in[i].Price,
			Prefectures: prefectures,
		}
	}
	return res, nil
}

func (h *handler) DeleteShipping(ctx *gin.Context) {
	in := &store.DeleteShippingInput{
		ShippingID: util.GetParam(ctx, "shippingId"),
	}
	if err := h.store.DeleteShipping(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) multiGetShippings(ctx context.Context, shippingIDs []string) (service.Shippings, error) {
	if len(shippingIDs) == 0 {
		return service.Shippings{}, nil
	}
	in := &store.MultiGetShippingsInput{
		ShippingIDs: shippingIDs,
	}
	shippings, err := h.store.MultiGetShippings(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewShippings(shippings)
}

func (h *handler) getShipping(ctx context.Context, shippingID string) (*service.Shipping, error) {
	in := &store.GetShippingInput{
		ShippingID: shippingID,
	}
	shipping, err := h.store.GetShipping(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewShipping(shipping)
}
