package entity

import (
	"iter"

	"github.com/and-period/furumaru/api/pkg/set"
)

// All はインデックスと要素のペアを返すイテレーターを返す。
// slices.All と同様の振る舞いで、for i, p := range ps.All() の形で使用できる。
func (ps Products) All() iter.Seq2[int, *Product] {
	return func(yield func(int, *Product) bool) {
		for i, p := range ps {
			if !yield(i, p) {
				return
			}
		}
	}
}

// IterFilter は指定された商品IDに一致する商品を返すイテレーターを返す。
func (ps Products) IterFilter(productIDs ...string) iter.Seq[*Product] {
	s := set.New(productIDs...)
	return FilterIter(ps, func(p *Product) bool {
		return s.Contains(p.ID)
	})
}

// IterFilterByProducerID は指定された生産者IDに一致する商品を返すイテレーターを返す。
func (ps Products) IterFilterByProducerID(producerID string) iter.Seq[*Product] {
	return FilterIter(ps, func(p *Product) bool {
		return p.ProducerID == producerID
	})
}

// IterFilterBySales は販売中の商品を返すイテレーターを返す。
func (ps Products) IterFilterBySales() iter.Seq[*Product] {
	return FilterIter(ps, func(p *Product) bool {
		return p.Status == ProductStatusForSale
	})
}

// IterFilterByPublished は公開中の商品を返すイテレーターを返す。
// 非公開・アーカイブ済みの商品は除外される。
func (ps Products) IterFilterByPublished() iter.Seq[*Product] {
	return FilterIter(ps, func(p *Product) bool {
		return p.Status != ProductStatusPrivate && p.Status != ProductStatusArchived
	})
}

// IterMap は商品IDをキー、商品を値とするイテレーターを返す。
func (ps Products) IterMap() iter.Seq2[string, *Product] {
	return MapIter(ps, func(p *Product) (string, *Product) {
		return p.ID, p
	})
}

// IterMapByRevision はリビジョンIDをキー、商品を値とするイテレーターを返す。
func (ps Products) IterMapByRevision() iter.Seq2[int64, *Product] {
	return MapIter(ps, func(p *Product) (int64, *Product) {
		return p.ProductRevision.ID, p
	})
}
