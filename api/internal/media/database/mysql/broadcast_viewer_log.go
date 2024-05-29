package mysql

import (
	"context"
	"fmt"
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

func (l *broadcastViewerLog) Aggregate(
	ctx context.Context, params *database.AggregateBroadcastViewerLogsParams,
) (entity.AggregatedBroadcastViewerLogs, error) {
	var logs entity.AggregatedBroadcastViewerLogs

	fields := []string{
		"broadcast_id",
		fmt.Sprintf("DATE_FORMAT(created_at, '%s') AS timestamp", params.Interval),
		"count(DISTINCT(user_id)) AS total",
	}
	stmt := l.db.Statement(ctx, l.db.DB, broadcastViewerLogTable, fields...).
		Where("broadcast_id = ?", params.BroadcastID).
		Where("user_agent NOT IN (?)", entity.ExcludeUserAgentLogs)
	if !params.CreatedAtGte.IsZero() {
		stmt = stmt.Where("created_at >= ?", params.CreatedAtGte)
	}
	if !params.CreatedAtLt.IsZero() {
		stmt = stmt.Where("created_at < ?", params.CreatedAtLt)
	}
	stmt = stmt.Group("broadcast_id, timestamp").Order("timestamp ASC")

	err := stmt.Find(&logs).Error
	return logs, dbError(err)
}
