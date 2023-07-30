package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
)

const broadcastTable = "broadcasts"

type broadcast struct {
	db  *database.Client
	now func() time.Time
}

func NewBroadcast(db *database.Client) Broadcast {
	return &broadcast{
		db:  db,
		now: jst.Now,
	}
}

func (b *broadcast) GetByScheduleID(
	ctx context.Context, scheduleID string, fields ...string,
) (*entity.Broadcast, error) {
	var broadcast *entity.Broadcast

	stmt := b.db.Statement(ctx, b.db.DB, broadcastTable, fields...).
		Where("schedule_id = ?", scheduleID)

	if err := stmt.First(&broadcast).Error; err != nil {
		return nil, err
	}
	return broadcast, nil
}
