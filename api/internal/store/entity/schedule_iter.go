package entity

import "iter"

// All はインデックスと開催スケジュールのペアを返すイテレーターを返す。
func (ss Schedules) All() iter.Seq2[int, *Schedule] {
	return func(yield func(int, *Schedule) bool) {
		for i, s := range ss {
			if !yield(i, s) {
				return
			}
		}
	}
}

// IterMap はスケジュールIDをキー、スケジュールを値とするイテレーターを返す。
func (ss Schedules) IterMap() iter.Seq2[string, *Schedule] {
	return MapIter(ss, func(s *Schedule) (string, *Schedule) {
		return s.ID, s
	})
}
