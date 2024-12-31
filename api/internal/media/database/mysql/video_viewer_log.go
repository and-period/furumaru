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

const videoViewerLogTable = "video_viewer_logs"

type videoViewerLog struct {
	db  *mysql.Client
	now func() time.Time
}

func NewVideoViewerLog(db *mysql.Client) database.VideoViewerLog {
	return &videoViewerLog{
		db:  db,
		now: jst.Now,
	}
}

func (l *videoViewerLog) Create(ctx context.Context, log *entity.VideoViewerLog) error {
	now := l.now()
	log.CreatedAt, log.UpdatedAt = now, now

	err := l.db.DB.WithContext(ctx).Table(videoViewerLogTable).Create(&log).Error
	return dbError(err)
}

func (l *videoViewerLog) GetTotal(ctx context.Context, params *database.GetVideoTotalViewersParams) (int64, error) {
	var total int64

	const field = "COUNT(DISTINCT(session_id)) AS total"
	stmt := l.db.Statement(ctx, l.db.DB, videoViewerLogTable, field).
		Where("video_id = ?", params.VideoID).
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

func (l *videoViewerLog) Aggregate(
	ctx context.Context, params *database.AggregateVideoViewerLogsParams,
) (entity.AggregatedVideoViewerLogs, error) {
	var logs internalAggregatedVideoViewerLogs

	fields := []string{
		"video_id",
		fmt.Sprintf("DATE_FORMAT(created_at, '%s') AS reported_at", params.Interval),
		"COUNT(DISTINCT(session_id)) AS total",
	}
	stmt := l.db.Statement(ctx, l.db.DB, videoViewerLogTable, fields...).
		Where("video_id = ?", params.VideoID).
		Where("user_agent NOT IN (?)", entity.ExcludeUserAgentLogs)
	if !params.CreatedAtGte.IsZero() {
		stmt = stmt.Where("created_at >= ?", params.CreatedAtGte)
	}
	if !params.CreatedAtLt.IsZero() {
		stmt = stmt.Where("created_at < ?", params.CreatedAtLt)
	}
	stmt = stmt.Group("video_id, reported_at").Order("reported_at ASC")

	if err := stmt.Scan(&logs).Error; err != nil {
		return nil, dbError(err)
	}
	return logs.Entities(), nil
}

type internalAggregatedVideoViewerLog struct {
	VideoID    string
	ReportedAt string
	Total      int64
}

type internalAggregatedVideoViewerLogs []*internalAggregatedVideoViewerLog

func (l *internalAggregatedVideoViewerLog) Entity() *entity.AggregatedVideoViewerLog {
	reportedAt, _ := jst.Parse("2006-01-02 15:04:05", l.ReportedAt)
	return &entity.AggregatedVideoViewerLog{
		VideoID:    l.VideoID,
		ReportedAt: reportedAt,
		Total:      l.Total,
	}
}

func (ls internalAggregatedVideoViewerLogs) Entities() entity.AggregatedVideoViewerLogs {
	res := make(entity.AggregatedVideoViewerLogs, len(ls))
	for i := range ls {
		res[i] = ls[i].Entity()
	}
	return res
}
