package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
)

// ProductReviewReactionType - 商品レビューのリアクション種別
type ProductReviewReactionType int32

const (
	ProductReviewReactionTypeUnknown ProductReviewReactionType = 0
	ProductReviewReactionTypeLike    ProductReviewReactionType = 1 // いいね
	ProductReviewReactionTypeDislike ProductReviewReactionType = 2 // いまいち
)

func NewProductReviewReactionType(typ entity.ProductReviewReactionType) ProductReviewReactionType {
	switch typ {
	case entity.ProductReviewReactionTypeLike:
		return ProductReviewReactionTypeLike
	case entity.ProductReviewReactionTypeDislike:
		return ProductReviewReactionTypeDislike
	default:
		return ProductReviewReactionTypeUnknown
	}
}

func NewProductReviewReactionTypeFromRequest(typ int32) (ProductReviewReactionType, bool) {
	switch typ {
	case 1:
		return ProductReviewReactionTypeLike, true
	case 2:
		return ProductReviewReactionTypeDislike, true
	default:
		return ProductReviewReactionTypeUnknown, false
	}
}

func (t ProductReviewReactionType) StoreEntity() entity.ProductReviewReactionType {
	switch t {
	case ProductReviewReactionTypeLike:
		return entity.ProductReviewReactionTypeLike
	case ProductReviewReactionTypeDislike:
		return entity.ProductReviewReactionTypeDislike
	default:
		return entity.ProductReviewReactionTypeUnknown
	}
}

func (t ProductReviewReactionType) Response() int32 {
	return int32(t)
}

type ProductReview struct {
	response.ProductReview
}

type ProductReviews []*ProductReview

func NewProductReview(review *entity.ProductReview, user *uentity.User) *ProductReview {
	return &ProductReview{
		ProductReview: response.ProductReview{
			ID:           review.ID,
			ProductID:    review.ProductID,
			UserID:       user.ID,
			Username:     user.Username(),
			AccountID:    user.AccountID,
			ThumbnailURL: user.ThumbnailURL,
			Rate:         review.Rate,
			Title:        review.Title,
			Comment:      review.Comment,
			PublishedAt:  review.CreatedAt.Unix(),
			LikeTotal:    review.Reactions[entity.ProductReviewReactionTypeLike],
			DislikeTotal: review.Reactions[entity.ProductReviewReactionTypeDislike],
		},
	}
}

func (r *ProductReview) Response() *response.ProductReview {
	return &r.ProductReview
}

func NewProductReviews(
	reviews entity.ProductReviews,
	users map[string]*uentity.User,
) ProductReviews {
	res := make(ProductReviews, 0, len(reviews))
	for _, review := range reviews {
		user, ok := users[review.UserID]
		if !ok {
			continue
		}
		res = append(res, NewProductReview(review, user))
	}
	return res
}

func (rs ProductReviews) Response() []*response.ProductReview {
	res := make([]*response.ProductReview, len(rs))
	for i := range rs {
		res[i] = rs[i].Response()
	}
	return res
}

type ProductReviewReaction struct {
	response.ProductReviewReaction
}

type ProductReviewReactions []*ProductReviewReaction

func NewProductReviewReaction(reaction *entity.ProductReviewReaction) *ProductReviewReaction {
	return &ProductReviewReaction{
		ProductReviewReaction: response.ProductReviewReaction{
			ReviewID:     reaction.ReviewID,
			ReactionType: NewProductReviewReactionType(reaction.ReactionType).Response(),
		},
	}
}

func (r *ProductReviewReaction) Response() *response.ProductReviewReaction {
	return &r.ProductReviewReaction
}

func NewProductReviewReactions(reactions entity.ProductReviewReactions) ProductReviewReactions {
	res := make(ProductReviewReactions, len(reactions))
	for i := range reactions {
		res[i] = NewProductReviewReaction(reactions[i])
	}
	return res
}

func (rs ProductReviewReactions) Response() []*response.ProductReviewReaction {
	res := make([]*response.ProductReviewReaction, len(rs))
	for i := range rs {
		res[i] = rs[i].Response()
	}
	return res
}
