package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *service) GetBroadcastByScheduleID(ctx context.Context, in *store.GetBroadcastByScheduleIDInput) (*entity.Broadcast, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	return broadcast, exception.InternalError(err)
}
