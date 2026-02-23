package entity

import "iter"

// All はインデックスとスポットのペアを返すイテレーターを返す。
func (cs Spots) All() iter.Seq2[int, *Spot] {
	return func(yield func(int, *Spot) bool) {
		for i, c := range cs {
			if !yield(i, c) {
				return
			}
		}
	}
}

// IterGroupByUserType は投稿者種別をキー、スポット一覧を値とするイテレーターを返す。
func (cs Spots) IterGroupByUserType() iter.Seq2[SpotUserType, Spots] {
	return func(yield func(SpotUserType, Spots) bool) {
		groups := cs.GroupByUserType()
		for k, v := range groups {
			if !yield(k, v) {
				return
			}
		}
	}
}
