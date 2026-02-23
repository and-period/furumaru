package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAggregatedBroadcastViewerLogs_All(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name string
		logs AggregatedBroadcastViewerLogs
	}{
		{
			name: "success",
			logs: AggregatedBroadcastViewerLogs{
				{BroadcastID: "broadcast-id01", ReportedAt: now, Total: 10},
				{BroadcastID: "broadcast-id02", ReportedAt: now, Total: 20},
			},
		},
		{
			name: "empty",
			logs: AggregatedBroadcastViewerLogs{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var broadcastIDs []string
			for i, l := range tt.logs.All() {
				indices = append(indices, i)
				broadcastIDs = append(broadcastIDs, l.BroadcastID)
			}
			for i, l := range tt.logs {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, l.BroadcastID, broadcastIDs[i])
				}
			}
			assert.Len(t, indices, len(tt.logs))
		})
	}
}

func TestAggregatedBroadcastViewerLogs_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	now := time.Now()
	logs := AggregatedBroadcastViewerLogs{
		{BroadcastID: "broadcast-id01", ReportedAt: now, Total: 10},
		{BroadcastID: "broadcast-id02", ReportedAt: now, Total: 20},
		{BroadcastID: "broadcast-id03", ReportedAt: now, Total: 30},
	}
	var count int
	for range logs.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestAggregatedBroadcastViewerLogs_IterMapByReportedAt(t *testing.T) {
	t.Parallel()
	now := time.Now().Truncate(time.Second)
	later := now.Add(time.Minute)
	tests := []struct {
		name   string
		logs   AggregatedBroadcastViewerLogs
		expect map[time.Time]string // reportedAt -> broadcastID
	}{
		{
			name: "success",
			logs: AggregatedBroadcastViewerLogs{
				{BroadcastID: "broadcast-id01", ReportedAt: now, Total: 10},
				{BroadcastID: "broadcast-id02", ReportedAt: later, Total: 20},
			},
			expect: map[time.Time]string{
				now:   "broadcast-id01",
				later: "broadcast-id02",
			},
		},
		{
			name:   "empty",
			logs:   AggregatedBroadcastViewerLogs{},
			expect: map[time.Time]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[time.Time]string)
			for k, v := range tt.logs.IterMapByReportedAt() {
				result[k] = v.BroadcastID
			}
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestAggregatedBroadcastViewerLogs_IterGroupByBroadcastID(t *testing.T) {
	t.Parallel()
	now := time.Now().Truncate(time.Second)
	later := now.Add(time.Minute)
	tests := []struct {
		name   string
		logs   AggregatedBroadcastViewerLogs
		expect map[string]int // broadcastID -> count
	}{
		{
			name: "success",
			logs: AggregatedBroadcastViewerLogs{
				{BroadcastID: "broadcast-id01", ReportedAt: now, Total: 10},
				{BroadcastID: "broadcast-id01", ReportedAt: later, Total: 20},
				{BroadcastID: "broadcast-id02", ReportedAt: now, Total: 30},
			},
			expect: map[string]int{
				"broadcast-id01": 2,
				"broadcast-id02": 1,
			},
		},
		{
			name:   "empty",
			logs:   AggregatedBroadcastViewerLogs{},
			expect: map[string]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]int)
			for k, v := range tt.logs.IterGroupByBroadcastID() {
				result[k] = len(v)
			}
			assert.Equal(t, tt.expect, result)
		})
	}
}
