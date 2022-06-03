package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestProducer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admin  *entity.Admin
		expect *Producer
	}{
		{
			name: "success",
			admin: &entity.Admin{
				ID:            "admin-id",
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				Role:          entity.AdminRoleProducer,
				CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &Producer{
				Producer: &response.Producer{
					ID:            "admin-id",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					StoreName:     "&.農園",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					Email:         "test-admin@and-period.jp",
					PhoneNumber:   "+819012345678",
					PostalCode:    "1000014",
					Prefecture:    "東京都",
					City:          "千代田区",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewProducer(tt.admin))
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
				Producer: &response.Producer{
					ID:            "admin-id",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					StoreName:     "&.農園",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					Email:         "test-admin@and-period.jp",
					PhoneNumber:   "+819012345678",
					PostalCode:    "1000014",
					Prefecture:    "東京都",
					City:          "千代田区",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
			},
			expect: &response.Producer{
				ID:            "admin-id",
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				CreatedAt:     1640962800,
				UpdatedAt:     1640962800,
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
		name   string
		admins entity.Admins
		expect Producers
	}{
		{
			name: "success",
			admins: entity.Admins{
				{
					ID:            "admin-id01",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					StoreName:     "&.農園",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					Email:         "test-admin@and-period.jp",
					PhoneNumber:   "+819012345678",
					PostalCode:    "1000014",
					Prefecture:    "東京都",
					City:          "千代田区",
					Role:          entity.AdminRoleProducer,
					CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				{
					ID:            "admin-id02",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					StoreName:     "&.農園",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					Email:         "test-admin@and-period.jp",
					PhoneNumber:   "+819012345678",
					PostalCode:    "1000014",
					Prefecture:    "東京都",
					City:          "千代田区",
					Role:          entity.AdminRoleProducer,
					CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: Producers{
				{
					Producer: &response.Producer{
						ID:            "admin-id01",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						StoreName:     "&.農園",
						ThumbnailURL:  "https://and-period.jp/thumbnail.png",
						Email:         "test-admin@and-period.jp",
						PhoneNumber:   "+819012345678",
						PostalCode:    "1000014",
						Prefecture:    "東京都",
						City:          "千代田区",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
				{
					Producer: &response.Producer{
						ID:            "admin-id02",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						StoreName:     "&.農園",
						ThumbnailURL:  "https://and-period.jp/thumbnail.png",
						Email:         "test-admin@and-period.jp",
						PhoneNumber:   "+819012345678",
						PostalCode:    "1000014",
						Prefecture:    "東京都",
						City:          "千代田区",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewProducers(tt.admins))
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
					Producer: &response.Producer{
						ID:            "admin-id01",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						StoreName:     "&.農園",
						ThumbnailURL:  "https://and-period.jp/thumbnail.png",
						Email:         "test-admin@and-period.jp",
						PhoneNumber:   "+819012345678",
						PostalCode:    "1000014",
						Prefecture:    "東京都",
						City:          "千代田区",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
				{
					Producer: &response.Producer{
						ID:            "admin-id02",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						StoreName:     "&.農園",
						ThumbnailURL:  "https://and-period.jp/thumbnail.png",
						Email:         "test-admin@and-period.jp",
						PhoneNumber:   "+819012345678",
						PostalCode:    "1000014",
						Prefecture:    "東京都",
						City:          "千代田区",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
			expect: []*response.Producer{
				{
					ID:            "admin-id01",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					StoreName:     "&.農園",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					Email:         "test-admin@and-period.jp",
					PhoneNumber:   "+819012345678",
					PostalCode:    "1000014",
					Prefecture:    "東京都",
					City:          "千代田区",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
				{
					ID:            "admin-id02",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					StoreName:     "&.農園",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					Email:         "test-admin@and-period.jp",
					PhoneNumber:   "+819012345678",
					PostalCode:    "1000014",
					Prefecture:    "東京都",
					City:          "千代田区",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
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
