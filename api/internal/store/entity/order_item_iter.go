package entity

import "iter"

// All はインデックスと注文商品情報のペアを返すイテレーターを返す。
func (is OrderItems) All() iter.Seq2[int, *OrderItem] {
	return func(yield func(int, *OrderItem) bool) {
		for i, item := range is {
			if !yield(i, item) {
				return
			}
		}
	}
}

// IterGroupByFulfillmentID は注文配送IDをキー、注文商品情報一覧を値とするイテレーターを返す。
func (is OrderItems) IterGroupByFulfillmentID() iter.Seq2[string, OrderItems] {
	return func(yield func(string, OrderItems) bool) {
		groups := is.GroupByFulfillmentID()
		for k, v := range groups {
			if !yield(k, v) {
				return
			}
		}
	}
}

// IterGroupByOrderID は注文履歴IDをキー、注文商品情報一覧を値とするイテレーターを返す。
func (is OrderItems) IterGroupByOrderID() iter.Seq2[string, OrderItems] {
	return func(yield func(string, OrderItems) bool) {
		groups := is.GroupByOrderID()
		for k, v := range groups {
			if !yield(k, v) {
				return
			}
		}
	}
}
