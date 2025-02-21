package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestCoordinators(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		coordinators uentity.Coordinators
		shops        map[string]*Shop
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
					MarcheName:        "&.マルシェ",
					Username:          "&.農園",
					Profile:           "紹介文です。",
					ProductTypeIDs:    []string{"product-type-ids"},
					ThumbnailURL:      "https://and-period.jp/thumbnail.png",
					HeaderURL:         "https://and-period.jp/header.png",
					PromotionVideoURL: "https://and-period.jp/promotion.mp4",
					BonusVideoURL:     "https://and-period.jp/bonus.mp4",
					InstagramID:       "instagram-id",
					FacebookID:        "facebook-id",
					PhoneNumber:       "+819012345678",
					PostalCode:        "1000014",
					PrefectureCode:    codes.PrefectureValues["tokyo"],
					City:              "千代田区",
					BusinessDays:      []time.Weekday{time.Monday, time.Wednesday, time.Friday},
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
					MarcheName:        "&.マルシェ",
					Username:          "&.農園",
					Profile:           "紹介文です。",
					ProductTypeIDs:    []string{"product-type-ids"},
					ThumbnailURL:      "https://and-period.jp/thumbnail.png",
					HeaderURL:         "https://and-period.jp/header.png",
					PromotionVideoURL: "https://and-period.jp/promotion.mp4",
					BonusVideoURL:     "https://and-period.jp/bonus.mp4",
					InstagramID:       "instagram-id",
					FacebookID:        "facebook-id",
					PhoneNumber:       "+819012345678",
					PostalCode:        "1000014",
					PrefectureCode:    codes.PrefectureValues["tokyo"],
					City:              "千代田区",
					BusinessDays:      []time.Weekday{time.Monday, time.Wednesday, time.Friday},
					CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			shops: map[string]*Shop{
				"coordinator-id01": {
					Shop: response.Shop{
						ID:             "shop-id01",
						CoordinatorID:  "coordinator-id01",
						ProducerIDs:    []string{"producer-id"},
						ProductTypeIDs: []string{"product-type-id"},
						BusinessDays:   []time.Weekday{time.Monday},
						Name:           "テスト店舗1",
						CreatedAt:      1640962800,
						UpdatedAt:      1640962800,
					},
				},
				"coordinator-id02": {
					Shop: response.Shop{
						ID:             "shop-id02",
						CoordinatorID:  "coordinator-id02",
						ProducerIDs:    []string{"producer-id"},
						ProductTypeIDs: []string{"product-type-id"},
						BusinessDays:   []time.Weekday{time.Monday},
						Name:           "テスト店舗2",
						CreatedAt:      1640962800,
						UpdatedAt:      1640962800,
					},
				},
			},
			expect: Coordinators{
				{
					Coordinator: response.Coordinator{
						ID:                "coordinator-id01",
						ShopID:            "shop-id01",
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
						Email:             "test-coordinator01@and-period.jp",
						PhoneNumber:       "+819012345678",
						PostalCode:        "1000014",
						PrefectureCode:    13,
						City:              "千代田区",
						CreatedAt:         1640962800,
						UpdatedAt:         1640962800,
					},
				},
				{
					Coordinator: response.Coordinator{
						ID:                "coordinator-id02",
						ShopID:            "shop-id02",
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
						Email:             "test-coordinator02@and-period.jp",
						PhoneNumber:       "+819012345678",
						PostalCode:        "1000014",
						PrefectureCode:    13,
						City:              "千代田区",
						CreatedAt:         1640962800,
						UpdatedAt:         1640962800,
					},
				},
			},
			response: []*response.Coordinator{
				{
					ID:                "coordinator-id01",
					ShopID:            "shop-id01",
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
					Email:             "test-coordinator01@and-period.jp",
					PhoneNumber:       "+819012345678",
					PostalCode:        "1000014",
					PrefectureCode:    13,
					City:              "千代田区",
					CreatedAt:         1640962800,
					UpdatedAt:         1640962800,
				},
				{
					ID:                "coordinator-id02",
					ShopID:            "shop-id02",
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
					Email:             "test-coordinator02@and-period.jp",
					PhoneNumber:       "+819012345678",
					PostalCode:        "1000014",
					PrefectureCode:    13,
					City:              "千代田区",
					CreatedAt:         1640962800,
					UpdatedAt:         1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewCoordinators(tt.coordinators, tt.shops)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.response, actual.Response())
		})
	}
}
