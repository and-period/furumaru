package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCart(t *testing.T) {
	t.Parallel()
	now := time.Now()
	params := &CartParams{
		SessionID: "session-id",
		Now:       now,
		TTL:       time.Hour,
	}
	cart := NewCart(params)
	t.Run("new", func(t *testing.T) {
		t.Parallel()
		expect := &Cart{
			SessionID: "session-id",
			Baskets:   []*CartBasket{},
			ExpiredAt: now.Add(time.Hour),
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.Equal(t, expect, cart)
	})
	t.Run("table name", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "carts", cart.TableName())
	})
	t.Run("primary key", func(t *testing.T) {
		t.Parallel()
		expect := map[string]interface{}{
			"session_id": "session-id",
		}
		assert.Equal(t, expect, cart.PrimaryKey())
	})
}

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
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
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

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.cart.Refresh(tt.products)
			assert.Equal(t, tt.hasErr, err != nil, err)
		})
	}
}

func TestCart_AddItem(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name      string
		cart      *Cart
		productID string
		quantity  int64
		expect    *Cart
	}{
		{
			name: "success",
			cart: &Cart{
				SessionID: "session-id",
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
				},
				ExpiredAt: now.AddDate(0, 0, 7),
				CreatedAt: now,
				UpdatedAt: now,
			},
			productID: "product-id",
			quantity:  1,
			expect: &Cart{
				SessionID: "session-id",
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
					{
						BoxNumber: 0,
						BoxType:   ShippingTypeUnknown,
						BoxSize:   ShippingSizeUnknown,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
				},
				ExpiredAt: now.AddDate(0, 0, 7),
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name:      "empty",
			cart:      &Cart{},
			productID: "product-id",
			quantity:  1,
			expect: &Cart{
				Baskets: []*CartBasket{
					{
						BoxNumber: 0,
						BoxType:   ShippingTypeUnknown,
						BoxSize:   ShippingSizeUnknown,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.cart.AddItem(tt.productID, tt.quantity)
			assert.Equal(t, tt.expect, tt.cart)
		})
	}
}

func TestCart_RemoveBaskets(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name       string
		cart       *Cart
		boxNumbers []int64
		expect     *Cart
	}{
		{
			name: "success",
			cart: &Cart{
				SessionID: "session-id",
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
				},
				ExpiredAt: now.AddDate(0, 0, 7),
				CreatedAt: now,
				UpdatedAt: now,
			},
			boxNumbers: []int64{1},
			expect: &Cart{
				SessionID: "session-id",
				Baskets:   []*CartBasket{},
				ExpiredAt: now.AddDate(0, 0, 7),
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.cart.RemoveBaskets(tt.boxNumbers...)
			assert.Equal(t, tt.expect, tt.cart)
		})
	}
}

func TestCart_RemoveItem(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name      string
		cart      *Cart
		productID string
		boxNumber int64
		expect    *Cart
	}{
		{
			name: "success with box number",
			cart: &Cart{
				SessionID: "session-id",
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
					{
						BoxNumber: 2,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
				},
				ExpiredAt: now.AddDate(0, 0, 7),
				CreatedAt: now,
				UpdatedAt: now,
			},
			productID: "product-id",
			boxNumber: 2,
			expect: &Cart{
				SessionID: "session-id",
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
					{
						BoxNumber:     2,
						BoxType:       ShippingTypeNormal,
						BoxSize:       ShippingSize60,
						Items:         CartItems{},
						CoordinatorID: "",
					},
				},
				ExpiredAt: now.AddDate(0, 0, 7),
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "success without box number",
			cart: &Cart{
				SessionID: "session-id",
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
					{
						BoxNumber: 2,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
				},
				ExpiredAt: now.AddDate(0, 0, 7),
				CreatedAt: now,
				UpdatedAt: now,
			},
			productID: "product-id",
			expect: &Cart{
				SessionID: "session-id",
				Baskets: []*CartBasket{
					{
						BoxNumber:     1,
						BoxType:       ShippingTypeNormal,
						BoxSize:       ShippingSize60,
						Items:         CartItems{},
						CoordinatorID: "",
					},
					{
						BoxNumber:     2,
						BoxType:       ShippingTypeNormal,
						BoxSize:       ShippingSize60,
						Items:         CartItems{},
						CoordinatorID: "",
					},
				},
				ExpiredAt: now.AddDate(0, 0, 7),
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "success without item",
			cart: &Cart{
				SessionID: "session-id",
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
				},
				ExpiredAt: now.AddDate(0, 0, 7),
				CreatedAt: now,
				UpdatedAt: now,
			},
			productID: "other-id",
			expect: &Cart{
				SessionID: "session-id",
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
				},
				ExpiredAt: now.AddDate(0, 0, 7),
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "success empty",
			cart: &Cart{
				SessionID: "session-id",
				Baskets:   []*CartBasket{},
				ExpiredAt: now.AddDate(0, 0, 7),
				CreatedAt: now,
				UpdatedAt: now,
			},
			productID: "product-id",
			boxNumber: 1,
			expect: &Cart{
				SessionID: "session-id",
				Baskets:   []*CartBasket{},
				ExpiredAt: now.AddDate(0, 0, 7),
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.cart.RemoveItem(tt.productID, tt.boxNumber)
			assert.Equal(t, tt.expect, tt.cart)
		})
	}
}

func TestCart_DecreaseItem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		cart      *Cart
		productID string
		quantity  int64
		expect    *Cart
	}{
		{
			name: "success 1つの買い物かごから削除",
			cart: &Cart{
				SessionID: "session-id",
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
					{
						BoxNumber: 2,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
				},
			},
			productID: "product-id",
			quantity:  1,
			expect: &Cart{
				SessionID: "session-id",
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 0},
						},
						CoordinatorID: "",
					},
					{
						BoxNumber: 2,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
				},
			},
		},
		{
			name: "success 複数の買い物かごから削除",
			cart: &Cart{
				SessionID: "session-id",
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
					{
						BoxNumber: 2,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
				},
			},
			productID: "product-id",
			quantity:  3,
			expect: &Cart{
				SessionID: "session-id",
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 0},
						},
						CoordinatorID: "",
					},
					{
						BoxNumber: 2,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 0},
						},
						CoordinatorID: "",
					},
				},
			},
		},
		{
			name: "success 対象の商品が存在しない",
			cart: &Cart{
				SessionID: "session-id",
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
					{
						BoxNumber: 2,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
				},
			},
			productID: "other-id",
			quantity:  3,
			expect: &Cart{
				SessionID: "session-id",
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
					{
						BoxNumber: 2,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: CartItems{
							{ProductID: "product-id", Quantity: 1},
						},
						CoordinatorID: "",
					},
				},
			},
		},
		{
			name: "success 買い物かごが空",
			cart: &Cart{
				SessionID: "session-id",
				Baskets:   []*CartBasket{},
			},
			productID: "product-id",
			quantity:  1,
			expect: &Cart{
				SessionID: "session-id",
				Baskets:   []*CartBasket{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.cart.DecreaseItem(tt.productID, tt.quantity)
			assert.Equal(t, tt.expect, tt.cart)
		})
	}
}

func TestCartBaskets_MergeByProductID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		baskets CartBaskets
		expect  CartItems
	}{
		{
			name: "success",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
					CoordinatorID: "",
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
					},
					CoordinatorID: "",
				},
			},
			expect: CartItems{
				{ProductID: "product-id01", Quantity: 2},
				{ProductID: "product-id02", Quantity: 2},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.baskets.MergeByProductID()
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestCartBaskets_AdjustItems(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		baskets  CartBaskets
		products map[string]*Product
		expect   CartItems
	}{
		{
			name: "success",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
					CoordinatorID: "",
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
					},
					CoordinatorID: "",
				},
			},
			products: map[string]*Product{
				"product-id01": {
					ID:        "product-id01",
					Inventory: 2,
				},
				"product-id02": {
					ID:        "product-id02",
					Inventory: 1,
				},
			},
			expect: CartItems{
				{ProductID: "product-id01", Quantity: 2},
				{ProductID: "product-id02", Quantity: 1},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.baskets.AdjustItems(tt.products)
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestCartBaskets_FilterByCoordinatorID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		baskets        CartBaskets
		coordinatorIDs []string
		expect         CartBaskets
	}{
		{
			name: "success all match",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
					CoordinatorID: "coordinator-id",
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			coordinatorIDs: []string{"coordinator-id"},
			expect: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
					CoordinatorID: "coordinator-id",
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
					},
					CoordinatorID: "coordinator-id",
				},
			},
		},
		{
			name: "success partial match",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
					CoordinatorID: "coordinator-id01",
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
					},
					CoordinatorID: "coordinator-id02",
				},
			},
			coordinatorIDs: []string{"coordinator-id01"},
			expect: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
					CoordinatorID: "coordinator-id01",
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.baskets.FilterByCoordinatorID(tt.coordinatorIDs...)
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestCartBaskets_FilterByBoxNumber(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		baskets    CartBaskets
		boxNumbers []int64
		expect     CartBaskets
	}{
		{
			name: "success with 0",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
					CoordinatorID: "",
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
					},
					CoordinatorID: "",
				},
			},
			boxNumbers: []int64{0},
			expect: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
					CoordinatorID: "",
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
					},
					CoordinatorID: "",
				},
			},
		},
		{
			name: "success without 0",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
					CoordinatorID: "",
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
					},
					CoordinatorID: "",
				},
			},
			boxNumbers: []int64{1},
			expect: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
					CoordinatorID: "",
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.baskets.FilterByBoxNumber(tt.boxNumbers...)
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestCartBaskets_VerifyQuantity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		baskets    CartBaskets
		additional int64
		product    *Product
		hasErr     bool
	}{
		{
			name: "success with item",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
				},
			},
			additional: 1,
			product: &Product{
				ID:        "product-id01",
				Inventory: 2,
			},
			hasErr: false,
		},
		{
			name: "success without item",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items:     CartItems{},
				},
			},
			additional: 2,
			product: &Product{
				ID:        "product-id01",
				Inventory: 2,
			},
			hasErr: false,
		},
		{
			name: "insufficient product stock",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
				},
			},
			additional: 2,
			product: &Product{
				ID:        "product-id01",
				Inventory: 2,
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.baskets.VerifyQuantity(tt.additional, tt.product)
			assert.Equal(t, tt.hasErr, err != nil, err)
		})
	}
}

func TestCartBaskets_TotalPrice(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		baskets   CartBaskets
		products  map[string]*Product
		expect    int64
		expectErr error
	}{
		{
			name: "success",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 3},
					},
				},
			},
			products: map[string]*Product{
				"product-id01": {
					ID: "product-id01",
					ProductRevision: ProductRevision{
						ID:        1,
						ProductID: "product-id01",
						Price:     500,
						Cost:      200,
					},
				},
				"product-id02": {
					ID: "product-id02",
					ProductRevision: ProductRevision{
						ID:        1,
						ProductID: "product-id02",
						Price:     1980,
						Cost:      500,
					},
				},
			},
			expect:    5960,
			expectErr: nil,
		},
		{
			name: "not found product",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 3},
					},
				},
			},
			products:  map[string]*Product{},
			expect:    0,
			expectErr: errNotFoundProduct,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tt.baskets.TotalPrice(tt.products)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
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

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.baskets.ProductIDs())
		})
	}
}

func TestCartBaskets_CoordinatorIDs(t *testing.T) {
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
					CoordinatorID: "coordinator-id",
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 3},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			expect: []string{"coordinator-id"},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.baskets.CoordinatorID())
		})
	}
}

func TestCartBaskets_BoxNumbers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		baskets CartBaskets
		expect  []int64
	}{
		{
			name: "success",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
					CoordinatorID: "coordinator-id",
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 3},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			expect: []int64{1, 2},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.baskets.BoxNumbers())
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

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.items.ProductIDs())
		})
	}
}

func TestCartItems_MapByProductID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		items  CartItems
		expect map[string]*CartItem
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
			},
			expect: map[string]*CartItem{
				"product-id01": {
					ProductID: "product-id01",
					Quantity:  1,
				},
				"product-id02": {
					ProductID: "product-id02",
					Quantity:  2,
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.items.MapByProductID())
		})
	}
}

func TestGenerateBascketKey(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		coordinatorID string
		shippingType  ShippingType
		expect        string
	}{
		{
			name:          "success",
			coordinatorID: "coordinator-id",
			shippingType:  ShippingTypeNormal,
			expect:        "coordinator-id:1",
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := generateCartBasketKey(tt.coordinatorID, tt.shippingType)
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
		expectShippingType  ShippingType
		hasErr              bool
	}{
		{
			name:                "success",
			key:                 "coordinator-id:1",
			expectCoordinatorID: "coordinator-id",
			expectShippingType:  ShippingTypeNormal,
			hasErr:              false,
		},
		{
			name:                "invalid format",
			key:                 "",
			expectCoordinatorID: "",
			expectShippingType:  ShippingTypeUnknown,
			hasErr:              true,
		},
		{
			name:                "unknotn delivery type",
			key:                 "coordinator-id:delivery-type",
			expectCoordinatorID: "",
			expectShippingType:  ShippingTypeUnknown,
			hasErr:              true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			coordinatorID, shippingType, err := parseCartBasketKey(tt.key)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expectCoordinatorID, coordinatorID)
			assert.Equal(t, tt.expectShippingType, shippingType)
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 1},
					},
					CoordinatorID: "coordinator-id",
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 1},
					},
					CoordinatorID: "coordinator-id",
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 1},
					},
					CoordinatorID: "coordinator-id",
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					BoxRate:   80,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  1,
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 同じ商品が2つ 箱のサイズ80",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 2},
					},
					CoordinatorID: "coordinator-id",
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize80,
					BoxRate:   100,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  2,
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 同じ商品が3つ 箱のサイズ100",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 3},
					},
					CoordinatorID: "coordinator-id",
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					BoxRate:   90,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  3,
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 同じ商品が4つ 箱のサイズ60と100",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 4},
					},
					CoordinatorID: "coordinator-id",
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					BoxRate:   90,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  3,
						},
					},
					CoordinatorID: "coordinator-id",
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					BoxRate:   80,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  1,
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 複数商品 同じ箱",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 1},
						{ProductID: "product-id02", Quantity: 2},
					},
					CoordinatorID: "coordinator-id",
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize80,
					BoxRate:   100,
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
					CoordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 複数商品 複数の箱",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 3},
						{ProductID: "product-id02", Quantity: 4},
						{ProductID: "product-id03", Quantity: 2},
					},
					CoordinatorID: "coordinator-id",
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					BoxRate:   100,
					Items: CartItems{
						{
							ProductID: "product-id02",
							Quantity:  2,
						},
					},
					CoordinatorID: "coordinator-id",
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					BoxRate:   100,
					Items: CartItems{
						{
							ProductID: "product-id02",
							Quantity:  2,
						},
					},
					CoordinatorID: "coordinator-id",
				},
				{
					BoxNumber: 3,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					BoxRate:   100,
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
					CoordinatorID: "coordinator-id",
				},
				{
					BoxNumber: 4,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					BoxRate:   30,
					Items: CartItems{
						{
							ProductID: "product-id03",
							Quantity:  1,
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 箱のサイズ60 重量上限",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 4},
					},
					CoordinatorID: "coordinator-id",
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					BoxRate:   80,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  4,
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 箱のサイズ80 重量上限",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 2},
					},
					CoordinatorID: "coordinator-id",
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize80,
					BoxRate:   20,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  2,
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 箱のサイズ100 重量上限",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 4},
					},
					CoordinatorID: "coordinator-id",
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					BoxRate:   20,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  4,
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 箱のサイズ100 重量超過",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id", Quantity: 5},
					},
					CoordinatorID: "coordinator-id",
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					BoxRate:   20,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  4,
						},
					},
					CoordinatorID: "coordinator-id",
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize80,
					BoxRate:   10,
					Items: CartItems{
						{
							ProductID: "product-id",
							Quantity:  1,
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
		{
			name: "success 別のコーディネータの商品",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 2},
						{ProductID: "product-id02", Quantity: 2},
					},
					CoordinatorID: "",
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize80,
					BoxRate:   100,
					Items: CartItems{
						{
							ProductID: "product-id01",
							Quantity:  2,
						},
					},
					CoordinatorID: "coordinator-id01",
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize60,
					BoxRate:   60,
					Items: CartItems{
						{
							ProductID: "product-id02",
							Quantity:  2,
						},
					},
					CoordinatorID: "coordinator-id02",
				},
			},
			hasErr: false,
		},
		{
			name: "success 別の配送方法の商品",
			baskets: CartBaskets{
				{
					BoxNumber: 1,
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize100,
					Items: CartItems{
						{ProductID: "product-id01", Quantity: 2},
						{ProductID: "product-id02", Quantity: 2},
					},
					CoordinatorID: "",
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
					BoxType:   ShippingTypeNormal,
					BoxSize:   ShippingSize80,
					BoxRate:   100,
					Items: CartItems{
						{
							ProductID: "product-id01",
							Quantity:  2,
						},
					},
					CoordinatorID: "coordinator-id",
				},
				{
					BoxNumber: 2,
					BoxType:   ShippingTypeFrozen,
					BoxSize:   ShippingSize60,
					BoxRate:   60,
					Items: CartItems{
						{
							ProductID: "product-id02",
							Quantity:  2,
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := refreshCart(tt.baskets, tt.products)
			assert.Equal(t, tt.hasErr, err != nil, err)
			require.Len(t, actual, len(tt.expect))
			for i := range actual {
				assert.ElementsMatch(t, tt.expect[i].Items, actual[i].Items)
				tt.expect[i].Items = actual[i].Items // 商品順は保証されていないため
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}
