package entity

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestLive(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewLiveParams
		expect *Live
	}{
		{
			name: "success",
			params: &NewLiveParams{
				ScheduleID:  "schedule-id",
				ProducerID:  "producer-id",
				Title:       "ライブのタイトル",
				Description: "ライブの説明",
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
			expect: &Live{
				ScheduleID:  "schedule-id",
				ProducerID:  "producer-id",
				Title:       "ライブのタイトル",
				Description: "ライブの説明",
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewLive(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestLive_Fill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		live     *Live
		products LiveProducts
		now      time.Time
		expect   *Live
	}{
		{
			name: "success canceled",
			live: &Live{
				ID:          "live-id",
				ScheduleID:  "schedule-id",
				ProducerID:  "producer-id",
				Title:       "ライブのタイトル",
				Description: "ライブの説明",
				Published:   false,
				Canceled:    true,
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
				CreatedAt:   jst.Date(2022, 7, 1, 0, 0, 0, 0),
				UpdatedAt:   jst.Date(2022, 7, 1, 0, 0, 0, 0),
			},
			products: LiveProducts{
				{
					LiveID:    "live-id",
					ProductID: "product-id",
					CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				},
			},
			now: jst.Date(2022, 7, 1, 0, 0, 0, 0),
			expect: &Live{
				ID:          "live-id",
				ScheduleID:  "schedule-id",
				ProducerID:  "producer-id",
				Title:       "ライブのタイトル",
				Description: "ライブの説明",
				Status:      LiveStatusCanceled,
				Published:   false,
				Canceled:    true,
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
				LiveProducts: LiveProducts{
					{
						LiveID:    "live-id",
						ProductID: "product-id",
						CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
						UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					},
				},
				CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
			},
		},
		{
			name: "success non published",
			live: &Live{
				ID:          "live-id",
				ScheduleID:  "schedule-id",
				ProducerID:  "producer-id",
				Title:       "ライブのタイトル",
				Description: "ライブの説明",
				Published:   false,
				Canceled:    false,
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
				CreatedAt:   jst.Date(2022, 7, 1, 0, 0, 0, 0),
				UpdatedAt:   jst.Date(2022, 7, 1, 0, 0, 0, 0),
			},
			products: LiveProducts{
				{
					LiveID:    "live-id",
					ProductID: "product-id",
					CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				},
			},
			now: jst.Date(2022, 7, 1, 0, 0, 0, 0),
			expect: &Live{
				ID:          "live-id",
				ScheduleID:  "schedule-id",
				ProducerID:  "producer-id",
				Title:       "ライブのタイトル",
				Description: "ライブの説明",
				Status:      LiveStatusWaiting,
				Published:   false,
				Canceled:    false,
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
				LiveProducts: LiveProducts{
					{
						LiveID:    "live-id",
						ProductID: "product-id",
						CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
						UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					},
				},
				CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
			},
		},
		{
			name: "success waiting",
			live: &Live{
				ID:          "live-id",
				ScheduleID:  "schedule-id",
				ProducerID:  "producer-id",
				Title:       "ライブのタイトル",
				Description: "ライブの説明",
				Published:   true,
				Canceled:    false,
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
				CreatedAt:   jst.Date(2022, 7, 1, 0, 0, 0, 0),
				UpdatedAt:   jst.Date(2022, 7, 1, 0, 0, 0, 0),
			},
			products: LiveProducts{
				{
					LiveID:    "live-id",
					ProductID: "product-id",
					CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				},
			},
			now: jst.Date(2022, 7, 1, 0, 0, 0, 0),
			expect: &Live{
				ID:          "live-id",
				ScheduleID:  "schedule-id",
				ProducerID:  "producer-id",
				Title:       "ライブのタイトル",
				Description: "ライブの説明",
				Status:      LiveStatusWaiting,
				Published:   true,
				Canceled:    false,
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
				LiveProducts: LiveProducts{
					{
						LiveID:    "live-id",
						ProductID: "product-id",
						CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
						UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					},
				},
				CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
			},
		},
		{
			name: "success opened",
			live: &Live{
				ID:          "live-id",
				ScheduleID:  "schedule-id",
				ProducerID:  "producer-id",
				Title:       "ライブのタイトル",
				Description: "ライブの説明",
				Published:   true,
				Canceled:    false,
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
				CreatedAt:   jst.Date(2022, 7, 1, 0, 0, 0, 0),
				UpdatedAt:   jst.Date(2022, 7, 1, 0, 0, 0, 0),
			},
			products: LiveProducts{
				{
					LiveID:    "live-id",
					ProductID: "product-id",
					CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				},
			},
			now: jst.Date(2022, 8, 15, 0, 0, 0, 0),
			expect: &Live{
				ID:          "live-id",
				ScheduleID:  "schedule-id",
				ProducerID:  "producer-id",
				Title:       "ライブのタイトル",
				Description: "ライブの説明",
				Status:      LiveStatusOpened,
				Published:   true,
				Canceled:    false,
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
				LiveProducts: LiveProducts{
					{
						LiveID:    "live-id",
						ProductID: "product-id",
						CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
						UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					},
				},
				CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
			},
		},
		{
			name: "success closed",
			live: &Live{
				ID:          "live-id",
				ScheduleID:  "schedule-id",
				ProducerID:  "producer-id",
				Title:       "ライブのタイトル",
				Description: "ライブの説明",
				Published:   true,
				Canceled:    false,
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
				CreatedAt:   jst.Date(2022, 7, 1, 0, 0, 0, 0),
				UpdatedAt:   jst.Date(2022, 7, 1, 0, 0, 0, 0),
			},
			products: LiveProducts{
				{
					LiveID:    "live-id",
					ProductID: "product-id",
					CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				},
			},
			now: jst.Date(2022, 9, 15, 0, 0, 0, 0),
			expect: &Live{
				ID:          "live-id",
				ScheduleID:  "schedule-id",
				ProducerID:  "producer-id",
				Title:       "ライブのタイトル",
				Description: "ライブの説明",
				Status:      LiveStatusClosed,
				Published:   true,
				Canceled:    false,
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
				LiveProducts: LiveProducts{
					{
						LiveID:    "live-id",
						ProductID: "product-id",
						CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
						UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					},
				},
				CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.live.Fill(tt.products, tt.now)
			assert.Equal(t, tt.expect, tt.live)
		})
	}
}

func TestLive_FillIVS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		live   *Live
		params FillLiveIvsParams
		expect *Live
	}{
		{
			name: "success",
			live: &Live{
				ID:           "live-id",
				ChannelArn:   "channel-arn",
				StreamKeyArn: "streamKey-arn",
			},
			params: FillLiveIvsParams{
				ChannelName:    "配信チャンネル",
				IngestEndpoint: "ingest-endpoint",
				StreamKey:      "streamKey-value",
				StreamID:       "stream-id",
				PlaybackURL:    "playback-url",
				ViewerCount:    100,
			},
			expect: &Live{
				ID:             "live-id",
				ChannelArn:     "channel-arn",
				StreamKeyArn:   "streamKey-arn",
				ChannelName:    "配信チャンネル",
				IngestEndpoint: "ingest-endpoint",
				StreamKey:      "streamKey-value",
				StreamID:       "stream-id",
				PlaybackURL:    "playback-url",
				ViewerCount:    100,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.live.FillIVS(tt.params)
			assert.Equal(t, tt.expect, tt.live)
		})
	}
}

func TestLives_ProducerIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		lives  Lives
		expect []string
	}{
		{
			name: "success",
			lives: Lives{
				{
					ID:          "live-id",
					ScheduleID:  "schedule-id",
					ProducerID:  "producer-id",
					Title:       "ライブのタイトル",
					Description: "ライブの説明",
					Status:      LiveStatusCanceled,
					Published:   false,
					Canceled:    true,
					StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
					LiveProducts: LiveProducts{
						{
							LiveID:    "live-id",
							ProductID: "product-id",
							CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
							UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
						},
					},
					CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				},
			},
			expect: []string{"producer-id"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.lives.ProducerIDs())
		})
	}
}

func TestLives_Fill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		lives    Lives
		products map[string]LiveProducts
		now      time.Time
		expect   Lives
	}{
		{
			name: "success",
			lives: Lives{
				{
					ID:          "live-id",
					ScheduleID:  "schedule-id",
					ProducerID:  "producer-id",
					Title:       "ライブのタイトル",
					Description: "ライブの説明",
					Published:   false,
					Canceled:    true,
					StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
					CreatedAt:   jst.Date(2022, 7, 1, 0, 0, 0, 0),
					UpdatedAt:   jst.Date(2022, 7, 1, 0, 0, 0, 0),
				},
			},
			products: map[string]LiveProducts{
				"live-id": {
					{
						LiveID:    "live-id",
						ProductID: "product-id",
						CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
						UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					},
				},
			},
			now: jst.Date(2022, 7, 1, 0, 0, 0, 0),
			expect: Lives{
				{
					ID:          "live-id",
					ScheduleID:  "schedule-id",
					ProducerID:  "producer-id",
					Title:       "ライブのタイトル",
					Description: "ライブの説明",
					Status:      LiveStatusCanceled,
					Published:   false,
					Canceled:    true,
					StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
					LiveProducts: LiveProducts{
						{
							LiveID:    "live-id",
							ProductID: "product-id",
							CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
							UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
						},
					},
					CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.lives.Fill(tt.products, tt.now)
			assert.Equal(t, tt.expect, tt.lives)
		})
	}
}
