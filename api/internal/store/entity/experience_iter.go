package entity

import "iter"

// All はインデックスと体験のペアを返すイテレーターを返す。
func (es Experiences) All() iter.Seq2[int, *Experience] {
	return func(yield func(int, *Experience) bool) {
		for i, e := range es {
			if !yield(i, e) {
				return
			}
		}
	}
}

// IterMap は体験IDをキー、体験を値とするイテレーターを返す。
func (es Experiences) IterMap() iter.Seq2[string, *Experience] {
	return MapIter(es, func(e *Experience) (string, *Experience) {
		return e.ID, e
	})
}

// IterFilterByPublished は公開中の体験を返すイテレーターを返す。
// 非公開・アーカイブ済みの体験は除外される。
func (es Experiences) IterFilterByPublished() iter.Seq[*Experience] {
	return FilterIter(es, func(e *Experience) bool {
		return e.Status != ExperienceStatusPrivate && e.Status != ExperienceStatusArchived
	})
}
