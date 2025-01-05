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

func TestLive_List(t *testing.T) {
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
	products := make(entity.Products, 1)
	products[0] = testProduct("product-id01", "type-id", "coordinator-id", "producer-id", []string{}, 1, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	for i := range products {
		err = db.DB.Create(&products[i].ProductRevision).Error
		require.NoError(t, err)
	}

	schedule := testSchedule("schedule-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	productIDs := []string{"product-id01"}
	lives := make(entity.Lives, 3)
	lives[0] = testLive("live-id01", "schedule-id", "producer-id", productIDs, now().Add(-time.Hour))
	lives[1] = testLive("live-id02", "schedule-id", "producer-id", productIDs, now())
	lives[2] = testLive("live-id03", "schedule-id", "producer-id", productIDs, now().Add(time.Hour))
	err = db.DB.Create(&lives).Error
	require.NoError(t, err)
	for _, live := range lives {
		err = db.DB.Create(&live.LiveProducts).Error
		require.NoError(t, err)
	}

	type args struct {
		params *database.ListLivesParams
	}
	type want struct {
		lives  entity.Lives
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
				params: &database.ListLivesParams{
					ScheduleIDs: []string{"schedule-id"},
					ProducerID:  "producer-id",
					Limit:       20,
					Offset:      0,
				},
			},
			want: want{
				lives:  lives,
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

			db := &live{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.lives, actual)
		})
	}
}

func TestLive_Count(t *testing.T) {
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
	products := make(entity.Products, 1)
	products[0] = testProduct("product-id01", "type-id", "coordinator-id", "producer-id", []string{}, 1, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	for i := range products {
		err = db.DB.Create(&products[i].ProductRevision).Error
		require.NoError(t, err)
	}

	schedule := testSchedule("schedule-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	productIDs := []string{"product-id01"}
	lives := make(entity.Lives, 3)
	lives[0] = testLive("live-id01", "schedule-id", "producer-id", productIDs, now())
	lives[1] = testLive("live-id02", "schedule-id", "producer-id", productIDs, now())
	lives[2] = testLive("live-id03", "schedule-id", "producer-id", productIDs, now())
	err = db.DB.Create(&lives).Error
	require.NoError(t, err)
	for _, live := range lives {
		err = db.DB.Create(&live.LiveProducts).Error
		require.NoError(t, err)
	}

	type args struct {
		params *database.ListLivesParams
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
				params: &database.ListLivesParams{
					ScheduleIDs: []string{"schedule-id"},
					Limit:       20,
					Offset:      1,
				},
			},
			want: want{
				total:  3,
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

			db := &live{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestLive_Get(t *testing.T) {
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
	products := make(entity.Products, 1)
	products[0] = testProduct("product-id01", "type-id", "coordinator-id", "producer-id", []string{}, 1, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	for i := range products {
		err = db.DB.Create(&products[i].ProductRevision).Error
		require.NoError(t, err)
	}

	schedule := testSchedule("schedule-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	productIDs := []string{"product-id01"}
	l := testLive("live-id", "schedule-id", "producer-id", productIDs, now())
	err = db.DB.Create(&l).Error
	require.NoError(t, err)
	err = db.DB.Create(&l.LiveProducts).Error
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
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				liveID: "live-id",
			},
			want: want{
				live:   l,
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				liveID: "",
			},
			want: want{
				live:   nil,
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

			tt.setup(ctx, t, db)

			db := &live{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.liveID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.live, actual)
		})
	}
}

func TestLive_Update(t *testing.T) {
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
	products := make(entity.Products, 3)
	products[0] = testProduct("product-id01", "type-id", "coordinator-id", "producer-id", []string{}, 1, now())
	products[1] = testProduct("product-id02", "type-id", "coordinator-id", "producer-id", []string{}, 2, now())
	products[2] = testProduct("product-id03", "type-id", "coordinator-id", "producer-id", []string{}, 3, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	for i := range products {
		err = db.DB.Create(&products[i].ProductRevision).Error
		require.NoError(t, err)
	}

	schedule := testSchedule("schedule-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	type args struct {
		liveID string
		params *database.UpdateLiveParams
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
				live := testLive("live-id", "schedule-id", "producer-id", []string{"product-id01", "product-id02"}, now())
				err = db.DB.Create(&live).Error
				require.NoError(t, err)
			},
			args: args{
				liveID: "live-id",
				params: &database.UpdateLiveParams{
					ProductIDs: []string{"product-id02", "product-id03"},
					Comment:    "よろしくお願いします",
					StartAt:    now().AddDate(0, -2, 0),
					EndAt:      now().AddDate(0, 2, 0),
				},
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

			err := delete(ctx, liveTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &live{db: db, now: now}
			err = db.Update(ctx, tt.args.liveID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestLive_Delete(t *testing.T) {
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
	products := make(entity.Products, 3)
	products[0] = testProduct("product-id01", "type-id", "coordinator-id", "producer-id", []string{}, 1, now())
	products[1] = testProduct("product-id02", "type-id", "coordinator-id", "producer-id", []string{}, 2, now())
	products[2] = testProduct("product-id03", "type-id", "coordinator-id", "producer-id", []string{}, 3, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	for i := range products {
		err = db.DB.Create(&products[i].ProductRevision).Error
		require.NoError(t, err)
	}

	schedule := testSchedule("schedule-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	type args struct {
		liveID string
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
				live := testLive("live-id", "schedule-id", "producer-id", []string{"product-id01", "product-id02"}, now())
				err = db.DB.Create(&live).Error
				require.NoError(t, err)
			},
			args: args{
				liveID: "live-id",
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

			err := delete(ctx, liveTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &live{db: db, now: now}
			err = db.Delete(ctx, tt.args.liveID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testLive(liveID, scheduleID, producerID string, productIDs []string, now time.Time) *entity.Live {
	products := make(entity.LiveProducts, len(productIDs))
	for i := range productIDs {
		products[i] = testLiveProduct(liveID, productIDs[i], now)
	}
	return &entity.Live{
		ID:           liveID,
		ScheduleID:   scheduleID,
		ProducerID:   producerID,
		ProductIDs:   productIDs,
		Comment:      "よろしくお願いします。",
		LiveProducts: products,
		StartAt:      now.AddDate(0, -1, 0),
		EndAt:        now.AddDate(0, 1, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
