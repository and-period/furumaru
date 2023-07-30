package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

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
	return &Broadcast{
		Broadcast: response.Broadcast{
			ID:         broadcast.ID,
			ScheduleID: broadcast.ScheduleID,
			Status:     NewBroadcastStatus(broadcast.Status).Response(),
			InputURL:   broadcast.InputURL,
			OutputURL:  broadcast.OutputURL,
			CreatedAt:  broadcast.CreatedAt.Unix(),
			UpdatedAt:  broadcast.CreatedAt.Unix(),
		},
	}
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
