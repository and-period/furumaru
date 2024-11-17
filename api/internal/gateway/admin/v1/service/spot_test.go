package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/stretchr/testify/assert"
)

func TestSpotUserType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		userType entity.SpotUserType
		expect   SpotUserType
		response int32
	}{
		{
			name:     "user",
			userType: entity.SpotUserTypeUser,
			expect:   SpotUserTypeUser,
			response: 1,
		},
		{
			name:     "admin",
			userType: entity.SpotUserTypeAdmin,
			expect:   SpotUserTypeAdmin,
			response: 2,
		},
		{
			name:     "unknown",
			userType: entity.SpotUserTypeUnknown,
			expect:   SpotUserTypeUnknown,
			response: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual1 := NewSpotUserType(tt.userType)
			assert.Equal(t, tt.expect, actual1)
			assert.Equal(t, tt.response, actual1.Response())

			actual2 := NewSpotUserTypeFromInt32(tt.response)
			assert.Equal(t, actual1, actual2)
		})
	}
}

func TestSpots(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		spots    entity.Spots
		expect   Spots
		response []*response.Spot
	}{
		{
			name: "success",
			spots: entity.Spots{
				{
					ID:              "spot-id",
					UserType:        entity.SpotUserTypeUser,
					UserID:          "user-id",
					Name:            "東京タワー",
					Description:     "東京タワーの説明",
					ThumbnailURL:    "https://example.com/thumbnail.jpg",
					Longitude:       139.732293,
					Latitude:        35.658580,
					Approved:        true,
					ApprovedAdminID: "admin-id",
					CreatedAt:       time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:       time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expect: Spots{
				{
					Spot: response.Spot{
						ID:           "spot-id",
						UserType:     1,
						UserID:       "user-id",
						Name:         "東京タワー",
						Description:  "東京タワーの説明",
						ThumbnailURL: "https://example.com/thumbnail.jpg",
						Longitude:    139.732293,
						Latitude:     35.658580,
						Approved:     true,
						CreatedAt:    1609459200,
						UpdatedAt:    1609459200,
					},
					userType: SpotUserTypeUser,
				},
			},
			response: []*response.Spot{
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
					CreatedAt:    1609459200,
					UpdatedAt:    1609459200,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewSpots(tt.spots)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.response, actual.Response())
		})
	}
}

func TestSpots_UserIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		spots  Spots
		expect []string
	}{
		{
			name: "success",
			spots: Spots{
				{
					Spot: response.Spot{
						ID:           "spot-id",
						UserType:     1,
						UserID:       "user-id",
						Name:         "東京タワー",
						Description:  "東京タワーの説明",
						ThumbnailURL: "https://example.com/thumbnail.jpg",
						Longitude:    139.732293,
						Latitude:     35.658580,
						Approved:     true,
						CreatedAt:    1609459200,
						UpdatedAt:    1609459200,
					},
					userType: SpotUserTypeUser,
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

	tests := []struct {
		name   string
		spots  Spots
		expect map[SpotUserType]Spots
	}{
		{
			name: "success",
			spots: Spots{
				{
					Spot: response.Spot{
						ID:           "spot-id",
						UserType:     1,
						UserID:       "user-id",
						Name:         "東京タワー",
						Description:  "東京タワーの説明",
						ThumbnailURL: "https://example.com/thumbnail.jpg",
						Longitude:    139.732293,
						Latitude:     35.658580,
						Approved:     true,
						CreatedAt:    1609459200,
						UpdatedAt:    1609459200,
					},
					userType: SpotUserTypeUser,
				},
			},
			expect: map[SpotUserType]Spots{
				SpotUserTypeUser: {
					{
						Spot: response.Spot{
							ID:           "spot-id",
							UserType:     1,
							UserID:       "user-id",
							Name:         "東京タワー",
							Description:  "東京タワーの説明",
							ThumbnailURL: "https://example.com/thumbnail.jpg",
							Longitude:    139.732293,
							Latitude:     35.658580,
							Approved:     true,
							CreatedAt:    1609459200,
							UpdatedAt:    1609459200,
						},
						userType: SpotUserTypeUser,
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
