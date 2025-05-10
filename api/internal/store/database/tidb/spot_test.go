package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpot(t *testing.T) {
	assert.NotNil(t, NewSpot(nil))
}

func TestSpot_List(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	spotType := testSpotType("spot-type-id", "観光地", now())
	err = db.DB.Create(&spotType).Error
	require.NoError(t, err)

	spots := make(entity.Spots, 3)
	spots[0] = testSpot("spot-id01", "spot-type-id", 35.658581, 139.745433, now())
	spots[1] = testSpot("spot-id02", "spot-type-id", 35.658581, 139.745433, now())
	spots[2] = testSpot("spot-id03", "spot-type-id", 35.658581, 139.745433, now())
	err = db.DB.Create(&spots).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListSpotsParams
	}
	type want struct {
		spots entity.Spots
		err   error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListSpotsParams{
					Name:            "東京",
					UserID:          "user-id",
					ExcludeApproved: false,
					ExcludeDisabled: false,
					Limit:           2,
					Offset:          1,
				},
			},
			want: want{
				spots: spots[1:],
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &spot{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.spots, actual)
		})
	}
}

func TestSpot_ListByGeolocation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	spotType := testSpotType("spot-type-id", "観光地", now())
	err = db.DB.Create(&spotType).Error
	require.NoError(t, err)

	spots := make(entity.Spots, 3)
	spots[0] = testSpot("spot-id01", "spot-type-id", 35.65861, 139.74545, now())
	spots[1] = testSpot("spot-id02", "spot-type-id", 0, 0, now())
	spots[2] = testSpot("spot-id03", "spot-type-id", 0, 0, now())
	err = db.DB.Create(&spots).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListSpotsByGeolocationParams
	}
	type want struct {
		spots entity.Spots
		err   error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "inside",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListSpotsByGeolocationParams{
					Longitude: 139.81083,
					Latitude:  35.71014,
					Radius:    9,
				},
			},
			want: want{
				spots: spots[:1],
				err:   nil,
			},
		},
		{
			name:  "inside longitude",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListSpotsByGeolocationParams{
					Longitude: 0.1,
					Latitude:  0.0,
					Radius:    12,
				},
			},
			want: want{
				spots: spots[1:],
				err:   nil,
			},
		},
		{
			name:  "inside latitude",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListSpotsByGeolocationParams{
					Longitude: 0.0,
					Latitude:  0.1,
					Radius:    12,
				},
			},
			want: want{
				spots: spots[1:],
				err:   nil,
			},
		},
		{
			name:  "inside",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListSpotsByGeolocationParams{
					Longitude: 139.81083,
					Latitude:  35.71014,
					Radius:    8,
				},
			},
			want: want{
				spots: entity.Spots{},
				err:   nil,
			},
		},
		{
			name:  "outside longitude",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListSpotsByGeolocationParams{
					Longitude: 0.1,
					Latitude:  0.0,
					Radius:    11,
				},
			},
			want: want{
				spots: entity.Spots{},
				err:   nil,
			},
		},
		{
			name:  "outside latitude",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListSpotsByGeolocationParams{
					Longitude: 0.0,
					Latitude:  0.1,
					Radius:    11,
				},
			},
			want: want{
				spots: entity.Spots{},
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &spot{db: db, now: now}
			actual, err := db.ListByGeolocation(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.spots, actual)
		})
	}
}

func TestSpot_Count(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	spotType := testSpotType("spot-type-id", "観光地", now())
	err = db.DB.Create(&spotType).Error
	require.NoError(t, err)

	spots := make(entity.Spots, 3)
	spots[0] = testSpot("spot-id01", "spot-type-id", 35.658581, 139.745433, now())
	spots[1] = testSpot("spot-id02", "spot-type-id", 35.658581, 139.745433, now())
	spots[2] = testSpot("spot-id03", "spot-type-id", 35.658581, 139.745433, now())
	err = db.DB.Create(&spots).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListSpotsParams
	}
	type want struct {
		total int64
		err   error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListSpotsParams{
					Name:            "東京",
					UserID:          "user-id",
					ExcludeApproved: false,
					ExcludeDisabled: false,
					Limit:           2,
					Offset:          1,
				},
			},
			want: want{
				total: 3,
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &spot{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestSpot_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	spotType := testSpotType("spot-type-id", "観光地", now())
	err = db.DB.Create(&spotType).Error
	require.NoError(t, err)

	s := testSpot("spot-id", "spot-type-id", 35.658581, 139.745433, now())
	err = db.DB.Create(&s).Error
	require.NoError(t, err)

	type args struct {
		spotID string
	}
	type want struct {
		spot *entity.Spot
		err  error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				spotID: "spot-id",
			},
			want: want{
				spot: s,
				err:  nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				spotID: "",
			},
			want: want{
				spot: nil,
				err:  database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &spot{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.spotID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.spot, actual)
		})
	}
}

func TestSpot_Create(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	spotType := testSpotType("spot-type-id", "観光地", now())
	err = db.DB.Create(&spotType).Error
	require.NoError(t, err)
	type args struct {
		spot *entity.Spot
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				spot: testSpot("spot-id", "spot-type-id", 35.658581, 139.745433, now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				spot := testSpot("spot-id", "spot-type-id", 35.658581, 139.745433, now())
				err := db.DB.Create(&spot).Error
				require.NoError(t, err)
			},
			args: args{
				spot: testSpot("spot-id", "spot-type-id", 35.658581, 139.745433, now()),
			},
			want: want{
				err: database.ErrAlreadyExists,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, spotTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &spot{db: db, now: now}
			err = db.Create(ctx, tt.args.spot)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestSpot_Update(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	spotType := testSpotType("spot-type-id", "観光地", now())
	err = db.DB.Create(&spotType).Error
	require.NoError(t, err)
	type args struct {
		spotID string
		params *database.UpdateSpotParams
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				spot := testSpot("spot-id", "spot-type-id", 35.658581, 139.745433, now())
				err := db.DB.Create(&spot).Error
				require.NoError(t, err)
			},
			args: args{
				spotID: "spot-id",
				params: &database.UpdateSpotParams{
					SpotTypeID:   "spot-type-id",
					Name:         "東京スカイツリー",
					Description:  "東京スカイツリーの説明",
					ThumbnailURL: "http://example.com/thumbnail.jpg",
					Longitude:    139.8107,
					Latitude:     35.7101,
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, spotTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &spot{db: db, now: now}
			err = db.Update(ctx, tt.args.spotID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestSpot_Delete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	type args struct {
		spotID string
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				spot := testSpot("spot-id", "", 35.658581, 139.745433, now())
				err := db.DB.Create(&spot).Error
				require.NoError(t, err)
			},
			args: args{
				spotID: "spot-id",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, spotTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &spot{db: db, now: now}
			err = db.Delete(ctx, tt.args.spotID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestSpot_Approve(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	type args struct {
		spotID string
		params *database.ApproveSpotParams
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				spot := testSpot("spot-id", "", 35.658581, 139.745433, now())
				err := db.DB.Create(&spot).Error
				require.NoError(t, err)
			},
			args: args{
				spotID: "spot-id",
				params: &database.ApproveSpotParams{
					Approved:        true,
					ApprovedAdminID: "admin-id",
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, spotTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &spot{db: db, now: now}
			err = db.Approve(ctx, tt.args.spotID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testSpot(spotID, typeID string, latitude, longitude float64, now time.Time) *entity.Spot {
	return &entity.Spot{
		ID:             spotID,
		TypeID:         typeID,
		UserType:       entity.SpotUserTypeUser,
		UserID:         "user-id",
		Name:           "東京タワー",
		Description:    "東京タワーの説明",
		ThumbnailURL:   "http://example.com/thumbnail.jpg",
		Longitude:      longitude,
		Latitude:       latitude,
		PostalCode:     "100-0001",
		Prefecture:     "東京都",
		PrefectureCode: 13,
		City:           "港区",
		AddressLine1:   "芝公園4-2-8",
		AddressLine2:   "",
		Approved:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
