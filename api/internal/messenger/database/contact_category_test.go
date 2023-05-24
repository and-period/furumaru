package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContactCategory(t *testing.T) {
	assert.NotNil(t, NewContactCategory(nil))
}

func TestContactCategory_Get(t *testing.T) {
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

	c := testContactCategory("category-id", "問い合わせ種別", now())
	err = db.DB.Create(&c).Error
	require.NoError(t, err)

	type args struct {
		categoryID string
	}
	type want struct {
		category *entity.ContactCategory
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

			db := &contactCategory{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.categoryID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.category, actual)
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
