package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVideos_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		videos Videos
	}{
		{
			name: "success",
			videos: Videos{
				{ID: "video-id01", Title: "じゃがいもの育て方"},
				{ID: "video-id02", Title: "にんじんの育て方"},
				{ID: "video-id03", Title: "たまねぎの育て方"},
			},
		},
		{
			name:   "empty",
			videos: Videos{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, v := range tt.videos.All() {
				indices = append(indices, i)
				ids = append(ids, v.ID)
			}
			for i, v := range tt.videos {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, v.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.videos))
		})
	}
}

func TestVideos_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	videos := Videos{
		{ID: "video-id01"},
		{ID: "video-id02"},
		{ID: "video-id03"},
	}
	var count int
	for range videos.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestVideos_IterIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		videos Videos
		expect []string
	}{
		{
			name: "success",
			videos: Videos{
				{ID: "video-id01"},
				{ID: "video-id02"},
			},
			expect: []string{"video-id01", "video-id02"},
		},
		{
			name:   "empty",
			videos: Videos{},
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var actual []string
			for id := range tt.videos.IterIDs() {
				actual = append(actual, id)
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideos_IterCoordinatorIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		videos Videos
		expect []string
	}{
		{
			name: "success",
			videos: Videos{
				{ID: "video-id01", CoordinatorID: "coordinator-id01"},
				{ID: "video-id02", CoordinatorID: "coordinator-id02"},
			},
			expect: []string{"coordinator-id01", "coordinator-id02"},
		},
		{
			name:   "empty",
			videos: Videos{},
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var actual []string
			for id := range tt.videos.IterCoordinatorIDs() {
				actual = append(actual, id)
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideos_IterProductIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		videos Videos
		expect []string
	}{
		{
			name: "success",
			videos: Videos{
				{ID: "video-id01", ProductIDs: []string{"product-id01", "product-id02"}},
				{ID: "video-id02", ProductIDs: []string{"product-id03"}},
			},
			expect: []string{"product-id01", "product-id02", "product-id03"},
		},
		{
			name:   "empty",
			videos: Videos{},
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var actual []string
			for id := range tt.videos.IterProductIDs() {
				actual = append(actual, id)
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideos_IterExperienceIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		videos Videos
		expect []string
	}{
		{
			name: "success",
			videos: Videos{
				{ID: "video-id01", ExperienceIDs: []string{"experience-id01"}},
				{ID: "video-id02", ExperienceIDs: []string{"experience-id02", "experience-id03"}},
			},
			expect: []string{"experience-id01", "experience-id02", "experience-id03"},
		},
		{
			name:   "empty",
			videos: Videos{},
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var actual []string
			for id := range tt.videos.IterExperienceIDs() {
				actual = append(actual, id)
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideos_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		videos    Videos
		expectIDs []string
	}{
		{
			name: "success",
			videos: Videos{
				{ID: "video-id01", Title: "じゃがいもの育て方"},
				{ID: "video-id02", Title: "にんじんの育て方"},
			},
			expectIDs: []string{"video-id01", "video-id02"},
		},
		{
			name:      "empty",
			videos:    Videos{},
			expectIDs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Video)
			for k, v := range tt.videos.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.videos))
			for _, id := range tt.expectIDs {
				assert.Contains(t, result, id)
				assert.Equal(t, id, result[id].ID)
			}
		})
	}
}
