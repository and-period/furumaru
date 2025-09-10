package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
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
			typ:      entity.OrderTypeProduct,
			expect:   OrderType(types.OrderTypeProduct),
			response: 1,
		},
		{
			name:     "experience",
			typ:      entity.OrderTypeExperience,
			expect:   OrderType(types.OrderTypeExperience),
			response: 2,
		},
		{
			name:     "unknown",
			typ:      entity.OrderTypeUnknown,
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
		name     string
		status   entity.OrderStatus
		expect   OrderStatus
		response int32
	}{
		{
			name:     "unpaid",
			status:   entity.OrderStatusUnpaid,
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
			status:   entity.OrderStatusPreparing,
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
			status:   entity.OrderStatusCompleted,
			expect:   OrderStatus(types.OrderStatusCompleted),
			response: 3,
		},
		{
			name:     "canceled",
			status:   entity.OrderStatusCanceled,
			expect:   OrderStatus(types.OrderStatusCanceled),
			response: 4,
		},
		{
			name:     "refunded",
			status:   entity.OrderStatusRefunded,
			expect:   OrderStatus(types.OrderStatusRefunded),
			response: 5,
		},
		{
			name:     "failed",
			status:   entity.OrderStatusFailed,
			expect:   OrderStatus(types.OrderStatusFailed),
			response: 6,
		},
		{
			name:     "unknown",
			status:   entity.OrderStatusUnknown,
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
		name     string
		order    *entity.Order
		products map[int64]*Product
		expect   *Order
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
				OrderMetadata: entity.OrderMetadata{
					PickupAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
					PickupLocation: "施設の入り口",
				},
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
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
			expect: &Order{
				Order: types.Order{
					ID:            "order-id",
					CoordinatorID: "coordinator-id",
					PromotionID:   "promotion-id",
					Type:          types.OrderTypeProduct,
					Status:        types.OrderStatusPreparing,
					Payment: &types.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    NewPaymentMethodType(entity.PaymentMethodTypeCreditCard).Response(),
						Status:        NewPaymentStatus(entity.PaymentStatusCaptured).Response(),
						Subtotal:      1980,
						Discount:      0,
						ShippingFee:   550,
						Total:         2530,
						OrderedAt:     1640962800,
						PaidAt:        1640962800,
					},
					Refund: &types.OrderRefund{
						Total:      0,
						Type:       NewRefundType(entity.RefundTypeNone).Response(),
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
					PickupAt:       1640962800,
					PickupLocation: "施設の入り口",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, NewOrder(tt.order, tt.products))
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
			order: &Order{Order: types.Order{
					ID:            "order-id",
					CoordinatorID: "coordinator-id",
					PromotionID:   "",
					Status:        types.OrderStatusPreparing,
					Payment: &types.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    NewPaymentMethodType(entity.PaymentMethodTypeCreditCard).Response(),
						Status:        NewPaymentStatus(entity.PaymentStatusCaptured).Response(),
						Subtotal:      1100,
						Discount:      0,
						ShippingFee:   500,
						Total:         1600,
						OrderedAt:     1640962800,
						PaidAt:        1640962800,
					},
					Refund: &types.OrderRefund{
						Total:      0,
						Type:       NewRefundType(entity.RefundTypeNone).Response(),
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
					PickupAt:       1640962800,
					PickupLocation: "施設の入り口",
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
			order: &Order{Order: types.Order{
					ID:            "order-id",
					CoordinatorID: "coordinator-id",
					PromotionID:   "promotion-id",
					Status:        types.OrderStatusPreparing,
					Payment: &types.OrderPayment{
						TransactionID: "transaction-id",
						MethodType:    NewPaymentMethodType(entity.PaymentMethodTypeCreditCard).Response(),
						Status:        NewPaymentStatus(entity.PaymentStatusCaptured).Response(),
						Subtotal:      1980,
						Discount:      0,
						ShippingFee:   550,
						Total:         2530,
						OrderedAt:     1640962800,
						PaidAt:        1640962800,
					},
					Refund: &types.OrderRefund{
						Total:      0,
						Type:       NewRefundType(entity.RefundTypeNone).Response(),
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
					PickupAt:       1640962800,
					PickupLocation: "施設の入り口",
				},
			},
			expect: &types.Order{
				ID:            "order-id",
				CoordinatorID: "coordinator-id",
				PromotionID:   "promotion-id",
				Status:        types.OrderStatusPreparing,
				Payment: &types.OrderPayment{
					TransactionID: "transaction-id",
					MethodType:    NewPaymentMethodType(entity.PaymentMethodTypeCreditCard).Response(),
					Status:        NewPaymentStatus(entity.PaymentStatusCaptured).Response(),
					Subtotal:      1980,
					Discount:      0,
					ShippingFee:   550,
					Total:         2530,
					OrderedAt:     1640962800,
					PaidAt:        1640962800,
				},
				Refund: &types.OrderRefund{
					Total:      0,
					Type:       NewRefundType(entity.RefundTypeNone).Response(),
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
				PickupAt:       1640962800,
				PickupLocation: "施設の入り口",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.order.Response())
		})
	}
}
