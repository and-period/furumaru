package entity

import "iter"

// All はインデックスと体験種別のペアを返すイテレーターを返す。
func (ts ExperienceTypes) All() iter.Seq2[int, *ExperienceType] {
	return func(yield func(int, *ExperienceType) bool) {
		for i, t := range ts {
			if !yield(i, t) {
				return
			}
		}
	}
}

// IterMap は体験種別IDをキー、体験種別を値とするイテレーターを返す。
func (ts ExperienceTypes) IterMap() iter.Seq2[string, *ExperienceType] {
	return MapIter(ts, func(t *ExperienceType) (string, *ExperienceType) {
		return t.ID, t
	})
}
