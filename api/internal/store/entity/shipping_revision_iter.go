package entity

import (
	"iter"

	"github.com/and-period/furumaru/api/pkg/collection"
)

// All はインデックスと配送設定変更履歴のペアを返すイテレーターを返す。
func (rs ShippingRevisions) All() iter.Seq2[int, *ShippingRevision] {
	return func(yield func(int, *ShippingRevision) bool) {
		for i, r := range rs {
			if !yield(i, r) {
				return
			}
		}
	}
}

// IterMapByShippingID は配送設定IDをキー、配送設定変更履歴を値とするイテレーターを返す。
func (rs ShippingRevisions) IterMapByShippingID() iter.Seq2[string, *ShippingRevision] {
	return collection.MapIter(rs, func(r *ShippingRevision) (string, *ShippingRevision) {
		return r.ShippingID, r
	})
}
