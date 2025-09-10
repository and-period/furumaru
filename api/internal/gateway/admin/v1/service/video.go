package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

// VideoStatus - オンデマンド配信状況
type VideoStatus types.VideoStatus

type Video struct {
	types.Video
}

type Videos []*Video

func NewVideoStatus(status entity.VideoStatus) VideoStatus {
	switch status {
	case entity.VideoStatusPrivate:
		return VideoStatus(types.VideoStatusPrivate)
	case entity.VideoStatusWaiting:
		return VideoStatus(types.VideoStatusWaiting)
	case entity.VideoStatusLimited:
		return VideoStatus(types.VideoStatusLimited)
	case entity.VideoStatusPublished:
		return VideoStatus(types.VideoStatusPublished)
	default:
		return VideoStatus(types.VideoStatusUnknown)
	}
}

func (s VideoStatus) Response() types.VideoStatus {
	return types.VideoStatus(s)
}

func NewVideo(video *entity.Video) *Video {
	return &Video{
		Video: types.Video{
			ID:                video.ID,
			CoordinatorID:     video.CoordinatorID,
			ProductIDs:        video.ProductIDs,
			ExperienceIDs:     video.ExperienceIDs,
			Title:             video.Title,
			Description:       video.Description,
			Status:            NewVideoStatus(video.Status).Response(),
			ThumbnailURL:      video.ThumbnailURL,
			VideoURL:          video.VideoURL,
			Public:            video.Public,
			Limited:           video.Limited,
			DisplayProduct:    video.DisplayProduct,
			DisplayExperience: video.DisplayExperience,
			PublishedAt:       video.PublishedAt.Unix(),
			CreatedAt:         video.CreatedAt.Unix(),
			UpdatedAt:         video.UpdatedAt.Unix(),
		},
	}
}

func (v *Video) Response() *types.Video {
	return &v.Video
}

func NewVideos(videos entity.Videos) Videos {
	res := make(Videos, len(videos))
	for i := range videos {
		res[i] = NewVideo(videos[i])
	}
	return res
}

func (vs Videos) Response() []*types.Video {
	res := make([]*types.Video, len(vs))
	for i := range vs {
		res[i] = vs[i].Response()
	}
	return res
}
