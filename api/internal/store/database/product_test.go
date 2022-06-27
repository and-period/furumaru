package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
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

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, productTable, productTypeTable, categoryTable)
	categories := make(entity.Categories, 2)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	err = m.db.DB.Create(&categories).Error
	require.NoError(t, err)
	productTypes := make(entity.ProductTypes, 2)
	productTypes[0] = testProductType("type-id01", "category-id01", "野菜", now())
	productTypes[1] = testProductType("type-id02", "category-id02", "果物", now())
	err = m.db.DB.Create(&productTypes).Error
	require.NoError(t, err)
	products := make(entity.Products, 3)
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "producer-id", "coordinator-id", now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "producer-id", "coordinator-id", now())
	products[2] = testProduct("product-id03", "type-id02", "category-id02", "producer-id", "coordinator-id", now())
	err = m.db.DB.Create(&products).Error
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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				params: &ListProductsParams{
					Name:       "いも",
					ProducerID: "producer-id",
					CreatedBy:  "coordinator-id",
					Limit:      2,
					Offset:     1,
				},
			},
			want: want{
				products: products[1:],
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

			tt.setup(ctx, t, m)

			db := &product{db: m.db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			fillIgnoreProductsField(actual, now())
			assert.Equal(t, tt.want.products, actual)
		})
	}
}

func TestProduct_Get(t *testing.T) {
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

	db := &product{db: m.db, now: now}
	actual, err := db.Get(ctx, "product-id")
	assert.Nil(t, actual)
	assert.ErrorIs(t, err, exception.ErrNotImplemented)
}

func TestProduct_Create(t *testing.T) {
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

	db := &product{db: m.db, now: now}
	err = db.Create(ctx, &entity.Product{})
	assert.ErrorIs(t, err, exception.ErrNotImplemented)
}

func TestProduct_Update(t *testing.T) {
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

	db := &product{db: m.db, now: now}
	err = db.Update(ctx, &entity.Product{})
	assert.ErrorIs(t, err, exception.ErrNotImplemented)
}

func TestProduct_Delete(t *testing.T) {
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

	db := &product{db: m.db, now: now}
	err = db.Delete(ctx, "product-id")
	assert.ErrorIs(t, err, exception.ErrNotImplemented)
}

func testProduct(id, typeID, categoryID, producerID, coordinatorID string, now time.Time) *entity.Product {
	p := &entity.Product{
		ID:              id,
		TypeID:          typeID,
		CategoryID:      categoryID,
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
		Price:            400,
		DeliveryType:     entity.DeliveryTypeNormal,
		Box60Rate:        50,
		Box80Rate:        40,
		Box100Rate:       30,
		OriginPrefecture: "滋賀県",
		OriginCity:       "彦根市",
		CreatedAt:        now,
		UpdatedAt:        now,
		CreatedBy:        coordinatorID,
		UpdatedBy:        coordinatorID,
	}
	_ = p.FillJSON()
	return p
}

func fillIgnoreProductField(p *entity.Product, now time.Time) {
	if p == nil {
		return
	}
	_ = p.FillJSON()
	p.CreatedAt = now
	p.UpdatedAt = now
}

func fillIgnoreProductsField(ps entity.Products, now time.Time) {
	for i := range ps {
		fillIgnoreProductField(ps[i], now)
	}
}
