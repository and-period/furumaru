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
				MessageType: ScheduleTypeNotification,
				MessageID:   "message-id",
				SentAt:      jst.Date(2022, 7, 18, 18, 30, 0, 0),
			},
			expect: &Schedule{
				MessageType: ScheduleTypeNotification,
				MessageID:   "message-id",
				Status:      ScheduleStatusWaiting,
				Count:       0,
				SentAt:      jst.Date(2022, 7, 18, 18, 30, 0, 0),
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewSchedule(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestSchedule_Executable(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		schedule *Schedule
		expect   bool
	}{
		{
			name: "deadline",
			schedule: &Schedule{
				Status:   ScheduleStatusWaiting,
				Count:    0,
				SentAt:   now.Add(-1 * time.Minute),
				Deadline: now.Add(-1 * time.Minute),
			},
			expect: false,
		},
		{
			name: "waiting",
			schedule: &Schedule{
				Status: ScheduleStatusWaiting,
				Count:  0,
				SentAt: now,
			},
			expect: true,
		},
		{
			name: "waiting to before sent at",
			schedule: &Schedule{
				Status: ScheduleStatusWaiting,
				Count:  0,
				SentAt: now.Add(time.Hour),
			},
			expect: false,
		},
		{
			name: "processing to 1 minute after last run",
			schedule: &Schedule{
				Status:    ScheduleStatusProcessing,
				Count:     1,
				UpdatedAt: now.Add(-1 * time.Minute),
			},
			expect: false,
		},
		{
			name: "processing to 15 minute after last run",
			schedule: &Schedule{
				Status:    ScheduleStatusProcessing,
				Count:     1,
				UpdatedAt: now.Add(-15 * time.Minute),
			},
			expect: true,
		},
		{
			name: "processing to re-executed",
			schedule: &Schedule{
				Status:    ScheduleStatusProcessing,
				Count:     2,
				UpdatedAt: now.Add(-15 * time.Minute),
			},
			expect: false,
		},
		{
			name: "done",
			schedule: &Schedule{
				Status:    ScheduleStatusDone,
				Count:     1,
				UpdatedAt: now.Add(-15 * time.Minute),
			},
			expect: false,
		},
		{
			name: "canceled",
			schedule: &Schedule{
				Status:    ScheduleStatusCanceled,
				Count:     2,
				UpdatedAt: now.Add(-15 * time.Minute),
			},
			expect: false,
		},
		{
			name: "unknown",
			schedule: &Schedule{
				Status: -1,
			},
			expect: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.schedule.Executable(now))
		})
	}
}

func TestSchedule_ShouldCancel(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name string

		schedule *Schedule
		expect   bool
	}{
		{
			name: "deadline",
			schedule: &Schedule{
				Status:   ScheduleStatusWaiting,
				Count:    0,
				SentAt:   now.Add(-1 * time.Minute),
				Deadline: now.Add(-1 * time.Minute),
			},
			expect: true,
		},
		{
			name: "waiting",
			schedule: &Schedule{
				Status: ScheduleStatusWaiting,
				Count:  0,
			},
			expect: false,
		},
		{
			name: "processing to 1 minute after last run",
			schedule: &Schedule{
				Status:    ScheduleStatusProcessing,
				Count:     1,
				UpdatedAt: now.Add(-1 * time.Minute),
			},
			expect: false,
		},
		{
			name: "processing to 15 minute after last run",
			schedule: &Schedule{
				Status:    ScheduleStatusProcessing,
				Count:     1,
				UpdatedAt: now.Add(-15 * time.Minute),
			},
			expect: false,
		},
		{
			name: "processing to re-executed",
			schedule: &Schedule{
				Status:    ScheduleStatusProcessing,
				Count:     2,
				UpdatedAt: now.Add(-15 * time.Minute),
			},
			expect: true,
		},
		{
			name: "done",
			schedule: &Schedule{
				Status:    ScheduleStatusDone,
				Count:     1,
				UpdatedAt: now.Add(-15 * time.Minute),
			},
			expect: false,
		},
		{
			name: "canceled",
			schedule: &Schedule{
				Status:    ScheduleStatusCanceled,
				Count:     2,
				UpdatedAt: now.Add(-15 * time.Minute),
			},
			expect: false,
		},
		{
			name: "unknown",
			schedule: &Schedule{
				Status: -1,
			},
			expect: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.schedule.ShouldCancel(now))
		})
	}
}

func TestSchedules_Map(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name      string
		schedules Schedules
		expect    map[string]*Schedule
	}{
		{
			name: "success",
			schedules: Schedules{
				{
					MessageType: ScheduleTypeNotification,
					MessageID:   "message-id",
					Status:      ScheduleStatusDone,
					Count:       1,
					SentAt:      now,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			expect: map[string]*Schedule{
				"message-id": {
					MessageType: ScheduleTypeNotification,
					MessageID:   "message-id",
					Status:      ScheduleStatusDone,
					Count:       1,
					SentAt:      now,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.schedules.Map())
		})
	}
}
