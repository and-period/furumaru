package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetCart(t *testing.T) {
	t.Parallel()

	now := time.Now()
	product := func(productID string, invenroty int64) *entity.Product {
		return &entity.Product{
			ID:            productID,
			CoordinatorID: "coordinator-id",
			DeliveryType:  entity.DeliveryTypeNormal,
			Inventory:     invenroty,
			Public:        true,
			Status:        entity.ProductStatusForSale,
			Weight:        500,
			WeightUnit:    entity.WeightUnitGram,
			Box60Rate:     80,
			Box80Rate:     50,
			Box100Rate:    30,
		}
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetCartInput
		expect    *entity.Cart
		expectErr error
	}{
		{
			name: "success first time",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().Get(ctx, &entity.Cart{SessionID: "session-id"}).Return(dynamodb.ErrNotFound)
				cart := &entity.Cart{
					SessionID: "session-id",
					Baskets:   []*entity.CartBasket{},
					ExpiredAt: now.Add(defaultCartTTL),
					CreatedAt: now,
					UpdatedAt: now,
				}
				mocks.cache.EXPECT().Insert(ctx, cart).Return(nil)
			},
			input: &store.GetCartInput{
				SessionID: "session-id",
			},
			expect: &entity.Cart{
				SessionID: "session-id",
				Baskets:   []*entity.CartBasket{},
				ExpiredAt: now.Add(defaultCartTTL),
				CreatedAt: now,
				UpdatedAt: now,
			},
			expectErr: nil,
		},
		{
			name: "success after second times with refresh",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.Cart{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, cart *entity.Cart) error {
						cart.Baskets = []*entity.CartBasket{{
							BoxNumber: 1,
							BoxType:   entity.ShippingTypeNormal,
							BoxSize:   entity.ShippingSize80,
							BoxRate:   100,
							Items: entity.CartItems{
								{ProductID: "product-id", Quantity: 2},
							},
							CoordinatorID: "coordinator-id",
						}}
						cart.ExpiredAt = now.AddDate(0, 0, 12)
						cart.CreatedAt = now.AddDate(0, 0, -2)
						cart.UpdatedAt = now.AddDate(0, 0, -2)
						return nil
					})
				products := entity.Products{product("product-id", 1)}
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products, nil)
				cart := &entity.Cart{
					SessionID: "session-id",
					Baskets: []*entity.CartBasket{
						{
							BoxNumber: 1,
							BoxType:   entity.ShippingTypeNormal,
							BoxSize:   entity.ShippingSize60,
							BoxRate:   80,
							Items: entity.CartItems{
								{ProductID: "product-id", Quantity: 1},
							},
							CoordinatorID: "coordinator-id",
						},
					},
					ExpiredAt: now.Add(defaultCartTTL),
					CreatedAt: now.AddDate(0, 0, -2),
					UpdatedAt: now,
				}
				mocks.cache.EXPECT().Insert(ctx, cart).Return(nil)
			},
			input: &store.GetCartInput{
				SessionID: "session-id",
			},
			expect: &entity.Cart{
				SessionID: "session-id",
				Baskets: []*entity.CartBasket{
					{
						BoxNumber: 1,
						BoxType:   entity.ShippingTypeNormal,
						BoxSize:   entity.ShippingSize60,
						BoxRate:   80,
						Items: entity.CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "coordinator-id",
					},
				},
				ExpiredAt: now.Add(defaultCartTTL),
				CreatedAt: now.AddDate(0, 0, -2),
				UpdatedAt: now,
			},
			expectErr: nil,
		},
		{
			name: "success after second times without refresh",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.Cart{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, cart *entity.Cart) error {
						cart.Baskets = []*entity.CartBasket{{
							BoxNumber: 1,
							BoxType:   entity.ShippingTypeNormal,
							BoxSize:   entity.ShippingSize80,
							BoxRate:   100,
							Items: entity.CartItems{
								{ProductID: "product-id", Quantity: 2},
							},
							CoordinatorID: "coordinator-id",
						}}
						cart.ExpiredAt = now.Add(defaultCartTTL)
						cart.CreatedAt = now
						cart.UpdatedAt = now
						return nil
					})
			},
			input: &store.GetCartInput{
				SessionID: "session-id",
			},
			expect: &entity.Cart{
				SessionID: "session-id",
				Baskets: []*entity.CartBasket{
					{
						BoxNumber: 1,
						BoxType:   entity.ShippingTypeNormal,
						BoxSize:   entity.ShippingSize80,
						BoxRate:   100,
						Items: entity.CartItems{
							{ProductID: "product-id", Quantity: 2},
						},
						CoordinatorID: "coordinator-id",
					},
				},
				ExpiredAt: now.Add(defaultCartTTL),
				CreatedAt: now,
				UpdatedAt: now,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetCartInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get cart",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().Get(ctx, &entity.Cart{SessionID: "session-id"}).Return(assert.AnError)
			},
			input: &store.GetCartInput{
				SessionID: "session-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create cart",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().Get(ctx, &entity.Cart{SessionID: "session-id"}).Return(dynamodb.ErrNotFound)
				cart := &entity.Cart{
					SessionID: "session-id",
					Baskets:   []*entity.CartBasket{},
					ExpiredAt: now.Add(defaultCartTTL),
					CreatedAt: now,
					UpdatedAt: now,
				}
				mocks.cache.EXPECT().Insert(ctx, cart).Return(assert.AnError)
			},
			input: &store.GetCartInput{
				SessionID: "session-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to multi get products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.Cart{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, cart *entity.Cart) error {
						cart.Baskets = []*entity.CartBasket{}
						cart.ExpiredAt = now.AddDate(0, 0, 12)
						cart.CreatedAt = now.AddDate(0, 0, -2)
						cart.UpdatedAt = now.AddDate(0, 0, -2)
						return nil
					})
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{}).Return(nil, assert.AnError)
			},
			input: &store.GetCartInput{
				SessionID: "session-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to update cart",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.Cart{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, cart *entity.Cart) error {
						cart.Baskets = []*entity.CartBasket{{
							BoxNumber: 1,
							BoxType:   entity.ShippingTypeNormal,
							BoxSize:   entity.ShippingSize80,
							BoxRate:   100,
							Items: entity.CartItems{
								{ProductID: "product-id", Quantity: 2},
							},
							CoordinatorID: "coordinator-id",
						}}
						cart.ExpiredAt = now.AddDate(0, 0, 12)
						cart.CreatedAt = now.AddDate(0, 0, -2)
						cart.UpdatedAt = now.AddDate(0, 0, -2)
						return nil
					})
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(entity.Products{}, nil)
				cart := &entity.Cart{
					SessionID: "session-id",
					Baskets:   []*entity.CartBasket{},
					ExpiredAt: now.Add(defaultCartTTL),
					CreatedAt: now.AddDate(0, 0, -2),
					UpdatedAt: now,
				}
				mocks.cache.EXPECT().Insert(ctx, cart).Return(assert.AnError)
			},
			input: &store.GetCartInput{
				SessionID: "session-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetCart(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}, withNow(now)))
	}
}

func TestCalcCart(t *testing.T) {
	t.Parallel()
	now := time.Now()
	shikoku := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int32, 0, 47-len(shikoku))
	for _, val := range codes.PrefectureValues {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := entity.ShippingRates{
		{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
		{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
	}
	cart := &entity.Cart{
		SessionID: "session-id",
		Baskets: entity.CartBaskets{{
			BoxNumber:     1,
			BoxType:       entity.ShippingTypeNormal,
			BoxSize:       entity.ShippingSize80,
			BoxRate:       80,
			Items:         entity.CartItems{{ProductID: "product-id", Quantity: 2}},
			CoordinatorID: "coordinator-id",
		}},
		ExpiredAt: now.Add(defaultCartTTL),
		CreatedAt: now,
		UpdatedAt: now,
	}
	cartmocks := func(mocks *mocks, sessionID string, cart *entity.Cart, err error) {
		fn := func(_ context.Context, in *entity.Cart) error {
			in.Baskets = cart.Baskets
			in.ExpiredAt = now.Add(defaultCartTTL)
			in.CreatedAt = now
			in.UpdatedAt = now
			return err
		}
		mocks.cache.EXPECT().Get(gomock.Any(), &entity.Cart{SessionID: sessionID}).DoAndReturn(fn)
	}
	shipping := &entity.Shipping{
		ShippingRevision: entity.ShippingRevision{
			ShippingID:        "coordinator-id",
			Box60Rates:        rates,
			Box60Frozen:       800,
			Box80Rates:        rates,
			Box80Frozen:       800,
			Box100Rates:       rates,
			Box100Frozen:      800,
			HasFreeShipping:   true,
			FreeShippingRates: 3000,
		},
		ID:            "coordinator-id",
		CoordinatorID: "coordinator-id",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	promotion := &entity.Promotion{
		ID:           "promotion-id",
		Title:        "プロモーションタイトル",
		Description:  "プロモーションの詳細です。",
		Public:       true,
		PublishedAt:  now.AddDate(0, -1, 0),
		DiscountType: entity.DiscountTypeRate,
		DiscountRate: 10,
		Code:         "code1234",
		CodeType:     entity.PromotionCodeTypeAlways,
		StartAt:      now.AddDate(0, -1, 0),
		EndAt:        now.AddDate(0, 1, 0),
	}
	products := func(inventory int64) entity.Products {
		return entity.Products{
			{
				ID:        "product-id",
				Name:      "じゃがいも",
				Inventory: inventory,
				Public:    true,
				Status:    entity.ProductStatusForSale,
				ProductRevision: entity.ProductRevision{
					ID:        1,
					ProductID: "product-id",
					Price:     500,
				},
			},
		}
	}
	tests := []struct {
		name          string
		setup         func(ctx context.Context, mocks *mocks)
		input         *store.CalcCartInput
		expectCart    *entity.Cart
		expectSummary *entity.OrderPaymentSummary
		expectErr     error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products(30), nil)
			},
			input: &store.CalcCartInput{
				SessionID:      "session-id",
				CoordinatorID:  "coordinator-id",
				BoxNumber:      0,
				PromotionCode:  "code1234",
				PrefectureCode: 13,
			},
			expectCart: cart,
			expectSummary: &entity.OrderPaymentSummary{
				Subtotal:    1000,
				Discount:    100,
				ShippingFee: 500,
				Tax:         140,
				TaxRate:     10,
				Total:       1540,
			},
			expectErr: nil,
		},
		{
			name: "success without shipping and promotion",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products(30), nil)
			},
			input: &store.CalcCartInput{
				SessionID:      "session-id",
				CoordinatorID:  "coordinator-id",
				BoxNumber:      0,
				PromotionCode:  "",
				PrefectureCode: 0,
			},
			expectCart: cart,
			expectSummary: &entity.OrderPaymentSummary{
				Subtotal:    1000,
				Discount:    0,
				ShippingFee: 0,
				Tax:         100,
				TaxRate:     10,
				Total:       1100,
			},
			expectErr: nil,
		},
		{
			name: "failed to get cart",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, assert.AnError)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
			},
			input: &store.CalcCartInput{
				SessionID:      "session-id",
				CoordinatorID:  "coordinator-id",
				BoxNumber:      0,
				PromotionCode:  "code1234",
				PrefectureCode: 13,
			},
			expectCart:    nil,
			expectSummary: nil,
			expectErr:     exception.ErrInternal,
		},
		{
			name: "failed to get shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(nil, assert.AnError)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
			},
			input: &store.CalcCartInput{
				SessionID:      "session-id",
				CoordinatorID:  "coordinator-id",
				BoxNumber:      0,
				PromotionCode:  "code1234",
				PrefectureCode: 13,
			},
			expectCart:    nil,
			expectSummary: nil,
			expectErr:     exception.ErrInternal,
		},
		{
			name: "failed to get promotion",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(nil, assert.AnError)
			},
			input: &store.CalcCartInput{
				SessionID:      "session-id",
				CoordinatorID:  "coordinator-id",
				BoxNumber:      0,
				PromotionCode:  "code1234",
				PrefectureCode: 13,
			},
			expectCart:    nil,
			expectSummary: nil,
			expectErr:     exception.ErrInternal,
		},
		{
			name: "failed to get products",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(entity.Products{}, assert.AnError)
			},
			input: &store.CalcCartInput{
				SessionID:      "session-id",
				CoordinatorID:  "coordinator-id",
				BoxNumber:      0,
				PromotionCode:  "code1234",
				PrefectureCode: 13,
			},
			expectCart:    nil,
			expectSummary: nil,
			expectErr:     exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			cart, summary, err := service.CalcCart(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expectCart, cart)
			assert.Equal(t, tt.expectSummary, summary)
		}, withNow(now)))
	}
}

func TestAddCartItem(t *testing.T) {
	t.Parallel()

	now := time.Now()
	product := func(productID string, invenroty int64) *entity.Product {
		return &entity.Product{
			ID:            productID,
			CoordinatorID: "coordinator-id",
			DeliveryType:  entity.DeliveryTypeNormal,
			Inventory:     invenroty,
			Public:        true,
			Status:        entity.ProductStatusForSale,
			Weight:        500,
			WeightUnit:    entity.WeightUnitGram,
			Box60Rate:     80,
			Box80Rate:     50,
			Box100Rate:    30,
		}
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.AddCartItemInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				product := product("product-id", 1)
				mocks.db.Product.EXPECT().Get(ctx, "product-id").Return(product, nil)
				mocks.cache.EXPECT().
					Get(ctx, &entity.Cart{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, cart *entity.Cart) error {
						cart.Baskets = entity.CartBaskets{}
						cart.ExpiredAt = now.AddDate(0, 0, 12)
						cart.CreatedAt = now.AddDate(0, 0, -2)
						cart.UpdatedAt = now.AddDate(0, 0, -2)
						return nil
					})
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(entity.Products{product}, nil)
				cart := &entity.Cart{
					SessionID: "session-id",
					Baskets: entity.CartBaskets{
						{
							BoxNumber: 1,
							BoxType:   entity.ShippingTypeNormal,
							BoxSize:   entity.ShippingSize60,
							BoxRate:   80,
							Items: entity.CartItems{
								{ProductID: "product-id", Quantity: 1},
							},
							CoordinatorID: "coordinator-id",
						},
					},
					ExpiredAt: now.Add(defaultCartTTL),
					CreatedAt: now.AddDate(0, 0, -2),
					UpdatedAt: now,
				}
				mocks.cache.EXPECT().Insert(ctx, cart).Return(nil)
			},
			input: &store.AddCartItemInput{
				SessionID: "session-id",
				ProductID: "product-id",
				Quantity:  1,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.AddCartItemInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get product",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().Get(ctx, "product-id").Return(nil, assert.AnError)
			},
			input: &store.AddCartItemInput{
				SessionID: "session-id",
				ProductID: "product-id",
				Quantity:  2,
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get cart",
			setup: func(ctx context.Context, mocks *mocks) {
				product := product("product-id", 1)
				mocks.db.Product.EXPECT().Get(ctx, "product-id").Return(product, nil)
				mocks.cache.EXPECT().Get(ctx, &entity.Cart{SessionID: "session-id"}).Return(assert.AnError)
			},
			input: &store.AddCartItemInput{
				SessionID: "session-id",
				ProductID: "product-id",
				Quantity:  1,
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "product is not published",
			setup: func(ctx context.Context, mocks *mocks) {
				product := product("product-id", 1)
				product.Public = false
				mocks.db.Product.EXPECT().Get(ctx, "product-id").Return(product, nil)
			},
			input: &store.AddCartItemInput{
				SessionID: "session-id",
				ProductID: "product-id",
				Quantity:  1,
			},
			expectErr: exception.ErrForbidden,
		},
		{
			name: "insufficient product stock",
			setup: func(ctx context.Context, mocks *mocks) {
				product := product("product-id", 1)
				mocks.db.Product.EXPECT().Get(ctx, "product-id").Return(product, nil)
				mocks.cache.EXPECT().
					Get(ctx, &entity.Cart{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, cart *entity.Cart) error {
						cart.Baskets = entity.CartBaskets{}
						cart.ExpiredAt = now.AddDate(0, 0, 12)
						cart.CreatedAt = now.AddDate(0, 0, -2)
						cart.UpdatedAt = now.AddDate(0, 0, -2)
						return nil
					})
			},
			input: &store.AddCartItemInput{
				SessionID: "session-id",
				ProductID: "product-id",
				Quantity:  2,
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to multi get product",
			setup: func(ctx context.Context, mocks *mocks) {
				product := product("product-id", 1)
				mocks.db.Product.EXPECT().Get(ctx, "product-id").Return(product, nil)
				mocks.cache.EXPECT().
					Get(ctx, &entity.Cart{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, cart *entity.Cart) error {
						cart.Baskets = entity.CartBaskets{}
						cart.ExpiredAt = now.AddDate(0, 0, 12)
						cart.CreatedAt = now.AddDate(0, 0, -2)
						cart.UpdatedAt = now.AddDate(0, 0, -2)
						return nil
					})
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(nil, assert.AnError)
			},
			input: &store.AddCartItemInput{
				SessionID: "session-id",
				ProductID: "product-id",
				Quantity:  1,
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to put cart",
			setup: func(ctx context.Context, mocks *mocks) {
				product := product("product-id", 1)
				mocks.db.Product.EXPECT().Get(ctx, "product-id").Return(product, nil)
				mocks.cache.EXPECT().
					Get(ctx, &entity.Cart{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, cart *entity.Cart) error {
						cart.Baskets = entity.CartBaskets{}
						cart.ExpiredAt = now.AddDate(0, 0, 12)
						cart.CreatedAt = now.AddDate(0, 0, -2)
						cart.UpdatedAt = now.AddDate(0, 0, -2)
						return nil
					})
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(entity.Products{product}, nil)
				cart := &entity.Cart{
					SessionID: "session-id",
					Baskets: entity.CartBaskets{
						{
							BoxNumber: 1,
							BoxType:   entity.ShippingTypeNormal,
							BoxSize:   entity.ShippingSize60,
							BoxRate:   80,
							Items: entity.CartItems{
								{ProductID: "product-id", Quantity: 1},
							},
							CoordinatorID: "coordinator-id",
						},
					},
					ExpiredAt: now.Add(defaultCartTTL),
					CreatedAt: now.AddDate(0, 0, -2),
					UpdatedAt: now,
				}
				mocks.cache.EXPECT().Insert(ctx, cart).Return(assert.AnError)
			},
			input: &store.AddCartItemInput{
				SessionID: "session-id",
				ProductID: "product-id",
				Quantity:  1,
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.AddCartItem(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}, withNow(now)))
	}
}

func TestRemoveCartItem(t *testing.T) {
	t.Parallel()

	now := time.Now()
	product := func(productID string, invenroty int64) *entity.Product {
		return &entity.Product{
			ID:            productID,
			CoordinatorID: "coordinator-id",
			DeliveryType:  entity.DeliveryTypeNormal,
			Inventory:     invenroty,
			Public:        true,
			Status:        entity.ProductStatusForSale,
			Weight:        500,
			WeightUnit:    entity.WeightUnitGram,
			Box60Rate:     80,
			Box80Rate:     50,
			Box100Rate:    30,
		}
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.RemoveCartItemInput
		expectErr error
	}{
		{
			name: "success when single box",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.Cart{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, cart *entity.Cart) error {
						cart.Baskets = entity.CartBaskets{
							{
								BoxNumber: 1,
								BoxType:   entity.ShippingTypeNormal,
								BoxSize:   entity.ShippingSize60,
								BoxRate:   80,
								Items: entity.CartItems{
									{ProductID: "product-id", Quantity: 1},
								},
								CoordinatorID: "coordinator-id",
							},
						}
						cart.ExpiredAt = now.AddDate(0, 0, 12)
						cart.CreatedAt = now.AddDate(0, 0, -2)
						cart.UpdatedAt = now.AddDate(0, 0, -2)
						return nil
					})
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{}).Return(entity.Products{}, nil)
				cart := &entity.Cart{
					SessionID: "session-id",
					Baskets:   entity.CartBaskets{},
					ExpiredAt: now.Add(defaultCartTTL),
					CreatedAt: now.AddDate(0, 0, -2),
					UpdatedAt: now,
				}
				mocks.cache.EXPECT().Insert(ctx, cart).Return(nil)
			},
			input: &store.RemoveCartItemInput{
				SessionID: "session-id",
				BoxNumber: 1,
				ProductID: "product-id",
			},
			expectErr: nil,
		},
		{
			name: "success when some boxes",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.Cart{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, cart *entity.Cart) error {
						cart.Baskets = entity.CartBaskets{
							{
								BoxNumber: 1,
								BoxType:   entity.ShippingTypeNormal,
								BoxSize:   entity.ShippingSize60,
								BoxRate:   80,
								Items: entity.CartItems{
									{ProductID: "product-id", Quantity: 1},
								},
								CoordinatorID: "coordinator-id",
							},
							{
								BoxNumber: 2,
								BoxType:   entity.ShippingTypeNormal,
								BoxSize:   entity.ShippingSize60,
								BoxRate:   80,
								Items: entity.CartItems{
									{ProductID: "product-id", Quantity: 1},
								},
								CoordinatorID: "coordinator-id",
							},
						}
						cart.ExpiredAt = now.AddDate(0, 0, 12)
						cart.CreatedAt = now.AddDate(0, 0, -2)
						cart.UpdatedAt = now.AddDate(0, 0, -2)
						return nil
					})
				products := entity.Products{product("product-id", 1)}
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products, nil)
				cart := &entity.Cart{
					SessionID: "session-id",
					Baskets: entity.CartBaskets{
						{
							BoxNumber: 1,
							BoxType:   entity.ShippingTypeNormal,
							BoxSize:   entity.ShippingSize60,
							BoxRate:   80,
							Items: entity.CartItems{
								{ProductID: "product-id", Quantity: 1},
							},
							CoordinatorID: "coordinator-id",
						},
					},
					ExpiredAt: now.Add(defaultCartTTL),
					CreatedAt: now.AddDate(0, 0, -2),
					UpdatedAt: now,
				}
				mocks.cache.EXPECT().Insert(ctx, cart).Return(nil)
			},
			input: &store.RemoveCartItemInput{
				SessionID: "session-id",
				BoxNumber: 1,
				ProductID: "product-id",
			},
			expectErr: nil,
		},
		{
			name: "success without item",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.Cart{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, cart *entity.Cart) error {
						cart.Baskets = entity.CartBaskets{}
						cart.ExpiredAt = now.AddDate(0, 0, 12)
						cart.CreatedAt = now.AddDate(0, 0, -2)
						cart.UpdatedAt = now.AddDate(0, 0, -2)
						return nil
					})
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{}).Return(entity.Products{}, nil)
				cart := &entity.Cart{
					SessionID: "session-id",
					Baskets:   entity.CartBaskets{},
					ExpiredAt: now.Add(defaultCartTTL),
					CreatedAt: now.AddDate(0, 0, -2),
					UpdatedAt: now,
				}
				mocks.cache.EXPECT().Insert(ctx, cart).Return(nil)
			},
			input: &store.RemoveCartItemInput{
				SessionID: "session-id",
				BoxNumber: 1,
				ProductID: "product-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.RemoveCartItemInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get cart",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().Get(ctx, &entity.Cart{SessionID: "session-id"}).Return(assert.AnError)
			},
			input: &store.RemoveCartItemInput{
				SessionID: "session-id",
				BoxNumber: 1,
				ProductID: "product-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to multi get products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.Cart{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, cart *entity.Cart) error {
						cart.Baskets = entity.CartBaskets{
							{
								BoxNumber: 1,
								BoxType:   entity.ShippingTypeNormal,
								BoxSize:   entity.ShippingSize60,
								BoxRate:   80,
								Items: entity.CartItems{
									{ProductID: "product-id", Quantity: 1},
								},
								CoordinatorID: "coordinator-id",
							},
						}
						cart.ExpiredAt = now.AddDate(0, 0, 12)
						cart.CreatedAt = now.AddDate(0, 0, -2)
						cart.UpdatedAt = now.AddDate(0, 0, -2)
						return nil
					})
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{}).Return(nil, assert.AnError)
			},
			input: &store.RemoveCartItemInput{
				SessionID: "session-id",
				BoxNumber: 1,
				ProductID: "product-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to put cart",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.Cart{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, cart *entity.Cart) error {
						cart.Baskets = entity.CartBaskets{
							{
								BoxNumber: 1,
								BoxType:   entity.ShippingTypeNormal,
								BoxSize:   entity.ShippingSize60,
								BoxRate:   80,
								Items: entity.CartItems{
									{ProductID: "product-id", Quantity: 1},
								},
								CoordinatorID: "coordinator-id",
							},
						}
						cart.ExpiredAt = now.AddDate(0, 0, 12)
						cart.CreatedAt = now.AddDate(0, 0, -2)
						cart.UpdatedAt = now.AddDate(0, 0, -2)
						return nil
					})
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{}).Return(entity.Products{}, nil)
				cart := &entity.Cart{
					SessionID: "session-id",
					Baskets:   entity.CartBaskets{},
					ExpiredAt: now.Add(defaultCartTTL),
					CreatedAt: now.AddDate(0, 0, -2),
					UpdatedAt: now,
				}
				mocks.cache.EXPECT().Insert(ctx, cart).Return(assert.AnError)
			},
			input: &store.RemoveCartItemInput{
				SessionID: "session-id",
				BoxNumber: 1,
				ProductID: "product-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.RemoveCartItem(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}, withNow(now)))
	}
}
