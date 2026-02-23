package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVideoExperiences_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		experiences VideoExperiences
	}{
		{
			name: "success",
			experiences: VideoExperiences{
				{VideoID: "video-id01", ExperienceID: "experience-id01", Priority: 1},
				{VideoID: "video-id01", ExperienceID: "experience-id02", Priority: 2},
			},
		},
		{
			name:        "empty",
			experiences: VideoExperiences{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var experienceIDs []string
			for i, e := range tt.experiences.All() {
				indices = append(indices, i)
				experienceIDs = append(experienceIDs, e.ExperienceID)
			}
			for i, e := range tt.experiences {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, e.ExperienceID, experienceIDs[i])
				}
			}
			assert.Len(t, indices, len(tt.experiences))
		})
	}
}

func TestVideoExperiences_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	experiences := VideoExperiences{
		{VideoID: "video-id01", ExperienceID: "experience-id01"},
		{VideoID: "video-id01", ExperienceID: "experience-id02"},
		{VideoID: "video-id01", ExperienceID: "experience-id03"},
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

func TestVideoExperiences_IterExperienceIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		experiences VideoExperiences
		expect      []string
	}{
		{
			name: "success",
			experiences: VideoExperiences{
				{VideoID: "video-id01", ExperienceID: "experience-id01"},
				{VideoID: "video-id01", ExperienceID: "experience-id02"},
			},
			expect: []string{"experience-id01", "experience-id02"},
		},
		{
			name:        "empty",
			experiences: VideoExperiences{},
			expect:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var actual []string
			for id := range tt.experiences.IterExperienceIDs() {
				actual = append(actual, id)
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideoExperiences_IterGroupByVideoID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		experiences VideoExperiences
		expect      map[string]int // videoID -> count
	}{
		{
			name: "success",
			experiences: VideoExperiences{
				{VideoID: "video-id01", ExperienceID: "experience-id01"},
				{VideoID: "video-id01", ExperienceID: "experience-id02"},
				{VideoID: "video-id02", ExperienceID: "experience-id03"},
			},
			expect: map[string]int{
				"video-id01": 2,
				"video-id02": 1,
			},
		},
		{
			name:        "empty",
			experiences: VideoExperiences{},
			expect:      map[string]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]int)
			for k, v := range tt.experiences.IterGroupByVideoID() {
				result[k] = len(v)
			}
			assert.Equal(t, tt.expect, result)
		})
	}
}
