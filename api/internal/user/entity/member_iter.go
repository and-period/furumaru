package entity

import (
	"iter"

	"github.com/and-period/furumaru/api/pkg/collection"
)

// All はインデックスと会員のペアを返すイテレーターを返す。
func (ms Members) All() iter.Seq2[int, *Member] {
	return func(yield func(int, *Member) bool) {
		for i, m := range ms {
			if !yield(i, m) {
				return
			}
		}
	}
}

// IterMap はユーザーIDをキー、会員を値とするイテレーターを返す。
func (ms Members) IterMap() iter.Seq2[string, *Member] {
	return collection.MapIter(ms, func(m *Member) (string, *Member) {
		return m.UserID, m
	})
}
