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
	"go.uber.org/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShipping(t *testing.T) {
	assert.NotNil(t, NewShipping(nil))
}

func TestShipping_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shippings := make(entity.Shippings, 2)
	shippings[0] = testShipping("shipping-id01", "shop-id01", "coordinator-id01", 1, now())
	shippings[1] = testShipping("shipping-id02", "shop-id02", "coordinator-id02", 2, now())
	err = db.DB.Create(&shippings).Error
	require.NoError(t, err)
	for i := range shippings {
		internal := newInternalShippingRevision(&shippings[i].ShippingRevision)
		err = db.DB.Table(shippingRevisionTable).Create(&internal).Error
		require.NoError(t, err)
	}

	type args struct {
		params *database.ListShippingsParams
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
				params: &database.ListShippingsParams{
					Limit:     10,
					Offset:    0,
					OnlyInUse: true,
				},
			},
			want: want{
				shippings: shippings,
				hasErr:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &shipping{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.shippings, actual)
		})
	}
}

func TestShipping_ListByShopIDs(t *testing.T) {
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

	shippings := make(entity.Shippings, 2)
	shippings[0] = testShipping("shipping-id01", "shop-id01", "coordinator-id01", 1, now())
	shippings[1] = testShipping("shipping-id02", "shop-id02", "coordinator-id02", 2, now())
	err = db.DB.Create(&shippings).Error
	require.NoError(t, err)
	for i := range shippings {
		internal := newInternalShippingRevision(&shippings[i].ShippingRevision)
		err = db.DB.Table(shippingRevisionTable).Create(&internal).Error
		require.NoError(t, err)
	}

	type args struct {
		shopIDs []string
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
				shopIDs: []string{"shop-id01", "shop-id02"},
			},
			want: want{
				shippings: shippings,
				hasErr:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &shipping{db: db, now: now}
			actual, err := db.ListByShopIDs(ctx, tt.args.shopIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.shippings, actual)
		})
	}
}

func TestShipping_MultiGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shippings := make(entity.Shippings, 2)
	shippings[0] = testShipping("shipping-id01", "shop-id01", "coordinator-id01", 1, now())
	shippings[1] = testShipping("shipping-id02", "shop-id02", "coordinator-id02", 2, now())
	err = db.DB.Create(&shippings).Error
	require.NoError(t, err)
	for i := range shippings {
		internal := newInternalShippingRevision(&shippings[i].ShippingRevision)
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &shipping{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.shippingIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.shippings, actual)
		})
	}
}

func TestShipping_MultiGetByRevision(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shippings := make(entity.Shippings, 2)
	shippings[0] = testShipping("shipping-id01", "shop-id01", "coordinator-id01", 1, now())
	shippings[1] = testShipping("shipping-id02", "shop-id02", "coordinator-id02", 2, now())
	err = db.DB.Create(&shippings).Error
	require.NoError(t, err)
	for i := range shippings {
		internal := newInternalShippingRevision(&shippings[i].ShippingRevision)
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &shipping{db: db, now: now}
			actual, err := db.MultiGetByRevision(ctx, tt.args.revisionIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.shippings, actual)
		})
	}
}

func TestShipping_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	s := testShipping("shipping-id", "shop-id", "coordinator-id", 1, now())
	err = db.DB.Create(&s).Error
	require.NoError(t, err)

	internal := newInternalShippingRevision(&s.ShippingRevision)
	err = db.DB.Table(shippingRevisionTable).Create(&internal).Error
	require.NoError(t, err)

	type args struct {
		shippingID string
	}
	type want struct {
		shipping *entity.Shipping
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
				shippingID: "shipping-id",
			},
			want: want{
				shipping: s,
				err:      nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				shippingID: "",
			},
			want: want{
				shipping: nil,
				err:      database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &shipping{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.shippingID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.shipping, actual)
		})
	}
}

func TestShipping_GetDefault(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	s := testShipping(entity.DefaultShippingID, "shop-id", "coordinator-id", 1, now())
	err = db.DB.Create(&s).Error
	require.NoError(t, err)
	internal := newInternalShippingRevision(&s.ShippingRevision)
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &shipping{db: db, now: now}
			actual, err := db.GetDefault(ctx)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.shipping, actual)
		})
	}
}

func TestShipping_GetByShopID(t *testing.T) {
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

	s := testShipping("shipping-id", "shop-id", "coordinator-id", 1, now())
	err = db.DB.Create(&s).Error
	require.NoError(t, err)
	internal := newInternalShippingRevision(&s.ShippingRevision)
	err = db.DB.Table(shippingRevisionTable).Create(&internal).Error
	require.NoError(t, err)

	type args struct {
		shopID string
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
				shopID: "shop-id",
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
				shopID: "",
			},
			want: want{
				shipping: nil,
				hasErr:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &shipping{db: db, now: now}
			actual, err := db.GetByShopID(ctx, tt.args.shopID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.shipping, actual)
		})
	}
}

func TestShipping_GetByCoordinatorID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	s := testShipping("shipping-id", "shop-id", "coordinator-id", 1, now())
	err = db.DB.Create(&s).Error
	require.NoError(t, err)
	internal := newInternalShippingRevision(&s.ShippingRevision)
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &shipping{db: db, now: now}
			actual, err := db.GetByCoordinatorID(ctx, tt.args.coordinatorID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.shipping, actual)
		})
	}
}

func TestShipping_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	s := testShipping("shipping-id", "shop-id", "coordinator-id", 1, now())

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
				internal := newInternalShippingRevision(&s.ShippingRevision)
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
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	s := testShipping("shipping-id", "shop-id", "coordinator-id", 1, now())

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
				internal := newInternalShippingRevision(&s.ShippingRevision)
				err = db.DB.Table(shippingRevisionTable).Create(&internal).Error
				require.NoError(t, err)
			},
			args: args{
				shippingID: "shipping-id",
				params: &database.UpdateShippingParams{
					Name:              "配送設定",
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
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, shippingRevisionTable, shippingTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &shipping{db: db, now: now}
			err = db.Update(ctx, tt.args.shippingID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestShipping_UpdateInUse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	type args struct {
		shopID     string
		shippingID string
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
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				s := testShipping("shipping-id", "shop-id", "coordinator-id", 1, now())
				s.InUse = false
				err := db.DB.Create(&s).Error
				require.NoError(t, err)
				internal := newInternalShippingRevision(&s.ShippingRevision)
				err = db.DB.Table(shippingRevisionTable).Create(&internal).Error
				require.NoError(t, err)
			},
			args: args{
				shopID:     "shop-id",
				shippingID: "shipping-id",
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already in use",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				s := testShipping("shipping-id", "shop-id", "coordinator-id", 1, now())
				s.InUse = true
				err := db.DB.Create(&s).Error
				require.NoError(t, err)
				internal := newInternalShippingRevision(&s.ShippingRevision)
				err = db.DB.Table(shippingRevisionTable).Create(&internal).Error
				require.NoError(t, err)
			},
			args: args{
				shopID:     "shop-id",
				shippingID: "shipping-id",
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "shop id mismatch",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				s := testShipping("shipping-id", "shop-id", "coordinator-id", 1, now())
				err := db.DB.Create(&s).Error
				require.NoError(t, err)
				internal := newInternalShippingRevision(&s.ShippingRevision)
				err = db.DB.Table(shippingRevisionTable).Create(&internal).Error
				require.NoError(t, err)
			},
			args: args{
				shopID:     "other-id",
				shippingID: "shipping-id",
			},
			want: want{
				err: database.ErrFailedPrecondition,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, shippingRevisionTable, shippingTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &shipping{db: db, now: now}
			err = db.UpdateInUse(ctx, tt.args.shopID, tt.args.shippingID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestShipping_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	type args struct {
		shippingID string
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
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				s := testShipping("shipping-id", "shop-id", "coordinator-id", 1, now())
				s.InUse = false
				err := db.DB.Create(&s).Error
				require.NoError(t, err)
				internal := newInternalShippingRevision(&s.ShippingRevision)
				err = db.DB.Table(shippingRevisionTable).Create(&internal).Error
				require.NoError(t, err)
			},
			args: args{
				shippingID: "shipping-id",
			},
			want: want{
				err: nil,
			},
		},
		{
			name:  "default shipping",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				shippingID: entity.DefaultShippingID,
			},
			want: want{
				err: database.ErrPermissionDenied,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				shippingID: "",
			},
			want: want{
				err: database.ErrNotFound,
			},
		},
		{
			name: "in use",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				s := testShipping("shipping-id", "shop-id", "coordinator-id", 1, now())
				s.InUse = true
				err := db.DB.Create(&s).Error
				require.NoError(t, err)
			},
			args: args{
				shippingID: "shipping-id",
			},
			want: want{
				err: database.ErrFailedPrecondition,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, shippingRevisionTable, shippingTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &shipping{db: db, now: now}
			err = db.Delete(ctx, tt.args.shippingID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testShipping(shippingID, shopID, coordinatorID string, revisionID int64, now time.Time) *entity.Shipping {
	internal := testShippingRevision(revisionID, shippingID, now)
	revision := internal.entity()
	return &entity.Shipping{
		ID:               shippingID,
		ShopID:           shopID,
		CoordinatorID:    coordinatorID,
		Name:             "配送設定",
		InUse:            true,
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
	internal := newInternalShippingRevision(revision)
	return internal
}
