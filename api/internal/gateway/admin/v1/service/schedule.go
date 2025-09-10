package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

// ScheduleStatus - 開催状況
type ScheduleStatus types.ScheduleStatus

func NewScheduleStatus(status entity.ScheduleStatus) ScheduleStatus {
	switch status {
	case entity.ScheduleStatusPrivate:
		return ScheduleStatus(types.ScheduleStatusPrivate)
	case entity.ScheduleStatusInProgress:
		return ScheduleStatus(types.ScheduleStatusInProgress)
	case entity.ScheduleStatusWaiting:
		return ScheduleStatus(types.ScheduleStatusWaiting)
	case entity.ScheduleStatusLive:
		return ScheduleStatus(types.ScheduleStatusLive)
	case entity.ScheduleStatusClosed:
		return ScheduleStatus(types.ScheduleStatusClosed)
	default:
		return ScheduleStatus(types.ScheduleStatusUnknown)
	}
}

func (s ScheduleStatus) Response() types.ScheduleStatus {
	return types.ScheduleStatus(s)
}

type Schedule struct {
	types.Schedule
}

type Schedules []*Schedule

func NewSchedule(schedule *entity.Schedule) *Schedule {
	return &Schedule{
		Schedule: types.Schedule{
			ID:              schedule.ID,
			ShopID:          schedule.ShopID,
			CoordinatorID:   schedule.CoordinatorID,
			Status:          NewScheduleStatus(schedule.Status).Response(),
			Title:           schedule.Title,
			Description:     schedule.Description,
			ThumbnailURL:    schedule.ThumbnailURL,
			ImageURL:        schedule.ImageURL,
			OpeningVideoURL: schedule.OpeningVideoURL,
			Public:          schedule.Public,
			Approved:        schedule.Approved,
			StartAt:         schedule.StartAt.Unix(),
			EndAt:           schedule.EndAt.Unix(),
			CreatedAt:       schedule.CreatedAt.Unix(),
			UpdatedAt:       schedule.UpdatedAt.Unix(),
		},
	}
}

func (s *Schedule) Response() *types.Schedule {
	return &s.Schedule
}

func NewSchedules(schedules entity.Schedules) Schedules {
	res := make(Schedules, len(schedules))
	for i := range schedules {
		res[i] = NewSchedule(schedules[i])
	}
	return res
}

func (ss Schedules) Response() []*types.Schedule {
	res := make([]*types.Schedule, len(ss))
	for i := range ss {
		res[i] = ss[i].Response()
	}
	return res
}
