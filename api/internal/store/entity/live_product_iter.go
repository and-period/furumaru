package entity

import "iter"

// All はインデックスとライブ配信関連商品のペアを返すイテレーターを返す。
func (ps LiveProducts) All() iter.Seq2[int, *LiveProduct] {
	return func(yield func(int, *LiveProduct) bool) {
		for i, p := range ps {
			if !yield(i, p) {
				return
			}
		}
	}
}

// IterGroupByLiveID はライブ配信IDをキー、ライブ配信関連商品一覧を値とするイテレーターを返す。
func (ps LiveProducts) IterGroupByLiveID() iter.Seq2[string, LiveProducts] {
	return func(yield func(string, LiveProducts) bool) {
		groups := ps.GroupByLiveID()
		for k, v := range groups {
			if !yield(k, v) {
				return
			}
		}
	}
}
