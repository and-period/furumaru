package entity

import (
	"iter"
	"time"
)

// All はインデックスと要素のペアを返すイテレーターを返す。
func (ls AggregatedVideoViewerLogs) All() iter.Seq2[int, *AggregatedVideoViewerLog] {
	return func(yield func(int, *AggregatedVideoViewerLog) bool) {
		for i, l := range ls {
			if !yield(i, l) {
				return
			}
		}
	}
}

// IterMapByReportedAt は集計日時をキー、集計情報を値とするイテレーターを返す。
func (ls AggregatedVideoViewerLogs) IterMapByReportedAt() iter.Seq2[time.Time, *AggregatedVideoViewerLog] {
	return MapIter(ls, func(l *AggregatedVideoViewerLog) (time.Time, *AggregatedVideoViewerLog) {
		return l.ReportedAt, l
	})
}

// IterGroupByVideoID はビデオIDをキー、AggregatedVideoViewerLogsを値とするイテレーターを返す。
func (ls AggregatedVideoViewerLogs) IterGroupByVideoID() iter.Seq2[string, AggregatedVideoViewerLogs] {
	return func(yield func(string, AggregatedVideoViewerLogs) bool) {
		groups := ls.GroupByVideoID()
		for k, v := range groups {
			if !yield(k, v) {
				return
			}
		}
	}
}
