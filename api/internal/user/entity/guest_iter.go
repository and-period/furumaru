package entity

import (
	"iter"

	"github.com/and-period/furumaru/api/pkg/collection"
)

// All はインデックスとゲストのペアを返すイテレーターを返す。
func (gs Guests) All() iter.Seq2[int, *Guest] {
	return func(yield func(int, *Guest) bool) {
		for i, g := range gs {
			if !yield(i, g) {
				return
			}
		}
	}
}

// IterMap はユーザーIDをキー、ゲストを値とするイテレーターを返す。
func (gs Guests) IterMap() iter.Seq2[string, *Guest] {
	return collection.MapIter(gs, func(g *Guest) (string, *Guest) {
		return g.UserID, g
	})
}
