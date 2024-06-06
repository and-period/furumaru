package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

func (s *service) CreateBroadcastViewerLog(ctx context.Context, in *media.CreateBroadcastViewerLogInput) error {
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
) (entity.AggregatedBroadcastViewerLogs, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return nil, internalError(err)
	}
	params := &database.AggregateBroadcastViewerLogsParams{
		BroadcastID:  broadcast.ID,
		Interval:     in.Interval,
		CreatedAtGte: in.CreatedAtGte,
		CreatedAtLt:  in.CreatedAtLt,
	}
	logs, err := s.db.BroadcastViewerLog.Aggregate(ctx, params)
	return logs, internalError(err)
}
