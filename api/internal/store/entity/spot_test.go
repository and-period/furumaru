package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSpotByUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *SpotParams
		expect *Spot
		hasErr bool
	}{
		{
			name: "success",
			params: &SpotParams{
				SpotTypeID:   "spot-type-id",
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "おすすめの観光地です。",
				ThumbnailURL: "https://example.com/image.jpg",
				Longitude:    139.6917,
				Latitude:     35.6895,
				PostalCode:   "100-0001",
				Prefecture:   "東京都",
				City:         "千代田区",
				AddressLine1: "千代田1-1",
				AddressLine2: "",
			},
			expect: &Spot{
				TypeID:         "spot-type-id",
				UserType:       SpotUserTypeUser,
				UserID:         "user-id",
				Name:           "東京タワー",
				Description:    "おすすめの観光地です。",
				ThumbnailURL:   "https://example.com/image.jpg",
				Approved:       true,
				Longitude:      139.6917,
				Latitude:       35.6895,
				PostalCode:     "100-0001",
				Prefecture:     "東京都",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "千代田1-1",
				AddressLine2:   "",
			},
			hasErr: false,
		},
		{
			name: "invalid prefecture",
			params: &SpotParams{
				SpotTypeID:   "spot-type-id",
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "おすすめの観光地です。",
				ThumbnailURL: "https://example.com/image.jpg",
				Longitude:    200.0,
				Latitude:     35.6895,
				PostalCode:   "",
				Prefecture:   "",
				City:         "",
				AddressLine1: "",
				AddressLine2: "",
			},
			expect: nil,
			hasErr: true,
		},
		{
			name: "invalid longitude",
			params: &SpotParams{
				SpotTypeID:   "spot-type-id",
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "おすすめの観光地です。",
				ThumbnailURL: "https://example.com/image.jpg",
				Longitude:    200.0,
				Latitude:     35.6895,
				PostalCode:   "100-0001",
				Prefecture:   "東京都",
				City:         "千代田区",
				AddressLine1: "千代田1-1",
				AddressLine2: "",
			},
			expect: nil,
			hasErr: true,
		},
		{
			name: "invalid latitude",
			params: &SpotParams{
				SpotTypeID:   "spot-type-id",
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "おすすめの観光地です。",
				ThumbnailURL: "https://example.com/image.jpg",
				Longitude:    139.6917,
				Latitude:     100.0,
				PostalCode:   "100-0001",
				Prefecture:   "東京都",
				City:         "千代田区",
				AddressLine1: "千代田1-1",
				AddressLine2: "",
			},
			expect: nil,
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewSpotByUser(tt.params)
			if err != nil {
				assert.True(t, tt.hasErr, err)
				return
			}
			assert.False(t, tt.hasErr)

			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestSpotByAdmin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *SpotParams
		expect *Spot
		hasErr bool
	}{
		{
			name: "success",
			params: &SpotParams{
				SpotTypeID:   "spot-type-id",
				UserType:     SpotUserTypeCoordinator,
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "おすすめの観光地です。",
				ThumbnailURL: "https://example.com/image.jpg",
				Longitude:    139.6917,
				Latitude:     35.6895,
				PostalCode:   "100-0001",
				Prefecture:   "東京都",
				City:         "千代田区",
				AddressLine1: "千代田1-1",
				AddressLine2: "",
			},
			expect: &Spot{
				TypeID:          "spot-type-id",
				UserType:        SpotUserTypeCoordinator,
				UserID:          "user-id",
				Name:            "東京タワー",
				Description:     "おすすめの観光地です。",
				ThumbnailURL:    "https://example.com/image.jpg",
				Approved:        true,
				ApprovedAdminID: "user-id",
				Longitude:       139.6917,
				Latitude:        35.6895,
				PostalCode:      "100-0001",
				Prefecture:      "東京都",
				PrefectureCode:  13,
				City:            "千代田区",
				AddressLine1:    "千代田1-1",
				AddressLine2:    "",
			},
			hasErr: false,
		},
		{
			name: "invalid prefecture",
			params: &SpotParams{
				SpotTypeID:   "spot-type-id",
				UserType:     SpotUserTypeCoordinator,
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "おすすめの観光地です。",
				ThumbnailURL: "https://example.com/image.jpg",
				Longitude:    200.0,
				Latitude:     35.6895,
				PostalCode:   "",
				Prefecture:   "",
				City:         "",
				AddressLine1: "",
				AddressLine2: "",
			},
			expect: nil,
			hasErr: true,
		},
		{
			name: "invalid longitude",
			params: &SpotParams{
				SpotTypeID:   "spot-type-id",
				UserType:     SpotUserTypeCoordinator,
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "おすすめの観光地です。",
				ThumbnailURL: "https://example.com/image.jpg",
				Longitude:    200.0,
				Latitude:     35.6895,
				PostalCode:   "100-0001",
				Prefecture:   "東京都",
				City:         "千代田区",
				AddressLine1: "千代田1-1",
				AddressLine2: "",
			},
			expect: nil,
			hasErr: true,
		},
		{
			name: "invalid latitude",
			params: &SpotParams{
				SpotTypeID:   "spot-type-id",
				UserType:     SpotUserTypeCoordinator,
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "おすすめの観光地です。",
				ThumbnailURL: "https://example.com/image.jpg",
				Longitude:    139.6917,
				Latitude:     100.0,
				PostalCode:   "100-0001",
				Prefecture:   "東京都",
				City:         "千代田区",
				AddressLine1: "千代田1-1",
				AddressLine2: "",
			},
			expect: nil,
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewSpotByAdmin(tt.params)
			if err != nil {
				assert.True(t, tt.hasErr, err)
				return
			}
			assert.False(t, tt.hasErr)

			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestSpots_TypeIDs(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name   string
		spots  Spots
		expect []string
	}{
		{
			name: "success",
			spots: Spots{
				{
					ID:           "spot-id",
					TypeID:       "spot-type-id",
					UserType:     SpotUserTypeUser,
					UserID:       "user-id",
					Name:         "東京タワー",
					Description:  "東京タワーの説明",
					ThumbnailURL: "https://example.com/thumbnail.jpg",
					Longitude:    139.732293,
					Latitude:     35.658580,
					Approved:     true,
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
			expect: []string{"spot-type-id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.spots.TypeIDs()
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestSpots_UserIDs(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name   string
		spots  Spots
		expect []string
	}{
		{
			name: "success",
			spots: Spots{
				{
					ID:           "spot-id",
					UserType:     SpotUserTypeUser,
					UserID:       "user-id",
					Name:         "東京タワー",
					Description:  "東京タワーの説明",
					ThumbnailURL: "https://example.com/thumbnail.jpg",
					Longitude:    139.732293,
					Latitude:     35.658580,
					Approved:     true,
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
			expect: []string{"user-id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.spots.UserIDs()
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestSpots_Fill(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name   string
		spots  Spots
		expect Spots
		hasErr bool
	}{
		{
			name: "success",
			spots: Spots{
				{
					TypeID:          "spot-type-id",
					UserType:        SpotUserTypeCoordinator,
					UserID:          "user-id",
					Name:            "東京タワー",
					Description:     "おすすめの観光地です。",
					ThumbnailURL:    "https://example.com/image.jpg",
					Approved:        true,
					ApprovedAdminID: "user-id",
					Longitude:       139.6917,
					Latitude:        35.6895,
					PostalCode:      "100-0001",
					Prefecture:      "",
					PrefectureCode:  13,
					City:            "千代田区",
					AddressLine1:    "千代田1-1",
					AddressLine2:    "",
					CreatedAt:       now,
					UpdatedAt:       now,
				},
			},
			expect: Spots{
				{
					TypeID:          "spot-type-id",
					UserType:        SpotUserTypeCoordinator,
					UserID:          "user-id",
					Name:            "東京タワー",
					Description:     "おすすめの観光地です。",
					ThumbnailURL:    "https://example.com/image.jpg",
					Approved:        true,
					ApprovedAdminID: "user-id",
					Longitude:       139.6917,
					Latitude:        35.6895,
					PostalCode:      "100-0001",
					Prefecture:      "東京都",
					PrefectureCode:  13,
					City:            "千代田区",
					AddressLine1:    "千代田1-1",
					AddressLine2:    "",
					CreatedAt:       now,
					UpdatedAt:       now,
				},
			},
			hasErr: false,
		},
		{
			name: "failed to fill prefecture",
			spots: Spots{
				{
					TypeID:          "spot-type-id",
					UserType:        SpotUserTypeCoordinator,
					UserID:          "user-id",
					Name:            "東京タワー",
					Description:     "おすすめの観光地です。",
					ThumbnailURL:    "https://example.com/image.jpg",
					Approved:        true,
					ApprovedAdminID: "user-id",
					Longitude:       139.6917,
					Latitude:        35.6895,
					PostalCode:      "100-0001",
					Prefecture:      "",
					PrefectureCode:  0,
					City:            "",
					AddressLine1:    "",
					AddressLine2:    "",
					CreatedAt:       now,
					UpdatedAt:       now,
				},
			},
			expect: Spots{
				{
					TypeID:          "spot-type-id",
					UserType:        SpotUserTypeCoordinator,
					UserID:          "user-id",
					Name:            "東京タワー",
					Description:     "おすすめの観光地です。",
					ThumbnailURL:    "https://example.com/image.jpg",
					Approved:        true,
					ApprovedAdminID: "user-id",
					Longitude:       139.6917,
					Latitude:        35.6895,
					PostalCode:      "100-0001",
					Prefecture:      "",
					PrefectureCode:  0,
					City:            "",
					AddressLine1:    "",
					AddressLine2:    "",
					CreatedAt:       now,
					UpdatedAt:       now,
				},
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.spots.Fill()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.spots)
		})
	}
}

func TestSpots_GroupByUserType(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name   string
		spots  Spots
		expect map[SpotUserType]Spots
	}{
		{
			name: "success",
			spots: Spots{
				{
					ID:           "spot-id",
					UserType:     1,
					UserID:       "user-id",
					Name:         "東京タワー",
					Description:  "東京タワーの説明",
					ThumbnailURL: "https://example.com/thumbnail.jpg",
					Longitude:    139.732293,
					Latitude:     35.658580,
					Approved:     true,
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
			expect: map[SpotUserType]Spots{
				SpotUserTypeUser: {
					{
						ID:           "spot-id",
						UserType:     1,
						UserID:       "user-id",
						Name:         "東京タワー",
						Description:  "東京タワーの説明",
						ThumbnailURL: "https://example.com/thumbnail.jpg",
						Longitude:    139.732293,
						Latitude:     35.658580,
						Approved:     true,
						CreatedAt:    now,
						UpdatedAt:    now,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.spots.GroupByUserType()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
