package service

import (
	"testing"
	"time"

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
					ID:               "broadcast-id",
					ScheduleID:       "schedule-id",
					Status:           int32(BroadcastStatusIdle),
					InputURL:         "rtmp://127.0.0.1:1935/app/instance",
					OutputURL:        "http://example.com/index.m3u8",
					ArchiveURL:       "http://example.com/index.mp4",
					YoutubeAccount:   "youtube-account",
					YoutubeViewerURL: "https://youtube.com/live/youtube-broadcast-id",
					YoutubeAdminURL:  "https://studio.youtube.com/video/youtube-broadcast-id/livestreaming",
					CreatedAt:        1640962800,
					UpdatedAt:        1640962800,
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

func TestGuestBroadcast(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		schedule    *Schedule
		shop        *Shop
		coordinator *Coordinator
		expect      *GuestBroadcast
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
			shop: &Shop{
				Shop: response.Shop{
					ID:             "shop-id",
					Name:           "&.マルシェ",
					CoordinatorID:  "coordinator-id",
					ProducerIDs:    []string{"producer-id"},
					ProductTypeIDs: []string{"product-type-id"},
					BusinessDays:   []time.Weekday{time.Monday},
					CreatedAt:      1640962800,
					UpdatedAt:      1640962800,
				},
			},
			coordinator: &Coordinator{
				Coordinator: response.Coordinator{
					ID:                "coordinator-id",
					Status:            int32(AdminStatusActivated),
					Lastname:          "&.",
					Firstname:         "管理者",
					LastnameKana:      "あんどどっと",
					FirstnameKana:     "かんりしゃ",
					Username:          "&.農園",
					Profile:           "紹介文です。",
					ThumbnailURL:      "https://and-period.jp/thumbnail.png",
					HeaderURL:         "https://and-period.jp/header.png",
					PromotionVideoURL: "https://and-period.jp/promotion.mp4",
					BonusVideoURL:     "https://and-period.jp/bonus.mp4",
					InstagramID:       "instagram-id",
					FacebookID:        "facebook-id",
					Email:             "test-coordinator@and-period.jp",
					PhoneNumber:       "+819012345678",
					PostalCode:        "1000014",
					PrefectureCode:    13,
					City:              "千代田区",
					CreatedAt:         1640962800,
					UpdatedAt:         1640962800,
				},
			},
			expect: &GuestBroadcast{
				GuestBroadcast: response.GuestBroadcast{
					Title:             "スケジュールタイトル",
					Description:       "スケジュールの詳細です。",
					StartAt:           1638284400,
					EndAt:             1643641200,
					CoordinatorMarche: "&.マルシェ",
					CoordinatorName:   "&.農園",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewGuestBroadcast(tt.schedule, tt.shop, tt.coordinator)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
