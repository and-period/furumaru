package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestLiveSummary(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	tests := []struct {
		name     string
		schedule *entity.Schedule
		products entity.Products
		expect   *LiveSummary
	}{
		{
			name: "success",
			schedule: &entity.Schedule{
				ID:            "schedule-id",
				CoordinatorID: "coordinator-id",
				Status:        entity.ScheduleStatusLive,
				Title:         "スケジュールタイトル",
				Description:   "スケジュールの詳細です。",
				ThumbnailURL:  "https://example.com/thumbnail.png",
				Thumbnails: common.Images{
					{URL: "https://example.com/thumbnail_240.png", Size: common.ImageSizeSmall},
					{URL: "https://example.com/thumbnail_675.png", Size: common.ImageSizeMedium},
					{URL: "https://example.com/thumbnail_900.png", Size: common.ImageSizeLarge},
				},
				ImageURL:        "https://example.com/image.png",
				OpeningVideoURL: "https://example.com/opening-video.mp4",
				Public:          true,
				Approved:        true,
				ApprovedAdminID: "admin-id",
				StartAt:         now.AddDate(0, -1, 0),
				EndAt:           now.AddDate(0, 1, 0),
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			products: entity.Products{
				{
					ID:              "product-id",
					TypeID:          "product-type-id",
					ProducerID:      "producer-id",
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Public:          true,
					Status:          entity.ProductStatusForSale,
					Inventory:       100,
					Weight:          1300,
					WeightUnit:      entity.WeightUnitGram,
					Item:            1,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: entity.MultiProductMedia{
						{
							URL:         "https://example.com/thumbnail01.png",
							IsThumbnail: true,
							Images: common.Images{
								{URL: "https://example.com/thumbnail01_240.png", Size: common.ImageSizeSmall},
								{URL: "https://example.com/thumbnail01_675.png", Size: common.ImageSizeMedium},
								{URL: "https://example.com/thumbnail01_900.png", Size: common.ImageSizeLarge},
							},
						},
						{
							URL:         "https://example.com/thumbnail02.png",
							IsThumbnail: false,
							Images: common.Images{
								{URL: "https://example.com/thumbnail02_240.png", Size: common.ImageSizeSmall},
								{URL: "https://example.com/thumbnail02_675.png", Size: common.ImageSizeMedium},
								{URL: "https://example.com/thumbnail02_900.png", Size: common.ImageSizeLarge},
							},
						},
					},
					DeliveryType:     entity.DeliveryTypeNormal,
					Box60Rate:        50,
					Box80Rate:        40,
					Box100Rate:       30,
					OriginPrefecture: "滋賀県",
					OriginCity:       "彦根市",
					StartAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
					EndAt:            jst.Date(2022, 1, 1, 0, 0, 0, 0),
					ProductRevision: entity.ProductRevision{
						ID:        1,
						ProductID: "product-id",
						Price:     400,
						Cost:      300,
						CreatedAt: now,
						UpdatedAt: now,
					},
					CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: &LiveSummary{
				LiveSummary: response.LiveSummary{
					ScheduleID:    "schedule-id",
					CoordinatorID: "coordinator-id",
					Status:        int32(ScheduleStatusLive),
					Title:         "スケジュールタイトル",
					ThumbnailURL:  "https://example.com/thumbnail.png",
					StartAt:       1638284400,
					EndAt:         1643641200,
					Products: []*response.LiveProduct{
						{
							ProductID:    "product-id",
							Name:         "新鮮なじゃがいも",
							Price:        400,
							Inventory:    100,
							ThumbnailURL: "https://example.com/thumbnail01.png",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewLiveSummary(tt.schedule, tt.products)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestLiveSummary_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		live   *LiveSummary
		expect *response.LiveSummary
	}{
		{
			name: "success",
			live: &LiveSummary{
				LiveSummary: response.LiveSummary{
					ScheduleID:    "schedule-id",
					CoordinatorID: "coordinator-id",
					Status:        int32(ScheduleStatusLive),
					Title:         "スケジュールタイトル",
					ThumbnailURL:  "https://example.com/thumbnail.png",
					StartAt:       1638284400,
					EndAt:         1643641200,
					Products: []*response.LiveProduct{
						{
							ProductID:    "product-id",
							Name:         "新鮮なじゃがいも",
							Price:        400,
							Inventory:    100,
							ThumbnailURL: "https://example.com/thumbnail01.png",
						},
					},
				},
			},
			expect: &response.LiveSummary{
				ScheduleID:    "schedule-id",
				CoordinatorID: "coordinator-id",
				Status:        int32(ScheduleStatusLive),
				Title:         "スケジュールタイトル",
				ThumbnailURL:  "https://example.com/thumbnail.png",
				StartAt:       1638284400,
				EndAt:         1643641200,
				Products: []*response.LiveProduct{
					{
						ProductID:    "product-id",
						Name:         "新鮮なじゃがいも",
						Price:        400,
						Inventory:    100,
						ThumbnailURL: "https://example.com/thumbnail01.png",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.live.Response())
		})
	}
}

func TestLiveSummaries(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	tests := []struct {
		name      string
		schedules entity.Schedules
		lives     entity.Lives
		products  entity.Products
		expect    LiveSummaries
	}{
		{
			name: "success",
			schedules: entity.Schedules{
				{
					ID:            "schedule-id",
					CoordinatorID: "coordinator-id",
					Status:        entity.ScheduleStatusLive,
					Title:         "スケジュールタイトル",
					Description:   "スケジュールの詳細です。",
					ThumbnailURL:  "https://example.com/thumbnail.png",
					Thumbnails: common.Images{
						{URL: "https://example.com/thumbnail_240.png", Size: common.ImageSizeSmall},
						{URL: "https://example.com/thumbnail_675.png", Size: common.ImageSizeMedium},
						{URL: "https://example.com/thumbnail_900.png", Size: common.ImageSizeLarge},
					},
					ImageURL:        "https://example.com/image.png",
					OpeningVideoURL: "https://example.com/opening-video.mp4",
					Public:          true,
					Approved:        true,
					ApprovedAdminID: "admin-id",
					StartAt:         now.AddDate(0, -1, 0),
					EndAt:           now.AddDate(0, 1, 0),
					CreatedAt:       now,
					UpdatedAt:       now,
				},
			},
			lives: entity.Lives{
				{
					ID:         "live-id",
					ScheduleID: "schedule-id",
					ProducerID: "producer-id",
					ProductIDs: []string{"product-id"},
					Comment:    "よろしくお願いします。",
					StartAt:    now.AddDate(0, -1, 0),
					EndAt:      now.AddDate(0, 1, 0),
					CreatedAt:  now,
					UpdatedAt:  now,
				},
			},
			products: entity.Products{
				{
					ID:              "product-id",
					TypeID:          "product-type-id",
					ProducerID:      "producer-id",
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Public:          true,
					Status:          entity.ProductStatusForSale,
					Inventory:       100,
					Weight:          1300,
					WeightUnit:      entity.WeightUnitGram,
					Item:            1,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: entity.MultiProductMedia{
						{
							URL:         "https://example.com/thumbnail01.png",
							IsThumbnail: true,
							Images: common.Images{
								{URL: "https://example.com/thumbnail01_240.png", Size: common.ImageSizeSmall},
								{URL: "https://example.com/thumbnail01_675.png", Size: common.ImageSizeMedium},
								{URL: "https://example.com/thumbnail01_900.png", Size: common.ImageSizeLarge},
							},
						},
						{
							URL:         "https://example.com/thumbnail02.png",
							IsThumbnail: false,
							Images: common.Images{
								{URL: "https://example.com/thumbnail02_240.png", Size: common.ImageSizeSmall},
								{URL: "https://example.com/thumbnail02_675.png", Size: common.ImageSizeMedium},
								{URL: "https://example.com/thumbnail02_900.png", Size: common.ImageSizeLarge},
							},
						},
					},
					DeliveryType:     entity.DeliveryTypeNormal,
					Box60Rate:        50,
					Box80Rate:        40,
					Box100Rate:       30,
					OriginPrefecture: "滋賀県",
					OriginCity:       "彦根市",
					StartAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
					EndAt:            jst.Date(2022, 1, 1, 0, 0, 0, 0),
					ProductRevision: entity.ProductRevision{
						ID:        1,
						ProductID: "product-id",
						Price:     400,
						Cost:      300,
						CreatedAt: now,
						UpdatedAt: now,
					},
					CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: LiveSummaries{
				{
					LiveSummary: response.LiveSummary{
						ScheduleID:    "schedule-id",
						CoordinatorID: "coordinator-id",
						Status:        int32(ScheduleStatusLive),
						Title:         "スケジュールタイトル",
						ThumbnailURL:  "https://example.com/thumbnail.png",
						StartAt:       1638284400,
						EndAt:         1643641200,
						Products: []*response.LiveProduct{
							{
								ProductID:    "product-id",
								Name:         "新鮮なじゃがいも",
								Price:        400,
								Inventory:    100,
								ThumbnailURL: "https://example.com/thumbnail01.png",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewLiveSummaries(tt.schedules, tt.lives, tt.products)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestLiveSummaries_CoordinatorIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		lives  LiveSummaries
		expect []string
	}{
		{
			name: "success",
			lives: LiveSummaries{
				{
					LiveSummary: response.LiveSummary{
						ScheduleID:    "schedule-id",
						CoordinatorID: "coordinator-id",
						Status:        int32(ScheduleStatusLive),
						Title:         "スケジュールタイトル",
						ThumbnailURL:  "https://example.com/thumbnail.png",
						StartAt:       1638284400,
						EndAt:         1643641200,
						Products: []*response.LiveProduct{
							{
								ProductID:    "product-id",
								Name:         "新鮮なじゃがいも",
								Price:        400,
								Inventory:    100,
								ThumbnailURL: "https://example.com/thumbnail01.png",
							},
						},
					},
				},
			},
			expect: []string{"coordinator-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.lives.CoordinatorIDs())
		})
	}
}

func TestLiveSummaries_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		lives  LiveSummaries
		expect []*response.LiveSummary
	}{
		{
			name: "success",
			lives: LiveSummaries{
				{
					LiveSummary: response.LiveSummary{
						ScheduleID:    "schedule-id",
						CoordinatorID: "coordinator-id",
						Status:        int32(ScheduleStatusLive),
						Title:         "スケジュールタイトル",
						ThumbnailURL:  "https://example.com/thumbnail.png",
						StartAt:       1638284400,
						EndAt:         1643641200,
						Products: []*response.LiveProduct{
							{
								ProductID:    "product-id",
								Name:         "新鮮なじゃがいも",
								Price:        400,
								Inventory:    100,
								ThumbnailURL: "https://example.com/thumbnail01.png",
							},
						},
					},
				},
			},
			expect: []*response.LiveSummary{
				{
					ScheduleID:    "schedule-id",
					CoordinatorID: "coordinator-id",
					Status:        int32(ScheduleStatusLive),
					Title:         "スケジュールタイトル",
					ThumbnailURL:  "https://example.com/thumbnail.png",
					StartAt:       1638284400,
					EndAt:         1643641200,
					Products: []*response.LiveProduct{
						{
							ProductID:    "product-id",
							Name:         "新鮮なじゃがいも",
							Price:        400,
							Inventory:    100,
							ThumbnailURL: "https://example.com/thumbnail01.png",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.lives.Response())
		})
	}
}
