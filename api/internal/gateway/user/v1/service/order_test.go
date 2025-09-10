package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestNewOrderType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		typ      entity.OrderType
		expect   OrderType
		response int32
	}{
		{
			name:     "product",
			typ:      entity.OrderType(types.OrderTypeProduct),
			expect:   OrderType(types.OrderTypeProduct),
			response: 1,
		},
		{
			name:     "experience",
			typ:      entity.OrderType(types.OrderTypeExperience),
			expect:   OrderType(types.OrderTypeExperience),
			response: 2,
		},
		{
			name:     "unknown",
			typ:      entity.OrderType(types.OrderTypeUnknown),
			expect:   OrderType(types.OrderTypeUnknown),
			response: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewOrderType(tt.typ)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.response, actual.Response())
		})
	}
}

func TestNewOrderTypeFromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		typ      string
		expect   OrderType
		response int32
	}{
		{
			name:     "product",
			typ:      "product",
			expect:   OrderType(types.OrderTypeProduct),
			response: 1,
		},
		{
			name:     "experience",
			typ:      "experience",
			expect:   OrderType(types.OrderTypeExperience),
			response: 2,
		},
		{
			name:     "unknown",
			typ:      "unknown",
			expect:   OrderType(types.OrderTypeUnknown),
			response: 0,
		},
		{
			name:     "empty",
			typ:      "",
			expect:   OrderType(types.OrderTypeUnknown),
			response: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewOrderTypeFromString(tt.typ)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.response, actual.Response())
		})
	}
}

func TestOrderType_StoreEntity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		typ    OrderType
		expect entity.OrderType
	}{
		{
			name:   "product",
			typ:    OrderType(types.OrderTypeProduct),
			expect: entity.OrderType(types.OrderTypeProduct),
		},
		{
			name:   "experience",
			typ:    OrderType(types.OrderTypeExperience),
			expect: entity.OrderType(types.OrderTypeExperience),
		},
		{
			name:   "unknown",
			typ:    OrderType(types.OrderTypeUnknown),
			expect: entity.OrderType(types.OrderTypeUnknown),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.typ.StoreEntity())
		})
	}
}

func TestNewOrderStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		status   entity.OrderStatus
		expect   OrderStatus
		response int32
	}{
		{
			name:     "unpaid",
			status:   entity.OrderStatus(types.OrderStatusUnpaid),
			expect:   OrderStatus(types.OrderStatusUnpaid),
			response: 1,
		},
		{
			name:     "waiting",
			status:   entity.OrderStatusWaiting,
			expect:   OrderStatus(types.OrderStatusPreparing),
			response: 2,
		},
		{
			name:     "preparing",
			status:   entity.OrderStatus(types.OrderStatusPreparing),
			expect:   OrderStatus(types.OrderStatusPreparing),
			response: 2,
		},
		{
			name:     "shipped",
			status:   entity.OrderStatusShipped,
			expect:   OrderStatus(types.OrderStatusPreparing),
			response: 2,
		},
		{
			name:     "completed",
			status:   entity.OrderStatus(types.OrderStatusCompleted),
			expect:   OrderStatus(types.OrderStatusCompleted),
			response: 3,
		},
		{
			name:     "canceled",
			status:   entity.OrderStatus(types.OrderStatusCanceled),
			expect:   OrderStatus(types.OrderStatusCanceled),
			response: 4,
		},
		{
			name:     "refunded",
			status:   entity.OrderStatus(types.OrderStatusRefunded),
			expect:   OrderStatus(types.OrderStatusRefunded),
			response: 5,
		},
		{
			name:     "failed",
			status:   entity.OrderStatus(types.OrderStatusFailed),
			expect:   OrderStatus(types.OrderStatusFailed),
			response: 6,
		},
		{
			name:     "unknown",
			status:   entity.OrderStatus(types.OrderStatusUnknown),
			expect:   OrderStatus(types.OrderStatusUnknown),
			response: 0,
		},
	}
	for _, tt := range tests {
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
				Type:          entity.OrderType(types.OrderTypeProduct),
				Status:        entity.OrderStatus(types.OrderStatusPreparing),
				OrderPayment: entity.OrderPayment{
					OrderID:           "order-id",
					AddressRevisionID: 1,
					TransactionID:     "transaction-id",
					Status:            entity.PaymentStatusCaptured,
					MethodType:        entity.PaymentMethodType(types.PaymentMethodTypeCreditCard),
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
						ShippingCarrier:   entity.ShippingCarrier(types.ShippingCarrierUnknown),
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
					Address: types.Address{
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
					Product: types.Product{
						ID:              "product-id",
						CoordinatorID:   "coordinator-id",
						ProducerID:      "producer-id",
						CategoryID:      "promotion-id",
						ProductTypeID:   "product-type-id",
						ProductTagIDs:   []string{"product-tag-id"},
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Status:          types.ProductStatusForSale,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*types.ProductMedia{
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
						StorageMethodType: types.StorageMethodTypeNormal,
						DeliveryType:      types.DeliveryTypeNormal,
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
					Experience: types.Experience{
						ID:               "experience-id",
						CoordinatorID:    "coordinator-id",
						ProducerID:       "producer-id",
						ExperienceTypeID: "experience-type-id",
						Title:            "じゃがいも収穫",
						Description:      "じゃがいもを収穫する体験です。",
						Status:           types.ExperienceStatusAccepting,
						Media: []*types.ExperienceMedia{
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
				Order: types.Order{
					ID:            "order-id",
					CoordinatorID: "coordinator-id",
					PromotionID:   "promotion-id",
					Type:          types.OrderTypeProduct,
					Status:        types.OrderStatusPreparing,
					Payment: &types.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    PaymentMethodType(types.PaymentMethodTypeCreditCard).Response(),
						Status:        PaymentStatus(types.PaymentStatusPaid).Response(),
						Subtotal:      1980,
						Discount:      0,
						ShippingFee:   550,
						Total:         2530,
						OrderedAt:     1640962800,
						PaidAt:        1640962800,
					},
					Fulfillments: []*types.OrderFulfillment{
						{
							FulfillmentID:   "fulfillment-id",
							TrackingNumber:  "",
							Status:          FulfillmentStatus(types.FulfillmentStatusFulfilled).Response(),
							ShippingCarrier: ShippingCarrier(types.ShippingCarrierUnknown).Response(),
							ShippingType:    ShippingType(types.ShippingTypeNormal).Response(),
							BoxNumber:       1,
							BoxSize:         ShippingSize(types.ShippingSize60).Response(),
							ShippedAt:       1640962800,
						},
					},
					Refund: &types.OrderRefund{
						Total:      0,
						Type:       RefundType(types.RefundTypeNone).Response(),
						Reason:     "",
						Canceled:   false,
						CanceledAt: 0,
					},
					Items: []*types.OrderItem{
						{
							FulfillmentID: "fulfillment-id",
							ProductID:     "product-id",
							Price:         400,
							Quantity:      1,
						},
					},
					Experience: &types.OrderExperience{
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
					BillingAddress: &types.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					ShippingAddress: &types.Address{
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
				Order: types.Order{
					ID:            "order-id",
					CoordinatorID: "coordinator-id",
					PromotionID:   "",
					Status:        types.OrderStatusPreparing,
					Payment: &types.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    PaymentMethodType(types.PaymentMethodTypeCreditCard).Response(),
						Status:        PaymentStatus(types.PaymentStatusPaid).Response(),
						Subtotal:      1100,
						Discount:      0,
						ShippingFee:   500,
						Total:         1600,
						OrderedAt:     1640962800,
						PaidAt:        1640962800,
					},
					Fulfillments: []*types.OrderFulfillment{
						{
							FulfillmentID:   "fulfillment-id",
							TrackingNumber:  "",
							Status:          FulfillmentStatus(types.FulfillmentStatusFulfilled).Response(),
							ShippingCarrier: ShippingCarrier(types.ShippingCarrierUnknown).Response(),
							ShippingType:    ShippingType(types.ShippingTypeNormal).Response(),
							BoxNumber:       1,
							BoxSize:         ShippingSize(types.ShippingSize60).Response(),
							ShippedAt:       1640962800,
						},
					},
					Refund: &types.OrderRefund{
						Total:      0,
						Type:       RefundType(types.RefundTypeNone).Response(),
						Reason:     "",
						Canceled:   false,
						CanceledAt: 0,
					},
					Items: []*types.OrderItem{
						{
							FulfillmentID: "fulfillment-id",
							ProductID:     "product-id",
							Price:         400,
							Quantity:      1,
						},
					},
					BillingAddress: &types.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					ShippingAddress: &types.Address{
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
		expect *types.Order
	}{
		{
			name: "success",
			order: &Order{
				Order: types.Order{
					ID:            "order-id",
					CoordinatorID: "coordinator-id",
					PromotionID:   "promotion-id",
					Status:        types.OrderStatusPreparing,
					Payment: &types.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    PaymentMethodType(types.PaymentMethodTypeCreditCard).Response(),
						Status:        PaymentStatus(types.PaymentStatusPaid).Response(),
						Subtotal:      1980,
						Discount:      0,
						ShippingFee:   550,
						Total:         2530,
						OrderedAt:     1640962800,
						PaidAt:        1640962800,
					},
					Fulfillments: []*types.OrderFulfillment{
						{
							FulfillmentID:   "fulfillment-id",
							TrackingNumber:  "",
							Status:          FulfillmentStatus(types.FulfillmentStatusFulfilled).Response(),
							ShippingCarrier: ShippingCarrier(types.ShippingCarrierUnknown).Response(),
							ShippingType:    ShippingType(types.ShippingTypeNormal).Response(),
							BoxNumber:       1,
							BoxSize:         ShippingSize(types.ShippingSize60).Response(),
							ShippedAt:       1640962800,
						},
					},
					Refund: &types.OrderRefund{
						Total:      0,
						Type:       RefundType(types.RefundTypeNone).Response(),
						Reason:     "",
						Canceled:   false,
						CanceledAt: 0,
					},
					Items: []*types.OrderItem{
						{
							FulfillmentID: "fulfillment-id",
							ProductID:     "product-id",
							Price:         400,
							Quantity:      1,
						},
					},
					BillingAddress: &types.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					ShippingAddress: &types.Address{
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
			expect: &types.Order{
				ID:            "order-id",
				CoordinatorID: "coordinator-id",
				PromotionID:   "promotion-id",
				Status:        types.OrderStatusPreparing,
				Payment: &types.OrderPayment{
					TransactionID: "transaction-id",
					MethodType:    PaymentMethodType(types.PaymentMethodTypeCreditCard).Response(),
					Status:        PaymentStatus(types.PaymentStatusPaid).Response(),
					Subtotal:      1980,
					Discount:      0,
					ShippingFee:   550,
					Total:         2530,
					OrderedAt:     1640962800,
					PaidAt:        1640962800,
				},
				Fulfillments: []*types.OrderFulfillment{
					{
						FulfillmentID:   "fulfillment-id",
						TrackingNumber:  "",
						Status:          FulfillmentStatus(types.FulfillmentStatusFulfilled).Response(),
						ShippingCarrier: ShippingCarrier(types.ShippingCarrierUnknown).Response(),
						ShippingType:    ShippingType(types.ShippingTypeNormal).Response(),
						BoxNumber:       1,
						BoxSize:         ShippingSize(types.ShippingSize60).Response(),
						ShippedAt:       1640962800,
					},
				},
				Refund: &types.OrderRefund{
					Total:      0,
					Type:       RefundType(types.RefundTypeNone).Response(),
					Reason:     "",
					Canceled:   false,
					CanceledAt: 0,
				},
				Items: []*types.OrderItem{
					{
						FulfillmentID: "fulfillment-id",
						ProductID:     "product-id",
						Price:         400,
						Quantity:      1,
					},
				},
				BillingAddress: &types.Address{
					Lastname:       "&.",
					Firstname:      "購入者",
					PostalCode:     "1000014",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "090-1234-1234",
				},
				ShippingAddress: &types.Address{
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
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.order.Response())
		})
	}
}
