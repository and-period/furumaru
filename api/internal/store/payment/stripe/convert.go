package stripe

import (
	"github.com/and-period/furumaru/api/internal/store/entity"
	lib "github.com/stripe/stripe-go/v82"
)

func convertPaymentIntentStatus(status lib.PaymentIntentStatus) entity.PaymentStatus {
	switch status {
	case lib.PaymentIntentStatusRequiresPaymentMethod,
		lib.PaymentIntentStatusRequiresConfirmation,
		lib.PaymentIntentStatusRequiresAction,
		lib.PaymentIntentStatusProcessing:
		return entity.PaymentStatusPending
	case lib.PaymentIntentStatusRequiresCapture:
		return entity.PaymentStatusAuthorized
	case lib.PaymentIntentStatusSucceeded:
		return entity.PaymentStatusCaptured
	case lib.PaymentIntentStatusCanceled:
		return entity.PaymentStatusCanceled
	default:
		return entity.PaymentStatusUnknown
	}
}
