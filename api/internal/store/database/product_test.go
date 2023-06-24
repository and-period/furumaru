package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProduct(t *testing.T) {
	assert.NotNil(t, NewProduct(nil))
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
	productTags[1] = testProductTag("tag-id01", "有機野菜", now())
	err = db.DB.Create(&productTags).Error
	require.NoError(t, err)
	products := make(entity.Products, 3)
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "producer-id", productTags.IDs(), now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "producer-id", productTags.IDs(), now())
	products[2] = testProduct("product-id03", "type-id02", "category-id02", "producer-id", productTags.IDs(), now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)

	type args struct {
		params *ListProductsParams
	}
	type want struct {
		products entity.Products
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
				params: &ListProductsParams{
					Name:        "いも",
					ProducerID:  "producer-id",
					ProducerIDs: []string{"producer-id"},
					Limit:       2,
					Offset:      1,
				},
			},
			want: want{
				products: products[1:],
				hasErr:   false,
			},
		},
		{
			name:  "success with sort",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &ListProductsParams{
					Orders: []*ListProductsOrder{
						{Key: entity.ProductOrderByName, OrderByASC: true},
						{Key: entity.ProductOrderByPrice, OrderByASC: false},
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
			fillIgnoreProductsField(actual, now())
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
	productTags[1] = testProductTag("tag-id01", "有機野菜", now())
	err = db.DB.Create(&productTags).Error
	require.NoError(t, err)
	products := make(entity.Products, 3)
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "producer-id", productTags.IDs(), now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "producer-id", productTags.IDs(), now())
	products[2] = testProduct("product-id03", "type-id02", "category-id02", "producer-id", productTags.IDs(), now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)

	type args struct {
		params *ListProductsParams
	}
	type want struct {
		total  int64
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
				params: &ListProductsParams{
					Name:       "いも",
					ProducerID: "producer-id",
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
	productTags[1] = testProductTag("tag-id01", "有機野菜", now())
	err = db.DB.Create(&productTags).Error
	require.NoError(t, err)
	products := make(entity.Products, 3)
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "producer-id", productTags.IDs(), now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "producer-id", productTags.IDs(), now())
	products[2] = testProduct("product-id03", "type-id03", "category-id02", "producer-id", productTags.IDs(), now())
	err = db.DB.Create(&products).Error
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
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
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
			fillIgnoreProductsField(actual, now())
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
	p := testProduct("product-id", "type-id", "category-id", "producer-id", []string{"tag-id"}, now())
	err = db.DB.Create(&p).Error
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
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
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
			fillIgnoreProductField(actual, now())
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
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				product: testProduct("product-id", "type-id", "category-id", "producer-id", []string{"tag-id"}, now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				product := testProduct("product-id", "type-id", "category-id", "producer-id", []string{"tag-id"}, now())
				err = db.DB.Create(&product).Error
				require.NoError(t, err)
			},
			args: args{
				product: testProduct("product-id", "type-id", "category-id", "producer-id", []string{"tag-id"}, now()),
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

			err := delete(ctx, productTable)
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
		params    *UpdateProductParams
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
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				product := testProduct("product-id", "type-id", "category-id", "producer-id", []string{"tag-id"}, now())
				err = db.DB.Create(&product).Error
				require.NoError(t, err)
			},
			args: args{
				productID: "product-id",
				params: &UpdateProductParams{
					ProducerID:      "producer-id",
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
					Price:             400,
					Cost:              300,
					StorageMethodType: entity.StorageMethodTypeNormal,
					DeliveryType:      entity.DeliveryTypeNormal,
					Box60Rate:         50,
					Box80Rate:         40,
					Box100Rate:        30,
					OriginPrefecture:  codes.PrefectureValues["shiga"],
					OriginCity:        "彦根市",
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				productID: "product-id",
				params:    &UpdateProductParams{},
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

			err := delete(ctx, productTable)
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
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				product := testProduct("product-id", "type-id", "category-id", "producer-id", []string{"product-id"}, now())
				err = db.DB.Create(&product).Error
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
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				productID: "product-id",
				set:       func(media entity.MultiProductMedia) bool { return false },
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "media is non existent",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				product := testProduct("product-id", "type-id", "category-id", "producer-id", []string{"tag-id"}, now())
				err = db.DB.Create(&product).Error
				require.NoError(t, err)
			},
			args: args{
				productID: "product-id",
				set:       func(media entity.MultiProductMedia) bool { return false },
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

			err := delete(ctx, productTable)
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
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				product := testProduct("product-id", "type-id", "category-id", "producer-id", []string{"tag-id"}, now())
				err = db.DB.Create(&product).Error
				require.NoError(t, err)
			},
			args: args{
				productID: "product-id",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				productID: "product-id",
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

			err := delete(ctx, productTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &product{db: db, now: now}
			err = db.Delete(ctx, tt.args.productID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testProduct(id, typeID, categoryID, producerID string, tagIDs []string, now time.Time) *entity.Product {
	p := &entity.Product{
		ID:              id,
		TypeID:          typeID,
		TagIDs:          tagIDs,
		ProducerID:      producerID,
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
		Price:             400,
		Cost:              300,
		StorageMethodType: entity.StorageMethodTypeNormal,
		DeliveryType:      entity.DeliveryTypeNormal,
		Box60Rate:         50,
		Box80Rate:         40,
		Box100Rate:        30,
		OriginPrefecture:  codes.PrefectureValues["shiga"],
		OriginCity:        "彦根市",
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	_ = p.FillJSON()
	return p
}

func fillIgnoreProductField(p *entity.Product, now time.Time) {
	if p == nil {
		return
	}
	_ = p.FillJSON()
}

func fillIgnoreProductsField(ps entity.Products, now time.Time) {
	for i := range ps {
		fillIgnoreProductField(ps[i], now)
	}
}
