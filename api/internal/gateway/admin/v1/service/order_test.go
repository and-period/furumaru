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
				CoordinatorID:     "coordinator-id",
				ScheduleID:        "schedule-id",
				PromotionID:       "",
				PaymentStatus:     entity.PaymentStatusPending,
				FulfillmentStatus: entity.FulfillmentStatusUnfulfilled,
				RefundReason:      "",
				CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
				OrderItems: entity.OrderItems{
					{
						OrderID:   "order-id",
						ProductID: "product-id01",
						Price:     100,
						Quantity:  1,
						CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
						UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					},
					{
						OrderID:   "order-id",
						ProductID: "product-id02",
						Price:     500,
						Quantity:  2,
						CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
						UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					},
				},
				Payment: entity.Payment{
					OrderID:       "order-id",
					AddressID:     "address-id",
					TransactionID: "transaction-id",
					MethodType:    entity.PaymentMethodTypeCreditCard,
					Subtotal:      1100,
					Discount:      0,
					ShippingFee:   500,
					Tax:           160,
					Total:         1760,
					CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				Fulfillment: entity.Fulfillment{
					OrderID:         "order-id",
					AddressID:       "address-id",
					TrackingNumber:  "",
					ShippingCarrier: entity.ShippingCarrierUnknown,
					ShippingMethod:  entity.DeliveryTypeNormal,
					BoxSize:         entity.ShippingSize60,
					CreatedAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				Activities: []*entity.Activity{
					{
						ID:        "event-id",
						OrderID:   "order-id",
						UserID:    "user-id",
						EventType: entity.OrderEventTypeUnknown,
						Detail:    "支払いが完了しました。",
					},
				},
			},
			expect: &Order{
				Order: response.Order{
					ID:          "order-id",
					UserID:      "user-id",
					ScheduleID:  "schedule-id",
					PromotionID: "",
					OrderedAt:   0,
					PaidAt:      0,
					DeliveredAt: 0,
					CanceledAt:  0,
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
				payment: &OrderPayment{
					OrderPayment: response.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    PaymentMethodTypeCreditCard.Response(),
						Status:        PaymentStatusPending.Response(),
						Subtotal:      1100,
						Discount:      0,
						ShippingFee:   500,
						Tax:           160,
						Total:         1760,
						AddressID:     "address-id",
					},
					orderID: "order-id",
				},
				fulfillment: &OrderFulfillment{
					OrderFulfillment: response.OrderFulfillment{
						TrackingNumber:  "",
						Status:          FulfillmentStatusUnfulfilled.Response(),
						ShippingCarrier: ShippingCarrierUnknown.Response(),
						ShippingMethod:  DeliveryTypeNormal.Response(),
						BoxSize:         ShippingSize60.Response(),
						AddressID:       "address-id",
					},
					orderID: "order-id",
				},
				refund: &OrderRefund{
					OrderRefund: response.OrderRefund{
						Canceled: false,
						Reason:   "",
						Total:    0,
					},
				},
				items: OrderItems{
					{
						OrderItem: response.OrderItem{
							ProductID: "product-id01",
							Price:     100,
							Quantity:  1,
						},
						orderID: "order-id",
					},
					{
						OrderItem: response.OrderItem{
							ProductID: "product-id02",
							Price:     500,
							Quantity:  2,
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
		name      string
		order     *Order
		user      *User
		products  map[string]*Product
		addresses map[string]*Address
		expect    *Order
	}{
		{
			name: "success",
			order: &Order{
				Order: response.Order{
					ID:          "order-id",
					UserID:      "user-id",
					ScheduleID:  "schedule-id",
					PromotionID: "",
					OrderedAt:   0,
					PaidAt:      0,
					DeliveredAt: 0,
					CanceledAt:  0,
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
				payment: &OrderPayment{
					OrderPayment: response.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    PaymentMethodTypeCreditCard.Response(),
						Status:        PaymentStatusPaid.Response(),
						Subtotal:      100,
						Discount:      0,
						ShippingFee:   500,
						Tax:           60,
						Total:         660,
						AddressID:     "address-id",
					},
					orderID: "order-id",
				},
				fulfillment: &OrderFulfillment{
					OrderFulfillment: response.OrderFulfillment{
						TrackingNumber:  "",
						Status:          FulfillmentStatusFulfilled.Response(),
						ShippingCarrier: ShippingCarrierUnknown.Response(),
						ShippingMethod:  DeliveryTypeNormal.Response(),
						BoxSize:         ShippingSize60.Response(),
						AddressID:       "address-id",
					},
					orderID: "order-id",
				},
				refund: &OrderRefund{
					OrderRefund: response.OrderRefund{
						Canceled: false,
						Reason:   "",
						Total:    0,
					},
				},
				items: OrderItems{
					{
						OrderItem: response.OrderItem{
							ProductID: "product-id",
							Price:     100,
							Quantity:  1,
						},
						orderID: "order-id",
					},
				},
				CoordinatorID: "coordinator-id",
			},
			user: &User{
				response.User{
					ID:            "user-id",
					Lastname:      "&.",
					Firstname:     "購入者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "こうにゅうしゃ",
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
						ProductTypeID:   "product-type-id",
						CategoryID:      "category-id",
						ProducerID:      "producer-id",
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
			addresses: map[string]*Address{
				"address-id": {
					Address: response.Address{
						Lastname:     "&.",
						Firstname:    "購入者",
						PostalCode:   "1000014",
						Prefecture:   "東京都",
						City:         "千代田区",
						AddressLine1: "永田町1-7-1",
						AddressLine2: "",
						PhoneNumber:  "+819012345678",
					},
					id: "address-id",
				},
			},
			expect: &Order{
				Order: response.Order{
					ID:          "order-id",
					UserID:      "user-id",
					UserName:    "&. 購入者",
					ScheduleID:  "schedule-id",
					PromotionID: "",
					OrderedAt:   0,
					PaidAt:      0,
					DeliveredAt: 0,
					CanceledAt:  0,
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
				payment: &OrderPayment{
					OrderPayment: response.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    PaymentMethodTypeCreditCard.Response(),
						Status:        PaymentStatusPaid.Response(),
						Subtotal:      100,
						Discount:      0,
						ShippingFee:   500,
						Tax:           60,
						Total:         660,
						AddressID:     "address-id",
						Address: &response.Address{
							Lastname:     "&.",
							Firstname:    "購入者",
							PostalCode:   "1000014",
							Prefecture:   "東京都",
							City:         "千代田区",
							AddressLine1: "永田町1-7-1",
							AddressLine2: "",
							PhoneNumber:  "+819012345678",
						},
					},
					orderID: "order-id",
				},
				fulfillment: &OrderFulfillment{
					OrderFulfillment: response.OrderFulfillment{
						TrackingNumber:  "",
						Status:          FulfillmentStatusFulfilled.Response(),
						ShippingCarrier: ShippingCarrierUnknown.Response(),
						ShippingMethod:  DeliveryTypeNormal.Response(),
						BoxSize:         ShippingSize60.Response(),
						AddressID:       "address-id",
						Address: &response.Address{
							Lastname:     "&.",
							Firstname:    "購入者",
							PostalCode:   "1000014",
							Prefecture:   "東京都",
							City:         "千代田区",
							AddressLine1: "永田町1-7-1",
							AddressLine2: "",
							PhoneNumber:  "+819012345678",
						},
					},
					orderID: "order-id",
				},
				refund: &OrderRefund{
					OrderRefund: response.OrderRefund{
						Canceled: false,
						Reason:   "",
						Total:    0,
					},
				},
				items: OrderItems{
					{
						OrderItem: response.OrderItem{
							ProductID: "product-id",
							Name:      "新鮮なじゃがいも",
							Price:     100,
							Quantity:  1,
							Weight:    1.3,
							Media: []*response.ProductMedia{
								{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
								{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
							},
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
			t.Parallel()
			tt.order.Fill(tt.user, tt.products, tt.addresses)
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
					ScheduleID:  "schedule-id",
					PromotionID: "",
					OrderedAt:   0,
					PaidAt:      0,
					DeliveredAt: 0,
					CanceledAt:  0,
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
				payment: &OrderPayment{
					OrderPayment: response.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    PaymentMethodTypeCreditCard.Response(),
						Status:        PaymentStatusPaid.Response(),
						Subtotal:      100,
						Discount:      0,
						ShippingFee:   500,
						Tax:           60,
						Total:         660,
						AddressID:     "address-id",
					},
					orderID: "order-id",
				},
				fulfillment: &OrderFulfillment{
					OrderFulfillment: response.OrderFulfillment{
						TrackingNumber:  "",
						Status:          FulfillmentStatusFulfilled.Response(),
						ShippingCarrier: ShippingCarrierUnknown.Response(),
						ShippingMethod:  DeliveryTypeNormal.Response(),
						BoxSize:         ShippingSize60.Response(),
						AddressID:       "address-id",
					},
					orderID: "order-id",
				},
				refund: &OrderRefund{
					OrderRefund: response.OrderRefund{
						Canceled: false,
						Reason:   "",
						Total:    0,
					},
				},
				items: OrderItems{
					{
						OrderItem: response.OrderItem{
							ProductID: "product-id01",
							Price:     100,
							Quantity:  1,
						},
						orderID: "order-id",
					},
					{
						OrderItem: response.OrderItem{
							ProductID: "product-id02",
							Price:     500,
							Quantity:  2,
						},
						orderID: "order-id",
					},
				},
				CoordinatorID: "coordinator-id",
			},
			expect: []string{"product-id01", "product-id02"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.expect, tt.order.ProductIDs())
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
					UserName:    "&. 購入者",
					ScheduleID:  "schedule-id",
					PromotionID: "",
					OrderedAt:   0,
					PaidAt:      0,
					DeliveredAt: 0,
					CanceledAt:  0,
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
				payment: &OrderPayment{
					OrderPayment: response.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    PaymentMethodTypeCreditCard.Response(),
						Status:        PaymentStatusPaid.Response(),
						Subtotal:      100,
						Discount:      0,
						ShippingFee:   500,
						Tax:           60,
						Total:         660,
						AddressID:     "address-id",
						Address: &response.Address{
							Lastname:     "&.",
							Firstname:    "購入者",
							PostalCode:   "1000014",
							Prefecture:   "東京都",
							City:         "千代田区",
							AddressLine1: "永田町1-7-1",
							AddressLine2: "",
							PhoneNumber:  "+819012345678",
						},
					},
					orderID: "order-id",
				},
				fulfillment: &OrderFulfillment{
					OrderFulfillment: response.OrderFulfillment{
						TrackingNumber:  "",
						Status:          FulfillmentStatusFulfilled.Response(),
						ShippingCarrier: ShippingCarrierUnknown.Response(),
						ShippingMethod:  DeliveryTypeNormal.Response(),
						BoxSize:         ShippingSize60.Response(),
						AddressID:       "address-id",
						Address: &response.Address{
							Lastname:     "&.",
							Firstname:    "購入者",
							PostalCode:   "1000014",
							Prefecture:   "東京都",
							City:         "千代田区",
							AddressLine1: "永田町1-7-1",
							AddressLine2: "",
							PhoneNumber:  "+819012345678",
						},
					},
					orderID: "order-id",
				},
				refund: &OrderRefund{
					OrderRefund: response.OrderRefund{
						Canceled: false,
						Reason:   "",
						Total:    0,
					},
				},
				items: OrderItems{
					{
						OrderItem: response.OrderItem{
							ProductID: "product-id",
							Name:      "新鮮なじゃがいも",
							Price:     100,
							Quantity:  1,
							Weight:    1.3,
							Media: []*response.ProductMedia{
								{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
								{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
							},
						},
						orderID: "order-id",
					},
				},
				CoordinatorID: "coordinator-id",
			},
			expect: &response.Order{
				ID:          "order-id",
				UserID:      "user-id",
				UserName:    "&. 購入者",
				ScheduleID:  "schedule-id",
				PromotionID: "",
				Payment: &response.OrderPayment{
					TransactionID: "transaction-id",
					MethodType:    PaymentMethodTypeCreditCard.Response(),
					Status:        PaymentStatusPaid.Response(),
					Subtotal:      100,
					Discount:      0,
					ShippingFee:   500,
					Tax:           60,
					Total:         660,
					AddressID:     "address-id",
					Address: &response.Address{
						Lastname:     "&.",
						Firstname:    "購入者",
						PostalCode:   "1000014",
						Prefecture:   "東京都",
						City:         "千代田区",
						AddressLine1: "永田町1-7-1",
						AddressLine2: "",
						PhoneNumber:  "+819012345678",
					},
				},
				Fulfillment: &response.OrderFulfillment{
					TrackingNumber:  "",
					Status:          FulfillmentStatusFulfilled.Response(),
					ShippingCarrier: ShippingCarrierUnknown.Response(),
					ShippingMethod:  DeliveryTypeNormal.Response(),
					BoxSize:         ShippingSize60.Response(),
					AddressID:       "address-id",
					Address: &response.Address{
						Lastname:     "&.",
						Firstname:    "購入者",
						PostalCode:   "1000014",
						Prefecture:   "東京都",
						City:         "千代田区",
						AddressLine1: "永田町1-7-1",
						AddressLine2: "",
						PhoneNumber:  "+819012345678",
					},
				},
				Refund: &response.OrderRefund{
					Canceled: false,
					Reason:   "",
					Total:    0,
				},
				Items: []*response.OrderItem{
					{
						ProductID: "product-id",
						Name:      "新鮮なじゃがいも",
						Price:     100,
						Quantity:  1,
						Weight:    1.3,
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
					CoordinatorID:     "coordinator-id",
					ScheduleID:        "schedule-id",
					PromotionID:       "",
					PaymentStatus:     entity.PaymentStatusPending,
					FulfillmentStatus: entity.FulfillmentStatusUnfulfilled,
					RefundReason:      "",
					CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
					OrderItems: []*entity.OrderItem{
						{
							OrderID:   "order-id",
							ProductID: "product-id01",
							Price:     100,
							Quantity:  1,
							CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
							UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
						},
						{
							OrderID:   "order-id",
							ProductID: "product-id02",
							Price:     500,
							Quantity:  2,
							CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
							UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
						},
					},
					Payment: entity.Payment{
						OrderID:       "order-id",
						AddressID:     "address-id",
						TransactionID: "transaction-id",
						MethodType:    entity.PaymentMethodTypeCreditCard,
						Subtotal:      1100,
						Discount:      0,
						ShippingFee:   500,
						Tax:           160,
						Total:         1760,
						CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
						UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
					},
					Fulfillment: entity.Fulfillment{
						OrderID:         "order-id",
						AddressID:       "address-id",
						TrackingNumber:  "",
						ShippingCarrier: entity.ShippingCarrierUnknown,
						ShippingMethod:  entity.DeliveryTypeNormal,
						BoxSize:         entity.ShippingSize60,
						CreatedAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
						UpdatedAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
					},
					Activities: []*entity.Activity{
						{
							ID:        "event-id",
							OrderID:   "order-id",
							UserID:    "user-id",
							EventType: entity.OrderEventTypeUnknown,
							Detail:    "支払いが完了しました。",
						},
					},
				},
			},
			expect: Orders{
				{
					Order: response.Order{
						ID:          "order-id",
						UserID:      "user-id",
						ScheduleID:  "schedule-id",
						PromotionID: "",
						OrderedAt:   0,
						PaidAt:      0,
						DeliveredAt: 0,
						CanceledAt:  0,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
					payment: &OrderPayment{
						OrderPayment: response.OrderPayment{
							TransactionID: "transaction-id",
							MethodType:    PaymentMethodTypeCreditCard.Response(),
							Status:        PaymentStatusPending.Response(),
							Subtotal:      1100,
							Discount:      0,
							ShippingFee:   500,
							Tax:           160,
							Total:         1760,
							AddressID:     "address-id",
						},
						orderID: "order-id",
					},
					fulfillment: &OrderFulfillment{
						OrderFulfillment: response.OrderFulfillment{
							TrackingNumber:  "",
							Status:          FulfillmentStatusUnfulfilled.Response(),
							ShippingCarrier: ShippingCarrierUnknown.Response(),
							ShippingMethod:  DeliveryTypeNormal.Response(),
							BoxSize:         ShippingSize60.Response(),
							AddressID:       "address-id",
						},
						orderID: "order-id",
					},
					refund: &OrderRefund{
						OrderRefund: response.OrderRefund{
							Canceled: false,
							Reason:   "",
							Total:    0,
						},
					},
					items: OrderItems{
						{
							OrderItem: response.OrderItem{
								ProductID: "product-id01",
								Price:     100,
								Quantity:  1,
							},
							orderID: "order-id",
						},
						{
							OrderItem: response.OrderItem{
								ProductID: "product-id02",
								Price:     500,
								Quantity:  2,
							},
							orderID: "order-id",
						},
					},
					CoordinatorID: "coordinator-id",
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
		name      string
		orders    Orders
		users     map[string]*User
		products  map[string]*Product
		addresses map[string]*Address
		expect    Orders
	}{
		{
			name: "success",
			orders: Orders{
				{
					Order: response.Order{
						ID:          "order-id",
						UserID:      "user-id",
						ScheduleID:  "schedule-id",
						PromotionID: "",
						OrderedAt:   0,
						PaidAt:      0,
						DeliveredAt: 0,
						CanceledAt:  0,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
					payment: &OrderPayment{
						OrderPayment: response.OrderPayment{
							TransactionID: "transaction-id",
							MethodType:    PaymentMethodTypeCreditCard.Response(),
							Status:        PaymentStatusPaid.Response(),
							Subtotal:      100,
							Discount:      0,
							ShippingFee:   500,
							Tax:           60,
							Total:         660,
							AddressID:     "address-id",
						},
						orderID: "order-id",
					},
					fulfillment: &OrderFulfillment{
						OrderFulfillment: response.OrderFulfillment{
							TrackingNumber:  "",
							Status:          FulfillmentStatusFulfilled.Response(),
							ShippingCarrier: ShippingCarrierUnknown.Response(),
							ShippingMethod:  DeliveryTypeNormal.Response(),
							BoxSize:         ShippingSize60.Response(),
							AddressID:       "address-id",
						},
						orderID: "order-id",
					},
					refund: &OrderRefund{
						OrderRefund: response.OrderRefund{
							Canceled: false,
							Reason:   "",
							Total:    0,
						},
					},
					items: OrderItems{
						{
							OrderItem: response.OrderItem{
								ProductID: "product-id",
								Price:     100,
								Quantity:  1,
							},
							orderID: "order-id",
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			users: map[string]*User{
				"user-id": {
					response.User{
						ID:            "user-id",
						Lastname:      "&.",
						Firstname:     "購入者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "こうにゅうしゃ",
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
						ProductTypeID:   "product-type-id",
						CategoryID:      "category-id",
						ProducerID:      "producer-id",
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
			addresses: map[string]*Address{
				"address-id": {
					Address: response.Address{
						Lastname:     "&.",
						Firstname:    "購入者",
						PostalCode:   "1000014",
						Prefecture:   "東京都",
						City:         "千代田区",
						AddressLine1: "永田町1-7-1",
						AddressLine2: "",
						PhoneNumber:  "+819012345678",
					},
					id: "address-id",
				},
			},
			expect: Orders{
				{
					Order: response.Order{
						ID:          "order-id",
						UserID:      "user-id",
						UserName:    "&. 購入者",
						ScheduleID:  "schedule-id",
						PromotionID: "",
						OrderedAt:   0,
						PaidAt:      0,
						DeliveredAt: 0,
						CanceledAt:  0,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
					payment: &OrderPayment{
						OrderPayment: response.OrderPayment{
							TransactionID: "transaction-id",
							MethodType:    PaymentMethodTypeCreditCard.Response(),
							Status:        PaymentStatusPaid.Response(),
							Subtotal:      100,
							Discount:      0,
							ShippingFee:   500,
							Tax:           60,
							Total:         660,
							AddressID:     "address-id",
							Address: &response.Address{
								Lastname:     "&.",
								Firstname:    "購入者",
								PostalCode:   "1000014",
								Prefecture:   "東京都",
								City:         "千代田区",
								AddressLine1: "永田町1-7-1",
								AddressLine2: "",
								PhoneNumber:  "+819012345678",
							},
						},
						orderID: "order-id",
					},
					fulfillment: &OrderFulfillment{
						OrderFulfillment: response.OrderFulfillment{
							TrackingNumber:  "",
							Status:          FulfillmentStatusFulfilled.Response(),
							ShippingCarrier: ShippingCarrierUnknown.Response(),
							ShippingMethod:  DeliveryTypeNormal.Response(),
							BoxSize:         ShippingSize60.Response(),
							AddressID:       "address-id",
							Address: &response.Address{
								Lastname:     "&.",
								Firstname:    "購入者",
								PostalCode:   "1000014",
								Prefecture:   "東京都",
								City:         "千代田区",
								AddressLine1: "永田町1-7-1",
								AddressLine2: "",
								PhoneNumber:  "+819012345678",
							},
						},
						orderID: "order-id",
					},
					refund: &OrderRefund{
						OrderRefund: response.OrderRefund{
							Canceled: false,
							Reason:   "",
							Total:    0,
						},
					},
					items: OrderItems{
						{
							OrderItem: response.OrderItem{
								ProductID: "product-id",
								Name:      "新鮮なじゃがいも",
								Price:     100,
								Quantity:  1,
								Weight:    1.3,
								Media: []*response.ProductMedia{
									{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
									{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
								},
							},
							orderID: "order-id",
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.orders.Fill(tt.users, tt.products, tt.addresses)
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
						UserName:    "&. 購入者",
						ScheduleID:  "schedule-id",
						PromotionID: "",
						OrderedAt:   0,
						PaidAt:      0,
						DeliveredAt: 0,
						CanceledAt:  0,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
					payment: &OrderPayment{
						OrderPayment: response.OrderPayment{
							TransactionID: "transaction-id",
							MethodType:    PaymentMethodTypeCreditCard.Response(),
							Status:        PaymentStatusPaid.Response(),
							Subtotal:      100,
							Discount:      0,
							ShippingFee:   500,
							Tax:           60,
							Total:         660,
							AddressID:     "address-id",
							Address: &response.Address{
								Lastname:     "&.",
								Firstname:    "購入者",
								PostalCode:   "1000014",
								Prefecture:   "東京都",
								City:         "千代田区",
								AddressLine1: "永田町1-7-1",
								AddressLine2: "",
								PhoneNumber:  "+819012345678",
							},
						},
						orderID: "order-id",
					},
					fulfillment: &OrderFulfillment{
						OrderFulfillment: response.OrderFulfillment{
							TrackingNumber:  "",
							Status:          FulfillmentStatusFulfilled.Response(),
							ShippingCarrier: ShippingCarrierUnknown.Response(),
							ShippingMethod:  DeliveryTypeNormal.Response(),
							BoxSize:         ShippingSize60.Response(),
							AddressID:       "address-id",
							Address: &response.Address{
								Lastname:     "&.",
								Firstname:    "購入者",
								PostalCode:   "1000014",
								Prefecture:   "東京都",
								City:         "千代田区",
								AddressLine1: "永田町1-7-1",
								AddressLine2: "",
								PhoneNumber:  "+819012345678",
							},
						},
						orderID: "order-id",
					},
					refund: &OrderRefund{
						OrderRefund: response.OrderRefund{
							Canceled: false,
							Reason:   "",
							Total:    0,
						},
					},
					items: OrderItems{
						{
							OrderItem: response.OrderItem{
								ProductID: "product-id",
								Name:      "新鮮なじゃがいも",
								Price:     100,
								Quantity:  1,
								Weight:    1.3,
								Media: []*response.ProductMedia{
									{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
									{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
								},
							},
							orderID: "order-id",
						},
					},
					CoordinatorID: "coordinator-id",
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
						UserName:    "&. 購入者",
						ScheduleID:  "schedule-id",
						PromotionID: "",
						OrderedAt:   0,
						PaidAt:      0,
						DeliveredAt: 0,
						CanceledAt:  0,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
					payment: &OrderPayment{
						OrderPayment: response.OrderPayment{
							TransactionID: "transaction-id",
							MethodType:    PaymentMethodTypeCreditCard.Response(),
							Status:        PaymentStatusPaid.Response(),
							Subtotal:      100,
							Discount:      0,
							ShippingFee:   500,
							Tax:           60,
							Total:         660,
							AddressID:     "address-id",
							Address: &response.Address{
								Lastname:     "&.",
								Firstname:    "購入者",
								PostalCode:   "1000014",
								Prefecture:   "東京都",
								City:         "千代田区",
								AddressLine1: "永田町1-7-1",
								AddressLine2: "",
								PhoneNumber:  "+819012345678",
							},
						},
						orderID: "order-id",
					},
					fulfillment: &OrderFulfillment{
						OrderFulfillment: response.OrderFulfillment{
							TrackingNumber:  "",
							Status:          FulfillmentStatusFulfilled.Response(),
							ShippingCarrier: ShippingCarrierUnknown.Response(),
							ShippingMethod:  DeliveryTypeNormal.Response(),
							BoxSize:         ShippingSize60.Response(),
							AddressID:       "address-id",
							Address: &response.Address{
								Lastname:     "&.",
								Firstname:    "購入者",
								PostalCode:   "1000014",
								Prefecture:   "東京都",
								City:         "千代田区",
								AddressLine1: "永田町1-7-1",
								AddressLine2: "",
								PhoneNumber:  "+819012345678",
							},
						},
						orderID: "order-id",
					},
					refund: &OrderRefund{
						OrderRefund: response.OrderRefund{
							Canceled: false,
							Reason:   "",
							Total:    0,
						},
					},
					items: OrderItems{
						{
							OrderItem: response.OrderItem{
								ProductID: "product-id",
								Name:      "新鮮なじゃがいも",
								Price:     100,
								Quantity:  1,
								Weight:    1.3,
								Media: []*response.ProductMedia{
									{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
									{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
								},
							},
							orderID: "order-id",
						},
					},
					CoordinatorID: "coordinator-id",
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

func TestOrders_AddressIDs(t *testing.T) {
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
						UserName:    "&. 購入者",
						ScheduleID:  "schedule-id",
						PromotionID: "",
						OrderedAt:   0,
						PaidAt:      0,
						DeliveredAt: 0,
						CanceledAt:  0,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
					payment: &OrderPayment{
						OrderPayment: response.OrderPayment{
							TransactionID: "transaction-id",
							MethodType:    PaymentMethodTypeCreditCard.Response(),
							Status:        PaymentStatusPaid.Response(),
							Subtotal:      100,
							Discount:      0,
							ShippingFee:   500,
							Tax:           60,
							Total:         660,
							AddressID:     "address-id",
							Address: &response.Address{
								Lastname:     "&.",
								Firstname:    "購入者",
								PostalCode:   "1000014",
								Prefecture:   "東京都",
								City:         "千代田区",
								AddressLine1: "永田町1-7-1",
								AddressLine2: "",
								PhoneNumber:  "+819012345678",
							},
						},
						orderID: "order-id",
					},
					fulfillment: &OrderFulfillment{
						OrderFulfillment: response.OrderFulfillment{
							TrackingNumber:  "",
							Status:          FulfillmentStatusFulfilled.Response(),
							ShippingCarrier: ShippingCarrierUnknown.Response(),
							ShippingMethod:  DeliveryTypeNormal.Response(),
							BoxSize:         ShippingSize60.Response(),
							AddressID:       "address-id",
							Address: &response.Address{
								Lastname:     "&.",
								Firstname:    "購入者",
								PostalCode:   "1000014",
								Prefecture:   "東京都",
								City:         "千代田区",
								AddressLine1: "永田町1-7-1",
								AddressLine2: "",
								PhoneNumber:  "+819012345678",
							},
						},
						orderID: "order-id",
					},
					refund: &OrderRefund{
						OrderRefund: response.OrderRefund{
							Canceled: false,
							Reason:   "",
							Total:    0,
						},
					},
					items: OrderItems{
						{
							OrderItem: response.OrderItem{
								ProductID: "product-id",
								Name:      "新鮮なじゃがいも",
								Price:     100,
								Quantity:  1,
								Weight:    1.3,
								Media: []*response.ProductMedia{
									{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
									{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
								},
							},
							orderID: "order-id",
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			expect: []string{"address-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.orders.AddressIDs())
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
						UserName:    "&. 購入者",
						ScheduleID:  "schedule-id",
						PromotionID: "",
						OrderedAt:   0,
						PaidAt:      0,
						DeliveredAt: 0,
						CanceledAt:  0,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
					payment: &OrderPayment{
						OrderPayment: response.OrderPayment{
							TransactionID: "transaction-id",
							MethodType:    PaymentMethodTypeCreditCard.Response(),
							Status:        PaymentStatusPaid.Response(),
							Subtotal:      100,
							Discount:      0,
							ShippingFee:   500,
							Tax:           60,
							Total:         660,
							AddressID:     "address-id",
							Address: &response.Address{
								Lastname:     "&.",
								Firstname:    "購入者",
								PostalCode:   "1000014",
								Prefecture:   "東京都",
								City:         "千代田区",
								AddressLine1: "永田町1-7-1",
								AddressLine2: "",
								PhoneNumber:  "+819012345678",
							},
						},
						orderID: "order-id",
					},
					fulfillment: &OrderFulfillment{
						OrderFulfillment: response.OrderFulfillment{
							TrackingNumber:  "",
							Status:          FulfillmentStatusFulfilled.Response(),
							ShippingCarrier: ShippingCarrierUnknown.Response(),
							ShippingMethod:  DeliveryTypeNormal.Response(),
							BoxSize:         ShippingSize60.Response(),
							AddressID:       "address-id",
							Address: &response.Address{
								Lastname:     "&.",
								Firstname:    "購入者",
								PostalCode:   "1000014",
								Prefecture:   "東京都",
								City:         "千代田区",
								AddressLine1: "永田町1-7-1",
								AddressLine2: "",
								PhoneNumber:  "+819012345678",
							},
						},
						orderID: "order-id",
					},
					refund: &OrderRefund{
						OrderRefund: response.OrderRefund{
							Canceled: false,
							Reason:   "",
							Total:    0,
						},
					},
					items: OrderItems{
						{
							OrderItem: response.OrderItem{
								ProductID: "product-id",
								Name:      "新鮮なじゃがいも",
								Price:     100,
								Quantity:  1,
								Weight:    1.3,
								Media: []*response.ProductMedia{
									{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
									{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
								},
							},
							orderID: "order-id",
						},
					},
					CoordinatorID: "coordinator-id",
				},
			},
			expect: []*response.Order{
				{
					ID:          "order-id",
					UserID:      "user-id",
					UserName:    "&. 購入者",
					ScheduleID:  "schedule-id",
					PromotionID: "",
					Payment: &response.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    PaymentMethodTypeCreditCard.Response(),
						Status:        PaymentStatusPaid.Response(),
						Subtotal:      100,
						Discount:      0,
						ShippingFee:   500,
						Tax:           60,
						Total:         660,
						AddressID:     "address-id",
						Address: &response.Address{
							Lastname:     "&.",
							Firstname:    "購入者",
							PostalCode:   "1000014",
							Prefecture:   "東京都",
							City:         "千代田区",
							AddressLine1: "永田町1-7-1",
							AddressLine2: "",
							PhoneNumber:  "+819012345678",
						},
					},
					Fulfillment: &response.OrderFulfillment{
						TrackingNumber:  "",
						Status:          FulfillmentStatusFulfilled.Response(),
						ShippingCarrier: ShippingCarrierUnknown.Response(),
						ShippingMethod:  DeliveryTypeNormal.Response(),
						BoxSize:         ShippingSize60.Response(),
						AddressID:       "address-id",
						Address: &response.Address{
							Lastname:     "&.",
							Firstname:    "購入者",
							PostalCode:   "1000014",
							Prefecture:   "東京都",
							City:         "千代田区",
							AddressLine1: "永田町1-7-1",
							AddressLine2: "",
							PhoneNumber:  "+819012345678",
						},
					},
					Refund: &response.OrderRefund{
						Canceled: false,
						Reason:   "",
						Total:    0,
					},
					Items: []*response.OrderItem{
						{
							ProductID: "product-id",
							Name:      "新鮮なじゃがいも",
							Price:     100,
							Quantity:  1,
							Weight:    1.3,
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
