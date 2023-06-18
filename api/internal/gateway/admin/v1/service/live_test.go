package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestLive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		live   *entity.Live
		expect *Live
	}{
		{
			name: "success",
			live: &entity.Live{
				LiveProducts: entity.LiveProducts{
					{
						LiveID:    "live-id",
						ProductID: "product-id",
						CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
						UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					},
				},
				ID:          "live-id",
				ScheduleID:  "schedule-id",
				ProducerID:  "producer-id",
				Title:       "配信タイトル",
				Description: "配信の説明",
				Status:      1,
				StartAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
				CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &Live{
				Live: response.Live{
					ID:          "live-id",
					ScheduleID:  "schedule-id",
					ProducerID:  "producer-id",
					Title:       "配信タイトル",
					Description: "配信の説明",
					Status:      1,
					StartAt:     1640962800,
					EndAt:       1640962800,
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
				ProductIDs: []string{"product-id"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewLive(tt.live))
		})
	}
}

func TestLive_Fill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		live     *Live
		producer *Producer
		shipping *Shipping
		products map[string]*Product
		expect   *Live
	}{
		{
			name: "success",
			live: &Live{
				Live: response.Live{
					ID:           "live-id",
					ScheduleID:   "schedule-id",
					Title:        "配信タイトル",
					Description:  "配信の説明",
					Status:       1,
					ProducerID:   "producer-id",
					ProducerName: "&. 管理者",
					StartAt:      1640962800,
					EndAt:        1640962800,
					CreatedAt:    1640962800,
					UpdatedAt:    1640962800,
				},
				ProductIDs: []string{"product-id"},
			},
			producer: &Producer{
				Producer: response.Producer{
					ID:            "producer-id",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Username:      "&.農園",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					HeaderURL:     "https://and-period.jp/header.png",
					Email:         "test-producer@and-period.jp",
					PhoneNumber:   "+819012345678",
					PostalCode:    "1000014",
					Prefecture:    "東京都",
					City:          "千代田区",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
			},
			shipping: &Shipping{
				Shipping: response.Shipping{
					ID:   "shipping-id",
					Name: "デフォルト配送設定",
					Box60Rates: []*response.ShippingRate{
						{Number: 1, Name: "東京都", Price: 0, Prefectures: []string{"tokyo"}},
					},
					Box60Refrigerated:  500,
					Box60Frozen:        800,
					Box80Rates:         []*response.ShippingRate{},
					Box80Refrigerated:  500,
					Box80Frozen:        800,
					Box100Rates:        []*response.ShippingRate{},
					Box100Refrigerated: 500,
					Box100Frozen:       800,
					HasFreeShipping:    true,
					FreeShippingRates:  3000,
					CreatedAt:          1640962800,
					UpdatedAt:          1640962800,
				},
			},
			products: map[string]*Product{
				"product-id": {
					Product: response.Product{
						ID:              "product-id",
						TypeID:          "product-type-id",
						TypeName:        "じゃがいも",
						TypeIconURL:     "https://and-period.jp/icon.png",
						CategoryID:      "category-id",
						CategoryName:    "野菜",
						ProducerID:      "producer-id",
						StoreName:       "&.農園",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:            400,
						DeliveryType:     1,
						Box60Rate:        50,
						Box80Rate:        40,
						Box100Rate:       30,
						OriginPrefecture: "滋賀県",
						OriginCity:       "彦根市",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
			expect: &Live{
				Live: response.Live{
					ID:          "live-id",
					ScheduleID:  "schedule-id",
					Title:       "配信タイトル",
					Description: "配信の説明",
					Status:      1,
					Products: []*response.Product{
						{
							ID:              "product-id",
							TypeID:          "product-type-id",
							TypeName:        "じゃがいも",
							TypeIconURL:     "https://and-period.jp/icon.png",
							CategoryID:      "category-id",
							CategoryName:    "野菜",
							ProducerID:      "producer-id",
							StoreName:       "&.農園",
							Name:            "新鮮なじゃがいも",
							Description:     "新鮮なじゃがいもをお届けします。",
							Public:          true,
							Inventory:       100,
							Weight:          1.3,
							ItemUnit:        "袋",
							ItemDescription: "1袋あたり100gのじゃがいも",
							Media: []*response.ProductMedia{
								{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
								{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
							},
							Price:            400,
							DeliveryType:     1,
							Box60Rate:        50,
							Box80Rate:        40,
							Box100Rate:       30,
							OriginPrefecture: "滋賀県",
							OriginCity:       "彦根市",
							CreatedAt:        1640962800,
							UpdatedAt:        1640962800,
						},
					},
					ProducerID:   "producer-id",
					ProducerName: "&. 管理者",
					StartAt:      1640962800,
					EndAt:        1640962800,
					CreatedAt:    1640962800,
					UpdatedAt:    1640962800,
				},
				ProductIDs: []string{"product-id"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.live.Fill(tt.producer, tt.products)
			assert.Equal(t, tt.expect, tt.live)
		})
	}
}
