package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVideoViewerLog(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewVideoViewerLogParams
		expect *VideoViewerLog
	}{
		{
			name: "success",
			params: &NewVideoViewerLogParams{
				VideoID:   "video-id",
				SessionID: "session-id",
				UserID:    "user-id",
				UserAgent: "user-agent",
				ClientIP:  "127.0.0.1",
			},
			expect: &VideoViewerLog{
				VideoID:   "video-id",
				SessionID: "session-id",
				UserID:    "user-id",
				UserAgent: "user-agent",
				ClientIP:  "127.0.0.1",
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewVideoViewerLog(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAggregatedVideoViewerLogs_MapByReportedAt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		logs   AggregatedVideoViewerLogs
		expect map[time.Time]*AggregatedVideoViewerLog
	}{
		{
			name: "success",
			logs: AggregatedVideoViewerLogs{
				{
					VideoID:    "video-id",
					ReportedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Total:      1,
				},
			},
			expect: map[time.Time]*AggregatedVideoViewerLog{
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC): {
					VideoID:    "video-id",
					ReportedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Total:      1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.logs.MapByReportedAt()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAggregatedVideoViewerLogs_GroupByVideoID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		logs   AggregatedVideoViewerLogs
		expect map[string]AggregatedVideoViewerLogs
	}{
		{
			name: "success",
			logs: AggregatedVideoViewerLogs{
				{
					VideoID:    "video-id",
					ReportedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Total:      1,
				},
			},
			expect: map[string]AggregatedVideoViewerLogs{
				"video-id": {
					{
						VideoID:    "video-id",
						ReportedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						Total:      1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.logs.GroupByVideoID()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
