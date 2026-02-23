package entity

import (
	"iter"

	"github.com/and-period/furumaru/api/pkg/collection"
)

// All はインデックスと管理者のペアを返すイテレーターを返す。
func (as Admins) All() iter.Seq2[int, *Admin] {
	return func(yield func(int, *Admin) bool) {
		for i, a := range as {
			if !yield(i, a) {
				return
			}
		}
	}
}

// IterMap は管理者IDをキー、管理者を値とするイテレーターを返す。
func (as Admins) IterMap() iter.Seq2[string, *Admin] {
	return collection.MapIter(as, func(a *Admin) (string, *Admin) {
		return a.ID, a
	})
}

// IterGroupByType は管理者種別をキー、管理者を値とするイテレーターを返す。
func (as Admins) IterGroupByType() iter.Seq2[AdminType, *Admin] {
	return collection.MapIter(as, func(a *Admin) (AdminType, *Admin) {
		return a.Type, a
	})
}
