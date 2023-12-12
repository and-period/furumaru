package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type Live struct {
	response.Live
}

type Lives []*Live

func NewLive(live *entity.Live) *Live {
	return &Live{
		Live: response.Live{
			ScheduleID: live.ScheduleID,
			ProducerID: live.ProducerID,
			ProductIDs: live.ProductIDs,
			Comment:    live.Comment,
			StartAt:    live.StartAt.Unix(),
			EndAt:      live.EndAt.Unix(),
		},
	}
}

func (l *Live) Response() *response.Live {
	return &l.Live
}

func NewLives(lives entity.Lives) Lives {
	res := make(Lives, len(lives))
	for i := range lives {
		res[i] = NewLive(lives[i])
	}
	return res
}

func (ls Lives) Response() []*response.Live {
	res := make([]*response.Live, len(ls))
	for i := range ls {
		res[i] = ls[i].Response()
	}
	return res
}
