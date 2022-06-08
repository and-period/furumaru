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
				ID:            "admin-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				HeaderURL:     "https://and-period.jp/header.png",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expect: &Producer{
				ID:            "admin-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				HeaderURL:     "https://and-period.jp/header.png",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProducer(tt.params)
			assert.Equal(t, tt.expect, actual)
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
			name:     "success",
			producer: &Producer{Lastname: "&.", Firstname: "スタッフ"},
			expect:   "&. スタッフ",
		},
		{
			name:     "success only lastname",
			producer: &Producer{Lastname: "&.", Firstname: ""},
			expect:   "&.",
		},
		{
			name:     "success only firstname",
			producer: &Producer{Lastname: "", Firstname: "スタッフ"},
			expect:   "スタッフ",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.producer.Name())
		})
	}
}
