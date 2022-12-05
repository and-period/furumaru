package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestOrderRefundType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		cancelType entity.CancelType
		expect     OrderRefundType
	}{
		{
			name:       "unknown",
			cancelType: entity.CancelTypeUnknown,
			expect:     OrderRefundTypeUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewOrderRefundType(tt.cancelType))
		})
	}
}

func TestOrderRefundType_IsCanceled(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		refundType OrderRefundType
		expect     bool
	}{
		{
			name:       "success non cancel",
			refundType: OrderRefundTypeUnknown,
			expect:     false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.refundType.IsCanceled())
		})
	}
}

func TestOrderRefundType_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		refundType OrderRefundType
		expect     int32
	}{
		{
			name:       "unknown",
			refundType: OrderRefundTypeUnknown,
			expect:     0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.refundType.Response())
		})
	}
}

func TestOrderRefund(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		order  *entity.Order
		expect *OrderRefund
	}{
		{
			name: "success",
			order: &entity.Order{
				ID:                "order-id",
				UserID:            "user-id",
				CoordinatorID:     "coordinator-id",
				ScheduleID:        "schedule-id",
				PromotionID:       "",
				PaymentStatus:     entity.PaymentStatusInitialized,
				FulfillmentStatus: entity.FulfillmentStatusUnfulfilled,
				CancelType:        entity.CancelTypeUnknown,
				CancelReason:      "",
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
					MethodType:    entity.PaymentMethodTypeCard,
					MethodID:      "payment-id",
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
			expect: &OrderRefund{
				OrderRefund: response.OrderRefund{
					Canceled: false,
					Type:     int32(OrderRefundTypeUnknown),
					Reason:   "",
					Total:    0,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewOrderRefund(tt.order))
		})
	}
}

func TestOrderRefund_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		refund *OrderRefund
		expect *response.OrderRefund
	}{
		{
			name: "success",
			refund: &OrderRefund{
				OrderRefund: response.OrderRefund{
					Canceled: false,
					Type:     OrderRefundTypeUnknown.Response(),
					Reason:   "",
					Total:    0,
				},
			},
			expect: &response.OrderRefund{
				Canceled: false,
				Type:     OrderRefundTypeUnknown.Response(),
				Reason:   "",
				Total:    0,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.refund.Response())
		})
	}
}
