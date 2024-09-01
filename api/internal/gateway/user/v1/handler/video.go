package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/media"
)

type listVideoSummariesParams struct {
	coordinatorID string
	limit         int64
	offset        int64
	noLimit       bool
}

func (h *handler) listProductVideoSummaries(ctx context.Context, params *listVideoSummariesParams) (service.VideoSummaries, error) {
	in := &media.ListVideosInput{
		CoordinatorID:      params.coordinatorID,
		OnlyPublished:      true,
		OnlyDisplayProduct: true,
		ExcludeLimited:     true,
		Limit:              params.limit,
		Offset:             params.offset,
		NoLimit:            params.noLimit,
	}
	videos, _, err := h.media.ListVideos(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewVideoSummaries(videos), nil
}

func (h *handler) listExperienceVideoSummaries(ctx context.Context, params *listVideoSummariesParams) (service.VideoSummaries, error) {
	in := &media.ListVideosInput{
		CoordinatorID:         params.coordinatorID,
		OnlyPublished:         true,
		OnlyDisplayExperience: true,
		ExcludeLimited:        true,
		Limit:                 params.limit,
		Offset:                params.offset,
		NoLimit:               params.noLimit,
	}
	videos, _, err := h.media.ListVideos(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewVideoSummaries(videos), nil
}
