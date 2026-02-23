package entity

import (
	"iter"

	"github.com/and-period/furumaru/api/pkg/collection"
)

// All はインデックスと注文支払い情報のペアを返すイテレーターを返す。
func (ps OrderPayments) All() iter.Seq2[int, *OrderPayment] {
	return func(yield func(int, *OrderPayment) bool) {
		for i, p := range ps {
			if !yield(i, p) {
				return
			}
		}
	}
}

// IterMapByOrderID は注文履歴IDをキー、注文支払い情報を値とするイテレーターを返す。
func (ps OrderPayments) IterMapByOrderID() iter.Seq2[string, *OrderPayment] {
	return collection.MapIter(ps, func(p *OrderPayment) (string, *OrderPayment) {
		return p.OrderID, p
	})
}
