package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const broadcastViewerLogTable = "broadcast_viewer_logs"

type broadcastViewerLog struct {
	db  *mysql.Client
	now func() time.Time
}

func newBroadcastViewerLog(db *mysql.Client) database.BroadcastViewerLog {
	return &broadcastViewerLog{
		db:  db,
		now: jst.Now,
	}
}

func (l *broadcastViewerLog) Create(ctx context.Context, log *entity.BroadcastViewerLog) error {
	now := l.now()
	log.CreatedAt, log.UpdatedAt = now, now

	err := l.db.DB.WithContext(ctx).Table(broadcastViewerLogTable).Create(&log).Error
	return dbError(err)
}
