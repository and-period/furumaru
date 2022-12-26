package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testLive(id, scheduleID, producerID string, productIDs []string, now time.Time) *entity.Live {
	l := &entity.Live{
		ID:          id,
		ScheduleID:  scheduleID,
		ProducerID:  producerID,
		Title:       "配信のタイトル",
		Description: "配信の説明",
		Published:   false,
		Canceled:    false,
		StartAt:     now,
		EndAt:       now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	ps := make(entity.LiveProducts, len(productIDs))
	for i := range productIDs {
		ps[i] = testLiveProduct(id, productIDs[i], now)
	}
	l.Fill(ps, now)
	return l
}

func testLives(id, scheduleID, producerID string, liveIDs []string, now time.Time, length int) entity.Lives {
	lives := make(entity.Lives, length)
	for i := 0; i < length; i++ {
		liveID := fmt.Sprintf("%s-%2d", id, i)
		lives[i] = testLive(liveID, scheduleID, producerID, liveIDs, now)
	}
	return lives
}

func TestLive_Get(t *testing.T) {
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
	schedule := testSchedule("schedule-id", now())
	err = m.db.DB.Create(&schedule).Error
	require.NoError(t, err)
	l := testLive("live-id", "schedule-id", "producer-id", productIDs, now())
	err = m.db.DB.Create(&l).Error
	require.NoError(t, err)
	liveProducts := make(entity.LiveProducts, 2)
	liveProducts[0] = testLiveProduct("live-id", "product-id01", now())
	liveProducts[1] = testLiveProduct("live-id", "product-id02", now())
	err = m.db.DB.Create(&liveProducts).Error
	require.NoError(t, err)

	type args struct {
		liveID string
	}
	type want struct {
		live   *entity.Live
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
				liveID: "live-id",
			},
			want: want{
				live:   l,
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

			tt.setup(ctx, t, m)

			db := &live{db: m.db, now: now}
			actual, err := db.Get(ctx, tt.args.liveID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			fillIgnoreLiveField(actual, now())
			assert.Equal(t, tt.want.live, actual)
		})
	}
}

func TestLive_Update(t *testing.T) {
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

	_ = m.dbDelete(ctx, liveTable, liveProductTable, productTable, productTypeTable, categoryTable)
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

	type args struct {
		liveID string
		params *UpdateLiveParams
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
				live := testLive("live-id", "schedule-id", "producer-id", []string{"product-id01", "product-id02"}, now())
				err = m.db.DB.Create(&live).Error
				require.NoError(t, err)
			},
			args: args{
				liveID: "live-id",
				params: &UpdateLiveParams{
					LiveProducts: entity.LiveProducts{
						{
							LiveID:    "live-id",
							ProductID: "product-id01",
						},
						{
							LiveID:    "live-id",
							ProductID: "product-id02",
						},
					},
					Title:       "じゃがいもの祭典",
					Description: "いろんなじゃがいも勢揃い",
					Published:   true,
					StartAt:     now(),
					EndAt:       now(),
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				liveID: "live-id",
				params: &UpdateLiveParams{},
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

			err := m.dbDelete(ctx, liveTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &live{db: m.db, now: now}
			err = db.Update(ctx, tt.args.liveID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func fillIgnoreLiveField(l *entity.Live, now time.Time) {
	if l == nil {
		return
	}
	l.StartAt = now
	l.EndAt = now
	l.CreatedAt = now
	l.UpdatedAt = now
	for i := range l.LiveProducts {
		l.LiveProducts[i].CreatedAt = now
		l.LiveProducts[i].UpdatedAt = now
	}
}

func fillIgnoreLivesField(ls entity.Lives, now time.Time) {
	for i := range ls {
		fillIgnoreLiveField(ls[i], now)
	}
}
