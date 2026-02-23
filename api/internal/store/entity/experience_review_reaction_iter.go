package entity

import "iter"

// All はインデックスと体験レビューリアクション集計情報のペアを返すイテレーターを返す。
func (rs AggregatedExperienceReviewReactions) All() iter.Seq2[int, *AggregatedExperienceReviewReaction] {
	return func(yield func(int, *AggregatedExperienceReviewReaction) bool) {
		for i, r := range rs {
			if !yield(i, r) {
				return
			}
		}
	}
}

// IterGroupByReviewID は体験レビューIDをキー、リアクション集計一覧を値とするイテレーターを返す。
func (rs AggregatedExperienceReviewReactions) IterGroupByReviewID() iter.Seq2[string, AggregatedExperienceReviewReactions] {
	return func(yield func(string, AggregatedExperienceReviewReactions) bool) {
		groups := rs.GroupByReviewID()
		for k, v := range groups {
			if !yield(k, v) {
				return
			}
		}
	}
}
