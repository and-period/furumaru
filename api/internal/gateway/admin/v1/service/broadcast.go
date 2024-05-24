package service

import (
	"fmt"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

const youtubeAdminURL = "https://studio.youtube.com/video/%s/livestreaming"

// BroadcastStatus - ライブ配信状況
type BroadcastStatus int32

const (
	BroadcastStatusUnknown  BroadcastStatus = 0
	BroadcastStatusDisabled BroadcastStatus = 1 // リソース未作成
	BroadcastStatusWaiting  BroadcastStatus = 2 // リソース作成/削除中
	BroadcastStatusIdle     BroadcastStatus = 3 // 停止中
	BroadcastStatusActive   BroadcastStatus = 4 // 配信中
)

type Broadcast struct {
	response.Broadcast
}

type Broadcasts []*Broadcast

type GuestBroadcast struct {
	response.GuestBroadcast
}

func NewBroadcastStatus(status entity.BroadcastStatus) BroadcastStatus {
	switch status {
	case entity.BroadcastStatusDisabled:
		return BroadcastStatusDisabled
	case entity.BroadcastStatusWaiting:
		return BroadcastStatusWaiting
	case entity.BroadcastStatusIdle:
		return BroadcastStatusIdle
	case entity.BroadcastStatusActive:
		return BroadcastStatusActive
	default:
		return BroadcastStatusUnknown
	}
}

func (s BroadcastStatus) Response() int32 {
	return int32(s)
}

func NewBroadcast(broadcast *entity.Broadcast) *Broadcast {
	res := &Broadcast{
		Broadcast: response.Broadcast{
			ID:             broadcast.ID,
			ScheduleID:     broadcast.ScheduleID,
			Status:         NewBroadcastStatus(broadcast.Status).Response(),
			InputURL:       broadcast.InputURL,
			OutputURL:      broadcast.OutputURL,
			ArchiveURL:     broadcast.ArchiveURL,
			YoutubeAccount: broadcast.YoutubeAccount,
			CreatedAt:      broadcast.CreatedAt.Unix(),
			UpdatedAt:      broadcast.CreatedAt.Unix(),
		},
	}
	if broadcast.YoutubeBroadcastID != "" {
		res.YoutubeAdminURL = fmt.Sprintf(youtubeAdminURL, broadcast.YoutubeBroadcastID)
	}
	return res
}

func (b *Broadcast) Response() *response.Broadcast {
	if b == nil {
		return nil
	}
	return &b.Broadcast
}

func NewBroadcasts(broadcasts entity.Broadcasts) Broadcasts {
	res := make(Broadcasts, len(broadcasts))
	for i := range broadcasts {
		res[i] = NewBroadcast(broadcasts[i])
	}
	return res
}

func (bs Broadcasts) Response() []*response.Broadcast {
	res := make([]*response.Broadcast, len(bs))
	for i := range bs {
		res[i] = bs[i].Response()
	}
	return res
}

func NewGuestBroadcast(schedule *Schedule, coordinator *Coordinator) *GuestBroadcast {
	return &GuestBroadcast{
		response.GuestBroadcast{
			Title:             schedule.Title,
			Description:       schedule.Description,
			StartAt:           schedule.StartAt,
			EndAt:             schedule.EndAt,
			CoordinatorMarche: coordinator.MarcheName,
			CoordinatorName:   coordinator.Username,
		},
	}
}

func (b *GuestBroadcast) Response() *response.GuestBroadcast {
	return &b.GuestBroadcast
}
