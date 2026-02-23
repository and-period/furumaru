package entity

import "iter"

// All はインデックスと商品タグのペアを返すイテレーターを返す。
func (ts ProductTags) All() iter.Seq2[int, *ProductTag] {
	return func(yield func(int, *ProductTag) bool) {
		for i, t := range ts {
			if !yield(i, t) {
				return
			}
		}
	}
}

// IterMap は商品タグIDをキー、商品タグを値とするイテレーターを返す。
func (ts ProductTags) IterMap() iter.Seq2[string, *ProductTag] {
	return MapIter(ts, func(t *ProductTag) (string, *ProductTag) {
		return t.ID, t
	})
}
