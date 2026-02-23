package entity

import (
	"iter"

	"github.com/and-period/furumaru/api/pkg/collection"
)

// All はインデックスと体験変更履歴のペアを返すイテレーターを返す。
func (rs ExperienceRevisions) All() iter.Seq2[int, *ExperienceRevision] {
	return func(yield func(int, *ExperienceRevision) bool) {
		for i, r := range rs {
			if !yield(i, r) {
				return
			}
		}
	}
}

// IterMapByExperienceID は体験IDをキー、体験変更履歴を値とするイテレーターを返す。
func (rs ExperienceRevisions) IterMapByExperienceID() iter.Seq2[string, *ExperienceRevision] {
	return collection.MapIter(rs, func(r *ExperienceRevision) (string, *ExperienceRevision) {
		return r.ExperienceID, r
	})
}
