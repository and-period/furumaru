package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Thread
// @tag.description スレッド関連
func (h *handler) threadRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/contacts/:contactId/threads", h.authentication)

	r.GET("", h.ListThreadsByContactID)
	r.POST("", h.CreateThread)
	r.GET("/:threadId", h.GetThread)
	r.PATCH("/:threadId", h.UpdateThread)
	r.DELETE("/:threadId", h.DeleteThread)
}

// @Summary     お問い合わせスレッド一覧取得
// @Description 指定されたお問い合わせのスレッド一覧を取得します。
// @Tags        Thread
// @Router      /v1/contacts/{contactId}/threads [get]
// @Security    bearerauth
// @Param       contactId path string true "お問い合わせID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Produce     json
// @Success     200 {object} types.ThreadsResponse
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
		res := &types.ThreadsResponse{
			Threads: []*types.Thread{},
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
	res := &types.ThreadsResponse{
		Threads: threads.Response(),
		Users:   users.Response(),
		Admins:  admins.Response(),
		Total:   total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     スレッド登録
// @Description 新しいスレッドを登録します。
// @Tags        Thread
// @Router      /v1/contacts/{contactId}/threads [post]
// @Security    bearerauth
// @Param       contactId path string true "お問い合わせID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.CreateThreadRequest true "スレッド情報"
// @Produce     json
// @Success     200 {object} types.ThreadResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateThread(ctx *gin.Context) {
	req := &types.CreateThreadRequest{}
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

	res := &types.ThreadResponse{
		Thread: thread.Response(),
		User:   user.Response(),
		Admin:  admin.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     スレッド取得
// @Description 指定されたスレッドの詳細情報を取得します。
// @Tags        Thread
// @Router      /v1/contacts/{contactId}/threads/{threadId} [get]
// @Security    bearerauth
// @Param       contactId path string true "お問い合わせID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       threadId path string true "スレッドID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.ThreadResponse
// @Failure     404 {object} util.ErrorResponse "スレッドが存在しない"
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

	res := &types.ThreadResponse{
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

// @Summary     スレッド更新
// @Description スレッドの情報を更新します。
// @Tags        Thread
// @Router      /v1/contacts/{contactId}/threads/{threadId} [patch]
// @Security    bearerauth
// @Param       contactId path string true "お問い合わせID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       threadId path string true "スレッドID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.UpdateThreadRequest true "スレッド情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     404 {object} util.ErrorResponse "スレッドが存在しない"
func (h *handler) UpdateThread(ctx *gin.Context) {
	req := &types.UpdateThreadRequest{}
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

// @Summary     スレッド削除
// @Description スレッドを削除します。
// @Tags        Thread
// @Router      /v1/contacts/{contactId}/threads/{threadId} [delete]
// @Security    bearerauth
// @Param       contactId path string true "お問い合わせID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       threadId path string true "スレッドID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     404 {object} util.ErrorResponse "スレッドが存在しない"
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
