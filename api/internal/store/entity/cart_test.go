package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCart_Refresh(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		cart     *Cart
		products Products
		hasErr   bool
	}{
		{
			name: "success",
			cart: &Cart{
				SessionID: "session-id",
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   DeliveryTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						coordinatorID: "",
					},
				},
				ExpiredAt: now.AddDate(0, 0, 7),
				CreatedAt: now,
				UpdatedAt: now,
			},
			products: Products{
				{
					ID:            "product-id",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     1,
					Weight:        500,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     80,
					Box80Rate:     50,
					Box100Rate:    30,
				},
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.cart.Refresh(tt.products)
			assert.Equal(t, tt.hasErr, err != nil, err)
		})
	}
}

func TestCartBaskets_ProductIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		baskets CartBaskets
		expect  []string
	}{
		{
			name: "success",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
				},
				{
					BoxNumber: 2,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 3},
					},
				},
			},
			expect: []string{
				"product-id01",
				"product-id02",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.baskets.ProductIDs())
		})
	}
}

func TestCartItem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		productID string
		quantity  int64
		expect    *CartItem
	}{
		{
			name:      "success",
			productID: "product-id",
			quantity:  1,
			expect: &CartItem{
				ProductID: "product-id",
				Quantity:  1,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewCartItem(tt.productID, tt.quantity)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestCartItems(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		items  map[string]int64
		expect CartItems
	}{
		{
			name: "success",
			items: map[string]int64{
				"product-id01": 4,
				"product-id02": 2,
			},
			expect: CartItems{
				{
					ProductID: "product-id01",
					Quantity:  4,
				},
				{
					ProductID: "product-id02",
					Quantity:  2,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewCartItems(tt.items)
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestCartItemsWithProducts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products Products
		expect   CartItems
	}{
		{
			name: "success",
			products: Products{
				{ID: "product-id01"},
				{ID: "product-id02"},
				{ID: "product-id01"},
				{ID: "product-id02"},
				{ID: "product-id02"},
			},
			expect: CartItems{
				{
					ProductID: "product-id01",
					Quantity:  2,
				},
				{
					ProductID: "product-id02",
					Quantity:  3,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewCartItemsWithProducts(tt.products)
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestCartItems_ProductIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		items  CartItems
		expect []string
	}{
		{
			name: "success",
			items: CartItems{
				{
					ProductID: "product-id01",
					Quantity:  1,
				},
				{
					ProductID: "product-id02",
					Quantity:  2,
				},
				{
					ProductID: "product-id01",
					Quantity:  3,
				},
			},
			expect: []string{"product-id01", "product-id02"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.items.ProductIDs())
		})
	}
}

func TestGenerateBascketKey(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		coordinatorID string
		deliveryType  DeliveryType
		expect        string
	}{
		{
			name:          "success",
			coordinatorID: "coordinator-id",
			deliveryType:  DeliveryTypeNormal,
			expect:        "coordinator-id:1",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := generateCartBasketKey(tt.coordinatorID, tt.deliveryType)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestParseBascketKey(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                string
		key                 string
		expectCoordinatorID string
		expectDeliveryType  DeliveryType
		hasErr              bool
	}{
		{
			name:                "success",
			key:                 "coordinator-id:1",
			expectCoordinatorID: "coordinator-id",
			expectDeliveryType:  DeliveryTypeNormal,
			hasErr:              false,
		},
		{
			name:                "invalid format",
			key:                 "",
			expectCoordinatorID: "",
			expectDeliveryType:  DeliveryTypeUnknown,
			hasErr:              true,
		},
		{
			name:                "unknotn delivery type",
			key:                 "coordinator-id:delivery-type",
			expectCoordinatorID: "",
			expectDeliveryType:  DeliveryTypeUnknown,
			hasErr:              true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			coordinatorID, deliveryType, err := parseCartBasketKey(tt.key)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expectCoordinatorID, coordinatorID)
			assert.Equal(t, tt.expectDeliveryType, deliveryType)
		})
	}
}

func TestRefreshCart(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		baskets  CartBaskets
		products map[string]*Product
		expect   CartBaskets
		hasErr   bool
	}{
		{
			name: "success 商品なし",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 1},
					},
					coordinatorID: "coordinator-id",
				},
			},
			products: map[string]*Product{},
			expect:   CartBaskets{},
			hasErr:   false,
		},
		{
			name: "success 在庫切れ",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 1},
					},
					coordinatorID: "coordinator-id",
				},
			},
			products: map[string]*Product{
				"product-id": {
					ID:            "product-id",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     0,
					Weight:        500,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     80,
					Box80Rate:     50,
					Box100Rate:    30,
				},
			},
			expect: CartBaskets{},
			hasErr: false,
		},
		{
			name: "success 商品が1つ 箱のサイズ60",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 1},
					},
					coordinatorID: "coordinator-id",
				},
			},
			products: map[string]*Product{
				"product-id": {
					ID:            "product-id",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     1,
					Weight:        500,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     80,
					Box80Rate:     50,
					Box100Rate:    30,
				},
			},
			expect: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  1,
						},
					},
					coordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 同じ商品が2つ 箱のサイズ80",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 2},
					},
					coordinatorID: "coordinator-id",
				},
			},
			products: map[string]*Product{
				"product-id": {
					ID:            "product-id",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     30,
					Weight:        500,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     80,
					Box80Rate:     50,
					Box100Rate:    30,
				},
			},
			expect: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize80,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  2,
						},
					},
					coordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 同じ商品が3つ 箱のサイズ100",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 3},
					},
					coordinatorID: "coordinator-id",
				},
			},
			products: map[string]*Product{
				"product-id": {
					ID:            "product-id",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     30,
					Weight:        500,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     80,
					Box80Rate:     50,
					Box100Rate:    30,
				},
			},
			expect: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  3,
						},
					},
					coordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 同じ商品が4つ 箱のサイズ60と100",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 4},
					},
					coordinatorID: "coordinator-id",
				},
			},
			products: map[string]*Product{
				"product-id": {
					ID:            "product-id",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     30,
					Weight:        500,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     80,
					Box80Rate:     50,
					Box100Rate:    30,
				},
			},
			expect: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  3,
						},
					},
					coordinatorID: "coordinator-id",
				},
				{
					BoxNumber: 2,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  1,
						},
					},
					coordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 複数商品 同じ箱",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
					coordinatorID: "coordinator-id",
				},
			},
			products: map[string]*Product{
				"product-id01": {
					ID:            "product-id01",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     30,
					Weight:        500,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     80,
					Box80Rate:     50,
					Box100Rate:    30,
				},
				"product-id02": {
					ID:            "product-id02",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     30,
					Weight:        100,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     30,
					Box80Rate:     25,
					Box100Rate:    10,
				},
			},
			expect: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize80,
					Items: CartItems{
						{
							ProductID: "product-id01",
							Quantity:  1,
						},
						{
							ProductID: "product-id02",
							Quantity:  2,
						},
					},
					coordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 複数商品 複数の箱",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 3},
						{ProductID: "product-id02", Quantity: 4},
						{ProductID: "product-id03", Quantity: 2},
					},
					coordinatorID: "coordinator-id",
				},
			},
			products: map[string]*Product{
				"product-id01": {
					ID:            "product-id01",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     30,
					Weight:        500,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     80,
					Box80Rate:     50,
					Box100Rate:    30,
				},
				"product-id02": {
					ID:            "product-id02",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     30,
					Weight:        2,
					WeightUnit:    WeightUnitKilogram,
					Box60Rate:     80,
					Box80Rate:     65,
					Box100Rate:    50,
				},
				"product-id03": {
					ID:            "product-id03",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     30,
					Weight:        100,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     30,
					Box80Rate:     25,
					Box100Rate:    10,
				},
			},
			expect: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{
							ProductID: "product-id02",
							Quantity:  2,
						},
					},
					coordinatorID: "coordinator-id",
				},
				{
					BoxNumber: 2,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{
							ProductID: "product-id02",
							Quantity:  2,
						},
					},
					coordinatorID: "coordinator-id",
				},
				{
					BoxNumber: 3,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{
							ProductID: "product-id01",
							Quantity:  3,
						},
						{
							ProductID: "product-id03",
							Quantity:  1,
						},
					},
					coordinatorID: "coordinator-id",
				},
				{
					BoxNumber: 4,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{
							ProductID: "product-id03",
							Quantity:  1,
						},
					},
					coordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 箱のサイズ60 重量上限",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 4},
					},
					coordinatorID: "coordinator-id",
				},
			},
			products: map[string]*Product{
				"product-id": {
					ID:            "product-id",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     30,
					Weight:        500,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     20,
					Box80Rate:     10,
					Box100Rate:    5,
				},
			},
			expect: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  4,
						},
					},
					coordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 箱のサイズ80 重量上限",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 2},
					},
					coordinatorID: "coordinator-id",
				},
			},
			products: map[string]*Product{
				"product-id": {
					ID:            "product-id",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     30,
					Weight:        2500,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     20,
					Box80Rate:     10,
					Box100Rate:    5,
				},
			},
			expect: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize80,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  2,
						},
					},
					coordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 箱のサイズ100 重量上限",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 4},
					},
					coordinatorID: "coordinator-id",
				},
			},
			products: map[string]*Product{
				"product-id": {
					ID:            "product-id",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     30,
					Weight:        2500,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     20,
					Box80Rate:     10,
					Box100Rate:    5,
				},
			},
			expect: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  4,
						},
					},
					coordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 箱のサイズ100 重量超過",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 5},
					},
					coordinatorID: "coordinator-id",
				},
			},
			products: map[string]*Product{
				"product-id": {
					ID:            "product-id",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     30,
					Weight:        2500,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     20,
					Box80Rate:     10,
					Box100Rate:    5,
				},
			},
			expect: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  4,
						},
					},
					coordinatorID: "coordinator-id",
				},
				{
					BoxNumber: 2,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize80,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  1,
						},
					},
					coordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 別のコーディネータの商品",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 2},
						{ProductID: "product-id02", Quantity: 2},
					},
					coordinatorID: "",
				},
			},
			products: map[string]*Product{
				"product-id01": {
					ID:            "product-id01",
					CoordinatorID: "coordinator-id01",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     30,
					Weight:        500,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     80,
					Box80Rate:     50,
					Box100Rate:    30,
				},
				"product-id02": {
					ID:            "product-id02",
					CoordinatorID: "coordinator-id02",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     30,
					Weight:        100,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     30,
					Box80Rate:     25,
					Box100Rate:    10,
				},
			},
			expect: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize80,
					Items: CartItems{
						{
							ProductID: "product-id01",
							Quantity:  2,
						},
					},
					coordinatorID: "coordinator-id01",
				},
				{
					BoxNumber: 2,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{
							ProductID: "product-id02",
							Quantity:  2,
						},
					},
					coordinatorID: "coordinator-id02",
				},
			},
			hasErr: false,
		},
		{
			name: "success 別の配送方法の商品",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 2},
						{ProductID: "product-id02", Quantity: 2},
					},
					coordinatorID: "",
				},
			},
			products: map[string]*Product{
				"product-id01": {
					ID:            "product-id01",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeNormal,
					Inventory:     30,
					Weight:        500,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     80,
					Box80Rate:     50,
					Box100Rate:    30,
				},
				"product-id02": {
					ID:            "product-id02",
					CoordinatorID: "coordinator-id",
					DeliveryType:  DeliveryTypeFrozen,
					Inventory:     30,
					Weight:        100,
					WeightUnit:    WeightUnitGram,
					Box60Rate:     30,
					Box80Rate:     25,
					Box100Rate:    10,
				},
			},
			expect: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   DeliveryTypeNormal,
					BoxSize:   ShippingSize80,
					Items: CartItems{
						{
							ProductID: "product-id01",
							Quantity:  2,
						},
					},
					coordinatorID: "coordinator-id",
				},
				{
					BoxNumber: 2,
					BoxType:   DeliveryTypeFrozen,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{
							ProductID: "product-id02",
							Quantity:  2,
						},
					},
					coordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := refreshCart(tt.baskets, tt.products)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
