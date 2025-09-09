package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestSpotUserType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		userType sentity.SpotUserType
		expect   SpotUserType
		response int32
	}{
		{
			name:     "user",
			userType: sentity.SpotUserTypeUser,
			expect:   SpotUserTypeUser,
			response: 1,
		},
		{
			name:     "coordinator",
			userType: sentity.SpotUserTypeCoordinator,
			expect:   SpotUserTypeCoordinator,
			response: 2,
		},
		{
			name:     "producer",
			userType: sentity.SpotUserTypeProducer,
			expect:   SpotUserTypeProducer,
			response: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewSpotUserType(tt.userType)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.response, actual.Response())
		})
	}
}

func TestSpotsByUser(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name     string
		spots    sentity.Spots
		users    map[string]*uentity.User
		expect   Spots
		response []*types.Spot
	}{
		{
			name: "success",
			spots: sentity.Spots{
				{
					ID:              "spot-id",
					UserType:        sentity.SpotUserTypeUser,
					UserID:          "user-id",
					Name:            "name",
					Description:     "description",
					ThumbnailURL:    "thumbnail",
					Longitude:       1.0,
					Latitude:        2.0,
					Approved:        true,
					ApprovedAdminID: "admin-id",
					CreatedAt:       now,
					UpdatedAt:       now,
				},
			},
			users: map[string]*uentity.User{
				"user-id": {
					ID:         "user-id",
					Registered: true,
					CreatedAt:  now,
					UpdatedAt:  now,
					Member: uentity.Member{
						Username:     "username",
						ThumbnailURL: "thumbnail",
					},
				},
			},
			expect: Spots{
				{
					Spot: types.Spot{
						ID:               "spot-id",
						Name:             "name",
						Description:      "description",
						ThumbnailURL:     "thumbnail",
						Longitude:        1.0,
						Latitude:         2.0,
						UserType:         int32(sentity.SpotUserTypeUser),
						UserID:           "user-id",
						Username:         "username",
						UserThumbnailURL: "thumbnail",
						CreatedAt:        now.Unix(),
						UpdatedAt:        now.Unix(),
					},
				},
			},
			response: []*types.Spot{
				{
					ID:               "spot-id",
					Name:             "name",
					Description:      "description",
					ThumbnailURL:     "thumbnail",
					Longitude:        1.0,
					Latitude:         2.0,
					UserType:         int32(sentity.SpotUserTypeUser),
					UserID:           "user-id",
					Username:         "username",
					UserThumbnailURL: "thumbnail",
					CreatedAt:        now.Unix(),
					UpdatedAt:        now.Unix(),
				},
			},
		},
		{
			name: "empty",
			spots: sentity.Spots{
				{
					ID:              "spot-id01",
					UserType:        sentity.SpotUserTypeUser,
					UserID:          "user-id",
					Name:            "name",
					Description:     "description",
					ThumbnailURL:    "thumbnail",
					Longitude:       1.0,
					Latitude:        2.0,
					Approved:        true,
					ApprovedAdminID: "admin-id",
					CreatedAt:       now,
					UpdatedAt:       now,
				},
				{
					ID:              "spot-id02",
					UserType:        sentity.SpotUserTypeUser,
					UserID:          "user-id",
					Name:            "name",
					Description:     "description",
					ThumbnailURL:    "thumbnail",
					Longitude:       1.0,
					Latitude:        2.0,
					Approved:        false,
					ApprovedAdminID: "",
					CreatedAt:       now,
					UpdatedAt:       now,
				},
			},
			users:    map[string]*uentity.User{},
			expect:   Spots{},
			response: []*types.Spot{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewSpotsByUser(tt.spots, tt.users)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.response, actual.Response())
		})
	}
}

func TestSpotsByCoordinator(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name     string
		spots    sentity.Spots
		users    map[string]*Coordinator
		expect   Spots
		response []*types.Spot
	}{
		{
			name: "success",
			spots: sentity.Spots{
				{
					ID:              "spot-id",
					UserType:        sentity.SpotUserTypeCoordinator,
					UserID:          "user-id",
					Name:            "name",
					Description:     "description",
					ThumbnailURL:    "thumbnail",
					Longitude:       1.0,
					Latitude:        2.0,
					Approved:        true,
					ApprovedAdminID: "admin-id",
					CreatedAt:       now,
					UpdatedAt:       now,
				},
			},
			users: map[string]*Coordinator{
				"user-id": {
					Coordinator: types.Coordinator{
						ID:           "user-id",
						Username:     "username",
						ThumbnailURL: "thumbnail",
					},
				},
			},
			expect: Spots{
				{
					Spot: types.Spot{
						ID:               "spot-id",
						Name:             "name",
						Description:      "description",
						ThumbnailURL:     "thumbnail",
						Longitude:        1.0,
						Latitude:         2.0,
						UserType:         int32(sentity.SpotUserTypeCoordinator),
						UserID:           "user-id",
						Username:         "username",
						UserThumbnailURL: "thumbnail",
						CreatedAt:        now.Unix(),
						UpdatedAt:        now.Unix(),
					},
				},
			},
			response: []*types.Spot{
				{
					ID:               "spot-id",
					Name:             "name",
					Description:      "description",
					ThumbnailURL:     "thumbnail",
					Longitude:        1.0,
					Latitude:         2.0,
					UserType:         int32(sentity.SpotUserTypeCoordinator),
					UserID:           "user-id",
					Username:         "username",
					UserThumbnailURL: "thumbnail",
					CreatedAt:        now.Unix(),
					UpdatedAt:        now.Unix(),
				},
			},
		},
		{
			name: "empty",
			spots: sentity.Spots{
				{
					ID:              "spot-id01",
					UserType:        sentity.SpotUserTypeCoordinator,
					UserID:          "user-id",
					Name:            "name",
					Description:     "description",
					ThumbnailURL:    "thumbnail",
					Longitude:       1.0,
					Latitude:        2.0,
					Approved:        true,
					ApprovedAdminID: "admin-id",
					CreatedAt:       now,
					UpdatedAt:       now,
				},
				{
					ID:              "spot-id02",
					UserType:        sentity.SpotUserTypeCoordinator,
					UserID:          "user-id",
					Name:            "name",
					Description:     "description",
					ThumbnailURL:    "thumbnail",
					Longitude:       1.0,
					Latitude:        2.0,
					Approved:        false,
					ApprovedAdminID: "",
					CreatedAt:       now,
					UpdatedAt:       now,
				},
			},
			users:    map[string]*Coordinator{},
			expect:   Spots{},
			response: []*types.Spot{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewSpotsByCoordinator(tt.spots, tt.users)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.response, actual.Response())
		})
	}
}

func TestSpotsByProducer(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name     string
		spots    sentity.Spots
		users    map[string]*Producer
		expect   Spots
		response []*types.Spot
	}{
		{
			name: "success",
			spots: sentity.Spots{
				{
					ID:              "spot-id",
					UserType:        sentity.SpotUserTypeProducer,
					UserID:          "user-id",
					Name:            "name",
					Description:     "description",
					ThumbnailURL:    "thumbnail",
					Longitude:       1.0,
					Latitude:        2.0,
					Approved:        true,
					ApprovedAdminID: "admin-id",
					CreatedAt:       now,
					UpdatedAt:       now,
				},
			},
			users: map[string]*Producer{
				"user-id": {
					Producer: types.Producer{
						ID:           "user-id",
						Username:     "username",
						ThumbnailURL: "thumbnail",
					},
				},
			},
			expect: Spots{
				{
					Spot: types.Spot{
						ID:               "spot-id",
						Name:             "name",
						Description:      "description",
						ThumbnailURL:     "thumbnail",
						Longitude:        1.0,
						Latitude:         2.0,
						UserType:         int32(sentity.SpotUserTypeProducer),
						UserID:           "user-id",
						Username:         "username",
						UserThumbnailURL: "thumbnail",
						CreatedAt:        now.Unix(),
						UpdatedAt:        now.Unix(),
					},
				},
			},
			response: []*types.Spot{
				{
					ID:               "spot-id",
					Name:             "name",
					Description:      "description",
					ThumbnailURL:     "thumbnail",
					Longitude:        1.0,
					Latitude:         2.0,
					UserType:         int32(sentity.SpotUserTypeProducer),
					UserID:           "user-id",
					Username:         "username",
					UserThumbnailURL: "thumbnail",
					CreatedAt:        now.Unix(),
					UpdatedAt:        now.Unix(),
				},
			},
		},
		{
			name: "empty",
			spots: sentity.Spots{
				{
					ID:              "spot-id01",
					UserType:        sentity.SpotUserTypeProducer,
					UserID:          "user-id",
					Name:            "name",
					Description:     "description",
					ThumbnailURL:    "thumbnail",
					Longitude:       1.0,
					Latitude:        2.0,
					Approved:        true,
					ApprovedAdminID: "admin-id",
					CreatedAt:       now,
					UpdatedAt:       now,
				},
				{
					ID:              "spot-id02",
					UserType:        sentity.SpotUserTypeProducer,
					UserID:          "user-id",
					Name:            "name",
					Description:     "description",
					ThumbnailURL:    "thumbnail",
					Longitude:       1.0,
					Latitude:        2.0,
					Approved:        false,
					ApprovedAdminID: "",
					CreatedAt:       now,
					UpdatedAt:       now,
				},
			},
			users:    map[string]*Producer{},
			expect:   Spots{},
			response: []*types.Spot{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewSpotsByProducer(tt.spots, tt.users)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.response, actual.Response())
		})
	}
}
