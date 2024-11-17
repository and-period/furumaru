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
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "おすすめの観光地です。",
				ThumbnailURL: "https://example.com/image.jpg",
				Longitude:    139.6917,
				Latitude:     35.6895,
			},
			expect: &Spot{
				UserType:     SpotUserTypeUser,
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "おすすめの観光地です。",
				ThumbnailURL: "https://example.com/image.jpg",
				Approved:     true,
				Longitude:    139.6917,
				Latitude:     35.6895,
			},
			hasErr: false,
		},
		{
			name: "invalid longitude",
			params: &SpotParams{
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "おすすめの観光地です。",
				ThumbnailURL: "https://example.com/image.jpg",
				Longitude:    200.0,
				Latitude:     35.6895,
			},
			expect: nil,
			hasErr: true,
		},
		{
			name: "invalid latitude",
			params: &SpotParams{
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "おすすめの観光地です。",
				ThumbnailURL: "https://example.com/image.jpg",
				Longitude:    139.6917,
				Latitude:     100.0,
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
				UserType:     SpotUserTypeCoordinator,
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "おすすめの観光地です。",
				ThumbnailURL: "https://example.com/image.jpg",
				Longitude:    139.6917,
				Latitude:     35.6895,
			},
			expect: &Spot{
				UserType:        SpotUserTypeCoordinator,
				UserID:          "user-id",
				Name:            "東京タワー",
				Description:     "おすすめの観光地です。",
				ThumbnailURL:    "https://example.com/image.jpg",
				Approved:        true,
				ApprovedAdminID: "user-id",
				Longitude:       139.6917,
				Latitude:        35.6895,
			},
			hasErr: false,
		},
		{
			name: "invalid longitude",
			params: &SpotParams{
				UserType:     SpotUserTypeCoordinator,
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "おすすめの観光地です。",
				ThumbnailURL: "https://example.com/image.jpg",
				Longitude:    200.0,
				Latitude:     35.6895,
			},
			expect: nil,
			hasErr: true,
		},
		{
			name: "invalid latitude",
			params: &SpotParams{
				UserType:     SpotUserTypeCoordinator,
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "おすすめの観光地です。",
				ThumbnailURL: "https://example.com/image.jpg",
				Longitude:    139.6917,
				Latitude:     100.0,
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
