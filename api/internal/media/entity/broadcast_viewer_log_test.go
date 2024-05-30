package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBroadcastViewerLog(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *BroadcastViewerLogParams
		expect *BroadcastViewerLog
	}{
		{
			name: "success",
			params: &BroadcastViewerLogParams{
				BroadcastID: "broadcast-id",
				SessionID:   "session-id",
				UserID:      "user-id",
				UserAgent:   "user-agent",
				ClientIP:    "127.0.0.1",
			},
			expect: &BroadcastViewerLog{
				BroadcastID: "broadcast-id",
				SessionID:   "session-id",
				UserID:      "user-id",
				UserAgent:   "user-agent",
				ClientIP:    "127.0.0.1",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewBroadcastViewerLog(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAggregatedBroadcastViewerLogs_MapByReportedAt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		logs   AggregatedBroadcastViewerLogs
		expect map[time.Time]*AggregatedBroadcastViewerLog
	}{
		{
			name: "success",
			logs: AggregatedBroadcastViewerLogs{
				{
					BroadcastID: "broadcast-id",
					ReportedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Total:       1,
				},
			},
			expect: map[time.Time]*AggregatedBroadcastViewerLog{
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC): {
					BroadcastID: "broadcast-id",
					ReportedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Total:       1,
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

func TestAggregatedBroadcastViewerLogs_GroupByBroadcastID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		logs   AggregatedBroadcastViewerLogs
		expect map[string]AggregatedBroadcastViewerLogs
	}{
		{
			name: "success",
			logs: AggregatedBroadcastViewerLogs{
				{
					BroadcastID: "broadcast-id",
					ReportedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Total:       1,
				},
			},
			expect: map[string]AggregatedBroadcastViewerLogs{
				"broadcast-id": {
					{
						BroadcastID: "broadcast-id",
						ReportedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						Total:       1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.logs.GroupByBroadcastID()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
