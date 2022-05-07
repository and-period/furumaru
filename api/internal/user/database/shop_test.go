package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShop(t *testing.T) {
	assert.NotNil(t, NewShop(nil))
}

func TestShop_MultiGet(t *testing.T) {
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

	_ = m.dbDelete(ctx, shopTable)
	shops := make(entity.Shops, 2)
	shops[0] = testShop("shop-id01", "test-shop01@and-period.jp", now())
	shops[1] = testShop("shop-id02", "test-shop02@and-period.jp", now())
	err = m.db.DB.Create(&shops).Error
	require.NoError(t, err)

	type args struct {
		shopIDs []string
	}
	type want struct {
		shops  entity.Shops
		hasErr bool
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
				shopIDs: []string{"shop-id01", "shop-id02"},
			},
			want: want{
				shops:  shops,
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

			tt.setup(ctx, t, m)

			db := &shop{db: m.db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.shopIDs)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreShopsField(actual, now())
			assert.ElementsMatch(t, tt.want.shops, actual)
		})
	}
}

func TestShop_Get(t *testing.T) {
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

	_ = m.dbDelete(ctx, shopTable)
	s := testShop("shop-id", "test-shop@and-period.jp", now())
	err = m.db.DB.Create(&s).Error
	require.NoError(t, err)

	type args struct {
		shopID string
	}
	type want struct {
		shop   *entity.Shop
		hasErr bool
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
				shopID: "shop-id",
			},
			want: want{
				shop:   s,
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				shopID: "",
			},
			want: want{
				shop:   nil,
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

			tt.setup(ctx, t, m)

			db := &shop{db: m.db, now: now}
			actual, err := db.Get(ctx, tt.args.shopID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreShopField(actual, now())
			assert.Equal(t, tt.want.shop, actual)
		})
	}
}

func testShop(id, email string, now time.Time) *entity.Shop {
	return &entity.Shop{
		ID:            id,
		CognitoID:     id,
		Lastname:      "&.",
		Firstname:     "スタッフ",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "すたっふ",
		Email:         email,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func fillIgnoreShopField(s *entity.Shop, now time.Time) {
	if s == nil {
		return
	}
	s.CreatedAt = now
	s.UpdatedAt = now
}

func fillIgnoreShopsField(ss entity.Shops, now time.Time) {
	for i := range ss {
		fillIgnoreShopField(ss[i], now)
	}
}
