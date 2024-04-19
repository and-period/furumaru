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
				PrefectureCode:    13,
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
				PrefectureCode:    0,
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
		tt := tt
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
				PrefectureCode: 13,
			},
			admin: &Admin{
				ID:        "admin-id",
				CognitoID: "cognito-id",
			},
			expect: &Producer{
				AdminID:        "admin-id",
				Prefecture:     "東京都",
				PrefectureCode: 13,
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
				AdminID: "admin-id",
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

func TestProducers_Unrelated(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		producers Producers
		expect    Producers
	}{
		{
			name: "success",
			producers: Producers{
				{
					AdminID:       "admin-id01",
					CoordinatorID: "coordinator-id",
				},
				{
					AdminID:       "admin-id02",
					CoordinatorID: "",
				},
			},
			expect: Producers{
				{
					AdminID:       "admin-id02",
					CoordinatorID: "",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.producers.Unrelated())
		})
	}
}

func TestProducers_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		producers Producers
		admins    map[string]*Admin
		expect    Producers
		hasErr    bool
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
					Role:      AdminRoleProducer,
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
						Role:      AdminRoleProducer,
					},
				},
				{
					AdminID:        "admin-id02",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					Admin: Admin{
						ID:   "admin-id02",
						Role: AdminRoleProducer,
					},
				},
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.producers.Fill(tt.admins)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.producers)
		})
	}
}
