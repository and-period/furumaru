package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) CreateBroadcastViewerLog(
	ctx context.Context,
	in *media.CreateBroadcastViewerLogInput,
) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return internalError(err)
	}
	params := &entity.BroadcastViewerLogParams{
		BroadcastID: broadcast.ID,
		SessionID:   in.SessionID,
		UserID:      in.UserID,
		UserAgent:   in.UserAgent,
		ClientIP:    in.ClientIP,
	}
	log := entity.NewBroadcastViewerLog(params)
	err = s.db.BroadcastViewerLog.Create(ctx, log)
	return internalError(err)
}

func (s *service) AggregateBroadcastViewerLogs(
	ctx context.Context, in *media.AggregateBroadcastViewerLogsInput,
) (entity.AggregatedBroadcastViewerLogs, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return nil, 0, internalError(err)
	}
	var (
		logs  entity.AggregatedBroadcastViewerLogs
		total int64
	)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		params := &database.AggregateBroadcastViewerLogsParams{
			BroadcastID:  broadcast.ID,
			Interval:     in.Interval,
			CreatedAtGte: in.CreatedAtGte,
			CreatedAtLt:  in.CreatedAtLt,
		}
		logs, err = s.db.BroadcastViewerLog.Aggregate(ctx, params)
		return
	})
	eg.Go(func() (err error) {
		params := &database.GetBroadcastTotalViewersParams{
			BroadcastID:  broadcast.ID,
			CreatedAtGte: in.CreatedAtGte,
			CreatedAtLt:  in.CreatedAtLt,
		}
		total, err = s.db.BroadcastViewerLog.GetTotal(ctx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return logs, total, nil
}
