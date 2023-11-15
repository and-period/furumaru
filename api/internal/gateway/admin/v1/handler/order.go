package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) orderRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/orders", h.authentication)

	r.GET("", h.ListOrders)
	r.GET("/:orderId", h.filterAccessOrder, h.GetOrder)
	r.POST("/:orderId/capture", h.filterAccessOrder, h.CaptureOrder)
	r.POST("/:orderId/cancel", h.filterAccessOrder, h.CancelOrder)
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
		producer: func(ctx *gin.Context) (bool, error) {
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

	in := &store.ListOrdersInput{
		Limit:  limit,
		Offset: offset,
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
			Orders: []*response.Order{},
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
		coordinators, err = h.multiGetCoordinators(ectx, orders.CoordinatorIDs())
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
		coordinator, err = h.getCoordinator(ectx, order.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		if order.PromotionID == "" {
			return nil
		}
		promotion, err = h.getPromotion(ctx, order.PromotionID)
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

func (h *handler) CaptureOrder(ctx *gin.Context) {
	in := &store.CaptureOrderInput{
		OrderID: util.GetParam(ctx, "orderId"),
	}
	if err := h.store.CaptureOrder(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) CancelOrder(ctx *gin.Context) {
	in := &store.CancelOrderInput{
		OrderID: util.GetParam(ctx, "orderId"),
	}
	if err := h.store.CancelOrder(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
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
