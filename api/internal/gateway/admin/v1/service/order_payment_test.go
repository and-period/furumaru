package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestPaymentType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		paymentType entity.PaymentType
		expect      PaymentType
	}{
		{
			name:        "cash",
			paymentType: entity.PaymentTypeCash,
			expect:      PaymentTypeCash,
		},
		{
			name:        "card",
			paymentType: entity.PaymentTypeCard,
			expect:      PaymentTypeCard,
		},
		{
			name:        "unknown",
			paymentType: entity.PaymentTypeUnknown,
			expect:      PaymentTypeUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewPaymentType(tt.paymentType))
		})
	}
}

func TestPaymentType_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		paymentType PaymentType
		expect      int32
	}{
		{
			name:        "success",
			paymentType: PaymentTypeCard,
			expect:      2,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.paymentType.Response())
		})
	}
}

func TestPaymentStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.PaymentStatus
		expect PaymentStatus
	}{
		{
			name:   "unpaid",
			status: entity.PaymentStatusInitialized,
			expect: PaymentStatusUnpaid,
		},
		{
			name:   "pending",
			status: entity.PaymentStatusPending,
			expect: PaymentStatusPending,
		},
		{
			name:   "authorized",
			status: entity.PaymentStatusAuthorized,
			expect: PaymentStatusAuthorized,
		},
		{
			name:   "paid",
			status: entity.PaymentStatusCaptured,
			expect: PaymentStatusPaid,
		},
		{
			name:   "refunded",
			status: entity.PaymentStatusCanceled,
			expect: PaymentStatusRefunded,
		},
		{
			name:   "expired",
			status: entity.PaymentStatusFailed,
			expect: PaymentStatusExpired,
		},
		{
			name:   "unknown",
			status: entity.PaymentStatusUnknown,
			expect: PaymentStatusUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewPaymentStatus(tt.status))
		})
	}
}

func TestPaymentStatus_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status PaymentStatus
		expect int32
	}{
		{
			name:   "success",
			status: PaymentStatusPaid,
			expect: 4,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.status.Response())
		})
	}
}

func TestOrderPayment(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payment *entity.OrderPayment
		status  entity.PaymentStatus
		expect  *OrderPayment
	}{
		{
			name: "success",
			payment: &entity.OrderPayment{
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
			status: entity.PaymentStatusCaptured,
			expect: &OrderPayment{
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
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewOrderPayment(tt.payment, tt.status))
		})
	}
}

func TestOrderPayment_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payment *OrderPayment
		expect  *response.OrderPayment
	}{
		{
			name: "success",
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
			expect: &response.OrderPayment{
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
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.payment.Response())
		})
	}
}
