package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContactCategory(t *testing.T) {
	assert.NotNil(t, NewContactCategory(nil))
}

func TestContactCategory_List(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	contactCategories := make(entity.ContactCategories, 3)
	contactCategories[0] = testContactCategory("contactCategory-id01", "種別1", now())
	contactCategories[1] = testContactCategory("contactCategory-id02", "種別2", now())
	contactCategories[2] = testContactCategory("contactCategory-id03", "種別3", now())
	err = db.DB.Create(&contactCategories).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListContactCategoriesParams
	}
	type want struct {
		contactCategories entity.ContactCategories
		err               error
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
				params: &database.ListContactCategoriesParams{
					Limit:  3,
					Offset: 0,
				},
			},
			want: want{
				contactCategories: contactCategories[:3],
				err:               nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &contactCategory{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.ElementsMatch(t, tt.want.contactCategories, actual)
		})
	}
}

func TestContactCategory_MultiGet(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	contactCategories := make(entity.ContactCategories, 3)
	contactCategories[0] = testContactCategory("contactCategory-id01", "種別1", now())
	contactCategories[1] = testContactCategory("contactCategory-id02", "種別2", now())
	contactCategories[2] = testContactCategory("contactCategory-id03", "種別3", now())
	err = db.DB.Create(&contactCategories).Error
	require.NoError(t, err)

	type args struct {
		categoryIDs []string
	}
	type want struct {
		contactCategories entity.ContactCategories
		err               error
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
				categoryIDs: []string{
					"contactCategory-id01",
					"contactCategory-id02",
					"contactCategory-id03",
				},
			},
			want: want{
				contactCategories: contactCategories,
				err:               nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &contactCategory{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.categoryIDs)
			assert.ErrorIs(t, err, tt.want.err)
			assert.ElementsMatch(t, tt.want.contactCategories, actual)
		})
	}
}

func TestContactCategory_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	c := testContactCategory("category-id", "問い合わせ種別", now())
	err = db.DB.Create(&c).Error
	require.NoError(t, err)

	type args struct {
		categoryID string
	}
	type want struct {
		category *entity.ContactCategory
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
				categoryID: "category-id",
			},
			want: want{
				category: c,
				err:      nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				categoryID: "other-id",
			},
			want: want{
				category: nil,
				err:      database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &contactCategory{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.categoryID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.category, actual)
		})
	}
}

func TestContactCategory_Create(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	c := testContactCategory("category-id", "問い合わせ種別", now())

	type args struct {
		category *entity.ContactCategory
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
				category: c,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				err := db.DB.Create(c).Error
				require.NoError(t, err)
			},
			args: args{
				category: c,
			},
			want: want{
				err: database.ErrAlreadyExists,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			err := delete(ctx, contactCategoryTable)
			require.NoError(t, err)
			tt.setup(ctx, t, db)

			db := &contactCategory{db: db, now: now}
			err = db.Create(ctx, tt.args.category)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testContactCategory(id, name string, now time.Time) *entity.ContactCategory {
	return &entity.ContactCategory{
		ID:        id,
		Title:     name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
