package entity

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExperiences_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		experiences Experiences
	}{
		{
			name: "success",
			experiences: Experiences{
				{ID: "exp-01", Title: "いちご狩り"},
				{ID: "exp-02", Title: "みかん狩り"},
			},
		},
		{
			name:        "empty",
			experiences: Experiences{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, e := range tt.experiences.All() {
				indices = append(indices, i)
				ids = append(ids, e.ID)
			}
			for i, e := range tt.experiences {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, e.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.experiences))
		})
	}
}

func TestExperiences_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	experiences := Experiences{
		{ID: "exp-01"},
		{ID: "exp-02"},
		{ID: "exp-03"},
	}
	var count int
	for range experiences.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestExperiences_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		experiences Experiences
		expectIDs   []string
	}{
		{
			name: "success",
			experiences: Experiences{
				{ID: "exp-01", Title: "いちご狩り"},
				{ID: "exp-02", Title: "みかん狩り"},
			},
			expectIDs: []string{"exp-01", "exp-02"},
		},
		{
			name:        "empty",
			experiences: Experiences{},
			expectIDs:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Experience)
			for k, v := range tt.experiences.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.experiences))
			for _, id := range tt.expectIDs {
				assert.Contains(t, result, id)
				assert.Equal(t, id, result[id].ID)
			}
		})
	}
}

func TestExperiences_IterFilterByPublished(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		experiences Experiences
		expectIDs   []string
	}{
		{
			name: "success",
			experiences: Experiences{
				{ID: "exp-01", Status: ExperienceStatusAccepting},
				{ID: "exp-02", Status: ExperienceStatusPrivate},
				{ID: "exp-03", Status: ExperienceStatusWaiting},
				{ID: "exp-04", Status: ExperienceStatusArchived},
			},
			expectIDs: []string{"exp-01", "exp-03"},
		},
		{
			name: "all private",
			experiences: Experiences{
				{ID: "exp-01", Status: ExperienceStatusPrivate},
			},
			expectIDs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			collected := slices.Collect(tt.experiences.IterFilterByPublished())
			if len(collected) == 0 {
				assert.Equal(t, tt.expectIDs, []string(nil))
				return
			}
			ids := make([]string, len(collected))
			for i, e := range collected {
				ids[i] = e.ID
			}
			assert.Equal(t, tt.expectIDs, ids)
		})
	}
}
