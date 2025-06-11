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

	r.GET("/lives", h.ListLiveSchedules)
	r.GET("/archives", h.ListArchiveSchedules)
	r.GET("/:scheduleId", h.createBroadcastViewerLog, h.GetSchedule)
}

func (h *handler) ListLiveSchedules(ctx *gin.Context) {
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
	var shopID string
	if coordinatorID := util.GetQuery(ctx, "coordinator", ""); coordinatorID != "" {
		coordinator, err := h.getCoordinator(ctx, coordinatorID)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		shopID = coordinator.ShopID
	}

	params := &listLiveSummariesParams{
		shopID:     shopID,
		producerID: util.GetQuery(ctx, "producer", ""),
		limit:      limit,
		offset:     offset,
	}
	lives, total, err := h.listLiveSummaries(ctx, params)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	coordinators, err := h.multiGetCoordinators(ctx, lives.CoordinatorIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.LiveSchedulesResponse{
		Lives:        lives.Response(),
		Coordinators: coordinators.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) ListArchiveSchedules(ctx *gin.Context) {
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
	var shopID string
	if coordinatorID := util.GetQuery(ctx, "coordinator", ""); coordinatorID != "" {
		coordinator, err := h.getCoordinator(ctx, coordinatorID)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		shopID = coordinator.ShopID
	}

	params := &listArchiveSummariesParams{
		shopID:     shopID,
		producerID: util.GetQuery(ctx, "producer", ""),
		limit:      limit,
		offset:     offset,
	}
	archives, total, err := h.listArchiveSummaries(ctx, params)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	coordinators, err := h.multiGetCoordinators(ctx, archives.CoordinatorIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ArchiveSchedulesResponse{
		Archives:     archives.Response(),
		Coordinators: coordinators.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
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
