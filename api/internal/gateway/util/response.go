package util

import (
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	case errors.Is(err, exception.ErrInvalidArgument):
		s = http.StatusBadRequest
	case errors.Is(err, exception.ErrUnauthenticated):
		s = http.StatusUnauthorized
	case errors.Is(err, exception.ErrNotFound):
		s = http.StatusNotFound
	case errors.Is(err, exception.ErrAlreadyExists):
		s = http.StatusConflict
	case errors.Is(err, exception.ErrFailedPrecondition):
		s = http.StatusPreconditionFailed
	case errors.Is(err, exception.ErrResourceExhausted):
		s = http.StatusTooManyRequests
	case errors.Is(err, exception.ErrNotImplemented):
		s = http.StatusNotImplemented
	case errors.Is(err, exception.ErrInternal):
		s = http.StatusInternalServerError
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
	case codes.Canceled:
		s = 499 // client closed request
	case codes.Internal, codes.DataLoss:
		s = http.StatusInternalServerError
	case codes.InvalidArgument, codes.OutOfRange:
		s = http.StatusBadRequest
	case codes.DeadlineExceeded:
		s = http.StatusGatewayTimeout
	case codes.NotFound:
		s = http.StatusNotFound
	case codes.AlreadyExists:
		s = http.StatusConflict
	case codes.PermissionDenied:
		s = http.StatusForbidden
	case codes.FailedPrecondition:
		s = http.StatusPreconditionFailed
	case codes.Aborted:
		s = http.StatusConflict
	case codes.ResourceExhausted:
		s = http.StatusTooManyRequests
	case codes.Unimplemented:
		s = http.StatusNotImplemented
	case codes.Unavailable:
		s = http.StatusBadGateway
	case codes.Unauthenticated:
		s = http.StatusUnauthorized
	default:
		return 0, false
	}

	return s, true
}
