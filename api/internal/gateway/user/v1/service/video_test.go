package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestVideos(t *testing.T) {
	t.Parallel()

	now := jst.Date(2024, 8, 24, 18, 30, 0, 0)

	tests := []struct {
		name   string
		videos entity.Videos
		expect Videos
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
			expect: Videos{
				{
					Video: response.Video{
						ID:            "video-id",
						CoordinatorID: "coordinator-id",
						ProductIDs:    []string{"product-id"},
						ExperienceIDs: []string{"experience-id"},
						Title:         "じゃがいもの育て方",
						Description:   "じゃがいもの育て方の動画です。",
						ThumbnailURL:  "https://example.com/thumbnail.jpg",
						VideoURL:      "https://example.com/video.mp4",
						PublishedAt:   now.AddDate(0, 0, -1).Unix(),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewVideos(tt.videos)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideos_Response(t *testing.T) {
	t.Parallel()

	now := jst.Date(2024, 8, 24, 18, 30, 0, 0)

	tests := []struct {
		name   string
		videos Videos
		expect []*response.Video
	}{
		{
			name: "success",
			videos: Videos{
				{
					Video: response.Video{
						ID:            "video-id",
						CoordinatorID: "coordinator-id",
						ProductIDs:    []string{"product-id"},
						ExperienceIDs: []string{"experience-id"},
						Title:         "じゃがいもの育て方",
						Description:   "じゃがいもの育て方の動画です。",
						ThumbnailURL:  "https://example.com/thumbnail.jpg",
						VideoURL:      "https://example.com/video.mp4",
						PublishedAt:   now.AddDate(0, 0, -1).Unix(),
					},
				},
			},
			expect: []*response.Video{
				{
					ID:            "video-id",
					CoordinatorID: "coordinator-id",
					ProductIDs:    []string{"product-id"},
					ExperienceIDs: []string{"experience-id"},
					Title:         "じゃがいもの育て方",
					Description:   "じゃがいもの育て方の動画です。",
					ThumbnailURL:  "https://example.com/thumbnail.jpg",
					VideoURL:      "https://example.com/video.mp4",
					PublishedAt:   now.AddDate(0, 0, -1).Unix(),
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
