package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
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
				AdminID:       "producer-id",
				CoordinatorID: "coordinator-id",
				Username:      "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				Thumbnails: common.Images{
					{URL: "https://and-period.jp/thumbnail_240.png", Size: common.ImageSizeSmall},
					{URL: "https://and-period.jp/thumbnail_675.png", Size: common.ImageSizeMedium},
					{URL: "https://and-period.jp/thumbnail_900.png", Size: common.ImageSizeLarge},
				},
				HeaderURL: "https://and-period.jp/header.png",
				Headers: common.Images{
					{URL: "https://and-period.jp/header_240.png", Size: common.ImageSizeSmall},
					{URL: "https://and-period.jp/header_675.png", Size: common.ImageSizeMedium},
					{URL: "https://and-period.jp/header_900.png", Size: common.ImageSizeLarge},
				},
				PhoneNumber: "+819012345678",
				PostalCode:  "1000014",
				Prefecture:  "東京都",
				City:        "千代田区",
				CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &Producer{
				Producer: response.Producer{
					ID:            "producer-id",
					Status:        entity.AdminStatusActivated,
					CoordinatorID: "coordinator-id",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Username:      "&.農園",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					Thumbnails: []*response.Image{
						{URL: "https://and-period.jp/thumbnail_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/thumbnail_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/thumbnail_900.png", Size: int32(ImageSizeLarge)},
					},
					HeaderURL: "https://and-period.jp/header.png",
					Headers: []*response.Image{
						{URL: "https://and-period.jp/header_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/header_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/header_900.png", Size: int32(ImageSizeLarge)},
					},
					Email:       "test-producer@and-period.jp",
					PhoneNumber: "+819012345678",
					PostalCode:  "1000014",
					Prefecture:  "東京都",
					City:        "千代田区",
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
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
					Status:        entity.AdminStatusActivated,
					CoordinatorID: "coordinator-id",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Username:      "&.農園",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					Thumbnails: []*response.Image{
						{URL: "https://and-period.jp/thumbnail_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/thumbnail_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/thumbnail_900.png", Size: int32(ImageSizeLarge)},
					},
					HeaderURL: "https://and-period.jp/header.png",
					Headers: []*response.Image{
						{URL: "https://and-period.jp/header_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/header_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/header_900.png", Size: int32(ImageSizeLarge)},
					},
					Email:       "test-producer@and-period.jp",
					PhoneNumber: "+819012345678",
					PostalCode:  "1000014",
					Prefecture:  "東京都",
					City:        "千代田区",
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
			},
			expect: &response.Producer{
				ID:            "producer-id",
				Status:        entity.AdminStatusActivated,
				CoordinatorID: "coordinator-id",
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Username:      "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				Thumbnails: []*response.Image{
					{URL: "https://and-period.jp/thumbnail_240.png", Size: int32(ImageSizeSmall)},
					{URL: "https://and-period.jp/thumbnail_675.png", Size: int32(ImageSizeMedium)},
					{URL: "https://and-period.jp/thumbnail_900.png", Size: int32(ImageSizeLarge)},
				},
				HeaderURL: "https://and-period.jp/header.png",
				Headers: []*response.Image{
					{URL: "https://and-period.jp/header_240.png", Size: int32(ImageSizeSmall)},
					{URL: "https://and-period.jp/header_675.png", Size: int32(ImageSizeMedium)},
					{URL: "https://and-period.jp/header_900.png", Size: int32(ImageSizeLarge)},
				},
				Email:       "test-producer@and-period.jp",
				PhoneNumber: "+819012345678",
				PostalCode:  "1000014",
				Prefecture:  "東京都",
				City:        "千代田区",
				CreatedAt:   1640962800,
				UpdatedAt:   1640962800,
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

func TestProducer_Name(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		producer *Producer
		expect   string
	}{
		{
			name: "success",
			producer: &Producer{
				Producer: response.Producer{
					ID:            "producer-id",
					Status:        entity.AdminStatusActivated,
					CoordinatorID: "coordinator-id",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Username:      "&.農園",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					Thumbnails: []*response.Image{
						{URL: "https://and-period.jp/thumbnail_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/thumbnail_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/thumbnail_900.png", Size: int32(ImageSizeLarge)},
					},
					HeaderURL: "https://and-period.jp/header.png",
					Headers: []*response.Image{
						{URL: "https://and-period.jp/header_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/header_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/header_900.png", Size: int32(ImageSizeLarge)},
					},
					Email:       "test-producer@and-period.jp",
					PhoneNumber: "+819012345678",
					PostalCode:  "1000014",
					Prefecture:  "東京都",
					City:        "千代田区",
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
			},
			expect: "&. 管理者",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.producer.Name()
			assert.Equal(t, tt.expect, actual)
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
					AdminID:       "producer-id01",
					CoordinatorID: "coordinator-id",
					Username:      "&.農園",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					Thumbnails: common.Images{
						{URL: "https://and-period.jp/thumbnail_240.png", Size: common.ImageSizeSmall},
						{URL: "https://and-period.jp/thumbnail_675.png", Size: common.ImageSizeMedium},
						{URL: "https://and-period.jp/thumbnail_900.png", Size: common.ImageSizeLarge},
					},
					HeaderURL: "https://and-period.jp/header.png",
					Headers: common.Images{
						{URL: "https://and-period.jp/header_240.png", Size: common.ImageSizeSmall},
						{URL: "https://and-period.jp/header_675.png", Size: common.ImageSizeMedium},
						{URL: "https://and-period.jp/header_900.png", Size: common.ImageSizeLarge},
					},
					PhoneNumber: "+819012345678",
					PostalCode:  "1000014",
					Prefecture:  "東京都",
					City:        "千代田区",
					CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
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
					AdminID:       "producer-id02",
					CoordinatorID: "coordinator-id",
					Username:      "&.農園",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					HeaderURL:     "https://and-period.jp/header.png",
					PhoneNumber:   "+819012345678",
					PostalCode:    "1000014",
					Prefecture:    "東京都",
					City:          "千代田区",
					CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: Producers{
				{
					Producer: response.Producer{
						ID:            "producer-id01",
						Status:        entity.AdminStatusActivated,
						CoordinatorID: "coordinator-id",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Username:      "&.農園",
						ThumbnailURL:  "https://and-period.jp/thumbnail.png",
						Thumbnails: []*response.Image{
							{URL: "https://and-period.jp/thumbnail_240.png", Size: int32(ImageSizeSmall)},
							{URL: "https://and-period.jp/thumbnail_675.png", Size: int32(ImageSizeMedium)},
							{URL: "https://and-period.jp/thumbnail_900.png", Size: int32(ImageSizeLarge)},
						},
						HeaderURL: "https://and-period.jp/header.png",
						Headers: []*response.Image{
							{URL: "https://and-period.jp/header_240.png", Size: int32(ImageSizeSmall)},
							{URL: "https://and-period.jp/header_675.png", Size: int32(ImageSizeMedium)},
							{URL: "https://and-period.jp/header_900.png", Size: int32(ImageSizeLarge)},
						},
						Email:       "test-producer01@and-period.jp",
						PhoneNumber: "+819012345678",
						PostalCode:  "1000014",
						Prefecture:  "東京都",
						City:        "千代田区",
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
				},
				{
					Producer: response.Producer{
						ID:            "producer-id02",
						Status:        entity.AdminStatusActivated,
						CoordinatorID: "coordinator-id",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Username:      "&.農園",
						ThumbnailURL:  "https://and-period.jp/thumbnail.png",
						Thumbnails:    []*response.Image{},
						HeaderURL:     "https://and-period.jp/header.png",
						Headers:       []*response.Image{},
						Email:         "test-producer02@and-period.jp",
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
			assert.Equal(t, tt.expect, NewProducers(tt.producers))
		})
	}
}

func TestProducers_IDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		producers Producers
		expect    []string
	}{
		{
			name: "success",
			producers: Producers{
				{
					Producer: response.Producer{
						ID:            "producer-id01",
						Status:        entity.AdminStatusActivated,
						CoordinatorID: "coordinator-id",
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
			},
			expect: []string{"producer-id01"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.producers.IDs())
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
					Producer: response.Producer{
						ID:            "producer-id01",
						Status:        entity.AdminStatusActivated,
						CoordinatorID: "coordinator-id",
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
				{
					Producer: response.Producer{
						ID:            "producer-id02",
						Status:        entity.AdminStatusActivated,
						CoordinatorID: "coordinator-id",
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
			},
			expect: map[string]*Producer{
				"producer-id01": {
					Producer: response.Producer{
						ID:            "producer-id01",
						Status:        entity.AdminStatusActivated,
						CoordinatorID: "coordinator-id",
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
				"producer-id02": {
					Producer: response.Producer{
						ID:            "producer-id02",
						Status:        entity.AdminStatusActivated,
						CoordinatorID: "coordinator-id",
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
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.producers.Map())
		})
	}
}

func TestProducers_Contains(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		producers   Producers
		producerIDs []string
		expect      bool
	}{
		{
			name: "contain",
			producers: Producers{
				{
					Producer: response.Producer{
						ID:            "producer-id01",
						Status:        entity.AdminStatusActivated,
						CoordinatorID: "coordinator-id",
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
			},
			producerIDs: []string{"producer-id01"},
			expect:      true,
		},
		{
			name: "not contain",
			producers: Producers{
				{
					Producer: response.Producer{
						ID:            "producer-id01",
						Status:        entity.AdminStatusActivated,
						CoordinatorID: "coordinator-id",
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
			},
			producerIDs: []string{"producer-id01", "producer-id02"},
			expect:      false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.producers.Contains(tt.producerIDs...))
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
						Status:        entity.AdminStatusActivated,
						CoordinatorID: "coordinator-id",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Username:      "&.農園",
						Thumbnails: []*response.Image{
							{URL: "https://and-period.jp/thumbnail_240.png", Size: int32(ImageSizeSmall)},
							{URL: "https://and-period.jp/thumbnail_675.png", Size: int32(ImageSizeMedium)},
							{URL: "https://and-period.jp/thumbnail_900.png", Size: int32(ImageSizeLarge)},
						},
						HeaderURL: "https://and-period.jp/header.png",
						Headers: []*response.Image{
							{URL: "https://and-period.jp/header_240.png", Size: int32(ImageSizeSmall)},
							{URL: "https://and-period.jp/header_675.png", Size: int32(ImageSizeMedium)},
							{URL: "https://and-period.jp/header_900.png", Size: int32(ImageSizeLarge)},
						},
						Email:       "test-producer@and-period.jp",
						PhoneNumber: "+819012345678",
						PostalCode:  "1000014",
						Prefecture:  "東京都",
						City:        "千代田区",
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
				},
				{
					Producer: response.Producer{
						ID:            "producer-id02",
						Status:        entity.AdminStatusActivated,
						CoordinatorID: "coordinator-id",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Username:      "&.農園",
						Thumbnails: []*response.Image{
							{URL: "https://and-period.jp/thumbnail_240.png", Size: int32(ImageSizeSmall)},
							{URL: "https://and-period.jp/thumbnail_675.png", Size: int32(ImageSizeMedium)},
							{URL: "https://and-period.jp/thumbnail_900.png", Size: int32(ImageSizeLarge)},
						},
						HeaderURL: "https://and-period.jp/header.png",
						Headers: []*response.Image{
							{URL: "https://and-period.jp/header_240.png", Size: int32(ImageSizeSmall)},
							{URL: "https://and-period.jp/header_675.png", Size: int32(ImageSizeMedium)},
							{URL: "https://and-period.jp/header_900.png", Size: int32(ImageSizeLarge)},
						},
						Email:       "test-producer@and-period.jp",
						PhoneNumber: "+819012345678",
						PostalCode:  "1000014",
						Prefecture:  "東京都",
						City:        "千代田区",
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
				},
			},
			expect: []*response.Producer{
				{
					ID:            "producer-id01",
					Status:        entity.AdminStatusActivated,
					CoordinatorID: "coordinator-id",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Username:      "&.農園",
					Thumbnails: []*response.Image{
						{URL: "https://and-period.jp/thumbnail_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/thumbnail_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/thumbnail_900.png", Size: int32(ImageSizeLarge)},
					},
					HeaderURL: "https://and-period.jp/header.png",
					Headers: []*response.Image{
						{URL: "https://and-period.jp/header_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/header_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/header_900.png", Size: int32(ImageSizeLarge)},
					},
					Email:       "test-producer@and-period.jp",
					PhoneNumber: "+819012345678",
					PostalCode:  "1000014",
					Prefecture:  "東京都",
					City:        "千代田区",
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
				{
					ID:            "producer-id02",
					Status:        entity.AdminStatusActivated,
					CoordinatorID: "coordinator-id",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Username:      "&.農園",
					Thumbnails: []*response.Image{
						{URL: "https://and-period.jp/thumbnail_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/thumbnail_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/thumbnail_900.png", Size: int32(ImageSizeLarge)},
					},
					HeaderURL: "https://and-period.jp/header.png",
					Headers: []*response.Image{
						{URL: "https://and-period.jp/header_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/header_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/header_900.png", Size: int32(ImageSizeLarge)},
					},
					Email:       "test-producer@and-period.jp",
					PhoneNumber: "+819012345678",
					PostalCode:  "1000014",
					Prefecture:  "東京都",
					City:        "千代田区",
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
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
