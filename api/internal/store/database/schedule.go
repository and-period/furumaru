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

func (s *schedule) Get(ctx context.Context, scheduleID string, fields ...string) (*entity.Schedule, error) {
	schedule, err := s.get(ctx, s.db.DB, scheduleID, fields...)
	return schedule, exception.InternalError(err)
}

func (s *schedule) Create(
	ctx context.Context, schedule *entity.Schedule, lives entity.Lives, products entity.LiveProducts,
) error {
	err := s.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := s.now()
		schedule.CreatedAt, schedule.UpdatedAt = now, now
		err := tx.WithContext(ctx).Table(scheduleTable).Create(&schedule).Error
		if err != nil {
			return err
		}
		for i := range lives {
			lives[i].CreatedAt, lives[i].UpdatedAt = now, now
		}
		if err := tx.WithContext(ctx).Table(liveTable).Create(&lives).Error; err != nil {
			return err
		}
		for i := range products {
			products[i].CreatedAt, products[i].UpdatedAt = now, now
		}
		if err := tx.WithContext(ctx).Table(liveProductTable).Create(&products).Error; err != nil {
			return err
		}
		lives.Fill(products.GroupByLiveID())
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

	return schedule, nil
}
