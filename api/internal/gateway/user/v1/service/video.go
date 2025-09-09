package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

type Video struct {
	types.Video
}

type Videos []*Video

func NewVideo(v *entity.Video) *Video {
	return &Video{
		Video: types.Video{
			ID:            v.ID,
			CoordinatorID: v.CoordinatorID,
			ProductIDs:    v.ProductIDs,
			ExperienceIDs: v.ExperienceIDs,
			Title:         v.Title,
			Description:   v.Description,
			ThumbnailURL:  v.ThumbnailURL,
			VideoURL:      v.VideoURL,
			PublishedAt:   v.PublishedAt.Unix(),
		},
	}
}

func (v *Video) Response() *types.Video {
	return &v.Video
}

func NewVideos(videos []*entity.Video) Videos {
	result := make(Videos, len(videos))
	for i, v := range videos {
		result[i] = NewVideo(v)
	}
	return result
}

func (vs Videos) Response() []*types.Video {
	result := make([]*types.Video, len(vs))
	for i, v := range vs {
		result[i] = v.Response()
	}
	return result
}
