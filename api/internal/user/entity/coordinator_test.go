package entity

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
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
				MarcheName:        "&.マルシェ",
				Username:          "&.農園",
				Profile:           "紹介文です。",
				ProductTypeIDs:    []string{"product-type-id"},
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
				BusinessDays:      []time.Weekday{time.Monday, time.Wednesday, time.Friday},
			},
			expect: &Coordinator{
				AdminID:           "admin-id",
				PhoneNumber:       "+819012345678",
				MarcheName:        "&.マルシェ",
				Username:          "&.農園",
				Profile:           "紹介文です。",
				ProductTypeIDs:    []string{"product-type-id"},
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
				BusinessDays:      []time.Weekday{time.Monday, time.Wednesday, time.Friday},
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
				MarcheName:        "&.マルシェ",
				Username:          "&.農園",
				Profile:           "紹介文です。",
				ProductTypeIDs:    []string{"product-type-id"},
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
				BusinessDays:      []time.Weekday{time.Monday, time.Wednesday, time.Friday},
			},
			expect: nil,
			hasErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
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

func TestCoordinator_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		coordinator *Coordinator
		admin       *Admin
		expect      *Coordinator
		hasErr      bool
	}{
		{
			name: "success",
			coordinator: &Coordinator{
				AdminID:            "admin-id",
				PrefectureCode:     13,
				ProductTypeIDsJSON: datatypes.JSON([]byte(`["product-type-id"]`)),
				BusinessDaysJSON:   datatypes.JSON([]byte(`[1,3,5]`)),
			},
			admin: &Admin{
				ID:        "admin-id",
				CognitoID: "cognito-id",
			},
			expect: &Coordinator{
				AdminID:            "admin-id",
				Prefecture:         "東京都",
				PrefectureCode:     13,
				ProductTypeIDsJSON: []byte(`["product-type-id"]`),
				ProductTypeIDs: []string{
					"product-type-id",
				},
				BusinessDays:     []time.Weekday{time.Monday, time.Wednesday, time.Friday},
				BusinessDaysJSON: datatypes.JSON([]byte(`[1,3,5]`)),
				Admin: Admin{
					ID:        "admin-id",
					CognitoID: "cognito-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success empty",
			coordinator: &Coordinator{
				AdminID: "admin-id",
			},
			admin: &Admin{
				ID:        "admin-id",
				CognitoID: "cognito-id",
			},
			expect: &Coordinator{
				AdminID:        "admin-id",
				ProductTypeIDs: []string{},
				BusinessDays:   []time.Weekday{},
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
			err := tt.coordinator.Fill(tt.admin)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.coordinator)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.coordinators.IDs())
		})
	}
}

func TestCoordinators_ProductTypeIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		coordinators Coordinators
		expect       []string
	}{
		{
			name: "success",
			coordinators: Coordinators{
				{
					AdminID:        "coordinator-id01",
					ProductTypeIDs: []string{"product-type-id01"},
				},
				{
					AdminID:        "coordinator-id02",
					ProductTypeIDs: []string{},
				},
				{
					AdminID:        "coordinator-id03",
					ProductTypeIDs: []string{"product-type-id01", "product-type-id02"},
				},
			},
			expect: []string{
				"product-type-id01",
				"product-type-id02",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.coordinators.ProductTypeIDs())
		})
	}
}

func TestCoordinators_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		coordinators Coordinators
		admins       map[string]*Admin
		expect       Coordinators
		hasErr       bool
	}{
		{
			name: "success",
			coordinators: Coordinators{
				{
					AdminID:            "admin-id01",
					PrefectureCode:     13,
					ProductTypeIDsJSON: datatypes.JSON([]byte(`["product-type-id"]`)),
					BusinessDaysJSON:   datatypes.JSON([]byte(`[1,3,5]`)),
				},
				{
					AdminID:            "admin-id02",
					PrefectureCode:     13,
					ProductTypeIDsJSON: datatypes.JSON([]byte(`["product-type-id"]`)),
					BusinessDaysJSON:   datatypes.JSON([]byte(`[1,3,5]`)),
				},
			},
			admins: map[string]*Admin{
				"admin-id01": {
					ID:        "admin-id01",
					CognitoID: "cognito-id",
					Type:      AdminTypeCoordinator,
				},
			},
			expect: Coordinators{
				{
					AdminID:            "admin-id01",
					Prefecture:         "東京都",
					PrefectureCode:     13,
					ProductTypeIDsJSON: []byte(`["product-type-id"]`),
					ProductTypeIDs: []string{
						"product-type-id",
					},
					BusinessDays:     []time.Weekday{time.Monday, time.Wednesday, time.Friday},
					BusinessDaysJSON: datatypes.JSON([]byte(`[1,3,5]`)),
					Admin: Admin{
						ID:        "admin-id01",
						CognitoID: "cognito-id",
						Type:      AdminTypeCoordinator,
					},
				},
				{
					AdminID:            "admin-id02",
					Prefecture:         "東京都",
					PrefectureCode:     13,
					ProductTypeIDsJSON: []byte(`["product-type-id"]`),
					ProductTypeIDs: []string{
						"product-type-id",
					},
					BusinessDays:     []time.Weekday{time.Monday, time.Wednesday, time.Friday},
					BusinessDaysJSON: datatypes.JSON([]byte(`[1,3,5]`)),
					Admin: Admin{
						ID:   "admin-id02",
						Type: AdminTypeCoordinator,
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
			err := tt.coordinators.Fill(tt.admins)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.coordinators)
		})
	}
}
