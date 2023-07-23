package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) threadRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("", h.ListThreadsByContactID)
	arg.GET("/:threadId", h.GetThread)
	arg.POST("", h.CreateThread)
	arg.PATCH("/:threadId", h.UpdateThread)
	arg.DELETE("/:threadId", h.DeleteThread)
}

func (h *handler) ListThreadsByContactID(ctx *gin.Context) {
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
	threadsIn := &messenger.ListThreadsByContactIDInput{
		ContactID: util.GetParam(ctx, "contactId"),
		Limit:     limit,
		Offset:    offset,
	}

	sthreads, total, err := h.messenger.ListThreadsByContactID(ctx, threadsIn)
	if err != nil {
		httpError(ctx, err)
	}
	if len(sthreads) == 0 {
		res := &response.ThreadsResponse{
			Threads: []*response.Thread{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	var (
		users  uentity.Users
		admins service.Admins
	)

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		usersIn := &user.MultiGetUsersInput{
			UserIDs: sthreads.UserIDs(),
		}
		users, err = h.user.MultiGetUsers(ectx, usersIn)
		return
	})
	eg.Go(func() (err error) {
		admins, err = h.multiGetAdmins(ectx, sthreads.AdminIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	threads := service.NewThreads(sthreads)
	res := &response.ThreadsResponse{
		Threads: threads.Response(),
		Users:   service.NewUsers(users).Response(),
		Admins:  admins.Response(),
		Total:   total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateThread(ctx *gin.Context) {
	req := &request.CreateThreadRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	var (
		radmin *service.Admin
		ruser  *uentity.User
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if req.UserType == 1 {
			radmin, err = h.getAdmin(ectx, req.UserID)
		}
		return
	})
	eg.Go(func() (err error) {
		if req.UserType == 2 {
			in := &user.GetUserInput{
				UserID: req.UserID,
			}
			ruser, err = h.user.GetUser(ectx, in)
		}
		return
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	in := &messenger.CreateThreadInput{
		ContactID: req.ContactID,
		UserID:    req.UserID,
		UserType:  service.ThreadUserType(req.UserType).StoreEntity(),
		Content:   req.Content,
	}
	sthread, err := h.messenger.CreateThread(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	thread := service.NewThread(sthread)

	res := &response.ThreadResponse{
		Thread: thread.Response(),
		User:   service.NewUser(ruser).Response(),
		Admin:  radmin.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetThread(ctx *gin.Context) {
	thread, err := h.getThread(ctx, util.GetParam(ctx, "threadId"))
	if err != nil {
		httpError(ctx, err)
		return
	}

	var (
		radmin *service.Admin
		ruser  *uentity.User
	)

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if thread.UserType == 1 {
			radmin, err = h.getAdmin(ectx, thread.UserID)
		}
		return
	})
	eg.Go(func() (err error) {
		if thread.UserType == 2 {
			in := &user.GetUserInput{
				UserID: thread.UserID,
			}
			ruser, err = h.user.GetUser(ectx, in)
		}
		return
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ThreadResponse{
		Thread: thread.Response(),
		User:   service.NewUser(ruser).Response(),
		Admin:  radmin.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) getThread(ctx context.Context, threadID string) (*service.Thread, error) {
	in := &messenger.GetThreadInput{
		ThreadID: threadID,
	}
	mthread, err := h.messenger.GetThread(ctx, in)
	if err != nil {
		return nil, err
	}
	thread := service.NewThread(mthread)
	return thread, nil
}

func (h *handler) UpdateThread(ctx *gin.Context) {
	req := &request.UpdateThreadRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &messenger.UpdateThreadInput{
		ThreadID: util.GetParam(ctx, "threadId"),
		Content:  req.Content,
		UserID:   req.UserID,
		UserType: service.ThreadUserType(req.UserType).StoreEntity(),
	}

	if err := h.messenger.UpdateThread(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) DeleteThread(ctx *gin.Context) {
	in := &messenger.DeleteThreadInput{
		ThreadID: util.GetParam(ctx, "threadId"),
	}
	if err := h.messenger.DeleteThread(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
