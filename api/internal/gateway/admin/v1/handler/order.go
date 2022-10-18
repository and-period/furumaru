package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
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
	// TODO: ソート設定

	in := &store.ListOrdersInput{
		Limit:  limit,
		Offset: offset,
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

	var (
		users    service.Users
		products service.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		in := &user.MultiGetUsersInput{
			UserIDs: orders.UserIDs(),
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
			ProductIDs: orders.ProductIDs(),
		}
		sproducts, err := h.store.MultiGetProducts(ectx, in)
		if err != nil {
			return err
		}
		products = service.NewProducts(sproducts)
		return nil
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	orders.Fill(users.Map(), products.Map())

	res := &response.OrdersResponse{
		Orders: orders.Response(),
		Total:  total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetOrder(ctx *gin.Context) {
	in := &store.GetOrderInput{
		OrderID: util.GetParam(ctx, "orderId"),
	}
	sorder, err := h.store.GetOrder(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	order := service.NewOrder(sorder)

	var (
		u        *service.User
		products service.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		in := &user.GetUserInput{
			UserID: order.UserID,
		}
		uuser, err := h.user.GetUser(ectx, in)
		if err != nil {
			return err
		}
		u = service.NewUser(uuser)
		return nil
	})
	eg.Go(func() error {
		in := &store.MultiGetProductsInput{
			ProductIDs: order.ProductIDs(),
		}
		sproducts, err := h.store.MultiGetProducts(ectx, in)
		if err != nil {
			return err
		}
		products = service.NewProducts(sproducts)
		return nil
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	order.Fill(u, products.Map())

	res := &response.OrderResponse{
		Order: order.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) getOrder(ctx context.Context, orderID string) (*service.Order, error) {
	in := &store.GetOrderInput{
		OrderID: orderID,
	}
	order, err := h.store.GetOrder(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewOrder(order), nil
}
