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

func TestProductTag(t *testing.T) {
	assert.NotNil(t, newProductTag(nil))
}

func TestProductTag_List(t *testing.T) {
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

	productTags := make(entity.ProductTags, 3)
	productTags[0] = testProductTag("tag-id01", "野菜", now())
	productTags[1] = testProductTag("tag-id02", "水産加工物", now())
	productTags[2] = testProductTag("tag-id03", "水産物", now())
	err = db.DB.Create(&productTags).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListProductTagsParams
	}
	type want struct {
		productTags entity.ProductTags
		hasErr      bool
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
				params: &database.ListProductTagsParams{
					Name:   "水産",
					Limit:  1,
					Offset: 1,
				},
			},
			want: want{
				productTags: productTags[2:],
				hasErr:      false,
			},
		},
		{
			name:  "success with sort asc",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListProductTagsParams{
					Orders: []*database.ListProductTagsOrder{
						{Key: entity.ProductTagOrderByName, OrderByASC: true},
					},
				},
			},
			want: want{
				productTags: productTags,
				hasErr:      false,
			},
		},
		{
			name:  "success with sort desc",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListProductTagsParams{
					Orders: []*database.ListProductTagsOrder{
						{Key: entity.ProductTagOrderByName, OrderByASC: false},
					},
				},
			},
			want: want{
				productTags: productTags,
				hasErr:      false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &productTag{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.productTags, actual)
		})
	}
}

func TestProductTag_Count(t *testing.T) {
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

	productTags := make(entity.ProductTags, 3)
	productTags[0] = testProductTag("tag-id01", "野菜", now())
	productTags[1] = testProductTag("tag-id02", "水産加工物", now())
	productTags[2] = testProductTag("tag-id03", "水産物", now())
	err = db.DB.Create(&productTags).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListProductTagsParams
	}
	type want struct {
		total  int64
		hasErr bool
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
				params: &database.ListProductTagsParams{
					Name: "水産",
				},
			},
			want: want{
				total:  2,
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &productTag{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestProductTag_MultiGet(t *testing.T) {
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

	productTags := make(entity.ProductTags, 3)
	productTags[0] = testProductTag("tag-id01", "野菜", now())
	productTags[1] = testProductTag("tag-id02", "果物", now())
	productTags[2] = testProductTag("tag-id03", "水産物", now())
	err = db.DB.Create(&productTags).Error
	require.NoError(t, err)

	type args struct {
		productTagIDs []string
	}
	type want struct {
		productTags entity.ProductTags
		hasErr      bool
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
				productTagIDs: []string{"tag-id01", "tag-id02"},
			},
			want: want{
				productTags: productTags[:2],
				hasErr:      false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &productTag{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.productTagIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.productTags, actual)
		})
	}
}

func TestProductTag_Get(t *testing.T) {
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

	c := testProductTag("tag-id", "野菜", now())
	err = db.DB.Create(&c).Error
	require.NoError(t, err)

	type args struct {
		productTagID string
	}
	type want struct {
		productTag *entity.ProductTag
		hasErr     bool
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
				productTagID: "tag-id",
			},
			want: want{
				productTag: c,
				hasErr:     false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				productTagID: "other-id",
			},
			want: want{
				productTag: nil,
				hasErr:     true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &productTag{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.productTagID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.productTag, actual)
		})
	}
}

func TestProductTag_Create(t *testing.T) {
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
		productTag *entity.ProductTag
	}
	type want struct {
		hasErr bool
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
				productTag: testProductTag("tag-id", "野菜", now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				productTag := testProductTag("tag-id", "野菜", now())
				err = db.DB.Create(&productTag).Error
				require.NoError(t, err)
			},
			args: args{
				productTag: testProductTag("tag-id", "野菜", now()),
			},
			want: want{
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, productTagTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &productTag{db: db, now: now}
			err = db.Create(ctx, tt.args.productTag)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestProductTag_Update(t *testing.T) {
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
		productTagID string
		name         string
	}
	type want struct {
		hasErr bool
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
				productTag := testProductTag("tag-id", "野菜", now())
				err = db.DB.Create(&productTag).Error
				require.NoError(t, err)
			},
			args: args{
				productTagID: "tag-id",
				name:         "魚介類",
			},
			want: want{
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, productTagTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &productTag{db: db, now: now}
			err = db.Update(ctx, tt.args.productTagID, tt.args.name)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestProductTag_Delete(t *testing.T) {
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
		productTagID string
	}
	type want struct {
		hasErr bool
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
				productTag := testProductTag("tag-id", "野菜", now())
				err = db.DB.Create(&productTag).Error
				require.NoError(t, err)
			},
			args: args{
				productTagID: "tag-id",
			},
			want: want{
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, productTagTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &productTag{db: db, now: now}
			err = db.Delete(ctx, tt.args.productTagID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testProductTag(id, name string, now time.Time) *entity.ProductTag {
	return &entity.ProductTag{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
