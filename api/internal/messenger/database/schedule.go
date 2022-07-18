package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/gorm"
)

const scheduleTable = "schedules"

var scheduleFields = []string{
	"message_type", "message_id", "status", "count",
	"sent_at", "created_at", "updated_at",
}

type schedule struct {
	db  *database.Client
	now func() time.Time
}

func NewSchedule(db *database.Client) Schedule {
	return &schedule{
		db:  db,
		now: jst.Now,
	}
}

func (s *schedule) List(ctx context.Context, params *ListSchedulesParams, fields ...string) (entity.Schedules, error) {
	var schedules entity.Schedules
	if len(fields) == 0 {
		fields = scheduleFields
	}

	stmt := s.db.DB.WithContext(ctx).Table(scheduleTable).Select(fields)
	stmt = params.stmt(stmt)

	err := stmt.Find(&schedules).Error
	return schedules, exception.InternalError(err)
}

func (s *schedule) UpsertProcessing(ctx context.Context, schedule *entity.Schedule) error {
	_, err := s.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		current, err := s.get(ctx, tx, schedule.MessageType, schedule.MessageID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		now := s.now()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			schedule.CreatedAt, schedule.UpdatedAt = now, now
		} else {
			if !current.Executable(now) {
				return nil, fmt.Errorf("database: schedule is not executable %w", exception.ErrFailedPrecondition)
			}
			schedule.UpdatedAt = now
		}
		schedule.Status = entity.ScheduleStatusProcessing
		schedule.Count++

		err = tx.WithContext(ctx).Table(scheduleTable).Save(&schedule).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (s *schedule) UpdateDone(ctx context.Context, messageType entity.ScheduleType, messageID string) error {
	_, err := s.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		current, err := s.get(ctx, tx, messageType, messageID, "status")
		if err != nil {
			return nil, err
		}

		if current.Status == entity.ScheduleStatusDone {
			return nil, fmt.Errorf("database: schedule is already done: %w", exception.ErrFailedPrecondition)
		}

		params := map[string]interface{}{
			"status":     entity.ScheduleStatusDone,
			"updated_at": s.now(),
		}
		err = tx.WithContext(ctx).
			Table(scheduleTable).
			Where("message_type = ?", messageType).
			Where("message_id = ?", messageID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (s *schedule) UpdateCancel(ctx context.Context, messageType entity.ScheduleType, messageID string) error {
	_, err := s.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		current, err := s.get(ctx, tx, messageType, messageID)
		if err != nil {
			return nil, err
		}

		now := s.now()
		if !current.ShouldCancel(now) {
			return nil, fmt.Errorf("database: schedule should not cancel: %w", exception.ErrFailedPrecondition)
		}

		params := map[string]interface{}{
			"status":     entity.ScheduleStatusCanceled,
			"updated_at": now,
		}
		err = tx.WithContext(ctx).
			Table(scheduleTable).
			Where("message_type = ?", messageType).
			Where("message_id = ?", messageID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (s *schedule) get(
	ctx context.Context, tx *gorm.DB, messageType entity.ScheduleType, messageID string, fields ...string,
) (*entity.Schedule, error) {
	var schedule *entity.Schedule
	if len(fields) == 0 {
		fields = scheduleFields
	}

	err := tx.WithContext(ctx).
		Table(scheduleTable).Select(fields).
		Where("message_type = ?", messageType).
		Where("message_id = ?", messageID).
		First(&schedule).Error
	return schedule, err
}
