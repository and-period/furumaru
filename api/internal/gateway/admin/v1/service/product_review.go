package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type ProductReview struct {
	types.ProductReview
}

type ProductReviews []*ProductReview

func NewProductReview(review *entity.ProductReview, user *User) *ProductReview {
	return &ProductReview{
		ProductReview: types.ProductReview{
			ID:           review.ID,
			ProductID:    review.ProductID,
			UserID:       user.ID,
			Username:     user.Username,
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

func (r *ProductReview) Response() *types.ProductReview {
	return &r.ProductReview
}

func NewProductReviews(reviews entity.ProductReviews, users map[string]*User) ProductReviews {
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

func (rs ProductReviews) Response() []*types.ProductReview {
	res := make([]*types.ProductReview, len(rs))
	for i := range rs {
		res[i] = rs[i].Response()
	}
	return res
}
