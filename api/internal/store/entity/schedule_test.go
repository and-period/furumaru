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
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				Title:           "スケジュールタイトル",
				Description:     "スケジュールの詳細です。",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
				Public:          true,
				StartAt:         jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:           jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
			expect: &Schedule{
				ID:              "",
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				Title:           "スケジュールタイトル",
				Description:     "スケジュールの詳細です。",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
				Public:          true,
				Approved:        true,
				ApprovedAdminID: "",
				StartAt:         jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:           jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewSchedule(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestSchedule_SetStatus(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.schedule.SetStatus(now)
			assert.Equal(t, tt.expect, tt.schedule.Status)
		})
	}
}

func TestSchedule_Published(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		schedule *Schedule
		expect   bool
	}{
		{
			name: "published",
			schedule: &Schedule{
				Public:   true,
				Approved: true,
			},
			expect: true,
		},
		{
			name: "unpublished",
			schedule: &Schedule{
				Public:   true,
				Approved: false,
			},
			expect: false,
		},
		{
			name:     "nil",
			schedule: nil,
			expect:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.schedule.Published())
		})
	}
}

func TestSchedules_FIll(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 8, 15, 0, 0, 0, 0)
	tests := []struct {
		name      string
		schedules Schedules
		expect    Schedules
		hasErr    bool
	}{
		{
			name: "success",
			schedules: Schedules{
				{
					ID:              "schedule-id",
					CoordinatorID:   "coordinator-id",
					Title:           "スケジュールタイトル",
					Description:     "スケジュールの詳細です。",
					ThumbnailURL:    "http://example.com/thumbnail.png",
					ImageURL:        "http://example.com/image.png",
					OpeningVideoURL: "http://example.jp/opening-video.mp4",
					Approved:        false,
					ApprovedAdminID: "",
					StartAt:         jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:           jst.Date(2022, 9, 1, 0, 0, 0, 0),
				},
			},
			expect: Schedules{
				{
					ID:              "schedule-id",
					CoordinatorID:   "coordinator-id",
					Status:          ScheduleStatusInProgress,
					Title:           "スケジュールタイトル",
					Description:     "スケジュールの詳細です。",
					ThumbnailURL:    "http://example.com/thumbnail.png",
					ImageURL:        "http://example.com/image.png",
					OpeningVideoURL: "http://example.jp/opening-video.mp4",
					Approved:        false,
					ApprovedAdminID: "",
					StartAt:         jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:           jst.Date(2022, 9, 1, 0, 0, 0, 0),
				},
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.schedules.Fill(now)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.schedules)
		})
	}
}

func TestSchedules_IDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		schedules Schedules
		expect    []string
	}{
		{
			name: "success",
			schedules: Schedules{
				{
					ID:              "schedule-id",
					CoordinatorID:   "coordinator-id",
					Title:           "スケジュールタイトル",
					Description:     "スケジュールの詳細です。",
					ThumbnailURL:    "https://and-period.jp/thumbnail.png",
					ImageURL:        "https://and-period.jp/image.png",
					OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
					Approved:        false,
					ApprovedAdminID: "",
					StartAt:         jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:           jst.Date(2022, 9, 1, 0, 0, 0, 0),
				},
			},
			expect: []string{"schedule-id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.schedules.IDs())
		})
	}
}

func TestSchedules_CoordinatorIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		schedules Schedules
		expect    []string
	}{
		{
			name: "success",
			schedules: Schedules{
				{
					ID:              "schedule-id",
					CoordinatorID:   "coordinator-id",
					Title:           "スケジュールタイトル",
					Description:     "スケジュールの詳細です。",
					ThumbnailURL:    "https://and-period.jp/thumbnail.png",
					ImageURL:        "https://and-period.jp/image.png",
					OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
					Approved:        false,
					ApprovedAdminID: "",
					StartAt:         jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:           jst.Date(2022, 9, 1, 0, 0, 0, 0),
				},
			},
			expect: []string{"coordinator-id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.schedules.CoordinatorIDs())
		})
	}
}
