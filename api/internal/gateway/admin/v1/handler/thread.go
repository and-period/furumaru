package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/gin-gonic/gin"
)

func (h *handler) ThreadRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("", h.ListThreadsByContactID)
	arg.GET("/:threadId", h.GetThread)
}

func (h *handler) ListThreadsByContactID(ctx *gin.Context) {
	threadsIn := &messenger.ListThreadsByContactIDInput{
		ContactID: util.GetParam(ctx, "contactId"),
	}

	sthreads, err := h.messenger.ListThreadsByContactID(ctx, threadsIn)
	if err != nil {
		httpError(ctx, err)
	}
	threads := service.NewThreads(sthreads)
	if len(threads) == 0 {
		res := &response.ThreadsResponse{
			Threads: []*response.Thread{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	res := &response.ThreadsResponse{
		Threads: threads.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetThread(ctx *gin.Context) {
	thread, err := h.getThread(ctx, util.GetParam(ctx, "threadId"))
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ThreadResponse{
		Thread: thread.Response(),
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
