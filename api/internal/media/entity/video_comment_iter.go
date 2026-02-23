package entity

import "iter"

// All はインデックスと要素のペアを返すイテレーターを返す。
func (cs VideoComments) All() iter.Seq2[int, *VideoComment] {
	return func(yield func(int, *VideoComment) bool) {
		for i, c := range cs {
			if !yield(i, c) {
				return
			}
		}
	}
}

// IterUserIDs はユーザーIDを返すイテレーターを返す。
// 空のユーザーIDはスキップされる。
func (cs VideoComments) IterUserIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, c := range cs {
			if c.UserID == "" {
				continue
			}
			if !yield(c.UserID) {
				return
			}
		}
	}
}
