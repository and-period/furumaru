package service

import (
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

// VideoViewerLogInterval - ライブ配信視聴ログ取得間隔
type VideoViewerLogInterval types.VideoViewerLogInterval

type VideoViewerLog struct {
	types.VideoViewerLog
}

type VideoViewerLogs []*VideoViewerLog

func NewVideoViewerLogIntervalFromRequest(interval string) VideoViewerLogInterval {
	return VideoViewerLogInterval(interval)
}

func (i VideoViewerLogInterval) Duration() time.Duration {
	switch types.VideoViewerLogInterval(i) {
	case types.VideoViewerLogIntervalSecond:
		return time.Second
	case types.VideoViewerLogIntervalMinute:
		return time.Minute
	case types.VideoViewerLogIntervalHour:
		return time.Hour
	default:
		return 0
	}
}

func (i VideoViewerLogInterval) MediaEntity() entity.AggregateVideoViewerLogInterval {
	switch types.VideoViewerLogInterval(i) {
	case types.VideoViewerLogIntervalSecond:
		return entity.AggregateVideoViewerLogIntervalSecond
	case types.VideoViewerLogIntervalMinute:
		return entity.AggregateVideoViewerLogIntervalMinute
	case types.VideoViewerLogIntervalHour:
		return entity.AggregateVideoViewerLogIntervalHour
	default:
		return ""
	}
}

func NewVideoViewerLog(aggregate *entity.AggregatedVideoViewerLog, interval time.Duration) *VideoViewerLog {
	return &VideoViewerLog{
		VideoViewerLog: types.VideoViewerLog{
			VideoID: aggregate.VideoID,
			StartAt: aggregate.ReportedAt.Unix(),
			EndAt:   aggregate.ReportedAt.Add(interval).Unix(),
			Total:   aggregate.Total,
		},
	}
}

func newEmptyVideoViewerLog(videoID string, startAt time.Time, interval time.Duration) *VideoViewerLog {
	aggregate := &entity.AggregatedVideoViewerLog{
		VideoID:    videoID,
		ReportedAt: startAt,
		Total:      0,
	}
	return NewVideoViewerLog(aggregate, interval)
}

func (l *VideoViewerLog) Response() *types.VideoViewerLog {
	return &l.VideoViewerLog
}

func NewVideoViewerLogs(
	interval VideoViewerLogInterval,
	startAt, endAt time.Time,
	aggregates entity.AggregatedVideoViewerLogs,
) VideoViewerLogs {
	duration := interval.Duration()
	if duration == 0 {
		return VideoViewerLogs{}
	}
	aggregatesMap := aggregates.GroupByVideoID()
	res := make(VideoViewerLogs, 0, len(aggregates)*int((endAt.Sub(startAt)/duration)+1))
	for videoID, aggregates := range aggregatesMap {
		aggregateMap := aggregates.MapByReportedAt()
		for ts := startAt; ts.Before(endAt); ts = ts.Add(duration) {
			if aggregate, ok := aggregateMap[ts]; ok {
				res = append(res, NewVideoViewerLog(aggregate, duration))
				continue
			}
			res = append(res, newEmptyVideoViewerLog(videoID, ts, duration))
		}
	}
	return res
}

func (ls VideoViewerLogs) Response() []*types.VideoViewerLog {
	res := make([]*types.VideoViewerLog, len(ls))
	for i := range ls {
		res[i] = ls[i].Response()
	}
	return res
}
