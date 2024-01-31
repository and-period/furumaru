package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/media"
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
	if broadcast.Status == entity.BroadcastStatusDisabled {
		return nil // 視聴ログを書き込む必要がないため、後続処理はしない
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
