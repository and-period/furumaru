package handler

import (
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
)

func (h *handler) shippingRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/shippings", h.authentication)
	r.GET("/default", h.GetDefaultShipping)
	r.PATCH("/default", h.UpdateDefaultShipping)

	cr := rg.Group("/coordinators/:coordinatorId/shippings", h.authentication, h.filterAccessShipping)
	cr.GET("", h.GetShipping)
	cr.PATCH("", h.UpsertShipping)
}

func (h *handler) filterAccessShipping(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			coordinatorID := util.GetParam(ctx, "coordinatorId")
			return currentAdmin(ctx, coordinatorID), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Next()
}

func (h *handler) GetDefaultShipping(ctx *gin.Context) {
	in := &store.GetDefaultShippingInput{}
	shipping, err := h.store.GetDefaultShipping(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.ShippingResponse{
		Shipping: service.NewShipping(shipping).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateDefaultShipping(ctx *gin.Context) {
	req := &request.UpdateDefaultShippingRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.UpdateDefaultShippingInput{
		Box60Rates:        h.newShippingRatesForUpdateDefault(req.Box60Rates),
		Box60Frozen:       req.Box60Frozen,
		Box80Rates:        h.newShippingRatesForUpdateDefault(req.Box80Rates),
		Box80Frozen:       req.Box80Frozen,
		Box100Rates:       h.newShippingRatesForUpdateDefault(req.Box100Rates),
		Box100Frozen:      req.Box100Frozen,
		HasFreeShipping:   req.HasFreeShipping,
		FreeShippingRates: req.FreeShippingRates,
	}
	if err := h.store.UpdateDefaultShipping(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) newShippingRatesForUpdateDefault(in []*request.UpdateDefaultShippingRate) []*store.UpdateDefaultShippingRate {
	res := make([]*store.UpdateDefaultShippingRate, len(in))
	for i := range in {
		res[i] = &store.UpdateDefaultShippingRate{
			Name:            in[i].Name,
			Price:           in[i].Price,
			PrefectureCodes: in[i].PrefectureCodes,
		}
	}
	return res
}

func (h *handler) GetShipping(ctx *gin.Context) {
	in := &store.GetShippingByCoordinatorIDInput{
		CoordinatorID: util.GetParam(ctx, "coordinatorId"),
	}
	shipping, err := h.store.GetShippingByCoordinatorID(ctx, in)
	if errors.Is(err, exception.ErrNotFound) {
		// 配送設定の登録をしていない場合、デフォルト設定を返却する
		in := &store.GetDefaultShippingInput{}
		shipping, err = h.store.GetDefaultShipping(ctx, in)
	}
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.ShippingResponse{
		Shipping: service.NewShipping(shipping).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpsertShipping(ctx *gin.Context) {
	req := &request.UpsertShippingRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	coordinatorID := util.GetParam(ctx, "coordinatorId")
	if getRole(ctx).IsCoordinator() && getAdminID(ctx) != coordinatorID {
		h.forbidden(ctx, errors.New("handler: not authorized this coordinator"))
		return
	}
	in := &store.UpsertShippingInput{
		CoordinatorID:     coordinatorID,
		Box60Rates:        h.newShippingRatesForUpsert(req.Box60Rates),
		Box60Frozen:       req.Box60Frozen,
		Box80Rates:        h.newShippingRatesForUpsert(req.Box80Rates),
		Box80Frozen:       req.Box80Frozen,
		Box100Rates:       h.newShippingRatesForUpsert(req.Box100Rates),
		Box100Frozen:      req.Box100Frozen,
		HasFreeShipping:   req.HasFreeShipping,
		FreeShippingRates: req.FreeShippingRates,
	}
	if err := h.store.UpsertShipping(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) newShippingRatesForUpsert(in []*request.UpsertShippingRate) []*store.UpsertShippingRate {
	res := make([]*store.UpsertShippingRate, len(in))
	for i := range in {
		res[i] = &store.UpsertShippingRate{
			Name:            in[i].Name,
			Price:           in[i].Price,
			PrefectureCodes: in[i].PrefectureCodes,
		}
	}
	return res
}
