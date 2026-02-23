package entity

import "iter"

// All はインデックスと住所のペアを返すイテレーターを返す。
func (as Addresses) All() iter.Seq2[int, *Address] {
	return func(yield func(int, *Address) bool) {
		for i, a := range as {
			if !yield(i, a) {
				return
			}
		}
	}
}

// IterMap は住所IDをキー、住所を値とするイテレーターを返す。
func (as Addresses) IterMap() iter.Seq2[string, *Address] {
	return MapIter(as, func(a *Address) (string, *Address) {
		return a.ID, a
	})
}

// IterMapByRevision はリビジョンIDをキー、住所を値とするイテレーターを返す。
func (as Addresses) IterMapByRevision() iter.Seq2[int64, *Address] {
	return MapIter(as, func(a *Address) (int64, *Address) {
		return a.AddressRevision.ID, a
	})
}

// IterMapByUserID はユーザーIDをキー、住所を値とするイテレーターを返す。
func (as Addresses) IterMapByUserID() iter.Seq2[string, *Address] {
	return MapIter(as, func(a *Address) (string, *Address) {
		return a.UserID, a
	})
}
