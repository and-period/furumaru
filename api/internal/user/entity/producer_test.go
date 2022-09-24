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
				StoreName:    "&.農園",
				ThumbnailURL: "https://and-period.jp/thumbnail.png",
				HeaderURL:    "https://and-period.jp/header.png",
				PhoneNumber:  "+819012345678",
				PostalCode:   "1000014",
				Prefecture:   "東京都",
				City:         "千代田区",
				AddressLine1: "永田町1-7-1",
				AddressLine2: "",
			},
			expect: &Producer{
				AdminID:      "admin-id",
				StoreName:    "&.農園",
				ThumbnailURL: "https://and-period.jp/thumbnail.png",
				HeaderURL:    "https://and-period.jp/header.png",
				PhoneNumber:  "+819012345678",
				PostalCode:   "1000014",
				Prefecture:   "東京都",
				City:         "千代田区",
				AddressLine1: "永田町1-7-1",
				AddressLine2: "",
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
	}{
		{
			name: "success",
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
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.producer.Fill(tt.admin)
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
