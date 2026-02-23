package entity

import (
	"iter"

	"github.com/and-period/furumaru/api/pkg/collection"
)

// All はインデックスと要素のペアを返すイテレーターを返す。
func (bs Broadcasts) All() iter.Seq2[int, *Broadcast] {
	return func(yield func(int, *Broadcast) bool) {
		for i, b := range bs {
			if !yield(i, b) {
				return
			}
		}
	}
}

// IterScheduleIDs はスケジュールIDを返すイテレーターを返す。
func (bs Broadcasts) IterScheduleIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, b := range bs {
			if !yield(b.ScheduleID) {
				return
			}
		}
	}
}

// IterMap はライブ配信IDをキー、ライブ配信を値とするイテレーターを返す。
func (bs Broadcasts) IterMap() iter.Seq2[string, *Broadcast] {
	return collection.MapIter(bs, func(b *Broadcast) (string, *Broadcast) {
		return b.ID, b
	})
}
