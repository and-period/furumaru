package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestBroadcastStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.BroadcastStatus
		expect BroadcastStatus
	}{
		{
			name:   "disabled",
			status: entity.BroadcastStatusDisabled,
			expect: BroadcastStatusDisabled,
		},
		{
			name:   "waiting",
			status: entity.BroadcastStatusWaiting,
			expect: BroadcastStatusWaiting,
		},
		{
			name:   "idle",
			status: entity.BroadcastStatusIdle,
			expect: BroadcastStatusIdle,
		},
		{
			name:   "active",
			status: entity.BroadcastStatusActive,
			expect: BroadcastStatusActive,
		},
		{
			name:   "unknown",
			status: entity.BroadcastStatusUnknown,
			expect: BroadcastStatusUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewBroadcastStatus(tt.status))
		})
	}
}

func TestBroadcast(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		broadcast *entity.Broadcast
		expect    *Broadcast
	}{
		{
			name: "success",
			broadcast: &entity.Broadcast{
				ID:                 "broadcast-id",
				ScheduleID:         "schedule-id",
				Status:             entity.BroadcastStatusIdle,
				InputURL:           "rtmp://127.0.0.1:1935/app/instance",
				OutputURL:          "http://example.com/index.m3u8",
				ArchiveURL:         "http://example.com/index.mp4",
				YoutubeAccount:     "youtube-account",
				YoutubeBroadcastID: "youtube-broadcast-id",
				YoutubeStreamID:    "youtube-stream-id",
				YoutubeStreamKey:   "youtube-stream-key",
				YoutubeStreamURL:   "rtmp://stream.example.com",
				YoutubeBackupURL:   "rtmp://backup.example.com",
				CreatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &Broadcast{
				Broadcast: response.Broadcast{
					ID:              "broadcast-id",
					ScheduleID:      "schedule-id",
					Status:          int32(BroadcastStatusIdle),
					InputURL:        "rtmp://127.0.0.1:1935/app/instance",
					OutputURL:       "http://example.com/index.m3u8",
					ArchiveURL:      "http://example.com/index.mp4",
					YouTubeAccount:  "youtube-account",
					YouTubeAdminURL: "https://studio.youtube.com/video/youtube-broadcast-id/livestreaming",
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
			assert.Equal(t, tt.expect, NewBroadcast(tt.broadcast))
		})
	}
}

func TestBroadcasts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		broadcasts entity.Broadcasts
		expect     Broadcasts
	}{
		{
			name: "success",
			broadcasts: entity.Broadcasts{
				{
					ID:         "broadcast-id",
					ScheduleID: "schedule-id",
					Status:     entity.BroadcastStatusIdle,
					InputURL:   "rtmp://127.0.0.1:1935/app/instance",
					OutputURL:  "http://example.com/index.m3u8",
					ArchiveURL: "http://example.com/index.mp4",
					CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: Broadcasts{
				{
					Broadcast: response.Broadcast{
						ID:         "broadcast-id",
						ScheduleID: "schedule-id",
						Status:     int32(BroadcastStatusIdle),
						InputURL:   "rtmp://127.0.0.1:1935/app/instance",
						OutputURL:  "http://example.com/index.m3u8",
						ArchiveURL: "http://example.com/index.mp4",
						CreatedAt:  1640962800,
						UpdatedAt:  1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewBroadcasts(tt.broadcasts))
		})
	}
}
