package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListSchedules(ctx context.Context, in *store.ListSchedulesInput) (entity.Schedules, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	params := &database.ListSchedulesParams{
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	var (
		schedules entity.Schedules
		total     int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		schedules, err = s.db.Schedule.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Schedule.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	return schedules, total, nil
}

func (s *service) GetSchedule(ctx context.Context, in *store.GetScheduleInput) (*entity.Schedule, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	schedule, err := s.db.Schedule.Get(ctx, in.ScheduleID)
	return schedule, exception.InternalError(err)
}

func (s *service) CreateSchedule(ctx context.Context, in *store.CreateScheduleInput) (*entity.Schedule, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &user.GetCoordinatorInput{
			CoordinatorID: in.CoordinatorID,
		}
		_, err = s.user.GetCoordinator(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		_, err = s.db.Shipping.Get(ectx, in.ShippingID)
		return
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		return nil, fmt.Errorf("api: invalid request: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, exception.InternalError(err)
	}
	sparams := &entity.NewScheduleParams{
		CoordinatorID:   in.CoordinatorID,
		ShippingID:      in.ShippingID,
		Title:           in.Title,
		Description:     in.Description,
		ThumbnailURL:    in.ThumbnailURL,
		ImageURL:        in.ImageURL,
		OpeningVideoURL: in.OpeningVideoURL,
		Public:          in.Public,
		StartAt:         in.StartAt,
		EndAt:           in.EndAt,
	}
	schedule := entity.NewSchedule(sparams)
	bparams := &entity.NewBroadcastParams{
		ScheduleID: schedule.ID,
	}
	broadcast := entity.NewBroadcast(bparams)
	if err := s.db.Schedule.Create(ctx, schedule, broadcast); err != nil {
		return nil, exception.InternalError(err)
	}
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		s.resizeSchedule(context.Background(), schedule.ID, in.ThumbnailURL)
	}()
	return schedule, nil
}

func (s *service) UpdateSchedule(ctx context.Context, in *store.UpdateScheduleInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	schedule, err := s.db.Schedule.Get(ctx, in.ScheduleID)
	if err != nil {
		return exception.InternalError(err)
	}
	_, err = s.db.Shipping.Get(ctx, in.ShippingID)
	if errors.Is(err, exception.ErrNotFound) {
		return fmt.Errorf("api: invalid request: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return exception.InternalError(err)
	}
	params := &database.UpdateScheduleParams{
		ShippingID:      in.ShippingID,
		Title:           in.Title,
		Description:     in.Description,
		ThumbnailURL:    in.ThumbnailURL,
		ImageURL:        in.ImageURL,
		OpeningVideoURL: in.OpeningVideoURL,
		Public:          in.Public,
		StartAt:         in.StartAt,
		EndAt:           in.EndAt,
	}
	if err := s.db.Schedule.Update(ctx, in.ScheduleID, params); err != nil {
		return exception.InternalError(err)
	}
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		var thumbnailURL string
		if schedule.ThumbnailURL != in.ThumbnailURL {
			thumbnailURL = in.ThumbnailURL
		}
		s.resizeSchedule(context.Background(), schedule.ID, thumbnailURL)
	}()
	return nil
}

func (s *service) UpdateScheduleThumbnails(ctx context.Context, in *store.UpdateScheduleThumbnailsInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.Schedule.UpdateThumbnails(ctx, in.ScheduleID, in.Thumbnails)
	return exception.InternalError(err)
}

func (s *service) resizeSchedule(ctx context.Context, scheduleID, thumbnailURL string) {
	if thumbnailURL == "" {
		return
	}
	in := &media.ResizeFileInput{
		TargetID: scheduleID,
		URLs:     []string{thumbnailURL},
	}
	if err := s.media.ResizeScheduleThumbnail(ctx, in); err != nil {
		s.logger.Error("Failed to resize schedule thumbnail",
			zap.String("scheduleId", scheduleID), zap.Error(err),
		)
	}
}
