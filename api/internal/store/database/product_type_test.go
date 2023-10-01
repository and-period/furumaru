package database

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func TestProductType(t *testing.T) {
	assert.NotNil(t, NewProductType(nil))
}

func TestProductType_List(t *testing.T) {
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
	productTypes[0] = testProductType("category-id01", "category-id01", "野菜", now())
	productTypes[1] = testProductType("category-id02", "category-id02", "果物", now())
	productTypes[2] = testProductType("category-id03", "category-id02", "水産物", now())
	err = db.DB.Create(&productTypes).Error
	require.NoError(t, err)

	type args struct {
		params *ListProductTypesParams
	}
	type want struct {
		productTypes entity.ProductTypes
		hasErr       bool
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
				params: &ListProductTypesParams{
					Name:       "物",
					CategoryID: "category-id02",
					Limit:      1,
					Offset:     1,
				},
			},
			want: want{
				productTypes: productTypes[2:],
				hasErr:       false,
			},
		},
		{
			name:  "success with asc",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &ListProductTypesParams{
					Orders: []*ListProductTypesOrder{
						{Key: entity.ProductTypeOrderByName, OrderByASC: true},
					},
				},
			},
			want: want{
				productTypes: productTypes,
				hasErr:       false,
			},
		},
		{
			name:  "success with desc",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &ListProductTypesParams{
					Orders: []*ListProductTypesOrder{
						{Key: entity.ProductTypeOrderByName, OrderByASC: false},
					},
				},
			},
			want: want{
				productTypes: productTypes,
				hasErr:       false,
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

			db := &productType{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.productTypes, actual)
		})
	}
}

func TestProductType_Count(t *testing.T) {
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
	productTypes[0] = testProductType("category-id01", "category-id01", "野菜", now())
	productTypes[1] = testProductType("category-id02", "category-id02", "果物", now())
	productTypes[2] = testProductType("category-id03", "category-id02", "水産物", now())
	err = db.DB.Create(&productTypes).Error
	require.NoError(t, err)

	type args struct {
		params *ListProductTypesParams
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
				params: &ListProductTypesParams{
					Name:       "物",
					CategoryID: "category-id02",
				},
			},
			want: want{
				total:  2,
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

			db := &productType{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestProductType_MultiGet(t *testing.T) {
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
	productTypes[0] = testProductType("category-id01", "category-id01", "野菜", now())
	productTypes[1] = testProductType("category-id02", "category-id02", "果物", now())
	productTypes[2] = testProductType("category-id03", "category-id02", "水産物", now())
	err = db.DB.Create(&productTypes).Error
	require.NoError(t, err)

	type args struct {
		productTypeIDs []string
	}
	type want struct {
		productTypes entity.ProductTypes
		hasErr       bool
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
				productTypeIDs: []string{"category-id01", "category-id02"},
			},
			want: want{
				productTypes: productTypes[:2],
				hasErr:       false,
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

			db := &productType{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.productTypeIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.productTypes, actual)
		})
	}
}

func TestProductType_Get(t *testing.T) {
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

	c := testCategory("category-id", "野菜", now())
	err = db.DB.Create(&c).Error
	require.NoError(t, err)
	p := testProductType("type-id", "category-id", "野菜", now())
	err = db.DB.Create(&p).Error
	require.NoError(t, err)

	type args struct {
		productTypeID string
	}
	type want struct {
		productType *entity.ProductType
		hasErr      bool
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
				productTypeID: "type-id",
			},
			want: want{
				productType: p,
				hasErr:      false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				productTypeID: "other-id",
			},
			want: want{
				productType: nil,
				hasErr:      true,
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

			db := &productType{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.productTypeID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.productType, actual)
		})
	}
}

func TestProductType_Create(t *testing.T) {
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

	type args struct {
		productType *entity.ProductType
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
				category := testCategory("category-id", "野菜", now())
				err := db.DB.Create(&category).Error
				require.NoError(t, err)
			},
			args: args{
				productType: testProductType("product-id", "category-id", "じゃがいも", now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found parant record",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				productType: testProductType("product-id", "category-id", "じゃがいも", now()),
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				category := testCategory("category-id", "野菜", now())
				err := db.DB.Create(&category).Error
				require.NoError(t, err)
				productType := testProductType("product-id", "category-id", "じゃがいも", now())
				err = db.DB.Create(&productType).Error
				require.NoError(t, err)
			},
			args: args{
				productType: testProductType("product-id", "category-id", "じゃがいも", now()),
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

			err := delete(ctx, productTypeTable, categoryTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &productType{db: db, now: now}
			err = db.Create(ctx, tt.args.productType)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestProductType_Update(t *testing.T) {
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

	type args struct {
		productTypeID string
		name          string
		iconURL       string
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
				category := testCategory("category-id", "野菜", now())
				err := db.DB.Create(&category).Error
				require.NoError(t, err)
				productType := testProductType("product-id", "category-id", "じゃがいも", now())
				err = db.DB.Create(&productType).Error
				require.NoError(t, err)
			},
			args: args{
				productTypeID: "product-id",
				name:          "さつまいも",
				iconURL:       "https://and-period.jp/icon.png",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				productTypeID: "product-id",
				name:          "さつまいも",
				iconURL:       "https://and-period.jp/icon.png",
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

			err := delete(ctx, productTypeTable, categoryTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &productType{db: db, now: now}
			err = db.Update(ctx, tt.args.productTypeID, tt.args.name, tt.args.iconURL)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestProductType_UpdateIcons(t *testing.T) {
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

	type args struct {
		productTypeID string
		icons         common.Images
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
				category := testCategory("category-id", "野菜", now())
				err := db.DB.Create(&category).Error
				require.NoError(t, err)
				productType := testProductType("product-id", "category-id", "じゃがいも", now())
				err = db.DB.Create(&productType).Error
				require.NoError(t, err)
			},
			args: args{
				productTypeID: "product-id",
				icons: common.Images{
					{
						Size: common.ImageSizeSmall,
						URL:  "https://and-period.jp/icon_240.png",
					},
					{
						Size: common.ImageSizeMedium,
						URL:  "https://and-period.jp/icon_675.png",
					},
					{
						Size: common.ImageSizeLarge,
						URL:  "https://and-period.jp/icon_900.png",
					},
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				productTypeID: "product-id",
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to empty icon url",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				category := testCategory("category-id", "野菜", now())
				err := db.DB.Create(&category).Error
				require.NoError(t, err)
				productType := testProductType("product-id", "category-id", "じゃがいも", now())
				productType.IconURL = ""
				err = db.DB.Create(&productType).Error
				require.NoError(t, err)
			},
			args: args{
				productTypeID: "product-id",
				icons: common.Images{
					{
						Size: common.ImageSizeSmall,
						URL:  "https://and-period.jp/icon_240.png",
					},
					{
						Size: common.ImageSizeMedium,
						URL:  "https://and-period.jp/icon_675.png",
					},
					{
						Size: common.ImageSizeLarge,
						URL:  "https://and-period.jp/icon_900.png",
					},
				},
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

			err := delete(ctx, productTypeTable, categoryTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &productType{db: db, now: now}
			err = db.UpdateIcons(ctx, tt.args.productTypeID, tt.args.icons)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestProductType_Delete(t *testing.T) {
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

	type args struct {
		productTypeID string
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
				category := testCategory("category-id", "野菜", now())
				err := db.DB.Create(&category).Error
				require.NoError(t, err)
				productType := testProductType("product-id", "category-id", "じゃがいも", now())
				err = db.DB.Create(&productType).Error
				require.NoError(t, err)
			},
			args: args{
				productTypeID: "product-id",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				productTypeID: "product-id",
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

			err := delete(ctx, productTypeTable, categoryTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &productType{db: db, now: now}
			err = db.Delete(ctx, tt.args.productTypeID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testProductType(id, categoryID, name string, now time.Time) *entity.ProductType {
	t := &entity.ProductType{
		ID:         id,
		Name:       name,
		IconURL:    "https://and-period.jp/icon.png",
		Icons:      common.Images{},
		CategoryID: categoryID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	fillProductTypeJSON(t)
	return t
}

func fillProductTypeJSON(t *entity.ProductType) {
	icons, _ := json.Marshal(t.Icons)
	t.IconsJSON = datatypes.JSON(icons)
}
