package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

// ScheduleStatus - 開催状況
type ScheduleStatus int32

const (
	ScheduleStatusUnknown    ScheduleStatus = 0
	ScheduleStatusPrivate    ScheduleStatus = 1 // 非公開
	ScheduleStatusInProgress ScheduleStatus = 2 // 管理者承認前
	ScheduleStatusWaiting    ScheduleStatus = 3 // 開催前
	ScheduleStatusLive       ScheduleStatus = 4 // 開催中
	ScheduleStatusClosed     ScheduleStatus = 5 // 開催終了
)

func NewScheduleStatus(status entity.ScheduleStatus) ScheduleStatus {
	switch status {
	case entity.ScheduleStatusPrivate:
		return ScheduleStatusPrivate
	case entity.ScheduleStatusInProgress:
		return ScheduleStatusInProgress
	case entity.ScheduleStatusWaiting:
		return ScheduleStatusWaiting
	case entity.ScheduleStatusLive:
		return ScheduleStatusLive
	case entity.ScheduleStatusClosed:
		return ScheduleStatusClosed
	default:
		return ScheduleStatusUnknown
	}
}

func (s ScheduleStatus) Response() int32 {
	return int32(s)
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
