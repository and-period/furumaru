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

func TestShop(t *testing.T) {
	assert.NotNil(t, NewShop(nil))
}

func TestShop_List(t *testing.T) {
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

	internal := testShop("shop-id01", "coordinator-id01", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&internal).Error
	require.NoError(t, err)
	s, err := internal.entity()
	require.NoError(t, err)
	s.ProducerIDs = []string{"producer-id01"}

	ps := make(entity.ShopProducers, 1)
	ps[0] = testShopProducer("shop-id01", "producer-id01", now())
	err = db.DB.Table(shopProducerTable).Create(&ps).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListShopsParams
	}
	type want struct {
		shops entity.Shops
		err   error
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
				params: &database.ListShopsParams{
					CoordinatorIDs: []string{"coordinator-id01"},
					ProducerIDs:    []string{"producer-id01"},
				},
			},
			want: want{
				shops: entity.Shops{s},
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &shop{db: db, now: now}
			shops, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.shops, shops)
		})
	}
}

func TestShop_Count(t *testing.T) {
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

	internal := testShop("shop-id01", "coordinator-id01", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&internal).Error
	require.NoError(t, err)
	s, err := internal.entity()
	require.NoError(t, err)
	s.ProducerIDs = []string{"producer-id01"}

	ps := make(entity.ShopProducers, 1)
	ps[0] = testShopProducer("shop-id01", "producer-id01", now())
	err = db.DB.Table(shopProducerTable).Create(&ps).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListShopsParams
	}
	type want struct {
		total int64
		err   error
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
				params: &database.ListShopsParams{
					CoordinatorIDs: []string{"coordinator-id01"},
					ProducerIDs:    []string{"producer-id01"},
				},
			},
			want: want{
				total: 1,
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &shop{db: db, now: now}
			total, err := db.Count(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.total, total)
		})
	}
}

func TestShop_MultiGet(t *testing.T) {
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

	internal := make(internalShops, 1)
	internal[0] = testShop("shop-id01", "coordinator-id01", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&internal).Error
	require.NoError(t, err)
	shops, err := internal.entities()
	require.NoError(t, err)
	shops[0].ProducerIDs = []string{"producer-id01"}

	ps := make(entity.ShopProducers, 1)
	ps[0] = testShopProducer("shop-id01", "producer-id01", now())
	err = db.DB.Table(shopProducerTable).Create(&ps).Error
	require.NoError(t, err)

	type args struct {
		shopIDs []string
	}
	type want struct {
		shops entity.Shops
		err   error
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
				shopIDs: []string{"shop-id01"},
			},
			want: want{
				shops: shops,
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &shop{db: db, now: now}
			shops, err := db.MultiGet(ctx, tt.args.shopIDs)
			assert.ErrorIs(t, err, tt.want.err)
			assert.ElementsMatch(t, tt.want.shops, shops)
		})
	}
}

func TestShop_Get(t *testing.T) {
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

	internal := testShop("shop-id01", "coordinator-id01", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&internal).Error
	require.NoError(t, err)
	s, err := internal.entity()
	require.NoError(t, err)
	s.ProducerIDs = []string{"producer-id01"}

	ps := make(entity.ShopProducers, 1)
	ps[0] = testShopProducer("shop-id01", "producer-id01", now())
	err = db.DB.Table(shopProducerTable).Create(&ps).Error
	require.NoError(t, err)

	type args struct {
		shopID string
	}
	type want struct {
		shop *entity.Shop
		err  error
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
				shopID: "shop-id01",
			},
			want: want{
				shop: s,
				err:  nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				shopID: "",
			},
			want: want{
				shop: nil,
				err:  database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &shop{db: db, now: now}
			shop, err := db.Get(ctx, tt.args.shopID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.shop, shop)
		})
	}
}

func TestShop_GetByCoordinatorID(t *testing.T) {
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

	internal := testShop("shop-id01", "coordinator-id01", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&internal).Error
	require.NoError(t, err)
	s, err := internal.entity()
	require.NoError(t, err)
	s.ProducerIDs = []string{"producer-id01"}

	ps := make(entity.ShopProducers, 1)
	ps[0] = testShopProducer("shop-id01", "producer-id01", now())
	err = db.DB.Table(shopProducerTable).Create(&ps).Error
	require.NoError(t, err)

	type args struct {
		coordinatorID string
	}
	type want struct {
		shop *entity.Shop
		err  error
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
				coordinatorID: "coordinator-id01",
			},
			want: want{
				shop: s,
				err:  nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				coordinatorID: "",
			},
			want: want{
				shop: nil,
				err:  database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &shop{db: db, now: now}
			shop, err := db.GetByCoordinatorID(ctx, tt.args.coordinatorID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.shop, shop)
		})
	}
}

func TestShop_Create(t *testing.T) {
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

	internal := testShop("shop-id01", "coordinator-id01", []string{"producer-id01"}, []string{}, now())
	s, err := internal.entity()
	require.NoError(t, err)

	type args struct {
		shop *entity.Shop
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
				shop: s,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				internal := testShop("shop-id01", "coordinator-id01", []string{"producer-id01"}, []string{}, now())
				err = db.DB.Table(shopTable).Create(&internal).Error
				require.NoError(t, err)
			},
			args: args{
				shop: s,
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

			err := delete(ctx, shopProducerTable, shopTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &shop{db: db, now: now}
			err = db.Create(ctx, tt.args.shop)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestShop_Update(t *testing.T) {
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

	type args struct {
		shopID string
		params *database.UpdateShopParams
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
				internal := testShop("shop-id01", "coordinator-id01", []string{"producer-id01"}, []string{}, now())
				err = db.DB.Table(shopTable).Create(&internal).Error
				require.NoError(t, err)
			},
			args: args{
				shopID: "shop-id01",
				params: &database.UpdateShopParams{
					Name:           "テスト店舗",
					ProductTypeIDs: []string{"product-type-id01"},
					BusinessDays:   []time.Weekday{time.Monday},
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			err := delete(ctx, shopProducerTable, shopTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &shop{db: db, now: now}
			err = db.Update(ctx, tt.args.shopID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestShop_Delete(t *testing.T) {
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

	type args struct {
		shopID string
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
				internal := testShop("shop-id01", "coordinator-id01", []string{"producer-id01"}, []string{}, now())
				err = db.DB.Table(shopTable).Create(&internal).Error
				require.NoError(t, err)
			},
			args: args{
				shopID: "shop-id01",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			err := delete(ctx, shopProducerTable, shopTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &shop{db: db, now: now}
			err = db.Delete(ctx, tt.args.shopID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestShop_RemoveProductType(t *testing.T) {
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

	type args struct {
		productTypeID string
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
				internal := testShop("shop-id01", "coordinator-id01", []string{"producer-id01"}, []string{"product-type-id01"}, now())
				err = db.DB.Table(shopTable).Create(&internal).Error
				require.NoError(t, err)
			},
			args: args{
				productTypeID: "product-type-id01",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			err := delete(ctx, shopProducerTable, shopTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &shop{db: db, now: now}
			err = db.RemoveProductType(ctx, tt.args.productTypeID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestShop_ListProducers(t *testing.T) {
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

	internal := testShop("shop-id01", "coordinator-id01", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&internal).Error
	require.NoError(t, err)
	s, err := internal.entity()
	require.NoError(t, err)
	s.ProducerIDs = []string{"producer-id01"}

	ps := make(entity.ShopProducers, 1)
	ps[0] = testShopProducer("shop-id01", "producer-id01", now())
	err = db.DB.Table(shopProducerTable).Create(&ps).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListShopProducersParams
	}
	type want struct {
		producerIDs []string
		err         error
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
				params: &database.ListShopProducersParams{
					CoordinatorID: "coordinator-id01",
				},
			},
			want: want{
				producerIDs: []string{"producer-id01"},
				err:         nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &shop{db: db, now: now}
			producerIDs, err := db.ListProducers(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.producerIDs, producerIDs)
		})
	}
}

func TestShop_RelateProducer(t *testing.T) {
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

	internal := testShop("shop-id01", "coordinator-id01", []string{"producer-id01"}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&internal).Error
	require.NoError(t, err)

	type args struct {
		shopID     string
		producerID string
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
				shopID:     "shop-id01",
				producerID: "producer-id01",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			err := delete(ctx, shopProducerTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &shop{db: db, now: now}
			err = db.RelateProducer(ctx, tt.args.shopID, tt.args.producerID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestShop_UnrelateProducer(t *testing.T) {
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

	internal := testShop("shop-id01", "coordinator-id01", []string{"producer-id01"}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&internal).Error
	require.NoError(t, err)

	type args struct {
		shopID     string
		producerID string
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
				ps := make(entity.ShopProducers, 1)
				ps[0] = testShopProducer("shop-id01", "producer-id01", now())
				err = db.DB.Table(shopProducerTable).Create(&ps).Error
				require.NoError(t, err)
			},
			args: args{
				shopID:     "shop-id01",
				producerID: "producer-id01",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			err := delete(ctx, shopProducerTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &shop{db: db, now: now}
			err = db.UnrelateProducer(ctx, tt.args.shopID, tt.args.producerID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testShop(shopID, coordinatorID string, producerIDs, productTypeIDs []string, now time.Time) *internalShop {
	shop := &entity.Shop{
		ID:             shopID,
		CoordinatorID:  coordinatorID,
		ProducerIDs:    producerIDs,
		ProductTypeIDs: productTypeIDs,
		BusinessDays:   []time.Weekday{time.Monday},
		Name:           "テスト店舗",
		Activated:      true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	internal, _ := newInternalShop(shop)
	return internal
}

func testShopProducer(shopID, producerID string, now time.Time) *entity.ShopProducer {
	return &entity.ShopProducer{
		ShopID:     shopID,
		ProducerID: producerID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
