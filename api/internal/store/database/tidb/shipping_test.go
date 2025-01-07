package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShipping(t *testing.T) {
	assert.NotNil(t, NewShipping(nil))
}

func TestShipping_ListByCoordinatorIDs(t *testing.T) {
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

	shippings := make(entity.Shippings, 2)
	shippings[0] = testShipping("shipping-id01", "coordinator-id01", 1, now())
	shippings[1] = testShipping("shipping-id02", "coordinator-id02", 2, now())
	err = db.DB.Create(&shippings).Error
	require.NoError(t, err)
	for i := range shippings {
		internal, err := newInternalShippingRevision(&shippings[i].ShippingRevision)
		require.NoError(t, err)
		err = db.DB.Table(shippingRevisionTable).Create(&internal).Error
		require.NoError(t, err)
	}

	type args struct {
		coordinatorIDs []string
	}
	type want struct {
		shippings entity.Shippings
		hasErr    bool
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
				coordinatorIDs: []string{"coordinator-id01", "coordinator-id02"},
			},
			want: want{
				shippings: shippings,
				hasErr:    false,
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

			db := &shipping{db: db, now: now}
			actual, err := db.ListByCoordinatorIDs(ctx, tt.args.coordinatorIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.shippings, actual)
		})
	}
}

func TestShipping_MultiGet(t *testing.T) {
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

	shippings := make(entity.Shippings, 2)
	shippings[0] = testShipping("shipping-id01", "coordinator-id01", 1, now())
	shippings[1] = testShipping("shipping-id02", "coordinator-id02", 2, now())
	err = db.DB.Create(&shippings).Error
	require.NoError(t, err)
	for i := range shippings {
		internal, err := newInternalShippingRevision(&shippings[i].ShippingRevision)
		require.NoError(t, err)
		err = db.DB.Table(shippingRevisionTable).Create(&internal).Error
		require.NoError(t, err)
	}

	type args struct {
		shippingIDs []string
	}
	type want struct {
		shippings entity.Shippings
		hasErr    bool
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
				shippingIDs: []string{"shipping-id01", "shipping-id02", "shipping-id03"},
			},
			want: want{
				shippings: shippings,
				hasErr:    false,
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

			db := &shipping{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.shippingIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.shippings, actual)
		})
	}
}

func TestShipping_MultiGetByRevision(t *testing.T) {
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

	shippings := make(entity.Shippings, 2)
	shippings[0] = testShipping("shipping-id01", "coordinator-id01", 1, now())
	shippings[1] = testShipping("shipping-id02", "coordinator-id02", 2, now())
	err = db.DB.Create(&shippings).Error
	require.NoError(t, err)
	for i := range shippings {
		internal, err := newInternalShippingRevision(&shippings[i].ShippingRevision)
		require.NoError(t, err)
		err = db.DB.Table(shippingRevisionTable).Create(&internal).Error
		require.NoError(t, err)
	}

	type args struct {
		revisionIDs []int64
	}
	type want struct {
		shippings entity.Shippings
		hasErr    bool
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
				revisionIDs: []int64{1, 2, 3},
			},
			want: want{
				shippings: shippings,
				hasErr:    false,
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

			db := &shipping{db: db, now: now}
			actual, err := db.MultiGetByRevision(ctx, tt.args.revisionIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.shippings, actual)
		})
	}
}

func TestShipping_GetDefault(t *testing.T) {
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

	s := testShipping(entity.DefaultShippingID, "", 1, now())
	err = db.DB.Create(&s).Error
	require.NoError(t, err)
	internal, err := newInternalShippingRevision(&s.ShippingRevision)
	require.NoError(t, err)
	err = db.DB.Table(shippingRevisionTable).Create(&internal).Error
	require.NoError(t, err)

	type args struct{}
	type want struct {
		shipping *entity.Shipping
		hasErr   bool
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
			args:  args{},
			want: want{
				shipping: s,
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

			db := &shipping{db: db, now: now}
			actual, err := db.GetDefault(ctx)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.shipping, actual)
		})
	}
}

func TestShipping_GetByCoordinatorID(t *testing.T) {
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

	s := testShipping("shipping-id", "coordinator-id", 1, now())
	err = db.DB.Create(&s).Error
	require.NoError(t, err)
	internal, err := newInternalShippingRevision(&s.ShippingRevision)
	require.NoError(t, err)
	err = db.DB.Table(shippingRevisionTable).Create(&internal).Error
	require.NoError(t, err)

	type args struct {
		coordinatorID string
	}
	type want struct {
		shipping *entity.Shipping
		hasErr   bool
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
				coordinatorID: "coordinator-id",
			},
			want: want{
				shipping: s,
				hasErr:   false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				coordinatorID: "",
			},
			want: want{
				shipping: nil,
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

			db := &shipping{db: db, now: now}
			actual, err := db.GetByCoordinatorID(ctx, tt.args.coordinatorID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.shipping, actual)
		})
	}
}

func TestShipping_Create(t *testing.T) {
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

	s := testShipping("shipping-id", "coordinator-id", 1, now())

	type args struct {
		shipping *entity.Shipping
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
				shipping: s,
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate key",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				err := db.DB.Create(&s).Error
				require.NoError(t, err)
				internal, err := newInternalShippingRevision(&s.ShippingRevision)
				require.NoError(t, err)
				err = db.DB.Table(shippingRevisionTable).Create(&internal).Error
				require.NoError(t, err)
			},
			args: args{
				shipping: s,
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

			err := delete(ctx, shippingRevisionTable, shippingTable)
			require.NoError(t, err)
			tt.setup(ctx, t, db)

			db := &shipping{db: db, now: now}
			err = db.Create(ctx, tt.args.shipping)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestShipping_Update(t *testing.T) {
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

	s := testShipping("shipping-id", "coordinator-id", 1, now())

	type args struct {
		shippingID string
		params     *database.UpdateShippingParams
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
				err := db.DB.Create(&s).Error
				require.NoError(t, err)
				internal, err := newInternalShippingRevision(&s.ShippingRevision)
				require.NoError(t, err)
				err = db.DB.Table(shippingRevisionTable).Create(&internal).Error
				require.NoError(t, err)
			},
			args: args{
				shippingID: "shipping-id",
				params: &database.UpdateShippingParams{
					Box60Rates:        s.Box60Rates,
					Box60Frozen:       800,
					Box80Rates:        s.Box80Rates,
					Box80Frozen:       800,
					Box100Rates:       s.Box100Rates,
					Box100Frozen:      800,
					HasFreeShipping:   true,
					FreeShippingRates: 3000,
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				shippingID: "shipping-id",
				params:     &database.UpdateShippingParams{},
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

			err := delete(ctx, shippingRevisionTable, shippingTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &shipping{db: db, now: now}
			err = db.Update(ctx, tt.args.shippingID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testShipping(shippingID, coordinatorID string, revisionID int64, now time.Time) *entity.Shipping {
	internal := testShippingRevision(revisionID, shippingID, now)
	revision, _ := internal.entity()
	return &entity.Shipping{
		ID:               shippingID,
		CoordinatorID:    coordinatorID,
		ShippingRevision: *revision,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func testShippingRevision(revisionID int64, shippingID string, now time.Time) *internalShippingRevision {
	shikoku := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int32, 0, 47-len(shikoku))
	for _, val := range codes.PrefectureValues {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := entity.ShippingRates{
		{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
		{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
	}
	revision := &entity.ShippingRevision{
		ID:                revisionID,
		ShippingID:        shippingID,
		Box60Rates:        rates,
		Box60Frozen:       800,
		Box80Rates:        rates,
		Box80Frozen:       800,
		Box100Rates:       rates,
		Box100Frozen:      800,
		HasFreeShipping:   true,
		FreeShippingRates: 3000,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	internal, _ := newInternalShippingRevision(revision)
	return internal
}
