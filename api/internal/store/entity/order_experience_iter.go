package entity

import "iter"

// All はインデックスと注文体験情報のペアを返すイテレーターを返す。
func (os OrderExperiences) All() iter.Seq2[int, *OrderExperience] {
	return func(yield func(int, *OrderExperience) bool) {
		for i, o := range os {
			if !yield(i, o) {
				return
			}
		}
	}
}

// IterMapByOrderID は注文履歴IDをキー、注文体験情報を値とするイテレーターを返す。
func (os OrderExperiences) IterMapByOrderID() iter.Seq2[string, *OrderExperience] {
	return MapIter(os, func(o *OrderExperience) (string, *OrderExperience) {
		return o.OrderID, o
	})
}
