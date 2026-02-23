package entity

import "iter"

// All はインデックスとスケジュールのペアを返すイテレーターを返す。
func (ss Schedules) All() iter.Seq2[int, *Schedule] {
	return func(yield func(int, *Schedule) bool) {
		for i, s := range ss {
			if !yield(i, s) {
				return
			}
		}
	}
}

// IterMap はメッセージIDをキー、スケジュールを値とするイテレーターを返す。
func (ss Schedules) IterMap() iter.Seq2[string, *Schedule] {
	return MapIter(ss, func(s *Schedule) (string, *Schedule) {
		return s.MessageID, s
	})
}
