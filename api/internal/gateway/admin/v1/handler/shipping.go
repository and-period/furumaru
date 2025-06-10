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
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) shippingRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/shippings", h.authentication)
	r.GET("/default", h.GetDefaultShipping)
	r.PATCH("/default", h.UpdateDefaultShipping)

	cr := rg.Group("/coordinators/:coordinatorId/shippings", h.authentication, h.filterAccessShipping)
	cr.GET("", h.ListShippings)
	cr.POST("", h.CreateShipping)
	cr.GET("/:shippingId", h.GetShipping)
	cr.PATCH("/:shippingId", h.UpdateShipping)
	cr.DELETE("/:shippingId", h.DeleteShipping)
	cr.PATCH("/:shippingId/activation", h.UpdateActiveShipping)
	cr.GET("/-/activation", h.GetActiveShipping) // Deprecated
	cr.PATCH("", h.UpsertShipping)               // Deprecated
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

	coordinator, err := h.getCoordinator(ctx, util.GetParam(ctx, "coordinatorId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	in := &store.ListShippingsByShopIDInput{
		ShopID: coordinator.ShopID,
		Limit:  limit,
		Offset: offset,
	}
	shippings, total, err := h.store.ListShippingsByShopID(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	if len(shippings) > 0 {
		res := &response.ShippingsResponse{
			Shippings:    service.NewShippings(shippings).Response(),
			Coordinators: []*response.Coordinator{coordinator.Response()},
			Total:        total,
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	// 1件も取得できない場合、デフォルトの配送設定を返す
	shipping, err := h.store.GetDefaultShipping(ctx, &store.GetDefaultShippingInput{})
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.ShippingsResponse{
		Shippings: []*response.Shipping{service.NewShipping(shipping).Response()},
		Total:     1,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetShipping(ctx *gin.Context) {
	var (
		shipping    *service.Shipping
		coordinator *service.Coordinator
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		shipping, err = h.getShipping(ectx, util.GetParam(ctx, "shippingId"))
		return
	})
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinator(ectx, util.GetParam(ctx, "coordinatorId"))
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}
	if shipping.ShopID != coordinator.ShopID {
		h.notFound(ctx, errors.New("handler: not found"))
		return
	}
	res := &response.ShippingResponse{
		Shipping:    shipping.Response(),
		Coordinator: coordinator.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateShipping(ctx *gin.Context) {
	req := &request.CreateShippingRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	coordinator, err := h.getCoordinator(ctx, util.GetParam(ctx, "coordinatorId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	in := &store.CreateShippingInput{
		Name:              req.Name,
		ShopID:            coordinator.ShopID,
		CoordinatorID:     coordinator.ID,
		Box60Rates:        h.newShippingRatesForCreate(req.Box60Rates),
		Box60Frozen:       req.Box60Frozen,
		Box80Rates:        h.newShippingRatesForCreate(req.Box80Rates),
		Box80Frozen:       req.Box80Frozen,
		Box100Rates:       h.newShippingRatesForCreate(req.Box100Rates),
		Box100Frozen:      req.Box100Frozen,
		HasFreeShipping:   req.HasFreeShipping,
		FreeShippingRates: req.FreeShippingRates,
	}
	shipping, err := h.store.CreateShipping(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.ShippingResponse{
		Shipping:    service.NewShipping(shipping).Response(),
		Coordinator: coordinator.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateShipping(ctx *gin.Context) {
	req := &request.UpdateShippingRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	coordinator, err := h.getCoordinator(ctx, util.GetParam(ctx, "coordinatorId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	shipping, err := h.getShipping(ctx, util.GetParam(ctx, "shippingId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if shipping.ShopID != coordinator.ShopID {
		h.notFound(ctx, errors.New("handler: not found"))
		return
	}
	in := &store.UpdateShippingInput{
		Name:              req.Name,
		ShippingID:        shipping.ID,
		Box60Rates:        h.newShippingRatesForUpdate(req.Box60Rates),
		Box60Frozen:       req.Box60Frozen,
		Box80Rates:        h.newShippingRatesForUpdate(req.Box80Rates),
		Box80Frozen:       req.Box80Frozen,
		Box100Rates:       h.newShippingRatesForUpdate(req.Box100Rates),
		Box100Frozen:      req.Box100Frozen,
		HasFreeShipping:   req.HasFreeShipping,
		FreeShippingRates: req.FreeShippingRates,
	}
	if err := h.store.UpdateShipping(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) UpdateActiveShipping(ctx *gin.Context) {
	coordinator, err := h.getCoordinator(ctx, util.GetParam(ctx, "coordinatorId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	shipping, err := h.getShipping(ctx, util.GetParam(ctx, "shippingId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if shipping.ShopID != coordinator.ShopID {
		h.notFound(ctx, errors.New("handler: not found"))
		return
	}
	in := &store.UpdateShippingInUseInput{
		ShopID:     coordinator.ShopID,
		ShippingID: shipping.ID,
	}
	if err := h.store.UpdateShippingInUse(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteShipping(ctx *gin.Context) {
	coordinator, err := h.getCoordinator(ctx, util.GetParam(ctx, "coordinatorId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	shipping, err := h.getShipping(ctx, util.GetParam(ctx, "shippingId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if shipping.ShopID != coordinator.ShopID {
		h.notFound(ctx, errors.New("handler: not found"))
		return
	}
	in := &store.DeleteShippingInput{
		ShippingID: shipping.ID,
	}
	if err := h.store.DeleteShipping(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// Deprecated
func (h *handler) GetActiveShipping(ctx *gin.Context) {
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
	if getAdminType(ctx).IsCoordinator() && getAdminID(ctx) != coordinatorID {
		h.forbidden(ctx, errors.New("handler: not authorized this coordinator"))
		return
	}
	shop, err := h.getShopByCoordinatorID(ctx, coordinatorID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	in := &store.UpsertShippingInput{
		ShopID:            shop.ID,
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
