package entity

import "iter"

// All はインデックスと要素のペアを返すイテレーターを返す。
func (vs Videos) All() iter.Seq2[int, *Video] {
	return func(yield func(int, *Video) bool) {
		for i, v := range vs {
			if !yield(i, v) {
				return
			}
		}
	}
}

// IterIDs はビデオIDを返すイテレーターを返す。
func (vs Videos) IterIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, v := range vs {
			if !yield(v.ID) {
				return
			}
		}
	}
}

// IterCoordinatorIDs はコーディネータIDを返すイテレーターを返す。
func (vs Videos) IterCoordinatorIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, v := range vs {
			if !yield(v.CoordinatorID) {
				return
			}
		}
	}
}

// IterProductIDs は各ビデオの商品IDを返すイテレーターを返す。
func (vs Videos) IterProductIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, v := range vs {
			for _, pid := range v.ProductIDs {
				if !yield(pid) {
					return
				}
			}
		}
	}
}

// IterExperienceIDs は各ビデオの体験IDを返すイテレーターを返す。
func (vs Videos) IterExperienceIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, v := range vs {
			for _, eid := range v.ExperienceIDs {
				if !yield(eid) {
					return
				}
			}
		}
	}
}

// IterMap はビデオIDをキー、ビデオを値とするイテレーターを返す。
func (vs Videos) IterMap() iter.Seq2[string, *Video] {
	return MapIter(vs, func(v *Video) (string, *Video) {
		return v.ID, v
	})
}
