package entity

import "iter"

// All はインデックスとコーディネータのペアを返すイテレーターを返す。
func (cs Coordinators) All() iter.Seq2[int, *Coordinator] {
	return func(yield func(int, *Coordinator) bool) {
		for i, c := range cs {
			if !yield(i, c) {
				return
			}
		}
	}
}

// IterMap は管理者IDをキー、コーディネータを値とするイテレーターを返す。
func (cs Coordinators) IterMap() iter.Seq2[string, *Coordinator] {
	return MapIter(cs, func(c *Coordinator) (string, *Coordinator) {
		return c.AdminID, c
	})
}
