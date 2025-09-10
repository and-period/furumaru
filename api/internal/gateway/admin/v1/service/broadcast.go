package service

import (
	"fmt"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

const (
	YoutubeViewerURL = "https://youtube.com/live/%s"
	youtubeAdminURL  = "https://studio.youtube.com/video/%s/livestreaming"
)

// BroadcastStatus - ライブ配信状況
type BroadcastStatus types.BroadcastStatus

type Broadcast struct {
	types.Broadcast
}

type Broadcasts []*Broadcast

type GuestBroadcast struct {
	types.GuestBroadcast
}

func NewBroadcastStatus(status entity.BroadcastStatus) BroadcastStatus {
	switch status {
	case entity.BroadcastStatusDisabled:
		return BroadcastStatus(types.BroadcastStatusDisabled)
	case entity.BroadcastStatusWaiting:
		return BroadcastStatus(types.BroadcastStatusWaiting)
	case entity.BroadcastStatusIdle:
		return BroadcastStatus(types.BroadcastStatusIdle)
	case entity.BroadcastStatusActive:
		return BroadcastStatus(types.BroadcastStatusActive)
	default:
		return BroadcastStatus(types.BroadcastStatusUnknown)
	}
}

func (s BroadcastStatus) Response() types.BroadcastStatus {
	return types.BroadcastStatus(s)
}

func NewBroadcast(broadcast *entity.Broadcast) *Broadcast {
	res := &Broadcast{
		Broadcast: types.Broadcast{
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
		res.YoutubeViewerURL = fmt.Sprintf(YoutubeViewerURL, broadcast.YoutubeBroadcastID)
		res.YoutubeAdminURL = fmt.Sprintf(youtubeAdminURL, broadcast.YoutubeBroadcastID)
	}
	return res
}

func (b *Broadcast) Response() *types.Broadcast {
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

func (bs Broadcasts) Response() []*types.Broadcast {
	res := make([]*types.Broadcast, len(bs))
	for i := range bs {
		res[i] = bs[i].Response()
	}
	return res
}

func NewGuestBroadcast(schedule *Schedule, shop *Shop, coordinator *Coordinator) *GuestBroadcast {
	return &GuestBroadcast{
		types.GuestBroadcast{
			Title:             schedule.Title,
			Description:       schedule.Description,
			StartAt:           schedule.StartAt,
			EndAt:             schedule.EndAt,
			CoordinatorMarche: shop.Name,
			CoordinatorName:   coordinator.Username,
		},
	}
}

func (b *GuestBroadcast) Response() *types.GuestBroadcast {
	return &b.GuestBroadcast
}
