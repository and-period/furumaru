package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestCoordinators(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		coordinators uentity.Coordinators
		shops        map[string]*sentity.Shop
		expect       Coordinators
		response     []*response.Coordinator
	}{
		{
			name: "success",
			coordinators: uentity.Coordinators{
				{
					Admin: uentity.Admin{
						ID:            "coordinator-id01",
						Type:          uentity.AdminTypeCoordinator,
						Status:        uentity.AdminStatusActivated,
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-coordinator01@and-period.jp",
					},
					AdminID:           "coordinator-id01",
					Username:          "&.農園",
					Profile:           "紹介文です。",
					ThumbnailURL:      "https://and-period.jp/thumbnail.png",
					HeaderURL:         "https://and-period.jp/header.png",
					PromotionVideoURL: "https://and-period.jp/promotion.mp4",
					BonusVideoURL:     "https://and-period.jp/bonus.mp4",
					InstagramID:       "instagram-id",
					FacebookID:        "facebook-id",
					PhoneNumber:       "+819012345678",
					PostalCode:        "1000014",
					Prefecture:        "東京都",
					PrefectureCode:    13,
					City:              "千代田区",
					CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				{
					Admin: uentity.Admin{
						ID:            "coordinator-id02",
						Type:          uentity.AdminTypeCoordinator,
						Status:        uentity.AdminStatusActivated,
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-coordinator02@and-period.jp",
					},
					AdminID:           "coordinator-id02",
					Username:          "&.農園",
					Profile:           "紹介文です。",
					ThumbnailURL:      "https://and-period.jp/thumbnail.png",
					HeaderURL:         "https://and-period.jp/header.png",
					PromotionVideoURL: "https://and-period.jp/promotion.mp4",
					BonusVideoURL:     "https://and-period.jp/bonus.mp4",
					InstagramID:       "instagram-id",
					FacebookID:        "facebook-id",
					PhoneNumber:       "+819012345678",
					PostalCode:        "1000014",
					Prefecture:        "東京都",
					PrefectureCode:    13,
					City:              "千代田区",
					CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			shops: map[string]*sentity.Shop{
				"coordinator-id01": {
					ID:             "shop-id01",
					CoordinatorID:  "coordinator-id01",
					ProducerIDs:    []string{"producer-id01"},
					ProductTypeIDs: []string{"product-type-ids"},
					BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
					Name:           "&.マルシェ",
					Activated:      true,
					CreatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				"coordinator-id02": {
					ID:             "shop-id02",
					CoordinatorID:  "coordinator-id02",
					ProducerIDs:    []string{"producer-id02"},
					ProductTypeIDs: []string{"product-type-ids"},
					BusinessDays:   []time.Weekday{time.Monday, time.Wednesday, time.Friday},
					Name:           "&.マルシェ",
					Activated:      true,
					CreatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: Coordinators{
				{
					Coordinator: response.Coordinator{
						ID:                "coordinator-id01",
						MarcheName:        "&.マルシェ",
						Username:          "&.農園",
						Profile:           "紹介文です。",
						ProductTypeIDs:    []string{"product-type-ids"},
						BusinessDays:      []time.Weekday{time.Monday, time.Wednesday, time.Friday},
						ThumbnailURL:      "https://and-period.jp/thumbnail.png",
						HeaderURL:         "https://and-period.jp/header.png",
						PromotionVideoURL: "https://and-period.jp/promotion.mp4",
						InstagramID:       "instagram-id",
						FacebookID:        "facebook-id",
						Prefecture:        "東京都",
						City:              "千代田区",
					},
				},
				{
					Coordinator: response.Coordinator{
						ID:                "coordinator-id02",
						MarcheName:        "&.マルシェ",
						Username:          "&.農園",
						Profile:           "紹介文です。",
						ProductTypeIDs:    []string{"product-type-ids"},
						BusinessDays:      []time.Weekday{time.Monday, time.Wednesday, time.Friday},
						ThumbnailURL:      "https://and-period.jp/thumbnail.png",
						HeaderURL:         "https://and-period.jp/header.png",
						PromotionVideoURL: "https://and-period.jp/promotion.mp4",
						InstagramID:       "instagram-id",
						FacebookID:        "facebook-id",
						Prefecture:        "東京都",
						City:              "千代田区",
					},
				},
			},
			response: []*response.Coordinator{
				{
					ID:                "coordinator-id01",
					MarcheName:        "&.マルシェ",
					Username:          "&.農園",
					Profile:           "紹介文です。",
					ProductTypeIDs:    []string{"product-type-ids"},
					BusinessDays:      []time.Weekday{time.Monday, time.Wednesday, time.Friday},
					ThumbnailURL:      "https://and-period.jp/thumbnail.png",
					HeaderURL:         "https://and-period.jp/header.png",
					PromotionVideoURL: "https://and-period.jp/promotion.mp4",
					InstagramID:       "instagram-id",
					FacebookID:        "facebook-id",
					Prefecture:        "東京都",
					City:              "千代田区",
				},
				{
					ID:                "coordinator-id02",
					MarcheName:        "&.マルシェ",
					Username:          "&.農園",
					Profile:           "紹介文です。",
					ProductTypeIDs:    []string{"product-type-ids"},
					BusinessDays:      []time.Weekday{time.Monday, time.Wednesday, time.Friday},
					ThumbnailURL:      "https://and-period.jp/thumbnail.png",
					HeaderURL:         "https://and-period.jp/header.png",
					PromotionVideoURL: "https://and-period.jp/promotion.mp4",
					InstagramID:       "instagram-id",
					FacebookID:        "facebook-id",
					Prefecture:        "東京都",
					City:              "千代田区",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewCoordinators(tt.coordinators, tt.shops)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.response, actual.Response())
		})
	}
}

func TestCoordinators_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		coordinators Coordinators
		expect       map[string]*Coordinator
	}{
		{
			name: "success",
			coordinators: Coordinators{
				{
					Coordinator: response.Coordinator{
						ID:                "coordinator-id01",
						MarcheName:        "&.マルシェ",
						Username:          "&.農園",
						Profile:           "紹介文です。",
						ProductTypeIDs:    []string{"product-type-ids"},
						BusinessDays:      []time.Weekday{time.Monday, time.Wednesday, time.Friday},
						ThumbnailURL:      "https://and-period.jp/thumbnail.png",
						HeaderURL:         "https://and-period.jp/header.png",
						PromotionVideoURL: "https://and-period.jp/promotion.mp4",
						InstagramID:       "instagram-id",
						FacebookID:        "facebook-id",
						Prefecture:        "東京都",
						City:              "千代田区",
					},
				},
				{
					Coordinator: response.Coordinator{
						ID:                "coordinator-id02",
						MarcheName:        "&.マルシェ",
						Username:          "&.農園",
						Profile:           "紹介文です。",
						ProductTypeIDs:    []string{"product-type-ids"},
						BusinessDays:      []time.Weekday{time.Monday, time.Wednesday, time.Friday},
						ThumbnailURL:      "https://and-period.jp/thumbnail.png",
						HeaderURL:         "https://and-period.jp/header.png",
						PromotionVideoURL: "https://and-period.jp/promotion.mp4",
						InstagramID:       "instagram-id",
						FacebookID:        "facebook-id",
						Prefecture:        "東京都",
						City:              "千代田区",
					},
				},
			},
			expect: map[string]*Coordinator{
				"coordinator-id01": {
					Coordinator: response.Coordinator{
						ID:                "coordinator-id01",
						MarcheName:        "&.マルシェ",
						Username:          "&.農園",
						Profile:           "紹介文です。",
						ProductTypeIDs:    []string{"product-type-ids"},
						BusinessDays:      []time.Weekday{time.Monday, time.Wednesday, time.Friday},
						ThumbnailURL:      "https://and-period.jp/thumbnail.png",
						HeaderURL:         "https://and-period.jp/header.png",
						PromotionVideoURL: "https://and-period.jp/promotion.mp4",
						InstagramID:       "instagram-id",
						FacebookID:        "facebook-id",
						Prefecture:        "東京都",
						City:              "千代田区",
					},
				},
				"coordinator-id02": {
					Coordinator: response.Coordinator{
						ID:                "coordinator-id02",
						MarcheName:        "&.マルシェ",
						Username:          "&.農園",
						Profile:           "紹介文です。",
						ProductTypeIDs:    []string{"product-type-ids"},
						BusinessDays:      []time.Weekday{time.Monday, time.Wednesday, time.Friday},
						ThumbnailURL:      "https://and-period.jp/thumbnail.png",
						HeaderURL:         "https://and-period.jp/header.png",
						PromotionVideoURL: "https://and-period.jp/promotion.mp4",
						InstagramID:       "instagram-id",
						FacebookID:        "facebook-id",
						Prefecture:        "東京都",
						City:              "千代田区",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.coordinators.Map()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
