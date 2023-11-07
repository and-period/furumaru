package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProduct(t *testing.T) {
	assert.NotNil(t, newProduct(nil))
}

func TestProduct_List(t *testing.T) {
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
	products := make(entity.Products, 3)
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "coordinator-id", "producer-id", productTags.IDs(), 1, now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "coordinator-id", "producer-id", productTags.IDs(), 2, now())
	products[2] = testProduct("product-id03", "type-id02", "category-id02", "coordinator-id", "producer-id", productTags.IDs(), 3, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	for i := range products {
		err := db.DB.Create(&products[i].ProductRevision).Error
		require.NoError(t, err)
	}

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
					Limit:         2,
					Offset:        1,
				},
			},
			want: want{
				products: products[1:],
				hasErr:   false,
			},
		},
		{
			name:  "success with sort",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListProductsParams{
					Orders: []*database.ListProductsOrder{
						{Key: entity.ProductOrderByName, OrderByASC: true},
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &product{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.products, actual)
		})
	}
}

func TestProduct_Count(t *testing.T) {
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
	products := make(entity.Products, 3)
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "coordinator-id", "producer-id", productTags.IDs(), 1, now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "coordinator-id", "producer-id", productTags.IDs(), 2, now())
	products[2] = testProduct("product-id03", "type-id02", "category-id02", "coordinator-id", "producer-id", productTags.IDs(), 3, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	for i := range products {
		err := db.DB.Create(&products[i].ProductRevision).Error
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
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
	products := make(entity.Products, 3)
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "coordinator-id", "producer-id", productTags.IDs(), 1, now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "coordinator-id", "producer-id", productTags.IDs(), 2, now())
	products[2] = testProduct("product-id03", "type-id03", "category-id02", "coordinator-id", "producer-id", productTags.IDs(), 3, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	for i := range products {
		err := db.DB.Create(&products[i].ProductRevision).Error
		require.NoError(t, err)
	}

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
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
	products := make(entity.Products, 3)
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "coordinator-id", "producer-id", productTags.IDs(), 1, now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "coordinator-id", "producer-id", productTags.IDs(), 2, now())
	products[2] = testProduct("product-id03", "type-id03", "category-id02", "coordinator-id", "producer-id", productTags.IDs(), 3, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	for i := range products {
		err := db.DB.Create(&products[i].ProductRevision).Error
		require.NoError(t, err)
	}

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &product{db: db, now: now}
			actual, err := db.MultiGetByRevision(ctx, tt.args.revisionIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.products, actual)
		})
	}
}

func TestProduct_Get(t *testing.T) {
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
	productTag := testProductTag("tag-id", "贈答品", now())
	err = db.DB.Create(&productTag).Error
	require.NoError(t, err)
	p := testProduct("product-id", "type-id", "category-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now())
	err = db.DB.Create(&p).Error
	require.NoError(t, err)
	err = db.DB.Create(&p.ProductRevision).Error
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
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
	productTag := testProductTag("tag-id", "贈答品", now())
	err = db.DB.Create(&productTag).Error
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
				product: testProduct("product-id", "type-id", "category-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				product := testProduct("product-id", "type-id", "category-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now())
				err = db.DB.Create(&product).Error
				require.NoError(t, err)
				err = db.DB.Create(&product.ProductRevision).Error
				require.NoError(t, err)
			},
			args: args{
				product: testProduct("product-id", "type-id", "category-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now()),
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
				product := testProduct("product-id", "type-id", "category-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now())
				err = db.DB.Create(&product).Error
				require.NoError(t, err)
				err = db.DB.Create(&product.ProductRevision).Error
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
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

func TestProduct_UpdateMedia(t *testing.T) {
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
	productTag := testProductTag("tag-id", "贈答品", now())
	err = db.DB.Create(&productTag).Error
	require.NoError(t, err)

	type args struct {
		productID string
		set       func(media entity.MultiProductMedia) bool
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
				product := testProduct("product-id", "type-id", "category-id", "coordinator-id", "producer-id", []string{"product-id"}, 1, now())
				err = db.DB.Create(&product).Error
				require.NoError(t, err)
				err = db.DB.Create(&product.ProductRevision).Error
				require.NoError(t, err)
			},
			args: args{
				productID: "product-id",
				set: func(media entity.MultiProductMedia) (exists bool) {
					resized := map[string]common.Images{
						"https://and-period.jp/thumbnail01.png": {{
							Size: common.ImageSizeSmall,
							URL:  "https://and-period.jp/thumbnail01_240.png",
						}},
					}
					for i := range media {
						images, ok := resized[media[i].URL]
						if !ok {
							continue
						}
						exists = true
						media[i].Images = images
					}
					return
				},
			},
			want: want{
				hasErr: false,
			},
		},
		// {
		// 	name:  "not found",
		// 	setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
		// 	args: args{
		// 		productID: "product-id",
		// 		set:       func(media entity.MultiProductMedia) bool { return false },
		// 	},
		// 	want: want{
		// 		hasErr: true,
		// 	},
		// },
		// {
		// 	name: "media is non existent",
		// 	setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
		// 		product := testProduct("product-id", "type-id", "category-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now())
		// 		err = db.DB.Create(&product).Error
		// 		require.NoError(t, err)
		// 		err = db.DB.Create(&product.ProductRevision).Error
		// 		require.NoError(t, err)
		// 	},
		// 	args: args{
		// 		productID: "product-id",
		// 		set:       func(media entity.MultiProductMedia) bool { return false },
		// 	},
		// 	want: want{
		// 		hasErr: true,
		// 	},
		// },
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, productRevisionTable, productTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &product{db: db, now: now}
			err = db.UpdateMedia(ctx, tt.args.productID, tt.args.set)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestProduct_Delete(t *testing.T) {
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
				product := testProduct("product-id", "type-id", "category-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now())
				err = db.DB.Create(&product).Error
				require.NoError(t, err)
				err = db.DB.Create(&product.ProductRevision).Error
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
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

func testProduct(productID, typeID, categoryID, coordinatorID, producerID string, tagIDs []string, revisionID int64, now time.Time) *entity.Product {
	p := &entity.Product{
		ID:              productID,
		TypeID:          typeID,
		TagIDs:          tagIDs,
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
		Media: entity.MultiProductMedia{
			{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
			{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
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
	_ = p.FillJSON()
	return p
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
