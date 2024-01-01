package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/store"
)

func (s *service) ReserveStartLive(ctx context.Context, in *messenger.ReserveStartLiveInput) error {
	const (
		messageType = entity.ScheduleTypeStartLive
		duration    = time.Hour // ライブ配信開始の１時間前通知
	)
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	scheduleIn := &store.GetScheduleInput{
		ScheduleID: in.ScheduleID,
	}
	schedule, err := s.store.GetSchedule(ctx, scheduleIn)
	if errors.Is(err, exception.ErrNotFound) {
		return fmt.Errorf("service: not found schedule: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return internalError(err)
	}
	params := &upsertScheduleParams{
		messageType: messageType,
		messageID:   schedule.ID,
		sentAt:      schedule.StartAt.Add(-duration),
		deadline:    schedule.StartAt,
	}
	return s.upsertSchedule(ctx, params)
}

func (s *service) ReserveNotification(ctx context.Context, in *messenger.ReserveNotificationInput) error {
	const messageType = entity.ScheduleTypeNotification
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	notification, err := s.db.Notification.Get(ctx, in.NotificationID)
	if errors.Is(err, database.ErrNotFound) {
		return fmt.Errorf("service: not found notification: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return internalError(err)
	}
	params := &upsertScheduleParams{
		messageType: messageType,
		messageID:   notification.ID,
		sentAt:      notification.PublishedAt,
	}
	return s.upsertSchedule(ctx, params)
}

type upsertScheduleParams struct {
	messageType entity.ScheduleType
	messageID   string
	sentAt      time.Time
	deadline    time.Time
}

func (s *service) upsertSchedule(ctx context.Context, params *upsertScheduleParams) error {
	p := &entity.NewScheduleParams{
		MessageType: params.messageType,
		MessageID:   params.messageID,
		SentAt:      params.sentAt,
		Deadline:    params.deadline,
	}
	schedule := entity.NewSchedule(p)
	if err := s.db.Schedule.Upsert(ctx, schedule); err != nil {
		return internalError(err)
	}
	return nil
}
