package service

import (
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

// BroadcastViewerLogInterval - ライブ配信視聴ログ取得間隔
type BroadcastViewerLogInterval string

const (
	BroadcastViewerLogIntervalSecond BroadcastViewerLogInterval = "second"
	BroadcastViewerLogIntervalMinute BroadcastViewerLogInterval = "minute"
	BroadcastViewerLogIntervalHour   BroadcastViewerLogInterval = "hour"
)

type BroadcastViewerLog struct {
	response.BroadcastViewerLog
}

type BroadcastViewerLogs []*BroadcastViewerLog

func NewBroadcastViewerLogIntervalFromRequest(interval string) BroadcastViewerLogInterval {
	return BroadcastViewerLogInterval(interval)
}

func (i BroadcastViewerLogInterval) Duration() time.Duration {
	switch i {
	case BroadcastViewerLogIntervalSecond:
		return time.Second
	case BroadcastViewerLogIntervalMinute:
		return time.Minute
	case BroadcastViewerLogIntervalHour:
		return time.Hour
	default:
		return 0
	}
}

func (i BroadcastViewerLogInterval) MediaEntity() entity.AggregateBroadcastViewerLogInterval {
	switch i {
	case BroadcastViewerLogIntervalSecond:
		return entity.AggregateBroadcastViewerLogIntervalSecond
	case BroadcastViewerLogIntervalMinute:
		return entity.AggregateBroadcastViewerLogIntervalMinute
	case BroadcastViewerLogIntervalHour:
		return entity.AggregateBroadcastViewerLogIntervalHour
	default:
		return ""
	}
}

func NewBroadcastViewerLog(aggregate *entity.AggregatedBroadcastViewerLog, interval time.Duration) *BroadcastViewerLog {
	return &BroadcastViewerLog{
		BroadcastViewerLog: response.BroadcastViewerLog{
			BroadcastID: aggregate.BroadcastID,
			StartAt:     aggregate.ReportedAt.Unix(),
			EndAt:       aggregate.ReportedAt.Add(interval).Unix(),
			Total:       aggregate.Total,
		},
	}
}

func newEmptyBroadcastViewerLog(broadcastID string, startAt time.Time, interval time.Duration) *BroadcastViewerLog {
	aggregate := &entity.AggregatedBroadcastViewerLog{
		BroadcastID: broadcastID,
		ReportedAt:  startAt,
		Total:       0,
	}
	return NewBroadcastViewerLog(aggregate, interval)
}

func (l *BroadcastViewerLog) Response() *response.BroadcastViewerLog {
	return &l.BroadcastViewerLog
}

func NewBroadcastViewerLogs(
	interval BroadcastViewerLogInterval,
	startAt, endAt time.Time,
	aggregates entity.AggregatedBroadcastViewerLogs,
) BroadcastViewerLogs {
	duration := interval.Duration()
	if duration == 0 {
		return BroadcastViewerLogs{}
	}
	aggregatesMap := aggregates.GroupByBroadcastID()
	res := make(BroadcastViewerLogs, 0, len(aggregates)*int((endAt.Sub(startAt)/duration)+1))
	for broadcastID, aggregates := range aggregatesMap {
		aggregateMap := aggregates.MapByReportedAt()
		for ts := startAt; ts.Before(endAt); ts = ts.Add(duration) {
			if aggregate, ok := aggregateMap[ts]; ok {
				res = append(res, NewBroadcastViewerLog(aggregate, duration))
				continue
			}
			res = append(res, newEmptyBroadcastViewerLog(broadcastID, ts, duration))
		}
	}
	return res
}

func (ls BroadcastViewerLogs) Response() []*response.BroadcastViewerLog {
	res := make([]*response.BroadcastViewerLog, len(ls))
	for i := range ls {
		res[i] = ls[i].Response()
	}
	return res
}
