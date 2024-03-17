package handler

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/store"
)

const liveEndDuration = time.Hour // 終了予定時間を過ぎたあとも配信しているケースを救うため

type listLiveSummariesParams struct {
	coordinatorID string
	producerID    string
	limit         int64
	offset        int64
	noLimit       bool
}

func (h *handler) listLiveSummaries(ctx context.Context, params *listLiveSummariesParams) (service.LiveSummaries, int64, error) {
	in := &store.ListSchedulesInput{
		CoordinatorID: params.coordinatorID,
		ProducerID:    params.producerID,
		EndAtGte:      h.now().Add(-liveEndDuration),
		OnlyPublished: true,
		Limit:         params.limit,
		Offset:        params.offset,
		NoLimit:       params.noLimit,
	}
	schedules, total, err := h.store.ListSchedules(ctx, in)
	if err != nil || len(schedules) == 0 {
		return service.LiveSummaries{}, 0, err
	}
	livesIn := &store.ListLivesInput{
		ScheduleIDs:   schedules.IDs(),
		ProducerID:    params.producerID,
		NoLimit:       true,
		OnlyPublished: true,
	}
	lives, _, err := h.store.ListLives(ctx, livesIn)
	if err != nil {
		return nil, 0, err
	}
	productsIn := &store.MultiGetProductsInput{
		ProductIDs: lives.ProductIDs(),
	}
	products, err := h.store.MultiGetProducts(ctx, productsIn)
	if err != nil {
		return nil, 0, err
	}
	return service.NewLiveSummaries(schedules, lives, products), total, nil
}

type listArchiveSummariesParams struct {
	coordinatorID string
	producerID    string
	limit         int64
	offset        int64
	noLimit       bool
}

func (h *handler) listArchiveSummaries(ctx context.Context, params *listArchiveSummariesParams) (service.ArchiveSummaries, int64, error) {
	in := &store.ListSchedulesInput{
		CoordinatorID: params.coordinatorID,
		ProducerID:    params.producerID,
		EndAtLt:       h.now().Add(-liveEndDuration),
		OnlyPublished: true,
		Limit:         params.limit,
		Offset:        params.offset,
		NoLimit:       params.noLimit,
	}
	archives, total, err := h.store.ListSchedules(ctx, in)
	if err != nil {
		return nil, 0, err
	}
	return service.NewArchiveSummaries(archives), total, nil
}
