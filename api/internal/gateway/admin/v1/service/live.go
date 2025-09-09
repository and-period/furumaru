package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type Live struct {
	types.Live
}

type Lives []*Live

func NewLive(live *entity.Live) *Live {
	return &Live{
		Live: types.Live{
			ID:         live.ID,
			ScheduleID: live.ScheduleID,
			ProducerID: live.ProducerID,
			ProductIDs: live.ProductIDs,
			Comment:    live.Comment,
			StartAt:    live.StartAt.Unix(),
			EndAt:      live.EndAt.Unix(),
			CreatedAt:  live.CreatedAt.Unix(),
			UpdatedAt:  live.UpdatedAt.Unix(),
		},
	}
}

func (l *Live) Response() *types.Live {
	return &l.Live
}

func NewLives(lives entity.Lives) Lives {
	res := make(Lives, len(lives))
	for i := range lives {
		res[i] = NewLive(lives[i])
	}
	return res
}

func (ls Lives) Response() []*types.Live {
	res := make([]*types.Live, len(ls))
	for i := range ls {
		res[i] = ls[i].Response()
	}
	return res
}
