package entity

import "time"

// ExperienceReviewReactionType - 体験レビューのリアクション種別
type ExperienceReviewReactionType int32

const (
	ExperienceReviewReactionTypeUnknown ExperienceReviewReactionType = 0
	ExperienceReviewReactionTypeLike    ExperienceReviewReactionType = 1 // いいね
	ExperienceReviewReactionTypeDislike ExperienceReviewReactionType = 2 // いまいち
)

// ExperienceReviewReaction - 体験レビューに対するリアクション
type ExperienceReviewReaction struct {
	ReviewID     string                       `gorm:"primaryKey;<-:create"` // 体験レビューID
	UserID       string                       `gorm:"primaryKey;<-:create"` // ユーザーID
	ReactionType ExperienceReviewReactionType `gorm:""`                     // リアクション種別
	CreatedAt    time.Time                    `gorm:"<-:create"`            // 作成日時
	UpdatedAt    time.Time                    `gorm:""`                     // 更新日時
}

type ExperienceReviewReactions []*ExperienceReviewReaction

type NewExperienceReviewReactionParams struct {
	ReviewID     string
	UserID       string
	ReactionType ExperienceReviewReactionType
}

func NewExperienceReviewReaction(params *NewExperienceReviewReactionParams) *ExperienceReviewReaction {
	return &ExperienceReviewReaction{
		ReviewID:     params.ReviewID,
		UserID:       params.UserID,
		ReactionType: params.ReactionType,
	}
}

// AggregatedExperienceReviewReaction - 体験レビューのリアクション
type AggregatedExperienceReviewReaction struct {
	ReviewID     string                       `gorm:""` // 体験レビューID
	ReactionType ExperienceReviewReactionType `gorm:""` // リアクション種別
	Total        int64                        `gorm:""` // リアクション数
}

type AggregatedExperienceReviewReactions []*AggregatedExperienceReviewReaction

func newEmptyExperienceReviewReactionTotal() map[ExperienceReviewReactionType]int64 {
	return map[ExperienceReviewReactionType]int64{
		ExperienceReviewReactionTypeLike:    0,
		ExperienceReviewReactionTypeDislike: 0,
	}
}

func (rs AggregatedExperienceReviewReactions) GetTotalByMap() map[ExperienceReviewReactionType]int64 {
	res := newEmptyExperienceReviewReactionTotal()
	for _, r := range rs {
		res[r.ReactionType] = r.Total
	}
	return res
}

func (rs AggregatedExperienceReviewReactions) GroupByReviewID() map[string]AggregatedExperienceReviewReactions {
	res := make(map[string]AggregatedExperienceReviewReactions, len(rs))
	for _, r := range rs {
		if _, ok := res[r.ReviewID]; !ok {
			res[r.ReviewID] = make(AggregatedExperienceReviewReactions, 0)
		}
		res[r.ReviewID] = append(res[r.ReviewID], r)
	}
	return res
}
