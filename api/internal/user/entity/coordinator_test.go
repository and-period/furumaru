package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/stretchr/testify/assert"
)

func TestCoordinator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *NewCoordinatorParams
		expect *Coordinator
		hasErr bool
	}{
		{
			name: "success",
			params: &NewCoordinatorParams{
				Admin: &Admin{
					ID:            "admin-id",
					CognitoID:     "cognito-id",
					Lastname:      "&.",
					Firstname:     "スタッフ",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "すたっふ",
					Email:         "test-coordinator@and-period.jp",
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
				PrefectureCode:    codes.PrefectureValues["tokyo"],
				City:              "千代田区",
				AddressLine1:      "永田町1-7-1",
				AddressLine2:      "",
			},
			expect: &Coordinator{
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
					Email:         "test-coordinator@and-period.jp",
				},
			},
			hasErr: false,
		},
		{
			name: "invalid prefecture code",
			params: &NewCoordinatorParams{
				Admin: &Admin{
					ID:            "admin-id",
					CognitoID:     "cognito-id",
					Lastname:      "&.",
					Firstname:     "スタッフ",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "すたっふ",
					Email:         "test-coordinator@and-period.jp",
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
			actual, err := NewCoordinator(tt.params)
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

func TestCoordinators_IDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		coordinators Coordinators
		expect       []string
	}{
		{
			name: "success",
			coordinators: Coordinators{
				{AdminID: "coordinator-id01"},
				{AdminID: "coordinator-id02"},
			},
			expect: []string{
				"coordinator-id01",
				"coordinator-id02",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.coordinators.IDs())
		})
	}
}

func TestCoordinators_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		coordinators Coordinators
		admins       map[string]*Admin
		groups       map[string]AdminGroupUsers
		expect       Coordinators
	}{
		{
			name: "success",
			coordinators: Coordinators{
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
					Type:      AdminTypeCoordinator,
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
			expect: Coordinators{
				{
					AdminID:        "admin-id01",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					Admin: Admin{
						ID:        "admin-id01",
						CognitoID: "cognito-id",
						Type:      AdminTypeCoordinator,
						Status:    AdminStatusInvited,
						GroupIDs:  []string{"group-id"},
					},
				},
				{
					AdminID:        "admin-id02",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					Admin: Admin{
						ID:       "admin-id02",
						Type:     AdminTypeCoordinator,
						Status:   AdminStatusInvited,
						GroupIDs: []string{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.coordinators.Fill(tt.admins, tt.groups)
			assert.Equal(t, tt.expect, tt.coordinators)
		})
	}
}
