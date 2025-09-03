package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/service"
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
	r.GET("/:orderId", h.GetOrder)
}

// @Summary     注文一覧取得
// @Description 注文の一覧を取得します。
// @Tags        Order
// @Router      /facilities/{facilityId}/orders [get]
// @Param       facilityId path string true "施設ID"
// @Security    bearerauth
// @Param       limit query int64 false "取得件数" default(20)
// @Param       offset query int64 false "取得開始位置" default(0)
// @Param       types query []int32 false "注文ステータス" collectionFormat(csv)
// @Produce     json
// @Success     200 {object} response.OrdersResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
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
	types, err := util.GetQueryInt32s(ctx, "types")
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	orderTypes := make([]sentity.OrderType, len(types))
	for i, t := range types {
		orderTypes[i] = sentity.OrderType(t)
	}
	orderStatuses := []sentity.OrderStatus{
		sentity.OrderStatusUnpaid,
		sentity.OrderStatusWaiting,
		sentity.OrderStatusPreparing,
		sentity.OrderStatusShipped,
		sentity.OrderStatusCompleted,
		sentity.OrderStatusCanceled,
		sentity.OrderStatusRefunded,
	}
	ordersIn := &store.ListOrdersInput{
		UserID:   h.getUserID(ctx),
		Limit:    limit,
		Offset:   offset,
		Types:    orderTypes,
		Statuses: orderStatuses,
	}
	orders, total, err := h.store.ListOrders(ctx, ordersIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(orders) == 0 {
		res := &response.OrdersResponse{
			Order:        []*response.Order{},
			Coordinators: []*response.Coordinator{},
			Promotions:   []*response.Promotion{},
			Products:     []*response.Product{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	var (
		coordinators service.Coordinators
		promotions   service.Promotions
		oproducts    service.Products
		cproducts    service.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinators, err = h.multiGetCoordinatorsWithDeleted(ectx, orders.CoordinatorIDs())
		return
	})
	eg.Go(func() (err error) {
		promotions, err = h.multiGetPromotion(ectx, orders.PromotionIDs())
		return
	})
	eg.Go(func() (err error) {
		oproducts, err = h.multiGetProductsByRevision(ectx, h.getProducerID(ctx), orders.ProductRevisionIDs())
		if err != nil {
			return
		}
		cproducts, err = h.multiGetProducts(ectx, h.getProducerID(ctx), oproducts.IDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.OrdersResponse{
		Order:        service.NewOrders(orders, oproducts.MapByRevision()).Response(),
		Coordinators: coordinators.Response(),
		Promotions:   promotions.Response(),
		Products:     cproducts.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     注文詳細取得
// @Description 注文の詳細を取得します。
// @Tags        Order
// @Router      /facilities/{facilityId}/orders/{orderId} [get]
// @Param       facilityId path string true "施設ID"
// @Param       orderId path string true "注文ID"
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} response.OrderResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     404 {object} util.ErrorResponse "注文が見つからない"
func (h *handler) GetOrder(ctx *gin.Context) {
	order, err := h.getOrder(ctx, h.getProducerID(ctx), h.getUserID(ctx), util.GetParam(ctx, "orderId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		coordinator *service.Coordinator
		promotion   *service.Promotion
		products    service.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinatorWithDeleted(ectx, order.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		if order.PromotionID == "" {
			return
		}
		promotion, err = h.getPromotion(ectx, order.PromotionID)
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProducts(ectx, h.getProducerID(ctx), order.ProductIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.OrderResponse{
		Order:       order.Response(),
		Coordinator: coordinator.Response(),
		Promotion:   promotion.Response(),
		Products:    products.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) getOrder(ctx context.Context, producerID, userID, orderID string) (*service.Order, error) {
	in := &store.GetOrderInput{
		OrderID: orderID,
	}
	order, err := h.store.GetOrder(ctx, in)
	if err != nil {
		return nil, err
	}
	if userID != order.UserID {
		// 不正の疑いがあるため、リクエスト情報をログ出力しておく
		slog.WarnContext(ctx, "UserId does not match order information", slog.String("userId", userID), slog.String("orderId", orderID))
		return nil, fmt.Errorf("%w: %w", exception.ErrNotFound, errNotFoundOrder)
	}
	products, err := h.multiGetProductsByRevision(ctx, producerID, order.ProductRevisionIDs())
	if err != nil {
		return nil, err
	}
	return service.NewOrder(order, products.MapByRevision()), nil
}
