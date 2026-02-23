package entity

import "iter"

// All はインデックスと店舗のペアを返すイテレーターを返す。
func (ss Shops) All() iter.Seq2[int, *Shop] {
	return func(yield func(int, *Shop) bool) {
		for i, s := range ss {
			if !yield(i, s) {
				return
			}
		}
	}
}

// IterMapByCoordinatorID はコーディネータIDをキー、店舗を値とするイテレーターを返す。
func (ss Shops) IterMapByCoordinatorID() iter.Seq2[string, *Shop] {
	return MapIter(ss, func(s *Shop) (string, *Shop) {
		return s.CoordinatorID, s
	})
}

// IterGroupByProducerID は生産者IDをキー、店舗を値とするイテレーターを返す。
// 各店舗は所属する生産者IDごとに展開されるため、1つの店舗が複数回出現する場合がある。
func (ss Shops) IterGroupByProducerID() iter.Seq2[string, *Shop] {
	return func(yield func(string, *Shop) bool) {
		for _, s := range ss {
			for _, producerID := range s.ProducerIDs {
				if !yield(producerID, s) {
					return
				}
			}
		}
	}
}
