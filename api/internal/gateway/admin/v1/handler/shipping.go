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
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Shipping
// @tag.description 配送設定関連
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

// @Summary     デフォルト配送設定取得
// @Description デフォルトの配送設定を取得します。
// @Tags        Shipping
// @Router      /v1/shippings/default [get]
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} types.ShippingResponse
func (h *handler) GetDefaultShipping(ctx *gin.Context) {
	in := &store.GetDefaultShippingInput{}
	shipping, err := h.store.GetDefaultShipping(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.ShippingResponse{
		Shipping: service.NewShipping(shipping).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     デフォルト配送設定更新
// @Description デフォルトの配送設定を更新します。
// @Tags        Shipping
// @Router      /v1/shippings/default [patch]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.UpdateDefaultShippingRequest true "デフォルト配送設定情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) UpdateDefaultShipping(ctx *gin.Context) {
	req := &types.UpdateDefaultShippingRequest{}
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

func (h *handler) newShippingRatesForUpdateDefault(in []*types.UpdateDefaultShippingRate) []*store.UpdateDefaultShippingRate {
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

// @Summary     配送設定一覧取得
// @Description 指定されたコーディネーターの配送設定一覧を取得します。ページネーションに対応しています。
// @Tags        Shipping
// @Router      /v1/coordinators/{coordinatorId}/shippings [get]
// @Security    bearerauth
// @Param       coordinatorId path string true "コーディネーターID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Produce     json
// @Success     200 {object} types.ShippingsResponse
// @Failure     403 {object} util.ErrorResponse "アクセス権限がない"
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
		res := &types.ShippingsResponse{
			Shippings:    service.NewShippings(shippings).Response(),
			Coordinators: []*types.Coordinator{coordinator.Response()},
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
	res := &types.ShippingsResponse{
		Shippings: []*types.Shipping{service.NewShipping(shipping).Response()},
		Total:     1,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     配送設定取得
// @Description 指定された配送設定の詳細情報を取得します。
// @Tags        Shipping
// @Router      /v1/coordinators/{coordinatorId}/shippings/{shippingId} [get]
// @Security    bearerauth
// @Param       coordinatorId path string true "コーディネーターID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       shippingId path string true "配送設定ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.ShippingResponse
// @Failure     403 {object} util.ErrorResponse "アクセス権限がない"
// @Failure     404 {object} util.ErrorResponse "配送設定が存在しない"
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
	res := &types.ShippingResponse{
		Shipping:    shipping.Response(),
		Coordinator: coordinator.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     配送設定登録
// @Description 新しい配送設定を登録します。
// @Tags        Shipping
// @Router      /v1/coordinators/{coordinatorId}/shippings [post]
// @Security    bearerauth
// @Param       coordinatorId path string true "コーディネーターID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.CreateShippingRequest true "配送設定情報"
// @Produce     json
// @Success     200 {object} types.ShippingResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "アクセス権限がない"
func (h *handler) CreateShipping(ctx *gin.Context) {
	req := &types.CreateShippingRequest{}
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
	res := &types.ShippingResponse{
		Shipping:    service.NewShipping(shipping).Response(),
		Coordinator: coordinator.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     配送設定更新
// @Description 配送設定の情報を更新します。
// @Tags        Shipping
// @Router      /v1/coordinators/{coordinatorId}/shippings/{shippingId} [patch]
// @Security    bearerauth
// @Param       coordinatorId path string true "コーディネーターID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       shippingId path string true "配送設定ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.UpdateShippingRequest true "配送設定情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "アクセス権限がない"
// @Failure     404 {object} util.ErrorResponse "配送設定が存在しない"
func (h *handler) UpdateShipping(ctx *gin.Context) {
	req := &types.UpdateShippingRequest{}
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

// @Summary     アクティブ配送設定更新
// @Description 指定した配送設定をアクティブに設定します。
// @Tags        Shipping
// @Router      /v1/coordinators/{coordinatorId}/shippings/{shippingId}/activation [patch]
// @Security    bearerauth
// @Param       coordinatorId path string true "コーディネーターID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       shippingId path string true "配送設定ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     403 {object} util.ErrorResponse "アクセス権限がない"
// @Failure     404 {object} util.ErrorResponse "配送設定が存在しない"
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

// @Summary     配送設定削除
// @Description 配送設定を削除します。
// @Tags        Shipping
// @Router      /v1/coordinators/{coordinatorId}/shippings/{shippingId} [delete]
// @Security    bearerauth
// @Param       coordinatorId path string true "コーディネーターID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       shippingId path string true "配送設定ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     403 {object} util.ErrorResponse "アクセス権限がない"
// @Failure     404 {object} util.ErrorResponse "配送設定が存在しない"
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

// @Summary     アクティブ配送設定取得
// @Description Deprecated.指定されたコーディネーターのアクティブ配送設定を取得します。
// @Tags        Shipping
// @Router      /v1/coordinators/{coordinatorId}/shippings/-/activation [get]
// @Security    bearerauth
// @Param       coordinatorId path string true "コーディネーターID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.ShippingResponse
// @Failure     403 {object} util.ErrorResponse "アクセス権限がない"
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
	res := &types.ShippingResponse{
		Shipping: service.NewShipping(shipping).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     配送設定更新
// @Description Deprecated.コーディネータの配送設定を更新します。
// @Tags        Shipping
// @Router      /v1/coordinators/{coordinatorId}/shippings [patch]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.UpsertShippingRequest true "配送設定情報"
// @Param       coordinatorId path string true "コーディネータID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "アクセス権限がない"
func (h *handler) UpsertShipping(ctx *gin.Context) {
	req := &types.UpsertShippingRequest{}
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

func (h *handler) newShippingRatesForCreate(in []*types.CreateShippingRate) []*store.CreateShippingRate {
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

func (h *handler) newShippingRatesForUpdate(in []*types.UpdateShippingRate) []*store.UpdateShippingRate {
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

func (h *handler) newShippingRatesForUpsert(in []*types.UpsertShippingRate) []*store.UpsertShippingRate {
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
