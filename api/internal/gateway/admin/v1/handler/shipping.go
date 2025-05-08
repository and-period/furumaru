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
	res := &response.ShippingsResponse{
		Shippings:    []*response.Shipping{},
		Coordinators: []*response.Coordinator{},
		Total:        0,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetShipping(ctx *gin.Context) {
	var (
		shopID      string
		shipping    *service.Shipping
		coordinator *service.Coordinator
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		shippingIn := &store.GetShippingInput{
			ShippingID: util.GetParam(ctx, "shippingId"),
		}
		sshipping, err := h.store.GetShipping(ectx, shippingIn)
		if err != nil {
			return err
		}
		shopID = sshipping.ShopID
		shipping = service.NewShipping(sshipping)
		return nil
	})
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinator(ctx, util.GetParam(ctx, "coordinatorId"))
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}
	if coordinator.ShopID != shopID {
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
	res := &response.ShippingResponse{
		Shipping:    &response.Shipping{},
		Coordinator: &response.Coordinator{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateShipping(ctx *gin.Context) {
	req := &request.UpdateShippingRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) UpdateActiveShipping(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteShipping(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

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
