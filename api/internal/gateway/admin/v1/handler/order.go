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
	"github.com/and-period/furumaru/api/internal/store/entity"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

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
	statuses, err := h.newOrderFileters(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	orderType := util.GetQuery(ctx, "type", "")

	in := &store.ListOrdersInput{
		Limit:    limit,
		Offset:   offset,
		Statuses: statuses,
	}
	if orderType != "" {
		in.Types = []entity.OrderType{service.NewOrderTypeFromString(orderType).StoreEntity()}
	}
	if getRole(ctx) == service.AdminRoleCoordinator {
		in.CoordinatorID = getAdminID(ctx)
	}
	orders, total, err := h.store.ListOrders(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(orders) == 0 {
		res := &response.OrdersResponse{
			Orders:       []*response.Order{},
			Users:        []*response.User{},
			Coordinators: []*response.Coordinator{},
			Promotions:   []*response.Promotion{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	var (
		users        service.Users
		coordinators service.Coordinators
		addresses    service.Addresses
		products     service.Products
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

	res := &response.OrdersResponse{
		Orders:       service.NewOrders(orders, addresses.MapByRevision(), products.MapByRevision()).Response(),
		Users:        users.Response(),
		Coordinators: coordinators.Response(),
		Promotions:   promotions.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newOrderFileters(ctx *gin.Context) ([]sentity.OrderStatus, error) {
	params, err := util.GetQueryInt32s(ctx, "status")
	if err != nil {
		return nil, err
	}
	if len(params) == 0 {
		res := []sentity.OrderStatus{
			sentity.OrderStatusWaiting,   // 受注待ち
			sentity.OrderStatusPreparing, // 発送準備中
			sentity.OrderStatusShipped,   // 発送完了
			sentity.OrderStatusCompleted, // 完了
		}
		return res, nil
	}
	res := make([]sentity.OrderStatus, len(params))
	for i := range params {
		res[i] = sentity.OrderStatus(params[i])
	}
	return res, nil
}

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
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.OrderResponse{
		Order:       order.Response(),
		User:        user.Response(),
		Coordinator: coordinator.Response(),
		Promotion:   promotion.Response(),
		Products:    products.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) DraftOrder(ctx *gin.Context) {
	req := &request.DraftOrderRequest{}
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

func (h *handler) CompleteOrder(ctx *gin.Context) {
	req := &request.CompleteOrderRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.CompleteOrderInput{
		OrderID:         util.GetParam(ctx, "orderId"),
		ShippingMessage: req.ShippingMessage,
	}
	if err := h.store.CompleteOrder(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

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

func (h *handler) RefundOrder(ctx *gin.Context) {
	req := &request.RefundOrderRequest{}
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
	req := &request.UpdateOrderFulfillmentRequest{}
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

func (h *handler) ExportOrders(ctx *gin.Context) {
	req := &request.ExportOrdersRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.ExportOrdersInput{
		ShippingCarrier: sentity.ShippingCarrier(req.ShippingCarrier),
		EncodingType:    codes.CharacterEncodingType(req.CharacterEncodingType),
	}
	if getRole(ctx) == service.AdminRoleCoordinator {
		in.CoordinatorID = getAdminID(ctx)
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
		addresses service.Addresses
		products  service.Products
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
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return service.NewOrder(order, addresses.MapByRevision(), products.MapByRevision()), nil
}
