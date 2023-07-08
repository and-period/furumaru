package database

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const scheduleTable = "schedules"

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

	stmt := s.db.Statement(ctx, s.db.DB, scheduleTable, fields...)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	if err := stmt.Find(&schedules).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	if err := schedules.Fill(s.now()); err != nil {
		return nil, exception.InternalError(err)
	}
	return schedules, nil
}

func (s *schedule) Count(ctx context.Context, _ *ListSchedulesParams) (int64, error) {
	total, err := s.db.Count(ctx, s.db.DB, &entity.Schedule{}, nil)
	return total, exception.InternalError(err)
}

func (s *schedule) Get(ctx context.Context, scheduleID string, fields ...string) (*entity.Schedule, error) {
	schedule, err := s.get(ctx, s.db.DB, scheduleID, fields...)
	return schedule, exception.InternalError(err)
}

func (s *schedule) Create(ctx context.Context, schedule *entity.Schedule) error {
	err := s.db.Transaction(ctx, func(tx *gorm.DB) error {
		if err := schedule.FillJSON(); err != nil {
			return err
		}

		now := s.now()
		schedule.CreatedAt, schedule.UpdatedAt = now, now

		return tx.WithContext(ctx).Table(scheduleTable).Create(&schedule).Error
	})
	return exception.InternalError(err)
}

func (s *schedule) UpdateThumbnails(ctx context.Context, scheduleID string, thumbnails common.Images) error {
	err := s.db.Transaction(ctx, func(tx *gorm.DB) error {
		schedule, err := s.get(ctx, tx, scheduleID, "thumbnail_url")
		if err != nil {
			return err
		}
		if schedule.ThumbnailURL == "" {
			return fmt.Errorf("database: thumbnail url is empty: %w", exception.ErrFailedPrecondition)
		}

		buf, err := thumbnails.Marshal()
		if err != nil {
			return err
		}
		params := map[string]interface{}{
			"thumbnails": datatypes.JSON(buf),
			"updated_at": s.now(),
		}

		err = tx.WithContext(ctx).
			Table(scheduleTable).
			Where("id = ?", scheduleID).
			Updates(params).Error
		return err
	})
	return exception.InternalError(err)
}

func (s *schedule) get(ctx context.Context, tx *gorm.DB, scheduleID string, fields ...string) (*entity.Schedule, error) {
	var schedule *entity.Schedule

	err := s.db.Statement(ctx, tx, scheduleTable, fields...).
		Where("id = ?", scheduleID).
		First(&schedule).Error
	if err != nil {
		return nil, err
	}
	if err := schedule.Fill(s.now()); err != nil {
		return nil, err
	}
	return schedule, nil
}
