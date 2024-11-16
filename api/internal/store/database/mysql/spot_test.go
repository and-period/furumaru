package mysql

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
	assert.NotNil(t, newSpot(nil))
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

	spots := make(entity.Spots, 3)
	spots[0] = testSpot("spot-id01", 35.658581, 139.745433, now())
	spots[1] = testSpot("spot-id02", 35.658581, 139.745433, now())
	spots[2] = testSpot("spot-id03", 35.658581, 139.745433, now())
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

	spots := make(entity.Spots, 3)
	spots[0] = testSpot("spot-id01", 35.658581, 139.745433, now())
	spots[1] = testSpot("spot-id02", 35.658581, 139.745433, now())
	spots[2] = testSpot("spot-id03", 35.658581, 139.745433, now())
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
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListSpotsByGeolocationParams{
					Longitude: 139.745433 + 0.018,
					Latitude:  35.658581 + 0.018,
					Radius:    2,
				},
			},
			want: want{
				spots: spots,
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

	spots := make(entity.Spots, 3)
	spots[0] = testSpot("spot-id01", 35.658581, 139.745433, now())
	spots[1] = testSpot("spot-id02", 35.658581, 139.745433, now())
	spots[2] = testSpot("spot-id03", 35.658581, 139.745433, now())
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

	s := testSpot("spot-id", 35.658581, 139.745433, now())
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
				spot: testSpot("spot-id", 35.658581, 139.745433, now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				spot := testSpot("spot-id", 35.658581, 139.745433, now())
				err := db.DB.Create(&spot).Error
				require.NoError(t, err)
			},
			args: args{
				spot: testSpot("spot-id", 35.658581, 139.745433, now()),
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
				spot := testSpot("spot-id", 35.658581, 139.745433, now())
				err := db.DB.Create(&spot).Error
				require.NoError(t, err)
			},
			args: args{
				spotID: "spot-id",
				params: &database.UpdateSpotParams{
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
				spot := testSpot("spot-id", 35.658581, 139.745433, now())
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
				spot := testSpot("spot-id", 35.658581, 139.745433, now())
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

func testSpot(spotID string, latitude, longitude float64, now time.Time) *entity.Spot {
	return &entity.Spot{
		ID:           spotID,
		UserType:     entity.SpotUserTypeUser,
		UserID:       "user-id",
		Name:         "東京タワー",
		Description:  "東京タワーの説明",
		ThumbnailURL: "http://example.com/thumbnail.jpg",
		Longitude:    longitude,
		Latitude:     latitude,
		Approved:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
