package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/gorm"
)

const scheduleTable = "schedules"

var scheduleFields = []string{
	"id", "title", "description", "thumbnail_url",
	"start_at", "end_at", "canceled",
	"created_at", "updated_at", "deleted_at",
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

func (s *schedule) Create(ctx context.Context, schedule *entity.Schedule) error {
	_, err := s.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := s.now()
		schedule.CreatedAt, schedule.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(scheduleTable).Create(&schedule).Error
		return nil, err
	})
	return exception.InternalError(err)
}
