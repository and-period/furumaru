package komoju

import (
	"errors"
	"fmt"
)

var (
	ErrTooManyRequest = errors.New("komoju: too many request")
	ErrBadGateway     = errors.New("komoju: bad gateway")
	ErrGatewayTimeout = errors.New("komoju: gateway timeout")
	ErrNotImplemented = errors.New("komoju: not implemented")
)

// KOMOJU エラーコード
// @see https://ja.doc.komoju.com/docs/errors#%E3%82%A8%E3%83%A9%E3%83%BC%E3%82%B3%E3%83%BC%E3%83%89
type ErrCode string

const (
	ErrCodeBadRequest           ErrCode = "bad_request"
	ErrCodeUnauthorized         ErrCode = "unauthorized"
	ErrCodeNotFound             ErrCode = "not_found"
	ErrCodeInternalServerError  ErrCode = "internal_server_error"
	ErrCodeForbidden            ErrCode = "forbidden"
	ErrCodeUnprocessableEntity  ErrCode = "unprocessable_entity"
	ErrCodeBadGateway           ErrCode = "bad_gateway"
	ErrCodeGatewayTimeout       ErrCode = "gateway_timeout"
	ErrCodeServiceunavailable   ErrCode = "service_unavailable"
	ErrCodeRequestFailed        ErrCode = "request_failed"
	ErrCodeInvalidPaymentType   ErrCode = "invalid_payment_type"
	ErrCodeInvalidToken         ErrCode = "invalid_token"
	ErrCodeInvalidCurrency      ErrCode = "invalid_currency"
	ErrCodeNotRefundable        ErrCode = "not_refundable"
	ErrCodeNotCapturable        ErrCode = "not_capturable"
	ErrCodeNotCancellable       ErrCode = "not_cancellable"
	ErrCodeFraudulent           ErrCode = "fraudulent"
	ErrCodeInvalidParameter     ErrCode = "invalid_parameter"
	ErrCodeMissingParameter     ErrCode = "missing_parameter"
	ErrCodeInsufficientFunds    ErrCode = "insufficient_funds"
	ErrCodeUsedNumber           ErrCode = "used_number"
	ErrCodeCardDeclined         ErrCode = "card_declined"
	ErrCodeInvalidPassword      ErrCode = "invalid_password"
	ErrCodeBadVerificationValue ErrCode = "bad_verification_value"
	ErrCodeExceedsLimit         ErrCode = "exceeds_limit"
	ErrCodeCardExpired          ErrCode = "card_expired"
	ErrCodeInvalidNumber        ErrCode = "invalid_number"
	ErrCodeInvalidAccount       ErrCode = "invalid_account"
	ErrCodeRestrictedAccount    ErrCode = "restricted_account"
	ErrCodeOtherError           ErrCode = "other_error"
	ErrCodeInvalidUserKey       ErrCode = "invalid_user_key"
	ErrCodeOtherInvalid         ErrCode = "other_invalid"
)

func NewErrCode(err error) ErrCode {
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	return ""
}

func IsSessionFailed(err error) bool {
	var e *Error
	if !errors.As(err, &e) {
		return false
	}
	// ユーザー起因のもののみ失敗と判断する
	return e.Status == 422
}

func IsRetryable(err error) bool {
	var e *Error
	if !errors.As(err, &e) {
		return false
	}
	switch e.Code {
	case ErrCodeInternalServerError,
		ErrCodeBadGateway,
		ErrCodeGatewayTimeout,
		ErrCodeServiceunavailable:
		return true
	}
	return false
}

type Error struct {
	Method  string
	Route   string
	Status  int
	Code    ErrCode
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"komoju: method=%s, route=%s, status=%d, code=%s, message=%s",
		e.Method,
		e.Route,
		e.Status,
		e.Code,
		e.Message,
	)
}
