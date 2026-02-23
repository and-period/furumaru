package entity

import (
	"iter"

	"github.com/and-period/furumaru/api/pkg/collection"
)

// All はインデックスと商品レビューのペアを返すイテレーターを返す。
func (rs ProductReviews) All() iter.Seq2[int, *ProductReview] {
	return func(yield func(int, *ProductReview) bool) {
		for i, r := range rs {
			if !yield(i, r) {
				return
			}
		}
	}
}

// IterMap は商品レビューIDをキー、商品レビューを値とするイテレーターを返す。
func (rs ProductReviews) IterMap() iter.Seq2[string, *ProductReview] {
	return collection.MapIter(rs, func(r *ProductReview) (string, *ProductReview) {
		return r.ID, r
	})
}

// All はインデックスと商品レビュー集計情報のペアを返すイテレーターを返す。
func (rs AggregatedProductReviews) All() iter.Seq2[int, *AggregatedProductReview] {
	return func(yield func(int, *AggregatedProductReview) bool) {
		for i, r := range rs {
			if !yield(i, r) {
				return
			}
		}
	}
}

// IterMap は商品IDをキー、商品レビュー集計情報を値とするイテレーターを返す。
func (rs AggregatedProductReviews) IterMap() iter.Seq2[string, *AggregatedProductReview] {
	return collection.MapIter(rs, func(r *AggregatedProductReview) (string, *AggregatedProductReview) {
		return r.ProductID, r
	})
}
