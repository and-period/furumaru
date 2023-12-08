package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestOrderItem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		item    *entity.OrderItem
		product *Product
		expect  *OrderItem
	}{
		{
			name: "success",
			item: &entity.OrderItem{
				FulfillmentID:     "fulfillment-id",
				OrderID:           "order-id",
				ProductRevisionID: 1,
				Quantity:          1,
				CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			product: &Product{
				Product: response.Product{
					ID:              "product-id",
					CoordinatorID:   "coordinator-id",
					ProducerID:      "producer-id",
					CategoryID:      "",
					ProductTypeID:   "product-type-id",
					ProductTagIDs:   []string{"product-tag-id"},
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Status:          int32(ProductStatusForSale),
					Inventory:       100,
					Weight:          1.3,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: []*response.ProductMedia{
						{
							URL:         "https://and-period.jp/thumbnail01.png",
							IsThumbnail: true,
							Images:      []*response.Image{},
						},
						{
							URL:         "https://and-period.jp/thumbnail02.png",
							IsThumbnail: false,
							Images:      []*response.Image{},
						},
					},
					Price:             400,
					RecommendedPoint1: "ポイント1",
					RecommendedPoint2: "ポイント2",
					RecommendedPoint3: "ポイント3",
					StorageMethodType: int32(StorageMethodTypeNormal),
					DeliveryType:      int32(DeliveryTypeNormal),
					Box60Rate:         50,
					Box80Rate:         40,
					Box100Rate:        30,
					OriginCity:        "彦根市",
					StartAt:           1640962800,
					EndAt:             1640962800,
				},
			},
			expect: &OrderItem{
				OrderItem: response.OrderItem{
					FulfillmentID: "fulfillment-id",
					ProductID:     "product-id",
					Price:         400,
					Quantity:      1,
				},
				orderID: "order-id",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewOrderItem(tt.item, tt.product))
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
					FulfillmentID: "fulfillment-id",
					ProductID:     "product-id",
					Price:         400,
					Quantity:      1,
				},
				orderID: "order-id",
			},
			expect: &response.OrderItem{
				FulfillmentID: "fulfillment-id",
				ProductID:     "product-id",
				Price:         400,
				Quantity:      1,
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
		name     string
		items    entity.OrderItems
		products map[int64]*Product
		expect   OrderItems
	}{
		{
			name: "success",
			items: entity.OrderItems{
				{
					FulfillmentID:     "fulfillment-id",
					OrderID:           "order-id",
					ProductRevisionID: 1,
					Quantity:          1,
					CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			products: map[int64]*Product{
				1: {
					Product: response.Product{
						ID:              "product-id",
						CoordinatorID:   "coordinator-id",
						ProducerID:      "producer-id",
						CategoryID:      "",
						ProductTypeID:   "product-type-id",
						ProductTagIDs:   []string{"product-tag-id"},
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Status:          int32(ProductStatusForSale),
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{
								URL:         "https://and-period.jp/thumbnail01.png",
								IsThumbnail: true,
								Images:      []*response.Image{},
							},
							{
								URL:         "https://and-period.jp/thumbnail02.png",
								IsThumbnail: false,
								Images:      []*response.Image{},
							},
						},
						Price:             400,
						RecommendedPoint1: "ポイント1",
						RecommendedPoint2: "ポイント2",
						RecommendedPoint3: "ポイント3",
						StorageMethodType: int32(StorageMethodTypeNormal),
						DeliveryType:      int32(DeliveryTypeNormal),
						Box60Rate:         50,
						Box80Rate:         40,
						Box100Rate:        30,
						OriginCity:        "彦根市",
						StartAt:           1640962800,
						EndAt:             1640962800,
					},
				},
			},
			expect: OrderItems{
				{
					OrderItem: response.OrderItem{
						FulfillmentID: "fulfillment-id",
						ProductID:     "product-id",
						Price:         400,
						Quantity:      1,
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
			assert.Equal(t, tt.expect, NewOrderItems(tt.items, tt.products))
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
						FulfillmentID: "fulfillment-id",
						ProductID:     "product-id",
						Price:         400,
						Quantity:      1,
					},
					orderID: "order-id",
				},
			},
			expect: []*response.OrderItem{
				{
					FulfillmentID: "fulfillment-id",
					ProductID:     "product-id",
					Price:         400,
					Quantity:      1,
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
