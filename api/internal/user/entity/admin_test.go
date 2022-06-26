package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestAdminFromAdministrator(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 25, 18, 30, 0, 0)
	tests := []struct {
		name          string
		administrator *Administrator
		expect        *Admin
	}{
		{
			name: "success",
			administrator: &Administrator{
				ID:            "admin-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			expect: &Admin{
				ID:            "admin-id",
				Role:          AdminRoleAdministrator,
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				CreatedAt:     now,
				UpdatedAt:     now,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAdminFromAdministrator(tt.administrator)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAdminFromCoordinator(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 25, 18, 30, 0, 0)
	tests := []struct {
		name        string
		coordinator *Coordinator
		expect      *Admin
	}{
		{
			name: "success",
			coordinator: &Coordinator{
				ID:               "admin-id",
				Lastname:         "&.",
				Firstname:        "スタッフ",
				LastnameKana:     "あんどぴりおど",
				FirstnameKana:    "すたっふ",
				StoreName:        "&.農園",
				ThumbnailURL:     "https://and-period.jp/thumbnail.png",
				HeaderURL:        "https://and-period.jp/header.png",
				TwitterAccount:   "twitter-account",
				InstagramAccount: "instagram-account",
				FacebookAccount:  "facebook-account",
				Email:            "test-admin@and-period.jp",
				PhoneNumber:      "+819012345678",
				PostalCode:       "1000014",
				Prefecture:       "東京都",
				City:             "千代田区",
				AddressLine1:     "永田町1-7-1",
				AddressLine2:     "",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			expect: &Admin{
				ID:            "admin-id",
				Role:          AdminRoleCoordinator,
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				CreatedAt:     now,
				UpdatedAt:     now,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAdminFromCoordinator(tt.coordinator)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAdminFromProducer(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 25, 18, 30, 0, 0)
	tests := []struct {
		name     string
		producer *Producer
		expect   *Admin
	}{
		{
			name: "success",
			producer: &Producer{
				ID:            "admin-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
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
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			expect: &Admin{
				ID:            "admin-id",
				Role:          AdminRoleProducer,
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				CreatedAt:     now,
				UpdatedAt:     now,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAdminFromProducer(tt.producer)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
