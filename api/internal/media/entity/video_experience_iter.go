package entity

import "iter"

// All はインデックスと要素のペアを返すイテレーターを返す。
func (es VideoExperiences) All() iter.Seq2[int, *VideoExperience] {
	return func(yield func(int, *VideoExperience) bool) {
		for i, e := range es {
			if !yield(i, e) {
				return
			}
		}
	}
}

// IterExperienceIDs は体験IDを返すイテレーターを返す。
func (es VideoExperiences) IterExperienceIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, e := range es {
			if !yield(e.ExperienceID) {
				return
			}
		}
	}
}

// IterGroupByVideoID はビデオIDをキー、VideoExperiencesを値とするイテレーターを返す。
func (es VideoExperiences) IterGroupByVideoID() iter.Seq2[string, VideoExperiences] {
	return func(yield func(string, VideoExperiences) bool) {
		groups := es.GroupByVideoID()
		for k, v := range groups {
			if !yield(k, v) {
				return
			}
		}
	}
}
