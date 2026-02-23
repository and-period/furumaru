package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExperienceReviews_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews ExperienceReviews
	}{
		{
			name: "success",
			reviews: ExperienceReviews{
				{ID: "review-01", ExperienceID: "exp-01"},
				{ID: "review-02", ExperienceID: "exp-02"},
			},
		},
		{
			name:    "empty",
			reviews: ExperienceReviews{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, r := range tt.reviews.All() {
				indices = append(indices, i)
				ids = append(ids, r.ID)
			}
			for i, r := range tt.reviews {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, r.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.reviews))
		})
	}
}

func TestExperienceReviews_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews ExperienceReviews
	}{
		{
			name: "success",
			reviews: ExperienceReviews{
				{ID: "review-01", ExperienceID: "exp-01"},
				{ID: "review-02", ExperienceID: "exp-02"},
			},
		},
		{
			name:    "empty",
			reviews: ExperienceReviews{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*ExperienceReview)
			for k, v := range tt.reviews.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.reviews))
			for _, r := range tt.reviews {
				assert.Contains(t, result, r.ID)
			}
		})
	}
}

func TestAggregatedExperienceReviews_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews AggregatedExperienceReviews
	}{
		{
			name: "success",
			reviews: AggregatedExperienceReviews{
				{ExperienceID: "exp-01", Count: 10, Average: 4.5},
				{ExperienceID: "exp-02", Count: 5, Average: 3.8},
			},
		},
		{
			name:    "empty",
			reviews: AggregatedExperienceReviews{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var count int
			for range tt.reviews.All() {
				count++
			}
			assert.Equal(t, len(tt.reviews), count)
		})
	}
}

func TestAggregatedExperienceReviews_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews AggregatedExperienceReviews
	}{
		{
			name: "success",
			reviews: AggregatedExperienceReviews{
				{ExperienceID: "exp-01", Count: 10, Average: 4.5},
				{ExperienceID: "exp-02", Count: 5, Average: 3.8},
			},
		},
		{
			name:    "empty",
			reviews: AggregatedExperienceReviews{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*AggregatedExperienceReview)
			for k, v := range tt.reviews.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.reviews))
			for _, r := range tt.reviews {
				assert.Contains(t, result, r.ExperienceID)
			}
		})
	}
}
