package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *service) CreateSchedule(ctx context.Context, in *store.CreateScheduleInput) (*entity.Schedule, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	params := &entity.NewScheduleParams{
		Title:        in.Title,
		Description:  in.Description,
		ThumbnailURL: in.ThumbnailURL,
		StartAt:      in.StartAt,
		EndAt:        in.EndAt,
	}
	schedule := entity.NewSchedule(params)
	if err := s.db.Schedule.Create(ctx, schedule); err != nil {
		return nil, exception.InternalError(err)
	}
	return schedule, nil
}
