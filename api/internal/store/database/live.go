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

func (l *live) Create(ctx context.Context, live *entity.Live) error {
	_, err := l.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if err := live.FillJSON(); err != nil {
			return nil, err
		}

		now := l.now()
		live.CreatedAt, live.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(liveTable).Create(&live).Error
		return nil, err
	})
	return exception.InternalError(err)
}
