package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExperienceRevisions_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		revisions ExperienceRevisions
	}{
		{
			name: "success",
			revisions: ExperienceRevisions{
				{ID: 1, ExperienceID: "exp-01"},
				{ID: 2, ExperienceID: "exp-02"},
			},
		},
		{
			name:      "empty",
			revisions: ExperienceRevisions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var count int
			for range tt.revisions.All() {
				count++
			}
			assert.Equal(t, len(tt.revisions), count)
		})
	}
}

func TestExperienceRevisions_IterMapByExperienceID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		revisions ExperienceRevisions
	}{
		{
			name: "success",
			revisions: ExperienceRevisions{
				{ID: 1, ExperienceID: "exp-01"},
				{ID: 2, ExperienceID: "exp-02"},
			},
		},
		{
			name:      "empty",
			revisions: ExperienceRevisions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*ExperienceRevision)
			for k, v := range tt.revisions.IterMapByExperienceID() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.revisions))
			for _, r := range tt.revisions {
				assert.Contains(t, result, r.ExperienceID)
			}
		})
	}
}
