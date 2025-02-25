package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestVideoStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status entity.VideoStatus
		expect VideoStatus
	}{
		{
			name:   "private",
			status: entity.VideoStatusPrivate,
			expect: VideoStatusPrivate,
		},
		{
			name:   "waiting",
			status: entity.VideoStatusWaiting,
			expect: VideoStatusWaiting,
		},
		{
			name:   "limited",
			status: entity.VideoStatusLimited,
			expect: VideoStatusLimited,
		},
		{
			name:   "published",
			status: entity.VideoStatusPublished,
			expect: VideoStatusPublished,
		},
		{
			name:   "unknown",
			status: entity.VideoStatusUnknown,
			expect: VideoStatusUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewVideoStatus(tt.status)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideoStatus_Response(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status VideoStatus
		expect int32
	}{
		{
			name:   "private",
			status: VideoStatusPrivate,
			expect: 1,
		},
		{
			name:   "waiting",
			status: VideoStatusWaiting,
			expect: 2,
		},
		{
			name:   "limited",
			status: VideoStatusLimited,
			expect: 3,
		},
		{
			name:   "published",
			status: VideoStatusPublished,
			expect: 4,
		},
		{
			name:   "unknown",
			status: VideoStatusUnknown,
			expect: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.status.Response()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideos(t *testing.T) {
	t.Parallel()

	now := jst.Date(2024, 8, 24, 18, 30, 0, 0)

	tests := []struct {
		name   string
		videos entity.Videos
		expect Videos
	}{
		{
			name: "name",
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
						ID:                "video-id",
						CoordinatorID:     "coordinator-id",
						ProductIDs:        []string{"product-id"},
						ExperienceIDs:     []string{"experience-id"},
						Title:             "じゃがいもの育て方",
						Description:       "じゃがいもの育て方の動画です。",
						Status:            int32(VideoStatusPublished),
						ThumbnailURL:      "https://example.com/thumbnail.jpg",
						VideoURL:          "https://example.com/video.mp4",
						Public:            true,
						Limited:           false,
						DisplayProduct:    true,
						DisplayExperience: true,
						PublishedAt:       now.AddDate(0, 0, -1).Unix(),
						CreatedAt:         now.Unix(),
						UpdatedAt:         now.Unix(),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			name: "name",
			videos: Videos{
				{
					Video: response.Video{
						ID:                "video-id",
						CoordinatorID:     "coordinator-id",
						ProductIDs:        []string{"product-id"},
						ExperienceIDs:     []string{"experience-id"},
						Title:             "じゃがいもの育て方",
						Description:       "じゃがいもの育て方の動画です。",
						Status:            int32(VideoStatusPublished),
						ThumbnailURL:      "https://example.com/thumbnail.jpg",
						VideoURL:          "https://example.com/video.mp4",
						Public:            true,
						Limited:           false,
						DisplayProduct:    true,
						DisplayExperience: true,
						PublishedAt:       now.AddDate(0, 0, -1).Unix(),
						CreatedAt:         now.Unix(),
						UpdatedAt:         now.Unix(),
					},
				},
			},
			expect: []*response.Video{
				{
					ID:                "video-id",
					CoordinatorID:     "coordinator-id",
					ProductIDs:        []string{"product-id"},
					ExperienceIDs:     []string{"experience-id"},
					Title:             "じゃがいもの育て方",
					Description:       "じゃがいもの育て方の動画です。",
					Status:            int32(VideoStatusPublished),
					ThumbnailURL:      "https://example.com/thumbnail.jpg",
					VideoURL:          "https://example.com/video.mp4",
					Public:            true,
					Limited:           false,
					DisplayProduct:    true,
					DisplayExperience: true,
					PublishedAt:       now.AddDate(0, 0, -1).Unix(),
					CreatedAt:         now.Unix(),
					UpdatedAt:         now.Unix(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.videos.Response())
		})
	}
}
