package util

import (
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorResponse(t *testing.T) {
	t.Parallel()

	const msg = "some error"
	tests := []struct {
		name         string
		err          error
		expect       *ErrorResponse
		expectStatus int
	}{
		{
			name: "internal error",
			err:  exception.ErrInternal,
			expect: &ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Internal Server Error",
				Detail:  exception.ErrInternal.Error(),
			},
			expectStatus: http.StatusInternalServerError,
		},
		{
			name: "grpc error",
			err:  status.Error(codes.Internal, msg),
			expect: &ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Internal Server Error",
				Detail:  status.Error(codes.Internal, msg).Error(),
			},
			expectStatus: http.StatusInternalServerError,
		},
		{
			name: "err is empty",
			err:  nil,
			expect: &ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "unknown error code",
				Detail:  "unknown error",
			},
			expectStatus: http.StatusInternalServerError,
		},
		{
			name: "unknown error",
			err:  assert.AnError,
			expect: &ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "unknown error code",
				Detail:  "assert.AnError general error for testing",
			},
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, code := NewErrorResponse(tt.err)
			assert.Equal(t, tt.expectStatus, code)
			assert.Equal(t, tt.expect, res)
		})
	}
}

func TestErrorResponse_GetDetail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		res    *ErrorResponse
		expect string
	}{
		{
			name: "success",
			res: &ErrorResponse{
				Detail: "some error",
			},
			expect: "some error",
		},
		{
			name:   "empty",
			res:    nil,
			expect: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.res.GetDetail())
		})
	}
}

func TestErrorResponse_InternalError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		err          error
		expect       *ErrorResponse
		expectStatus int
	}{
		{
			name: "invalid argument",
			err:  exception.ErrInvalidArgument,
			expect: &ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Bad Request",
				Detail:  exception.ErrInvalidArgument.Error(),
			},
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "unauthenticated",
			err:  exception.ErrUnauthenticated,
			expect: &ErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized",
				Detail:  exception.ErrUnauthenticated.Error(),
			},
			expectStatus: http.StatusUnauthorized,
		},
		{
			name: "not found",
			err:  exception.ErrNotFound,
			expect: &ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "Not Found",
				Detail:  exception.ErrNotFound.Error(),
			},
			expectStatus: http.StatusNotFound,
		},
		{
			name: "already exists",
			err:  exception.ErrAlreadyExists,
			expect: &ErrorResponse{
				Status:  http.StatusConflict,
				Message: "Conflict",
				Detail:  exception.ErrAlreadyExists.Error(),
			},
			expectStatus: http.StatusConflict,
		},
		{
			name: "failed precondition",
			err:  exception.ErrFailedPrecondition,
			expect: &ErrorResponse{
				Status:  http.StatusPreconditionFailed,
				Message: "Precondition Failed",
				Detail:  exception.ErrFailedPrecondition.Error(),
			},
			expectStatus: http.StatusPreconditionFailed,
		},
		{
			name: "too many requests",
			err:  exception.ErrResourceExhausted,
			expect: &ErrorResponse{
				Status:  http.StatusTooManyRequests,
				Message: "Too Many Requests",
				Detail:  exception.ErrResourceExhausted.Error(),
			},
			expectStatus: http.StatusTooManyRequests,
		},
		{
			name: "not implemented",
			err:  exception.ErrNotImplemented,
			expect: &ErrorResponse{
				Status:  http.StatusNotImplemented,
				Message: "Not Implemented",
				Detail:  exception.ErrNotImplemented.Error(),
			},
			expectStatus: http.StatusNotImplemented,
		},
		{
			name: "internal",
			err:  exception.ErrInternal,
			expect: &ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Internal Server Error",
				Detail:  exception.ErrInternal.Error(),
			},
			expectStatus: http.StatusInternalServerError,
		},
		{
			name: "not internal error",
			err:  nil,
			expect: &ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "unknown error code",
				Detail:  "unknown error",
			},
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, code := NewErrorResponse(tt.err)
			assert.Equal(t, tt.expectStatus, code)
			assert.Equal(t, tt.expect, res)
		})
	}
}

func TestErrorResponse_GRPCError(t *testing.T) {
	t.Parallel()

	const msg = "some error"
	tests := []struct {
		name         string
		err          error
		expect       *ErrorResponse
		expectStatus int
	}{
		{
			name: "cancel",
			err:  status.Error(codes.Canceled, msg),
			expect: &ErrorResponse{
				Status:  499,
				Message: "",
				Detail:  status.Error(codes.Canceled, msg).Error(),
			},
			expectStatus: 499,
		},
		{
			name: "internal",
			err:  status.Error(codes.Internal, msg),
			expect: &ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Internal Server Error",
				Detail:  status.Error(codes.Internal, msg).Error(),
			},
			expectStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid argument",
			err:  status.Error(codes.InvalidArgument, msg),
			expect: &ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Bad Request",
				Detail:  status.Error(codes.InvalidArgument, msg).Error(),
			},
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "deadline exceeded",
			err:  status.Error(codes.DeadlineExceeded, msg),
			expect: &ErrorResponse{
				Status:  http.StatusGatewayTimeout,
				Message: "Gateway Timeout",
				Detail:  status.Error(codes.DeadlineExceeded, msg).Error(),
			},
			expectStatus: http.StatusGatewayTimeout,
		},
		{
			name: "not found",
			err:  status.Error(codes.NotFound, msg),
			expect: &ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "Not Found",
				Detail:  status.Error(codes.NotFound, msg).Error(),
			},
			expectStatus: http.StatusNotFound,
		},
		{
			name: "already exists",
			err:  status.Error(codes.AlreadyExists, msg),
			expect: &ErrorResponse{
				Status:  http.StatusConflict,
				Message: "Conflict",
				Detail:  status.Error(codes.AlreadyExists, msg).Error(),
			},
			expectStatus: http.StatusConflict,
		},
		{
			name: "permission denied",
			err:  status.Error(codes.PermissionDenied, msg),
			expect: &ErrorResponse{
				Status:  http.StatusForbidden,
				Message: "Forbidden",
				Detail:  status.Error(codes.PermissionDenied, msg).Error(),
			},
			expectStatus: http.StatusForbidden,
		},
		{
			name: "failed precondition",
			err:  status.Error(codes.FailedPrecondition, msg),
			expect: &ErrorResponse{
				Status:  http.StatusPreconditionFailed,
				Message: "Precondition Failed",
				Detail:  status.Error(codes.FailedPrecondition, msg).Error(),
			},
			expectStatus: http.StatusPreconditionFailed,
		},
		{
			name: "aborted",
			err:  status.Error(codes.Aborted, msg),
			expect: &ErrorResponse{
				Status:  http.StatusConflict,
				Message: "Conflict",
				Detail:  status.Error(codes.Aborted, msg).Error(),
			},
			expectStatus: http.StatusConflict,
		},
		{
			name: "resource exhausted",
			err:  status.Error(codes.ResourceExhausted, msg),
			expect: &ErrorResponse{
				Status:  http.StatusTooManyRequests,
				Message: "Too Many Requests",
				Detail:  status.Error(codes.ResourceExhausted, msg).Error(),
			},
			expectStatus: http.StatusTooManyRequests,
		},
		{
			name: "unimplemented",
			err:  status.Error(codes.Unimplemented, msg),
			expect: &ErrorResponse{
				Status:  http.StatusNotImplemented,
				Message: "Not Implemented",
				Detail:  status.Error(codes.Unimplemented, msg).Error(),
			},
			expectStatus: http.StatusNotImplemented,
		},
		{
			name: "unavailable",
			err:  status.Error(codes.Unavailable, msg),
			expect: &ErrorResponse{
				Status:  http.StatusBadGateway,
				Message: "Bad Gateway",
				Detail:  status.Error(codes.Unavailable, msg).Error(),
			},
			expectStatus: http.StatusBadGateway,
		},
		{
			name: "unauthenticated",
			err:  status.Error(codes.Unauthenticated, msg),
			expect: &ErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized",
				Detail:  status.Error(codes.Unauthenticated, msg).Error(),
			},
			expectStatus: http.StatusUnauthorized,
		},
		{
			name: "not grpc error",
			err:  nil,
			expect: &ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "unknown error code",
				Detail:  "unknown error",
			},
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, code := NewErrorResponse(tt.err)
			assert.Equal(t, tt.expectStatus, code)
			assert.Equal(t, tt.expect, res)
		})
	}
}
