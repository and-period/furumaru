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

func TestProduct(t *testing.T) {
	assert.NotNil(t, NewProduct(nil))
}

func TestProduct_List(t *testing.T) {
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

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	categories := make(entity.Categories, 2)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	err = db.DB.Create(&categories).Error
	require.NoError(t, err)
	productTypes := make(entity.ProductTypes, 2)
	productTypes[0] = testProductType("type-id01", "category-id01", "野菜", now())
	productTypes[1] = testProductType("type-id02", "category-id02", "果物", now())
	err = db.DB.Create(&productTypes).Error
	require.NoError(t, err)
	productTags := make(entity.ProductTags, 2)
	productTags[0] = testProductTag("tag-id01", "贈答品", now())
	productTags[1] = testProductTag("tag-id02", "有機野菜", now())
	err = db.DB.Create(&productTags).Error
	require.NoError(t, err)
	internal := make(internalProducts, 3)
	internal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", productTags.IDs(), 1, now())
	internal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", productTags.IDs(), 2, now())
	internal[2] = testProduct("product-id03", "type-id02", "shop-id", "coordinator-id", "producer-id", productTags.IDs(), 3, now())
	err = db.DB.Table(productTable).Create(&internal).Error
	require.NoError(t, err)
	for i := range internal {
		err := db.DB.Create(&internal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	products, err := internal.entities()
	require.NoError(t, err)

	type args struct {
		params *database.ListProductsParams
	}
	type want struct {
		products entity.Products
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
				params: &database.ListProductsParams{
					Name:          "いも",
					CoordinatorID: "coordinator-id",
					ProducerID:    "producer-id",
					ProducerIDs:   []string{"producer-id"},
					Limit:         3,
					Offset:        0,
				},
			},
			want: want{
				products: products,
				hasErr:   false,
			},
		},
		{
			name:  "success with sort",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListProductsParams{
					Orders: []*database.ListProductsOrder{
						{Key: database.ListProductsOrderByName, OrderByASC: true},
					},
				},
			},
			want: want{
				products: products,
				hasErr:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &product{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.products, actual)
		})
	}
}

func TestProduct_Count(t *testing.T) {
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

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	categories := make(entity.Categories, 2)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	err = db.DB.Create(&categories).Error
	require.NoError(t, err)
	productTypes := make(entity.ProductTypes, 2)
	productTypes[0] = testProductType("type-id01", "category-id01", "野菜", now())
	productTypes[1] = testProductType("type-id02", "category-id02", "果物", now())
	err = db.DB.Create(&productTypes).Error
	require.NoError(t, err)
	productTags := make(entity.ProductTags, 2)
	productTags[0] = testProductTag("tag-id01", "贈答品", now())
	productTags[1] = testProductTag("tag-id02", "有機野菜", now())
	err = db.DB.Create(&productTags).Error
	require.NoError(t, err)
	internal := make(internalProducts, 3)
	internal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", productTags.IDs(), 1, now())
	internal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", productTags.IDs(), 2, now())
	internal[2] = testProduct("product-id03", "type-id02", "shop-id", "coordinator-id", "producer-id", productTags.IDs(), 3, now())
	err = db.DB.Table(productTable).Create(&internal).Error
	require.NoError(t, err)
	for i := range internal {
		err := db.DB.Create(&internal[i].ProductRevision).Error
		require.NoError(t, err)
	}

	type args struct {
		params *database.ListProductsParams
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
				params: &database.ListProductsParams{
					Name:          "いも",
					CoordinatorID: "coordinator-id",
					ProducerID:    "producer-id",
				},
			},
			want: want{
				total:  3,
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &product{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestProduct_MultiGet(t *testing.T) {
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

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	categories := make(entity.Categories, 2)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	err = db.DB.Create(&categories).Error
	require.NoError(t, err)
	productTypes := make(entity.ProductTypes, 3)
	productTypes[0] = testProductType("type-id01", "category-id01", "野菜", now())
	productTypes[1] = testProductType("type-id02", "category-id02", "果物", now())
	productTypes[2] = testProductType("type-id03", "category-id02", "水産物", now())
	err = db.DB.Create(&productTypes).Error
	require.NoError(t, err)
	productTags := make(entity.ProductTags, 2)
	productTags[0] = testProductTag("tag-id01", "贈答品", now())
	productTags[1] = testProductTag("tag-id02", "有機野菜", now())
	err = db.DB.Create(&productTags).Error
	require.NoError(t, err)
	internal := make(internalProducts, 3)
	internal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", productTags.IDs(), 1, now())
	internal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", productTags.IDs(), 2, now())
	internal[2] = testProduct("product-id03", "type-id02", "shop-id", "coordinator-id", "producer-id", productTags.IDs(), 3, now())
	err = db.DB.Table(productTable).Create(&internal).Error
	require.NoError(t, err)
	for i := range internal {
		err := db.DB.Create(&internal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	products, err := internal.entities()
	require.NoError(t, err)

	type args struct {
		productIDs []string
	}
	type want struct {
		products entity.Products
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
				productIDs: []string{"product-id01", "product-id02", "product-id03"},
			},
			want: want{
				products: products[:3],
				hasErr:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &product{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.productIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.products, actual)
		})
	}
}

func TestProduct_MultiGetByRevision(t *testing.T) {
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

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	categories := make(entity.Categories, 2)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	err = db.DB.Create(&categories).Error
	require.NoError(t, err)
	productTypes := make(entity.ProductTypes, 3)
	productTypes[0] = testProductType("type-id01", "category-id01", "野菜", now())
	productTypes[1] = testProductType("type-id02", "category-id02", "果物", now())
	productTypes[2] = testProductType("type-id03", "category-id02", "水産物", now())
	err = db.DB.Create(&productTypes).Error
	require.NoError(t, err)
	productTags := make(entity.ProductTags, 2)
	productTags[0] = testProductTag("tag-id01", "贈答品", now())
	productTags[1] = testProductTag("tag-id02", "有機野菜", now())
	err = db.DB.Create(&productTags).Error
	require.NoError(t, err)
	internal := make(internalProducts, 3)
	internal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", productTags.IDs(), 1, now())
	internal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", productTags.IDs(), 2, now())
	internal[2] = testProduct("product-id03", "type-id02", "shop-id", "coordinator-id", "producer-id", productTags.IDs(), 3, now())
	err = db.DB.Table(productTable).Create(&internal).Error
	require.NoError(t, err)
	for i := range internal {
		err := db.DB.Create(&internal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	products, err := internal.entities()
	require.NoError(t, err)

	type args struct {
		revisionIDs []int64
	}
	type want struct {
		products entity.Products
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
				revisionIDs: []int64{1, 2, 3},
			},
			want: want{
				products: products,
				hasErr:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &product{db: db, now: now}
			actual, err := db.MultiGetByRevision(ctx, tt.args.revisionIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.products, actual)
		})
	}
}

func TestProduct_Get(t *testing.T) {
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

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	category := testCategory("category-id", "野菜", now())
	err = db.DB.Create(&category).Error
	require.NoError(t, err)
	productType := testProductType("type-id", "category-id", "野菜", now())
	err = db.DB.Create(&productType).Error
	require.NoError(t, err)
	productTag := testProductTag("tag-id", "贈答品", now())
	err = db.DB.Create(&productTag).Error
	require.NoError(t, err)
	internal := testProduct("product-id", "type-id", "shop-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now())
	err = db.DB.Table(productTable).Create(&internal).Error
	require.NoError(t, err)
	err = db.DB.Create(&internal.ProductRevision).Error
	require.NoError(t, err)
	p, err := internal.entity()
	require.NoError(t, err)

	type args struct {
		productID string
	}
	type want struct {
		product *entity.Product
		hasErr  bool
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
				productID: "product-id",
			},
			want: want{
				product: p,
				hasErr:  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &product{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.productID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.product, actual)
		})
	}
}

func TestProduct_Create(t *testing.T) {
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

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	category := testCategory("category-id", "野菜", now())
	err = db.DB.Create(&category).Error
	require.NoError(t, err)
	productType := testProductType("type-id", "category-id", "野菜", now())
	err = db.DB.Create(&productType).Error
	require.NoError(t, err)
	productTag := testProductTag("tag-id", "贈答品", now())
	err = db.DB.Create(&productTag).Error
	require.NoError(t, err)

	internal := testProduct("product-id", "type-id", "shop-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now())
	p, err := internal.entity()
	require.NoError(t, err)

	type args struct {
		product *entity.Product
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
				product: p,
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				internal := testProduct("product-id", "type-id", "shop-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now())
				err := db.DB.Table(productTable).Create(&internal).Error
				require.NoError(t, err)
				err = db.DB.Create(&internal.ProductRevision).Error
				require.NoError(t, err)
			},
			args: args{
				product: p,
			},
			want: want{
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			err := delete(ctx, productRevisionTable, productTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &product{db: db, now: now}
			err = db.Create(ctx, tt.args.product)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestProduct_Update(t *testing.T) {
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

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	category := testCategory("category-id", "野菜", now())
	err = db.DB.Create(&category).Error
	require.NoError(t, err)
	productType := testProductType("type-id", "category-id", "野菜", now())
	err = db.DB.Create(&productType).Error
	require.NoError(t, err)
	productTag := testProductTag("tag-id", "贈答品", now())
	err = db.DB.Create(&productTag).Error
	require.NoError(t, err)

	type args struct {
		productID string
		params    *database.UpdateProductParams
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
				internal := testProduct("product-id", "type-id", "shop-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now())
				err = db.DB.Table(productTable).Create(&internal).Error
				require.NoError(t, err)
				err = db.DB.Create(&internal.ProductRevision).Error
				require.NoError(t, err)
			},
			args: args{
				productID: "product-id",
				params: &database.UpdateProductParams{
					TypeID:          "type-id",
					TagIDs:          []string{"tag-id"},
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Public:          true,
					Inventory:       100,
					Weight:          100,
					WeightUnit:      entity.WeightUnitGram,
					Item:            1,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: entity.MultiProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
					Price:                400,
					Cost:                 300,
					StorageMethodType:    entity.StorageMethodTypeNormal,
					DeliveryType:         entity.DeliveryTypeNormal,
					Box60Rate:            50,
					Box80Rate:            40,
					Box100Rate:           30,
					OriginPrefectureCode: 25,
					OriginCity:           "彦根市",
					StartAt:              now().AddDate(0, -1, 0),
					EndAt:                now().AddDate(0, 1, 0),
				},
			},
			want: want{
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			err := delete(ctx, productRevisionTable, productTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &product{db: db, now: now}
			err = db.Update(ctx, tt.args.productID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestProduct_DescreaseInventory(t *testing.T) {
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

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	category := testCategory("category-id", "野菜", now())
	err = db.DB.Create(&category).Error
	require.NoError(t, err)
	productType := testProductType("type-id", "category-id", "野菜", now())
	err = db.DB.Create(&productType).Error
	require.NoError(t, err)
	productTag := testProductTag("tag-id", "贈答品", now())
	err = db.DB.Create(&productTag).Error
	require.NoError(t, err)

	type args struct {
		revisionID int64
		quantity   int64
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
				internal := testProduct("product-id", "type-id", "shop-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now())
				err = db.DB.Table(productTable).Create(&internal).Error
				require.NoError(t, err)
				err = db.DB.Create(&internal.ProductRevision).Error
				require.NoError(t, err)
			},
			args: args{
				revisionID: 1,
				quantity:   2,
			},
			want: want{
				err: nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				revisionID: 1,
				quantity:   2,
			},
			want: want{
				err: database.ErrNotFound,
			},
		},
		{
			name: "already empty",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				internal := testProduct("product-id", "type-id", "shop-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now())
				internal.Inventory = 0
				err = db.DB.Table(productTable).Create(&internal).Error
				require.NoError(t, err)
				err = db.DB.Create(&internal.ProductRevision).Error
				require.NoError(t, err)
			},
			args: args{
				revisionID: 1,
				quantity:   2,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "less then 0",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				internal := testProduct("product-id", "type-id", "shop-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now())
				internal.Inventory = 1
				err = db.DB.Table(productTable).Create(&internal).Error
				require.NoError(t, err)
				err = db.DB.Create(&internal.ProductRevision).Error
				require.NoError(t, err)
			},
			args: args{
				revisionID: 1,
				quantity:   2,
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

			err := delete(ctx, productRevisionTable, productTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &product{db: db, now: now}
			err = db.DecreaseInventory(ctx, tt.args.revisionID, tt.args.quantity)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestProduct_Delete(t *testing.T) {
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

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	category := testCategory("category-id", "野菜", now())
	err = db.DB.Create(&category).Error
	require.NoError(t, err)
	productType := testProductType("type-id", "category-id", "野菜", now())
	err = db.DB.Create(&productType).Error
	require.NoError(t, err)
	productTag := testProductTag("tag-id", "贈答品", now())
	err = db.DB.Create(&productTag).Error
	require.NoError(t, err)

	type args struct {
		productID string
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
				internal := testProduct("product-id", "type-id", "shop-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now())
				err = db.DB.Table(productTable).Create(&internal).Error
				require.NoError(t, err)
				err = db.DB.Create(&internal.ProductRevision).Error
				require.NoError(t, err)
			},
			args: args{
				productID: "product-id",
			},
			want: want{
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			err := delete(ctx, productRevisionTable, productTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &product{db: db, now: now}
			err = db.Delete(ctx, tt.args.productID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testProduct(productID, typeID, shopID, coordinatorID, producerID string, tagIDs []string, revisionID int64, now time.Time) *internalProduct {
	product := &entity.Product{
		ID:              productID,
		TypeID:          typeID,
		TagIDs:          tagIDs,
		ShopID:          shopID,
		CoordinatorID:   coordinatorID,
		ProducerID:      producerID,
		Name:            "新鮮なじゃがいも",
		Description:     "新鮮なじゃがいもをお届けします。",
		Public:          true,
		Status:          entity.ProductStatusForSale,
		Inventory:       100,
		Weight:          100,
		WeightUnit:      entity.WeightUnitGram,
		Item:            1,
		ItemUnit:        "袋",
		ItemDescription: "1袋あたり100gのじゃがいも",
		ThumbnailURL:    "https://and-period.jp/thumbnail01.png",
		Media: entity.MultiProductMedia{
			{
				URL:         "https://and-period.jp/thumbnail01.png",
				IsThumbnail: true,
			},
			{
				URL:         "https://and-period.jp/thumbnail02.png",
				IsThumbnail: false,
			},
		},
		ExpirationDate:       7,
		StorageMethodType:    entity.StorageMethodTypeNormal,
		DeliveryType:         entity.DeliveryTypeNormal,
		Box60Rate:            50,
		Box80Rate:            40,
		Box100Rate:           30,
		OriginPrefecture:     "滋賀県",
		OriginPrefectureCode: 25,
		OriginCity:           "彦根市",
		StartAt:              now.AddDate(0, -1, 0),
		EndAt:                now.AddDate(0, 1, 0),
		ProductRevision:      *testProductRevision(revisionID, productID, now),
		CreatedAt:            now,
		UpdatedAt:            now,
	}
	internal, _ := newInternalProduct(product)
	return internal
}

func testProductRevision(revisionID int64, productID string, now time.Time) *entity.ProductRevision {
	return &entity.ProductRevision{
		ID:        revisionID,
		ProductID: productID,
		Price:     400,
		Cost:      300,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
