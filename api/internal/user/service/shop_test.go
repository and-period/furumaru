package service

import (
	"context"
	"testing"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestMultiGetShops(t *testing.T) {
	t.Parallel()

	shops := entity.Shops{}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *MultiGetShopsInput
		expect    entity.Shops
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().MultiGet(ctx, []string{"shop-id"}).Return(shops, nil)
			},
			input: &MultiGetShopsInput{
				ShopIDs: []string{"shop-id"},
			},
			expect:    shops,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &MultiGetShopsInput{},
			expect:    nil,
			expectErr: ErrInvalidArgument,
		},
		{
			name: "failed to get shops",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().MultiGet(ctx, []string{"shop-id"}).Return(nil, errmock)
			},
			input: &MultiGetShopsInput{
				ShopIDs: []string{"shop-id"},
			},
			expect:    nil,
			expectErr: ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, service *userService) {
			actual, err := service.MultiGetShops(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetShop(t *testing.T) {
	t.Parallel()

	shop := &entity.Shop{}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *GetShopInput
		expect    *entity.Shop
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(ctx, "shop-id").Return(shop, nil)
			},
			input: &GetShopInput{
				ShopID: "shop-id",
			},
			expect:    shop,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &GetShopInput{},
			expect:    nil,
			expectErr: ErrInvalidArgument,
		},
		{
			name: "failed to get shop",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(ctx, "shop-id").Return(nil, errmock)
			},
			input: &GetShopInput{
				ShopID: "shop-id",
			},
			expect:    nil,
			expectErr: ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, service *userService) {
			actual, err := service.GetShop(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}
