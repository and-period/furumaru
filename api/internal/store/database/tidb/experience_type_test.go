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

func TestExperienceType(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, NewExperienceType(nil))
}

func TestExperienceType_List(t *testing.T) {
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

	types := make(entity.ExperienceTypes, 3)
	types[0] = testExperienceType("experience-type-id01", "じゃがいも収穫", now())
	types[1] = testExperienceType("experience-type-id02", "トマト収穫", now())
	types[2] = testExperienceType("experience-type-id03", "キャベツ収穫", now())
	err = db.DB.Create(&types).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListExperienceTypesParams
	}
	type want struct {
		experienceTypes entity.ExperienceTypes
		err             error
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
				params: &database.ListExperienceTypesParams{
					Limit:  1,
					Offset: 1,
				},
			},
			want: want{
				experienceTypes: types[1:2],
				err:             nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &experienceType{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.ElementsMatch(t, tt.want.experienceTypes, actual)
		})
	}
}

func TestExperienceType_Count(t *testing.T) {
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

	types := make(entity.ExperienceTypes, 3)
	types[0] = testExperienceType("experience-type-id01", "じゃがいも収穫", now())
	types[1] = testExperienceType("experience-type-id02", "トマト収穫", now())
	types[2] = testExperienceType("experience-type-id03", "キャベツ収穫", now())
	err = db.DB.Create(&types).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListExperienceTypesParams
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
				params: &database.ListExperienceTypesParams{
					Limit:  1,
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

			db := &experienceType{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestExperienceType_MultiGet(t *testing.T) {
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

	types := make(entity.ExperienceTypes, 3)
	types[0] = testExperienceType("experience-type-id01", "じゃがいも収穫", now())
	types[1] = testExperienceType("experience-type-id02", "トマト収穫", now())
	types[2] = testExperienceType("experience-type-id03", "キャベツ収穫", now())
	err = db.DB.Create(&types).Error
	require.NoError(t, err)

	type args struct {
		experienceTypeIDs []string
	}
	type want struct {
		experienceTypes entity.ExperienceTypes
		err             error
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
				experienceTypeIDs: []string{"experience-type-id01", "experience-type-id02"},
			},
			want: want{
				experienceTypes: types[:2],
				err:             nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &experienceType{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.experienceTypeIDs)
			assert.ErrorIs(t, err, tt.want.err)
			assert.ElementsMatch(t, tt.want.experienceTypes, actual)
		})
	}
}

func TestExperienceType_Get(t *testing.T) {
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

	typ := testExperienceType("experience-type-id", "じゃがいも収穫", now())
	err = db.DB.Create(&typ).Error
	require.NoError(t, err)

	type args struct {
		experienceTypeID string
	}
	type want struct {
		experienceType *entity.ExperienceType
		err            error
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
				experienceTypeID: "experience-type-id",
			},
			want: want{
				experienceType: typ,
				err:            nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				experienceTypeID: "",
			},
			want: want{
				experienceType: nil,
				err:            database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &experienceType{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.experienceTypeID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.experienceType, actual)
		})
	}
}

func TestExperienceType_Create(t *testing.T) {
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
		experienceType *entity.ExperienceType
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
				experienceType: testExperienceType("experience-type-id", "じゃがいも収穫", now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				typ := testExperienceType("experience-type-id", "じゃがいも収穫", now())
				err := db.DB.Create(&typ).Error
				require.NoError(t, err)
			},
			args: args{
				experienceType: testExperienceType("experience-type-id", "じゃがいも収穫", now()),
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

			err := delete(ctx, experienceTypeTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &experienceType{db: db, now: now}
			err = db.Create(ctx, tt.args.experienceType)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestExperienceType_Update(t *testing.T) {
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
		experienceTypeID string
		params           *database.UpdateExperienceTypeParams
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
				typ := testExperienceType("experience-type-id", "じゃがいも収穫", now())
				err := db.DB.Create(&typ).Error
				require.NoError(t, err)
			},
			args: args{
				experienceTypeID: "experience-type-id",
				params: &database.UpdateExperienceTypeParams{
					Name: "じゃがいも収穫",
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

			err := delete(ctx, experienceTypeTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &experienceType{db: db, now: now}
			err = db.Update(ctx, tt.args.experienceTypeID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestExperienceType_Delete(t *testing.T) {
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
		experienceTypeID string
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
				typ := testExperienceType("experience-type-id", "じゃがいも収穫", now())
				err := db.DB.Create(&typ).Error
				require.NoError(t, err)
			},
			args: args{
				experienceTypeID: "experience-type-id",
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

			err := delete(ctx, experienceTypeTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &experienceType{db: db, now: now}
			err = db.Delete(ctx, tt.args.experienceTypeID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testExperienceType(id, name string, now time.Time) *entity.ExperienceType {
	return &entity.ExperienceType{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
