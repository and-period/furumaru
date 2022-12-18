package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) orderRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("", h.ListOrders)
	arg.GET("/:orderId", h.filterAccessOrder, h.GetOrder)
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
		httpError(ctx, err)
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
		badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	os, err := h.newOrderOrders(ctx)
	if err != nil {
		badRequest(ctx, err)
		return
	}

	in := &store.ListOrdersInput{
		Limit:  limit,
		Offset: offset,
		Orders: os,
	}
	if getRole(ctx) == service.AdminRoleCoordinator {
		in.CoordinatorID = getAdminID(ctx)
	}
	sorders, total, err := h.store.ListOrders(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	orders := service.NewOrders(sorders)
	if len(orders) == 0 {
		res := &response.OrdersResponse{
			Orders: orders.Response(),
		}
		ctx.JSON(http.StatusOK, res)
		return
	}
	if err := h.getOrderDetails(ctx, orders...); err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.OrdersResponse{
		Orders: orders.Response(),
		Total:  total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newOrderOrders(ctx *gin.Context) ([]*store.ListOrdersOrder, error) {
	orders := map[string]sentity.OrderOrderBy{
		"paymentStatus":     sentity.OrderOrderByPaymentStatus,
		"fulfillmentStatus": sentity.OrderOrderByFulfillmentStatus,
		"orderedAt":         sentity.OrderOrderByOrderedAt,
		"paidAt":            sentity.OrderOrderByConfirmedAt,
		"deliveredAt":       sentity.OrderOrderByDeliveredAt,
		"canceledAt":        sentity.OrderOrderByCanceledAt,
		"createdAt":         sentity.OrderOrderByCreatedAt,
		"updatedAt":         sentity.OrderOrderByUpdatedAt,
	}
	params := util.GetOrders(ctx)
	res := make([]*store.ListOrdersOrder, len(params))
	for i, p := range params {
		key, ok := orders[p.Key]
		if !ok {
			return nil, fmt.Errorf("handler: unknown order key. key=%s: %w", p.Key, errInvalidOrderkey)
		}
		res[i] = &store.ListOrdersOrder{
			Key:        key,
			OrderByASC: p.Direction == util.OrderByASC,
		}
	}
	return res, nil
}

func (h *handler) GetOrder(ctx *gin.Context) {
	order, err := h.getOrder(ctx, util.GetParam(ctx, "orderId"))
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.OrderResponse{
		Order: order.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) getOrder(ctx context.Context, orderID string) (*service.Order, error) {
	in := &store.GetOrderInput{
		OrderID: orderID,
	}
	sorder, err := h.store.GetOrder(ctx, in)
	if err != nil {
		return nil, err
	}
	order := service.NewOrder(sorder)
	if err := h.getOrderDetails(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (h *handler) getOrderDetails(ctx context.Context, orders ...*service.Order) error {
	os := service.Orders(orders)
	var (
		users     service.Users
		products  service.Products
		addresses service.Addresses
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		in := &user.MultiGetUsersInput{
			UserIDs: os.UserIDs(),
		}
		uusers, err := h.user.MultiGetUsers(ectx, in)
		if err != nil {
			return err
		}
		users = service.NewUsers(uusers)
		return nil
	})
	eg.Go(func() error {
		in := &store.MultiGetProductsInput{
			ProductIDs: os.ProductIDs(),
		}
		sproducts, err := h.store.MultiGetProducts(ectx, in)
		if err != nil {
			return err
		}
		products = service.NewProducts(sproducts)
		return nil
	})
	eg.Go(func() error {
		in := &store.MultiGetAddressesInput{
			AddressIDs: os.AddressIDs(),
		}
		saddresses, err := h.store.MultiGetAddresses(ectx, in)
		if err != nil {
			return err
		}
		addresses = service.NewAddresses(saddresses)
		return nil
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	os.Fill(users.Map(), products.Map(), addresses.Map())
	return nil
}
