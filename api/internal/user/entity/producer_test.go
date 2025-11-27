package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProducer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *NewProducerParams
		expect *Producer
		hasErr bool
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
				PrefectureCode:    13,
				City:              "千代田区",
				AddressLine1:      "永田町1-7-1",
				AddressLine2:      "",
			},
			expect: &Producer{
				AdminID:           "admin-id",
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
				Prefecture:        "東京都",
				PrefectureCode:    13,
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
			hasErr: false,
		},
		{
			name: "success without prefecture code",
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
				PrefectureCode:    0,
				City:              "千代田区",
				AddressLine1:      "永田町1-7-1",
				AddressLine2:      "",
			},
			expect: &Producer{
				AdminID:           "admin-id",
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
				Prefecture:        "",
				PrefectureCode:    0,
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
			hasErr: false,
		},
		{
			name: "invalid prefecture code",
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
				PrefectureCode:    -1,
				City:              "千代田区",
				AddressLine1:      "永田町1-7-1",
				AddressLine2:      "",
			},
			expect: nil,
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewProducer(tt.params)
			if tt.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.producers.IDs())
		})
	}
}

func TestProducers_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		producers Producers
		admins    map[string]*Admin
		groups    map[string]AdminGroupUsers
		expect    Producers
	}{
		{
			name: "success",
			producers: Producers{
				{
					AdminID:        "admin-id01",
					PrefectureCode: 13,
				},
				{
					AdminID:        "admin-id02",
					PrefectureCode: 13,
				},
			},
			admins: map[string]*Admin{
				"admin-id01": {
					ID:        "admin-id01",
					CognitoID: "cognito-id",
					Type:      AdminTypeProducer,
				},
			},
			groups: map[string]AdminGroupUsers{
				"admin-id01": {
					{
						GroupID: "group-id",
						AdminID: "admin-id01",
					},
				},
			},
			expect: Producers{
				{
					AdminID:        "admin-id01",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					Admin: Admin{
						ID:        "admin-id01",
						CognitoID: "cognito-id",
						Type:      AdminTypeProducer,
						Status:    AdminStatusDeactivated,
						GroupIDs:  []string{"group-id"},
					},
				},
				{
					AdminID:        "admin-id02",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					Admin: Admin{
						ID:       "admin-id02",
						Type:     AdminTypeProducer,
						Status:   AdminStatusDeactivated,
						GroupIDs: []string{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.producers.Fill(tt.admins, tt.groups)
			assert.Equal(t, tt.expect, tt.producers)
		})
	}
}
