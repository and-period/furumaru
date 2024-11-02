package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/stretchr/testify/assert"
)

func TestVideoViewerLogInterval(t *testing.T) {
	t.Parallel()
	type want struct {
		interval    VideoViewerLogInterval
		duration    time.Duration
		mediaEntity entity.AggregateVideoViewerLogInterval
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "second",
			request: "second",
			want: want{
				interval:    VideoViewerLogIntervalSecond,
				duration:    time.Second,
				mediaEntity: entity.AggregateVideoViewerLogIntervalSecond,
			},
		},
		{
			name:    "minute",
			request: "minute",
			want: want{
				interval:    VideoViewerLogIntervalMinute,
				duration:    time.Minute,
				mediaEntity: entity.AggregateVideoViewerLogIntervalMinute,
			},
		},
		{
			name:    "hour",
			request: "hour",
			want: want{
				interval:    VideoViewerLogIntervalHour,
				duration:    time.Hour,
				mediaEntity: entity.AggregateVideoViewerLogIntervalHour,
			},
		},
		{
			name:    "unknown",
			request: "unknown",
			want: want{
				interval:    "unknown",
				mediaEntity: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewVideoViewerLogIntervalFromRequest(tt.request)
			assert.Equal(t, tt.want.interval, actual)
			assert.Equal(t, tt.want.duration, actual.Duration())
			assert.Equal(t, tt.want.mediaEntity, actual.MediaEntity())
		})
	}
}

func TestVideoViewerLog(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		aggregate *entity.AggregatedVideoViewerLog
		interval  time.Duration
		expect    *VideoViewerLog
	}{
		{
			name: "success",
			aggregate: &entity.AggregatedVideoViewerLog{
				VideoID:    "video-id",
				ReportedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Total:      1,
			},
			interval: time.Minute,
			expect: &VideoViewerLog{
				VideoViewerLog: response.VideoViewerLog{
					StartAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
					EndAt:   time.Date(2021, 1, 1, 0, 1, 0, 0, time.UTC).Unix(),
					VideoID: "video-id",
					Total:   1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewVideoViewerLog(tt.aggregate, tt.interval)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideoViewerLog_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		log    *VideoViewerLog
		expect *response.VideoViewerLog
	}{
		{
			name: "success",
			log: &VideoViewerLog{
				VideoViewerLog: response.VideoViewerLog{
					StartAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
					EndAt:   time.Date(2021, 1, 1, 0, 1, 0, 0, time.UTC).Unix(),
					VideoID: "video-id",
					Total:   1,
				},
			},
			expect: &response.VideoViewerLog{
				StartAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				EndAt:   time.Date(2021, 1, 1, 0, 1, 0, 0, time.UTC).Unix(),
				VideoID: "video-id",
				Total:   1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.log.Response())
		})
	}
}

func TestVideoViewerLogs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		startAt    time.Time
		endAt      time.Time
		interval   VideoViewerLogInterval
		aggregates entity.AggregatedVideoViewerLogs
		expect     VideoViewerLogs
	}{
		{
			name:     "success",
			startAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			endAt:    time.Date(2021, 1, 1, 0, 2, 0, 0, time.UTC),
			interval: VideoViewerLogIntervalMinute,
			aggregates: entity.AggregatedVideoViewerLogs{
				{
					VideoID:    "video-id",
					ReportedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Total:      1,
				},
			},
			expect: VideoViewerLogs{
				{
					VideoViewerLog: response.VideoViewerLog{
						StartAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
						EndAt:   time.Date(2021, 1, 1, 0, 1, 0, 0, time.UTC).Unix(),
						VideoID: "video-id",
						Total:   1,
					},
				},
				{
					VideoViewerLog: response.VideoViewerLog{
						StartAt: time.Date(2021, 1, 1, 0, 1, 0, 0, time.UTC).Unix(),
						EndAt:   time.Date(2021, 1, 1, 0, 2, 0, 0, time.UTC).Unix(),
						VideoID: "video-id",
						Total:   0,
					},
				},
			},
		},
		{
			name:       "duration is zero",
			startAt:    time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			endAt:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			interval:   VideoViewerLogInterval(""),
			aggregates: nil,
			expect:     VideoViewerLogs{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewVideoViewerLogs(tt.interval, tt.startAt, tt.endAt, tt.aggregates)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideoViewerLogs_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		logs   VideoViewerLogs
		expect []*response.VideoViewerLog
	}{
		{
			name: "success",
			logs: VideoViewerLogs{
				{
					VideoViewerLog: response.VideoViewerLog{
						StartAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
						EndAt:   time.Date(2021, 1, 1, 0, 1, 0, 0, time.UTC).Unix(),
						VideoID: "video-id",
						Total:   1,
					},
				},
			},
			expect: []*response.VideoViewerLog{
				{
					StartAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
					EndAt:   time.Date(2021, 1, 1, 0, 1, 0, 0, time.UTC).Unix(),
					VideoID: "video-id",
					Total:   1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.logs.Response())
		})
	}
}
