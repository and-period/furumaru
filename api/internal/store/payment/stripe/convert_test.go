package stripe

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/store/entity"
	lib "github.com/stripe/stripe-go/v82"
	"github.com/stretchr/testify/assert"
)

func TestConvertPaymentIntentStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status lib.PaymentIntentStatus
		expect entity.PaymentStatus
	}{
		{
			name:   "requires_payment_method",
			status: lib.PaymentIntentStatusRequiresPaymentMethod,
			expect: entity.PaymentStatusPending,
		},
		{
			name:   "requires_confirmation",
			status: lib.PaymentIntentStatusRequiresConfirmation,
			expect: entity.PaymentStatusPending,
		},
		{
			name:   "requires_action",
			status: lib.PaymentIntentStatusRequiresAction,
			expect: entity.PaymentStatusPending,
		},
		{
			name:   "processing",
			status: lib.PaymentIntentStatusProcessing,
			expect: entity.PaymentStatusPending,
		},
		{
			name:   "requires_capture",
			status: lib.PaymentIntentStatusRequiresCapture,
			expect: entity.PaymentStatusAuthorized,
		},
		{
			name:   "succeeded",
			status: lib.PaymentIntentStatusSucceeded,
			expect: entity.PaymentStatusCaptured,
		},
		{
			name:   "canceled",
			status: lib.PaymentIntentStatusCanceled,
			expect: entity.PaymentStatusCanceled,
		},
		{
			name:   "unknown",
			status: "unknown_status",
			expect: entity.PaymentStatusUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := convertPaymentIntentStatus(tt.status)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
