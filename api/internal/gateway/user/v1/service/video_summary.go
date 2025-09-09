package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/set"
)

type VideoSummary struct {
	types.VideoSummary
}

type VideoSummaries []*VideoSummary

func NewVideoSummary(video *entity.Video) *VideoSummary {
	return &VideoSummary{
		VideoSummary: types.VideoSummary{
			ID:            video.ID,
			CoordinatorID: video.CoordinatorID,
			Title:         video.Title,
			ThumbnailURL:  video.ThumbnailURL,
			PublishedAt:   video.PublishedAt.Unix(),
		},
	}
}

func (s *VideoSummary) Response() *types.VideoSummary {
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

func (ss VideoSummaries) Response() []*types.VideoSummary {
	res := make([]*types.VideoSummary, len(ss))
	for i := range ss {
		res[i] = ss[i].Response()
	}
	return res
}
