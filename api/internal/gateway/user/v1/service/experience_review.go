package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
)

// ExperienceReviewReactionType - 体験レビューのリアクション種別
type ExperienceReviewReactionType int32

const (
	ExperienceReviewReactionTypeUnknown ExperienceReviewReactionType = 0
	ExperienceReviewReactionTypeLike    ExperienceReviewReactionType = 1 // いいね
	ExperienceReviewReactionTypeDislike ExperienceReviewReactionType = 2 // いまいち
)

func NewExperienceReviewReactionType(typ entity.ExperienceReviewReactionType) ExperienceReviewReactionType {
	switch typ {
	case entity.ExperienceReviewReactionTypeLike:
		return ExperienceReviewReactionTypeLike
	case entity.ExperienceReviewReactionTypeDislike:
		return ExperienceReviewReactionTypeDislike
	default:
		return ExperienceReviewReactionTypeUnknown
	}
}

func NewExperienceReviewReactionTypeFromRequest(typ int32) (ExperienceReviewReactionType, bool) {
	switch typ {
	case 1:
		return ExperienceReviewReactionTypeLike, true
	case 2:
		return ExperienceReviewReactionTypeDislike, true
	default:
		return ExperienceReviewReactionTypeUnknown, false
	}
}

func (t ExperienceReviewReactionType) StoreEntity() entity.ExperienceReviewReactionType {
	switch t {
	case ExperienceReviewReactionTypeLike:
		return entity.ExperienceReviewReactionTypeLike
	case ExperienceReviewReactionTypeDislike:
		return entity.ExperienceReviewReactionTypeDislike
	default:
		return entity.ExperienceReviewReactionTypeUnknown
	}
}

func (t ExperienceReviewReactionType) Response() int32 {
	return int32(t)
}

type ExperienceReview struct {
	types.ExperienceReview
}

type ExperienceReviews []*ExperienceReview

func NewExperienceReview(review *entity.ExperienceReview, user *uentity.User) *ExperienceReview {
	return &ExperienceReview{
		ExperienceReview: types.ExperienceReview{
			ID:           review.ID,
			ExperienceID: review.ExperienceID,
			UserID:       user.ID,
			Username:     user.Username(),
			AccountID:    user.AccountID,
			ThumbnailURL: user.ThumbnailURL,
			Rate:         review.Rate,
			Title:        review.Title,
			Comment:      review.Comment,
			PublishedAt:  review.CreatedAt.Unix(),
			LikeTotal:    review.Reactions[entity.ExperienceReviewReactionTypeLike],
			DislikeTotal: review.Reactions[entity.ExperienceReviewReactionTypeDislike],
		},
	}
}

func (r *ExperienceReview) Response() *types.ExperienceReview {
	return &r.ExperienceReview
}

func NewExperienceReviews(reviews entity.ExperienceReviews, users map[string]*uentity.User) ExperienceReviews {
	res := make(ExperienceReviews, 0, len(reviews))
	for _, review := range reviews {
		user, ok := users[review.UserID]
		if !ok {
			continue
		}
		res = append(res, NewExperienceReview(review, user))
	}
	return res
}

func (rs ExperienceReviews) Response() []*types.ExperienceReview {
	res := make([]*types.ExperienceReview, len(rs))
	for i := range rs {
		res[i] = rs[i].Response()
	}
	return res
}

type ExperienceReviewReaction struct {
	types.ExperienceReviewReaction
}

type ExperienceReviewReactions []*ExperienceReviewReaction

func NewExperienceReviewReaction(reaction *entity.ExperienceReviewReaction) *ExperienceReviewReaction {
	return &ExperienceReviewReaction{
		ExperienceReviewReaction: types.ExperienceReviewReaction{
			ReviewID:     reaction.ReviewID,
			ReactionType: NewExperienceReviewReactionType(reaction.ReactionType).Response(),
		},
	}
}

func (r *ExperienceReviewReaction) Response() *types.ExperienceReviewReaction {
	return &r.ExperienceReviewReaction
}

func NewExperienceReviewReactions(reactions entity.ExperienceReviewReactions) ExperienceReviewReactions {
	res := make(ExperienceReviewReactions, len(reactions))
	for i := range reactions {
		res[i] = NewExperienceReviewReaction(reactions[i])
	}
	return res
}

func (rs ExperienceReviewReactions) Response() []*types.ExperienceReviewReaction {
	res := make([]*types.ExperienceReviewReaction, len(rs))
	for i := range rs {
		res[i] = rs[i].Response()
	}
	return res
}
