package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductType(t *testing.T) {
	assert.NotNil(t, NewProductType(nil))
}

func TestProductType_List(t *testing.T) {
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

	_ = m.dbDelete(ctx, productTypeTable, categoryTable)
	categories := make(entity.Categories, 2)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	err = m.db.DB.Create(&categories).Error
	require.NoError(t, err)
	productTypes := make(entity.ProductTypes, 3)
	productTypes[0] = testProductType("category-id01", "category-id01", "野菜", now())
	productTypes[1] = testProductType("category-id02", "category-id02", "果物", now())
	productTypes[2] = testProductType("category-id03", "category-id02", "水産物", now())
	err = m.db.DB.Create(&productTypes).Error
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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, m)

			db := &productType{db: m.db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreProductTypesField(actual, now())
			assert.Equal(t, tt.want.productTypes, actual)
		})
	}
}

func TestProductType_Create(t *testing.T) {
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

	_ = m.dbDelete(ctx, productTypeTable, categoryTable)

	type args struct {
		productType *entity.ProductType
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
				category := testCategory("category-id", "野菜", now())
				err := m.db.DB.Create(&category).Error
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
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				productType: testProductType("product-id", "category-id", "じゃがいも", now()),
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				category := testCategory("category-id", "野菜", now())
				err := m.db.DB.Create(&category).Error
				require.NoError(t, err)
				productType := testProductType("product-id", "category-id", "じゃがいも", now())
				err = m.db.DB.Create(&productType).Error
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

			err := m.dbDelete(ctx, productTypeTable, categoryTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &productType{db: m.db, now: now}
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

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, productTypeTable, categoryTable)

	type args struct {
		productTypeID string
		name          string
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
				category := testCategory("category-id", "野菜", now())
				err := m.db.DB.Create(&category).Error
				require.NoError(t, err)
				productType := testProductType("product-id", "category-id", "じゃがいも", now())
				err = m.db.DB.Create(&productType).Error
				require.NoError(t, err)
			},
			args: args{
				productTypeID: "product-id",
				name:          "さつまいも",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				productTypeID: "product-id",
				name:          "さつまいも",
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

			err := m.dbDelete(ctx, productTypeTable, categoryTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &productType{db: m.db, now: now}
			err = db.Update(ctx, tt.args.productTypeID, tt.args.name)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestProductType_Delete(t *testing.T) {
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

	_ = m.dbDelete(ctx, productTypeTable, categoryTable)

	type args struct {
		productTypeID string
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
				category := testCategory("category-id", "野菜", now())
				err := m.db.DB.Create(&category).Error
				require.NoError(t, err)
				productType := testProductType("product-id", "category-id", "じゃがいも", now())
				err = m.db.DB.Create(&productType).Error
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
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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

			err := m.dbDelete(ctx, productTypeTable, categoryTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &productType{db: m.db, now: now}
			err = db.Delete(ctx, tt.args.productTypeID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testProductType(id, categoryID, name string, now time.Time) *entity.ProductType {
	return &entity.ProductType{
		ID:         id,
		Name:       name,
		CategoryID: categoryID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func fillIgnoreProductTypeField(t *entity.ProductType, now time.Time) {
	if t == nil {
		return
	}
	t.CreatedAt = now
	t.UpdatedAt = now
}

func fillIgnoreProductTypesField(ts entity.ProductTypes, now time.Time) {
	for i := range ts {
		fillIgnoreProductTypeField(ts[i], now)
	}
}
