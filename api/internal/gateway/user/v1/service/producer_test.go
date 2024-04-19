package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestProducer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		producer *entity.Producer
		expect   *Producer
	}{
		{
			name: "success",
			producer: &entity.Producer{
				Admin: entity.Admin{
					ID:            "producer-id",
					Role:          entity.AdminRoleProducer,
					Status:        entity.AdminStatusActivated,
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Email:         "test-producer@and-period.jp",
				},
				AdminID:        "producer-id",
				CoordinatorID:  "coordinator-id",
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
			expect: &Producer{
				Producer: response.Producer{
					ID:            "producer-id",
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewProducer(tt.producer))
		})
	}
}

func TestProducer_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		producer *Producer
		expect   *response.Producer
	}{
		{
			name: "success",
			producer: &Producer{
				Producer: response.Producer{
					ID:            "producer-id",
					CoordinatorID: "coordinator-id",
					Username:      "&.農園",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					HeaderURL:     "https://and-period.jp/header.png",
					Prefecture:    "東京都",
					City:          "千代田区",
				},
			},
			expect: &response.Producer{
				ID:            "producer-id",
				CoordinatorID: "coordinator-id",
				Username:      "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				HeaderURL:     "https://and-period.jp/header.png",
				Prefecture:    "東京都",
				City:          "千代田区",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.producer.Response())
		})
	}
}

func TestProducers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		producers entity.Producers
		expect    Producers
	}{
		{
			name: "success",
			producers: entity.Producers{
				{
					Admin: entity.Admin{
						ID:            "producer-id01",
						Role:          entity.AdminRoleProducer,
						Status:        entity.AdminStatusActivated,
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-producer01@and-period.jp",
					},
					AdminID:        "producer-id01",
					CoordinatorID:  "coordinator-id",
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
						Role:          entity.AdminRoleProducer,
						Status:        entity.AdminStatusActivated,
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-producer02@and-period.jp",
					},
					AdminID:        "producer-id02",
					CoordinatorID:  "coordinator-id",
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
			expect: Producers{
				{
					Producer: response.Producer{
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
					Producer: response.Producer{
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
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewProducers(tt.producers))
		})
	}
}

func TestProducers_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		producers Producers
		expect    []*response.Producer
	}{
		{
			name: "success",
			producers: Producers{
				{
					Producer: response.Producer{
						ID:            "producer-id01",
						CoordinatorID: "coordinator-id",
						Username:      "&.農園",
						HeaderURL:     "https://and-period.jp/header.png",
						Prefecture:    "東京都",
						City:          "千代田区",
					},
				},
				{
					Producer: response.Producer{
						ID:            "producer-id02",
						CoordinatorID: "coordinator-id",
						Username:      "&.農園",
						HeaderURL:     "https://and-period.jp/header.png",
						Prefecture:    "東京都",
						City:          "千代田区",
					},
				},
			},
			expect: []*response.Producer{
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.producers.Response())
		})
	}
}
