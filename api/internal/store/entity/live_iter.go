package entity

import (
	"iter"

	"github.com/and-period/furumaru/api/pkg/collection"
)

// All はインデックスとライブ配信のペアを返すイテレーターを返す。
func (ls Lives) All() iter.Seq2[int, *Live] {
	return func(yield func(int, *Live) bool) {
		for i, l := range ls {
			if !yield(i, l) {
				return
			}
		}
	}
}

// IterMap はライブ配信IDをキー、ライブ配信を値とするイテレーターを返す。
func (ls Lives) IterMap() iter.Seq2[string, *Live] {
	return collection.MapIter(ls, func(l *Live) (string, *Live) {
		return l.ID, l
	})
}

// IterGroupByScheduleID は開催スケジュールIDをキー、ライブ配信一覧を値とするイテレーターを返す。
func (ls Lives) IterGroupByScheduleID() iter.Seq2[string, Lives] {
	return func(yield func(string, Lives) bool) {
		groups := ls.GroupByScheduleID()
		for k, v := range groups {
			if !yield(k, v) {
				return
			}
		}
	}
}
