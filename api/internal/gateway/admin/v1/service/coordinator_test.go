package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestCoordinator(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		coordinator *entity.Coordinator
		expect      *Coordinator
	}{
		{
			name: "success",
			coordinator: &entity.Coordinator{
				Admin: entity.Admin{
					ID:            "coordinator-id",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Email:         "test-coordinator@and-period.jp",
				},
				AdminID:          "coordinator-id",
				CompanyName:      "&.株式会社",
				StoreName:        "&.農園",
				ThumbnailURL:     "https://and-period.jp/thumbnail.png",
				HeaderURL:        "https://and-period.jp/header.png",
				TwitterAccount:   "twitter-id",
				InstagramAccount: "instagram-id",
				FacebookAccount:  "facebook-id",
				PhoneNumber:      "+819012345678",
				PostalCode:       "1000014",
				Prefecture:       "東京都",
				City:             "千代田区",
				CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &Coordinator{
				Coordinator: response.Coordinator{
					ID:               "coordinator-id",
					Lastname:         "&.",
					Firstname:        "管理者",
					LastnameKana:     "あんどどっと",
					FirstnameKana:    "かんりしゃ",
					CompanyName:      "&.株式会社",
					StoreName:        "&.農園",
					ThumbnailURL:     "https://and-period.jp/thumbnail.png",
					HeaderURL:        "https://and-period.jp/header.png",
					TwitterAccount:   "twitter-id",
					InstagramAccount: "instagram-id",
					FacebookAccount:  "facebook-id",
					Email:            "test-coordinator@and-period.jp",
					PhoneNumber:      "+819012345678",
					PostalCode:       "1000014",
					Prefecture:       "東京都",
					City:             "千代田区",
					CreatedAt:        1640962800,
					UpdatedAt:        1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewCoordinator(tt.coordinator))
		})
	}
}

func TestCoordinator_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		coordinator *Coordinator
		expect      *response.Coordinator
	}{
		{
			name: "success",
			coordinator: &Coordinator{
				Coordinator: response.Coordinator{
					ID:               "coordinator-id",
					Lastname:         "&.",
					Firstname:        "管理者",
					LastnameKana:     "あんどどっと",
					FirstnameKana:    "かんりしゃ",
					CompanyName:      "&.株式会社",
					StoreName:        "&.農園",
					ThumbnailURL:     "https://and-period.jp/thumbnail.png",
					HeaderURL:        "https://and-period.jp/header.png",
					TwitterAccount:   "twitter-id",
					InstagramAccount: "instagram-id",
					FacebookAccount:  "facebook-id",
					Email:            "test-coordinator@and-period.jp",
					PhoneNumber:      "+819012345678",
					PostalCode:       "1000014",
					Prefecture:       "東京都",
					City:             "千代田区",
					CreatedAt:        1640962800,
					UpdatedAt:        1640962800,
				},
			},
			expect: &response.Coordinator{
				ID:               "coordinator-id",
				Lastname:         "&.",
				Firstname:        "管理者",
				LastnameKana:     "あんどどっと",
				FirstnameKana:    "かんりしゃ",
				CompanyName:      "&.株式会社",
				StoreName:        "&.農園",
				ThumbnailURL:     "https://and-period.jp/thumbnail.png",
				HeaderURL:        "https://and-period.jp/header.png",
				TwitterAccount:   "twitter-id",
				InstagramAccount: "instagram-id",
				FacebookAccount:  "facebook-id",
				Email:            "test-coordinator@and-period.jp",
				PhoneNumber:      "+819012345678",
				PostalCode:       "1000014",
				Prefecture:       "東京都",
				City:             "千代田区",
				CreatedAt:        1640962800,
				UpdatedAt:        1640962800,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.coordinator.Response())
		})
	}
}

func TestCoordinators(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		coordinators entity.Coordinators
		expect       Coordinators
	}{
		{
			name: "success",
			coordinators: entity.Coordinators{
				{
					Admin: entity.Admin{
						ID:            "coordinator-id01",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-coordinator01@and-period.jp",
					},
					AdminID:          "coordinator-id01",
					CompanyName:      "&.株式会社",
					StoreName:        "&.農園",
					ThumbnailURL:     "https://and-period.jp/thumbnail.png",
					HeaderURL:        "https://and-period.jp/header.png",
					TwitterAccount:   "twitter-id",
					InstagramAccount: "instagram-id",
					FacebookAccount:  "facebook-id",
					PhoneNumber:      "+819012345678",
					PostalCode:       "1000014",
					Prefecture:       "東京都",
					City:             "千代田区",
					CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				{
					Admin: entity.Admin{
						ID:            "coordinator-id02",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-coordinator02@and-period.jp",
					},
					AdminID:          "coordinator-id02",
					CompanyName:      "&.株式会社",
					StoreName:        "&.農園",
					ThumbnailURL:     "https://and-period.jp/thumbnail.png",
					HeaderURL:        "https://and-period.jp/header.png",
					TwitterAccount:   "twitter-id",
					InstagramAccount: "instagram-id",
					FacebookAccount:  "facebook-id",
					PhoneNumber:      "+819012345678",
					PostalCode:       "1000014",
					Prefecture:       "東京都",
					City:             "千代田区",
					CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: Coordinators{
				{
					Coordinator: response.Coordinator{
						ID:               "coordinator-id01",
						Lastname:         "&.",
						Firstname:        "管理者",
						LastnameKana:     "あんどどっと",
						FirstnameKana:    "かんりしゃ",
						CompanyName:      "&.株式会社",
						StoreName:        "&.農園",
						ThumbnailURL:     "https://and-period.jp/thumbnail.png",
						HeaderURL:        "https://and-period.jp/header.png",
						TwitterAccount:   "twitter-id",
						InstagramAccount: "instagram-id",
						FacebookAccount:  "facebook-id",
						Email:            "test-coordinator01@and-period.jp",
						PhoneNumber:      "+819012345678",
						PostalCode:       "1000014",
						Prefecture:       "東京都",
						City:             "千代田区",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
				{
					Coordinator: response.Coordinator{
						ID:               "coordinator-id02",
						Lastname:         "&.",
						Firstname:        "管理者",
						LastnameKana:     "あんどどっと",
						FirstnameKana:    "かんりしゃ",
						CompanyName:      "&.株式会社",
						StoreName:        "&.農園",
						ThumbnailURL:     "https://and-period.jp/thumbnail.png",
						HeaderURL:        "https://and-period.jp/header.png",
						TwitterAccount:   "twitter-id",
						InstagramAccount: "instagram-id",
						FacebookAccount:  "facebook-id",
						Email:            "test-coordinator02@and-period.jp",
						PhoneNumber:      "+819012345678",
						PostalCode:       "1000014",
						Prefecture:       "東京都",
						City:             "千代田区",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewCoordinators(tt.coordinators))
		})
	}
}

func TestCoordinators_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		coordinators Coordinators
		expect       []*response.Coordinator
	}{
		{
			name: "success",
			coordinators: Coordinators{
				{
					Coordinator: response.Coordinator{
						ID:               "coordinator-id01",
						Lastname:         "&.",
						Firstname:        "管理者",
						LastnameKana:     "あんどどっと",
						FirstnameKana:    "かんりしゃ",
						CompanyName:      "&.株式会社",
						StoreName:        "&.農園",
						ThumbnailURL:     "https://and-period.jp/thumbnail.png",
						HeaderURL:        "https://and-period.jp/header.png",
						TwitterAccount:   "twitter-id",
						InstagramAccount: "instagram-id",
						FacebookAccount:  "facebook-id",
						Email:            "test-coordinator@and-period.jp",
						PhoneNumber:      "+819012345678",
						PostalCode:       "1000014",
						Prefecture:       "東京都",
						City:             "千代田区",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
				{
					Coordinator: response.Coordinator{
						ID:               "coordinator-id02",
						Lastname:         "&.",
						Firstname:        "管理者",
						LastnameKana:     "あんどどっと",
						FirstnameKana:    "かんりしゃ",
						CompanyName:      "&.株式会社",
						StoreName:        "&.農園",
						ThumbnailURL:     "https://and-period.jp/thumbnail.png",
						HeaderURL:        "https://and-period.jp/header.png",
						TwitterAccount:   "twitter-id",
						InstagramAccount: "instagram-id",
						FacebookAccount:  "facebook-id",
						Email:            "test-coordinator@and-period.jp",
						PhoneNumber:      "+819012345678",
						PostalCode:       "1000014",
						Prefecture:       "東京都",
						City:             "千代田区",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
			expect: []*response.Coordinator{
				{
					ID:               "coordinator-id01",
					Lastname:         "&.",
					Firstname:        "管理者",
					LastnameKana:     "あんどどっと",
					FirstnameKana:    "かんりしゃ",
					CompanyName:      "&.株式会社",
					StoreName:        "&.農園",
					ThumbnailURL:     "https://and-period.jp/thumbnail.png",
					HeaderURL:        "https://and-period.jp/header.png",
					TwitterAccount:   "twitter-id",
					InstagramAccount: "instagram-id",
					FacebookAccount:  "facebook-id",
					Email:            "test-coordinator@and-period.jp",
					PhoneNumber:      "+819012345678",
					PostalCode:       "1000014",
					Prefecture:       "東京都",
					City:             "千代田区",
					CreatedAt:        1640962800,
					UpdatedAt:        1640962800,
				},
				{
					ID:               "coordinator-id02",
					Lastname:         "&.",
					Firstname:        "管理者",
					LastnameKana:     "あんどどっと",
					FirstnameKana:    "かんりしゃ",
					CompanyName:      "&.株式会社",
					StoreName:        "&.農園",
					ThumbnailURL:     "https://and-period.jp/thumbnail.png",
					HeaderURL:        "https://and-period.jp/header.png",
					TwitterAccount:   "twitter-id",
					InstagramAccount: "instagram-id",
					FacebookAccount:  "facebook-id",
					Email:            "test-coordinator@and-period.jp",
					PhoneNumber:      "+819012345678",
					PostalCode:       "1000014",
					Prefecture:       "東京都",
					City:             "千代田区",
					CreatedAt:        1640962800,
					UpdatedAt:        1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.coordinators.Response())
		})
	}
}
