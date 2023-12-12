package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestTopCommonLive(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	tests := []struct {
		name     string
		schedule *entity.Schedule
		products entity.Products
		expect   *TopCommonLive
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
			expect: &TopCommonLive{
				TopCommonLive: response.TopCommonLive{
					ScheduleID:    "schedule-id",
					CoordinatorID: "coordinator-id",
					Status:        int32(ScheduleStatusLive),
					Title:         "スケジュールタイトル",
					ThumbnailURL:  "https://example.com/thumbnail.png",
					Thumbnails: []*response.Image{
						{URL: "https://example.com/thumbnail_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://example.com/thumbnail_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://example.com/thumbnail_900.png", Size: int32(ImageSizeLarge)},
					},
					StartAt: 1638284400,
					EndAt:   1643641200,
					Products: []*response.TopCommonLiveProduct{
						{
							ProductID:    "product-id",
							Name:         "新鮮なじゃがいも",
							Price:        400,
							Inventory:    100,
							ThumbnailURL: "https://example.com/thumbnail01.png",
							Thumbnails: []*response.Image{
								{URL: "https://example.com/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
								{URL: "https://example.com/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
								{URL: "https://example.com/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
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
			actual := NewTopCommonLive(tt.schedule, tt.products)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestTopCommonLive_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		live   *TopCommonLive
		expect *response.TopCommonLive
	}{
		{
			name: "success",
			live: &TopCommonLive{
				TopCommonLive: response.TopCommonLive{
					ScheduleID:    "schedule-id",
					CoordinatorID: "coordinator-id",
					Status:        int32(ScheduleStatusLive),
					Title:         "スケジュールタイトル",
					ThumbnailURL:  "https://example.com/thumbnail.png",
					Thumbnails: []*response.Image{
						{URL: "https://example.com/thumbnail_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://example.com/thumbnail_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://example.com/thumbnail_900.png", Size: int32(ImageSizeLarge)},
					},
					StartAt: 1638284400,
					EndAt:   1643641200,
					Products: []*response.TopCommonLiveProduct{
						{
							ProductID:    "product-id",
							Name:         "新鮮なじゃがいも",
							Price:        400,
							Inventory:    100,
							ThumbnailURL: "https://example.com/thumbnail01.png",
							Thumbnails: []*response.Image{
								{URL: "https://example.com/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
								{URL: "https://example.com/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
								{URL: "https://example.com/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
							},
						},
					},
				},
			},
			expect: &response.TopCommonLive{
				ScheduleID:    "schedule-id",
				CoordinatorID: "coordinator-id",
				Status:        int32(ScheduleStatusLive),
				Title:         "スケジュールタイトル",
				ThumbnailURL:  "https://example.com/thumbnail.png",
				Thumbnails: []*response.Image{
					{URL: "https://example.com/thumbnail_240.png", Size: int32(ImageSizeSmall)},
					{URL: "https://example.com/thumbnail_675.png", Size: int32(ImageSizeMedium)},
					{URL: "https://example.com/thumbnail_900.png", Size: int32(ImageSizeLarge)},
				},
				StartAt: 1638284400,
				EndAt:   1643641200,
				Products: []*response.TopCommonLiveProduct{
					{
						ProductID:    "product-id",
						Name:         "新鮮なじゃがいも",
						Price:        400,
						Inventory:    100,
						ThumbnailURL: "https://example.com/thumbnail01.png",
						Thumbnails: []*response.Image{
							{URL: "https://example.com/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
							{URL: "https://example.com/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
							{URL: "https://example.com/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
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
			assert.Equal(t, tt.expect, tt.live.Response())
		})
	}
}

func TestTopCommonLives(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	tests := []struct {
		name      string
		schedules entity.Schedules
		lives     entity.Lives
		products  entity.Products
		expect    TopCommonLives
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
			expect: TopCommonLives{
				{
					TopCommonLive: response.TopCommonLive{
						ScheduleID:    "schedule-id",
						CoordinatorID: "coordinator-id",
						Status:        int32(ScheduleStatusLive),
						Title:         "スケジュールタイトル",
						ThumbnailURL:  "https://example.com/thumbnail.png",
						Thumbnails: []*response.Image{
							{URL: "https://example.com/thumbnail_240.png", Size: int32(ImageSizeSmall)},
							{URL: "https://example.com/thumbnail_675.png", Size: int32(ImageSizeMedium)},
							{URL: "https://example.com/thumbnail_900.png", Size: int32(ImageSizeLarge)},
						},
						StartAt: 1638284400,
						EndAt:   1643641200,
						Products: []*response.TopCommonLiveProduct{
							{
								ProductID:    "product-id",
								Name:         "新鮮なじゃがいも",
								Price:        400,
								Inventory:    100,
								ThumbnailURL: "https://example.com/thumbnail01.png",
								Thumbnails: []*response.Image{
									{URL: "https://example.com/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
									{URL: "https://example.com/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
									{URL: "https://example.com/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
								},
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
			actual := NewTopCommonLives(tt.schedules, tt.lives, tt.products)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestTopCommonLives_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		lives  TopCommonLives
		expect []*response.TopCommonLive
	}{
		{
			name: "success",
			lives: TopCommonLives{
				{
					TopCommonLive: response.TopCommonLive{
						ScheduleID:    "schedule-id",
						CoordinatorID: "coordinator-id",
						Status:        int32(ScheduleStatusLive),
						Title:         "スケジュールタイトル",
						ThumbnailURL:  "https://example.com/thumbnail.png",
						Thumbnails: []*response.Image{
							{URL: "https://example.com/thumbnail_240.png", Size: int32(ImageSizeSmall)},
							{URL: "https://example.com/thumbnail_675.png", Size: int32(ImageSizeMedium)},
							{URL: "https://example.com/thumbnail_900.png", Size: int32(ImageSizeLarge)},
						},
						StartAt: 1638284400,
						EndAt:   1643641200,
						Products: []*response.TopCommonLiveProduct{
							{
								ProductID:    "product-id",
								Name:         "新鮮なじゃがいも",
								Price:        400,
								Inventory:    100,
								ThumbnailURL: "https://example.com/thumbnail01.png",
								Thumbnails: []*response.Image{
									{URL: "https://example.com/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
									{URL: "https://example.com/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
									{URL: "https://example.com/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
								},
							},
						},
					},
				},
			},
			expect: []*response.TopCommonLive{
				{
					ScheduleID:    "schedule-id",
					CoordinatorID: "coordinator-id",
					Status:        int32(ScheduleStatusLive),
					Title:         "スケジュールタイトル",
					ThumbnailURL:  "https://example.com/thumbnail.png",
					Thumbnails: []*response.Image{
						{URL: "https://example.com/thumbnail_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://example.com/thumbnail_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://example.com/thumbnail_900.png", Size: int32(ImageSizeLarge)},
					},
					StartAt: 1638284400,
					EndAt:   1643641200,
					Products: []*response.TopCommonLiveProduct{
						{
							ProductID:    "product-id",
							Name:         "新鮮なじゃがいも",
							Price:        400,
							Inventory:    100,
							ThumbnailURL: "https://example.com/thumbnail01.png",
							Thumbnails: []*response.Image{
								{URL: "https://example.com/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
								{URL: "https://example.com/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
								{URL: "https://example.com/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
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
			assert.Equal(t, tt.expect, tt.lives.Response())
		})
	}
}

func TestTopCommonArchive(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	tests := []struct {
		name     string
		schedule *entity.Schedule
		expect   *TopCommonArchive
	}{
		{
			name: "success",
			schedule: &entity.Schedule{
				ID:            "schedule-id",
				CoordinatorID: "coordinator-id",
				Status:        entity.ScheduleStatusClosed,
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
			expect: &TopCommonArchive{
				TopCommonArchive: response.TopCommonArchive{
					ScheduleID:    "schedule-id",
					CoordinatorID: "coordinator-id",
					Title:         "スケジュールタイトル",
					ThumbnailURL:  "https://example.com/thumbnail.png",
					Thumbnails: []*response.Image{
						{URL: "https://example.com/thumbnail_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://example.com/thumbnail_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://example.com/thumbnail_900.png", Size: int32(ImageSizeLarge)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewTopCommonArchive(tt.schedule)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestTopCommonArchive_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		archive *TopCommonArchive
		expect  *response.TopCommonArchive
	}{
		{
			name: "success",
			archive: &TopCommonArchive{
				TopCommonArchive: response.TopCommonArchive{
					ScheduleID:    "schedule-id",
					CoordinatorID: "coordinator-id",
					Title:         "スケジュールタイトル",
					ThumbnailURL:  "https://example.com/thumbnail.png",
					Thumbnails: []*response.Image{
						{URL: "https://example.com/thumbnail_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://example.com/thumbnail_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://example.com/thumbnail_900.png", Size: int32(ImageSizeLarge)},
					},
				},
			},
			expect: &response.TopCommonArchive{
				ScheduleID:    "schedule-id",
				CoordinatorID: "coordinator-id",
				Title:         "スケジュールタイトル",
				ThumbnailURL:  "https://example.com/thumbnail.png",
				Thumbnails: []*response.Image{
					{URL: "https://example.com/thumbnail_240.png", Size: int32(ImageSizeSmall)},
					{URL: "https://example.com/thumbnail_675.png", Size: int32(ImageSizeMedium)},
					{URL: "https://example.com/thumbnail_900.png", Size: int32(ImageSizeLarge)},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.archive.Response())
		})
	}
}

func TestTopCommonArchives(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	tests := []struct {
		name      string
		schedules entity.Schedules
		expect    TopCommonArchives
	}{
		{
			name: "success",
			schedules: entity.Schedules{
				{
					ID:            "schedule-id",
					CoordinatorID: "coordinator-id",
					Status:        entity.ScheduleStatusClosed,
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
			expect: TopCommonArchives{
				{
					TopCommonArchive: response.TopCommonArchive{
						ScheduleID:    "schedule-id",
						CoordinatorID: "coordinator-id",
						Title:         "スケジュールタイトル",
						ThumbnailURL:  "https://example.com/thumbnail.png",
						Thumbnails: []*response.Image{
							{URL: "https://example.com/thumbnail_240.png", Size: int32(ImageSizeSmall)},
							{URL: "https://example.com/thumbnail_675.png", Size: int32(ImageSizeMedium)},
							{URL: "https://example.com/thumbnail_900.png", Size: int32(ImageSizeLarge)},
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
			actual := NewTopCommonArchives(tt.schedules)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestTopCommonArchives_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		archives TopCommonArchives
		expect   []*response.TopCommonArchive
	}{
		{
			name: "success",
			archives: TopCommonArchives{
				{
					TopCommonArchive: response.TopCommonArchive{
						ScheduleID:    "schedule-id",
						CoordinatorID: "coordinator-id",
						Title:         "スケジュールタイトル",
						ThumbnailURL:  "https://example.com/thumbnail.png",
						Thumbnails: []*response.Image{
							{URL: "https://example.com/thumbnail_240.png", Size: int32(ImageSizeSmall)},
							{URL: "https://example.com/thumbnail_675.png", Size: int32(ImageSizeMedium)},
							{URL: "https://example.com/thumbnail_900.png", Size: int32(ImageSizeLarge)},
						},
					},
				},
			},
			expect: []*response.TopCommonArchive{
				{
					ScheduleID:    "schedule-id",
					CoordinatorID: "coordinator-id",
					Title:         "スケジュールタイトル",
					ThumbnailURL:  "https://example.com/thumbnail.png",
					Thumbnails: []*response.Image{
						{URL: "https://example.com/thumbnail_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://example.com/thumbnail_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://example.com/thumbnail_900.png", Size: int32(ImageSizeLarge)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.archives.Response())
		})
	}
}
