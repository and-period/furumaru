package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) userRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/users", h.authentication)

	r.GET("", h.ListUsers)
	r.GET("/:userId", h.GetUser)
	r.DELETE("/:userId", h.DeleteUser)
	r.GET("/:userId/orders", h.ListUserOrders)
}

func (h *handler) ListUsers(ctx *gin.Context) {
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

	in := &user.ListUsersInput{
		Limit:       limit,
		Offset:      offset,
		WithDeleted: true,
	}
	users, total, err := h.user.ListUsers(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(users) == 0 {
		res := &response.UsersResponse{
			Users: []*response.UserToList{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	var (
		orders    sentity.AggregatedOrders
		addresses uentity.Addresses
	)
	userIDs := users.IDs()
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &store.AggregateOrdersInput{
			UserIDs: userIDs,
		}
		if getRole(ctx) == service.AdminRoleCoordinator {
			in.CoordinatorID = getAdminID(ctx)
		}
		orders, err = h.store.AggregateOrders(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &user.ListDefaultAddressesInput{
			UserIDs: userIDs,
		}
		addresses, err = h.user.ListDefaultAddresses(ectx, in)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	susers := service.NewUsers(users, addresses.MapByUserID())
	res := &response.UsersResponse{
		Users: service.NewUsersToList(susers, orders.Map()).Response(),
		Total: total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetUser(ctx *gin.Context) {
	user, err := h.getUser(ctx, util.GetParam(ctx, "userId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.UserResponse{
		User:    user.Response(),
		Address: user.Address().Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) DeleteUser(ctx *gin.Context) {
	in := &user.DeleteUserInput{
		UserID: util.GetParam(ctx, "userId"),
	}
	if err := h.user.DeleteUser(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) ListUserOrders(ctx *gin.Context) {
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
	userID := util.GetParam(ctx, "userId")
	if userID == "" {
		h.badRequest(ctx, errors.New("handler: userId is required"))
		return
	}

	var (
		orders          sentity.Orders
		aggregatedOrder *sentity.AggregatedOrder
		total           int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &store.ListOrdersInput{
			UserID: userID,
			Limit:  limit,
			Offset: offset,
		}
		if getRole(ctx) == service.AdminRoleCoordinator {
			in.CoordinatorID = getAdminID(ctx)
		}
		orders, total, err = h.store.ListOrders(ectx, in)
		return
	})
	eg.Go(func() error {
		in := &store.AggregateOrdersInput{
			UserIDs: []string{userID},
		}
		if getRole(ctx) == service.AdminRoleCoordinator {
			in.CoordinatorID = getAdminID(ctx)
		}
		aggregate, err := h.store.AggregateOrders(ectx, in)
		if err != nil {
			return err
		}
		order, ok := aggregate.Map()[userID]
		if !ok {
			order = &sentity.AggregatedOrder{}
		}
		aggregatedOrder = order
		return nil
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.UserOrdersResponse{
		Orders:             service.NewUserOrders(orders).Response(),
		OrderTotalCount:    total,
		PaymentTotalCount:  aggregatedOrder.OrderCount,
		ProductTotalAmount: aggregatedOrder.Subtotal,
		PaymentTotalAmount: aggregatedOrder.Total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) multiGetUsers(ctx context.Context, userIDs []string) (service.Users, error) {
	if len(userIDs) == 0 {
		return service.Users{}, nil
	}
	var (
		users     uentity.Users
		addresses uentity.Addresses
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &user.MultiGetUsersInput{
			UserIDs: userIDs,
		}
		users, err = h.user.MultiGetUsers(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &user.ListDefaultAddressesInput{
			UserIDs: userIDs,
		}
		addresses, err = h.user.ListDefaultAddresses(ectx, in)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return service.NewUsers(users, addresses.MapByUserID()), nil
}

func (h *handler) getUser(ctx context.Context, userID string) (*service.User, error) {
	userIn := &user.GetUserInput{
		UserID: userID,
	}
	u, err := h.user.GetUser(ctx, userIn)
	if err != nil {
		return nil, err
	}
	addressIn := &user.GetDefaultAddressInput{
		UserID: userID,
	}
	address, err := h.user.GetDefaultAddress(ctx, addressIn)
	if err != nil && !errors.Is(err, exception.ErrNotFound) {
		return nil, err
	}
	return service.NewUser(u, address), nil
}
