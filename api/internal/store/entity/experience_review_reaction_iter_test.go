package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregatedExperienceReviewReactions_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		reactions AggregatedExperienceReviewReactions
	}{
		{
			name: "success",
			reactions: AggregatedExperienceReviewReactions{
				{ReviewID: "review-01", ReactionType: ExperienceReviewReactionTypeLike, Total: 5},
				{ReviewID: "review-01", ReactionType: ExperienceReviewReactionTypeDislike, Total: 1},
			},
		},
		{
			name:      "empty",
			reactions: AggregatedExperienceReviewReactions{},
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

func TestAggregatedExperienceReviewReactions_IterGroupByReviewID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		reactions  AggregatedExperienceReviewReactions
		expectKeys int
	}{
		{
			name: "success",
			reactions: AggregatedExperienceReviewReactions{
				{ReviewID: "review-01", ReactionType: ExperienceReviewReactionTypeLike, Total: 5},
				{ReviewID: "review-01", ReactionType: ExperienceReviewReactionTypeDislike, Total: 1},
				{ReviewID: "review-02", ReactionType: ExperienceReviewReactionTypeLike, Total: 3},
			},
			expectKeys: 2,
		},
		{
			name:       "empty",
			reactions:  AggregatedExperienceReviewReactions{},
			expectKeys: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]AggregatedExperienceReviewReactions)
			for k, v := range tt.reactions.IterGroupByReviewID() {
				result[k] = v
			}
			assert.Len(t, result, tt.expectKeys)
		})
	}
}
