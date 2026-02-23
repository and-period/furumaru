package entity

import (
	"iter"

	"github.com/and-period/furumaru/api/pkg/collection"
)

// All はインデックスとシステム管理者のペアを返すイテレーターを返す。
func (as Administrators) All() iter.Seq2[int, *Administrator] {
	return func(yield func(int, *Administrator) bool) {
		for i, a := range as {
			if !yield(i, a) {
				return
			}
		}
	}
}

// IterMap は管理者IDをキー、システム管理者を値とするイテレーターを返す。
func (as Administrators) IterMap() iter.Seq2[string, *Administrator] {
	return collection.MapIter(as, func(a *Administrator) (string, *Administrator) {
		return a.AdminID, a
	})
}
