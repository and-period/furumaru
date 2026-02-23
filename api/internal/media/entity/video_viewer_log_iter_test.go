package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAggregatedVideoViewerLogs_All(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name string
		logs AggregatedVideoViewerLogs
	}{
		{
			name: "success",
			logs: AggregatedVideoViewerLogs{
				{VideoID: "video-id01", ReportedAt: now, Total: 10},
				{VideoID: "video-id02", ReportedAt: now, Total: 20},
			},
		},
		{
			name: "empty",
			logs: AggregatedVideoViewerLogs{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var videoIDs []string
			for i, l := range tt.logs.All() {
				indices = append(indices, i)
				videoIDs = append(videoIDs, l.VideoID)
			}
			for i, l := range tt.logs {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, l.VideoID, videoIDs[i])
				}
			}
			assert.Len(t, indices, len(tt.logs))
		})
	}
}

func TestAggregatedVideoViewerLogs_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	now := time.Now()
	logs := AggregatedVideoViewerLogs{
		{VideoID: "video-id01", ReportedAt: now, Total: 10},
		{VideoID: "video-id02", ReportedAt: now, Total: 20},
		{VideoID: "video-id03", ReportedAt: now, Total: 30},
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

func TestAggregatedVideoViewerLogs_IterMapByReportedAt(t *testing.T) {
	t.Parallel()
	now := time.Now().Truncate(time.Second)
	later := now.Add(time.Minute)
	tests := []struct {
		name   string
		logs   AggregatedVideoViewerLogs
		expect map[time.Time]string // reportedAt -> videoID
	}{
		{
			name: "success",
			logs: AggregatedVideoViewerLogs{
				{VideoID: "video-id01", ReportedAt: now, Total: 10},
				{VideoID: "video-id02", ReportedAt: later, Total: 20},
			},
			expect: map[time.Time]string{
				now:   "video-id01",
				later: "video-id02",
			},
		},
		{
			name:   "empty",
			logs:   AggregatedVideoViewerLogs{},
			expect: map[time.Time]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[time.Time]string)
			for k, v := range tt.logs.IterMapByReportedAt() {
				result[k] = v.VideoID
			}
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestAggregatedVideoViewerLogs_IterGroupByVideoID(t *testing.T) {
	t.Parallel()
	now := time.Now().Truncate(time.Second)
	later := now.Add(time.Minute)
	tests := []struct {
		name   string
		logs   AggregatedVideoViewerLogs
		expect map[string]int // videoID -> count
	}{
		{
			name: "success",
			logs: AggregatedVideoViewerLogs{
				{VideoID: "video-id01", ReportedAt: now, Total: 10},
				{VideoID: "video-id01", ReportedAt: later, Total: 20},
				{VideoID: "video-id02", ReportedAt: now, Total: 30},
			},
			expect: map[string]int{
				"video-id01": 2,
				"video-id02": 1,
			},
		},
		{
			name:   "empty",
			logs:   AggregatedVideoViewerLogs{},
			expect: map[string]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]int)
			for k, v := range tt.logs.IterGroupByVideoID() {
				result[k] = len(v)
			}
			assert.Equal(t, tt.expect, result)
		})
	}
}
