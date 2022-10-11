package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestOrderItem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		item   *entity.OrderItem
		expect *OrderItem
	}{
		{
			name: "success",
			item: &entity.OrderItem{
				ID:         "item-id",
				OrderID:    "order-id",
				ProductID:  "product-id",
				Price:      100,
				Quantity:   1,
				Weight:     1000,
				WeightUnit: entity.WeightUnitGram,
				CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &OrderItem{
				OrderItem: response.OrderItem{
					ProductID: "product-id",
					Price:     100,
					Quantity:  1,
					Weight:    1.0,
				},
				orderID: "order-id",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewOrderItem(tt.item))
		})
	}
}

func TestOrderItem_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		item    *OrderItem
		product *Product
		expect  *OrderItem
	}{
		{
			name: "success",
			item: &OrderItem{
				OrderItem: response.OrderItem{
					ProductID: "product-id",
					Price:     100,
					Quantity:  1,
					Weight:    1.0,
				},
				orderID: "order-id",
			},
			product: &Product{
				Product: response.Product{
					ID:              "product-id",
					TypeID:          "product-type-id",
					TypeName:        "",
					TypeIconURL:     "",
					CategoryID:      "category-id",
					CategoryName:    "",
					ProducerID:      "producer-id",
					StoreName:       "",
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Public:          true,
					Inventory:       100,
					Weight:          1.3,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: []*response.ProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
					Price:            400,
					DeliveryType:     int32(DeliveryTypeNormal),
					Box60Rate:        50,
					Box80Rate:        40,
					Box100Rate:       30,
					OriginPrefecture: "滋賀県",
					OriginCity:       "彦根市",
					CreatedAt:        1640962800,
					UpdatedAt:        1640962800,
					CreatedBy:        "coordinator-id",
					UpdatedBy:        "coordinator-id",
				},
			},
			expect: &OrderItem{
				OrderItem: response.OrderItem{
					ProductID: "product-id",
					Name:      "新鮮なじゃがいも",
					Price:     100,
					Quantity:  1,
					Weight:    1.0,
					Media: []*response.ProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
				},
				orderID: "order-id",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.item.Fill(tt.product)
			assert.Equal(t, tt.expect, tt.item)
		})
	}
}

func TestOrderItem_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		item   *OrderItem
		expect *response.OrderItem
	}{
		{
			name: "success",
			item: &OrderItem{
				OrderItem: response.OrderItem{
					ProductID: "product-id",
					Name:      "新鮮なじゃがいも",
					Price:     100,
					Quantity:  1,
					Weight:    1.0,
					Media: []*response.ProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
				},
				orderID: "order-id",
			},
			expect: &response.OrderItem{
				ProductID: "product-id",
				Name:      "新鮮なじゃがいも",
				Price:     100,
				Quantity:  1,
				Weight:    1.0,
				Media: []*response.ProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.item.Response())
		})
	}
}

func TestOrderItems(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		items  entity.OrderItems
		expect OrderItems
	}{
		{
			name: "success",
			items: entity.OrderItems{
				{
					ID:         "item-id",
					OrderID:    "order-id",
					ProductID:  "product-id",
					Price:      100,
					Quantity:   1,
					Weight:     1000,
					WeightUnit: entity.WeightUnitGram,
					CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: OrderItems{
				{
					OrderItem: response.OrderItem{
						ProductID: "product-id",
						Price:     100,
						Quantity:  1,
						Weight:    1.0,
					},
					orderID: "order-id",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewOrderItems(tt.items))
		})
	}
}

func TestOrderItems_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		items    OrderItems
		products map[string]*Product
		expect   OrderItems
	}{
		{
			name: "success",
			items: OrderItems{
				{
					OrderItem: response.OrderItem{
						ProductID: "product-id",
						Price:     100,
						Quantity:  1,
						Weight:    1.0,
					},
					orderID: "order-id",
				},
			},
			products: map[string]*Product{
				"product-id": {
					Product: response.Product{
						ID:              "product-id",
						TypeID:          "product-type-id",
						TypeName:        "",
						TypeIconURL:     "",
						CategoryID:      "category-id",
						CategoryName:    "",
						ProducerID:      "producer-id",
						StoreName:       "",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:            400,
						DeliveryType:     int32(DeliveryTypeNormal),
						Box60Rate:        50,
						Box80Rate:        40,
						Box100Rate:       30,
						OriginPrefecture: "滋賀県",
						OriginCity:       "彦根市",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
						CreatedBy:        "coordinator-id",
						UpdatedBy:        "coordinator-id",
					},
				},
			},
			expect: OrderItems{
				{
					OrderItem: response.OrderItem{
						ProductID: "product-id",
						Name:      "新鮮なじゃがいも",
						Price:     100,
						Quantity:  1,
						Weight:    1.0,
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
					},
					orderID: "order-id",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.items.Fill(tt.products)
			assert.Equal(t, tt.expect, tt.items)
		})
	}
}

func TestOrderItems_ProductIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		items  OrderItems
		expect []string
	}{
		{
			name: "success",
			items: OrderItems{
				{
					OrderItem: response.OrderItem{
						ProductID: "product-id",
						Name:      "新鮮なじゃがいも",
						Price:     100,
						Quantity:  1,
						Weight:    1.0,
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
					},
					orderID: "order-id",
				},
			},
			expect: []string{"product-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.items.ProductIDs())
		})
	}
}

func TestOrderItems_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		items  OrderItems
		expect []*response.OrderItem
	}{
		{
			name: "success",
			items: OrderItems{
				{
					OrderItem: response.OrderItem{
						ProductID: "product-id",
						Name:      "新鮮なじゃがいも",
						Price:     100,
						Quantity:  1,
						Weight:    1.0,
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
					},
					orderID: "order-id",
				},
			},
			expect: []*response.OrderItem{
				{
					ProductID: "product-id",
					Name:      "新鮮なじゃがいも",
					Price:     100,
					Quantity:  1,
					Weight:    1.0,
					Media: []*response.ProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.items.Response())
		})
	}
}
