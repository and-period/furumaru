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
			expect: &OrderRefund{
				OrderRefund: response.OrderRefund{
					Canceled: false,
					Type:     OrderRefundTypeUnknown.Response(),
					Reason:   "",
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
				},
			},
			expect: &response.OrderRefund{
				Canceled: false,
				Type:     OrderRefundTypeUnknown.Response(),
				Reason:   "",
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
