package entity

import "iter"

// All はインデックスと注文配送情報のペアを返すイテレーターを返す。
func (fs OrderFulfillments) All() iter.Seq2[int, *OrderFulfillment] {
	return func(yield func(int, *OrderFulfillment) bool) {
		for i, f := range fs {
			if !yield(i, f) {
				return
			}
		}
	}
}

// IterGroupByOrderID は注文履歴IDをキー、注文配送情報一覧を値とするイテレーターを返す。
func (fs OrderFulfillments) IterGroupByOrderID() iter.Seq2[string, OrderFulfillments] {
	return func(yield func(string, OrderFulfillments) bool) {
		groups := fs.GroupByOrderID()
		for k, v := range groups {
			if !yield(k, v) {
				return
			}
		}
	}
}
