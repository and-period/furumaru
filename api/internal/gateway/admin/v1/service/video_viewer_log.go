package service

import (
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

// VideoViewerLogInterval - ライブ配信視聴ログ取得間隔
type VideoViewerLogInterval string

const (
	VideoViewerLogIntervalSecond VideoViewerLogInterval = "second"
	VideoViewerLogIntervalMinute VideoViewerLogInterval = "minute"
	VideoViewerLogIntervalHour   VideoViewerLogInterval = "hour"
)

type VideoViewerLog struct {
	response.VideoViewerLog
}

type VideoViewerLogs []*VideoViewerLog

func NewVideoViewerLogIntervalFromRequest(interval string) VideoViewerLogInterval {
	return VideoViewerLogInterval(interval)
}

func (i VideoViewerLogInterval) Duration() time.Duration {
	switch i {
	case VideoViewerLogIntervalSecond:
		return time.Second
	case VideoViewerLogIntervalMinute:
		return time.Minute
	case VideoViewerLogIntervalHour:
		return time.Hour
	default:
		return 0
	}
}

func (i VideoViewerLogInterval) MediaEntity() entity.AggregateVideoViewerLogInterval {
	switch i {
	case VideoViewerLogIntervalSecond:
		return entity.AggregateVideoViewerLogIntervalSecond
	case VideoViewerLogIntervalMinute:
		return entity.AggregateVideoViewerLogIntervalMinute
	case VideoViewerLogIntervalHour:
		return entity.AggregateVideoViewerLogIntervalHour
	default:
		return ""
	}
}

func NewVideoViewerLog(aggregate *entity.AggregatedVideoViewerLog, interval time.Duration) *VideoViewerLog {
	return &VideoViewerLog{
		VideoViewerLog: response.VideoViewerLog{
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

func (l *VideoViewerLog) Response() *response.VideoViewerLog {
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

func (ls VideoViewerLogs) Response() []*response.VideoViewerLog {
	res := make([]*response.VideoViewerLog, len(ls))
	for i := range ls {
		res[i] = ls[i].Response()
	}
	return res
}
