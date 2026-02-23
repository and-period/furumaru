package entity

import (
	"iter"

	"github.com/and-period/furumaru/api/pkg/collection"
)

// All はインデックスと注文付加情報のペアを返すイテレーターを返す。
func (ms MultiOrderMetadata) All() iter.Seq2[int, *OrderMetadata] {
	return func(yield func(int, *OrderMetadata) bool) {
		for i, m := range ms {
			if !yield(i, m) {
				return
			}
		}
	}
}

// IterMapByOrderID は注文履歴IDをキー、注文付加情報を値とするイテレーターを返す。
func (ms MultiOrderMetadata) IterMapByOrderID() iter.Seq2[string, *OrderMetadata] {
	return collection.MapIter(ms, func(m *OrderMetadata) (string, *OrderMetadata) {
		return m.OrderID, m
	})
}
