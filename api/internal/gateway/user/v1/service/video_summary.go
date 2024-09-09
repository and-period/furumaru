package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/set"
)

type VideoSummary struct {
	response.VideoSummary
}

type VideoSummaries []*VideoSummary

func NewVideoSummary(video *entity.Video) *VideoSummary {
	return &VideoSummary{
		VideoSummary: response.VideoSummary{
			ID:            video.ID,
			CoordinatorID: video.CoordinatorID,
			Title:         video.Title,
			ThumbnailURL:  video.ThumbnailURL,
			PublishedAt:   video.PublishedAt.Unix(),
		},
	}
}

func (s *VideoSummary) Response() *response.VideoSummary {
	return &s.VideoSummary
}

func NewVideoSummaries(videos entity.Videos) VideoSummaries {
	res := make(VideoSummaries, len(videos))
	for i := range videos {
		res[i] = NewVideoSummary(videos[i])
	}
	return res
}

func (ss VideoSummaries) CoordinatorIDs() []string {
	return set.UniqBy(ss, func(s *VideoSummary) string {
		return s.CoordinatorID
	})
}

func (ss VideoSummaries) Response() []*response.VideoSummary {
	res := make([]*response.VideoSummary, len(ss))
	for i := range ss {
		res[i] = ss[i].Response()
	}
	return res
}
