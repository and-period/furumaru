package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
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
							BoxType:   entity.DeliveryTypeNormal,
							BoxSize:   entity.ShippingSize80,
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
							BoxType:   entity.DeliveryTypeNormal,
							BoxSize:   entity.ShippingSize60,
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
						BoxType:   entity.DeliveryTypeNormal,
						BoxSize:   entity.ShippingSize60,
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
							BoxType:   entity.DeliveryTypeNormal,
							BoxSize:   entity.ShippingSize80,
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
						BoxType:   entity.DeliveryTypeNormal,
						BoxSize:   entity.ShippingSize80,
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
							BoxType:   entity.DeliveryTypeNormal,
							BoxSize:   entity.ShippingSize80,
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
							BoxType:   entity.DeliveryTypeNormal,
							BoxSize:   entity.ShippingSize60,
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
							BoxType:   entity.DeliveryTypeNormal,
							BoxSize:   entity.ShippingSize60,
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
								BoxType:   entity.DeliveryTypeNormal,
								BoxSize:   entity.ShippingSize60,
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
								BoxType:   entity.DeliveryTypeNormal,
								BoxSize:   entity.ShippingSize60,
								Items: entity.CartItems{
									{ProductID: "product-id", Quantity: 1},
								},
								CoordinatorID: "coordinator-id",
							},
							{
								BoxNumber: 2,
								BoxType:   entity.DeliveryTypeNormal,
								BoxSize:   entity.ShippingSize60,
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
							BoxType:   entity.DeliveryTypeNormal,
							BoxSize:   entity.ShippingSize60,
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
								BoxType:   entity.DeliveryTypeNormal,
								BoxSize:   entity.ShippingSize60,
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
								BoxType:   entity.DeliveryTypeNormal,
								BoxSize:   entity.ShippingSize60,
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
