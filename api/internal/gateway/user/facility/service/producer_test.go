package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestProducers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		producers entity.Producers
		shops     map[string]entity.Shops
		expect    Producers
		response  []*types.Producer
	}{
		{
			name: "success",
			producers: entity.Producers{
				{
					Admin: entity.Admin{
						ID:            "producer-id01",
						Type:          entity.AdminTypeProducer,
						Status:        entity.AdminStatusActivated,
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-producer01@and-period.jp",
					},
					AdminID:        "producer-id01",
					Username:       "&.農園",
					ThumbnailURL:   "https://and-period.jp/thumbnail.png",
					HeaderURL:      "https://and-period.jp/header.png",
					PhoneNumber:    "+819012345678",
					PostalCode:     "1000014",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					City:           "千代田区",
					CreatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				{
					Admin: entity.Admin{
						ID:            "producer-id02",
						Type:          entity.AdminTypeProducer,
						Status:        entity.AdminStatusActivated,
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-producer02@and-period.jp",
					},
					AdminID:        "producer-id02",
					Username:       "&.農園",
					ThumbnailURL:   "https://and-period.jp/thumbnail.png",
					HeaderURL:      "https://and-period.jp/header.png",
					PhoneNumber:    "+819012345678",
					PostalCode:     "1000014",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					City:           "千代田区",
					CreatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			shops: map[string]entity.Shops{
				"producer-id01": {
					{
						ID:            "shop-id01",
						CoordinatorID: "coordinator-id",
						ProducerIDs:   []string{"producer-id01", "producer-id02"},
						Name:          "&.農園",
					},
				},
				"producer-id02": {
					{
						ID:            "shop-id01",
						CoordinatorID: "coordinator-id",
						ProducerIDs:   []string{"producer-id01", "producer-id02"},
						Name:          "&.農園",
					},
				},
			},
			expect: Producers{
				{
					Producer: types.Producer{
						ID:            "producer-id01",
						CoordinatorID: "coordinator-id",
						Username:      "&.農園",
						ThumbnailURL:  "https://and-period.jp/thumbnail.png",
						HeaderURL:     "https://and-period.jp/header.png",
						Prefecture:    "東京都",
						City:          "千代田区",
					},
				},
				{
					Producer: types.Producer{
						ID:            "producer-id02",
						CoordinatorID: "coordinator-id",
						Username:      "&.農園",
						ThumbnailURL:  "https://and-period.jp/thumbnail.png",
						HeaderURL:     "https://and-period.jp/header.png",
						Prefecture:    "東京都",
						City:          "千代田区",
					},
				},
			},
			response: []*types.Producer{
				{
					ID:            "producer-id01",
					CoordinatorID: "coordinator-id",
					Username:      "&.農園",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					HeaderURL:     "https://and-period.jp/header.png",
					Prefecture:    "東京都",
					City:          "千代田区",
				},
				{
					ID:            "producer-id02",
					CoordinatorID: "coordinator-id",
					Username:      "&.農園",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					HeaderURL:     "https://and-period.jp/header.png",
					Prefecture:    "東京都",
					City:          "千代田区",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProducers(tt.producers, tt.shops)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.response, actual.Response())
		})
	}
}

func TestProducers_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		producers Producers
		expect    map[string]*Producer
	}{
		{
			name: "success",
			producers: Producers{
				{
					Producer: types.Producer{
						ID:            "producer-id01",
						CoordinatorID: "coordinator-id",
						Username:      "&.農園",
						HeaderURL:     "https://and-period.jp/header.png",
						Prefecture:    "東京都",
						City:          "千代田区",
					},
				},
				{
					Producer: types.Producer{
						ID:            "producer-id02",
						CoordinatorID: "coordinator-id",
						Username:      "&.農園",
						HeaderURL:     "https://and-period.jp/header.png",
						Prefecture:    "東京都",
						City:          "千代田区",
					},
				},
			},
			expect: map[string]*Producer{
				"producer-id01": {
					Producer: types.Producer{
						ID:            "producer-id01",
						CoordinatorID: "coordinator-id",
						Username:      "&.農園",
						HeaderURL:     "https://and-period.jp/header.png",
						Prefecture:    "東京都",
						City:          "千代田区",
					},
				},
				"producer-id02": {
					Producer: types.Producer{
						ID:            "producer-id02",
						CoordinatorID: "coordinator-id",
						Username:      "&.農園",
						HeaderURL:     "https://and-period.jp/header.png",
						Prefecture:    "東京都",
						City:          "千代田区",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.producers.Map()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProducers_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		producers Producers
		expect    []*types.Producer
	}{
		{
			name: "success",
			producers: Producers{
				{
					Producer: types.Producer{
						ID:            "producer-id01",
						CoordinatorID: "coordinator-id",
						Username:      "&.農園",
						HeaderURL:     "https://and-period.jp/header.png",
						Prefecture:    "東京都",
						City:          "千代田区",
					},
				},
				{
					Producer: types.Producer{
						ID:            "producer-id02",
						CoordinatorID: "coordinator-id",
						Username:      "&.農園",
						HeaderURL:     "https://and-period.jp/header.png",
						Prefecture:    "東京都",
						City:          "千代田区",
					},
				},
			},
			expect: []*types.Producer{
				{
					ID:            "producer-id01",
					CoordinatorID: "coordinator-id",
					Username:      "&.農園",
					HeaderURL:     "https://and-period.jp/header.png",
					Prefecture:    "東京都",
					City:          "千代田区",
				},
				{
					ID:            "producer-id02",
					CoordinatorID: "coordinator-id",
					Username:      "&.農園",
					HeaderURL:     "https://and-period.jp/header.png",
					Prefecture:    "東京都",
					City:          "千代田区",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.producers.Response())
		})
	}
}
