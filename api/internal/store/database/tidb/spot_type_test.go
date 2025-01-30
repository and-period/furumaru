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

func TestSpotType(t *testing.T) {
	assert.NotNil(t, NewSpotType(nil))
}

func TestSpotType_List(t *testing.T) {
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

	spotTypes := make(entity.SpotTypes, 3)
	spotTypes[0] = testSpotType("spot-type-id01", "スポット1", now())
	spotTypes[1] = testSpotType("spot-type-id02", "スポット2", now())
	spotTypes[2] = testSpotType("spot-type-id03", "スポット3", now())
	err = db.DB.Create(&spotTypes).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListSpotTypesParams
	}
	type want struct {
		spotTypes entity.SpotTypes
		err       error
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
				params: &database.ListSpotTypesParams{
					Name:   "スポット",
					Limit:  2,
					Offset: 1,
				},
			},
			want: want{
				spotTypes: spotTypes[1:],
				err:       nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &spotType{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.spotTypes, actual)
		})
	}
}

func TestSpotType_Count(t *testing.T) {
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

	spotTypes := make(entity.SpotTypes, 3)
	spotTypes[0] = testSpotType("spot-type-id01", "スポット1", now())
	spotTypes[1] = testSpotType("spot-type-id02", "スポット2", now())
	spotTypes[2] = testSpotType("spot-type-id03", "スポット3", now())
	err = db.DB.Create(&spotTypes).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListSpotTypesParams
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
				params: &database.ListSpotTypesParams{
					Name:   "スポット",
					Limit:  2,
					Offset: 1,
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

			db := &spotType{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestSpotType_MultiGet(t *testing.T) {
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

	spotTypes := make(entity.SpotTypes, 3)
	spotTypes[0] = testSpotType("spot-type-id01", "スポット1", now())
	spotTypes[1] = testSpotType("spot-type-id02", "スポット2", now())
	spotTypes[2] = testSpotType("spot-type-id03", "スポット3", now())
	err = db.DB.Create(&spotTypes).Error
	require.NoError(t, err)

	type args struct {
		spotTypeIDs []string
	}
	type want struct {
		spotTypes entity.SpotTypes
		err       error
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
				spotTypeIDs: []string{
					"spot-type-id01",
					"spot-type-id02",
					"spot-type-id03",
				},
			},
			want: want{
				spotTypes: spotTypes,
				err:       nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &spotType{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.spotTypeIDs)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.spotTypes, actual)
		})
	}
}

func TestSpotType_Get(t *testing.T) {
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

	typ := testSpotType("spot-type-id", "スポット1", now())
	err = db.DB.Create(&typ).Error
	require.NoError(t, err)

	type args struct {
		spotTypeID string
	}
	type want struct {
		spotType *entity.SpotType
		err      error
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
				spotTypeID: "spot-type-id",
			},
			want: want{
				spotType: typ,
				err:      nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				spotTypeID: "",
			},
			want: want{
				spotType: nil,
				err:      database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &spotType{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.spotTypeID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.spotType, actual)
		})
	}
}

func TestSpotType_Create(t *testing.T) {
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
		spotType *entity.SpotType
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
				spotType: testSpotType("spot-type-id", "スポット1", now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				typ := testSpotType("spot-type-id", "スポット1", now())
				err = db.DB.Create(&typ).Error
				require.NoError(t, err)
			},
			args: args{
				spotType: testSpotType("spot-type-id", "スポット1", now()),
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

			err := delete(ctx, spotTypeTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &spotType{db: db, now: now}
			err = db.Create(ctx, tt.args.spotType)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestSpotType_Update(t *testing.T) {
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
		spotTypeID string
		params     *database.UpdateSpotTypeParams
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
				spotType := testSpotType("spot-type-id", "スポット1", now())
				err = db.DB.Create(&spotType).Error
				require.NoError(t, err)
			},
			args: args{
				spotTypeID: "spot-type-id",
				params: &database.UpdateSpotTypeParams{
					Name: "スポット2",
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

			err := delete(ctx, spotTypeTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &spotType{db: db, now: now}
			err = db.Update(ctx, tt.args.spotTypeID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestSpotType_Delete(t *testing.T) {
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
		spotTypeID string
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
				spotType := testSpotType("spot-type-id", "スポット1", now())
				err = db.DB.Create(&spotType).Error
				require.NoError(t, err)
			},
			args: args{
				spotTypeID: "spot-type-id",
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

			err := delete(ctx, spotTypeTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &spotType{db: db, now: now}
			err = db.Delete(ctx, tt.args.spotTypeID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testSpotType(spotID, name string, now time.Time) *entity.SpotType {
	return &entity.SpotType{
		ID:        spotID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
