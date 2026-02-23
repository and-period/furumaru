package entity

import (
	"iter"

	"github.com/and-period/furumaru/api/pkg/collection"
)

// All はインデックスと生産者のペアを返すイテレーターを返す。
func (ps Producers) All() iter.Seq2[int, *Producer] {
	return func(yield func(int, *Producer) bool) {
		for i, p := range ps {
			if !yield(i, p) {
				return
			}
		}
	}
}

// IterMap は管理者IDをキー、生産者を値とするイテレーターを返す。
func (ps Producers) IterMap() iter.Seq2[string, *Producer] {
	return collection.MapIter(ps, func(p *Producer) (string, *Producer) {
		return p.AdminID, p
	})
}
