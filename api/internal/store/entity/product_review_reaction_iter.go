package entity

import "iter"

// All はインデックスと商品レビューリアクション集計情報のペアを返すイテレーターを返す。
func (rs AggregatedProductReviewReactions) All() iter.Seq2[int, *AggregatedProductReviewReaction] {
	return func(yield func(int, *AggregatedProductReviewReaction) bool) {
		for i, r := range rs {
			if !yield(i, r) {
				return
			}
		}
	}
}

// IterGroupByReviewID は商品レビューIDをキー、リアクション集計一覧を値とするイテレーターを返す。
func (rs AggregatedProductReviewReactions) IterGroupByReviewID() iter.Seq2[string, AggregatedProductReviewReactions] {
	return func(yield func(string, AggregatedProductReviewReactions) bool) {
		groups := rs.GroupByReviewID()
		for k, v := range groups {
			if !yield(k, v) {
				return
			}
		}
	}
}
