package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategory(t *testing.T) {
	assert.NotNil(t, NewCategory(nil))
}

func TestCategory_List(t *testing.T) {
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

	categories := make(entity.Categories, 3)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	categories[2] = testCategory("category-id03", "水産物", now())
	err = db.DB.Create(&categories).Error
	require.NoError(t, err)

	type args struct {
		params *ListCategoriesParams
	}
	type want struct {
		categories entity.Categories
		hasErr     bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &ListCategoriesParams{
					Name:   "物",
					Limit:  1,
					Offset: 1,
				},
			},
			want: want{
				categories: categories[2:],
				hasErr:     false,
			},
		},
		{
			name:  "success with sort asc",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &ListCategoriesParams{
					Orders: []*ListCategoriesOrder{
						{Key: entity.CategoryOrderByName, OrderByASC: true},
					},
				},
			},
			want: want{
				categories: categories,
				hasErr:     false,
			},
		},
		{
			name:  "success with sort desc",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &ListCategoriesParams{
					Orders: []*ListCategoriesOrder{
						{Key: entity.CategoryOrderByName, OrderByASC: false},
					},
				},
			},
			want: want{
				categories: categories,
				hasErr:     false,
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

			db := &category{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.categories, actual)
		})
	}
}

func TestCategory_Count(t *testing.T) {
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

	categories := make(entity.Categories, 3)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	categories[2] = testCategory("category-id03", "水産物", now())
	err = db.DB.Create(&categories).Error
	require.NoError(t, err)

	type args struct {
		params *ListCategoriesParams
	}
	type want struct {
		total  int64
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &ListCategoriesParams{
					Name: "物",
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

			db := &category{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestCategory_MultiGet(t *testing.T) {
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

	categories := make(entity.Categories, 3)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	categories[2] = testCategory("category-id03", "水産物", now())
	err = db.DB.Create(&categories).Error
	require.NoError(t, err)

	type args struct {
		categoryIDs []string
	}
	type want struct {
		categories entity.Categories
		hasErr     bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				categoryIDs: []string{"category-id01", "category-id02"},
			},
			want: want{
				categories: categories[:2],
				hasErr:     false,
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

			db := &category{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.categoryIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.categories, actual)
		})
	}
}

func TestCategory_Get(t *testing.T) {
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

	c := testCategory("category-id", "野菜", now())
	err = db.DB.Create(&c).Error
	require.NoError(t, err)

	type args struct {
		categoryID string
	}
	type want struct {
		category *entity.Category
		hasErr   bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				categoryID: "category-id",
			},
			want: want{
				category: c,
				hasErr:   false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				categoryID: "other-id",
			},
			want: want{
				category: nil,
				hasErr:   true,
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

			db := &category{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.categoryID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.category, actual)
		})
	}
}

func TestCategory_Create(t *testing.T) {
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
		category *entity.Category
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				category: testCategory("category-id", "野菜", now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				category := testCategory("category-id", "野菜", now())
				err = db.DB.Create(&category).Error
				require.NoError(t, err)
			},
			args: args{
				category: testCategory("category-id", "野菜", now()),
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

			err := delete(ctx, categoryTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &category{db: db, now: now}
			err = db.Create(ctx, tt.args.category)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestCategory_Update(t *testing.T) {
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
		categoryID string
		name       string
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				category := testCategory("category-id", "野菜", now())
				err = db.DB.Create(&category).Error
				require.NoError(t, err)
			},
			args: args{
				categoryID: "category-id",
				name:       "魚介類",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				categoryID: "category-id",
				name:       "魚介類",
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

			err := delete(ctx, categoryTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &category{db: db, now: now}
			err = db.Update(ctx, tt.args.categoryID, tt.args.name)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestCategory_Delete(t *testing.T) {
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
		categoryID string
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				category := testCategory("category-id", "野菜", now())
				err = db.DB.Create(&category).Error
				require.NoError(t, err)
			},
			args: args{
				categoryID: "category-id",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				categoryID: "category-id",
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

			err := delete(ctx, categoryTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &category{db: db, now: now}
			err = db.Delete(ctx, tt.args.categoryID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testCategory(id, name string, now time.Time) *entity.Category {
	return &entity.Category{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
