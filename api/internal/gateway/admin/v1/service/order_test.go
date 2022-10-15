package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		order  *entity.Order
		expect *Order
	}{
		{
			name: "success",
			order: &entity.Order{
				ID:                "order-id",
				UserID:            "user-id",
				ScheduleID:        "schedule-id",
				CoordinatorID:     "coordinator-id",
				PaymentStatus:     entity.PaymentStatusCaptured,
				FulfillmentStatus: entity.FulfillmentStatusFulfilled,
				CancelType:        entity.CancelTypeUnknown,
				CancelReason:      "",
				OrderItems: entity.OrderItems{
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
				OrderPayment: entity.OrderPayment{
					ID:             "payment-id",
					TransactionID:  "transaction-id",
					OrderID:        "order-id",
					PromotionID:    "promotion-id",
					PaymentID:      "payment-id",
					PaymentType:    entity.PaymentTypeCard,
					Subtotal:       100,
					Discount:       0,
					ShippingCharge: 500,
					Tax:            60,
					Total:          660,
					Lastname:       "&.",
					Firstname:      "スタッフ",
					PostalCode:     "1000014",
					Prefecture:     "東京都",
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "+819012345678",
					CreatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				OrderFulfillment: entity.OrderFulfillment{
					ID:              "fulfillment-id",
					OrderID:         "order-id",
					ShippingID:      "shipping-id",
					TrackingNumber:  "",
					ShippingCarrier: entity.ShippingCarrierUnknown,
					ShippingMethod:  entity.DeliveryTypeNormal,
					BoxSize:         entity.ShippingSize60,
					BoxCount:        1,
					WeightTotal:     1000,
					Lastname:        "&.",
					Firstname:       "スタッフ",
					PostalCode:      "1000014",
					Prefecture:      "東京都",
					City:            "千代田区",
					AddressLine1:    "永田町1-7-1",
					AddressLine2:    "",
					PhoneNumber:     "+819012345678",
					CreatedAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				OrderActivities: entity.OrderActivities{
					{
						ID:        "event-id",
						OrderID:   "order-id",
						UserID:    "user-id",
						EventType: entity.OrderEventTypeUnknown,
						Detail:    "支払いが完了しました。",
					},
				},
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &Order{
				Order: response.Order{
					ID:          "order-id",
					UserID:      "user-id",
					ScheduleID:  "schedule-id",
					OrderedAt:   0,
					PaidAt:      0,
					DeliveredAt: 0,
					CanceledAt:  0,
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
				payment: &OrderPayment{
					OrderPayment: response.OrderPayment{
						TransactionID:  "transaction-id",
						PromotionID:    "promotion-id",
						PaymentID:      "payment-id",
						PaymentType:    PaymentTypeCard.Response(),
						Status:         PaymentStatusPaid.Response(),
						Subtotal:       100,
						Discount:       0,
						ShippingCharge: 500,
						Tax:            60,
						Total:          660,
						Lastname:       "&.",
						Firstname:      "スタッフ",
						PostalCode:     "1000014",
						Prefecture:     "東京都",
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
					id:      "payment-id",
					orderID: "order-id",
				},
				fulfillment: &OrderFulfillment{
					OrderFulfillment: response.OrderFulfillment{
						TrackingNumber:  "",
						Status:          FulfillmentStatusFulfilled.Response(),
						ShippingCarrier: ShippingCarrierUnknown.Response(),
						ShippingMethod:  DeliveryTypeNormal.Response(),
						BoxSize:         ShippingSize60.Response(),
						BoxCount:        1,
						WeightTotal:     1.0,
						Lastname:        "&.",
						Firstname:       "スタッフ",
						PostalCode:      "1000014",
						Prefecture:      "東京都",
						City:            "千代田区",
						AddressLine1:    "永田町1-7-1",
						AddressLine2:    "",
						PhoneNumber:     "+819012345678",
					},
					id:         "fulfillment-id",
					orderID:    "order-id",
					shippingID: "shipping-id",
				},
				refund: &OrderRefund{
					OrderRefund: response.OrderRefund{
						Canceled: false,
						Type:     OrderRefundTypeUnknown.Response(),
						Reason:   "",
					},
				},
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
				CoordinatorID: "coordinator-id",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, NewOrder(tt.order))
		})
	}
}

func TestOrder_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		order    *Order
		user     *User
		products map[string]*Product
		expect   *Order
	}{
		{
			name: "success",
			order: &Order{
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
			},
			user: &User{
				response.User{
					ID:            "user-id",
					Lastname:      "&.",
					Firstname:     "スタッフ",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "すたっふ",
					Registered:    true,
					Email:         "test-user@and-period.jp",
					PhoneNumber:   "+819012345678",
					PostalCode:    "1000014",
					Prefecture:    "東京都",
					City:          "千代田区",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
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
					},
				},
			},
			expect: &Order{
				Order: response.Order{
					UserName: "&. スタッフ",
				},
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
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.order.Fill(tt.user, tt.products)
			assert.Equal(t, tt.expect, tt.order)
		})
	}
}

func TestOrder_ProductIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		order  *Order
		expect []string
	}{
		{
			name: "success",
			order: &Order{
				Order: response.Order{
					ID:          "order-id",
					UserID:      "user-id",
					UserName:    "&. スタッフ",
					OrderedAt:   0,
					PaidAt:      0,
					DeliveredAt: 0,
					CanceledAt:  0,
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
				payment: &OrderPayment{
					OrderPayment: response.OrderPayment{
						TransactionID:  "transaction-id",
						PromotionID:    "promotion-id",
						PaymentID:      "payment-id",
						PaymentType:    PaymentTypeCard.Response(),
						Status:         PaymentStatusPaid.Response(),
						Subtotal:       100,
						Discount:       0,
						ShippingCharge: 500,
						Tax:            60,
						Total:          660,
						Lastname:       "&.",
						Firstname:      "スタッフ",
						PostalCode:     "1000014",
						Prefecture:     "東京都",
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
					id:      "payment-id",
					orderID: "order-id",
				},
				fulfillment: &OrderFulfillment{
					OrderFulfillment: response.OrderFulfillment{
						TrackingNumber:  "",
						Status:          FulfillmentStatusFulfilled.Response(),
						ShippingCarrier: ShippingCarrierUnknown.Response(),
						ShippingMethod:  DeliveryTypeNormal.Response(),
						BoxSize:         ShippingSize60.Response(),
						BoxCount:        1,
						WeightTotal:     1.0,
						Lastname:        "&.",
						Firstname:       "スタッフ",
						PostalCode:      "1000014",
						Prefecture:      "東京都",
						City:            "千代田区",
						AddressLine1:    "永田町1-7-1",
						AddressLine2:    "",
						PhoneNumber:     "+819012345678",
					},
					id:         "fulfillment-id",
					orderID:    "order-id",
					shippingID: "shipping-id",
				},
				refund: &OrderRefund{
					OrderRefund: response.OrderRefund{
						Canceled: false,
						Type:     OrderRefundTypeUnknown.Response(),
						Reason:   "",
					},
				},
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
			},
			expect: []string{"product-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.order.ProductIDs())
		})
	}
}

func TestOrder_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		order  *Order
		expect *response.Order
	}{
		{
			name: "success",
			order: &Order{
				Order: response.Order{
					ID:          "order-id",
					UserID:      "user-id",
					UserName:    "&. スタッフ",
					OrderedAt:   0,
					PaidAt:      0,
					DeliveredAt: 0,
					CanceledAt:  0,
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
				payment: &OrderPayment{
					OrderPayment: response.OrderPayment{
						TransactionID:  "transaction-id",
						PromotionID:    "promotion-id",
						PaymentID:      "payment-id",
						PaymentType:    PaymentTypeCard.Response(),
						Status:         PaymentStatusPaid.Response(),
						Subtotal:       100,
						Discount:       0,
						ShippingCharge: 500,
						Tax:            60,
						Total:          660,
						Lastname:       "&.",
						Firstname:      "スタッフ",
						PostalCode:     "1000014",
						Prefecture:     "東京都",
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
					id:      "payment-id",
					orderID: "order-id",
				},
				fulfillment: &OrderFulfillment{
					OrderFulfillment: response.OrderFulfillment{
						TrackingNumber:  "",
						Status:          FulfillmentStatusFulfilled.Response(),
						ShippingCarrier: ShippingCarrierUnknown.Response(),
						ShippingMethod:  DeliveryTypeNormal.Response(),
						BoxSize:         ShippingSize60.Response(),
						BoxCount:        1,
						WeightTotal:     1.0,
						Lastname:        "&.",
						Firstname:       "スタッフ",
						PostalCode:      "1000014",
						Prefecture:      "東京都",
						City:            "千代田区",
						AddressLine1:    "永田町1-7-1",
						AddressLine2:    "",
						PhoneNumber:     "+819012345678",
					},
					id:         "fulfillment-id",
					orderID:    "order-id",
					shippingID: "shipping-id",
				},
				refund: &OrderRefund{
					OrderRefund: response.OrderRefund{
						Canceled: false,
						Type:     OrderRefundTypeUnknown.Response(),
						Reason:   "",
					},
				},
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
			},
			expect: &response.Order{
				ID:       "order-id",
				UserID:   "user-id",
				UserName: "&. スタッフ",
				Payment: &response.OrderPayment{
					TransactionID:  "transaction-id",
					PromotionID:    "promotion-id",
					PaymentID:      "payment-id",
					PaymentType:    PaymentTypeCard.Response(),
					Status:         PaymentStatusPaid.Response(),
					Subtotal:       100,
					Discount:       0,
					ShippingCharge: 500,
					Tax:            60,
					Total:          660,
					Lastname:       "&.",
					Firstname:      "スタッフ",
					PostalCode:     "1000014",
					Prefecture:     "東京都",
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "+819012345678",
				},
				Fulfillment: &response.OrderFulfillment{
					TrackingNumber:  "",
					Status:          FulfillmentStatusFulfilled.Response(),
					ShippingCarrier: ShippingCarrierUnknown.Response(),
					ShippingMethod:  DeliveryTypeNormal.Response(),
					BoxSize:         ShippingSize60.Response(),
					BoxCount:        1,
					WeightTotal:     1.0,
					Lastname:        "&.",
					Firstname:       "スタッフ",
					PostalCode:      "1000014",
					Prefecture:      "東京都",
					City:            "千代田区",
					AddressLine1:    "永田町1-7-1",
					AddressLine2:    "",
					PhoneNumber:     "+819012345678",
				},
				Refund: &response.OrderRefund{
					Canceled: false,
					Type:     OrderRefundTypeUnknown.Response(),
					Reason:   "",
				},
				Items: []*response.OrderItem{
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
				OrderedAt:   0,
				PaidAt:      0,
				DeliveredAt: 0,
				CanceledAt:  0,
				CreatedAt:   1640962800,
				UpdatedAt:   1640962800,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.order.Response())
		})
	}
}

func TestOrders(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders entity.Orders
		expect Orders
	}{
		{
			name: "success",
			orders: entity.Orders{
				{
					ID:                "order-id",
					UserID:            "user-id",
					PaymentStatus:     entity.PaymentStatusCaptured,
					FulfillmentStatus: entity.FulfillmentStatusFulfilled,
					CancelType:        entity.CancelTypeUnknown,
					CancelReason:      "",
					OrderItems: entity.OrderItems{
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
					OrderPayment: entity.OrderPayment{
						ID:             "payment-id",
						TransactionID:  "transaction-id",
						OrderID:        "order-id",
						PromotionID:    "promotion-id",
						PaymentID:      "payment-id",
						PaymentType:    entity.PaymentTypeCard,
						Subtotal:       100,
						Discount:       0,
						ShippingCharge: 500,
						Tax:            60,
						Total:          660,
						Lastname:       "&.",
						Firstname:      "スタッフ",
						PostalCode:     "1000014",
						Prefecture:     "東京都",
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
						CreatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
						UpdatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
					},
					OrderFulfillment: entity.OrderFulfillment{
						ID:              "fulfillment-id",
						OrderID:         "order-id",
						ShippingID:      "shipping-id",
						TrackingNumber:  "",
						ShippingCarrier: entity.ShippingCarrierUnknown,
						ShippingMethod:  entity.DeliveryTypeNormal,
						BoxSize:         entity.ShippingSize60,
						BoxCount:        1,
						WeightTotal:     1000,
						Lastname:        "&.",
						Firstname:       "スタッフ",
						PostalCode:      "1000014",
						Prefecture:      "東京都",
						City:            "千代田区",
						AddressLine1:    "永田町1-7-1",
						AddressLine2:    "",
						PhoneNumber:     "+819012345678",
						CreatedAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
						UpdatedAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
					},
					OrderActivities: entity.OrderActivities{
						{
							ID:        "event-id",
							OrderID:   "order-id",
							UserID:    "user-id",
							EventType: entity.OrderEventTypeUnknown,
							Detail:    "支払いが完了しました。",
						},
					},
					CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: Orders{
				{
					Order: response.Order{
						ID:          "order-id",
						UserID:      "user-id",
						OrderedAt:   0,
						PaidAt:      0,
						DeliveredAt: 0,
						CanceledAt:  0,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
					payment: &OrderPayment{
						OrderPayment: response.OrderPayment{
							TransactionID:  "transaction-id",
							PromotionID:    "promotion-id",
							PaymentID:      "payment-id",
							PaymentType:    PaymentTypeCard.Response(),
							Status:         PaymentStatusPaid.Response(),
							Subtotal:       100,
							Discount:       0,
							ShippingCharge: 500,
							Tax:            60,
							Total:          660,
							Lastname:       "&.",
							Firstname:      "スタッフ",
							PostalCode:     "1000014",
							Prefecture:     "東京都",
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "+819012345678",
						},
						id:      "payment-id",
						orderID: "order-id",
					},
					fulfillment: &OrderFulfillment{
						OrderFulfillment: response.OrderFulfillment{
							TrackingNumber:  "",
							Status:          FulfillmentStatusFulfilled.Response(),
							ShippingCarrier: ShippingCarrierUnknown.Response(),
							ShippingMethod:  DeliveryTypeNormal.Response(),
							BoxSize:         ShippingSize60.Response(),
							BoxCount:        1,
							WeightTotal:     1.0,
							Lastname:        "&.",
							Firstname:       "スタッフ",
							PostalCode:      "1000014",
							Prefecture:      "東京都",
							City:            "千代田区",
							AddressLine1:    "永田町1-7-1",
							AddressLine2:    "",
							PhoneNumber:     "+819012345678",
						},
						id:         "fulfillment-id",
						orderID:    "order-id",
						shippingID: "shipping-id",
					},
					refund: &OrderRefund{
						OrderRefund: response.OrderRefund{
							Canceled: false,
							Type:     OrderRefundTypeUnknown.Response(),
							Reason:   "",
						},
					},
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
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, NewOrders(tt.orders))
		})
	}
}

func TestOrders_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		orders   Orders
		users    map[string]*User
		products map[string]*Product
		expect   Orders
	}{
		{
			name: "success",
			orders: Orders{
				{
					Order: response.Order{
						UserID: "user-id",
					},
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
				},
			},
			users: map[string]*User{
				"user-id": {
					response.User{
						ID:            "user-id",
						Lastname:      "&.",
						Firstname:     "スタッフ",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "すたっふ",
						Registered:    true,
						Email:         "test-user@and-period.jp",
						PhoneNumber:   "+819012345678",
						PostalCode:    "1000014",
						Prefecture:    "東京都",
						City:          "千代田区",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
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
					},
				},
			},
			expect: Orders{
				{
					Order: response.Order{
						UserID:   "user-id",
						UserName: "&. スタッフ",
					},
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
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.orders.Fill(tt.users, tt.products)
			assert.Equal(t, tt.expect, tt.orders)
		})
	}
}

func TestOrders_UserIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders Orders
		expect []string
	}{
		{
			name: "success",
			orders: Orders{
				{
					Order: response.Order{
						ID:          "order-id",
						UserID:      "user-id",
						UserName:    "&. スタッフ",
						OrderedAt:   0,
						PaidAt:      0,
						DeliveredAt: 0,
						CanceledAt:  0,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
					payment: &OrderPayment{
						OrderPayment: response.OrderPayment{
							TransactionID:  "transaction-id",
							PromotionID:    "promotion-id",
							PaymentID:      "payment-id",
							PaymentType:    PaymentTypeCard.Response(),
							Status:         PaymentStatusPaid.Response(),
							Subtotal:       100,
							Discount:       0,
							ShippingCharge: 500,
							Tax:            60,
							Total:          660,
							Lastname:       "&.",
							Firstname:      "スタッフ",
							PostalCode:     "1000014",
							Prefecture:     "東京都",
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "+819012345678",
						},
						id:      "payment-id",
						orderID: "order-id",
					},
					fulfillment: &OrderFulfillment{
						OrderFulfillment: response.OrderFulfillment{
							TrackingNumber:  "",
							Status:          FulfillmentStatusFulfilled.Response(),
							ShippingCarrier: ShippingCarrierUnknown.Response(),
							ShippingMethod:  DeliveryTypeNormal.Response(),
							BoxSize:         ShippingSize60.Response(),
							BoxCount:        1,
							WeightTotal:     1.0,
							Lastname:        "&.",
							Firstname:       "スタッフ",
							PostalCode:      "1000014",
							Prefecture:      "東京都",
							City:            "千代田区",
							AddressLine1:    "永田町1-7-1",
							AddressLine2:    "",
							PhoneNumber:     "+819012345678",
						},
						id:         "fulfillment-id",
						orderID:    "order-id",
						shippingID: "shipping-id",
					},
					refund: &OrderRefund{
						OrderRefund: response.OrderRefund{
							Canceled: false,
							Type:     OrderRefundTypeUnknown.Response(),
							Reason:   "",
						},
					},
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
				},
			},
			expect: []string{"user-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.orders.UserIDs())
		})
	}
}

func TestOrders_ProductIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders Orders
		expect []string
	}{
		{
			name: "success",
			orders: Orders{
				{
					Order: response.Order{
						ID:          "order-id",
						UserID:      "user-id",
						UserName:    "&. スタッフ",
						OrderedAt:   0,
						PaidAt:      0,
						DeliveredAt: 0,
						CanceledAt:  0,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
					payment: &OrderPayment{
						OrderPayment: response.OrderPayment{
							TransactionID:  "transaction-id",
							PromotionID:    "promotion-id",
							PaymentID:      "payment-id",
							PaymentType:    PaymentTypeCard.Response(),
							Status:         PaymentStatusPaid.Response(),
							Subtotal:       100,
							Discount:       0,
							ShippingCharge: 500,
							Tax:            60,
							Total:          660,
							Lastname:       "&.",
							Firstname:      "スタッフ",
							PostalCode:     "1000014",
							Prefecture:     "東京都",
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "+819012345678",
						},
						id:      "payment-id",
						orderID: "order-id",
					},
					fulfillment: &OrderFulfillment{
						OrderFulfillment: response.OrderFulfillment{
							TrackingNumber:  "",
							Status:          FulfillmentStatusFulfilled.Response(),
							ShippingCarrier: ShippingCarrierUnknown.Response(),
							ShippingMethod:  DeliveryTypeNormal.Response(),
							BoxSize:         ShippingSize60.Response(),
							BoxCount:        1,
							WeightTotal:     1.0,
							Lastname:        "&.",
							Firstname:       "スタッフ",
							PostalCode:      "1000014",
							Prefecture:      "東京都",
							City:            "千代田区",
							AddressLine1:    "永田町1-7-1",
							AddressLine2:    "",
							PhoneNumber:     "+819012345678",
						},
						id:         "fulfillment-id",
						orderID:    "order-id",
						shippingID: "shipping-id",
					},
					refund: &OrderRefund{
						OrderRefund: response.OrderRefund{
							Canceled: false,
							Type:     OrderRefundTypeUnknown.Response(),
							Reason:   "",
						},
					},
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
				},
			},
			expect: []string{"product-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.orders.ProductIDs())
		})
	}
}

func TestOrders_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders Orders
		expect []*response.Order
	}{
		{
			name: "success",
			orders: Orders{
				{
					Order: response.Order{
						ID:          "order-id",
						UserID:      "user-id",
						UserName:    "&. スタッフ",
						OrderedAt:   0,
						PaidAt:      0,
						DeliveredAt: 0,
						CanceledAt:  0,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
					payment: &OrderPayment{
						OrderPayment: response.OrderPayment{
							TransactionID:  "transaction-id",
							PromotionID:    "promotion-id",
							PaymentID:      "payment-id",
							PaymentType:    PaymentTypeCard.Response(),
							Status:         PaymentStatusPaid.Response(),
							Subtotal:       100,
							Discount:       0,
							ShippingCharge: 500,
							Tax:            60,
							Total:          660,
							Lastname:       "&.",
							Firstname:      "スタッフ",
							PostalCode:     "1000014",
							Prefecture:     "東京都",
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "+819012345678",
						},
						id:      "payment-id",
						orderID: "order-id",
					},
					fulfillment: &OrderFulfillment{
						OrderFulfillment: response.OrderFulfillment{
							TrackingNumber:  "",
							Status:          FulfillmentStatusFulfilled.Response(),
							ShippingCarrier: ShippingCarrierUnknown.Response(),
							ShippingMethod:  DeliveryTypeNormal.Response(),
							BoxSize:         ShippingSize60.Response(),
							BoxCount:        1,
							WeightTotal:     1.0,
							Lastname:        "&.",
							Firstname:       "スタッフ",
							PostalCode:      "1000014",
							Prefecture:      "東京都",
							City:            "千代田区",
							AddressLine1:    "永田町1-7-1",
							AddressLine2:    "",
							PhoneNumber:     "+819012345678",
						},
						id:         "fulfillment-id",
						orderID:    "order-id",
						shippingID: "shipping-id",
					},
					refund: &OrderRefund{
						OrderRefund: response.OrderRefund{
							Canceled: false,
							Type:     OrderRefundTypeUnknown.Response(),
							Reason:   "",
						},
					},
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
				},
			},
			expect: []*response.Order{
				{
					ID:       "order-id",
					UserID:   "user-id",
					UserName: "&. スタッフ",
					Payment: &response.OrderPayment{
						TransactionID:  "transaction-id",
						PromotionID:    "promotion-id",
						PaymentID:      "payment-id",
						PaymentType:    PaymentTypeCard.Response(),
						Status:         PaymentStatusPaid.Response(),
						Subtotal:       100,
						Discount:       0,
						ShippingCharge: 500,
						Tax:            60,
						Total:          660,
						Lastname:       "&.",
						Firstname:      "スタッフ",
						PostalCode:     "1000014",
						Prefecture:     "東京都",
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
					Fulfillment: &response.OrderFulfillment{
						TrackingNumber:  "",
						Status:          FulfillmentStatusFulfilled.Response(),
						ShippingCarrier: ShippingCarrierUnknown.Response(),
						ShippingMethod:  DeliveryTypeNormal.Response(),
						BoxSize:         ShippingSize60.Response(),
						BoxCount:        1,
						WeightTotal:     1.0,
						Lastname:        "&.",
						Firstname:       "スタッフ",
						PostalCode:      "1000014",
						Prefecture:      "東京都",
						City:            "千代田区",
						AddressLine1:    "永田町1-7-1",
						AddressLine2:    "",
						PhoneNumber:     "+819012345678",
					},
					Refund: &response.OrderRefund{
						Canceled: false,
						Type:     OrderRefundTypeUnknown.Response(),
						Reason:   "",
					},
					Items: []*response.OrderItem{
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
					OrderedAt:   0,
					PaidAt:      0,
					DeliveredAt: 0,
					CanceledAt:  0,
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.orders.Response())
		})
	}
}
