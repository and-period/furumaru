package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestPaymentMethodType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		PaymentMethodType entity.PaymentMethodType
		expect            PaymentMethodType
	}{
		{
			name:              "cash",
			PaymentMethodType: entity.PaymentMethodTypeCash,
			expect:            PaymentMethodTypeCash,
		},
		{
			name:              "card",
			PaymentMethodType: entity.PaymentMethodTypeCreditCard,
			expect:            PaymentMethodTypeCreditCard,
		},
		{
			name:              "unknown",
			PaymentMethodType: entity.PaymentMethodTypeUnknown,
			expect:            PaymentMethodTypeUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewPaymentMethodType(tt.PaymentMethodType))
		})
	}
}

func TestPaymentMethodType_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		PaymentMethodType PaymentMethodType
		expect            int32
	}{
		{
			name:              "success",
			PaymentMethodType: PaymentMethodTypeCreditCard,
			expect:            2,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.PaymentMethodType.Response())
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
			status: entity.PaymentStatusRefunded,
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
		payment *entity.Payment
		status  entity.PaymentStatus
		expect  *OrderPayment
	}{
		{
			name: "success",
			payment: &entity.Payment{
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
			status: entity.PaymentStatusCaptured,
			expect: &OrderPayment{
				OrderPayment: response.OrderPayment{
					TransactionID: "transaction-id",
					MethodType:    PaymentMethodTypeCreditCard.Response(),
					Status:        PaymentStatusPaid.Response(),
					Subtotal:      1100,
					Discount:      0,
					ShippingFee:   500,
					Tax:           160,
					Total:         1760,
					AddressID:     "address-id",
				},
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

func TestOrderPayment_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payment *OrderPayment
		address *Address
		expect  *OrderPayment
	}{
		{
			name: "success",
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
			address: &Address{
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
			expect: &OrderPayment{
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
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.payment.Fill(tt.address)
			assert.Equal(t, tt.expect, tt.payment)
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
			expect: &response.OrderPayment{
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
