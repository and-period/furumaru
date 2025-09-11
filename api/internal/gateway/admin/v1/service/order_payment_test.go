package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestPaymentMethodType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		methodType entity.PaymentMethodType
		expect     PaymentMethodType
	}{
		{
			name:       "cash",
			methodType: entity.PaymentMethodTypeCash,
			expect:     PaymentMethodType(types.PaymentMethodTypeCash),
		},
		{
			name:       "credit card",
			methodType: entity.PaymentMethodTypeCreditCard,
			expect:     PaymentMethodType(types.PaymentMethodTypeCreditCard),
		},
		{
			name:       "konbini",
			methodType: entity.PaymentMethodTypeKonbini,
			expect:     PaymentMethodType(types.PaymentMethodTypeKonbini),
		},
		{
			name:       "bank transfer",
			methodType: entity.PaymentMethodTypeBankTransfer,
			expect:     PaymentMethodType(types.PaymentMethodTypeBankTransfer),
		},
		{
			name:       "paypay",
			methodType: entity.PaymentMethodTypePayPay,
			expect:     PaymentMethodType(types.PaymentMethodTypePayPay),
		},
		{
			name:       "line pay",
			methodType: entity.PaymentMethodTypeLinePay,
			expect:     PaymentMethodType(types.PaymentMethodTypeLinePay),
		},
		{
			name:       "merpay",
			methodType: entity.PaymentMethodTypeMerpay,
			expect:     PaymentMethodType(types.PaymentMethodTypeMerpay),
		},
		{
			name:       "rakuten pay",
			methodType: entity.PaymentMethodTypeRakutenPay,
			expect:     PaymentMethodType(types.PaymentMethodTypeRakutenPay),
		},
		{
			name:       "au pay",
			methodType: entity.PaymentMethodTypeAUPay,
			expect:     PaymentMethodType(types.PaymentMethodTypeAUPay),
		},
		{
			name:       "paidy",
			methodType: entity.PaymentMethodTypePaidy,
			expect:     PaymentMethodType(types.PaymentMethodTypePaidy),
		},
		{
			name:       "pay easy",
			methodType: entity.PaymentMethodTypePayEasy,
			expect:     PaymentMethodType(types.PaymentMethodTypePayEasy),
		},
		{
			name:       "none",
			methodType: entity.PaymentMethodTypeNone,
			expect:     PaymentMethodType(types.PaymentMethodTypeNone),
		},
		{
			name:       "unknown",
			methodType: entity.PaymentMethodTypeUnknown,
			expect:     PaymentMethodType(types.PaymentMethodTypeUnknown),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewPaymentMethodType(tt.methodType))
		})
	}
}

func TestPaymentMethodType_StoreEntity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		methodType PaymentMethodType
		expect     entity.PaymentMethodType
	}{
		{
			name:       "cash",
			methodType: PaymentMethodType(types.PaymentMethodTypeCash),
			expect:     entity.PaymentMethodTypeCash,
		},
		{
			name:       "credit card",
			methodType: PaymentMethodType(types.PaymentMethodTypeCreditCard),
			expect:     entity.PaymentMethodTypeCreditCard,
		},
		{
			name:       "konbini",
			methodType: PaymentMethodType(types.PaymentMethodTypeKonbini),
			expect:     entity.PaymentMethodTypeKonbini,
		},
		{
			name:       "bank transfer",
			methodType: PaymentMethodType(types.PaymentMethodTypeBankTransfer),
			expect:     entity.PaymentMethodTypeBankTransfer,
		},
		{
			name:       "paypay",
			methodType: PaymentMethodType(types.PaymentMethodTypePayPay),
			expect:     entity.PaymentMethodTypePayPay,
		},
		{
			name:       "line pay",
			methodType: PaymentMethodType(types.PaymentMethodTypeLinePay),
			expect:     entity.PaymentMethodTypeLinePay,
		},
		{
			name:       "merpay",
			methodType: PaymentMethodType(types.PaymentMethodTypeMerpay),
			expect:     entity.PaymentMethodTypeMerpay,
		},
		{
			name:       "rakuten pay",
			methodType: PaymentMethodType(types.PaymentMethodTypeRakutenPay),
			expect:     entity.PaymentMethodTypeRakutenPay,
		},
		{
			name:       "au pay",
			methodType: PaymentMethodType(types.PaymentMethodTypeAUPay),
			expect:     entity.PaymentMethodTypeAUPay,
		},
		{
			name:       "paidy",
			methodType: PaymentMethodType(types.PaymentMethodTypePaidy),
			expect:     entity.PaymentMethodTypePaidy,
		},
		{
			name:       "pay easy",
			methodType: PaymentMethodType(types.PaymentMethodTypePayEasy),
			expect:     entity.PaymentMethodTypePayEasy,
		},
		{
			name:       "none",
			methodType: PaymentMethodType(types.PaymentMethodTypeNone),
			expect:     entity.PaymentMethodTypeNone,
		},
		{
			name:       "unknown",
			methodType: PaymentMethodType(types.PaymentMethodTypeUnknown),
			expect:     entity.PaymentMethodTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.methodType.StoreEntity())
		})
	}
}

func TestPaymentMethodType_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		PaymentMethodType PaymentMethodType
		expect            types.PaymentMethodType
	}{
		{
			name:              "success",
			PaymentMethodType: PaymentMethodType(types.PaymentMethodTypeCreditCard),
			expect:            types.PaymentMethodTypeCreditCard,
		},
	}
	for _, tt := range tests {
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
			expect: PaymentStatus(types.PaymentStatusUnpaid),
		},
		{
			name:   "authorized",
			status: entity.PaymentStatusAuthorized,
			expect: PaymentStatus(types.PaymentStatusAuthorized),
		},
		{
			name:   "paid",
			status: entity.PaymentStatusCaptured,
			expect: PaymentStatus(types.PaymentStatusPaid),
		},
		{
			name:   "canceled",
			status: entity.PaymentStatusCanceled,
			expect: PaymentStatus(types.PaymentStatusCanceled),
		},
		{
			name:   "refunded",
			status: entity.PaymentStatusRefunded,
			expect: PaymentStatus(types.PaymentStatusCanceled),
		},
		{
			name:   "faield",
			status: entity.PaymentStatusFailed,
			expect: PaymentStatus(types.PaymentStatusFailed),
		},
		{
			name:   "expired",
			status: entity.PaymentStatusExpired,
			expect: PaymentStatus(types.PaymentStatusFailed),
		},
		{
			name:   "unknown",
			status: entity.PaymentStatusUnknown,
			expect: PaymentStatus(types.PaymentStatusUnknown),
		},
	}
	for _, tt := range tests {
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
		expect types.PaymentStatus
	}{
		{
			name:   "success",
			status: PaymentStatus(types.PaymentStatusPaid),
			expect: types.PaymentStatusPaid,
		},
	}
	for _, tt := range tests {
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
		address *Address
		expect  *OrderPayment
	}{
		{
			name: "success",
			payment: &entity.OrderPayment{
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
			address: &Address{
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
			expect: &OrderPayment{
				OrderPayment: types.OrderPayment{
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
				orderID: "order-id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewOrderPayment(tt.payment, tt.address))
		})
	}
}

func TestOrderPayment_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payment *OrderPayment
		expect  *types.OrderPayment
	}{
		{
			name: "success",
			payment: &OrderPayment{
				OrderPayment: types.OrderPayment{
					TransactionID: "transaction-id",
					MethodType:    types.PaymentMethodTypeCreditCard,
					Status:        types.PaymentStatusPaid,
					Subtotal:      1100,
					Discount:      0,
					ShippingFee:   500,
					Total:         1600,
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
				orderID: "order-id",
			},
			expect: &types.OrderPayment{
				TransactionID: "transaction-id",
				MethodType:    types.PaymentMethodTypeCreditCard,
				Status:        types.PaymentStatusPaid,
				Subtotal:      1100,
				Discount:      0,
				ShippingFee:   500,
				Total:         1600,
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.payment.Response())
		})
	}
}
