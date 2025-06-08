package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExperienceReview(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewExperienceReviewParams
		expect *ExperienceReview
	}{
		{
			name: "success",
			params: &NewExperienceReviewParams{
				ExperienceID: "experience-id",
				UserID:       "user-id",
				Rate:         5,
				Title:        "最高の体験",
				Comment:      "おすすめできる体験です",
			},
			expect: &ExperienceReview{
				ExperienceID: "experience-id",
				UserID:       "user-id",
				Rate:         5,
				Title:        "最高の体験",
				Comment:      "おすすめできる体験です",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewExperienceReview(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestExperienceReviews_SetReactions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		reviews   ExperienceReviews
		reactions map[string]AggregatedExperienceReviewReactions
		expect    ExperienceReviews
	}{
		{
			name: "success",
			reviews: ExperienceReviews{
				{
					ID:           "preview-id01",
					ExperienceID: "experience-id",
					UserID:       "user-id",
					Rate:         5,
					Title:        "最高の体験",
					Comment:      "おすすめできる体験です",
				},
				{
					ID:           "preview-id02",
					ExperienceID: "experience-id",
					UserID:       "user-id",
					Rate:         5,
					Title:        "最高の体験",
					Comment:      "おすすめできる体験です",
				},
			},
			reactions: map[string]AggregatedExperienceReviewReactions{
				"preview-id01": {
					{
						ReviewID:     "preview-id01",
						ReactionType: ExperienceReviewReactionTypeLike,
						Total:        1,
					},
					{
						ReviewID:     "preview-id01",
						ReactionType: ExperienceReviewReactionTypeDislike,
						Total:        2,
					},
				},
			},
			expect: ExperienceReviews{
				{
					ID:           "preview-id01",
					ExperienceID: "experience-id",
					UserID:       "user-id",
					Rate:         5,
					Title:        "最高の体験",
					Comment:      "おすすめできる体験です",
					Reactions: map[ExperienceReviewReactionType]int64{
						ExperienceReviewReactionTypeLike:    1,
						ExperienceReviewReactionTypeDislike: 2,
					},
				},
				{
					ID:           "preview-id02",
					ExperienceID: "experience-id",
					UserID:       "user-id",
					Rate:         5,
					Title:        "最高の体験",
					Comment:      "おすすめできる体験です",
					Reactions: map[ExperienceReviewReactionType]int64{
						ExperienceReviewReactionTypeLike:    0,
						ExperienceReviewReactionTypeDislike: 0,
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

func TestExperienceReviews_IDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews ExperienceReviews
		expect  []string
	}{
		{
			name: "success",
			reviews: ExperienceReviews{
				{
					ID:           "preview-id",
					ExperienceID: "experience-id",
					UserID:       "user-id",
					Rate:         5,
					Title:        "最高の体験",
					Comment:      "おすすめできる体験です",
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

func TestExperienceReviews_UserIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews ExperienceReviews
		expect  []string
	}{
		{
			name: "success",
			reviews: ExperienceReviews{
				{
					ID:           "preview-id",
					ExperienceID: "experience-id",
					UserID:       "user-id",
					Rate:         5,
					Title:        "最高の体験",
					Comment:      "おすすめできる体験です",
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

func TestAggregatedExperienceReviews_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews AggregatedExperienceReviews
		expect  map[string]*AggregatedExperienceReview
	}{
		{
			name: "success",
			reviews: AggregatedExperienceReviews{
				{
					ExperienceID: "experience-id",
					Count:        4,
					Average:      2.5,
					Rate1:        2,
					Rate2:        0,
					Rate3:        1,
					Rate4:        0,
					Rate5:        1,
				},
			},
			expect: map[string]*AggregatedExperienceReview{
				"experience-id": {
					ExperienceID: "experience-id",
					Count:        4,
					Average:      2.5,
					Rate1:        2,
					Rate2:        0,
					Rate3:        1,
					Rate4:        0,
					Rate5:        1,
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
