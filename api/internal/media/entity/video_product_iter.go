package entity

import "iter"

// All はインデックスと要素のペアを返すイテレーターを返す。
func (ps VideoProducts) All() iter.Seq2[int, *VideoProduct] {
	return func(yield func(int, *VideoProduct) bool) {
		for i, p := range ps {
			if !yield(i, p) {
				return
			}
		}
	}
}

// IterProductIDs は商品IDを返すイテレーターを返す。
func (ps VideoProducts) IterProductIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, p := range ps {
			if !yield(p.ProductID) {
				return
			}
		}
	}
}

// IterGroupByVideoID はビデオIDをキー、VideoProductsを値とするイテレーターを返す。
func (ps VideoProducts) IterGroupByVideoID() iter.Seq2[string, VideoProducts] {
	return func(yield func(string, VideoProducts) bool) {
		groups := ps.GroupByVideoID()
		for k, v := range groups {
			if !yield(k, v) {
				return
			}
		}
	}
}
