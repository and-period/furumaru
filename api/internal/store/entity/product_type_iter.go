package entity

import "iter"

// All はインデックスと品目のペアを返すイテレーターを返す。
func (ts ProductTypes) All() iter.Seq2[int, *ProductType] {
	return func(yield func(int, *ProductType) bool) {
		for i, t := range ts {
			if !yield(i, t) {
				return
			}
		}
	}
}

// IterMap は品目IDをキー、品目を値とするイテレーターを返す。
func (ts ProductTypes) IterMap() iter.Seq2[string, *ProductType] {
	return MapIter(ts, func(t *ProductType) (string, *ProductType) {
		return t.ID, t
	})
}
