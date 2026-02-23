package entity

import (
	"iter"
	"time"

	"github.com/and-period/furumaru/api/pkg/collection"
)

// All はインデックスと要素のペアを返すイテレーターを返す。
func (ls AggregatedBroadcastViewerLogs) All() iter.Seq2[int, *AggregatedBroadcastViewerLog] {
	return func(yield func(int, *AggregatedBroadcastViewerLog) bool) {
		for i, l := range ls {
			if !yield(i, l) {
				return
			}
		}
	}
}

// IterMapByReportedAt は集計日時をキー、集計情報を値とするイテレーターを返す。
func (ls AggregatedBroadcastViewerLogs) IterMapByReportedAt() iter.Seq2[time.Time, *AggregatedBroadcastViewerLog] {
	return collection.MapIter(ls, func(l *AggregatedBroadcastViewerLog) (time.Time, *AggregatedBroadcastViewerLog) {
		return l.ReportedAt, l
	})
}

// IterGroupByBroadcastID はライブ配信IDをキー、AggregatedBroadcastViewerLogsを値とするイテレーターを返す。
func (ls AggregatedBroadcastViewerLogs) IterGroupByBroadcastID() iter.Seq2[string, AggregatedBroadcastViewerLogs] {
	return func(yield func(string, AggregatedBroadcastViewerLogs) bool) {
		groups := ls.GroupByBroadcastID()
		for k, v := range groups {
			if !yield(k, v) {
				return
			}
		}
	}
}
