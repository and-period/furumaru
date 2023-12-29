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
}

func (h *handler) listLiveSummaries(ctx context.Context, params *listLiveSummariesParams) (service.LiveSummaries, error) {
	in := &store.ListSchedulesInput{
		CoordinatorID: params.coordinatorID,
		ProducerID:    params.producerID,
		EndAtGte:      h.now().Add(-liveEndDuration),
		OnlyPublished: true,
		NoLimit:       true,
	}
	schedules, _, err := h.store.ListSchedules(ctx, in)
	if err != nil || len(schedules) == 0 {
		return service.LiveSummaries{}, err
	}
	livesIn := &store.ListLivesInput{
		ScheduleIDs:   schedules.IDs(),
		ProducerID:    params.producerID,
		NoLimit:       true,
		OnlyPublished: true,
	}
	lives, _, err := h.store.ListLives(ctx, livesIn)
	if err != nil {
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
	producerID    string
	limit         int64
	offset        int64
	noLimit       bool
}

func (h *handler) listArchiveSummaries(ctx context.Context, params *listArchiveSummariesParams) (service.ArchiveSummaries, error) {
	in := &store.ListSchedulesInput{
		CoordinatorID: params.coordinatorID,
		ProducerID:    params.producerID,
		EndAtLt:       h.now().Add(-liveEndDuration),
		OnlyPublished: true,
		NoLimit:       true,
	}
	archives, _, err := h.store.ListSchedules(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewArchiveSummaries(archives), nil
}
