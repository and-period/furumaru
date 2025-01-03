package entity

import (
	"sort"
	"time"
)

// オンデマンド配信関連体験情報
type VideoExperience struct {
	VideoID      string    `gorm:"primaryKey;<-:create"` // オンデマンド動画ID
	ExperienceID string    `gorm:"primaryKey;<-:create"` // 体験ID
	Priority     int64     `gorm:"default:0"`            // 表示優先度
	CreatedAt    time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt    time.Time `gorm:""`                     // 更新日時
}

type VideoExperiences []*VideoExperience

type NewVideoExperienceParams struct {
	VideoID      string
	ExperienceID string
	Priority     int64
}

func NewVideoExperience(params *NewVideoExperienceParams) *VideoExperience {
	return &VideoExperience{
		VideoID:      params.VideoID,
		ExperienceID: params.ExperienceID,
		Priority:     params.Priority,
	}
}

func NewVideoExperiences(videoID string, experienceIDs []string) VideoExperiences {
	res := make(VideoExperiences, len(experienceIDs))
	for i := range experienceIDs {
		params := &NewVideoExperienceParams{
			VideoID:      videoID,
			ExperienceID: experienceIDs[i],
			Priority:     int64(i + 1),
		}
		res[i] = NewVideoExperience(params)
	}
	return res
}

func (es VideoExperiences) ExperienceIDs() []string {
	res := make([]string, len(es))
	for i := range es {
		res[i] = es[i].ExperienceID
	}
	return res
}

func (es VideoExperiences) GroupByVideoID() map[string]VideoExperiences {
	res := make(map[string]VideoExperiences, len(es))
	for _, e := range es {
		if _, ok := res[e.VideoID]; !ok {
			res[e.VideoID] = make(VideoExperiences, 0, len(es))
		}
		res[e.VideoID] = append(res[e.VideoID], e)
	}
	return res
}

func (es VideoExperiences) SortByPriority() VideoExperiences {
	sort.SliceStable(es, func(i, j int) bool {
		return es[i].Priority < es[j].Priority
	})
	return es
}
