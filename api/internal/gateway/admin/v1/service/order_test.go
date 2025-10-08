package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestNewOrderType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		typ    entity.OrderType
		expect types.OrderType
	}{
		{
			name:   "product",
			typ:    entity.OrderTypeProduct,
			expect: types.OrderTypeProduct,
		},
		{
			name:   "experience",
			typ:    entity.OrderTypeExperience,
			expect: types.OrderTypeExperience,
		},
		{
			name:   "unknown",
			typ:    entity.OrderTypeUnknown,
			expect: types.OrderTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewOrderType(tt.typ)
			assert.Equal(t, tt.expect, actual.Response())
		})
	}
}

func TestNewOrderTypeFromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		typ    string
		expect types.OrderType
	}{
		{
			name:   "product",
			typ:    "product",
			expect: types.OrderTypeProduct,
		},
		{
			name:   "experience",
			typ:    "experience",
			expect: types.OrderTypeExperience,
		},
		{
			name:   "unknown",
			typ:    "unknown",
			expect: types.OrderTypeUnknown,
		},
		{
			name:   "empty",
			typ:    "",
			expect: types.OrderTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewOrderTypeFromString(tt.typ)
			assert.Equal(t, tt.expect, actual.Response())
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
			expect: entity.OrderTypeProduct,
		},
		{
			name:   "experience",
			typ:    OrderType(types.OrderTypeExperience),
			expect: entity.OrderTypeExperience,
		},
		{
			name:   "unknown",
			typ:    OrderType(types.OrderTypeUnknown),
			expect: entity.OrderTypeUnknown,
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
		name   string
		status entity.OrderStatus
		expect types.OrderStatus
	}{
		{
			name:   "unpaid",
			status: entity.OrderStatusUnpaid,
			expect: types.OrderStatusUnpaid,
		},
		{
			name:   "waiting",
			status: entity.OrderStatusWaiting,
			expect: types.OrderStatusWaiting,
		},
		{
			name:   "preparing",
			status: entity.OrderStatusPreparing,
			expect: types.OrderStatusPreparing,
		},
		{
			name:   "shipped",
			status: entity.OrderStatusShipped,
			expect: types.OrderStatusShipped,
		},
		{
			name:   "completed",
			status: entity.OrderStatusCompleted,
			expect: types.OrderStatusCompleted,
		},
		{
			name:   "canceled",
			status: entity.OrderStatusCanceled,
			expect: types.OrderStatusCanceled,
		},
		{
			name:   "refunded",
			status: entity.OrderStatusRefunded,
			expect: types.OrderStatusRefunded,
		},
		{
			name:   "failed",
			status: entity.OrderStatusFailed,
			expect: types.OrderStatusFailed,
		},
		{
			name:   "unknown",
			status: entity.OrderStatusUnknown,
			expect: types.OrderStatusUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewOrderStatus(tt.status)
			assert.Equal(t, tt.expect, actual.Response())
		})
	}
}

func TestNewOrderShippingType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		typ    entity.OrderShippingType
		expect types.OrderShippingType
	}{
		{
			name:   "none",
			typ:    entity.OrderShippingTypeNone,
			expect: types.OrderShippingTypeNone,
		},
		{
			name:   "standard",
			typ:    entity.OrderShippingTypeStandard,
			expect: types.OrderShippingTypeStandard,
		},
		{
			name:   "pickup",
			typ:    entity.OrderShippingTypePickup,
			expect: types.OrderShippingTypePickup,
		},
		{
			name:   "unknown",
			typ:    entity.OrderShippingTypeUnknown,
			expect: types.OrderShippingTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewOrderShippingType(tt.typ)
			assert.Equal(t, tt.expect, actual.Response())
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
				Type:          entity.OrderTypeProduct,
				Status:        entity.OrderStatusShipped,
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
					JuniorHighSchoolCount: 0,
					ElementarySchoolCount: 0,
					PreschoolCount:        0,
					SeniorCount:           0,
					Remarks: entity.OrderExperienceRemarks{
						Transportation: "電車",
						RequestedDate:  jst.Date(2024, 1, 2, 0, 0, 0, 0),
						RequestedTime:  jst.Date(0, 1, 1, 18, 30, 0, 0),
					},
					CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
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
						PhoneNumber:    "090-1234-5678",
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
						Scope:           types.ProductScopePublic,
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
						Price:                400,
						Cost:                 300,
						RecommendedPoint1:    "ポイント1",
						RecommendedPoint2:    "ポイント2",
						RecommendedPoint3:    "ポイント3",
						StorageMethodType:    types.StorageMethodTypeNormal,
						DeliveryType:         types.DeliveryTypeNormal,
						Box60Rate:            50,
						Box80Rate:            40,
						Box100Rate:           30,
						OriginPrefectureCode: 25,
						OriginCity:           "彦根市",
						StartAt:              1640962800,
						EndAt:                1640962800,
						CreatedAt:            1640962800,
						UpdatedAt:            1640962800,
					},
					revisionID: 1,
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
						Public:           true,
						SoldOut:          false,
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
						HostPrefectureCode:    25,
						HostCity:              "彦根市",
						HostAddressLine1:      "金亀町１−１",
						HostAddressLine2:      "",
						StartAt:               1640962800,
						EndAt:                 1640962800,
						CreatedAt:             1640962800,
						UpdatedAt:             1640962800,
					},
					revisionID: 1,
				},
			},
			expect: &Order{
				Order: types.Order{
					ID:              "order-id",
					UserID:          "user-id",
					CoordinatorID:   "coordinator-id",
					PromotionID:     "promotion-id",
					ManagementID:    1,
					ShippingMessage: "",
					Type:            types.OrderTypeProduct,
					Status:          types.OrderStatusShipped,
					CreatedAt:       1640962800,
					UpdatedAt:       1640962800,
					CompletedAt:     0,
					Metadata: &types.OrderMetadata{
						PickupAt:       0,
						PickupLocation: "",
					},
					Payment: &types.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    types.PaymentMethodTypeCreditCard,
						Status:        types.PaymentStatusPaid,
						Subtotal:      1980,
						Discount:      0,
						ShippingFee:   550,
						Total:         2530,
						OrderedAt:     1640962800,
						PaidAt:        1640962800,
						Address: &types.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							PostalCode:     "1000014",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "090-1234-5678",
						},
					},
					Fulfillments: []*types.OrderFulfillment{
						{
							FulfillmentID:   "fulfillment-id",
							TrackingNumber:  "",
							Status:          types.FulfillmentStatusFulfilled,
							ShippingCarrier: types.ShippingCarrierUnknown,
							ShippingType:    types.ShippingTypeNormal,
							BoxNumber:       1,
							BoxSize:         types.ShippingSize60,
							ShippedAt:       1640962800,
							Address: &types.Address{
								Lastname:       "&.",
								Firstname:      "購入者",
								PostalCode:     "1000014",
								PrefectureCode: 13,
								City:           "千代田区",
								AddressLine1:   "永田町1-7-1",
								AddressLine2:   "",
								PhoneNumber:    "090-1234-5678",
							},
						},
					},
					Refund: &types.OrderRefund{
						Total:      0,
						Type:       types.RefundTypeNone,
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
						JuniorHighSchoolCount: 0,
						JuniorHighSchoolPrice: 800,
						ElementarySchoolCount: 0,
						ElementarySchoolPrice: 600,
						PreschoolCount:        0,
						PreschoolPrice:        400,
						SeniorCount:           0,
						SeniorPrice:           700,
						Remarks: &types.OrderExperienceRemarks{
							Transportation: "電車",
							RequestedDate:  "20240102",
							RequestedTime:  "1830",
						},
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
					ID:              "order-id",
					UserID:          "user-id",
					CoordinatorID:   "coordinator-id",
					PromotionID:     "",
					ShippingMessage: "",
					Status:          types.OrderStatusShipped,
					CreatedAt:       1640962800,
					UpdatedAt:       1640962800,
					Payment: &types.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    types.PaymentMethodTypeCreditCard,
						Status:        types.PaymentStatusPaid,
						Subtotal:      1100,
						Discount:      0,
						ShippingFee:   500,
						Total:         1760,
						OrderedAt:     1640962800,
						PaidAt:        1640962800,
						Address: &types.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							PostalCode:     "1000014",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "090-1234-5678",
						},
					},
					Fulfillments: []*types.OrderFulfillment{
						{
							FulfillmentID:   "fulfillment-id",
							TrackingNumber:  "",
							Status:          types.FulfillmentStatusFulfilled,
							ShippingCarrier: types.ShippingCarrierUnknown,
							ShippingType:    types.ShippingTypeNormal,
							BoxNumber:       1,
							BoxSize:         types.ShippingSize60,
							ShippedAt:       1640962800,
							Address: &types.Address{
								Lastname:       "&.",
								Firstname:      "購入者",
								PostalCode:     "1000014",
								PrefectureCode: 13,
								City:           "千代田区",
								AddressLine1:   "永田町1-7-1",
								AddressLine2:   "",
								PhoneNumber:    "090-1234-5678",
							},
						},
					},
					Refund: &types.OrderRefund{
						Total:      0,
						Type:       types.RefundTypeNone,
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
					ID:              "order-id",
					UserID:          "user-id",
					CoordinatorID:   "coordinator-id",
					PromotionID:     "",
					ShippingMessage: "",
					Status:          types.OrderStatusShipped,
					CreatedAt:       1640962800,
					UpdatedAt:       1640962800,
					Payment: &types.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    types.PaymentMethodTypeCreditCard,
						Status:        types.PaymentStatusPaid,
						Subtotal:      1100,
						Discount:      0,
						ShippingFee:   500,
						Total:         1760,
						OrderedAt:     1640962800,
						PaidAt:        1640962800,
						Address: &types.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							PostalCode:     "1000014",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "090-1234-5678",
						},
					},
					Fulfillments: []*types.OrderFulfillment{
						{
							FulfillmentID:   "fulfillment-id",
							TrackingNumber:  "",
							Status:          types.FulfillmentStatusFulfilled,
							ShippingCarrier: types.ShippingCarrierUnknown,
							ShippingType:    types.ShippingTypeNormal,
							BoxNumber:       1,
							BoxSize:         types.ShippingSize60,
							ShippedAt:       1640962800,
							Address: &types.Address{
								Lastname:       "&.",
								Firstname:      "購入者",
								PostalCode:     "1000014",
								PrefectureCode: 13,
								City:           "千代田区",
								AddressLine1:   "永田町1-7-1",
								AddressLine2:   "",
								PhoneNumber:    "090-1234-5678",
							},
						},
					},
					Refund: &types.OrderRefund{
						Total:      0,
						Type:       types.RefundTypeNone,
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
				},
			},
			expect: &types.Order{
				ID:              "order-id",
				UserID:          "user-id",
				CoordinatorID:   "coordinator-id",
				PromotionID:     "",
				ShippingMessage: "",
				Status:          types.OrderStatusShipped,
				CreatedAt:       1640962800,
				UpdatedAt:       1640962800,
				Payment: &types.OrderPayment{
					TransactionID: "transaction-id",
					MethodType:    types.PaymentMethodTypeCreditCard,
					Status:        types.PaymentStatusPaid,
					Subtotal:      1100,
					Discount:      0,
					ShippingFee:   500,
					Total:         1760,
					OrderedAt:     1640962800,
					PaidAt:        1640962800,
					Address: &types.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-5678",
					},
				},
				Fulfillments: []*types.OrderFulfillment{
					{
						FulfillmentID:   "fulfillment-id",
						TrackingNumber:  "",
						Status:          types.FulfillmentStatusFulfilled,
						ShippingCarrier: types.ShippingCarrierUnknown,
						ShippingType:    types.ShippingTypeNormal,
						BoxNumber:       1,
						BoxSize:         types.ShippingSize60,
						ShippedAt:       1640962800,
						Address: &types.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							PostalCode:     "1000014",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "090-1234-5678",
						},
					},
				},
				Refund: &types.OrderRefund{
					Total:      0,
					Type:       types.RefundTypeNone,
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.order.Response())
		})
	}
}

func TestOrders(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		orders      entity.Orders
		addresses   map[int64]*Address
		products    map[int64]*Product
		experiences map[int64]*Experience
		expect      Orders
	}{
		{
			name: "success",
			orders: entity.Orders{
				{
					ID:            "product-order-id",
					Type:          entity.OrderTypeProduct,
					UserID:        "user-id",
					CoordinatorID: "coordinator-id",
					PromotionID:   "promotion-id",
					Status:        entity.OrderStatusShipped,
					OrderPayment: entity.OrderPayment{
						OrderID:           "product-order-id",
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
							OrderID:           "product-order-id",
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
							OrderID:           "product-order-id",
							ProductRevisionID: 1,
							Quantity:          1,
							CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
							UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
						},
					},
					CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				{
					ID:            "experience-order-id",
					Type:          entity.OrderTypeExperience,
					UserID:        "user-id",
					CoordinatorID: "coordinator-id",
					PromotionID:   "promotion-id",
					Status:        entity.OrderStatusShipped,
					OrderPayment: entity.OrderPayment{
						OrderID:           "experience-order-id",
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
					OrderExperience: entity.OrderExperience{
						OrderID:               "order-experience-id",
						ExperienceRevisionID:  1,
						AdultCount:            2,
						JuniorHighSchoolCount: 2,
						ElementarySchoolCount: 0,
						PreschoolCount:        0,
						SeniorCount:           0,
						Remarks: entity.OrderExperienceRemarks{
							Transportation: "電車",
							RequestedDate:  jst.Date(2024, 1, 2, 0, 0, 0, 0),
							RequestedTime:  jst.Date(0, 1, 1, 18, 30, 0, 0),
						},
						CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
						UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					},
					CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
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
						PhoneNumber:    "090-1234-5678",
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
						CategoryID:      "",
						ProductTypeID:   "product-type-id",
						ProductTagIDs:   []string{"product-tag-id"},
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Scope:           types.ProductScopePublic,
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
						Price:                400,
						Cost:                 300,
						RecommendedPoint1:    "ポイント1",
						RecommendedPoint2:    "ポイント2",
						RecommendedPoint3:    "ポイント3",
						StorageMethodType:    types.StorageMethodTypeNormal,
						DeliveryType:         types.DeliveryTypeNormal,
						Box60Rate:            50,
						Box80Rate:            40,
						Box100Rate:           30,
						OriginPrefectureCode: 25,
						OriginCity:           "彦根市",
						StartAt:              1640962800,
						EndAt:                1640962800,
						CreatedAt:            1640962800,
						UpdatedAt:            1640962800,
					},
					revisionID: 1,
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
						Public:           true,
						SoldOut:          false,
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
						HostPrefectureCode:    25,
						HostCity:              "彦根市",
						HostAddressLine1:      "金亀町１−１",
						HostAddressLine2:      "",
						StartAt:               1640962800,
						EndAt:                 1640962800,
						CreatedAt:             1640962800,
						UpdatedAt:             1640962800,
					},
					revisionID: 1,
				},
			},
			expect: Orders{
				{
					Order: types.Order{
						ID:              "product-order-id",
						UserID:          "user-id",
						CoordinatorID:   "coordinator-id",
						PromotionID:     "promotion-id",
						ShippingMessage: "",
						Type:            types.OrderTypeProduct,
						Status:          types.OrderStatusShipped,
						CreatedAt:       1640962800,
						UpdatedAt:       1640962800,
						CompletedAt:     0,
						Metadata: &types.OrderMetadata{
							PickupAt:       0,
							PickupLocation: "",
						},
						Payment: &types.OrderPayment{
							TransactionID: "transaction-id",
							MethodType:    types.PaymentMethodTypeCreditCard,
							Status:        types.PaymentStatusPaid,
							Subtotal:      1980,
							Discount:      0,
							ShippingFee:   550,
							Total:         2530,
							OrderedAt:     1640962800,
							PaidAt:        1640962800,
							Address: &types.Address{
								Lastname:       "&.",
								Firstname:      "購入者",
								PostalCode:     "1000014",
								PrefectureCode: 13,
								City:           "千代田区",
								AddressLine1:   "永田町1-7-1",
								AddressLine2:   "",
								PhoneNumber:    "090-1234-5678",
							},
						},
						Fulfillments: []*types.OrderFulfillment{
							{
								FulfillmentID:   "fulfillment-id",
								TrackingNumber:  "",
								Status:          types.FulfillmentStatusFulfilled,
								ShippingCarrier: types.ShippingCarrierUnknown,
								ShippingType:    types.ShippingTypeNormal,
								BoxNumber:       1,
								BoxSize:         types.ShippingSize60,
								ShippedAt:       1640962800,
								Address: &types.Address{
									Lastname:       "&.",
									Firstname:      "購入者",
									PostalCode:     "1000014",
									PrefectureCode: 13,
									City:           "千代田区",
									AddressLine1:   "永田町1-7-1",
									AddressLine2:   "",
									PhoneNumber:    "090-1234-5678",
								},
							},
						},
						Refund: &types.OrderRefund{
							Total:      0,
							Type:       types.RefundTypeNone,
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
							Remarks: &types.OrderExperienceRemarks{},
						},
					},
				},
				{
					Order: types.Order{
						ID:              "experience-order-id",
						UserID:          "user-id",
						CoordinatorID:   "coordinator-id",
						PromotionID:     "promotion-id",
						ShippingMessage: "",
						Type:            types.OrderTypeExperience,
						Status:          types.OrderStatusShipped,
						CreatedAt:       1640962800,
						UpdatedAt:       1640962800,
						CompletedAt:     0,
						Metadata: &types.OrderMetadata{
							PickupAt:       0,
							PickupLocation: "",
						},
						Payment: &types.OrderPayment{
							TransactionID: "transaction-id",
							MethodType:    types.PaymentMethodTypeCreditCard,
							Status:        types.PaymentStatusPaid,
							Subtotal:      1980,
							Discount:      0,
							ShippingFee:   550,
							Total:         2530,
							OrderedAt:     1640962800,
							PaidAt:        1640962800,
							Address: &types.Address{
								Lastname:       "&.",
								Firstname:      "購入者",
								PostalCode:     "1000014",
								PrefectureCode: 13,
								City:           "千代田区",
								AddressLine1:   "永田町1-7-1",
								AddressLine2:   "",
								PhoneNumber:    "090-1234-5678",
							},
						},
						Refund:       &types.OrderRefund{},
						Fulfillments: []*types.OrderFulfillment{},
						Items:        []*types.OrderItem{},
						Experience: &types.OrderExperience{
							ExperienceID:          "experience-id",
							AdultCount:            2,
							AdultPrice:            1000,
							JuniorHighSchoolCount: 2,
							JuniorHighSchoolPrice: 800,
							ElementarySchoolCount: 0,
							ElementarySchoolPrice: 600,
							PreschoolCount:        0,
							PreschoolPrice:        400,
							SeniorCount:           0,
							SeniorPrice:           700,
							Remarks: &types.OrderExperienceRemarks{
								Transportation: "電車",
								RequestedDate:  "20240102",
								RequestedTime:  "1830",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, NewOrders(tt.orders, tt.addresses, tt.products, tt.experiences))
		})
	}
}

func TestOrders_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders Orders
		expect []*types.Order
	}{
		{
			name: "success",
			orders: Orders{
				{
					Order: types.Order{
						ID:              "order-id",
						UserID:          "user-id",
						CoordinatorID:   "coordinator-id",
						PromotionID:     "",
						ShippingMessage: "",
						Status:          types.OrderStatusShipped,
						CreatedAt:       1640962800,
						UpdatedAt:       1640962800,
						Payment: &types.OrderPayment{
							TransactionID: "transaction-id",
							MethodType:    types.PaymentMethodTypeCreditCard,
							Status:        types.PaymentStatusPaid,
							Subtotal:      1100,
							Discount:      0,
							ShippingFee:   500,
							Total:         1760,
							OrderedAt:     1640962800,
							PaidAt:        1640962800,
							Address: &types.Address{
								Lastname:       "&.",
								Firstname:      "購入者",
								PostalCode:     "1000014",
								PrefectureCode: 13,
								City:           "千代田区",
								AddressLine1:   "永田町1-7-1",
								AddressLine2:   "",
								PhoneNumber:    "090-1234-5678",
							},
						},
						Fulfillments: []*types.OrderFulfillment{
							{
								FulfillmentID:   "fulfillment-id",
								TrackingNumber:  "",
								Status:          types.FulfillmentStatusFulfilled,
								ShippingCarrier: types.ShippingCarrierUnknown,
								ShippingType:    types.ShippingTypeNormal,
								BoxNumber:       1,
								BoxSize:         types.ShippingSize60,
								ShippedAt:       1640962800,
								Address: &types.Address{
									Lastname:       "&.",
									Firstname:      "購入者",
									PostalCode:     "1000014",
									PrefectureCode: 13,
									City:           "千代田区",
									AddressLine1:   "永田町1-7-1",
									AddressLine2:   "",
									PhoneNumber:    "090-1234-5678",
								},
							},
						},
						Refund: &types.OrderRefund{
							Total:      0,
							Type:       types.RefundTypeNone,
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
					},
				},
			},
			expect: []*types.Order{
				{
					ID:              "order-id",
					UserID:          "user-id",
					CoordinatorID:   "coordinator-id",
					PromotionID:     "",
					ShippingMessage: "",
					Status:          types.OrderStatusShipped,
					CreatedAt:       1640962800,
					UpdatedAt:       1640962800,
					Payment: &types.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    types.PaymentMethodTypeCreditCard,
						Status:        types.PaymentStatusPaid,
						Subtotal:      1100,
						Discount:      0,
						ShippingFee:   500,
						Total:         1760,
						OrderedAt:     1640962800,
						PaidAt:        1640962800,
						Address: &types.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							PostalCode:     "1000014",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "090-1234-5678",
						},
					},
					Fulfillments: []*types.OrderFulfillment{
						{
							FulfillmentID:   "fulfillment-id",
							TrackingNumber:  "",
							Status:          types.FulfillmentStatusFulfilled,
							ShippingCarrier: types.ShippingCarrierUnknown,
							ShippingType:    types.ShippingTypeNormal,
							BoxNumber:       1,
							BoxSize:         types.ShippingSize60,
							ShippedAt:       1640962800,
							Address: &types.Address{
								Lastname:       "&.",
								Firstname:      "購入者",
								PostalCode:     "1000014",
								PrefectureCode: 13,
								City:           "千代田区",
								AddressLine1:   "永田町1-7-1",
								AddressLine2:   "",
								PhoneNumber:    "090-1234-5678",
							},
						},
					},
					Refund: &types.OrderRefund{
						Total:      0,
						Type:       types.RefundTypeNone,
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
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.orders.Response())
		})
	}
}
