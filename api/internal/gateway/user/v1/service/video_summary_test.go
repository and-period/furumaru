package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestVideoSummaries(t *testing.T) {
	t.Parallel()

	now := jst.Date(2024, 8, 24, 18, 30, 0, 0)

	tests := []struct {
		name   string
		videos entity.Videos
		expect VideoSummaries
	}{
		{
			name: "success",
			videos: entity.Videos{
				{
					ID:                "video-id",
					CoordinatorID:     "coordinator-id",
					ProductIDs:        []string{"product-id"},
					ExperienceIDs:     []string{"experience-id"},
					Title:             "じゃがいもの育て方",
					Description:       "じゃがいもの育て方の動画です。",
					Status:            entity.VideoStatusPublished,
					ThumbnailURL:      "https://example.com/thumbnail.jpg",
					VideoURL:          "https://example.com/video.mp4",
					Public:            true,
					Limited:           false,
					DisplayProduct:    true,
					DisplayExperience: true,
					VideoProducts: entity.VideoProducts{{
						VideoID:   "video-id",
						ProductID: "product-id",
						Priority:  1,
						CreatedAt: now,
						UpdatedAt: now,
					}},
					VideoExperiences: entity.VideoExperiences{{
						VideoID:      "video-id",
						ExperienceID: "experience-id",
						Priority:     1,
						CreatedAt:    now,
						UpdatedAt:    now,
					}},
					PublishedAt: now.AddDate(0, 0, -1),
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			expect: VideoSummaries{
				{
					VideoSummary: response.VideoSummary{
						ID:            "video-id",
						CoordinatorID: "coordinator-id",
						Title:         "じゃがいもの育て方",
						ThumbnailURL:  "https://example.com/thumbnail.jpg",
						PublishedAt:   now.AddDate(0, 0, -1).Unix(),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewVideoSummaries(tt.videos)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideoSummaries_CoordinatorIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		videos VideoSummaries
		expect []string
	}{
		{
			name: "success",
			videos: VideoSummaries{
				{
					VideoSummary: response.VideoSummary{
						ID:            "video-id",
						CoordinatorID: "coordinator-id",
						Title:         "じゃがいもの育て方",
						ThumbnailURL:  "https://example.com/thumbnail.jpg",
						PublishedAt:   0,
					},
				},
			},
			expect: []string{"coordinator-id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.videos.CoordinatorIDs()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideoSummaries_Response(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		videos VideoSummaries
		expect []*response.VideoSummary
	}{
		{
			name: "success",
			videos: VideoSummaries{
				{
					VideoSummary: response.VideoSummary{
						ID:            "video-id",
						CoordinatorID: "coordinator-id",
						Title:         "じゃがいもの育て方",
						ThumbnailURL:  "https://example.com/thumbnail.jpg",
						PublishedAt:   0,
					},
				},
			},
			expect: []*response.VideoSummary{
				{
					ID:            "video-id",
					CoordinatorID: "coordinator-id",
					Title:         "じゃがいもの育て方",
					ThumbnailURL:  "https://example.com/thumbnail.jpg",
					PublishedAt:   0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.videos.Response()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
