package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestScheduleStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.ScheduleStatus
		expect ScheduleStatus
	}{
		{
			name:   "private",
			status: entity.ScheduleStatusPrivate,
			expect: ScheduleStatusPrivate,
		},
		{
			name:   "in progress",
			status: entity.ScheduleStatusInProgress,
			expect: ScheduleStatusInProgress,
		},
		{
			name:   "waiting",
			status: entity.ScheduleStatusWaiting,
			expect: ScheduleStatusWaiting,
		},
		{
			name:   "live",
			status: entity.ScheduleStatusLive,
			expect: ScheduleStatusLive,
		},
		{
			name:   "closed",
			status: entity.ScheduleStatusClosed,
			expect: ScheduleStatusClosed,
		},
		{
			name:   "unknown",
			status: entity.ScheduleStatusUnknown,
			expect: ScheduleStatusUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewScheduleStatus(tt.status))
		})
	}
}

func TestSchedule(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	tests := []struct {
		name     string
		schedule *entity.Schedule
		expect   *Schedule
	}{
		{
			name: "success",
			schedule: &entity.Schedule{
				ID:              "schedule-id",
				CoordinatorID:   "coordinator-id",
				Status:          entity.ScheduleStatusLive,
				Title:           "スケジュールタイトル",
				Description:     "スケジュールの詳細です。",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
				Public:          true,
				Approved:        true,
				ApprovedAdminID: "admin-id",
				StartAt:         now.AddDate(0, -1, 0),
				EndAt:           now.AddDate(0, 1, 0),
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			expect: &Schedule{
				Schedule: response.Schedule{
					ID:              "schedule-id",
					CoordinatorID:   "coordinator-id",
					Status:          ScheduleStatusLive.Response(),
					Title:           "スケジュールタイトル",
					Description:     "スケジュールの詳細です。",
					ThumbnailURL:    "https://and-period.jp/thumbnail.png",
					ImageURL:        "https://and-period.jp/image.png",
					OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
					Public:          true,
					Approved:        true,
					StartAt:         1638284400,
					EndAt:           1643641200,
					CreatedAt:       1640962800,
					UpdatedAt:       1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewSchedule(tt.schedule))
		})
	}
}

func TestSchedule_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		schedule *Schedule
		expect   *response.Schedule
	}{
		{
			name: "success",
			schedule: &Schedule{
				Schedule: response.Schedule{
					ID:              "schedule-id",
					CoordinatorID:   "coordinator-id",
					Status:          ScheduleStatusLive.Response(),
					Title:           "スケジュールタイトル",
					Description:     "スケジュールの詳細です。",
					ThumbnailURL:    "https://and-period.jp/thumbnail.png",
					ImageURL:        "https://and-period.jp/image.png",
					OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
					Public:          true,
					Approved:        true,
					StartAt:         1638284400,
					EndAt:           1643641200,
					CreatedAt:       1640962800,
					UpdatedAt:       1640962800,
				},
			},
			expect: &response.Schedule{
				ID:              "schedule-id",
				CoordinatorID:   "coordinator-id",
				Status:          ScheduleStatusLive.Response(),
				Title:           "スケジュールタイトル",
				Description:     "スケジュールの詳細です。",
				ThumbnailURL:    "https://and-period.jp/thumbnail.png",
				ImageURL:        "https://and-period.jp/image.png",
				OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
				Public:          true,
				Approved:        true,
				StartAt:         1638284400,
				EndAt:           1643641200,
				CreatedAt:       1640962800,
				UpdatedAt:       1640962800,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.schedule.Response())
		})
	}
}

func TestSchedules(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	tests := []struct {
		name      string
		schedules entity.Schedules
		expect    Schedules
	}{
		{
			name: "success",
			schedules: entity.Schedules{
				{
					ID:              "schedule-id",
					CoordinatorID:   "coordinator-id",
					Status:          entity.ScheduleStatusLive,
					Title:           "スケジュールタイトル",
					Description:     "スケジュールの詳細です。",
					ThumbnailURL:    "https://and-period.jp/thumbnail.png",
					ImageURL:        "https://and-period.jp/image.png",
					OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
					Public:          true,
					Approved:        true,
					ApprovedAdminID: "admin-id",
					StartAt:         now.AddDate(0, -1, 0),
					EndAt:           now.AddDate(0, 1, 0),
					CreatedAt:       now,
					UpdatedAt:       now,
				},
			},
			expect: Schedules{
				{
					Schedule: response.Schedule{
						ID:              "schedule-id",
						CoordinatorID:   "coordinator-id",
						Status:          ScheduleStatusLive.Response(),
						Title:           "スケジュールタイトル",
						Description:     "スケジュールの詳細です。",
						ThumbnailURL:    "https://and-period.jp/thumbnail.png",
						ImageURL:        "https://and-period.jp/image.png",
						OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
						Public:          true,
						Approved:        true,
						StartAt:         1638284400,
						EndAt:           1643641200,
						CreatedAt:       1640962800,
						UpdatedAt:       1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewSchedules(tt.schedules))
		})
	}
}

func TestSchedules_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		schedules Schedules
		expect    []*response.Schedule
	}{
		{
			name: "success",
			schedules: Schedules{
				{
					Schedule: response.Schedule{
						ID:              "schedule-id",
						CoordinatorID:   "coordinator-id",
						Status:          ScheduleStatusLive.Response(),
						Title:           "スケジュールタイトル",
						Description:     "スケジュールの詳細です。",
						ThumbnailURL:    "https://and-period.jp/thumbnail.png",
						ImageURL:        "https://and-period.jp/image.png",
						OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
						Public:          true,
						Approved:        true,
						StartAt:         1638284400,
						EndAt:           1643641200,
						CreatedAt:       1640962800,
						UpdatedAt:       1640962800,
					},
				},
			},
			expect: []*response.Schedule{
				{
					ID:              "schedule-id",
					CoordinatorID:   "coordinator-id",
					Status:          ScheduleStatusLive.Response(),
					Title:           "スケジュールタイトル",
					Description:     "スケジュールの詳細です。",
					ThumbnailURL:    "https://and-period.jp/thumbnail.png",
					ImageURL:        "https://and-period.jp/image.png",
					OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
					Public:          true,
					Approved:        true,
					StartAt:         1638284400,
					EndAt:           1643641200,
					CreatedAt:       1640962800,
					UpdatedAt:       1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.schedules.Response())
		})
	}
}
