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

const liveTable = "lives"

type live struct {
	db  *database.Client
	now func() time.Time
}

func NewLive(db *database.Client) Live {
	return &live{
		db:  db,
		now: jst.Now,
	}
}

func (l *live) Get(ctx context.Context, liveID string, fields ...string) (*entity.Live, error) {
	live, err := l.get(ctx, l.db.DB, liveID, fields...)
	return live, exception.InternalError(err)
}

func (l *live) Update(ctx context.Context, liveID string, params *UpdateLiveParams) error {
	err := l.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := l.now()
		products := params.LiveProducts
		if _, err := l.get(ctx, tx, liveID); err != nil {
			return err
		}
		updates := map[string]interface{}{
			"producer_id": params.ProducerID,
			"title":       params.Title,
			"description": params.Description,
			"published":   params.Published,
			"canceled":    params.Canceled,
			"start_at":    params.StartAt,
			"end_at":      params.EndAt,
			"updated_at":  now,
		}
		err := tx.WithContext(ctx).
			Table(liveProductTable).
			Where("live_id = ?", liveID).
			Delete(&entity.LiveProduct{}).Error
		if err != nil {
			return err
		}
		for i := range products {
			products[i].CreatedAt, products[i].UpdatedAt = now, now
			err = tx.WithContext(ctx).Table(liveProductTable).Create(&products[i]).Error
			if err != nil {
				return err
			}
		}
		err = tx.WithContext(ctx).
			Table(liveTable).
			Where("id = ?", liveID).
			Updates(updates).Error
		return err
	})
	return exception.InternalError(err)
}

func (l *live) get(ctx context.Context, tx *gorm.DB, liveID string, fields ...string) (*entity.Live, error) {
	var (
		live         *entity.Live
		liveProducts entity.LiveProducts
	)

	err := l.db.Statement(ctx, tx, liveTable, fields...).
		Where("id = ?", liveID).
		First(&live).Error
	if err != nil {
		return nil, err
	}
	err = l.db.Statement(ctx, tx, liveProductTable).
		Where("live_id IN (?)", liveID).
		Find(&liveProducts).Error
	if err != nil {
		return nil, err
	}

	live.Fill(liveProducts, jst.Now())
	return live, nil
}
