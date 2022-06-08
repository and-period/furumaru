package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
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

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, categoryTable)
	categories := make(entity.Categories, 3)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	categories[2] = testCategory("category-id03", "水産物", now())
	err = m.db.DB.Create(&categories).Error
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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, m)

			db := &category{db: m.db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreCategoriesField(actual, now())
			assert.Equal(t, tt.want.categories, actual)
		})
	}
}

func TestCategory_Create(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, categoryTable)

	type args struct {
		category *entity.Category
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				category: testCategory("category-id", "野菜", now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				category := testCategory("category-id", "野菜", now())
				err = m.db.DB.Create(&category).Error
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
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := m.dbDelete(ctx, categoryTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &category{db: m.db, now: now}
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

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, categoryTable)

	type args struct {
		categoryID string
		name       string
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				category := testCategory("category-id", "野菜", now())
				err = m.db.DB.Create(&category).Error
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
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := m.dbDelete(ctx, categoryTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &category{db: m.db, now: now}
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

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, categoryTable)

	type args struct {
		categoryID string
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				category := testCategory("category-id", "野菜", now())
				err = m.db.DB.Create(&category).Error
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
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := m.dbDelete(ctx, categoryTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &category{db: m.db, now: now}
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

func fillIgnoreCategoryField(c *entity.Category, now time.Time) {
	if c == nil {
		return
	}
	c.CreatedAt = now
	c.UpdatedAt = now
}

func fillIgnoreCategoriesField(cs entity.Categories, now time.Time) {
	for i := range cs {
		fillIgnoreCategoryField(cs[i], now)
	}
}
