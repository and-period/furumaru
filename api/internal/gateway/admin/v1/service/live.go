package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

// 配信ステータス
type LiveStatus int32

const (
	LiveStatusUnknown  LiveStatus = 0
	LiveStatusWaiting  LiveStatus = 1 // 開始前
	LiveStatusOpened   LiveStatus = 2 // 配信中
	LiveStatusClosed   LiveStatus = 3 // 配信終了
	LiveStatusCanceled LiveStatus = 4 // 配信中止
)

type Live struct {
	response.Live
	productIDs []string
}

type Lives []*Live

func NewLiveStatus(sta entity.LiveStatus) LiveStatus {
	switch sta {
	case entity.LiveStatusWaiting:
		return LiveStatusWaiting
	case entity.LiveStatusOpened:
		return LiveStatusOpened
	case entity.LiveStatusClosed:
		return LiveStatusClosed
	case entity.LiveStatusCanceled:
		return LiveStatusCanceled
	default:
		return LiveStatusUnknown
	}
}

func (s LiveStatus) Response() int32 {
	return int32(s)
}

func NewLive(live *entity.Live) *Live {
	return &Live{
		Live: response.Live{
			ID:          live.ID,
			ScheduleID:  live.ScheduleID,
			Title:       live.Title,
			Description: live.Description,
			ProducerID:  live.ProducerID,
			StartAt:     live.StartAt.Unix(),
			EndAt:       live.EndAt.Unix(),
			Canceled:    live.Canceled,
			Status:      NewLiveStatus(live.Status).Response(),
			CreatedAt:   live.CreatedAt.Unix(),
			UpdatedAt:   live.UpdatedAt.Unix(),
		},
	}
}

func (l *Live) Fill(producer *Producer, products map[string]*Product) {
	if producer != nil {
		l.ProducerName = producer.Name()
	}
	ps := make(Products, len(l.productIDs))
	for i := range l.productIDs {
		if products[l.productIDs[i]] != nil {
			ps[i] = products[l.productIDs[i]]
		}
	}
	l.Products = ps.Response()
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

func (ls Lives) Fill(
	producers Producers,
	products map[string]*Product,
) {
	for i := range ls {
		ls[i].Fill(producers[i], products)
	}
}

func (ls Lives) Response() []*response.Live {
	res := make([]*response.Live, len(ls))
	for i := range ls {
		res[i] = ls[i].Response()
	}
	return res
}
