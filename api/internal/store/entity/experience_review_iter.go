package entity

import "iter"

// All はインデックスと体験レビューのペアを返すイテレーターを返す。
func (rs ExperienceReviews) All() iter.Seq2[int, *ExperienceReview] {
	return func(yield func(int, *ExperienceReview) bool) {
		for i, r := range rs {
			if !yield(i, r) {
				return
			}
		}
	}
}

// IterMap は体験レビューIDをキー、体験レビューを値とするイテレーターを返す。
func (rs ExperienceReviews) IterMap() iter.Seq2[string, *ExperienceReview] {
	return MapIter(rs, func(r *ExperienceReview) (string, *ExperienceReview) {
		return r.ID, r
	})
}

// All はインデックスと体験レビュー集計情報のペアを返すイテレーターを返す。
func (rs AggregatedExperienceReviews) All() iter.Seq2[int, *AggregatedExperienceReview] {
	return func(yield func(int, *AggregatedExperienceReview) bool) {
		for i, r := range rs {
			if !yield(i, r) {
				return
			}
		}
	}
}

// IterMap は体験IDをキー、体験レビュー集計情報を値とするイテレーターを返す。
func (rs AggregatedExperienceReviews) IterMap() iter.Seq2[string, *AggregatedExperienceReview] {
	return MapIter(rs, func(r *AggregatedExperienceReview) (string, *AggregatedExperienceReview) {
		return r.ExperienceID, r
	})
}
