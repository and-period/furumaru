package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type ArchiveSummary struct {
	response.ArchiveSummary
}

type ArchiveSummaries []*ArchiveSummary

func NewArchiveSummary(schedule *entity.Schedule) *ArchiveSummary {
	return &ArchiveSummary{
		ArchiveSummary: response.ArchiveSummary{
			ScheduleID:    schedule.ID,
			CoordinatorID: schedule.CoordinatorID,
			Title:         schedule.Title,
			ThumbnailURL:  schedule.ThumbnailURL,
			Thumbnails:    NewImages(schedule.Thumbnails).Response(),
		},
	}
}

func (a *ArchiveSummary) Response() *response.ArchiveSummary {
	return &a.ArchiveSummary
}

func NewArchiveSummaries(schedules entity.Schedules) ArchiveSummaries {
	res := make(ArchiveSummaries, len(schedules))
	for i := range schedules {
		res[i] = NewArchiveSummary(schedules[i])
	}
	return res
}

func (as ArchiveSummaries) Response() []*response.ArchiveSummary {
	res := make([]*response.ArchiveSummary, len(as))
	for i := range as {
		res[i] = as[i].Response()
	}
	return res
}
