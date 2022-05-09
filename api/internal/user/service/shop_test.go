package service

import (
	"context"
	"testing"

	"github.com/and-period/marche/api/internal/user/database"
	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListShops(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	shops := entity.Shops{
		{
			ID:            "shop-id",
			CognitoID:     "cognito-id",
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "すたっふ",
			Email:         "test-shop@and-period.jp",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}
	params := &database.ListShopsParams{
		Limit:  30,
		Offset: 0,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *ListShopsInput
		expect    entity.Shops
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().List(ctx, params).Return(shops, nil)
			},
			input: &ListShopsInput{
				Limit:  30,
				Offset: 0,
			},
			expect:    shops,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &ListShopsInput{},
			expect:    nil,
			expectErr: ErrInvalidArgument,
		},
		{
			name: "failed to get shops",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().List(ctx, params).Return(nil, errmock)
			},
			input: &ListShopsInput{
				Limit:  30,
				Offset: 0,
			},
			expect:    nil,
			expectErr: ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *userService) {
			actual, err := service.ListShops(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestMultiGetShops(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	shops := entity.Shops{
		{
			ID:            "shop-id",
			CognitoID:     "cognito-id",
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "すたっふ",
			Email:         "test-shop@and-period.jp",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

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
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *userService) {
			actual, err := service.MultiGetShops(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetShop(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	shop := &entity.Shop{
		ID:            "shop-id",
		CognitoID:     "cognito-id",
		Lastname:      "&.",
		Firstname:     "スタッフ",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "すたっふ",
		Email:         "test-shop@and-period.jp",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

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
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *userService) {
			actual, err := service.GetShop(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateShop(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *CreateShopInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.shopAuth.EXPECT().AdminCreateUser(ctx, gomock.Any()).Return(nil)
			},
			input: &CreateShopInput{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				Email:         "test-shop@and-period.jp",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &CreateShopInput{},
			expectErr: ErrInvalidArgument,
		},
		{
			name: "failed to create shop",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Create(ctx, gomock.Any()).Return(errmock)
			},
			input: &CreateShopInput{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				Email:         "test-shop@and-period.jp",
			},
			expectErr: ErrInternal,
		},
		{
			name: "failed to create auth shop",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.shopAuth.EXPECT().AdminCreateUser(ctx, gomock.Any()).Return(errmock)
			},
			input: &CreateShopInput{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				Email:         "test-shop@and-period.jp",
			},
			expectErr: ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *userService) {
			_, err := service.CreateShop(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
