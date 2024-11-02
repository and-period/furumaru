package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestOrderStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		status   entity.OrderStatus
		expect   OrderStatus
		response int32
	}{
		{
			name:     "unpaid",
			status:   entity.OrderStatusUnpaid,
			expect:   OrderStatusUnpaid,
			response: 1,
		},
		{
			name:     "preparing",
			status:   entity.OrderStatusPreparing,
			expect:   OrderStatusPreparing,
			response: 2,
		},
		{
			name:     "completed",
			status:   entity.OrderStatusCompleted,
			expect:   OrderStatusCompleted,
			response: 3,
		},
		{
			name:     "canceled",
			status:   entity.OrderStatusCanceled,
			expect:   OrderStatusCanceled,
			response: 4,
		},
		{
			name:     "refunded",
			status:   entity.OrderStatusRefunded,
			expect:   OrderStatusRefunded,
			response: 5,
		},
		{
			name:     "failed",
			status:   entity.OrderStatusFailed,
			expect:   OrderStatusFailed,
			response: 6,
		},
		{
			name:     "unknown",
			status:   entity.OrderStatusUnknown,
			expect:   OrderStatusUnknown,
			response: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewOrderStatus(tt.status)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.response, actual.Response())
		})
	}
}

func TestOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		order       *entity.Order
		addresses   map[int64]*Address
		products    map[int64]*Product
		experiences map[int64]*Experience
		expect      *Order
	}{
		{
			name: "success",
			order: &entity.Order{
				ID:            "order-id",
				UserID:        "user-id",
				CoordinatorID: "coordinator-id",
				PromotionID:   "promotion-id",
				ManagementID:  1,
				Status:        entity.OrderStatusPreparing,
				OrderPayment: entity.OrderPayment{
					OrderID:           "order-id",
					AddressRevisionID: 1,
					TransactionID:     "transaction-id",
					Status:            entity.PaymentStatusCaptured,
					MethodType:        entity.PaymentMethodTypeCreditCard,
					Subtotal:          1980,
					Discount:          0,
					ShippingFee:       550,
					Tax:               230,
					Total:             2530,
					RefundTotal:       0,
					RefundType:        entity.RefundTypeNone,
					RefundReason:      "",
					OrderedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
					PaidAt:            jst.Date(2022, 1, 1, 0, 0, 0, 0),
					RefundedAt:        time.Time{},
					CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				OrderFulfillments: entity.OrderFulfillments{
					{
						ID:                "fulfillment-id",
						OrderID:           "order-id",
						AddressRevisionID: 1,
						TrackingNumber:    "",
						Status:            entity.FulfillmentStatusFulfilled,
						ShippingCarrier:   entity.ShippingCarrierUnknown,
						ShippingType:      entity.ShippingTypeNormal,
						BoxNumber:         1,
						BoxSize:           entity.ShippingSize60,
						CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
						UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
						ShippedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
					},
				},
				OrderItems: entity.OrderItems{
					{
						FulfillmentID:     "fulfillment-id",
						OrderID:           "order-id",
						ProductRevisionID: 1,
						Quantity:          1,
						CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
						UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
					},
				},
				OrderExperience: entity.OrderExperience{
					ExperienceRevisionID:  1,
					OrderID:               "order-id",
					AdultCount:            2,
					JuniorHighSchoolCount: 1,
					ElementarySchoolCount: 0,
					PreschoolCount:        0,
					SeniorCount:           0,
					Remarks: entity.OrderExperienceRemarks{
						Transportation: "電車",
						RequestedDate:  jst.Date(2024, 1, 2, 0, 0, 0, 0),
						RequestedTime:  jst.Date(0, 1, 1, 18, 30, 0, 0),
					},
				},
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			addresses: map[int64]*Address{
				1: {
					Address: response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					revisionID: 1,
				},
			},
			products: map[int64]*Product{
				1: {
					Product: response.Product{
						ID:              "product-id",
						CoordinatorID:   "coordinator-id",
						ProducerID:      "producer-id",
						CategoryID:      "promotion-id",
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
							},
							{
								URL:         "https://and-period.jp/thumbnail02.png",
								IsThumbnail: false,
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
			experiences: map[int64]*Experience{
				1: {
					Experience: response.Experience{
						ID:               "experience-id",
						CoordinatorID:    "coordinator-id",
						ProducerID:       "producer-id",
						ExperienceTypeID: "experience-type-id",
						Title:            "じゃがいも収穫",
						Description:      "じゃがいもを収穫する体験です。",
						Status:           int32(ExperienceStatusAccepting),
						Media: []*response.ExperienceMedia{
							{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
							{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
						},
						PriceAdult:            1000,
						PriceJuniorHighSchool: 800,
						PriceElementarySchool: 600,
						PricePreschool:        400,
						PriceSenior:           700,
						RecommendedPoint1:     "じゃがいもを収穫する楽しさを体験できます。",
						RecommendedPoint2:     "新鮮なじゃがいもを持ち帰ることができます。",
						RecommendedPoint3:     "じゃがいもの美味しさを再認識できます。",
						PromotionVideoURL:     "http://example.com/promotion.mp4",
						Duration:              60,
						Direction:             "彦根駅から徒歩10分",
						BusinessOpenTime:      "1000",
						BusinessCloseTime:     "1800",
						HostPostalCode:        "5220061",
						HostCity:              "彦根市",
						HostAddressLine1:      "金亀町１−１",
						HostAddressLine2:      "",
						StartAt:               1640962800,
						EndAt:                 1640962800,
					},
					revisionID: 1,
				},
			},
			expect: &Order{
				Order: response.Order{
					ID:            "order-id",
					CoordinatorID: "coordinator-id",
					PromotionID:   "promotion-id",
					Status:        int32(OrderStatusPreparing),
					Payment: &response.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    PaymentMethodTypeCreditCard.Response(),
						Status:        PaymentStatusPaid.Response(),
						Subtotal:      1980,
						Discount:      0,
						ShippingFee:   550,
						Total:         2530,
						OrderedAt:     1640962800,
						PaidAt:        1640962800,
					},
					Fulfillments: []*response.OrderFulfillment{
						{
							FulfillmentID:   "fulfillment-id",
							TrackingNumber:  "",
							Status:          FulfillmentStatusFulfilled.Response(),
							ShippingCarrier: ShippingCarrierUnknown.Response(),
							ShippingType:    ShippingTypeNormal.Response(),
							BoxNumber:       1,
							BoxSize:         ShippingSize60.Response(),
							ShippedAt:       1640962800,
						},
					},
					Refund: &response.OrderRefund{
						Total:      0,
						Type:       RefundTypeNone.Response(),
						Reason:     "",
						Canceled:   false,
						CanceledAt: 0,
					},
					Items: []*response.OrderItem{
						{
							FulfillmentID: "fulfillment-id",
							ProductID:     "product-id",
							Price:         400,
							Quantity:      1,
						},
					},
					Experience: &response.OrderExperience{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						AdultPrice:            1000,
						JuniorHighSchoolCount: 1,
						JuniorHighSchoolPrice: 800,
						ElementarySchoolCount: 0,
						ElementarySchoolPrice: 600,
						PreschoolCount:        0,
						PreschoolPrice:        400,
						SeniorCount:           0,
						SeniorPrice:           700,
						Transportation:        "電車",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					BillingAddress: &response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					ShippingAddress: &response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, NewOrder(tt.order, tt.addresses, tt.products, tt.experiences))
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
					ID:            "order-id",
					CoordinatorID: "coordinator-id",
					PromotionID:   "",
					Status:        int32(OrderStatusPreparing),
					Payment: &response.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    PaymentMethodTypeCreditCard.Response(),
						Status:        PaymentStatusPaid.Response(),
						Subtotal:      1100,
						Discount:      0,
						ShippingFee:   500,
						Total:         1600,
						OrderedAt:     1640962800,
						PaidAt:        1640962800,
					},
					Fulfillments: []*response.OrderFulfillment{
						{
							FulfillmentID:   "fulfillment-id",
							TrackingNumber:  "",
							Status:          FulfillmentStatusFulfilled.Response(),
							ShippingCarrier: ShippingCarrierUnknown.Response(),
							ShippingType:    ShippingTypeNormal.Response(),
							BoxNumber:       1,
							BoxSize:         ShippingSize60.Response(),
							ShippedAt:       1640962800,
						},
					},
					Refund: &response.OrderRefund{
						Total:      0,
						Type:       RefundTypeNone.Response(),
						Reason:     "",
						Canceled:   false,
						CanceledAt: 0,
					},
					Items: []*response.OrderItem{
						{
							FulfillmentID: "fulfillment-id",
							ProductID:     "product-id",
							Price:         400,
							Quantity:      1,
						},
					},
					BillingAddress: &response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					ShippingAddress: &response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
				},
			},
			expect: []string{"product-id"},
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
					ID:            "order-id",
					CoordinatorID: "coordinator-id",
					PromotionID:   "promotion-id",
					Status:        int32(OrderStatusPreparing),
					Payment: &response.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    PaymentMethodTypeCreditCard.Response(),
						Status:        PaymentStatusPaid.Response(),
						Subtotal:      1980,
						Discount:      0,
						ShippingFee:   550,
						Total:         2530,
						OrderedAt:     1640962800,
						PaidAt:        1640962800,
					},
					Fulfillments: []*response.OrderFulfillment{
						{
							FulfillmentID:   "fulfillment-id",
							TrackingNumber:  "",
							Status:          FulfillmentStatusFulfilled.Response(),
							ShippingCarrier: ShippingCarrierUnknown.Response(),
							ShippingType:    ShippingTypeNormal.Response(),
							BoxNumber:       1,
							BoxSize:         ShippingSize60.Response(),
							ShippedAt:       1640962800,
						},
					},
					Refund: &response.OrderRefund{
						Total:      0,
						Type:       RefundTypeNone.Response(),
						Reason:     "",
						Canceled:   false,
						CanceledAt: 0,
					},
					Items: []*response.OrderItem{
						{
							FulfillmentID: "fulfillment-id",
							ProductID:     "product-id",
							Price:         400,
							Quantity:      1,
						},
					},
					BillingAddress: &response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					ShippingAddress: &response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
				},
			},
			expect: &response.Order{
				ID:            "order-id",
				CoordinatorID: "coordinator-id",
				PromotionID:   "promotion-id",
				Status:        int32(OrderStatusPreparing),
				Payment: &response.OrderPayment{
					TransactionID: "transaction-id",
					MethodType:    PaymentMethodTypeCreditCard.Response(),
					Status:        PaymentStatusPaid.Response(),
					Subtotal:      1980,
					Discount:      0,
					ShippingFee:   550,
					Total:         2530,
					OrderedAt:     1640962800,
					PaidAt:        1640962800,
				},
				Fulfillments: []*response.OrderFulfillment{
					{
						FulfillmentID:   "fulfillment-id",
						TrackingNumber:  "",
						Status:          FulfillmentStatusFulfilled.Response(),
						ShippingCarrier: ShippingCarrierUnknown.Response(),
						ShippingType:    ShippingTypeNormal.Response(),
						BoxNumber:       1,
						BoxSize:         ShippingSize60.Response(),
						ShippedAt:       1640962800,
					},
				},
				Refund: &response.OrderRefund{
					Total:      0,
					Type:       RefundTypeNone.Response(),
					Reason:     "",
					Canceled:   false,
					CanceledAt: 0,
				},
				Items: []*response.OrderItem{
					{
						FulfillmentID: "fulfillment-id",
						ProductID:     "product-id",
						Price:         400,
						Quantity:      1,
					},
				},
				BillingAddress: &response.Address{
					Lastname:       "&.",
					Firstname:      "購入者",
					PostalCode:     "1000014",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "090-1234-1234",
				},
				ShippingAddress: &response.Address{
					Lastname:       "&.",
					Firstname:      "購入者",
					PostalCode:     "1000014",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "090-1234-1234",
				},
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
