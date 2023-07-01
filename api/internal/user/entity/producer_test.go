package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestProducer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *NewProducerParams
		expect *Producer
	}{
		{
			name: "success",
			params: &NewProducerParams{
				Admin: &Admin{
					ID:            "admin-id",
					CognitoID:     "cognito-id",
					Lastname:      "&.",
					Firstname:     "スタッフ",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "すたっふ",
					Email:         "test-admin@and-period.jp",
				},
				CoordinatorID:     "coordinator-id",
				PhoneNumber:       "+819012345678",
				Username:          "&.農園",
				Profile:           "紹介文です。",
				ThumbnailURL:      "https://and-period.jp/thumbnail.png",
				HeaderURL:         "https://and-period.jp/header.png",
				PromotionVideoURL: "https://and-period.jp/promotion.mp4",
				BonusVideoURL:     "https://and-period.jp/bonus.mp4",
				InstagramID:       "instagram-id",
				FacebookID:        "facebook-id",
				PostalCode:        "1000014",
				Prefecture:        codes.PrefectureValues["tokyo"],
				City:              "千代田区",
				AddressLine1:      "永田町1-7-1",
				AddressLine2:      "",
			},
			expect: &Producer{
				AdminID:           "admin-id",
				CoordinatorID:     "coordinator-id",
				PhoneNumber:       "+819012345678",
				Username:          "&.農園",
				Profile:           "紹介文です。",
				ThumbnailURL:      "https://and-period.jp/thumbnail.png",
				HeaderURL:         "https://and-period.jp/header.png",
				PromotionVideoURL: "https://and-period.jp/promotion.mp4",
				BonusVideoURL:     "https://and-period.jp/bonus.mp4",
				InstagramID:       "instagram-id",
				FacebookID:        "facebook-id",
				PostalCode:        "1000014",
				Prefecture:        codes.PrefectureValues["tokyo"],
				City:              "千代田区",
				AddressLine1:      "永田町1-7-1",
				AddressLine2:      "",
				Admin: Admin{
					CognitoID:     "cognito-id",
					Lastname:      "&.",
					Firstname:     "スタッフ",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "すたっふ",
					Email:         "test-admin@and-period.jp",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProducer(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProducer_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		producer *Producer
		admin    *Admin
		expect   *Producer
		hasErr   bool
	}{
		{
			name: "success",
			producer: &Producer{
				AdminID:        "admin-id",
				ThumbnailsJSON: []byte(`[{"url":"http://example.com/media.png","size":1}]`),
				HeadersJSON:    []byte(`[{"url":"http://example.com/media.png","size":1}]`),
			},
			admin: &Admin{
				ID:        "admin-id",
				CognitoID: "cognito-id",
			},
			expect: &Producer{
				AdminID:        "admin-id",
				ThumbnailsJSON: []byte(`[{"url":"http://example.com/media.png","size":1}]`),
				Thumbnails: common.Images{
					{Size: common.ImageSizeSmall, URL: "http://example.com/media.png"},
				},
				HeadersJSON: []byte(`[{"url":"http://example.com/media.png","size":1}]`),
				Headers: common.Images{
					{Size: common.ImageSizeSmall, URL: "http://example.com/media.png"},
				},
				Admin: Admin{
					ID:        "admin-id",
					CognitoID: "cognito-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success empty",
			producer: &Producer{
				AdminID: "admin-id",
			},
			admin: &Admin{
				ID:        "admin-id",
				CognitoID: "cognito-id",
			},
			expect: &Producer{
				AdminID:    "admin-id",
				Thumbnails: common.Images{},
				Headers:    common.Images{},
				Admin: Admin{
					ID:        "admin-id",
					CognitoID: "cognito-id",
				},
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.producer.Fill(tt.admin)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.producer)
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
				{AdminID: "producer-id01"},
				{AdminID: "producer-id02"},
			},
			expect: []string{
				"producer-id01",
				"producer-id02",
			},
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

func TestProducers_CoordinatorIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		producers Producers
		expect    []string
	}{
		{
			name: "success",
			producers: Producers{
				{CoordinatorID: "coordinator-id"},
				{CoordinatorID: "coordinator-id"},
			},
			expect: []string{
				"coordinator-id",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.producers.CoordinatorIDs())
		})
	}
}
