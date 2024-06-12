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

func (l *broadcastViewerLog) GetTotal(ctx context.Context, params *database.GetBroadcastTotalViewersParams) (int64, error) {
	var total int64

	const field = "COUNT(DISTINCT(session_id)) AS total"
	stmt := l.db.Statement(ctx, l.db.DB, broadcastViewerLogTable, field).
		Where("broadcast_id = ?", params.BroadcastID).
		Where("user_agent NOT IN (?)", entity.ExcludeUserAgentLogs)
	if !params.CreatedAtGte.IsZero() {
		stmt = stmt.Where("created_at >= ?", params.CreatedAtGte)
	}
	if !params.CreatedAtLt.IsZero() {
		stmt = stmt.Where("created_at < ?", params.CreatedAtLt)
	}

	err := stmt.Count(&total).Error
	return total, dbError(err)
}

func (l *broadcastViewerLog) Aggregate(
	ctx context.Context, params *database.AggregateBroadcastViewerLogsParams,
) (entity.AggregatedBroadcastViewerLogs, error) {
	var logs internalAggregatedBroadcastViewerLogs

	fields := []string{
		"broadcast_id",
		fmt.Sprintf("DATE_FORMAT(created_at, '%s') AS reported_at", params.Interval),
		"COUNT(DISTINCT(session_id)) AS total",
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
	stmt = stmt.Group("broadcast_id, reported_at").Order("reported_at ASC")

	if err := stmt.Scan(&logs).Error; err != nil {
		return nil, dbError(err)
	}
	return logs.Entities(), nil
}

type internalAggregatedBroadcastViewerLog struct {
	BroadcastID string
	ReportedAt  string
	Total       int64
}

type internalAggregatedBroadcastViewerLogs []*internalAggregatedBroadcastViewerLog

func (l *internalAggregatedBroadcastViewerLog) Entity() *entity.AggregatedBroadcastViewerLog {
	reportedAt, _ := jst.Parse("2006-01-02 15:04:05", l.ReportedAt)
	return &entity.AggregatedBroadcastViewerLog{
		BroadcastID: l.BroadcastID,
		ReportedAt:  reportedAt,
		Total:       l.Total,
	}
}

func (ls internalAggregatedBroadcastViewerLogs) Entities() entity.AggregatedBroadcastViewerLogs {
	res := make(entity.AggregatedBroadcastViewerLogs, len(ls))
	for i := range ls {
		res[i] = ls[i].Entity()
	}
	return res
}
