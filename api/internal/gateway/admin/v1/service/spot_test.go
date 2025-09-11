package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/stretchr/testify/assert"
)

func TestSpotUserType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		userType entity.SpotUserType
		expect   types.SpotUserType
	}{
		{
			name:     "user",
			userType: entity.SpotUserTypeUser,
			expect:   types.SpotUserTypeUser,
		},
		{
			name:     "coordinator",
			userType: entity.SpotUserTypeCoordinator,
			expect:   types.SpotUserTypeCoordinator,
		},
		{
			name:     "producer",
			userType: entity.SpotUserTypeProducer,
			expect:   types.SpotUserTypeProducer,
		},
		{
			name:     "unknown",
			userType: entity.SpotUserTypeUnknown,
			expect:   types.SpotUserTypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual1 := NewSpotUserType(tt.userType)
			assert.Equal(t, tt.expect, actual1.Response())
		})
	}
}

func TestSpots(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		spots    entity.Spots
		expect   Spots
		response []*types.Spot
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
					Spot: types.Spot{
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
					UserType: SpotUserType(types.SpotUserTypeUser),
				},
			},
			response: []*types.Spot{
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
