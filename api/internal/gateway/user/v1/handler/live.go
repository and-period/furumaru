package handler

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/store"
)

const liveEndDuration = time.Hour // 終了予定時間を過ぎたあとも配信しているケースを救うため

type listLiveSummariesParams struct {
	coordinatorID string
}

func (h *handler) listLiveSummaries(ctx context.Context, params *listLiveSummariesParams) (service.LiveSummaries, error) {
	in := &store.ListSchedulesInput{
		CoordinatorID: params.coordinatorID,
		EndAtGte:      h.now().Add(-liveEndDuration),
		OnlyPublished: true,
		NoLimit:       true,
	}
	schedules, _, err := h.store.ListSchedules(ctx, in)
	if err != nil || len(schedules) == 0 {
		return service.LiveSummaries{}, err
	}
	livesIn := &store.ListLivesInput{
		ScheduleIDs: schedules.IDs(),
		NoLimit:     true,
	}
	lives, _, err := h.store.ListLives(ctx, livesIn)
	if err != nil || len(lives) == 0 {
		return service.LiveSummaries{}, err
	}
	productsIn := &store.MultiGetProductsInput{
		ProductIDs: lives.ProductIDs(),
	}
	products, err := h.store.MultiGetProducts(ctx, productsIn)
	if err != nil {
		return nil, err
	}
	return service.NewLiveSummaries(schedules, lives, products), nil
}

type listArchiveSummariesParams struct {
	coordinatorID string
	limit         int64
	offset        int64
	noLimit       bool
}

func (h *handler) listArchiveSummaries(ctx context.Context, params *listArchiveSummariesParams) (service.ArchiveSummaries, error) {
	broadcastsIn := &media.ListBroadcastsInput{
		CoordinatorID: params.coordinatorID,
		OnlyArchived:  true,
		Limit:         params.limit,
		Offset:        params.offset,
		NoLimit:       params.noLimit,
	}
	broadcasts, _, err := h.media.ListBroadcasts(ctx, broadcastsIn)
	if err != nil || len(broadcasts) == 0 {
		return service.ArchiveSummaries{}, err
	}
	schedulesIn := &store.MultiGetSchedulesInput{
		ScheduleIDs: broadcasts.ScheduleIDs(),
	}
	archives, err := h.store.MultiGetSchedules(ctx, schedulesIn)
	if err != nil {
		return nil, err
	}
	return service.NewArchiveSummaries(archives), nil
}
