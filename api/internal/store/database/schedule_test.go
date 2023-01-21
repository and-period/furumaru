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

func TestSchedule(t *testing.T) {
	assert.NotNil(t, NewSchedule(nil))
}

func TestSchedule_Get(t *testing.T) {
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

	shipping := testShipping("shipping-id", now())
	err = db.DB.Create(&shipping).Error
	require.NoError(t, err)
	s := testSchedule("schedule-id", now())
	err = db.DB.Create(&s).Error
	require.NoError(t, err)

	type args struct {
		scheduleID string
	}
	type want struct {
		schedule *entity.Schedule
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
				scheduleID: "schedule-id",
			},
			want: want{
				schedule: s,
				hasErr:   false,
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

			db := &schedule{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.scheduleID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.schedule, actual)
		})
	}
}

func TestSchedule_Create(t *testing.T) {
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

	category := testCategory("category-id", "野菜", now())
	err = db.DB.Create(&category).Error
	require.NoError(t, err)
	productType := testProductType("type-id", "category-id", "野菜", now())
	err = db.DB.Create(&productType).Error
	require.NoError(t, err)
	products := make(entity.Products, 2)
	products[0] = testProduct("product-id01", "type-id", "category-id", "producer-id", now())
	products[1] = testProduct("product-id02", "type-id", "category-id", "producer-id", now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	shipping := testShipping("shipping-id", now())
	err = db.DB.Create(&shipping).Error
	require.NoError(t, err)

	productIDs := []string{"product-id01", "product-id02"}
	s := testSchedule("schedule-id", now())
	lives := testLives("live-id", "schedule-id", "producer-id", productIDs, now(), 3)
	lproducts := make(entity.LiveProducts, 0)
	for i := range lives {
		lproducts = append(lproducts, lives[i].LiveProducts...)
	}

	type args struct {
		schedule *entity.Schedule
		lives    entity.Lives
		products entity.LiveProducts
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
				schedule: s,
				lives:    lives,
				products: lproducts,
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				schedule := testSchedule("schedule-id", now())
				err = db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				schedule: s,
				lives:    lives,
				products: lproducts,
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

			err := delete(ctx, scheduleTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &schedule{db: db, now: now}
			err = db.Create(ctx, tt.args.schedule, tt.args.lives, tt.args.products)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testSchedule(id string, now time.Time) *entity.Schedule {
	return &entity.Schedule{
		ID:           id,
		ShippingID:   "shipping-id",
		Title:        "旬の夏野菜配信",
		Description:  "旬の夏野菜特集",
		ThumbnailURL: "https://and-period.jp/thumbnail01.png",
		StartAt:      now,
		EndAt:        now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
