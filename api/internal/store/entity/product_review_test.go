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
			actual := NewProductReview(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAggregatedProductReviews_UserIDs(t *testing.T) {
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
			actual := tt.reviews.Map()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
