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
