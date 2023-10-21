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
	"github.com/and-period/furumaru/api/pkg/backoff"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListSchedules(ctx context.Context, in *store.ListSchedulesInput) (entity.Schedules, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListSchedulesParams{
		StartAtGte: in.StartAtGte,
		StartAtLt:  in.StartAtLt,
		EndAtGte:   in.EndAtGte,
		EndAtLt:    in.EndAtLt,
		Statuses:   in.Statuses,
		Limit:      int(in.Limit),
		Offset:     int(in.Offset),
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
		return nil, 0, internalError(err)
	}
	return schedules, total, nil
}

func (s *service) MultiGetSchedules(ctx context.Context, in *store.MultiGetSchedulesInput) (entity.Schedules, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	schedules, err := s.db.Schedule.MultiGet(ctx, in.ScheduleIDs)
	return schedules, internalError(err)
}

func (s *service) GetSchedule(ctx context.Context, in *store.GetScheduleInput) (*entity.Schedule, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	schedule, err := s.db.Schedule.Get(ctx, in.ScheduleID)
	return schedule, internalError(err)
}

func (s *service) CreateSchedule(ctx context.Context, in *store.CreateScheduleInput) (*entity.Schedule, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
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
		return nil, internalError(err)
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
	if err := s.db.Schedule.Create(ctx, schedule); err != nil {
		return nil, internalError(err)
	}
	s.waitGroup.Add(2)
	go func() {
		defer s.waitGroup.Done()
		s.resizeSchedule(context.Background(), schedule.ID, in.ThumbnailURL)
	}()
	go func() {
		defer s.waitGroup.Done()
		const maxRetries = 3
		ctx := context.Background()
		createFn := func() error {
			in := &media.CreateBroadcastInput{
				ScheduleID: schedule.ID,
			}
			_, err := s.media.CreateBroadcast(ctx, in)
			return err
		}
		retry := backoff.NewExponentialBackoff(maxRetries)
		if err := backoff.Retry(ctx, retry, createFn, backoff.WithRetryablel(s.isRetryable)); err != nil {
			s.logger.Error("Failed to create broadcast", zap.String("scheduleId", schedule.ID), zap.Error(err))
		}
	}()
	return schedule, nil
}

func (s *service) UpdateSchedule(ctx context.Context, in *store.UpdateScheduleInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	schedule, err := s.db.Schedule.Get(ctx, in.ScheduleID)
	if err != nil {
		return internalError(err)
	}
	_, err = s.db.Shipping.Get(ctx, in.ShippingID)
	if errors.Is(err, exception.ErrNotFound) {
		return fmt.Errorf("api: invalid request: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return internalError(err)
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
		return internalError(err)
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
		return internalError(err)
	}
	err := s.db.Schedule.UpdateThumbnails(ctx, in.ScheduleID, in.Thumbnails)
	return internalError(err)
}

func (s *service) ApproveSchedule(ctx context.Context, in *store.ApproveScheduleInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	adminIn := &user.GetAdministratorInput{
		AdministratorID: in.AdminID,
	}
	_, err := s.user.GetAdministrator(ctx, adminIn)
	if errors.Is(err, exception.ErrNotFound) {
		return fmt.Errorf("api: invalid request: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return internalError(err)
	}
	params := &database.ApproveScheduleParams{
		Approved:        in.Approved,
		ApprovedAdminID: in.AdminID,
	}
	err = s.db.Schedule.Approve(ctx, in.ScheduleID, params)
	return internalError(err)
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
