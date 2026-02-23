package entity

import "iter"

// All はインデックスとカテゴリのペアを返すイテレーターを返す。
func (cs Categories) All() iter.Seq2[int, *Category] {
	return func(yield func(int, *Category) bool) {
		for i, c := range cs {
			if !yield(i, c) {
				return
			}
		}
	}
}

// IterMapByName はカテゴリ名をキー、カテゴリを値とするイテレーターを返す。
func (cs Categories) IterMapByName() iter.Seq2[string, *Category] {
	return MapIter(cs, func(c *Category) (string, *Category) {
		return c.Name, c
	})
}
