package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListShops(t *testing.T) {
	t.Parallel()

	now := time.Now()
	params := &database.ListShopsParams{
		CoordinatorIDs: []string{"coordinator-id"},
		ProducerIDs:    []string{"producer-id"},
		Limit:          20,
		Offset:         0,
	}
	shops := entity.Shops{
		{
			ID:             "shop-id",
			Name:           "テスト店舗",
			CoordinatorID:  "coordinator-id",
			ProducerIDs:    []string{"producer-id"},
			ProductTypeIDs: []string{"product-type-id"},
			Activated:      true,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListShopsInput
		expect      entity.Shops
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().List(gomock.Any(), params).Return(shops, nil)
				mocks.db.Shop.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListShopsInput{
				CoordinatorIDs: []string{"coordinator-id"},
				ProducerIDs:    []string{"producer-id"},
				Limit:          20,
			},
			expect:      shops,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.ListShopsInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list shops",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Shop.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListShopsInput{
				CoordinatorIDs: []string{"coordinator-id"},
				ProducerIDs:    []string{"producer-id"},
				Limit:          20,
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to count shops",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().List(gomock.Any(), params).Return(shops, nil)
				mocks.db.Shop.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &store.ListShopsInput{
				CoordinatorIDs: []string{"coordinator-id"},
				ProducerIDs:    []string{"producer-id"},
				Limit:          20,
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListShops(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestListShopProducers(t *testing.T) {
	t.Parallel()

	params := &database.ListShopProducersParams{
		CoordinatorID: "coordinator-id",
		Limit:         20,
		Offset:        0,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.ListShopProducersInput
		expect    []string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().ListProducers(ctx, params).Return([]string{"producer-id"}, nil)
			},
			input: &store.ListShopProducersInput{
				CoordinatorID: "coordinator-id",
				Limit:         20,
			},
			expect:    []string{"producer-id"},
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.ListShopProducersInput{
				Limit: -1,
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list producers",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().ListProducers(ctx, params).Return(nil, assert.AnError)
			},
			input: &store.ListShopProducersInput{
				CoordinatorID: "coordinator-id",
				Limit:         20,
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.ListShopProducers(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestMultiGetShops(t *testing.T) {
	t.Parallel()

	now := time.Now()
	shops := entity.Shops{
		{
			ID:             "shop-id",
			Name:           "テスト店舗",
			CoordinatorID:  "coordinator-id",
			ProducerIDs:    []string{"producer-id"},
			ProductTypeIDs: []string{"product-type-id"},
			Activated:      true,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetShopsInput
		expect    entity.Shops
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().MultiGet(ctx, []string{"shop-id"}).Return(shops, nil)
			},
			input: &store.MultiGetShopsInput{
				ShopIDs: []string{"shop-id"},
			},
			expect:    shops,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetShopsInput{
				ShopIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multi get shops",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().MultiGet(ctx, []string{"shop-id"}).Return(nil, assert.AnError)
			},
			input: &store.MultiGetShopsInput{
				ShopIDs: []string{"shop-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetShops(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestGetShop(t *testing.T) {
	t.Parallel()

	now := time.Now()
	shop := &entity.Shop{
		ID:            "shop-id",
		CoordinatorID: "coordinator-id",
		ProducerIDs:   []string{"producer-id"},
		Name:          "テスト店舗",
		Activated:     true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetShopInput
		expect    *entity.Shop
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(ctx, "shop-id").Return(shop, nil)
			},
			input: &store.GetShopInput{
				ShopID: "shop-id",
			},
			expect:    shop,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetShopInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get shop",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(ctx, "shop-id").Return(nil, assert.AnError)
			},
			input: &store.GetShopInput{
				ShopID: "shop-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetShop(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestGetShopByCoordinatorID(t *testing.T) {
	t.Parallel()

	now := time.Now()
	shop := &entity.Shop{
		ID:            "shop-id",
		CoordinatorID: "coordinator-id",
		ProducerIDs:   []string{"producer-id"},
		Name:          "テスト店舗",
		Activated:     true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetShopByCoordinatorIDInput
		expect    *entity.Shop
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().GetByCoordinatorID(ctx, "coordinator-id").Return(shop, nil)
			},
			input: &store.GetShopByCoordinatorIDInput{
				CoordinatorID: "coordinator-id",
			},
			expect:    shop,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetShopByCoordinatorIDInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get shop",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().GetByCoordinatorID(ctx, "coordinator-id").Return(nil, assert.AnError)
			},
			input: &store.GetShopByCoordinatorIDInput{
				CoordinatorID: "coordinator-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetShopByCoordinatorID(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateShop(t *testing.T) {
	t.Parallel()

	productTypes := entity.ProductTypes{
		{
			ID:         "product-type-id",
			CategoryID: "category-id",
			Name:       "テスト商品種別",
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateShopInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().MultiGet(ctx, []string{"product-type-id"}).Return(productTypes, nil)
				mocks.db.Shop.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, shop *entity.Shop) error {
						expect := &entity.Shop{
							ID:             shop.ID, // ignore
							CoordinatorID:  "coordinator-id",
							Name:           "テスト店舗",
							ProductTypeIDs: []string{"product-type-id"},
							BusinessDays:   []time.Weekday{time.Monday, time.Tuesday},
							Activated:      true,
						}
						assert.Equal(t, expect, shop)
						return nil
					})
			},
			input: &store.CreateShopInput{
				CoordinatorID:  "coordinator-id",
				Name:           "テスト店舗",
				ProductTypeIDs: []string{"product-type-id"},
				BusinessDays:   []time.Weekday{time.Monday, time.Tuesday},
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateShopInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get product types",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().MultiGet(ctx, []string{"product-type-id"}).Return(nil, assert.AnError)
			},
			input: &store.CreateShopInput{
				CoordinatorID:  "coordinator-id",
				Name:           "テスト店舗",
				ProductTypeIDs: []string{"product-type-id"},
				BusinessDays:   []time.Weekday{time.Monday, time.Tuesday},
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "contains invalid product type ids",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().MultiGet(ctx, []string{"product-type-id"}).Return(entity.ProductTypes{}, nil)
			},
			input: &store.CreateShopInput{
				CoordinatorID:  "coordinator-id",
				Name:           "テスト店舗",
				ProductTypeIDs: []string{"product-type-id"},
				BusinessDays:   []time.Weekday{time.Monday, time.Tuesday},
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create shop",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().MultiGet(ctx, []string{"product-type-id"}).Return(productTypes, nil)
				mocks.db.Shop.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateShopInput{
				CoordinatorID:  "coordinator-id",
				Name:           "テスト店舗",
				ProductTypeIDs: []string{"product-type-id"},
				BusinessDays:   []time.Weekday{time.Monday, time.Tuesday},
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateShop(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateShop(t *testing.T) {
	t.Parallel()

	productTypes := entity.ProductTypes{
		{
			ID:         "product-type-id",
			CategoryID: "category-id",
			Name:       "テスト商品種別",
		},
	}
	params := &database.UpdateShopParams{
		Name:           "テスト店舗",
		ProductTypeIDs: []string{"product-type-id"},
		BusinessDays:   []time.Weekday{time.Monday, time.Tuesday},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateShopInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().MultiGet(ctx, []string{"product-type-id"}).Return(productTypes, nil)
				mocks.db.Shop.EXPECT().Update(ctx, "shop-id", params).Return(nil)
			},
			input: &store.UpdateShopInput{
				ShopID:         "shop-id",
				Name:           "テスト店舗",
				ProductTypeIDs: []string{"product-type-id"},
				BusinessDays:   []time.Weekday{time.Monday, time.Tuesday},
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateShopInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get product types",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().MultiGet(ctx, []string{"product-type-id"}).Return(nil, assert.AnError)
			},
			input: &store.UpdateShopInput{
				ShopID:         "shop-id",
				Name:           "テスト店舗",
				ProductTypeIDs: []string{"product-type-id"},
				BusinessDays:   []time.Weekday{time.Monday, time.Tuesday},
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "contains invalid product type ids",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().MultiGet(ctx, []string{"product-type-id"}).Return(entity.ProductTypes{}, nil)
			},
			input: &store.UpdateShopInput{
				ShopID:         "shop-id",
				Name:           "テスト店舗",
				ProductTypeIDs: []string{"product-type-id"},
				BusinessDays:   []time.Weekday{time.Monday, time.Tuesday},
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update shop",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ProductType.EXPECT().MultiGet(ctx, []string{"product-type-id"}).Return(productTypes, nil)
				mocks.db.Shop.EXPECT().Update(ctx, "shop-id", params).Return(assert.AnError)
			},
			input: &store.UpdateShopInput{
				ShopID:         "shop-id",
				Name:           "テスト店舗",
				ProductTypeIDs: []string{"product-type-id"},
				BusinessDays:   []time.Weekday{time.Monday, time.Tuesday},
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateShop(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteShop(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteShopInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Delete(ctx, "shop-id").Return(nil)
			},
			input: &store.DeleteShopInput{
				ShopID: "shop-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteShopInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete shop",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Delete(ctx, "shop-id").Return(assert.AnError)
			},
			input: &store.DeleteShopInput{
				ShopID: "shop-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteShop(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestRelateShopProducer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.RelateShopProducerInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().RelateProducer(ctx, "shop-id", "producer-id").Return(nil)
			},
			input: &store.RelateShopProducerInput{
				ShopID:     "shop-id",
				ProducerID: "producer-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.RelateShopProducerInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to relate shop producer",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().RelateProducer(ctx, "shop-id", "producer-id").Return(assert.AnError)
			},
			input: &store.RelateShopProducerInput{
				ShopID:     "shop-id",
				ProducerID: "producer-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.RelateShopProducer(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUnrelateShopProducer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UnrelateShopProducerInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().UnrelateProducer(ctx, "shop-id", "producer-id").Return(nil)
			},
			input: &store.UnrelateShopProducerInput{
				ShopID:     "shop-id",
				ProducerID: "producer-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UnrelateShopProducerInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to unrelate shop producer",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().UnrelateProducer(ctx, "shop-id", "producer-id").Return(assert.AnError)
			},
			input: &store.UnrelateShopProducerInput{
				ShopID:     "shop-id",
				ProducerID: "producer-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UnrelateShopProducer(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
