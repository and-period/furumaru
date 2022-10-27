package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type Schedule struct {
	response.Schedule
}

func NewSchedule(schedule *entity.Schedule) *Schedule {
	return &Schedule{
		Schedule: response.Schedule{
			ID:            schedule.ID,
			CoordinatorID: schedule.CoordinatorID,
			ShippingID:    schedule.ShippingID,
			Title:         schedule.Title,
			Description:   schedule.Description,
			ThumbnailURL:  schedule.ThumbnailURL,
			StartAt:       schedule.StartAt.Unix(),
			EndAt:         schedule.EndAt.Unix(),
			Canceled:      schedule.Canceled,
			CreatedAt:     schedule.CreatedAt.Unix(),
			UpdatedAt:     schedule.UpdatedAt.Unix(),
		},
	}
}

func (s *Schedule) Response() *response.Schedule {
	return &s.Schedule
}

func (s *Schedule) Fill(shipping *Shipping) {
	if shipping != nil {
		s.ShippingName = shipping.Name
	}
}
