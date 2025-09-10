package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Order
// @tag.description 注文関連
func (h *handler) orderRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/orders", h.authentication)

	r.GET("", h.ListOrders)
	r.POST("/-/export", h.ExportOrders)
	r.GET("/:orderId", h.filterAccessOrder, h.GetOrder)
	r.POST("/:orderId/draft", h.filterAccessOrder, h.DraftOrder)
	r.POST("/:orderId/capture", h.filterAccessOrder, h.CaptureOrder)
	r.POST("/:orderId/complete", h.filterAccessOrder, h.CompleteOrder)
	r.POST("/:orderId/cancel", h.filterAccessOrder, h.CancelOrder)
	r.POST("/:orderId/refund", h.filterAccessOrder, h.RefundOrder)
	r.PATCH("/:orderId/fulfillments/:fulfillmentId", h.filterAccessOrder, h.UpdateOrderFulfillment)
}

func (h *handler) filterAccessOrder(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			order, err := h.getOrder(ctx, util.GetParam(ctx, "orderId"))
			if err != nil {
				return false, err
			}
			return currentAdmin(ctx, order.CoordinatorID), nil
		},
		producer: func(_ *gin.Context) (bool, error) {
			// TODO: フィルタリング実装までは全て拒否
			return false, nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Next()
}

// @Summary     注文一覧取得
// @Description 注文の一覧を取得します。コーディネータは自分の店舗の注文のみ取得できます。
// @Tags        Order
// @Router      /v1/orders [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Param       statuses query []int32 false "注文ステータスフィルタ" collectionFormat(csv)
// @Param       types query []int32 false "注文タイプフィルタ" collectionFormat(csv)
// @Produce     json
// @Success     200 {object} types.OrdersResponse
func (h *handler) ListOrders(ctx *gin.Context) {
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
	statuses, otypes, err := h.newOrderFileters(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.ListOrdersInput{
		ShopID:   getShopID(ctx),
		Limit:    limit,
		Offset:   offset,
		Statuses: statuses,
		Types:    otypes,
	}
	orders, total, err := h.store.ListOrders(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(orders) == 0 {
		res := &types.OrdersResponse{
			Orders:       []*types.Order{},
			Users:        []*types.User{},
			Coordinators: []*types.Coordinator{},
			Promotions:   []*types.Promotion{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	var (
		users        service.Users
		coordinators service.Coordinators
		addresses    service.Addresses
		products     service.Products
		experiences  service.Experiences
		promotions   service.Promotions
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		users, err = h.multiGetUsers(ectx, orders.UserIDs())
		return
	})
	eg.Go(func() (err error) {
		coordinators, err = h.multiGetCoordinatorsWithDeleted(ectx, orders.CoordinatorIDs())
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProductsByRevision(ectx, orders.ProductRevisionIDs())
		return
	})
	eg.Go(func() (err error) {
		experiences, err = h.multiGetExperiencesByRevision(ectx, orders.ExperienceRevisionIDs())
		return
	})
	eg.Go(func() (err error) {
		addresses, err = h.multiGetAddressesByRevision(ectx, orders.AddressRevisionIDs())
		return
	})
	eg.Go(func() (err error) {
		promotions, err = h.multiGetPromotions(ectx, orders.PromotionIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.OrdersResponse{
		Orders:       service.NewOrders(orders, addresses.MapByRevision(), products.MapByRevision(), experiences.MapByRevision()).Response(),
		Users:        users.Response(),
		Coordinators: coordinators.Response(),
		Promotions:   promotions.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newOrderFileters(ctx *gin.Context) ([]sentity.OrderStatus, []sentity.OrderType, error) {
	sparams, err := util.GetQueryInt32s(ctx, "statuses")
	if err != nil {
		return nil, nil, fmt.Errorf("handler: failed to get status query params: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	tparams, err := util.GetQueryInt32s(ctx, "types")
	if err != nil {
		return nil, nil, fmt.Errorf("handler: failed to get type query params: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}

	statuses := make([]sentity.OrderStatus, len(sparams))
	for i := range sparams {
		statuses[i] = sentity.OrderStatus(sparams[i])
	}
	if len(statuses) == 0 {
		statuses = []sentity.OrderStatus{
			sentity.OrderStatusWaiting,   // 受注待ち
			sentity.OrderStatusPreparing, // 発送準備中
			sentity.OrderStatusShipped,   // 発送完了
			sentity.OrderStatusCompleted, // 完了
		}
	}

	types := make([]sentity.OrderType, len(tparams))
	for i := range tparams {
		types[i] = sentity.OrderType(tparams[i])
	}
	if len(types) == 0 {
		types = []sentity.OrderType{
			sentity.OrderTypeProduct,    // 商品
			sentity.OrderTypeExperience, // 体験
		}
	}

	return statuses, types, nil
}

// @Summary     注文取得
// @Description 指定された注文の詳細情報を取得します。
// @Tags        Order
// @Router      /v1/orders/{orderId} [get]
// @Security    bearerauth
// @Param       orderId path string true "注文ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.OrderResponse
// @Failure     404 {object} util.ErrorResponse "注文が存在しない"
func (h *handler) GetOrder(ctx *gin.Context) {
	order, err := h.getOrder(ctx, util.GetParam(ctx, "orderId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	var (
		user        *service.User
		coordinator *service.Coordinator
		promotion   *service.Promotion
		products    service.Products
		experience  *service.Experience
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		user, err = h.getUser(ectx, order.UserID)
		return
	})
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinatorWithDeleted(ectx, order.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		if order.PromotionID == "" {
			return nil
		}
		promotion, err = h.getPromotion(ectx, order.PromotionID)
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProducts(ectx, order.ProductIDs())
		return
	})
	eg.Go(func() (err error) {
		if order.Experience == nil || order.Experience.ExperienceID == "" {
			return nil
		}
		experience, err = h.getExperience(ectx, order.Experience.ExperienceID)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.OrderResponse{
		Order:       order.Response(),
		User:        user.Response(),
		Coordinator: coordinator.Response(),
		Promotion:   promotion.Response(),
		Products:    products.Response(),
		Experience:  experience.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     注文下書き保存
// @Description 注文の下書きを保存します。
// @Tags        Order
// @Router      /v1/orders/{orderId}/draft [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.DraftOrderRequest true "注文下書き保存"
// @Param       orderId path string true "注文ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     403 {object} util.ErrorResponse "操作の権限がない"
// @Failure     404 {object} util.ErrorResponse "注文が存在しない"
func (h *handler) DraftOrder(ctx *gin.Context) {
	req := &types.DraftOrderRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.DraftOrderInput{
		OrderID:         util.GetParam(ctx, "orderId"),
		ShippingMessage: req.ShippingMessage,
	}
	if err := h.store.DraftOrder(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     注文確定
// @Description 注文を確定します。
// @Tags        Order
// @Router      /v1/orders/{orderId}/capture [post]
// @Security    bearerauth
// @Param       orderId path string true "注文ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     403 {object} util.ErrorResponse "注文の確定権限がない"
// @Failure     404 {object} util.ErrorResponse "注文が存在しない"
func (h *handler) CaptureOrder(ctx *gin.Context) {
	in := &store.CaptureOrderInput{
		OrderID: util.GetParam(ctx, "orderId"),
	}
	if err := h.store.CaptureOrder(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     注文対応完了
// @Description 注文対応を完了します。
// @Tags        Order
// @Router      /v1/orders/{orderId}/complete [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.CompleteOrderRequest true "注文対応完了"
// @Param       orderId path string true "注文ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     403 {object} util.ErrorResponse "操作の権限がない"
// @Failure     404 {object} util.ErrorResponse "注文が存在しない"
func (h *handler) CompleteOrder(ctx *gin.Context) {
	req := &types.CompleteOrderRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	order, err := h.getOrder(ctx, util.GetParam(ctx, "orderId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	switch order.Type {
	case types.OrderTypeProduct:
		in := &store.CompleteProductOrderInput{
			OrderID:         order.ID,
			ShippingMessage: req.ShippingMessage,
		}
		err = h.store.CompleteProductOrder(ctx, in)
	case types.OrderTypeExperience:
		in := &store.CompleteExperienceOrderInput{
			OrderID: order.ID,
		}
		err = h.store.CompleteExperienceOrder(ctx, in)
	default:
		err = fmt.Errorf("handler: unknown order type: %w", exception.ErrInternal)
	}
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     注文キャンセル
// @Description 注文をキャンセルします。
// @Tags        Order
// @Router      /v1/orders/{orderId}/cancel [post]
// @Security    bearerauth
// @Param       orderId path string true "注文ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     403 {object} util.ErrorResponse "注文のキャンセル権限がない"
// @Failure     404 {object} util.ErrorResponse "注文が存在しない"
func (h *handler) CancelOrder(ctx *gin.Context) {
	in := &store.CancelOrderInput{
		OrderID: util.GetParam(ctx, "orderId"),
	}
	if err := h.store.CancelOrder(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     注文の返金依頼
// @Description 注文の返金を依頼します。
// @Tags        Order
// @Router      /v1/orders/{orderId}/refund [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.RefundOrderRequest true "注文の返金依頼"
// @Param       orderId path string true "注文ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     403 {object} util.ErrorResponse "操作の権限がない"
// @Failure     404 {object} util.ErrorResponse "注文が存在しない"
func (h *handler) RefundOrder(ctx *gin.Context) {
	req := &types.RefundOrderRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.RefundOrderInput{
		OrderID:     util.GetParam(ctx, "orderId"),
		Description: req.Description,
	}
	if err := h.store.RefundOrder(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) UpdateOrderFulfillment(ctx *gin.Context) {
	req := &types.UpdateOrderFulfillmentRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.UpdateOrderFulfillmentInput{
		OrderID:         util.GetParam(ctx, "orderId"),
		FulfillmentID:   util.GetParam(ctx, "fulfillmentId"),
		ShippingCarrier: sentity.ShippingCarrier(req.ShippingCarrier),
		TrackingNumber:  req.TrackingNumber,
	}
	if err := h.store.UpdateOrderFulfillment(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     注文履歴のCSV出力
// @Description 注文履歴をCSV形式で出力します。
// @Tags        Order
// @Router      /v1/orders/-/export [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.ExportOrdersRequest true "注文履歴のCSV出力"
// @Produce     json
// @Success     204
// @Failure     403 {object} util.ErrorResponse "操作の権限がない"
func (h *handler) ExportOrders(ctx *gin.Context) {
	req := &types.ExportOrdersRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.ExportOrdersInput{
		ShopID:          getShopID(ctx),
		ShippingCarrier: sentity.ShippingCarrier(req.ShippingCarrier),
		EncodingType:    codes.CharacterEncodingType(req.CharacterEncodingType),
	}
	value, err := h.store.ExportOrders(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	filename := fmt.Sprintf("orders_%s.csv", h.now().Format("20060102150405"))
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Writer.Header().Set("Content-Type", "text/csv")
	if _, err := ctx.Writer.Write(value); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *handler) getOrder(ctx context.Context, orderID string) (*service.Order, error) {
	in := &store.GetOrderInput{
		OrderID: orderID,
	}
	order, err := h.store.GetOrder(ctx, in)
	if err != nil {
		return nil, err
	}
	var (
		addresses   service.Addresses
		products    service.Products
		experiences service.Experiences
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		addresses, err = h.multiGetAddressesByRevision(ectx, order.AddressRevisionIDs())
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProductsByRevision(ectx, order.ProductRevisionIDs())
		return
	})
	eg.Go(func() (err error) {
		if order.ExperienceRevisionID == 0 {
			return
		}
		experiences, err = h.multiGetExperiencesByRevision(ectx, []int64{order.ExperienceRevisionID})
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return service.NewOrder(order, addresses.MapByRevision(), products.MapByRevision(), experiences.MapByRevision()), nil
}
