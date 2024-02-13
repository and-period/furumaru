package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) scheduleRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/schedules")

	// TODO: メモリリークの原因切り分けのため、一時的にコメントアウト
	// r.GET("/:scheduleId", h.createBroadcastViewerLog, h.GetSchedule)

	r.GET("/:scheduleId", h.GetSchedule)
}

func (h *handler) GetSchedule(ctx *gin.Context) {
	schedule, err := h.getSchedule(ctx, util.GetParam(ctx, "scheduleId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	in := &store.ListLivesInput{
		ScheduleIDs: []string{schedule.ID},
		NoLimit:     true,
	}
	lives, _, err := h.store.ListLives(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		coordinator *service.Coordinator
		producers   service.Producers
		products    service.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinator(ectx, schedule.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		producers, err = h.multiGetProducers(ectx, lives.ProducerIDs())
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProducts(ectx, lives.ProductIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ScheduleResponse{
		Schedule:    schedule.Response(),
		Coordinator: coordinator.Response(),
		Lives:       service.NewLives(lives).Response(),
		Producers:   producers.Response(),
		Products:    products.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) getSchedule(ctx context.Context, scheduleID string) (*service.Schedule, error) {
	scheduleIn := &store.GetScheduleInput{
		ScheduleID: scheduleID,
	}
	schedule, err := h.store.GetSchedule(ctx, scheduleIn)
	if err != nil {
		return nil, err
	}
	broadcastIn := &media.GetBroadcastByScheduleIDInput{
		ScheduleID: schedule.ID,
	}
	broadcast, err := h.media.GetBroadcastByScheduleID(ctx, broadcastIn)
	if err != nil && !errors.Is(err, exception.ErrNotFound) {
		return nil, err
	}
	return service.NewSchedule(schedule, broadcast), nil
}
