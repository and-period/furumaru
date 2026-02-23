package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregatedProductReviewReactions_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		reactions AggregatedProductReviewReactions
	}{
		{
			name: "success",
			reactions: AggregatedProductReviewReactions{
				{ReviewID: "review-01", ReactionType: ProductReviewReactionTypeLike, Total: 5},
				{ReviewID: "review-01", ReactionType: ProductReviewReactionTypeDislike, Total: 1},
			},
		},
		{
			name:      "empty",
			reactions: AggregatedProductReviewReactions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var count int
			for range tt.reactions.All() {
				count++
			}
			assert.Equal(t, len(tt.reactions), count)
		})
	}
}

func TestAggregatedProductReviewReactions_IterGroupByReviewID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		reactions  AggregatedProductReviewReactions
		expectKeys int
	}{
		{
			name: "success",
			reactions: AggregatedProductReviewReactions{
				{ReviewID: "review-01", ReactionType: ProductReviewReactionTypeLike, Total: 5},
				{ReviewID: "review-01", ReactionType: ProductReviewReactionTypeDislike, Total: 1},
				{ReviewID: "review-02", ReactionType: ProductReviewReactionTypeLike, Total: 3},
			},
			expectKeys: 2,
		},
		{
			name:       "empty",
			reactions:  AggregatedProductReviewReactions{},
			expectKeys: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]AggregatedProductReviewReactions)
			for k, v := range tt.reactions.IterGroupByReviewID() {
				result[k] = v
			}
			assert.Len(t, result, tt.expectKeys)
		})
	}
}
