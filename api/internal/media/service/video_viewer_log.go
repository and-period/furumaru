package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) CreateVideoViewerLog(
	ctx context.Context,
	in *media.CreateVideoViewerLogInput,
) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	video, err := s.db.Video.Get(ctx, in.VideoID)
	if err != nil {
		return internalError(err)
	}
	if !video.Published() {
		return fmt.Errorf("%w: video is not published", exception.ErrFailedPrecondition)
	}
	params := &entity.NewVideoViewerLogParams{
		VideoID:   video.ID,
		SessionID: in.SessionID,
		UserID:    in.UserID,
		UserAgent: in.UserAgent,
		ClientIP:  in.ClientIP,
	}
	log := entity.NewVideoViewerLog(params)
	err = s.db.VideoViewerLog.Create(ctx, log)
	return internalError(err)
}

func (s *service) AggregateVideoViewerLogs(
	ctx context.Context, in *media.AggregateVideoViewerLogsInput,
) (entity.AggregatedVideoViewerLogs, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	var (
		logs  entity.AggregatedVideoViewerLogs
		total int64
	)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		params := &database.AggregateVideoViewerLogsParams{
			VideoID:      in.VideoID,
			Interval:     in.Interval,
			CreatedAtGte: in.CreatedAtGte,
			CreatedAtLt:  in.CreatedAtLt,
		}
		logs, err = s.db.VideoViewerLog.Aggregate(ctx, params)
		return
	})
	eg.Go(func() (err error) {
		params := &database.GetVideoTotalViewersParams{
			VideoID:      in.VideoID,
			CreatedAtGte: in.CreatedAtGte,
			CreatedAtLt:  in.CreatedAtLt,
		}
		total, err = s.db.VideoViewerLog.GetTotal(ctx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return logs, total, nil
}
