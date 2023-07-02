package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
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
	response.Schedule
}

type Schedules []*Schedule

func NewSchedule(schedule *entity.Schedule) *Schedule {
	return &Schedule{
		Schedule: response.Schedule{
			ID:                   schedule.ID,
			CoordinatorID:        schedule.CoordinatorID,
			CoordinatorName:      "",
			ShippingID:           schedule.ShippingID,
			ShippingName:         "",
			Status:               NewScheduleStatus(schedule.Status).Response(),
			Title:                schedule.Title,
			Description:          schedule.Description,
			ThumbnailURL:         schedule.ThumbnailURL,
			OpeningVideoURL:      schedule.OpeningVideoURL,
			IntermissionVideoURL: schedule.IntermissionVideoURL,
			Public:               schedule.Public,
			Approved:             schedule.Approved,
			StartAt:              schedule.StartAt.Unix(),
			EndAt:                schedule.EndAt.Unix(),
			CreatedAt:            schedule.CreatedAt.Unix(),
			UpdatedAt:            schedule.UpdatedAt.Unix(),
		},
	}
}

func (s *Schedule) Fill(shipping *Shipping, coordinator *Coordinator) {
	if shipping != nil {
		s.ShippingName = shipping.Name
	}
	if coordinator != nil {
		s.CoordinatorName = coordinator.Username
	}
}

func (s *Schedule) Response() *response.Schedule {
	return &s.Schedule
}

func NewSchedules(schedules entity.Schedules) Schedules {
	res := make(Schedules, len(schedules))
	for i := range schedules {
		res[i] = NewSchedule(schedules[i])
	}
	return res
}

func (ss Schedules) Fill(shippings map[string]*Shipping, coordinators map[string]*Coordinator) {
	for _, s := range ss {
		s.Fill(shippings[s.ShippingID], coordinators[s.CoordinatorID])
	}
}

func (ss Schedules) Response() []*response.Schedule {
	res := make([]*response.Schedule, len(ss))
	for i := range ss {
		res[i] = ss[i].Response()
	}
	return res
}
