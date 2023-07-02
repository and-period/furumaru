package entity

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestSchedule(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewScheduleParams
		expect *Schedule
	}{
		{
			name: "success",
			params: &NewScheduleParams{
				CoordinatorID:        "coordinator-id",
				ShippingID:           "shipping-id",
				Title:                "スケジュールタイトル",
				Description:          "スケジュールの詳細です。",
				ThumbnailURL:         "https://and-period.jp/thumbnail.png",
				OpeningVideoURL:      "https://and-period.jp/opening-video.mp4",
				IntermissionVideoURL: "https://and-period.jp/intermission-video.mp4",
				StartAt:              jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:                jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
			expect: &Schedule{
				ID:                   "",
				CoordinatorID:        "coordinator-id",
				ShippingID:           "shipping-id",
				Title:                "スケジュールタイトル",
				Description:          "スケジュールの詳細です。",
				ThumbnailURL:         "https://and-period.jp/thumbnail.png",
				OpeningVideoURL:      "https://and-period.jp/opening-video.mp4",
				IntermissionVideoURL: "https://and-period.jp/intermission-video.mp4",
				Approved:             false,
				ApprovedAdminID:      "",
				StartAt:              jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:                jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewSchedule(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestSchedule_Fill(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		schedule *Schedule
		expect   ScheduleStatus
	}{
		{
			name: "in progress",
			schedule: &Schedule{
				Public:   false,
				Approved: false,
				StartAt:  time.Time{},
				EndAt:    time.Time{},
			},
			expect: ScheduleStatusInProgress,
		},
		{
			name: "private",
			schedule: &Schedule{
				Public:   false,
				Approved: true,
				StartAt:  time.Time{},
				EndAt:    time.Time{},
			},
			expect: ScheduleStatusPrivate,
		},
		{
			name: "waiting",
			schedule: &Schedule{
				Public:   true,
				Approved: true,
				StartAt:  now.AddDate(0, 1, 0),
				EndAt:    now.AddDate(0, 1, 0),
			},
			expect: ScheduleStatusWaiting,
		},
		{
			name: "live",
			schedule: &Schedule{
				Public:   true,
				Approved: true,
				StartAt:  now.AddDate(0, -1, 0),
				EndAt:    now.AddDate(0, 1, 0),
			},
			expect: ScheduleStatusLive,
		},
		{
			name: "closed",
			schedule: &Schedule{
				Public:   true,
				Approved: true,
				StartAt:  now.AddDate(0, -1, 0),
				EndAt:    now.AddDate(0, -1, 0),
			},
			expect: ScheduleStatusClosed,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.schedule.Fill(now)
			assert.Equal(t, tt.expect, tt.schedule.Status)
		})
	}
}
