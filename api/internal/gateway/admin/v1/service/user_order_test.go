package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestUserOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		order  *entity.Order
		expect *UserOrder
	}{
		{
			name: "success",
			order: &entity.Order{
				ID:            "order-id",
				UserID:        "user-id",
				CoordinatorID: "coordinator-id",
				PromotionID:   "promotion-id",
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
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &UserOrder{
				UserOrder: response.UserOrder{
					OrderID:   "order-id",
					Status:    int32(PaymentStatusPaid),
					SubTotal:  1980,
					Total:     2530,
					OrderedAt: 1640962800,
					PaidAt:    1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewUserOrder(tt.order)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestUserOrder_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		order  *UserOrder
		expect *response.UserOrder
	}{
		{
			name: "success",
			order: &UserOrder{
				UserOrder: response.UserOrder{
					OrderID:   "order-id",
					Status:    int32(PaymentStatusPaid),
					SubTotal:  1980,
					Total:     2530,
					OrderedAt: 1640962800,
					PaidAt:    1640962800,
				},
			},
			expect: &response.UserOrder{
				OrderID:   "order-id",
				Status:    int32(PaymentStatusPaid),
				SubTotal:  1980,
				Total:     2530,
				OrderedAt: 1640962800,
				PaidAt:    1640962800,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.order.Response())
		})
	}
}

func TestUserOrders(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders entity.Orders
		expect UserOrders
	}{
		{
			name: "success",
			orders: entity.Orders{
				{
					ID:            "order-id",
					UserID:        "user-id",
					CoordinatorID: "coordinator-id",
					PromotionID:   "promotion-id",
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
						CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0), UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
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
					CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: UserOrders{
				{
					UserOrder: response.UserOrder{
						OrderID:   "order-id",
						Status:    int32(PaymentStatusPaid),
						SubTotal:  1980,
						Total:     2530,
						OrderedAt: 1640962800,
						PaidAt:    1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewUserOrders(tt.orders)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestUserOrders_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders UserOrders
		expect []*response.UserOrder
	}{
		{
			name: "success",
			orders: UserOrders{
				{
					UserOrder: response.UserOrder{
						OrderID:   "order-id",
						Status:    int32(PaymentStatusPaid),
						SubTotal:  1980,
						Total:     2530,
						OrderedAt: 1640962800,
						PaidAt:    1640962800,
					},
				},
			},
			expect: []*response.UserOrder{
				{
					OrderID:   "order-id",
					Status:    int32(PaymentStatusPaid),
					SubTotal:  1980,
					Total:     2530,
					OrderedAt: 1640962800,
					PaidAt:    1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.orders.Response())
		})
	}
}
