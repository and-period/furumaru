package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductReviewReaction(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewProductReviewReactionParams
		expect *ProductReviewReaction
	}{
		{
			name: "success",
			params: &NewProductReviewReactionParams{
				ReviewID:     "review-id",
				UserID:       "user-id",
				ReactionType: ProductReviewReactionTypeLike,
			},
			expect: &ProductReviewReaction{
				ReviewID:     "review-id",
				UserID:       "user-id",
				ReactionType: ProductReviewReactionTypeLike,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProductReviewReaction(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAggregatedProductReviewReactions_GetTotalByMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews AggregatedProductReviewReactions
		expect  map[ProductReviewReactionType]int64
	}{
		{
			name: "success",
			reviews: AggregatedProductReviewReactions{
				{
					ReviewID:     "review-id01",
					ReactionType: ProductReviewReactionTypeLike,
					Total:        1,
				},
				{
					ReviewID:     "review-id01",
					ReactionType: ProductReviewReactionTypeDislike,
					Total:        2,
				},
			},
			expect: map[ProductReviewReactionType]int64{
				ProductReviewReactionTypeLike:    1,
				ProductReviewReactionTypeDislike: 2,
			},
		},
		{
			name:    "empty",
			reviews: AggregatedProductReviewReactions{},
			expect: map[ProductReviewReactionType]int64{
				ProductReviewReactionTypeLike:    0,
				ProductReviewReactionTypeDislike: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.reviews.GetTotalByMap()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAggregatedProductReviewReactions_GroupByReviewID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews AggregatedProductReviewReactions
		expect  map[string]AggregatedProductReviewReactions
	}{
		{
			name: "success",
			reviews: AggregatedProductReviewReactions{
				{
					ReviewID:     "review-id01",
					ReactionType: ProductReviewReactionTypeLike,
					Total:        1,
				},
				{
					ReviewID:     "review-id01",
					ReactionType: ProductReviewReactionTypeDislike,
					Total:        2,
				},
				{
					ReviewID:     "review-id02",
					ReactionType: ProductReviewReactionTypeLike,
					Total:        3,
				},
				{
					ReviewID:     "review-id02",
					ReactionType: ProductReviewReactionTypeDislike,
					Total:        4,
				},
			},
			expect: map[string]AggregatedProductReviewReactions{
				"review-id01": {
					{
						ReviewID:     "review-id01",
						ReactionType: ProductReviewReactionTypeLike,
						Total:        1,
					},
					{
						ReviewID:     "review-id01",
						ReactionType: ProductReviewReactionTypeDislike,
						Total:        2,
					},
				},
				"review-id02": {
					{
						ReviewID:     "review-id02",
						ReactionType: ProductReviewReactionTypeLike,
						Total:        3,
					},
					{
						ReviewID:     "review-id02",
						ReactionType: ProductReviewReactionTypeDislike,
						Total:        4,
					},
				},
			},
		},
		{
			name:    "empty",
			reviews: AggregatedProductReviewReactions{},
			expect:  map[string]AggregatedProductReviewReactions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.reviews.GroupByReviewID()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
