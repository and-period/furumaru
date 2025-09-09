package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/stretchr/testify/assert"
)

func TestCart(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		basket *entity.CartBasket
		expect *Cart
	}{
		{
			name: "success",
			basket: &entity.CartBasket{
				BoxNumber: 1,
				BoxType:   entity.ShippingTypeNormal,
				BoxSize:   entity.ShippingSize60,
				BoxRate:   80,
				Items: []*entity.CartItem{
					{
						ProductID: "product-id",
						Quantity:  1,
					},
				},
				CoordinatorID: "coordinator-id",
			},
			expect: &Cart{
				Cart: types.Cart{
					Number: 1,
					Type:   ShippingTypeNormal.Response(),
					Size:   ShippingSize60.Response(),
					Rate:   80,
					Items: []*types.CartItem{
						{
							ProductID: "product-id",
							Quantity:  1,
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewCart(tt.basket)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestCart_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		cart   *Cart
		expect *types.Cart
	}{
		{
			name: "success",
			cart: &Cart{
				Cart: types.Cart{
					Number: 1,
					Type:   ShippingTypeNormal.Response(),
					Size:   ShippingSize60.Response(),
					Rate:   80,
					Items: []*types.CartItem{
						{
							ProductID: "product-id",
							Quantity:  1,
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			expect: &types.Cart{
				Number: 1,
				Type:   ShippingTypeNormal.Response(),
				Size:   ShippingSize60.Response(),
				Rate:   80,
				Items: []*types.CartItem{
					{
						ProductID: "product-id",
						Quantity:  1,
					},
				},
				CoordinatorID: "coordinator-id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.cart.Response())
		})
	}
}

func TestCarts(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name   string
		cart   *entity.Cart
		expect Carts
	}{
		{
			name: "success",
			cart: &entity.Cart{
				SessionID: "session-id",
				Baskets: entity.CartBaskets{
					{
						BoxNumber: 1,
						BoxType:   entity.ShippingTypeNormal,
						BoxSize:   entity.ShippingSize60,
						BoxRate:   80,
						Items: []*entity.CartItem{
							{
								ProductID: "product-id",
								Quantity:  1,
							},
						},
						CoordinatorID: "coordinator-id",
					},
				},
				ExpiredAt: now.AddDate(0, 0, 14),
				CreatedAt: now,
				UpdatedAt: now,
			},
			expect: Carts{
				{
					Cart: types.Cart{
						Number: 1,
						Type:   ShippingTypeNormal.Response(),
						Size:   ShippingSize60.Response(),
						Rate:   80,
						Items: []*types.CartItem{
							{
								ProductID: "product-id",
								Quantity:  1,
							},
						},
						CoordinatorID: "coordinator-id",
					},
				},
			},
		},
		{
			name:   "empty",
			cart:   nil,
			expect: Carts{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewCarts(tt.cart)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestCarts_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		carts  Carts
		expect []*types.Cart
	}{
		{
			name: "success",
			carts: Carts{
				{
					Cart: types.Cart{
						Number: 1,
						Type:   ShippingTypeNormal.Response(),
						Size:   ShippingSize60.Response(),
						Rate:   80,
						Items: []*types.CartItem{
							{
								ProductID: "product-id",
								Quantity:  1,
							},
						},
						CoordinatorID: "coordinator-id",
					},
				},
			},
			expect: []*types.Cart{
				{
					Number: 1,
					Type:   ShippingTypeNormal.Response(),
					Size:   ShippingSize60.Response(),
					Rate:   80,
					Items: []*types.CartItem{
						{
							ProductID: "product-id",
							Quantity:  1,
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.carts.Response())
		})
	}
}

func TestCartItem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		item   *entity.CartItem
		expect *CartItem
	}{
		{
			name: "success",
			item: &entity.CartItem{
				ProductID: "product-id",
				Quantity:  1,
			},
			expect: &CartItem{
				CartItem: types.CartItem{
					ProductID: "product-id",
					Quantity:  1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewCartItem(tt.item)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestCartItem_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		item   *CartItem
		expect *types.CartItem
	}{
		{
			name: "success",
			item: &CartItem{
				CartItem: types.CartItem{
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			expect: &types.CartItem{
				ProductID: "product-id",
				Quantity:  1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.item.Response())
		})
	}
}

func TestCartItems(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		items  entity.CartItems
		expect CartItems
	}{
		{
			name: "success",
			items: entity.CartItems{
				{
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			expect: CartItems{
				{
					CartItem: types.CartItem{
						ProductID: "product-id",
						Quantity:  1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewCartItems(tt.items)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestCartItems_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		items  CartItems
		expect []*types.CartItem
	}{
		{
			name: "success",
			items: CartItems{
				{
					CartItem: types.CartItem{
						ProductID: "product-id",
						Quantity:  1,
					},
				},
			},
			expect: []*types.CartItem{
				{
					ProductID: "product-id",
					Quantity:  1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.items.Response())
		})
	}
}
