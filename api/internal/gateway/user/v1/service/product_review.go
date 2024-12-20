package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
)

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
			Username:     user.Username,
			AccountID:    user.AccountID,
			ThumbnailURL: user.ThumbnailURL,
			Rate:         review.Rate,
			Title:        review.Title,
			Comment:      review.Comment,
			PublishedAt:  review.CreatedAt.Unix(),
		},
	}
}

func (r *ProductReview) Response() *response.ProductReview {
	return &r.ProductReview
}

func NewProductReviews(reviews entity.ProductReviews, users map[string]*uentity.User) ProductReviews {
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
