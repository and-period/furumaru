package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

type Video struct {
	response.Video
}

type Videos []*Video

func NewVideo(v *entity.Video) *Video {
	return &Video{
		Video: response.Video{
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

func (v *Video) Response() *response.Video {
	return &v.Video
}

func NewVideos(videos []*entity.Video) Videos {
	result := make(Videos, len(videos))
	for i, v := range videos {
		result[i] = NewVideo(v)
	}
	return result
}

func (vs Videos) Response() []*response.Video {
	result := make([]*response.Video, len(vs))
	for i, v := range vs {
		result[i] = v.Response()
	}
	return result
}
