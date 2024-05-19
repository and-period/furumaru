package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	mentity "github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/messenger"
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
		CoordinatorID: in.CoordinatorID,
		ProducerID:    in.ProducerID,
		StartAtGte:    in.StartAtGte,
		StartAtLt:     in.StartAtLt,
		EndAtGte:      in.EndAtGte,
		EndAtLt:       in.EndAtLt,
		OnlyPublished: in.OnlyPublished,
		Limit:         int(in.Limit),
		Offset:        int(in.Offset),
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
	coordinatorIn := &user.GetCoordinatorInput{
		CoordinatorID: in.CoordinatorID,
	}
	_, err := s.user.GetCoordinator(ctx, coordinatorIn)
	if errors.Is(err, exception.ErrNotFound) {
		return nil, fmt.Errorf("api: invalid request: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, internalError(err)
	}
	sparams := &entity.NewScheduleParams{
		CoordinatorID:   in.CoordinatorID,
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
		const maxRetries = 5
		ctx := context.Background()
		createFn := func() error {
			in := &media.CreateBroadcastInput{
				ScheduleID:    schedule.ID,
				CoordinatorID: schedule.CoordinatorID,
			}
			_, err := s.media.CreateBroadcast(ctx, in)
			return err
		}
		retry := backoff.NewExponentialBackoff(maxRetries)
		if err := backoff.Retry(ctx, retry, createFn, backoff.WithRetryablel(s.isRetryable)); err != nil {
			s.logger.Error("Failed to create broadcast", zap.String("scheduleId", schedule.ID), zap.Error(err))
		}
	}()
	go func() {
		defer s.waitGroup.Done()
		in := &messenger.ReserveStartLiveInput{
			ScheduleID: schedule.ID,
		}
		if err := s.messenger.ReserveStartLive(context.Background(), in); err != nil {
			s.logger.Error("Failed to reserve start live", zap.String("scheduleId", schedule.ID), zap.Error(err))
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
	params := &database.UpdateScheduleParams{
		Title:           in.Title,
		Description:     in.Description,
		ThumbnailURL:    in.ThumbnailURL,
		ImageURL:        in.ImageURL,
		OpeningVideoURL: in.OpeningVideoURL,
		StartAt:         in.StartAt,
		EndAt:           in.EndAt,
	}
	if err := s.db.Schedule.Update(ctx, in.ScheduleID, params); err != nil {
		return internalError(err)
	}
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		in := &messenger.ReserveStartLiveInput{
			ScheduleID: schedule.ID,
		}
		if err := s.messenger.ReserveStartLive(context.Background(), in); err != nil {
			s.logger.Error("Failed to reserve start live", zap.String("scheduleId", schedule.ID), zap.Error(err))
		}
	}()
	return nil
}

func (s *service) DeleteSchedule(ctx context.Context, in *store.DeleteScheduleInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	schedule, err := s.db.Schedule.Get(ctx, in.ScheduleID)
	if err != nil {
		return internalError(err)
	}
	broadcastIn := &media.GetBroadcastByScheduleIDInput{
		ScheduleID: schedule.ID,
	}
	broadcast, err := s.media.GetBroadcastByScheduleID(ctx, broadcastIn)
	if err != nil && !errors.Is(err, exception.ErrNotFound) {
		return internalError(err)
	}
	if broadcast != nil && broadcast.Status != mentity.BroadcastStatusDisabled {
		return fmt.Errorf("api: invalid request: broadcast is not disabled: %w", exception.ErrFailedPrecondition)
	}
	err = s.db.Schedule.Delete(ctx, in.ScheduleID)
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

func (s *service) PublishSchedule(ctx context.Context, in *store.PublishScheduleInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Schedule.Publish(ctx, in.ScheduleID, in.Public)
	return internalError(err)
}
