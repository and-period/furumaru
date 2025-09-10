package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/stretchr/testify/assert"
)

func TestBroadcastViewerLogInterval(t *testing.T) {
	t.Parallel()
	type want struct {
		interval    BroadcastViewerLogInterval
		duration    time.Duration
		mediaEntity entity.AggregateBroadcastViewerLogInterval
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
				interval:    BroadcastViewerLogInterval(types.BroadcastViewerLogIntervalSecond),
				duration:    time.Second,
				mediaEntity: entity.AggregateBroadcastViewerLogInterval(types.BroadcastViewerLogIntervalSecond),
			},
		},
		{
			name:    "minute",
			request: "minute",
			want: want{
				interval:    BroadcastViewerLogInterval(types.BroadcastViewerLogIntervalMinute),
				duration:    time.Minute,
				mediaEntity: entity.AggregateBroadcastViewerLogInterval(types.BroadcastViewerLogIntervalMinute),
			},
		},
		{
			name:    "hour",
			request: "hour",
			want: want{
				interval:    BroadcastViewerLogInterval(types.BroadcastViewerLogIntervalHour),
				duration:    time.Hour,
				mediaEntity: entity.AggregateBroadcastViewerLogInterval(types.BroadcastViewerLogIntervalHour),
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
			actual := NewBroadcastViewerLogInterval(types.BroadcastViewerLogIntervalFromRequest)(tt.request)
			assert.Equal(t, tt.want.interval, actual)
			assert.Equal(t, tt.want.duration, actual.Duration())
			assert.Equal(t, tt.want.mediaEntity, actual.MediaEntity())
		})
	}
}

func TestBroadcastViewerLog(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		aggregate *entity.AggregatedBroadcastViewerLog
		interval  time.Duration
		expect    *BroadcastViewerLog
	}{
		{
			name: "success",
			aggregate: &entity.AggregatedBroadcastViewerLog{
				BroadcastID: "broadcast-id",
				ReportedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Total:       1,
			},
			interval: time.Minute,
			expect: &BroadcastViewerLog{
				BroadcastViewerLog: types.BroadcastViewerLog{
					StartAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
					EndAt:       time.Date(2021, 1, 1, 0, 1, 0, 0, time.UTC).Unix(),
					BroadcastID: "broadcast-id",
					Total:       1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewBroadcastViewerLog(tt.aggregate, tt.interval)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestBroadcastViewerLog_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		log    *BroadcastViewerLog
		expect *types.BroadcastViewerLog
	}{
		{
			name: "success",
			log: &BroadcastViewerLog{
				BroadcastViewerLog: types.BroadcastViewerLog{
					StartAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
					EndAt:       time.Date(2021, 1, 1, 0, 1, 0, 0, time.UTC).Unix(),
					BroadcastID: "broadcast-id",
					Total:       1,
				},
			},
			expect: &types.BroadcastViewerLog{
				StartAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				EndAt:       time.Date(2021, 1, 1, 0, 1, 0, 0, time.UTC).Unix(),
				BroadcastID: "broadcast-id",
				Total:       1,
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

func TestBroadcastViewerLogs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		startAt    time.Time
		endAt      time.Time
		interval   BroadcastViewerLogInterval
		aggregates entity.AggregatedBroadcastViewerLogs
		expect     BroadcastViewerLogs
	}{
		{
			name:     "success",
			startAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			endAt:    time.Date(2021, 1, 1, 0, 2, 0, 0, time.UTC),
			interval: BroadcastViewerLogInterval(types.BroadcastViewerLogIntervalMinute),
			aggregates: entity.AggregatedBroadcastViewerLogs{
				{
					BroadcastID: "broadcast-id",
					ReportedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Total:       1,
				},
			},
			expect: BroadcastViewerLogs{
				{
					BroadcastViewerLog: types.BroadcastViewerLog{
						StartAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
						EndAt:       time.Date(2021, 1, 1, 0, 1, 0, 0, time.UTC).Unix(),
						BroadcastID: "broadcast-id",
						Total:       1,
					},
				},
				{
					BroadcastViewerLog: types.BroadcastViewerLog{
						StartAt:     time.Date(2021, 1, 1, 0, 1, 0, 0, time.UTC).Unix(),
						EndAt:       time.Date(2021, 1, 1, 0, 2, 0, 0, time.UTC).Unix(),
						BroadcastID: "broadcast-id",
						Total:       0,
					},
				},
			},
		},
		{
			name:       "duration is zero",
			startAt:    time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			endAt:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			interval:   BroadcastViewerLogInterval(""),
			aggregates: nil,
			expect:     BroadcastViewerLogs{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewBroadcastViewerLogs(tt.interval, tt.startAt, tt.endAt, tt.aggregates)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestBroadcastViewerLogs_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		logs   BroadcastViewerLogs
		expect []*types.BroadcastViewerLog
	}{
		{
			name: "success",
			logs: BroadcastViewerLogs{
				{
					BroadcastViewerLog: types.BroadcastViewerLog{
						StartAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
						EndAt:       time.Date(2021, 1, 1, 0, 1, 0, 0, time.UTC).Unix(),
						BroadcastID: "broadcast-id",
						Total:       1,
					},
				},
			},
			expect: []*types.BroadcastViewerLog{
				{
					StartAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
					EndAt:       time.Date(2021, 1, 1, 0, 1, 0, 0, time.UTC).Unix(),
					BroadcastID: "broadcast-id",
					Total:       1,
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
