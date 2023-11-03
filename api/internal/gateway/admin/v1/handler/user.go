package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

func (h *handler) userRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/users", h.authentication)

	r.GET("", h.ListUsers)
	r.GET("/:userId", h.GetUser)
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

	usersIn := &user.ListUsersInput{
		Limit:  limit,
		Offset: offset,
	}
	uusers, total, err := h.user.ListUsers(ctx, usersIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	users := service.NewUsers(uusers)
	if len(users) == 0 {
		res := &response.UsersResponse{
			Users: []*response.UserSummary{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	ordersIn := &store.AggregateOrdersInput{
		UserIDs: uusers.IDs(),
	}
	sorders, err := h.store.AggregateOrders(ctx, ordersIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.UsersResponse{
		Users: service.NewUserSummaries(users, sorders.Map()).Response(),
		Total: total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetUser(ctx *gin.Context) {
	in := &user.GetUserInput{
		UserID: util.GetParam(ctx, "userId"),
	}
	uuser, err := h.user.GetUser(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.UserResponse{
		User: service.NewUser(uuser).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
