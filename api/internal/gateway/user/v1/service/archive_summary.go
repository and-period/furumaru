package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/set"
)

type ArchiveSummary struct {
	types.ArchiveSummary
}

type ArchiveSummaries []*ArchiveSummary

func NewArchiveSummary(schedule *entity.Schedule) *ArchiveSummary {
	return &ArchiveSummary{
		ArchiveSummary: types.ArchiveSummary{
			ScheduleID:    schedule.ID,
			CoordinatorID: schedule.CoordinatorID,
			Title:         schedule.Title,
			StartAt:       schedule.StartAt.Unix(),
			EndAt:         schedule.EndAt.Unix(),
			ThumbnailURL:  schedule.ThumbnailURL,
		},
	}
}

func (a *ArchiveSummary) Response() *types.ArchiveSummary {
	return &a.ArchiveSummary
}

func NewArchiveSummaries(schedules entity.Schedules) ArchiveSummaries {
	res := make(ArchiveSummaries, len(schedules))
	for i := range schedules {
		res[i] = NewArchiveSummary(schedules[i])
	}
	return res
}

func (as ArchiveSummaries) CoordinatorIDs() []string {
	return set.UniqBy(as, func(a *ArchiveSummary) string {
		return a.CoordinatorID
	})
}

func (as ArchiveSummaries) Response() []*types.ArchiveSummary {
	res := make([]*types.ArchiveSummary, len(as))
	for i := range as {
		res[i] = as[i].Response()
	}
	return res
}
