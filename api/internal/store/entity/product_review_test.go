package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductReview(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewProductReviewParams
		expect *ProductReview
	}{
		{
			name: "success",
			params: &NewProductReviewParams{
				ProductID: "product-id",
				UserID:    "user-id",
				Rate:      5,
				Title:     "最高の商品",
				Comment:   "おすすめできる商品です",
			},
			expect: &ProductReview{
				ProductID: "product-id",
				UserID:    "user-id",
				Rate:      5,
				Title:     "最高の商品",
				Comment:   "おすすめできる商品です",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProductReview(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProductReviews_SetReactions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		reviews   ProductReviews
		reactions map[string]AggregatedProductReviewReactions
		expect    ProductReviews
	}{
		{
			name: "success",
			reviews: ProductReviews{
				{
					ID:        "preview-id01",
					ProductID: "product-id",
					UserID:    "user-id",
					Rate:      5,
					Title:     "最高の商品",
					Comment:   "おすすめできる商品です",
				},
				{
					ID:        "preview-id02",
					ProductID: "product-id",
					UserID:    "user-id",
					Rate:      5,
					Title:     "最高の商品",
					Comment:   "おすすめできる商品です",
				},
			},
			reactions: map[string]AggregatedProductReviewReactions{
				"preview-id01": {
					{
						ReviewID:     "preview-id01",
						ReactionType: ProductReviewReactionTypeLike,
						Total:        1,
					},
					{
						ReviewID:     "preview-id01",
						ReactionType: ProductReviewReactionTypeDislike,
						Total:        2,
					},
				},
			},
			expect: ProductReviews{
				{
					ID:        "preview-id01",
					ProductID: "product-id",
					UserID:    "user-id",
					Rate:      5,
					Title:     "最高の商品",
					Comment:   "おすすめできる商品です",
					Reactions: map[ProductReviewReactionType]int64{
						ProductReviewReactionTypeLike:    1,
						ProductReviewReactionTypeDislike: 2,
					},
				},
				{
					ID:        "preview-id02",
					ProductID: "product-id",
					UserID:    "user-id",
					Rate:      5,
					Title:     "最高の商品",
					Comment:   "おすすめできる商品です",
					Reactions: map[ProductReviewReactionType]int64{
						ProductReviewReactionTypeLike:    0,
						ProductReviewReactionTypeDislike: 0,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.reviews.SetReactions(tt.reactions)
			assert.ElementsMatch(t, tt.expect, tt.reviews)
		})
	}
}

func TestProductReviews_IDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews ProductReviews
		expect  []string
	}{
		{
			name: "success",
			reviews: ProductReviews{
				{
					ID:        "preview-id",
					ProductID: "product-id",
					UserID:    "user-id",
					Rate:      5,
					Title:     "最高の商品",
					Comment:   "おすすめできる商品です",
				},
			},
			expect: []string{"preview-id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.reviews.IDs()
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestProductReviews_UserIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews ProductReviews
		expect  []string
	}{
		{
			name: "success",
			reviews: ProductReviews{
				{
					ID:        "preview-id",
					ProductID: "product-id",
					UserID:    "user-id",
					Rate:      5,
					Title:     "最高の商品",
					Comment:   "おすすめできる商品です",
				},
			},
			expect: []string{"user-id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.reviews.UserIDs()
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestAggregatedProductReviews_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews AggregatedProductReviews
		expect  map[string]*AggregatedProductReview
	}{
		{
			name: "success",
			reviews: AggregatedProductReviews{
				{
					ProductID: "product-id",
					Count:     4,
					Average:   2.5,
					Rate1:     2,
					Rate2:     0,
					Rate3:     1,
					Rate4:     0,
					Rate5:     1,
				},
			},
			expect: map[string]*AggregatedProductReview{
				"product-id": {
					ProductID: "product-id",
					Count:     4,
					Average:   2.5,
					Rate1:     2,
					Rate2:     0,
					Rate3:     1,
					Rate4:     0,
					Rate5:     1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.reviews.Map()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
