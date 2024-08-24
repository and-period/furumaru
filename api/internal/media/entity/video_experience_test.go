package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVideoExperiences_ExperienceIDs(t *testing.T) {
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
				{VideoID: "video-id02", ExperienceID: "experience-id02"},
				{VideoID: "video-id02", ExperienceID: "experience-id03"},
			},
			expect: []string{"experience-id01", "experience-id02", "experience-id03"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.experiences.ExperienceIDs()
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestVideoExperiences_GroupByVideoID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		experiences VideoExperiences
		expect      map[string]VideoExperiences
	}{
		{
			name: "success",
			experiences: VideoExperiences{
				{VideoID: "video-id01", ExperienceID: "experience-id01"},
				{VideoID: "video-id02", ExperienceID: "experience-id02"},
				{VideoID: "video-id02", ExperienceID: "experience-id03"},
			},
			expect: map[string]VideoExperiences{
				"video-id01": {
					{VideoID: "video-id01", ExperienceID: "experience-id01"},
				},
				"video-id02": {
					{VideoID: "video-id02", ExperienceID: "experience-id02"},
					{VideoID: "video-id02", ExperienceID: "experience-id03"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.experiences.GroupByVideoID()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideoExperiences_SortByPrority(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		experiences VideoExperiences
		expect      VideoExperiences
	}{
		{
			name: "success",
			experiences: VideoExperiences{
				{VideoID: "video-id01", Priority: 1, ExperienceID: "experience-id01"},
				{VideoID: "video-id02", Priority: 3, ExperienceID: "experience-id03"},
				{VideoID: "video-id02", Priority: 2, ExperienceID: "experience-id02"},
			},
			expect: VideoExperiences{
				{VideoID: "video-id01", Priority: 1, ExperienceID: "experience-id01"},
				{VideoID: "video-id02", Priority: 2, ExperienceID: "experience-id02"},
				{VideoID: "video-id02", Priority: 3, ExperienceID: "experience-id03"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.experiences.SortByPriority()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
