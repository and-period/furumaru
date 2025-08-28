package util

import (
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	StatusClientClosedRequest = 499
)

type ErrorResponse struct {
	Status  int    `json:"status"`  // ステータスコード
	Message string `json:"message"` // エラー概要
	Detail  string `json:"detail"`  // エラー詳細
}

func NewErrorResponse(err error) (*ErrorResponse, int) {
	if status, ok := internalError(err); ok {
		return newErrorResponse(status, err), status
	}
	if status, ok := grpcError(err); ok {
		return newErrorResponse(status, err), status
	}

	if err == nil {
		err = errors.New("unknown error")
	}

	res := &ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: "unknown error code",
		Detail:  err.Error(),
	}
	return res, http.StatusInternalServerError
}

func (r *ErrorResponse) GetDetail() string {
	if r == nil {
		return ""
	}
	return r.Detail
}

func newErrorResponse(status int, err error) *ErrorResponse {
	return &ErrorResponse{
		Status:  status,
		Message: http.StatusText(status),
		Detail:  err.Error(),
	}
}

func internalError(err error) (int, bool) {
	if err == nil {
		return 0, false
	}

	var s int
	switch {
	// 4xx
	case errors.Is(err, exception.ErrInvalidArgument), errors.Is(err, exception.ErrOutOfRange):
		s = http.StatusBadRequest
	case errors.Is(err, exception.ErrUnauthenticated):
		s = http.StatusUnauthorized
	case errors.Is(err, exception.ErrForbidden):
		s = http.StatusForbidden
	case errors.Is(err, exception.ErrNotFound):
		s = http.StatusNotFound
	case errors.Is(err, exception.ErrAlreadyExists):
		s = http.StatusConflict
	case errors.Is(err, exception.ErrFailedPrecondition):
		s = http.StatusPreconditionFailed
	case errors.Is(err, exception.ErrUnprocessableEntity):
		s = http.StatusUnprocessableEntity
	case errors.Is(err, exception.ErrResourceExhausted):
		s = http.StatusTooManyRequests
	case errors.Is(err, exception.ErrCanceled):
		s = StatusClientClosedRequest
	// 5xx
	case errors.Is(err, exception.ErrInternal):
		s = http.StatusInternalServerError
	case errors.Is(err, exception.ErrNotImplemented):
		s = http.StatusNotImplemented
	case errors.Is(err, exception.ErrUnavailable):
		s = http.StatusBadGateway
	case errors.Is(err, exception.ErrDeadlineExceeded):
		s = http.StatusGatewayTimeout
	default:
		return 0, false
	}

	return s, true
}

func grpcError(err error) (int, bool) {
	if err == nil {
		return 0, false
	}

	var s int
	switch status.Code(err) {
	// 4xx
	case codes.InvalidArgument, codes.OutOfRange:
		s = http.StatusBadRequest
	case codes.Unauthenticated:
		s = http.StatusUnauthorized
	case codes.PermissionDenied:
		s = http.StatusForbidden
	case codes.NotFound:
		s = http.StatusNotFound
	case codes.AlreadyExists, codes.Aborted:
		s = http.StatusConflict
	case codes.FailedPrecondition:
		s = http.StatusPreconditionFailed
	case codes.ResourceExhausted:
		s = http.StatusTooManyRequests
	case codes.Canceled:
		s = StatusClientClosedRequest
	// 5xx
	case codes.Internal, codes.DataLoss:
		s = http.StatusInternalServerError
	case codes.Unimplemented:
		s = http.StatusNotImplemented
	case codes.Unavailable:
		s = http.StatusBadGateway
	case codes.DeadlineExceeded:
		s = http.StatusGatewayTimeout
	default:
		return 0, false
	}

	return s, true
}
