package entity

import (
	"iter"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/collection"
)

// All はインデックスと買い物かごのペアを返すイテレーターを返す。
func (bs CartBaskets) All() iter.Seq2[int, *CartBasket] {
	return func(yield func(int, *CartBasket) bool) {
		for i, b := range bs {
			if !yield(i, b) {
				return
			}
		}
	}
}

// IterFilterByCoordinatorID は指定されたコーディネータIDに一致する買い物かごを返すイテレーターを返す。
func (bs CartBaskets) IterFilterByCoordinatorID(coordinatorIDs ...string) iter.Seq[*CartBasket] {
	s := set.New(coordinatorIDs...)
	return collection.FilterIter(bs, func(b *CartBasket) bool {
		return s.Contains(b.CoordinatorID)
	})
}

// IterFilterByBoxNumber は指定された箱の通番に一致する買い物かごを返すイテレーターを返す。
// 0を含む場合、すべての買い物かごを対象とする。
func (bs CartBaskets) IterFilterByBoxNumber(targets ...int64) iter.Seq[*CartBasket] {
	s := set.New(targets...)
	if s.Contains(0) {
		return collection.FilterIter(bs, func(_ *CartBasket) bool {
			return true
		})
	}
	return collection.FilterIter(bs, func(b *CartBasket) bool {
		return s.Contains(b.BoxNumber)
	})
}

// All はインデックスとカートアイテムのペアを返すイテレーターを返す。
func (is CartItems) All() iter.Seq2[int, *CartItem] {
	return func(yield func(int, *CartItem) bool) {
		for i, item := range is {
			if !yield(i, item) {
				return
			}
		}
	}
}

// IterMapByProductID は商品IDをキー、カートアイテムを値とするイテレーターを返す。
func (is CartItems) IterMapByProductID() iter.Seq2[string, *CartItem] {
	return collection.MapIter(is, func(item *CartItem) (string, *CartItem) {
		return item.ProductID, item
	})
}
