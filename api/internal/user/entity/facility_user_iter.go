package entity

import (
	"iter"

	"github.com/and-period/furumaru/api/pkg/collection"
)

// All はインデックスと施設利用者のペアを返すイテレーターを返す。
func (fs FacilityUsers) All() iter.Seq2[int, *FacilityUser] {
	return func(yield func(int, *FacilityUser) bool) {
		for i, f := range fs {
			if !yield(i, f) {
				return
			}
		}
	}
}

// IterMap はユーザーIDをキー、施設利用者を値とするイテレーターを返す。
func (fs FacilityUsers) IterMap() iter.Seq2[string, *FacilityUser] {
	return collection.MapIter(fs, func(f *FacilityUser) (string, *FacilityUser) {
		return f.UserID, f
	})
}
