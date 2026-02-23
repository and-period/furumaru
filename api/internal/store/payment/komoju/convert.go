package komoju

import (
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func convertPaymentStatus(status PaymentStatus) entity.PaymentStatus {
	switch status {
	case PaymentStatusPending:
		return entity.PaymentStatusPending
	case PaymentStatusAuthorized:
		return entity.PaymentStatusAuthorized
	case PaymentStatusCaptured:
		return entity.PaymentStatusCaptured
	case PaymentStatusCancelled:
		return entity.PaymentStatusCanceled
	case PaymentStatusRefunded:
		return entity.PaymentStatusRefunded
	case PaymentStatusExpired, PaymentStatusFailed:
		return entity.PaymentStatusFailed
	default:
		return entity.PaymentStatusUnknown
	}
}

func paymentTypesFromMethodType(methodType entity.PaymentMethodType) []PaymentType {
	switch methodType {
	case entity.PaymentMethodTypeCreditCard:
		return []PaymentType{PaymentTypeCreditCard}
	case entity.PaymentMethodTypeKonbini:
		return []PaymentType{PaymentTypeKonbini}
	case entity.PaymentMethodTypeBankTransfer:
		return []PaymentType{PaymentTypeBankTransfer}
	case entity.PaymentMethodTypePayPay:
		return []PaymentType{PaymentTypePayPay}
	case entity.PaymentMethodTypeLinePay:
		return []PaymentType{PaymentTypeLinePay}
	case entity.PaymentMethodTypeMerpay:
		return []PaymentType{PaymentTypeMerpay}
	case entity.PaymentMethodTypeRakutenPay:
		return []PaymentType{PaymentTypeRakutenPay}
	case entity.PaymentMethodTypeAUPay:
		return []PaymentType{PaymentTypeAUPay}
	case entity.PaymentMethodTypePaidy:
		return []PaymentType{PaymentTypePaidy}
	case entity.PaymentMethodTypePayEasy:
		return []PaymentType{PaymentTypePayEasy}
	default:
		return []PaymentType{}
	}
}
