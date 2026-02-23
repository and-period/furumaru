package entity

import (
	"iter"

	"github.com/and-period/furumaru/api/pkg/collection"
)

// All はインデックスとプロモーションのペアを返すイテレーターを返す。
func (ps Promotions) All() iter.Seq2[int, *Promotion] {
	return func(yield func(int, *Promotion) bool) {
		for i, p := range ps {
			if !yield(i, p) {
				return
			}
		}
	}
}

// IterMap はプロモーションIDをキー、プロモーションを値とするイテレーターを返す。
func (ps Promotions) IterMap() iter.Seq2[string, *Promotion] {
	return collection.MapIter(ps, func(p *Promotion) (string, *Promotion) {
		return p.ID, p
	})
}
