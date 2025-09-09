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
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        User
// @tag.description 購入者関連
func (h *handler) userRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/users", h.authentication)

	r.GET("", h.ListUsers)
	r.GET("/:userId", h.GetUser)
	r.DELETE("/:userId", h.DeleteUser)
	r.GET("/:userId/orders", h.ListUserOrders)
}

// @Summary     購入者一覧取得
// @Description 購入者の一覧を取得します。管理者は全購入者、コーディネーターは注文実績のある購入者のみ取得可能です。
// @Tags        User
// @Router      /v1/users [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Produce     json
// @Success     200 {object} types.UsersResponse
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

	var (
		users service.Users
		total int64
	)
	switch getAdminType(ctx) {
	case service.AdminTypeAdministrator:
		// 管理者の場合、すべての購入者情報を取得する
		usersIn := &user.ListUsersInput{
			Limit:       limit,
			Offset:      offset,
			WithDeleted: true,
		}
		var us uentity.Users
		us, total, err = h.user.ListUsers(ctx, usersIn)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		if len(us) == 0 {
			break
		}
		addressesIn := &user.ListDefaultAddressesInput{
			UserIDs: us.IDs(),
		}
		addresses, err := h.user.ListDefaultAddresses(ctx, addressesIn)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		users = service.NewUsers(us, addresses.MapByUserID())
	case service.AdminTypeCoordinator:
		// コーディネータの場合、注文した購入者のみを取得する
		in := &store.ListOrderUserIDsInput{
			ShopID: getShopID(ctx),
			Limit:  limit,
			Offset: offset,
		}
		var userIDs []string
		userIDs, total, err = h.store.ListOrderUserIDs(ctx, in)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		if len(userIDs) == 0 {
			break
		}
		users, err = h.multiGetUsers(ctx, userIDs)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
	default:
		h.forbidden(ctx, errors.New("handler: forbidden"))
		return
	}
	if len(users) == 0 {
		res := &types.UsersResponse{
			Users: []*types.UserToList{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	in := &store.AggregateOrdersByUserInput{
		ShopID:  getShopID(ctx),
		UserIDs: users.IDs(),
	}
	orders, err := h.store.AggregateOrdersByUser(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.UsersResponse{
		Users: service.NewUsersToList(users, orders.Map()).Response(),
		Total: total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     購入者取得
// @Description 指定された購入者の詳細情報を取得します。
// @Tags        User
// @Router      /v1/users/{userId} [get]
// @Security    bearerauth
// @Param       userId path string true "購入者ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.UserResponse
// @Failure     404 {object} util.ErrorResponse "購入者が存在しない"
func (h *handler) GetUser(ctx *gin.Context) {
	user, err := h.getUser(ctx, util.GetParam(ctx, "userId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.UserResponse{
		User:    user.Response(),
		Address: user.Address().Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     購入者削除
// @Description 購入者を削除します。管理者のみ実行可能です。
// @Tags        User
// @Router      /v1/users/{userId} [delete]
// @Security    bearerauth
// @Param       userId path string true "購入者ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     404 {object} util.ErrorResponse "購入者が存在しない"
func (h *handler) DeleteUser(ctx *gin.Context) {
	in := &user.DeleteUserInput{
		UserID: util.GetParam(ctx, "userId"),
	}
	if err := h.user.DeleteUser(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     購入者注文履歴取得
// @Description 指定された購入者の注文履歴と注文統計情報を取得します。
// @Tags        User
// @Router      /v1/users/{userId}/orders [get]
// @Security    bearerauth
// @Param       userId path string true "購入者ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Produce     json
// @Success     200 {object} types.UserOrdersResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     404 {object} util.ErrorResponse "購入者が存在しない"
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
		aggregatedOrder *sentity.AggregatedUserOrder
		total           int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &store.ListOrdersInput{
			ShopID: getShopID(ctx),
			UserID: userID,
			Limit:  limit,
			Offset: offset,
		}
		orders, total, err = h.store.ListOrders(ectx, in)
		return
	})
	eg.Go(func() error {
		in := &store.AggregateOrdersByUserInput{
			ShopID:  getShopID(ctx),
			UserIDs: []string{userID},
		}
		aggregate, err := h.store.AggregateOrdersByUser(ectx, in)
		if err != nil {
			return err
		}
		order, ok := aggregate.Map()[userID]
		if !ok {
			order = &sentity.AggregatedUserOrder{}
		}
		aggregatedOrder = order
		return nil
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.UserOrdersResponse{
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
