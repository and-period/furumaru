package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/set"
)

type LiveSummary struct {
	types.LiveSummary
}

type LiveSummaries []*LiveSummary

func NewLiveSummary(schedule *entity.Schedule, products entity.Products) *LiveSummary {
	return &LiveSummary{
		LiveSummary: types.LiveSummary{
			ScheduleID:    schedule.ID,
			CoordinatorID: schedule.CoordinatorID,
			Status:        int32(NewScheduleStatus(schedule.Status, false)),
			Title:         schedule.Title,
			ThumbnailURL:  schedule.ThumbnailURL,
			StartAt:       schedule.StartAt.Unix(),
			EndAt:         schedule.EndAt.Unix(),
			Products:      NewLiveProducts(products).Response(),
		},
	}
}

func (l *LiveSummary) Response() *types.LiveSummary {
	return &l.LiveSummary
}

func NewLiveSummaries(schedules entity.Schedules, lives entity.Lives, products entity.Products) LiveSummaries {
	livesMap := lives.GroupByScheduleID()
	res := make(LiveSummaries, len(schedules))
	for i := range schedules {
		productIDs := livesMap[schedules[i].ID].ProductIDs()
		res[i] = NewLiveSummary(schedules[i], products.Filter(productIDs...))
	}
	return res
}

func (ls LiveSummaries) CoordinatorIDs() []string {
	return set.UniqBy(ls, func(l *LiveSummary) string {
		return l.CoordinatorID
	})
}

func (ls LiveSummaries) Response() []*types.LiveSummary {
	res := make([]*types.LiveSummary, len(ls))
	for i := range ls {
		res[i] = ls[i].Response()
	}
	return res
}
