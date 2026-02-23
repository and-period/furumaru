package entity

import "iter"

// All はインデックスと配送設定のペアを返すイテレーターを返す。
func (ss Shippings) All() iter.Seq2[int, *Shipping] {
	return func(yield func(int, *Shipping) bool) {
		for i, s := range ss {
			if !yield(i, s) {
				return
			}
		}
	}
}

// IterMap は配送設定IDをキー、配送設定を値とするイテレーターを返す。
func (ss Shippings) IterMap() iter.Seq2[string, *Shipping] {
	return MapIter(ss, func(s *Shipping) (string, *Shipping) {
		return s.ID, s
	})
}
