package mysql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const scheduleTable = "schedules"

type schedule struct {
	db  *mysql.Client
	now func() time.Time
}

func newSchedule(db *mysql.Client) database.Schedule {
	return &schedule{
		db:  db,
		now: jst.Now,
	}
}

type listSchedulesParams database.ListSchedulesParams

func (p listSchedulesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if len(p.Types) > 0 {
		stmt = stmt.Where("message_type IN (?)", p.Types)
	}
	if len(p.Statuses) > 0 {
		stmt = stmt.Where("status IN (?)", p.Statuses)
	}
	if !p.Since.IsZero() {
		stmt = stmt.Where("sent_at >= ?", p.Since)
	}
	if !p.Until.IsZero() {
		stmt = stmt.Where("sent_at <= ?", p.Until)
	}
	return stmt
}

func (s *schedule) List(ctx context.Context, params *database.ListSchedulesParams, fields ...string) (entity.Schedules, error) {
	var schedules entity.Schedules

	p := listSchedulesParams(*params)

	stmt := s.db.Statement(ctx, s.db.DB, scheduleTable, fields...)
	stmt = p.stmt(stmt)

	err := stmt.Find(&schedules).Error
	return schedules, dbError(err)
}

func (s *schedule) UpsertProcessing(ctx context.Context, schedule *entity.Schedule) error {
	err := s.db.Transaction(ctx, func(tx *gorm.DB) error {
		current, err := s.get(ctx, tx, schedule.MessageType, schedule.MessageID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		now := s.now()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			schedule.CreatedAt, schedule.UpdatedAt = now, now
		} else {
			if !current.Executable(now) {
				return fmt.Errorf("database: schedule is not executable %w", database.ErrFailedPrecondition)
			}
			schedule.UpdatedAt = now
		}
		schedule.Status = entity.ScheduleStatusProcessing
		schedule.Count++

		err = tx.WithContext(ctx).Table(scheduleTable).Save(&schedule).Error
		return err
	})
	return dbError(err)
}

func (s *schedule) UpdateDone(ctx context.Context, messageType entity.ScheduleType, messageID string) error {
	err := s.db.Transaction(ctx, func(tx *gorm.DB) error {
		current, err := s.get(ctx, tx, messageType, messageID, "status")
		if err != nil {
			return err
		}

		if current.Status == entity.ScheduleStatusDone {
			return fmt.Errorf("database: schedule is already done: %w", database.ErrFailedPrecondition)
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
		return err
	})
	return dbError(err)
}

func (s *schedule) UpdateCancel(ctx context.Context, messageType entity.ScheduleType, messageID string) error {
	err := s.db.Transaction(ctx, func(tx *gorm.DB) error {
		current, err := s.get(ctx, tx, messageType, messageID)
		if err != nil {
			return err
		}

		now := s.now()
		if !current.ShouldCancel(now) {
			return fmt.Errorf("database: schedule should not cancel: %w", database.ErrFailedPrecondition)
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
		return err
	})
	return dbError(err)
}

func (s *schedule) get(
	ctx context.Context, tx *gorm.DB, messageType entity.ScheduleType, messageID string, fields ...string,
) (*entity.Schedule, error) {
	var schedule *entity.Schedule

	stmt := s.db.Statement(ctx, tx, scheduleTable, fields...).
		Where("message_type = ?", messageType).
		Where("message_id = ?", messageID)

	if err := stmt.First(&schedule).Error; err != nil {
		return nil, err
	}
	return schedule, nil
}
