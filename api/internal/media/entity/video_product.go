package entity

import (
	"sort"
	"time"
)

// オンデマンド配信関連商品情報
type VideoProduct struct {
	VideoID   string    `gorm:"primaryKey;<-:create"` // オンデマンド動画ID
	ProductID string    `gorm:"primaryKey;<-:create"` // 商品ID
	Priority  int64     `gorm:"default:0"`            // 表示優先度
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type VideoProducts []*VideoProduct

func (ps VideoProducts) ProductIDs() []string {
	res := make([]string, len(ps))
	for i := range ps {
		res[i] = ps[i].ProductID
	}
	return res
}

func (ps VideoProducts) GroupByVideoID() map[string]VideoProducts {
	res := make(map[string]VideoProducts, len(ps))
	for _, p := range ps {
		if _, ok := res[p.VideoID]; !ok {
			res[p.VideoID] = make(VideoProducts, 0, len(ps))
		}
		res[p.VideoID] = append(res[p.VideoID], p)
	}
	return res
}

func (ps VideoProducts) SortByPriority() VideoProducts {
	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].Priority < ps[j].Priority
	})
	return ps
}
