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

func TestSchedule(t *testing.T) {
	assert.NotNil(t, NewSchedule(nil))
}

func TestSchedule_Create(t *testing.T) {
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

	_ = m.dbDelete(ctx, liveProductTable, liveTable, scheduleTable, productTable, productTypeTable, categoryTable, shippingTable)
	category := testCategory("category-id", "野菜", now())
	err = m.db.DB.Create(&category).Error
	require.NoError(t, err)
	productType := testProductType("type-id", "category-id", "野菜", now())
	err = m.db.DB.Create(&productType).Error
	require.NoError(t, err)
	products := make(entity.Products, 2)
	products[0] = testProduct("product-id01", "type-id", "category-id", "producer-id", now())
	products[1] = testProduct("product-id02", "type-id", "category-id", "producer-id", now())
	err = m.db.DB.Create(&products).Error
	require.NoError(t, err)
	shipping := testShipping("shipping-id", now())
	err = m.db.DB.Create(&shipping).Error
	require.NoError(t, err)

	productIDs := []string{"product-id01", "product-id02"}
	s := testSchedule("schedule-id", now())
	lives := testLives("live-id", "schedule-id", "shipping-id", "producer-id", productIDs, now(), 3)
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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				schedule := testSchedule("schedule-id", now())
				err = m.db.DB.Create(&schedule).Error
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

			err := m.dbDelete(ctx, scheduleTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &schedule{db: m.db, now: now}
			err = db.Create(ctx, tt.args.schedule, tt.args.lives, tt.args.products)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testSchedule(id string, now time.Time) *entity.Schedule {
	return &entity.Schedule{
		ID:           id,
		Title:        "旬の夏野菜配信",
		Description:  "旬の夏野菜特集",
		ThumbnailURL: "https://and-period.jp/thumbnail01.png",
		StartAt:      now,
		EndAt:        now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func fillIgnoreScheduleField(s *entity.Schedule, now time.Time) {
	if s == nil {
		return
	}
	s.CreatedAt = now
	s.UpdatedAt = now
}

func fillIgnoreSchedulesField(ss entity.Schedules, now time.Time) {
	for i := range ss {
		fillIgnoreScheduleField(ss[i], now)
	}
}
