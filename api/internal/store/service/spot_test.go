package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListSpots(t *testing.T) {
	t.Parallel()

	now := jst.Date(2024, 11, 17, 18, 30, 0, 0)
	params := &database.ListSpotsParams{
		Name:            "東京",
		UserID:          "user-id",
		ExcludeApproved: false,
		ExcludeDisabled: false,
		Limit:           20,
		Offset:          0,
	}
	spots := entity.Spots{
		{
			ID:              "spot-id",
			UserType:        entity.SpotUserTypeUser,
			UserID:          "user-id",
			Name:            "東京タワー",
			Description:     "東京タワーの説明",
			ThumbnailURL:    "https://example.com/thumbnail.jpg",
			Longitude:       139.74545,
			Latitude:        35.65861,
			Approved:        true,
			ApprovedAdminID: "admin-id",
			CreatedAt:       now,
			UpdatedAt:       now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListSpotsInput
		expect      entity.Spots
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Spot.EXPECT().List(gomock.Any(), params).Return(spots, nil)
				mocks.db.Spot.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListSpotsInput{
				Name:            "東京",
				UserID:          "user-id",
				ExcludeApproved: false,
				ExcludeDisabled: false,
				Limit:           20,
				Offset:          0,
			},
			expect:      spots,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &store.ListSpotsInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list spots",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Spot.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Spot.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListSpotsInput{
				Name:            "東京",
				UserID:          "user-id",
				ExcludeApproved: false,
				ExcludeDisabled: false,
				Limit:           20,
				Offset:          0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count spots",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Spot.EXPECT().List(gomock.Any(), params).Return(spots, nil)
				mocks.db.Spot.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &store.ListSpotsInput{
				Name:            "東京",
				UserID:          "user-id",
				ExcludeApproved: false,
				ExcludeDisabled: false,
				Limit:           20,
				Offset:          0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			spots, total, err := service.ListSpots(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, spots)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestListSpotsByGeolocation(t *testing.T) {
	t.Parallel()

	now := jst.Date(2024, 11, 17, 18, 30, 0, 0)
	params := &database.ListSpotsByGeolocationParams{
		Longitude: 139.81083,
		Latitude:  35.71014,
		Radius:    9,
	}
	spots := entity.Spots{
		{
			ID:              "spot-id",
			UserType:        entity.SpotUserTypeUser,
			UserID:          "user-id",
			Name:            "東京タワー",
			Description:     "東京タワーの説明",
			ThumbnailURL:    "https://example.com/thumbnail.jpg",
			Longitude:       139.74545,
			Latitude:        35.65861,
			Approved:        true,
			ApprovedAdminID: "admin-id",
			CreatedAt:       now,
			UpdatedAt:       now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.ListSpotsByGeolocationInput
		expect    entity.Spots
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Spot.EXPECT().ListByGeolocation(ctx, params).Return(spots, nil)
			},
			input: &store.ListSpotsByGeolocationInput{
				Longitude: 139.81083,
				Latitude:  35.71014,
				Radius:    9,
			},
			expect:    spots,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.ListSpotsByGeolocationInput{
				Longitude: 139.81083,
				Latitude:  35.71014,
				Radius:    -1,
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list spots",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Spot.EXPECT().ListByGeolocation(ctx, params).Return(nil, assert.AnError)
			},
			input: &store.ListSpotsByGeolocationInput{
				Longitude: 139.81083,
				Latitude:  35.71014,
				Radius:    9,
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			spots, err := service.ListSpotsByGeolocation(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, spots)
		}))
	}
}

func TestGetSpot(t *testing.T) {
	t.Parallel()

	now := jst.Date(2024, 11, 17, 18, 30, 0, 0)
	spot := &entity.Spot{
		ID:              "spot-id",
		UserType:        entity.SpotUserTypeUser,
		UserID:          "user-id",
		Name:            "東京タワー",
		Description:     "東京タワーの説明",
		ThumbnailURL:    "https://example.com/thumbnail.jpg",
		Longitude:       139.74545,
		Latitude:        35.65861,
		Approved:        true,
		ApprovedAdminID: "admin-id",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetSpotInput
		expect    *entity.Spot
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Spot.EXPECT().Get(ctx, "spot-id").Return(spot, nil)
			},
			input: &store.GetSpotInput{
				SpotID: "spot-id",
			},
			expect:    spot,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetSpotInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get spot",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Spot.EXPECT().Get(ctx, "spot-id").Return(nil, assert.AnError)
			},
			input: &store.GetSpotInput{
				SpotID: "spot-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			spot, err := service.GetSpot(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, spot)
		}))
	}
}

func TestCreateSpotByUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateSpotByUserInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Spot.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, spot *entity.Spot) error {
						expect := &entity.Spot{
							ID:              spot.ID, // ignore
							UserType:        entity.SpotUserTypeUser,
							UserID:          "user-id",
							Name:            "東京タワー",
							Description:     "東京タワーの説明",
							ThumbnailURL:    "https://example.com/thumbnail.jpg",
							Longitude:       139.74545,
							Latitude:        35.65861,
							Approved:        true,
							ApprovedAdminID: "",
						}
						assert.Equal(t, expect, spot)
						return nil
					})
			},
			input: &store.CreateSpotByUserInput{
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "東京タワーの説明",
				ThumbnailURL: "https://example.com/thumbnail.jpg",
				Longitude:    139.74545,
				Latitude:     35.65861,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateSpotByUserInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create spot",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Spot.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateSpotByUserInput{
				UserID:       "user-id",
				Name:         "東京タワー",
				Description:  "東京タワーの説明",
				ThumbnailURL: "https://example.com/thumbnail.jpg",
				Longitude:    139.74545,
				Latitude:     35.65861,
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateSpotByUser(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestCreateSpotByAdmin(t *testing.T) {
	t.Parallel()

	adminIn := &user.GetAdminInput{
		AdminID: "admin-id",
	}
	admin := func(role uentity.AdminRole) *uentity.Admin {
		return &uentity.Admin{
			ID:   "admin-id",
			Role: role,
		}
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateSpotByAdminInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(ctx, adminIn).Return(admin(uentity.AdminRoleCoordinator), nil)
				mocks.db.Spot.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, spot *entity.Spot) error {
						expect := &entity.Spot{
							ID:              spot.ID, // ignore
							UserType:        entity.SpotUserTypeCoordinator,
							UserID:          "admin-id",
							Name:            "東京タワー",
							Description:     "東京タワーの説明",
							ThumbnailURL:    "https://example.com/thumbnail.jpg",
							Longitude:       139.74545,
							Latitude:        35.65861,
							Approved:        true,
							ApprovedAdminID: "admin-id",
						}
						assert.Equal(t, expect, spot)
						return nil
					})
			},
			input: &store.CreateSpotByAdminInput{
				AdminID:      "admin-id",
				Name:         "東京タワー",
				Description:  "東京タワーの説明",
				ThumbnailURL: "https://example.com/thumbnail.jpg",
				Longitude:    139.74545,
				Latitude:     35.65861,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateSpotByAdminInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "not found admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(ctx, adminIn).Return(nil, exception.ErrNotFound)
			},
			input: &store.CreateSpotByAdminInput{
				AdminID:      "admin-id",
				Name:         "東京タワー",
				Description:  "東京タワーの説明",
				ThumbnailURL: "https://example.com/thumbnail.jpg",
				Longitude:    139.74545,
				Latitude:     35.65861,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(ctx, adminIn).Return(nil, assert.AnError)
			},
			input: &store.CreateSpotByAdminInput{
				AdminID:      "admin-id",
				Name:         "東京タワー",
				Description:  "東京タワーの説明",
				ThumbnailURL: "https://example.com/thumbnail.jpg",
				Longitude:    139.74545,
				Latitude:     35.65861,
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to unsppoted",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(ctx, adminIn).Return(admin(uentity.AdminRoleAdministrator), nil)
			},
			input: &store.CreateSpotByAdminInput{
				AdminID:      "admin-id",
				Name:         "東京タワー",
				Description:  "東京タワーの説明",
				ThumbnailURL: "https://example.com/thumbnail.jpg",
				Longitude:    139.74545,
				Latitude:     35.65861,
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to create spot",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(ctx, adminIn).Return(admin(uentity.AdminRoleCoordinator), nil)
				mocks.db.Spot.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateSpotByAdminInput{
				AdminID:      "admin-id",
				Name:         "東京タワー",
				Description:  "東京タワーの説明",
				ThumbnailURL: "https://example.com/thumbnail.jpg",
				Longitude:    139.74545,
				Latitude:     35.65861,
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateSpotByAdmin(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateSpot(t *testing.T) {
	t.Parallel()

	params := &database.UpdateSpotParams{
		Name:         "東京スカイツリー",
		Description:  "東京スカイツリーの説明",
		ThumbnailURL: "https://example.com/thumbnail.jpg",
		Longitude:    139.81083,
		Latitude:     35.71014,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateSpotInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Spot.EXPECT().Update(ctx, "spot-id", params).Return(nil)
			},
			input: &store.UpdateSpotInput{
				SpotID:       "spot-id",
				Name:         "東京スカイツリー",
				Description:  "東京スカイツリーの説明",
				ThumbnailURL: "https://example.com/thumbnail.jpg",
				Longitude:    139.81083,
				Latitude:     35.71014,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateSpotInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update spot",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Spot.EXPECT().Update(ctx, "spot-id", params).Return(assert.AnError)
			},
			input: &store.UpdateSpotInput{
				SpotID:       "spot-id",
				Name:         "東京スカイツリー",
				Description:  "東京スカイツリーの説明",
				ThumbnailURL: "https://example.com/thumbnail.jpg",
				Longitude:    139.81083,
				Latitude:     35.71014,
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateSpot(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteSpot(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteSpotInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Spot.EXPECT().Delete(ctx, "spot-id").Return(nil)
			},
			input: &store.DeleteSpotInput{
				SpotID: "spot-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteSpotInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete spot",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Spot.EXPECT().Delete(ctx, "spot-id").Return(assert.AnError)
			},
			input: &store.DeleteSpotInput{
				SpotID: "spot-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteSpot(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestApproveSpot(t *testing.T) {
	t.Parallel()

	params := &database.ApproveSpotParams{
		Approved:        true,
		ApprovedAdminID: "admin-id",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.ApproveSpotInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Spot.EXPECT().Approve(ctx, "spot-id", params).Return(nil)
			},
			input: &store.ApproveSpotInput{
				SpotID:   "spot-id",
				AdminID:  "admin-id",
				Approved: true,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.ApproveSpotInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to Approve spot",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Spot.EXPECT().Approve(ctx, "spot-id", params).Return(assert.AnError)
			},
			input: &store.ApproveSpotInput{
				SpotID:   "spot-id",
				AdminID:  "admin-id",
				Approved: true,
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ApproveSpot(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
