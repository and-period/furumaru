package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) threadRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/contacts/:contactId/threads", h.authentication)

	r.GET("", h.ListThreadsByContactID)
	r.POST("", h.CreateThread)
	r.GET("/:threadId", h.GetThread)
	r.PATCH("/:threadId", h.UpdateThread)
	r.DELETE("/:threadId", h.DeleteThread)
}

func (h *handler) ListThreadsByContactID(ctx *gin.Context) {
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
	threadsIn := &messenger.ListThreadsInput{
		ContactID: util.GetParam(ctx, "contactId"),
		Limit:     limit,
		Offset:    offset,
	}

	sthreads, total, err := h.messenger.ListThreads(ctx, threadsIn)
	if err != nil {
		h.httpError(ctx, err)
	}
	if len(sthreads) == 0 {
		res := &response.ThreadsResponse{
			Threads: []*response.Thread{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	var (
		users  service.Users
		admins service.Admins
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		users, err = h.multiGetUsers(ectx, sthreads.UserIDs())
		return
	})
	eg.Go(func() (err error) {
		admins, err = h.multiGetAdmins(ectx, sthreads.AdminIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	threads := service.NewThreads(sthreads)
	res := &response.ThreadsResponse{
		Threads: threads.Response(),
		Users:   users.Response(),
		Admins:  admins.Response(),
		Total:   total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateThread(ctx *gin.Context) {
	req := &request.CreateThreadRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	var (
		admin *service.Admin
		user  *service.User
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if req.UserType != 1 {
			return
		}
		admin, err = h.getAdmin(ectx, req.UserID)
		return
	})
	eg.Go(func() (err error) {
		if req.UserType != 2 {
			return
		}
		user, err = h.getUser(ectx, req.UserID)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
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
		h.httpError(ctx, err)
		return
	}
	thread := service.NewThread(sthread)

	res := &response.ThreadResponse{
		Thread: thread.Response(),
		User:   user.Response(),
		Admin:  admin.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetThread(ctx *gin.Context) {
	thread, err := h.getThread(ctx, util.GetParam(ctx, "threadId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		admin *service.Admin
		user  *service.User
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if thread.UserType != 1 {
			return
		}
		admin, err = h.getAdmin(ectx, thread.UserID)
		return
	})
	eg.Go(func() (err error) {
		if thread.UserType != 2 {
			return
		}
		user, err = h.getUser(ectx, thread.UserID)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ThreadResponse{
		Thread: thread.Response(),
		User:   user.Response(),
		Admin:  admin.Response(),
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
		h.badRequest(ctx, err)
		return
	}

	in := &messenger.UpdateThreadInput{
		ThreadID: util.GetParam(ctx, "threadId"),
		Content:  req.Content,
		UserID:   req.UserID,
		UserType: service.ThreadUserType(req.UserType).StoreEntity(),
	}

	if err := h.messenger.UpdateThread(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteThread(ctx *gin.Context) {
	in := &messenger.DeleteThreadInput{
		ThreadID: util.GetParam(ctx, "threadId"),
	}
	if err := h.messenger.DeleteThread(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
