package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/datatypes"
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
	if p.CoordinatorID != "" {
		stmt = stmt.Where("coordinator_id = ?", p.CoordinatorID)
	}
	if p.ProducerID != "" {
		stmt = stmt.Where("EXISTS (SELECT * FROM lives WHERE lives.schedule_id = schedules.id AND lives.producer_id = ?)", p.ProducerID)
	}
	if !p.StartAtGte.IsZero() {
		stmt = stmt.Where("start_at >= ?", p.StartAtGte)
	}
	if !p.StartAtLt.IsZero() {
		stmt = stmt.Where("start_at < ?", p.StartAtLt)
	}
	if !p.EndAtGte.IsZero() {
		stmt = stmt.Where("end_at >= ?", p.EndAtGte)
	}
	if !p.EndAtLt.IsZero() {
		stmt = stmt.Where("end_at < ?", p.EndAtLt)
	}
	if p.OnlyPublished {
		stmt = stmt.Where("public = ?", true).Where("approved = ?", true)
	}
	return stmt
}

func (p listSchedulesParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (s *schedule) List(ctx context.Context, params *database.ListSchedulesParams, fields ...string) (entity.Schedules, error) {
	var schedules entity.Schedules

	p := listSchedulesParams(*params)

	stmt := s.db.Statement(ctx, s.db.DB, scheduleTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&schedules).Error; err != nil {
		return nil, dbError(err)
	}
	if err := schedules.Fill(s.now()); err != nil {
		return nil, dbError(err)
	}
	return schedules, nil
}

func (s *schedule) Count(ctx context.Context, params *database.ListSchedulesParams) (int64, error) {
	p := listSchedulesParams(*params)

	total, err := s.db.Count(ctx, s.db.DB, &entity.Schedule{}, p.stmt)
	return total, dbError(err)
}

func (s *schedule) MultiGet(ctx context.Context, scheduleIDs []string, fields ...string) (entity.Schedules, error) {
	var schedules entity.Schedules

	stmt := s.db.Statement(ctx, s.db.DB, scheduleTable, fields...).
		Where("id IN (?)", scheduleIDs)

	if err := stmt.Find(&schedules).Error; err != nil {
		return nil, dbError(err)
	}
	if err := schedules.Fill(s.now()); err != nil {
		return nil, dbError(err)
	}
	return schedules, nil
}

func (s *schedule) Get(ctx context.Context, scheduleID string, fields ...string) (*entity.Schedule, error) {
	schedule, err := s.get(ctx, s.db.DB, scheduleID, fields...)
	return schedule, dbError(err)
}

func (s *schedule) Create(ctx context.Context, schedule *entity.Schedule) error {
	if err := schedule.FillJSON(); err != nil {
		return err
	}

	now := s.now()
	schedule.CreatedAt, schedule.UpdatedAt = now, now

	err := s.db.DB.WithContext(ctx).Table(scheduleTable).Create(&schedule).Error
	return dbError(err)
}

func (s *schedule) Update(ctx context.Context, scheduleID string, params *database.UpdateScheduleParams) error {
	err := s.db.Transaction(ctx, func(tx *gorm.DB) error {
		schedule, err := s.get(ctx, tx, scheduleID, "start_at")
		if err != nil {
			return err
		}
		if s.now().After(schedule.StartAt) {
			return fmt.Errorf("database: this schedule has already started: %w", database.ErrFailedPrecondition)
		}

		update := map[string]interface{}{
			"title":             params.Title,
			"description":       params.Description,
			"thumbnail_url":     params.ThumbnailURL,
			"image_url":         params.ImageURL,
			"opening_video_url": params.OpeningVideoURL,
			"start_at":          params.StartAt,
			"end_at":            params.EndAt,
			"updated_at":        s.now(),
		}

		err = s.db.DB.WithContext(ctx).
			Table(scheduleTable).
			Where("id = ?", scheduleID).
			Updates(update).Error
		return err
	})
	return dbError(err)
}

func (s *schedule) UpdateThumbnails(ctx context.Context, scheduleID string, thumbnails common.Images) error {
	err := s.db.Transaction(ctx, func(tx *gorm.DB) error {
		schedule, err := s.get(ctx, tx, scheduleID, "thumbnail_url")
		if err != nil {
			return err
		}
		if schedule.ThumbnailURL == "" {
			return fmt.Errorf("database: thumbnail url is empty: %w", database.ErrFailedPrecondition)
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
	return dbError(err)
}

func (s *schedule) Approve(ctx context.Context, scheduleID string, params *database.ApproveScheduleParams) error {
	var approvedAdminID *string
	if params.Approved {
		approvedAdminID = &params.ApprovedAdminID
	}
	update := map[string]interface{}{
		"approved":          params.Approved,
		"approved_admin_id": approvedAdminID,
		"updated_at":        s.now(),
	}

	err := s.db.DB.WithContext(ctx).
		Table(scheduleTable).
		Where("id = ?", scheduleID).
		Updates(update).Error
	return dbError(err)
}

func (s *schedule) Publish(ctx context.Context, scheduleID string, public bool) error {
	update := map[string]interface{}{
		"public":     public,
		"updated_at": s.now(),
	}
	err := s.db.DB.WithContext(ctx).
		Table(scheduleTable).
		Where("id = ?", scheduleID).
		Updates(update).Error
	return dbError(err)
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
