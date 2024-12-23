package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExperienceReviewReaction(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewExperienceReviewReactionParams
		expect *ExperienceReviewReaction
	}{
		{
			name: "success",
			params: &NewExperienceReviewReactionParams{
				ReviewID:     "review-id",
				UserID:       "user-id",
				ReactionType: ExperienceReviewReactionTypeLike,
			},
			expect: &ExperienceReviewReaction{
				ReviewID:     "review-id",
				UserID:       "user-id",
				ReactionType: ExperienceReviewReactionTypeLike,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewExperienceReviewReaction(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAggregatedExperienceReviewReactions_GetTotalByMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews AggregatedExperienceReviewReactions
		expect  map[ExperienceReviewReactionType]int64
	}{
		{
			name: "success",
			reviews: AggregatedExperienceReviewReactions{
				{
					ReviewID:     "review-id01",
					ReactionType: ExperienceReviewReactionTypeLike,
					Total:        1,
				},
				{
					ReviewID:     "review-id01",
					ReactionType: ExperienceReviewReactionTypeDislike,
					Total:        2,
				},
			},
			expect: map[ExperienceReviewReactionType]int64{
				ExperienceReviewReactionTypeLike:    1,
				ExperienceReviewReactionTypeDislike: 2,
			},
		},
		{
			name:    "empty",
			reviews: AggregatedExperienceReviewReactions{},
			expect: map[ExperienceReviewReactionType]int64{
				ExperienceReviewReactionTypeLike:    0,
				ExperienceReviewReactionTypeDislike: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.reviews.GetTotalByMap()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAggregatedExperienceReviewReactions_GroupByReviewID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews AggregatedExperienceReviewReactions
		expect  map[string]AggregatedExperienceReviewReactions
	}{
		{
			name: "success",
			reviews: AggregatedExperienceReviewReactions{
				{
					ReviewID:     "review-id01",
					ReactionType: ExperienceReviewReactionTypeLike,
					Total:        1,
				},
				{
					ReviewID:     "review-id01",
					ReactionType: ExperienceReviewReactionTypeDislike,
					Total:        2,
				},
				{
					ReviewID:     "review-id02",
					ReactionType: ExperienceReviewReactionTypeLike,
					Total:        3,
				},
				{
					ReviewID:     "review-id02",
					ReactionType: ExperienceReviewReactionTypeDislike,
					Total:        4,
				},
			},
			expect: map[string]AggregatedExperienceReviewReactions{
				"review-id01": {
					{
						ReviewID:     "review-id01",
						ReactionType: ExperienceReviewReactionTypeLike,
						Total:        1,
					},
					{
						ReviewID:     "review-id01",
						ReactionType: ExperienceReviewReactionTypeDislike,
						Total:        2,
					},
				},
				"review-id02": {
					{
						ReviewID:     "review-id02",
						ReactionType: ExperienceReviewReactionTypeLike,
						Total:        3,
					},
					{
						ReviewID:     "review-id02",
						ReactionType: ExperienceReviewReactionTypeDislike,
						Total:        4,
					},
				},
			},
		},
		{
			name:    "empty",
			reviews: AggregatedExperienceReviewReactions{},
			expect:  map[string]AggregatedExperienceReviewReactions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.reviews.GroupByReviewID()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
