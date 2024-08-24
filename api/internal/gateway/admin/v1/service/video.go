package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

// VideoStatus - オンデマンド配信状況
type VideoStatus int32

const (
	VideoStatusUnknown   VideoStatus = 0
	VideoStatusPrivate   VideoStatus = 1 // 非公開
	VideoStatusWaiting   VideoStatus = 2 // 公開前
	VideoStatusLimited   VideoStatus = 3 // 限定公開
	VideoStatusPublished VideoStatus = 4 // 公開済み
)

type Video struct {
	response.Video
}

type Videos []*Video

func NewVideoStatus(status entity.VideoStatus) VideoStatus {
	switch status {
	case entity.VideoStatusPrivate:
		return VideoStatusPrivate
	case entity.VideoStatusWaiting:
		return VideoStatusWaiting
	case entity.VideoStatusLimited:
		return VideoStatusLimited
	case entity.VideoStatusPublished:
		return VideoStatusPublished
	default:
		return VideoStatusUnknown
	}
}

func (s VideoStatus) Response() int32 {
	return int32(s)
}

func NewVideo(video *entity.Video) *Video {
	return &Video{
		Video: response.Video{
			ID:            video.ID,
			CoordinatorID: video.CoordinatorID,
			ProductIDs:    video.ProductIDs,
			ExperienceIDs: video.ExperienceIDs,
			Title:         video.Title,
			Description:   video.Description,
			Status:        NewVideoStatus(video.Status).Response(),
			ThumbnailURL:  video.ThumbnailURL,
			VideoURL:      video.VideoURL,
			Public:        video.Public,
			Limited:       video.Limited,
			PublishedAt:   video.PublishedAt.Unix(),
			CreatedAt:     video.CreatedAt.Unix(),
			UpdatedAt:     video.UpdatedAt.Unix(),
		},
	}
}

func (v *Video) Response() *response.Video {
	return &v.Video
}

func NewVideos(videos entity.Videos) Videos {
	res := make(Videos, len(videos))
	for i := range videos {
		res[i] = NewVideo(videos[i])
	}
	return res
}

func (vs Videos) Response() []*response.Video {
	res := make([]*response.Video, len(vs))
	for i := range vs {
		res[i] = vs[i].Response()
	}
	return res
}
