package entity

import "iter"

// All はインデックスと商品変更履歴のペアを返すイテレーターを返す。
func (rs ProductRevisions) All() iter.Seq2[int, *ProductRevision] {
	return func(yield func(int, *ProductRevision) bool) {
		for i, r := range rs {
			if !yield(i, r) {
				return
			}
		}
	}
}

// IterMapByProductID は商品IDをキー、商品変更履歴を値とするイテレーターを返す。
func (rs ProductRevisions) IterMapByProductID() iter.Seq2[string, *ProductRevision] {
	return MapIter(rs, func(r *ProductRevision) (string, *ProductRevision) {
		return r.ProductID, r
	})
}
